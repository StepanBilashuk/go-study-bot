# Deployment & release strategies

You can have perfect code and still cause an outage at deploy time. This topic is
about shipping change *safely* — decoupling "deployed" from "released", and always
having a fast way back.

## The core idea: separate deploy from release
- **Deploy** = new code is running on servers.
- **Release** = new code is serving user traffic.
Feature flags let you deploy dark and release later (even to 1% of users first).

## Strategies
| Strategy | How | Rollback | Cost |
|---|---|---|---|
| **Recreate** | stop old, start new | redeploy old | downtime |
| **Rolling** | replace instances N at a time | roll back batch by batch | slow, mixed versions live |
| **Blue-green** | full parallel env, flip the LB | flip back instantly | 2× infra during switch |
| **Canary** | send 1%→5%→50%→100%, watch metrics | stop + shift back | needs good metrics |
| **Shadow / mirror** | copy prod traffic to new version, discard responses | n/a (no user impact) | double compute |

## Feature flags
- Decouple release from deploy; kill-switch a bad feature without redeploying.
- Enable trunk-based development (merge incomplete work behind a flag).
- Watch out: flag debt. Every flag is a branch in prod — remove them after rollout.

## The hard part: database migrations without downtime
Old and new code run **simultaneously** during any rolling/canary/blue-green deploy,
so the schema must be compatible with both. Use **expand → migrate → contract**:
1. **Expand** — add the new column/table (nullable, backward-compatible). Deploy.
2. **Migrate** — backfill data; new code writes both old + new. Deploy.
3. **Contract** — stop writing the old column; drop it. Deploy.
Never rename a column in one step. Never make a column NOT NULL before backfilling.

## Interview probes
- "How do you deploy without downtime?" → rolling/blue-green + backward-compatible schema.
- "You shipped a bug to 100% of users. What went wrong in your process?" → no canary, no metrics gate, no flag to disable it.
- "How do you migrate a 500M-row table's column type live?" → expand/contract + dual-write + background backfill in batches.

## Further reading
- Martin Fowler — BlueGreenDeployment: https://martinfowler.com/bliki/BlueGreenDeployment.html
- Martin Fowler — CanaryRelease: https://martinfowler.com/bliki/CanaryRelease.html
- Martin Fowler — FeatureToggles (flags): https://martinfowler.com/articles/feature-toggles.html
- GitHub — Move fast and fix things (online schema migrations): https://github.blog/2016-05-16-move-fast-and-fix-things/
- PlanetScale — Online schema changes / expand-contract: https://planetscale.com/blog/backward-compatible-databases-changes
- Google SRE Book — Release Engineering: https://sre.google/sre-book/release-engineering/
