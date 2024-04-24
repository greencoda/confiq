package main

import (
	"log"
	"net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/greencoda/confiq"

	confiqjson "github.com/greencoda/confiq/loaders/json"
)

// ServiceSettings is a struct to test the config set.
type ServiceSettings struct {
	URL     *url.URL      `cfg:"url,default=https://default-service-url.com"`
	Timeout time.Duration `cfg:"timeout,default=45s"`
}

func main() {
	// Create a new config set with default options (struct tag name is "cfg" by default).
	configSet := confiq.New()

	// Load the JSON file into the config set.
	if err := configSet.Load(
		confiqjson.Load().FromFile("./servicesConfig.json"),
		confiq.WithPrefix(""),
	); err != nil {
		log.Fatal(err)
	}

	// Create a new ConfigStruct instance.
	var (
		birthdayService ServiceSettings
		addressService  ServiceSettings
		missingService  ServiceSettings
	)

	// Decode the birthday service settings from the config set.
	if err := configSet.Decode(&birthdayService, confiq.FromPrefix("birthday_service")); err != nil {
		log.Fatal(err)
	}

	// Decode the address service settings from the config set.
	if err := configSet.Decode(&addressService, confiq.FromPrefix("address_service")); err != nil {
		log.Fatal(err)
	}

	// Decode the parcel service settings from the config set.
	if err := configSet.Decode(&missingService, confiq.FromPrefix("missing_service")); err != nil {
		log.Fatal(err)
	}

	// Print the decoded service settings.
	spew.Dump(birthdayService, addressService, missingService)
}
