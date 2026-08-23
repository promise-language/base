# base — the shared layer of Bounded-Autonomy Software Engineering

> **Status: scaffolding.** The contracts are specified in the
> [Reactor design docs](https://github.com/promise-language/reactor/tree/main/docs) and are being
> written here as they settle. `wire/` holds the identity types; the rest is still specification.

**BASE** — Bounded-Autonomy Software Engineering — is a way of building software in which
agents resolve work items unattended, bounded by what they are *able* to do rather than by
what they are asked to do. [Reactor](https://github.com/promise-language/reactor) is the
orchestrator that schedules and executes that work. This repo is the layer beneath both:
**the contracts every participant speaks, and the reusable machinery built on them.**

## What lives here

**Every in/out type in the BASE implementation.** Each of these crosses a process boundary,
and each is versioned here so that both sides of the boundary agree on it:

| Contract | Spoken between | Carries |
|---|---|---|
| **Identity types** | everyone | the references that cross an owner boundary — a project, an item, a step, a step run, a question, an article |
| **Wire types** | a flow and Reactor | every operation a flow performs, [proxied by its runner](https://github.com/promise-language/reactor/blob/main/docs/base-engineering.md#the-flow-contract) — item reads and annotations, artifact and checkpoint writes, holds, agent runs, gate runs, and proxied VCS |
| **Feed article** | any component and Reactor | one call on human attention, with its calls to action |
| **Gate manifest** | a project and Reactor | which gates exist, what each blocks, what each measures |
| **Gate output envelope** | a gate and Reactor | one gate run's result, as JSON on stdout — measurements, each carrying the type it was measured in, and the reason it measured less than usual if it did; never a verdict, and never whether the run finished, which malformed JSON already says |
| **Flow self-description** | a flow and Reactor | item types, eligibility, exclusions, session and arena hints |
| **Authority config** | a companion repo and Reactor | roles, step grants, the capability vocabulary, read scope |
| **Arena context** | the runner and the workspace setup tool | the arena's purpose, its paths, and the inputs worktree materialization needs, read over loopback |

Two of those are deliberately **language-neutral**. The gate manifest and the output envelope
are a JSON contract over a subprocess rather than an SDK interface, so a project satisfies
them by printing JSON — no BASE library, no Promise, and no code generation. That is the one
place another language legitimately enters the system, because a gate is built from the tree
it measures. Everything else here is Promise.

## What will live here

Reusable, domain-agnostic machinery — the pre-canned implementations every adopting project
would otherwise rebuild:

- **Flow common library** — the app skeleton, step execution, commit handling, push leases,
  artifact extraction
- **Gate SDK** — a convenience for Promise gates. A gate that depends on its *existence*
  has broken the contract above
- **Arena provisioning**, worktree materialization, and flow delivery
- **Dev-tooling conventions** — successor to the Go
  [forge](https://github.com/promise-language/forge) blueprint
- **Ratcheting baselines**

## What does not live here

Anything specific to one project: step composition, item types, prompts, gate
implementations, metrics, thresholds, schedules. Gates live in the project's own repo,
because a gate measures the tree and so must come from the tree. Everything else
project-specific lives in that project's companion BASE repo, out of reach of the agents it
constrains.

Keeping that boundary firm is what makes this a reusable layer rather than a place where a
directory per orchestrated project accumulates.

**Reactor's persistence layer is also not here**, and its absence is a decision rather than an
omission. Nothing outside the Reactor server ever speaks to a store — flows reach items through
the Flow API and gates emit an envelope — so both ends of that interface live in one address
space. It also keeps the two representations of an item apart: a flow sees the wire type
published here, while the server stores whatever shape suits it. Unify them and changing storage
becomes a breaking change for every flow.

The same test decides the rest. **A contract belongs here when its two ends are built by
different owners.** A contract whose ends share an owner stays with that owner, even when the two
sides deploy separately and therefore still need versioning — the runner and the Reactor server
being the case in point.

## Engineering guide

[docs/engineering-guide.md](docs/engineering-guide.md) is how Promise code under BASE is written —
naming, shape, testing, visibility, and what to do when the platform is in the way. It is **the
source other BASE repositories vendor from**: they hold a copy, because a rule kept in another repo
is not in an agent's context at the moment it has to be followed, and when a copy disagrees with
this one, this one is right.

## Design

The architecture is specified in the Reactor repo:

- [base-engineering.md](https://github.com/promise-language/reactor/blob/main/docs/base-engineering.md) — the project-facing layer: the
  invariants, gate discovery, and what lives where
- [design.md](https://github.com/promise-language/reactor/blob/main/docs/design.md) — Reactor's half: authority, reliability, topology
- [WHITEPAPER.md](https://github.com/promise-language/reactor/blob/main/WHITEPAPER.md) — the methodology

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Signing the Promise Lang CLA is required before a
pull request can be merged.

## License

Dual-licensed under [Apache 2.0](LICENSE-APACHE) or [MIT](LICENSE-MIT), at your option.
