# Roadmap

This file tracks implementation milestones for AI-assisted development.

Use milestone specs in `.ai/milestones/` as the source of truth for current work.

## Milestones

| Milestone | Status | Goal |
| --- | --- | --- |
| M1: Backend Rule Parser Foundation | planned | Start the Go backend with typed rule parser contracts, fixtures, and deterministic parser tests. |
| M2: Selector Matching | planned | Match persisted card selectors against fixture target cards without database dependencies. |
| M3: Relationship Preprocessing | planned | Generate selector-target relationship records from extracted effects and matched targets. |
| M4: Backend Persistence | planned | Add PostgreSQL migrations, sqlc queries, repositories, and idempotent writes for cards, selectors, effects, and relationships. |
| M5: Public Backend API | planned | Add health, card search, card detail, searcher results, and report endpoints. |

## Status Values

- `planned`: defined but not started.
- `in_progress`: current active work.
- `blocked`: waiting on a decision or dependency.
- `done`: delivered and verified.
- `deferred`: intentionally moved out of current scope.
