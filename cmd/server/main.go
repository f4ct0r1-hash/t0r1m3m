package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// OS별 임시 폴더에 안전하게 소켓 파일 경로 생성 (윈도우/리눅스 공통)
func getSocketPath() string {
	if runtimeOSIsWindows() {
		// 윈도우의 경우 이름드 파이프 경로 형태로 자동 매핑됨
		return `\\.\pipe\T0r1m3m-sock`
	}
	return filepath.Join(os.TempDir(), "T0r1m3m.sock")
}

// 윈도우 환경 체크 헬퍼
func runtimeOSIsWindows() bool {
	return os.PathSeparator == '\\'
}

func main() {
	socketPath := getSocketPath()

	// 만약 이전에 비정상 종료되면서 남아있던 소켓 파일이 있다면 삭제
	if !runtimeOSIsWindows() {
		if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Failed to remove old socket file: %v\n", err)
		}
	}

	// 유닉스 도메인 소켓 리스너 오픈
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Printf("Failed to start T0r1m3m daemon: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("T0r1m3m Daemon is running on [%s]\n", socketPath)

	// 안전한 종료(Ctrl+C 등)를 위한 시그널 처리
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down T0r1m3m daemon...")
		listener.Close()
		if !runtimeOSIsWindows() {
			os.Remove(socketPath)
		}
		os.Exit(0)
	}()

	// 클라이언트 접속 대기 루프
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			break
		}
		// 각 클라이언트를 고루틴으로 병렬 처리
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 클라이언트가 보낸 메시지를 줄바꿈(\n) 단위로 읽기
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Printf("📥 Received from client: %s", message)

	// 서버 응답 전송 (Echo)
	response := "ACK: " + message
	_, _ = conn.Write([]byte(response))
}
