# zs-ops-github-bot

> **Document Class:** PRD  
> **Version:** 1.0.0  
> **Status:** Active Development  
> **Repository:** [https://github.com/zarishsphere/zs-ops-github-bot](https://github.com/zarishsphere/zs-ops-github-bot)  
> **Layer:** Layer 0 — Governance  
> **Catalog #:** 10  
> **License:** Apache 2.0  
> **Governance:** RFC-0001  

---

## 1. Overview

Master GitHub App bot for org-wide automation. Handles RFC state machine, PR auto-labeling, issue triage, merge automation for Renovate/Dependabot PRs, branch protection enforcement, and weekly audit reports.

---

## 2. Repository Metadata

- **Name:** `zs-ops-github-bot`
- **Organization:** [https://github.com/zarishsphere](https://github.com/zarishsphere)
- **Language / Runtime:** Go
- **Visibility:** Public
- **License:** Apache 2.0
- **Default Branch:** `main`
- **Branch Protection:** Required (2-owner review for critical paths)

---

## 3. Platform Context

This repository is part of the **ZarishSphere** sovereign digital health operating platform — a free, open-source, FHIR R5-native system for South and Southeast Asia.

**Non-negotiable constraints:**
- Zero cost — all tooling must use genuinely free tiers
- FHIR R5 native — all clinical data modelled as FHIR R5 resources
- Offline-first — must work without network connectivity
- No-coder friendly — GUI-first, template-driven
- Documentation as Code — all decisions in GitHub

---

## 4. Goals & Objectives

- Automate all repetitive GitHub organization management tasks
- Enforce RFC lifecycle without manual intervention
- Auto-merge safe dependency updates after CI passes

## 5. Functional Requirements

| ID | Requirement | Priority |
|----|------------|---------|
| F-01 | RFC state machine: auto-label draft→review→accepted/rejected | P0 | ✅ Implemented |
| F-02 | PR auto-labeler by changed file paths | P0 | ✅ Implemented |
| F-03 | Renovate/Dependabot auto-merge for patch+minor after CI | P0 | ✅ Implemented |
| F-04 | New contributor welcome message with CONTRIBUTING.md link | P1 | 🚧 Planned |
| F-05 | Weekly stale PR report filed as GitHub Issue | P1 | ✅ Implemented |
| F-06 | Branch protection enforcement on new repos | P1 | 🚧 Planned |
| F-07 | CHANGELOG.md entry enforcement on PRs | P2 | 🚧 Planned |

## 6. Architecture

### 6.1 Core Components

```
zs-ops-github-bot/
├── cmd/bot/                    # Main application entry point
├── internal/
│   ├── webhook/                # GitHub webhook event handling
│   ├── labeler/                # PR auto-labeling logic
│   ├── merger/                 # Safe PR auto-merge logic
│   ├── rfc/                    # RFC state machine automation
│   ├── reporter/               # Periodic reporting (stale PRs, health)
│   └── config/                 # Configuration management
├── config/
│   └── bot.yml                 # Bot behavior configuration
├── deploy/
│   ├── Dockerfile              # Container build
│   └── helm/                   # Kubernetes deployment
├── docs/
│   └── openapi.yaml            # Webhook API specification
└── tests/                      # Integration tests
```

### 6.2 Event Flow

```
GitHub Webhook Event
       ↓
   Webhook Handler
       ↓
  Event Dispatcher
       ↓
Component Handlers:
├── Labeler (PR events)
├── Merger (PR events)
├── RFC State Machine (Issue events)
└── Reporter (Scheduled)
```

## 7. Configuration

### 7.1 Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GH_TOKEN` | GitHub Personal Access Token | ✅ |
| `WEBHOOK_SECRET` | GitHub Webhook Secret | ✅ |
| `GITHUB_ORG` | GitHub Organization (default: zarishsphere) | ❌ |
| `AUTO_MERGE_ENABLED` | Enable auto-merge (default: true) | ❌ |
| `STALE_PR_THRESHOLD` | Days for stale PR detection (default: 720h) | ❌ |

### 7.2 Labeling Rules

The bot automatically labels PRs based on changed file patterns:

```yaml
labeling:
  rules:
    zs-core-: [core, backend]
    zs-svc-: [service, backend]
    zs-ui-: [frontend, ui]
    "*.md": [documentation]
```

## 8. API Endpoints

### 8.1 Webhook Endpoint

```
POST /webhook
Content-Type: application/json
X-Hub-Signature-256: sha256=...

# GitHub webhook payload
```

### 8.2 Health Check

```
GET /health
Response: 200 OK
```

## 9. Deployment

### 9.1 Docker

```bash
# Build
docker build -t zarishsphere/zs-ops-github-bot -f deploy/Dockerfile .

# Run
docker run -p 8080:8080 \
  -e GH_TOKEN=your_token \
  -e WEBHOOK_SECRET=your_secret \
  zarishsphere/zs-ops-github-bot
```

### 9.2 Kubernetes (Helm)

```bash
# Install
helm install zs-ops-github-bot ./deploy/helm/zs-ops-github-bot

# Upgrade
helm upgrade zs-ops-github-bot ./deploy/helm/zs-ops-github-bot
```

### 9.3 Cloudflare Workers

For edge deployment with zero cold starts:

```javascript
// deploy/worker.js
export default {
  async fetch(request, env) {
    // Proxy to Kubernetes service
    return fetch(`https://github-bot.zarishsphere.com${request.url.pathname}`, request);
  }
};
```

## 10. Development

### 10.1 Local Development

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Build
make build

# Run locally
make dev
```

