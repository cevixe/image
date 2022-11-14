package project

import (
	"log"

	"gopkg.in/yaml.v3"
)

func unmarshal(buffer []byte, file interface{}) {

	err := yaml.Unmarshal(buffer, file)
	if err != nil {
		log.Fatalf("invalid configuration file format: %v", err)
	}
}
