package main

import (
	"bufio"
	"encoding/json"
	"log"
	"lsplearning/lsp"
	"lsplearning/rpc"
	"os"
)

func main() {
	logger := getLogger("./log.txt")
	logger.Println("Started!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()

		header := rpc.GetHeaderMsg(msg)
		jsonPretty, err := rpc.PretifyJsonMsg(msg)

		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
		}

		if err != nil {
			panic("hello error in the for")
		}

		handleMessage(logger, header, jsonPretty, method, content)
	}
}

func handleMessage(logger *log.Logger, header string, jsonS string, method string, contents []byte) {
	logger.Println(header)
	logger.Printf("Received message with method: %s", method)
	logger.Println(jsonS)

	switch method {
	case "initialize":
		var request lsp.InitialeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse: %s", err)
		}

		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		// start with the reply process
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Print("Sent the reply")
	}
}

func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("not a good file")
	}

	return log.New(logFile, "[loggin] ", log.Ldate|log.Ltime|log.Lshortfile)
}
