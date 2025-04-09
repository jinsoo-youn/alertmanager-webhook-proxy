# alertmanager-webhook-proxy

**Alertmanager Webhook Proxy Service**  
Prometheus Alertmanager로부터 Webhook 형태의 알림(JSON)을 수신하여, 지정된 포맷으로 변환한 후 최종 대상 서버로 HTTP POST 방식으로 전달하는 경량화된 Go webhook proxy 서버입니다.

---

## 주요 기능

- Prometheus Alertmanager의 Webhook JSON 수신
- 사용자 정의 포맷(JSON)으로 데이터 변환
- 변환된 데이터의 외부 HTTP 서버로 전송 (예: Ward, Dooray)
  - ward (클라우드시스템실 통합 모니터링 시스템)
  - dooray (NHN Dooray! 서비스)

---

## 프로젝트 구조

```
alertmanager-webhook-proxy/
├── build
├── chart
│    └── alertmanager-webhook-proxy
│        ├── charts
│        └── templates
│            └── monitoring
├── cmd
│    └── alertmanager-webhook-proxy
├── internal
│    ├── config             # 설정 파일
│    ├── handler            # HTTP 요청 처리 함수
│    ├── logging            # slog 기반 로깅
│    ├── sender             # 목적지 서버로 메시지 전송
│    ├── server             # HTTP 서버 - gin 
│    └── templating         # 템플릿 구조체 
├── pkg
│    └── models
├── templates               # JSON 템플릿
└── test                    # E2E 테스트
```

---

## ⚙️ 환경 변수

| 환경 변수명            | 설명                                           | 기본값 (예시)                                               |
|----------------------|------------------------------------------------|--------------------------------------------------------|
| `LISTEN_ADDRESS`     | HTTP 수신 주소                                 | `0.0.0.0:8080`                                         |
| `LOG_LEVEL`          | 로그 레벨 (`debug`, `info`, `warn`, `error`)   | `info`                                                 |
| `STAGE`              | 환경(stage) 값 (예: dev, beta, prod 등)        | `beta`                                                 |
| `REGION`             | 리전(region) 이름                              | `kr2`                                                  |
| `WARD_ENABLE`        | Ward 전송 여부 (`true`/`false`)                | `false`                                                |
| `WARD_EVENT_URL`     | Ward 이벤트 수신 API 주소                      | `https://ward_URL/events`                              |
| `WARD_ACTOR`         | Ward에 표기할 actor 이름                        | `buoy`                                                 |
| `DOORAY_ENABLE`      | Dooray 전송 여부 (`true`/`false`)              | `false`                                                |
| `DOORAY_WEBHOOK_URL` | Dooray Webhook 주소                             | `https://nhnent.dooray.com/services/xxxxx/yyyyy/zzzzz` |

---

## 🛠️ 빌드 방법

### 바이너리 빌드 (리눅스용)
```bash
make build
```

Docker 이미지 생성
```bash
make docker-build
```

⸻

## 실행 방법

### 1. 환경 변수로 실행
```bash
export LISTEN_ADDRESS="0.0.0.0:8080"
export LOG_LEVEL="info"
export STAGE="beta"
export REGION="kr2"

# Ward 설정
export WARD_ENABLE="true"
export WARD_EVENT_URL="https://ward_URL/event"
export WARD_ACTOR="buoy"

# Dooray 설정
export DOORAY_ENABLE="true"
export DOORAY_WEBHOOK_URL="https://nhnent.dooray.com/services/xxxxx/yyyyy/zzzzz"

make run
```

2. 플래그 기반 실행
```bash
./build/alertmanager-webhook-proxy \
  -listen="0.0.0.0:8080" \
  -log-level="info" \
  -stage="beta" \
  -region="kr2" \
  -ward-enable \
  -ward-url="https://ward_URL/event" \
  -ward-actor="buoy" \
  -dooray-enable \
  -dooray-webhook="https://nhnent.dooray.com/services/xxxxx/yyyyy/zzzzz"
```

3. Docker로 실행
```bash
make docker-run
```
