package project

import (
	"fmt"
	"log"
	"os"

	"github.com/cevixe/cdk/common/file"
)

func specBytes() []byte {

	dir := os.Getenv("CEVIXE_MOD_HOME")
	if dir == "" {
		log.Fatalf("CEVIXE_MOD_HOME not configured")
	}

	options := []string{
		fmt.Sprintf("%s/cevixe.yaml", dir),
		fmt.Sprintf("%s/cevixe.yml", dir),
	}

	for _, opt := range options {
		if file.Exists(opt) {
			return *file.GetBytes(opt)
		}
	}

	log.Fatalf("not found cevixe.yaml file on dir: %v", dir)
	return nil
}
