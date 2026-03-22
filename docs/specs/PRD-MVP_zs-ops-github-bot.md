# PRD-MVP — `zs-ops-github-bot`

> **Document:** Product Requirements (MVP) | **Version:** 1.0.0-mvp
> **Repository:** [https://github.com/zarishsphere/zs-ops-github-bot](https://github.com/zarishsphere/zs-ops-github-bot)
> **Layer:** Layer 9 — Agents | **Catalog #:** 10
> **Language:** Go 1.26.1 / GitHub Actions | **License:** Apache 2.0

---

## Executive Summary

**GitHub App bot for org-wide automation — RFC state machine, PR labeler, and merge automation.**

This document defines the **Minimum Viable Product (MVP)** scope for `zs-ops-github-bot` within the ZarishSphere sovereign digital health platform. It covers what must be built first, acceptance criteria, user stories, and the complete repository file structure.


### Platform Non-Negotiables (apply to every repository)

| Constraint | Rule |
|-----------|------|
| **Zero Cost** | All tooling, hosting, and services must use genuinely free tiers |
| **Open Source** | Apache 2.0 license; all code public |
| **FHIR R5 Native** | All clinical data modelled as FHIR R5 resources |
| **Offline-First** | Must function without network connectivity |
| **No-Coder Friendly** | GUI-first, template-driven, automatable |
| **Documentation as Code** | All decisions in GitHub via RFC/ADR |
| **Multi-tenant** | tenant_id scoping on all data operations |
| **HIPAA/GDPR** | AuditEvent on all PHI access; field-level encryption |

---

## Problem Statement

Without this bot, platform owners must manually label every RFC, manually triage every PR, and manually merge safe dependency updates — consuming hours of valuable time every week.

## MVP Goals

1. Core automation function implemented and tested
2. GitHub API integration working
3. Deployable to Cloudflare Workers or GitHub Actions
4. Error handling: all failures logged and surfaced as GitHub notifications

## Triggers

- PR opened in any org repo
- Issue created with rfc:draft label
- Renovate PR opened with passing CI
- Weekly cron on Monday 9am

## Outputs

- PR auto-labeled by changed file paths
- RFC issue labeled draft→review→accepted/rejected
- Safe Renovate PRs (patch/minor) auto-merged
- Weekly stale PR report filed as GitHub Issue

## MVP Functional Requirements

| ID | Requirement | Acceptance Criteria | Priority |
|----|------------|---------------------|---------|
| M-01 | Core automation function works end-to-end | Manual trigger produces expected output | P0 |
| M-02 | GitHub API auth via GitHub App private key | API calls succeed with App JWT | P0 |
| M-03 | Unit tests cover >80% of logic | `go test -cover` shows >80% | P1 |
| M-04 | No silent failures — all errors logged | zerolog error logs appear on failure | P0 |
| M-05 | Configuration via env vars / secrets | No hardcoded values | P0 |

## MVP Complete Repository Tree

```
zs-ops-github-bot/
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
├── CHANGELOG.md
├── .github/
│   ├── CODEOWNERS
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml                 # Deploy to Cloudflare Workers
├── cmd/
│   └── bot/
│       └── main.go
├── internal/
│   ├── webhook/
│   │   ├── handler.go
│   │   └── handler_test.go
│   ├── rfc/
│   │   ├── state_machine.go
│   │   └── state_machine_test.go
│   ├── labeler/
│   │   ├── labeler.go
│   │   └── labeler_test.go
│   └── merger/
│       ├── auto_merge.go
│       └── auto_merge_test.go
├── config/
│   ├── labeler.yml
│   └── bot.yml
└── deploy/
    └── wrangler.toml
```

---


## Owners & Governance

| Role | GitHub Handle | Responsibility |
|------|--------------|----------------|
| Platform Lead | `@arwa-zarish` | Final approval, RFC votes |
| Technical Lead | `@code-and-brain` | Architecture, Go/TS review |
| DevOps Lead | `@DevOps-Ariful-Islam` | CI/CD, infra, deployment |
| Health Programs | `@BGD-Health-Program` | Clinical content, country programs |

**PR Policy:** All changes via Pull Request. Minimum 1 owner review. CI must pass. No direct commits to `main`.


---

## MVP Acceptance Checklist

- [ ] All MVP files exist in repository with real content (not placeholders)
- [ ] CI pipeline passes on `main` branch
- [ ] No secrets, credentials, or PHI committed
- [ ] README.md reflects current state with setup instructions
- [ ] CODEOWNERS file present
- [ ] All MVP functional requirements verified manually or via automated tests
- [ ] Linked to `CATALOGS.md` and `TODO.md` in `zs-docs-platform`

---

*This document is the authoritative MVP specification for `zs-ops-github-bot`.*
*Changes require a Pull Request with at least 1 owner approval.*
