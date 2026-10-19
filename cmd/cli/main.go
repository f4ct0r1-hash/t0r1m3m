package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
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

	// 서버로 보낼 테스트 메시지 (\n 필수)
	message := "Hello Tori, Set up connection!\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
		return
	}

	// 서버로부터 응답 읽기
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	fmt.Printf("✨ Response from server: %s", response)
}
