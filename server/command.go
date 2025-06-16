package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

type Handler struct {
	client *pluginapi.Client
}

type Command interface {
	Handle(args *model.CommandArgs, p *NtfyPlugin) (*model.CommandResponse, error)
	executeNtfyCommand(args *model.CommandArgs, p *NtfyPlugin) *model.CommandResponse
}

const ntfyCommandTrigger = "ntfy"

func NewCommandHandler(client *pluginapi.Client) Command {
	err := client.SlashCommand.Register(&model.Command{
		Trigger:          ntfyCommandTrigger,
		DisplayName:      "Ntfy",
		Description:      "Turn ntfy notifications on or off for a channel or set the topic",
		AutoComplete:     true,
		AutoCompleteDesc: "Turn ntfy notifications on or off for a channel or set the topic",
		AutoCompleteHint: "on|off|topic [topic]",
		IconURL:          "https://ntfy.sh/static/images/favicon.ico",
	})
	if err != nil {
		client.Log.Error("Failed to register slash command", "error", err)
	}
	// Return command handler
	return &Handler{
		client: client,
	}
}

func (c *Handler) Handle(args *model.CommandArgs, p *NtfyPlugin) (*model.CommandResponse, error) {
	fields := strings.Fields(args.Command)
	trigger := strings.TrimPrefix(fields[0], "/")
	//if trigger != ntfyCommandTrigger {

	if trigger != ntfyCommandTrigger {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: %s", args.Command),
		}, nil
	}
	return c.executeNtfyCommand(args, p), nil
}

func (c *Handler) executeNtfyCommand(args *model.CommandArgs, p *NtfyPlugin) *model.CommandResponse {

	var subscription *SubscriptionDetails
	subscription_preference, err := p.API.GetPreferenceForUser(args.UserId, "ntfy_subscribed", args.ChannelId)
	if err != nil {
		subscription = &SubscriptionDetails{
			Active: false,
			Topic:  "",
		}
	} else {

		if unmarshalErr := json.Unmarshal([]byte(subscription_preference.Value), &subscription); unmarshalErr != nil {
			subscription = &SubscriptionDetails{
				Active: false,
				Topic:  "",
			}
		}
	}

	fields := strings.Fields(args.Command)
	if len(fields) == 2 && (strings.ToLower(fields[1]) == "on" || strings.ToLower(fields[1]) == "off") {
		// on or off command
		if strings.ToLower(fields[1]) == "on" {
			// Turn on notifications
			subscription.Active = true
		} else {
			// Turn off notifications
			subscription.Active = false
		}
	} else if len(fields) == 2 && strings.ToLower(fields[1]) == "topic" {
		subscription.Topic = ""
	} else if len(fields) == 3 && strings.ToLower(fields[1]) == "topic" {
		// topic command
		sanitizedTopic := ""
		for _, r := range fields[2] {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
				sanitizedTopic += string(r)
			}
		}
		subscription.Topic = sanitizedTopic
	} else {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Invalid command format. Use `/ntfy [on|off|topic [topic]]`.",
		}
	}

	userid := args.UserId
	channelid := args.ChannelId

	subscriptionJSON, err_s := json.Marshal(subscription)
	if err_s != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to serialize subscription details.",
		}
	}
	preferences := model.Preference{
		UserId:   userid,
		Category: "ntfy_subscribed",
		Name:     channelid,
		Value:    string(subscriptionJSON),
	}

	err2 := p.API.UpdatePreferencesForUser(userid, []model.Preference{preferences})
	if err2 != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to update user preferences.",
		}
	}

	responseText := "Ntfy notifications have been turned "
	if subscription.Active {
		responseText += "on "
	} else {
		responseText += "off "
	}
	if subscription.Topic != "" {
		responseText += fmt.Sprintf("with topic '%s' ", subscription.Topic)
	} else {
		responseText += "with the default topic "
	}
	responseText += "for this channel."

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         responseText,
		Username:     "Ntfy Plugin",
		ChannelId:    channelid,
	}
}
