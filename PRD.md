# PRD — `zs-ops-github-bot`

> **Document Class:** PRD | **Version:** 1.0.0 | **Status:** Bootstrapping
> **Repository:** [https://github.com/zarishsphere/zs-ops-github-bot](https://github.com/zarishsphere/zs-ops-github-bot)
> **Layer:** Layer 0 — Governance | **Catalog #:** 10
> **License:** Apache 2.0 | **Governance:** RFC-0001

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
| F-01 | RFC state machine: auto-label draft→review→accepted/rejected | P0 |
| F-02 | PR auto-labeler by changed file paths | P0 |
| F-03 | Renovate/Dependabot auto-merge for patch+minor after CI | P0 |
| F-04 | New contributor welcome message with CONTRIBUTING.md link | P1 |
| F-05 | Weekly stale PR report filed as GitHub Issue | P1 |
| F-06 | Branch protection enforcement on new repos | P1 |
| F-07 | CHANGELOG.md entry enforcement on PRs | P2 |

## 6. Repository Tree

```
zs-ops-github-bot/
├── README.md
├── LICENSE
├── go.mod
├── go.sum
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
│   │   ├── auto_merge.go              # Safe PR auto-merge logic
│   │   └── auto_merge_test.go
│   ├── reporter/
│   │   ├── weekly_report.go           # Stale PR reporter
│   │   └── weekly_report_test.go
│   └── config/
│       └── config.go                  # App configuration
├── config/
│   ├── labeler.yml                    # Label rules by file path
│   └── bot.yml                        # Bot behavior configuration
├── deploy/
│   └── wrangler.toml                  # Cloudflare Workers config
└── docs/
    └── openapi.yaml                   # Webhook endpoint spec
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

- [ ] All listed files exist with content
- [ ] CI pipeline passes (test + lint + security)
- [ ] README.md reflects current state
- [ ] OpenAPI / AsyncAPI spec present (services only)
- [ ] At least 1 integration test using testcontainers-go (Go) or Playwright (UI)
- [ ] No secrets committed (GitGuardian verified)
- [ ] CODEOWNERS file present
- [ ] Linked to CATALOGS.md and TODO.md

---

*This PRD is the canonical source of truth for this repository's purpose, structure, and requirements.*
*Changes require a PR against this file with owner approval.*
