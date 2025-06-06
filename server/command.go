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
		Description:      "Send notifications to ntfy.sh",
		AutoComplete:     true,
		AutoCompleteDesc: "Send a notification to ntfy.sh",
		AutoCompleteHint: "[message]",
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
	fields := strings.Fields(args.Command)
	if len(fields) != 2 || (strings.ToLower(fields[1]) != "on" && strings.ToLower(fields[1]) != "off") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Invalid command format. Use `/ntfy [on|off]`.",
		}
	}

	userid := args.UserId
	channelid := args.ChannelId
	var subscription *SubscriptionDetails
	if strings.ToLower(fields[1]) == "on" {
		subscription = &SubscriptionDetails{
			Active: true,
		}
	} else {
		subscription = &SubscriptionDetails{
			Active: false,
		}
	}
	subscriptionJSON, err := json.Marshal(subscription)
	if err != nil {
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

	responseText := fmt.Sprintf("Ntfy notifications have been turned %s for channel %s.", fields[1], args.ChannelId)
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         responseText,
		Username:     "Ntfy Plugin",
		ChannelId:    channelid,
	}
}
