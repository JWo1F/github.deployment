package main

import (
	"errors"
	"log"
	"os"
)

const tmpl = `secret: secret
port: 3001
env:
	NODE_ENV: development
actions:
  - name: deployment
    if: "[[ {repository.name} = repo ]] && exit 0 || exit 1"
    env:
      NODE_ENV: production
    run: echo repo is {repository.name}`

func configTempate() {
	const configFile = "config.yml"

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		os.WriteFile("config.yml", []byte(tmpl), 0770)
		return
	}

	log.Fatalf("%s already exists. Remove or move it before creating new config from template\n", configFile)
}
