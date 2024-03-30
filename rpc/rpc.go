package rpc

import (
	"encoding/json"
	"fmt"
)

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
