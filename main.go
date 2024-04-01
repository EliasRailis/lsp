package main

import (
	"bufio"
	"encoding/json"
	"log"
	"lsplearning/analysis"
	"lsplearning/lsp"
	"lsplearning/rpc"
	"os"
)

func main() {
	logger := getLogger("./log.txt")
	logger.Println("Started!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()

		header := rpc.GetHeaderMsg(msg)
		jsonPretty, err := rpc.PretifyJsonMsg(msg)

		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}

		handleMessage(logger, header, jsonPretty, state, method, content)
	}
}

func handleMessage(logger *log.Logger, header string, jsonPretty string, state analysis.State,
	method string, contents []byte) {
	logger.Println(header)
	logger.Printf("Received message with method: %s", method)
	logger.Println(jsonPretty)

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
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("didOpen: %s", err)
			return
		}

		logger.Printf("Opened: %s", request.Params.TextDocument.URI)

		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("didChange: %s", err)
			return
		}

		logger.Printf("Changed: %s", request.Params.TextDocument.URI)

		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(
				request.Params.TextDocument.URI,
				change.Text)
		}
	}
}

func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("not a good file")
	}

	return log.New(logFile, "[loggin] ", log.Ldate|log.Ltime|log.Lshortfile)
}
