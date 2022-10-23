package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
)

func execCommand(cmd string, env map[string]string) bool {
	var dir, _ = os.Getwd()

	var header = "#!/bin/bash\n\ncd " + dir + "\n"
	var contents = []byte(header + cmd)
	var hash = md5.Sum(contents)
	var filename = "deploy_" + hex.EncodeToString(hash[:]) + ".sh"
	var path = os.TempDir() + filename

	if err := os.WriteFile(path, contents, 0777); err != nil {
		return false
	}

	defer os.Remove(path)

	var command = exec.Command("bash", path)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	applyEnv(command, env)

	if err := command.Run(); err != nil {
		return false
	}

	return true
}

func applyEnv(cmd *exec.Cmd, env map[string]string) {
	var config = readConfig()

	cmd.Env = os.Environ()
	for key, value := range config.Env {
		cmd.Env = append(cmd.Env, key + "=" + value)
	}
	for key, value := range env {
		cmd.Env = append(cmd.Env, key + "=" + value)
	}
}