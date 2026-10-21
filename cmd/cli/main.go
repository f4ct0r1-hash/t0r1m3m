package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// 서버와 동일한 소켓 경로 생성 함수
func getSocketPath() string {
	if os.PathSeparator == '\\' {
		return `\\.\pipe\T0r1m3m-sock`
	}
	return filepath.Join(os.TempDir(), "T0r1m3m.sock")
}

func main() {
	socketPath := getSocketPath()

	// 데몬 서버의 유닉스 소켓에 연결
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Printf("❌ Failed to connect to T0r1m3m daemon (Is it running?): %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to t0r1m3m daemon!")
	fmt.Println("Type 'quit' or 'exit' to close the client.")
	fmt.Println("--------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Print("t0r1m3m> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// 종료 처리
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Good Bye~!")
			break
		}

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Printf("connection lost: %v\n", err)
			break
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read response: %v\n", err)
			break
		}
		fmt.Print(response)
	}
}
