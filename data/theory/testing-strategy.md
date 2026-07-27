# Testing strategy for distributed systems

Tests exist to let you change code fast without fear. In a microservices world the
naive answer ("spin up everything and click around") is slow and flaky. The skill is
choosing the *cheapest* test that catches the *most likely* failure.

## The test pyramid (still the default)
```
        /\    E2E        few, slow, brittle — real user journeys only
       /  \   Integration  service + its real DB / queue
      /----\  Unit        many, fast, pure logic
```
Most tests are unit; a thin layer of integration; a *handful* of E2E. An inverted
pyramid (mostly E2E) is the classic mistake — slow, flaky, hard to debug.

## Where mocks lie
Unit tests mock the network, the DB, the queue. That's fast, but a mock encodes
*your assumption* of how the dependency behaves. When the real service changes its
contract, your mock still passes — green tests, broken prod. Two fixes:
- **Integration tests** against a real dependency (Testcontainers: real Postgres/Kafka in a container).
- **Contract tests** (Pact): the consumer publishes the requests it makes + responses
  it expects; the provider's CI verifies it still satisfies them. Catches breaking API
  changes *before* deploy, without running both services together.

## The layers, concretely
| Layer | Tests | Tool feel |
|---|---|---|
| Unit | pure functions, business rules, edge cases | fast, in-process |
| Integration | repository ↔ real DB, handler ↔ real queue | Testcontainers |
| Contract | service A's expectations of service B | Pact / consumer-driven |
| E2E | one critical journey across services | Playwright/Cypress, staging |
| Non-functional | load, chaos, soak | k6/Gatling, Chaos Monkey |

## Testing in production (yes, really)
You cannot reproduce prod's scale/data/traffic in staging. Complement pre-prod tests with:
- **Canary + metrics gates** (see /learn deployment-release).
- **Synthetic monitoring** — bots run the critical journey every minute.
- **Chaos engineering** — inject failure (kill nodes, add latency) to prove resilience.
- **Shadow traffic** — replay prod traffic at the new version.

## Go specifics
- Table-driven tests are idiomatic; `t.Run` for subtests.
- `httptest` for handlers; interfaces at the boundary so you can fake the DB in unit tests
  but still run Testcontainers integration tests.
- Race detector: `go test -race` — non-negotiable for concurrent code.

## Interview probes
- "Your services pass all tests but break in prod. Why?" → mocks drift from real contracts; add contract/integration tests.
- "How do you test a service that talks to 5 others?" → don't boot all 5; contract-test the boundaries + one E2E for the golden path.
- "How do you know a deploy is safe?" → canary + metric gate + synthetic checks, not just a green CI.

## Further reading
- Martin Fowler — TestPyramid: https://martinfowler.com/bliki/TestPyramid.html
- Martin Fowler — Practical Test Pyramid: https://martinfowler.com/articles/practical-test-pyramid.html
- Pact — Contract testing: https://docs.pact.io/
- Testcontainers for Go: https://golang.testcontainers.org/
- Google Testing Blog — Testing on the Toilet: https://testing.googleblog.com/
- Principles of Chaos Engineering: https://principlesofchaos.org/
