package main

import (
	"github.com/mattermost/mattermost/server/public/plugin"
)

type NtfyPlugin struct {
	plugin.MattermostPlugin
}

func (p *NtfyPlugin) OnActivate() error {
	return nil
}

func (p *NtfyPlugin) OnDeactivate() error {
	return nil
}
