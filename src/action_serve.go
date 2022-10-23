package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type ValidateBodyError struct {
	Code    int
	Message string
}

func (err ValidateBodyError) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Message)
}

var queue = make(chan []byte, 10)
var tmplRegex = regexp.MustCompile(`{{?[^{}]+}}?`)

func actionServe() {
	var config = readConfig()
	var mux = http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		log.Println("A new connection requested")

		var body, err = validateRequest(&config, req)
		if err != nil {
			var code int

			if validateError, ok := err.(ValidateBodyError); ok == true {
				code = validateError.Code
			} else {
				code = http.StatusInternalServerError
			}

			writer.WriteHeader(code)
			writer.Write([]byte(err.Error()))
			return
		}

		writer.WriteHeader(http.StatusOK)
		queue <- body
		log.Println("Request added to queue. Current size: " + strconv.Itoa(len(queue)))
	})

	go depoyWorker()
	log.Println("Server started at port", config.Port)
	http.ListenAndServe("0.0.0.0:"+config.Port, mux)
}

func checkSignature(config *Config, original string, payload []byte) bool {
	var digest = hmac.New(sha256.New, []byte(config.Secret))
	digest.Write(payload)
	var signature = "sha256=" + hex.EncodeToString(digest.Sum(nil))
	return original == signature
}

func validateRequest(config *Config, req *http.Request) ([]byte, error) {
	var body, _ = io.ReadAll(req.Body)
	var signature = req.Header.Get("X-Hub-Signature-256")

	if !gjson.ValidBytes(body) {
		return nil, ValidateBodyError{http.StatusBadRequest, "Invalid request body accepted. Close connection"}
	}

	if !checkSignature(config, signature, body) {
		return nil, ValidateBodyError{http.StatusForbidden, "Invalid signature accepted. Close connection"}
	}

	return body, nil
}

func depoyWorker() {
	for body := range queue {
		log.Println("Starting working on task")

		var config = readConfig()
		var replacer = func(s string) string {
			return gjson.GetBytes(body, s[1:len(s)-1]).String()
		}

		for _, action := range config.Actions {
			var condition = tmplRegex.ReplaceAllStringFunc(action.If, replacer)

			if condition != "" {
				if !execCommand(condition, action.Env) {
					continue
				}
			}

			var run = tmplRegex.ReplaceAllStringFunc(action.Run, replacer)
			execCommand(run, action.Env)
		}
	}
}
