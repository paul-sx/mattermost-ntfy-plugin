package main

import (
	"encoding/json"
	"sync"

	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

type NtfyPlugin struct {
	plugin.MattermostPlugin
	commandHandler    Command
	client            *pluginapi.Client
	configurationLock sync.RWMutex
	configuration     *Configuration
}

func (p *NtfyPlugin) OnActivate() error {
	p.client = pluginapi.NewClient(p.API, p.Driver)
	p.commandHandler = NewCommandHandler(p.client)
	return nil
}

func (p *NtfyPlugin) OnDeactivate() error {
	return nil
}

type SubscriptionDetails struct {
	Active bool `json:"active"`
}

// MessageHasBeenPosted is called when a message has been posted in a channel.
func (p *NtfyPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	configuration := p.getConfiguration()
	if configuration == nil || !configuration.Active {
		return
	}

	// Get the users who are subscribed to the channel
	subscribers, err := p.API.GetUsersInChannel(post.ChannelId, "username", 0, 100)
	if err != nil {
		p.API.LogError("Failed to get channel subscribers", "error", err.Error())
		return
	}

	for _, user := range subscribers {
		if user.Id == post.UserId {
			// Skip the user who posted the message
			continue
		}
		pref, err := p.API.GetPreferenceForUser(user.Id, "ntfy_subscribed", post.ChannelId)
		if err != nil {
			p.API.LogError("Failed to get user preference", "user_id", user.Id, "error", err.Error())
			continue
		}
		var details SubscriptionDetails
		if err := json.Unmarshal([]byte(pref.Value), &details); err != nil {
			p.API.LogError("Failed to unmarshal user preference", "user_id", user.Id, "error", err.Error())
			continue
		}
		if details.Active {
			notification := post.Message
			if len(notification) > 50 {
				notification = notification[:50]
			}
			payload := map[string]string{
				"title":   "New Mattermost Message",
				"message": notification,
			}
			payloadBytes, _ := json.Marshal(payload)
			// Example topic: "mattermost-" + post.ChannelId + "-" + userID
			topic := configuration.Topic
			url := configuration.ServerURL + "/" + topic

			// Replace these with your actual username and password variables
			username := configuration.Username
			password := configuration.Password
			auth := username + ":" + password

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
			if err != nil {
				p.API.LogError("Failed to create ntfy.sh request", "error", err.Error())
				return
			}
			req.Header.Set("Content-Type", "application/json")

			encoded := base64.StdEncoding.EncodeToString([]byte(auth))
			req.Header.Set("Authorization", "Basic "+encoded)

			client := &http.Client{}
			_, err2 := client.Do(req)
			if err2 != nil {
				p.API.LogError("Failed to send to ntfy", "error", err2.Error())
				return
			}
		}
	}

}
func (p *NtfyPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	response, err := p.commandHandler.Handle(args, p)
	if err != nil {
		return nil, model.NewAppError("ExecuteCommand", "plugin.command.execute_command.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return response, nil
}
