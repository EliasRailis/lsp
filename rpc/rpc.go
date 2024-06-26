package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

func EncodeMessage(msg any) string {
	// serialize the msg value to json
	content, err := json.Marshal(msg)

	// panic if there's any errors
	if err != nil {
		panic(err)
	}

	// print out the content length by using the len() function
	// and also return the content in json format
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return "", nil, errors.New("Did not find separator")
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))

	if err != nil {
		return "", nil, err
	}

	var baseMessage BaseMessage

	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))

	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}

func GetHeaderMsg(msg []byte) string {
	header, _, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return ""
	}

	return string(header)
}

func PretifyJsonMsg(msg []byte) (string, error) {
	_, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return "", nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return "", err
	}

	prettyJson, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return "", err
	}

	return string(prettyJson), nil
}
