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
	if len(os.Args) < 2 {
		if os.Args[1] == "SET" {
			fmt.Println("error 1")
			return
		}
		if os.Args[1] == "GET" {
			fmt.Println("error 2")
			return
		}
	}

	command := strings.Join(os.Args[1:], " ") + "\n"
	socketPath := getSocketPath()

	// 데몬 서버의 유닉스 소켓에 연결
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Printf("❌ Failed to connect to T0r1m3m daemon (Is it running?): %v\n", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Printf("Failed to send command: %v\n", err)
		return
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("failed to read response: %v\n", err)
		return
	}
	fmt.Print(response)
}
