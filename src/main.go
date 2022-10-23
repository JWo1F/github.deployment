package main

import (
	"log"
	"os"
)

func main() {
	var action, err = getAction()

	if err != nil {
		log.Fatalln(err.Error())
	}

	switch action {
	case "serve":
		actionServe()
	case "init":
		configTempate()
	}
}

func getAction() (string, error) {
	var args = os.Args

	if len(args) == 2 {
		return args[1], nil
	}

	return "serve", nil
}
