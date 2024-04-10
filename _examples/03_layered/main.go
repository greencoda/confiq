package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/greencoda/confiq"
)

// DBSettings is a struct which holds the database settings.
type DBSettings struct {
	DatabaseDNS            string `cfg:"dns"`
	DatabaseMaxConnections int    `cfg:"settings.maxConnections"`
	DatabaseRunMigrations  bool   `cfg:"settings.runMigrations"`
}

// ConfigStruct is a struct to test the config set.
type ConfigStruct struct {
	DBSettings DBSettings `cfg:"dbSettings"`
}

func main() {
	// Create a new config set with default options (struct tag name is "cfg" by default).
	configSet := confiq.New()

	// Load the base DB settings from a JSON file into the config set, with a prefix of "dbSettings".
	if err := configSet.LoadJSONFromFile("./dbSettingsBase.json", confiq.WithPrefix("dbSettings")); err != nil {
		log.Fatal(err)
	}

	log.Println("Please celect db profile. [prod, dev]")

	var profileChoice string

	if _, err := fmt.Scanln(&profileChoice); err != nil {
		log.Fatalf("Failed to scan profile choice: %s", err)
	}

	switch profileChoice {
	case "prod":
		// Load the production DB settings from a JSON file into the config set, with a prefix of "dbSettings".
		loadDBConfig(configSet, "./dbSettingsProduction.json")

	case "dev":
		// Load the development DB settings from a JSON file into the config set, with a prefix of "dbSettings".
		loadDBConfig(configSet, "./dbSettingsDevelopment.json")

	default:
		log.Fatalf("Invalid profile choice: %s", profileChoice)
	}

	// Create a new ConfigStruct instance.
	var cfg ConfigStruct

	// Decode the config set into the ConfigStruct instance.
	if err := configSet.StrictDecode(&cfg); err != nil {
		log.Fatal(err)
	}

	// Print the ConfigStruct instance.
	spew.Dump(cfg)
}

func loadDBConfig(configSet *confiq.ConfigSet, configFile string) {
	if err := configSet.LoadYAMLFromFile(configFile, confiq.WithPrefix("dbSettings")); err != nil {
		log.Fatal(err)
	}
}
