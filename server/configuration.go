package main

import (
	"reflect"

	"github.com/pkg/errors"
)

type Configuration struct {
	ServerURL string `json:"server_url"`
	Topic     string `json:"topic"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Active    bool   `json:"active"`
}

func (c *Configuration) Clone() *Configuration {
	clone := *c
	return &clone
}

func (p *NtfyPlugin) getConfiguration() *Configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()
	if p.configuration == nil {
		p.configuration = &Configuration{
			ServerURL: "https://ntfy.sh",
			Topic:     "",
			Username:  "",
			Password:  "",
			Active:    false,
		}
	}
	return p.configuration
}

func (p *NtfyPlugin) setConfiguration(newConfig *Configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if newConfig != nil && p.configuration == newConfig {
		if reflect.ValueOf(p.configuration).NumField() == 0 {
			return // No changes, do nothing
		}
		panic("setConfiguration called with the same configuration instance")
	}
	p.configuration = newConfig
}

func (p *NtfyPlugin) OnConfigurationChange() error {
	var configuration = new(Configuration)
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}
	p.setConfiguration(configuration)
	return nil
}
