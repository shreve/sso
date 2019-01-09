package main

import "encoding/json"

type Message struct {
	Data string `json:"data"`
}

var errJSON = func(msg string) []byte {
	return []byte("{\"error\": \"" + msg + "\"}")
}

func Render(obj interface{}, err error) []byte {
	if err != nil {
		return errJSON(err.Error())
	}
	payload, err := json.Marshal(obj)
	if err != nil {
		return errJSON(err.Error())
	}
	return payload
}