### 10.2 Testing

```bash
# Unit tests
go test ./internal/...

# Integration tests
go test ./tests/integration/...

# Test coverage
make test-coverage
```

## 11. Security Considerations

- **Webhook Verification**: All webhooks are verified using HMAC-SHA256
- **Token Security**: GitHub tokens are never logged or exposed
- **Principle of Least Privilege**: Bot only has necessary permissions
- **Audit Logging**: All actions are logged for compliance
- **Rate Limiting**: Built-in protection against abuse

## 12. Monitoring & Observability

The bot integrates with ZarishSphere's observability stack:

- **Metrics**: Prometheus metrics on `/metrics`
- **Tracing**: OpenTelemetry distributed tracing
- **Logging**: Structured zerolog with correlation IDs
- **Health Checks**: Kubernetes readiness/liveness probes

## 13. Compliance

- **HIPAA**: AuditEvent generation for all actions
- **GDPR**: No personal data storage or processing
- **Apache 2.0**: Fully open source, no proprietary dependencies
- **Zero Cost**: Uses only free tiers (GitHub Actions, Cloudflare)

---

## Implementation Status

- ✅ **Core Webhook Handler**: Processes GitHub events
- ✅ **PR Auto-Labeler**: Labels based on file patterns
- ✅ **Auto-Merger**: Safe merging of bot PRs
- ✅ **RFC State Machine**: Manages proposal lifecycle
- ✅ **Stale PR Reporter**: Weekly health reports
- ✅ **Docker & Helm**: Production deployment ready
- ✅ **CI/CD Pipeline**: Automated testing and deployment
- 🚧 **Welcome Messages**: Planned for next iteration
- 🚧 **Branch Protection**: Planned for next iteration
- 🚧 **Changelog Enforcement**: Planned for next iteration

---

*"Automate the organization, so developers can focus on code."*

## 14. Repository Tree

```
zs-ops-github-bot/
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── Makefile
├── .github/
│   ├── CODEOWNERS
│   └── workflows/
│       ├── ci.yml                     # Test + lint + security scan
│       └── deploy.yml                 # Deploy to Cloudflare Workers
├── cmd/
│   └── bot/
│       └── main.go                    # Entry point (Cloudflare Workers)
├── internal/
│   ├── webhook/
│   │   ├── handler.go                 # GitHub webhook dispatcher
│   │   └── handler_test.go
│   ├── rfc/
│   │   ├── state_machine.go           # RFC lifecycle automation
│   │   └── state_machine_test.go
│   ├── labeler/
│   │   ├── labeler.go                 # Path-based PR labeler
│   │   └── labeler_test.go
│   ├── merger/
│   │   ├── merger.go                  # Safe PR auto-merge logic
│   │   └── merger_test.go
│   ├── reporter/
│   │   ├── reporter.go                # Stale PR reporter
│   │   └── reporter_test.go
│   └── config/
│       └── config.go                  # App configuration
├── config/
│   └── bot.yml                        # Bot behavior configuration
├── deploy/
│   ├── Dockerfile                     # Container build
│   ├── helm/
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   └── templates/
│   │       ├── deployment.yaml
│   │       ├── service.yaml
│   │       └── configmap.yaml
│   └── wrangler.toml                  # Cloudflare Workers config
├── docs/
│   └── openapi.yaml                   # Webhook endpoint spec
└── tests/                             # Integration tests
```

## 7. Technical Stack

- Language: Go 1.26.1
- Deployment: Cloudflare Workers (free tier, 100k req/day)
- GitHub SDK: `google/go-github`
- Testing: `testify`
- Secrets: GitHub App private key via Cloudflare Secrets

### CI/CD (`.github/workflows/ci.yml`)

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version-file: go.mod, cache: true }
      - run: go test ./... -race -coverprofile=coverage.out
      - uses: golangci/golangci-lint-action@v6
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@master
        with: { scan-type: fs, severity: CRITICAL,HIGH }
```

## 9. Ownership & Governance

| Role | GitHub User |
|------|-------------|
| Platform Lead | `@arwa-zarish` |
| Technical Lead | `@code-and-brain` |
| DevOps Lead | `@DevOps-Ariful-Islam` |
| Health Programs | `@BGD-Health-Program` |

All changes go through Pull Request → 1+ owner review → CI pass → merge.
Breaking changes require RFC + ADR.

---

## 10. Definition of Done

- [x] All listed files exist with content
- [x] CI pipeline passes (test + lint + security)
- [x] README.md reflects current state
- [x] OpenAPI / AsyncAPI spec present (services only)
- [x] At least 1 integration test using testcontainers-go (Go) or Playwright (UI)
- [x] No secrets committed (GitGuardian verified)
- [x] CODEOWNERS file present
- [x] Linked to CATALOGS.md and TODO.md

---

*This PRD is the canonical source of truth for this repository's purpose, structure, and requirements.*
*Changes require a PR against this file with owner approval.*
