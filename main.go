package main

import (
	"bufio"
	"log"
	"lsplearning/rpc"
	"os"
)

func main() {
	logger := getLogger("./log.txt")
	logger.Println("Started!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(msg)
	}
}

func handleMessage(_ any) {}

func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("not a good file")
	}

	return log.New(logFile, "[loggin] ", log.Ldate|log.Ltime|log.Lshortfile)
}
