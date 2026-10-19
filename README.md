# t0r1m3m
A custom in-memory database solution written in "Go"

---

## [ 1. 프로젝트 진행 배경 ]
> **취미 및 공부 목적**
> * 개인적인 관심사를 바탕으로 깊이 있는 학습과 프로그래밍 역량을 기르기 위해 진행한 프로젝트입니다.
> * 핵심 시스템 구조와 내부 동작 원리를 직접 구현하며 이해를 높이고자 하였습니다.

---

## [ 2. 프로젝트 사용 언어 ]
이 프로젝트에서 사용된 프로그래밍 언어입니다.

* **Language:** ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)

---

## [ 3. 사용된 기술 ]
프로젝트 구현에 핵심적으로 사용된 기술 및 아키텍처입니다.

* **Core Features & Architecture:**
  * **IPC (Inter-Process Communication):** 프로세스 간 통신 구현
  * **메모리 관리 (Memory Management):** 효율적인 자원 할당 및 관리
  * **프로토콜 처리 (Protocol Handling):** 커스텀 데이터 통신 규격 및 패킷 처리
  * **CLI 레이어 (Command-Line Interface):** 사용자 상호작용을 위한 터미널 인터페이스 구현

---

## [ 4. 파일 설명 ]
이 프로젝트에 대한 경로 별 설명입니다.

```text

├── cmd/
│   ├── cli
│        └── main.go     # CLI 레이어 진입점      
│   └── server
│        └── main.go     # (리눅스)데몬 / (윈도우)서비스 진입점
├── internal/             # 내부 핵심 로직(차차 추가될 예정)
└── README.md
```

---

## [ 5. 진행 로그 ]
이 프로젝트에 대한 진행 로그입니다.

### [Day 1] 로컬 데몬과 CLI 간의 통신 뼈대 구축 (2026-10-19 12:38)
- **Goal:** 네트워크 포트 오픈 방식 대신 OS 커널을 직접 활용하는 **Unix Domain Socket(UDS)** 기반 통신 파이프라인 구현
- **What I Did:**
  - `cmd/server`: 클라이언트의 접속을 대기하고 메시지를 에코(Echo)하는 데몬 서버 뼈대 작성
  - `cmd/cli`: 소켓 파일에 접속해 명령어를 전달하고 응답을 받는 CLI 클라이언트 제작
- **💡 Troubleshooting & Learnings:**
  - **OS 호환성 이슈:** 리눅스는 `/tmp/t0r1m3m.sock` 경로를 그대로 쓰면 되지만, 윈도우 환경을 고려해 `os.TempDir()` 및 파이프 경로 분기 처리를 적용했습니다.
  - **소켓 파일 잔여 문제:** 데몬이 비정상 종료되었을 때 이전에 쓰던 소켓 파일이 남아있어 `bind: address already in use` 에러가 발생하는 것을 확인, 서버 구동 시 기존 소켓 파일을 안전하게 삭제하는 로직을 추가했습니다.