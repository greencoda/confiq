package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/greencoda/confiq"
)

// APISettings is a struct which holds the API settings.
type APISettings struct {
	URL     string        `cfg:"url"`
	Timeout time.Duration `cfg:"timeout"`
}

// DBSettings is a struct which holds the database settings.
type DBSettings struct {
	DatabaseDNS            string `cfg:"dns"`
	DatabaseMaxConnections int    `cfg:"settings.maxConnections"`
	DatabaseRunMigrations  bool   `cfg:"settings.runMigrations"`
}

// OSSettings is a struct which holds the database settings.
type OSSettings struct {
	OSArch string `cfg:"OS_ARCH"`
}

// UserSettings is a struct which holds the cache settings.
type UserSettings struct {
	Username         string    `cfg:"username"`
	Password         string    `cfg:"password"`
	RegistrationDate time.Time `cfg:"registrationDate"`
}

// ConfigStruct is a struct to test the config set.
type ConfigStruct struct {
	APISettings  APISettings  `cfg:"apiSettings"`
	DBSettings   DBSettings   `cfg:"dbSettings"`
	OSSettings   OSSettings   `cfg:"osSettings"`
	UserSettings UserSettings `cfg:"userSettings"`
}

func main() {
	// Create a new config set with default options (struct tag name is "cfg" by default).
	configSet := confiq.New()

	// Load the API settings from a TOML file into the config set, with a prefix of "apiSettings".
	if err := configSet.LoadTOMLFromFile("./apiSettings.toml", confiq.WithPrefix("apiSettings")); err != nil {
		log.Fatal(err)
	}

	// Load the DB settings from a JSON file into the config set.
	if err := configSet.LoadJSONFromFile("./dbSettings.json", confiq.WithPrefix("dbSettings")); err != nil {
		log.Fatal(err)
	}

	// Load the DB settings from an Env file into the config set.
	if err := configSet.LoadEnvFromFile("./osSettings.env", confiq.WithPrefix("osSettings")); err != nil {
		log.Fatal(err)
	}

	// Load the DB settings from a YAML file into the config set.
	if err := configSet.LoadYAMLFromFile("./userSettings.yaml", confiq.WithPrefix("userSettings")); err != nil {
		log.Fatal(err)
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
