# fiber-skeleton 반복 개선 워크플로우

## 작업 사이클

```text
1. 코드 리뷰  →  2. 이슈 리스트업  →  3. 수정  →  4. 커밋  →  5. CLAUDE.md 업데이트
     ↑                                                              │
     └──────────────────────────────────────────────────────────────┘
```

## 종료 조건

코드 리뷰 결과 Medium 이상의 이슈가 발견되지 않으면 해당 사이클을 종료한다.
Low 항목만 남은 경우 모아서 한 번에 처리하거나, 무시할 수 있다.

---

## 1단계: 코드 리뷰

전체 코드베이스를 읽고 아래 3가지 관점에서 이슈를 찾는다.

### 탐지 분류

#### 1. 개선/오류 수정

현재 동작하고 있는 코드가 오동작할 여지가 있거나 필요한 처리가 누락된 경우.

- 동시성 버그: 포인터 복사 없이 반환, 외부 포인터를 map에 직접 저장, 락 범위
- 로직 버그: 조건문 누락, 정수 오버플로, 에러 전파 누락, timeout/cancel 처리
- 보안: 시크릿 검증, 입력 바운드, injection, OWASP Top 10
- 설정 안정성: 파싱 실패 시 경고/기본값, 환경별 검증 누락
- 방어적 프로그래밍: nil 체크, mock 가드, 엣지 케이스
- 테스트 누락: 응답 바디 미검증, 누락된 케이스, 경계값

#### 2. 리팩토링

코드가 동작하지만 불필요하게 복잡하거나 가독성/유지보수성을 해치는 경우.

- 중복 코드: 같은 패턴이 반복되는 경우 → 헬퍼 추출
- 매직 넘버/문자열: 상수로 정의되지 않은 리터럴
- Go 컨벤션 위반: 네이밍, 에러 처리 패턴, 패키지 구조
- 불필요한 의존성: stdlib으로 대체 가능한 외부 패키지
- API 비일관성: 응답 구조, 에러코드 의미, HTTP 메서드/상태 불일치

#### 3. 기능 추가

현재 프로젝트에 추가되면 유용한 기능. 오버 엔지니어링이 되지 않도록 주의.

- 판단 기준: 현재 스켈레톤 수준에서 실제로 필요한가? 추가 시 복잡도 대비 가치가 있는가?
- 제안 예시: graceful shutdown 개선, 구조화된 로깅 필드 추가, 헬스체크 의존성 확인 등

### 오탐 필터링 기준

- Go 버전 특성 (예: Go 1.22+ 루프 변수 캡처, nil 슬라이스 JSON 직렬화)
- 프레임워크/라이브러리가 보장하는 동작
- 의도적 설계 결정

## 2단계: 이슈 리스트업

발견된 이슈를 아래 형식으로 정리한다.

```markdown
| # | 파일:라인 | 이슈 설명 | 분류 | 심각도 |
|---|-----------|-----------|------|--------|
| 1 | src/repository/example.go:45 | FindByID가 포인터를 직접 반환 | 개선/오류 수정 | High |
| 2 | src/handler/example.go:30 | bindJSON/bindQuery 중복 패턴 | 리팩토링 | Medium |
| 3 | - | graceful shutdown 시 DB 연결 정리 누락 | 기능 추가 | Low |
```

심각도 기준:

- **Critical**: 데이터 손실, 보안 취약점
- **High**: 동시성 버그, 로직 오류
- **Medium**: API 비일관성, 테스트 누락, 주요 리팩토링
- **Low**: 유지보수성, 코드 스타일, 선택적 기능 추가

## 3단계: 수정

- 항목별로 파일 수정 → `go build ./...` → `go test ./...`
- 각 이슈를 **개별 커밋**으로 분리

## 4단계: 커밋

커밋 메시지 컨벤션:

```
<type>: <설명>

type: feat | fix | refactor | test | chore | docs
```

## 5단계: CLAUDE.md 업데이트

변경된 패턴, 구조, 설정을 `CLAUDE.md`에 반영한다.

## 6단계: 검증

```bash
go build ./...        # 컴파일
go test ./...         # 테스트
golangci-lint run     # 린트
```

---

## 완료된 라운드 이력

### Round 1 — 초기 기능 구현 (18개)

godotenv 로드, JWT 검증, DTO 정리, 검증 에러 구조화, 로거 필드 추가,
Rate Limiter, CORS, BodyLimit, testutil 추출, 테스트 타임아웃,
통합 테스트, 페이지네이션, Health 분리, Prometheus, Makefile, Docker, air

### Round 2 — 1차 코드 리뷰 (11개)

Mock nil 가드, Repository 정렬/값 복사, Metrics 카디널리티,
PaginatedResp 래핑, 404 catch-all, bindJSON/bindQuery 추출,
Logger 2xx Debug, CORS production 경고, CreateExample 전체 반환, 통합 테스트 추가

### Round 3 — 2차 코드 리뷰 (6개)

FindByID 값 복사, Create/Update 값 복사, Timeout ctx.Err() 독립 체크,
CORS 콤마 split, NotFound 상수, REQUEST_TIMEOUT 파싱 경고

### Round 4 — 3차 코드 리뷰 (9개)

Page/Content max 제약, bindQuery ErrBadRequest, RequestID 순서,
테스트 바디 검증, ErrTooManyRequests 상수, Route() nil 체크,
JWT 32바이트 최소, CORS MaxAge

### Round 5 — 문서/인프라 정리

CLAUDE.md 전면 업데이트, golangci-lint v2 마이그레이션

### Round 6 — 4차 코드 리뷰 (2개)

CORS origins 빈 문자열 필터링, List 핸들러 유닛 테스트 추가.
Low 이하 이슈만 잔존 (go-errors→stdlib 전환) — 종료 조건 근접.

### Round 7 — 5차 코드 리뷰 (3개)

repository 미사용 ctx 파라미터 일관성(`ctx`→`_`), `"requestid"` 상수 추출, max pagination 경계값 테스트 추가.
오탐 2건 필터링 (동일 패키지 상수 공유, Fatal 프로세스 종료). Low 이슈만 잔존 — 종료 권고.

### Round 8 — 6차 코드 리뷰 (2개) — **종료**

`3600`(CORS MaxAge) 및 `"Bearer "` 매직 리터럴 상수화.
Medium 이상 이슈 없음 → **종료 조건 도달. 사이클 종료.**
