package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func readMessage(reader *bufio.Reader) ([]byte, error) {
	// Read headers
	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // Empty line marks end of headers
		}
		if strings.HasPrefix(line, "Content-Length:") {
			parts := strings.Split(line, ":")
			contentLength, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
		}
	}

	// Read body
	body := make([]byte, contentLength)
	_, err := io.ReadFull(reader, body)
	return body, err
}

func main() {
	cmd := exec.Command(".\\lsp-server.exe")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()

	reader := bufio.NewReader(stdout)

	// Send initialize request
	initReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"processId":    os.Getpid(),
			"rootUri":      "file:///test",
			"capabilities": map[string]interface{}{},
		},
	}

	body, _ := json.Marshal(initReq)
	msg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body)
	stdin.Write([]byte(msg))

	// Read response
	response, _ := readMessage(reader)
	fmt.Printf("Response: %s\n", response)

	// Send initialized notification
	initializedNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialized",
		"params":  map[string]interface{}{},
	}
	body2, _ := json.Marshal(initializedNotif)
	msg2 := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body2), body2)
	stdin.Write([]byte(msg2))

	cmd.Wait()
}
