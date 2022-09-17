package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
	cfg "press-test/config"
)

type Generic interface {
	PV | Press | Resource | Backfill | WRK | interface{}
}

type generic[gen Generic] struct{}

func GenericFactory() *generic[interface{}] {
	return &generic[interface{}]{}
}

func (g *generic[gen]) SendWechat(message string, xargs ...interface{}) error {
	data := map[string]interface{}{}
	mentioned_list := []string{"@all"}
	text := make(map[string]interface{})
	if len(xargs) != 0 {
		text["content"] = message
	} else {
		text["content"] = message
	}
	text["mentioned_list"] = mentioned_list
	data["msgtype"] = "text"
	data["text"] = text
	body, _ := json.Marshal(data)
	response, httpError := http.Post(cfg.WEBHOOK_URL, "text", bytes.NewReader(body))
	if httpError != nil {
		return httpError
	}
	if response.StatusCode != 200 {
		return errors.New("send wechat failed: " + message)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	return nil
}

func (g *generic[gen]) Command(cmd string) (string, error) {
	command = exec.Command("bash", "-c", cmd)
	output, err = command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
