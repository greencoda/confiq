package main

import (
	"errors"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/greencoda/confiq"
)

var errValueIsNotString = errors.New("value is not a string")

// Plumbus is a custom type which implements the decoder interface.
type Plumbus struct {
	Schleem string
}

func (p *Plumbus) Decode(value any) error {
	if stringValue, ok := value.(string); !ok {
		return errValueIsNotString
	} else {
		decryptor := func(input, key string) string {
			const c1, c2 = 32, 94

			var result string

			for i := range len(input) {
				char := input[i]
				shift := key[i%len(key)]

				decryptedChar := byte(((int(char)-c1)-(int(shift)-c1)+c2)%c2 + c1)

				result += string(decryptedChar)
			}

			return result
		}

		p.Schleem = decryptor(stringValue, "c-137")

		return nil
	}
}

// ConfigStruct is a struct to test the config set.
type ConfigStruct struct {
	Plumbus Plumbus `cfg:"plumbus"`
}

func main() {
	// Create a new config set with default options (struct tag name is "cfg" by default).
	configSet := confiq.New()

	// Load the DB settings from a JSON file into the config set.
	if err := configSet.LoadJSONFromFile("./plumbus.json"); err != nil {
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
