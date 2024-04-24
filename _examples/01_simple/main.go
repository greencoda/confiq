package main

import (
	"encoding/json"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/greencoda/confiq"

	confiqjson "github.com/greencoda/confiq/loaders/json"
)

// ConfigSettingAPIGroup is a struct to test the config set.
type ConfigSettingAPIGroup struct {
	URL     *url.URL      `cfg:"url"`
	Timeout time.Duration `cfg:"timeout"`
}

// ConfigStruct is a struct to test the config set.
type ConfigStruct struct {
	APITokens              []string                `cfg:"apiTokens"`
	OutputAPIs             []ConfigSettingAPIGroup `cfg:"outputAPIs"`
	SettingName            string                  `cfg:"settingName"`
	SettingCfg             json.RawMessage         `cfg:"settingCfg"`
	SettingIP              net.IP                  `cfg:"settingIP"`
	SettingPointer         *string                 `cfg:"settingPointer"`
	SettingCount           int                     `cfg:"settingCount"`
	SettingEnabled         bool                    `cfg:"settingEnabled,required"`
	SettingExpiration      time.Time               `cfg:"settingExpiration"`
	SettingAPIGroup        ConfigSettingAPIGroup   `cfg:"settingAPIGroup"`
	DatabaseDNS            string                  `cfg:"settingDBGroup.dns"`
	DatabaseMaxConnections int                     `cfg:"settingDBGroup.maxConnections"`
	DatabaseRunMigrations  bool                    `cfg:"settingDBGroup.runMigrations"`
}

func main() {
	// Create a new config set with default options (struct tag name is "cfg" by default).
	configSet := confiq.New()

	// Load the JSON file into the config set.
	if err := configSet.Load(
		confiqjson.Load().FromFile("./config.json"),
		confiq.WithPrefix(""),
	); err != nil {
		log.Fatal(err)
	}

	// Create a new ConfigStruct instance.
	var cfg ConfigStruct

	// Decode the config set into the ConfigStruct instance.
	if err := configSet.Decode(&cfg, confiq.AsStrict()); err != nil {
		log.Fatal(err)
	}

	// Print the ConfigStruct instance.
	spew.Dump(cfg)
}
