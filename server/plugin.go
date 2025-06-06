package main

import (
	"encoding/json"
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type NtfyPlugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex
	configuration     *Configuration
}

func (p *NtfyPlugin) OnActivate() error {
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

		}

	}
}
