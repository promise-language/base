# base — the shared layer of Bounded-Autonomy Software Engineering

> **Status: scaffolding.** The repo exists; the contracts it will hold are still being
> specified in the [Reactor design docs](https://github.com/promise-language/reactor/tree/main/docs).
> There is no implementation here yet.

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
| **Wire types** | a flow and Reactor | claim and release, load state, resolve artifact, worktree coordination |
| **Gate manifest** | a project and Reactor | which gates exist, what each blocks, what each measures |
| **Gate output envelope** | a gate and Reactor | one gate run's result, as JSON on stdout |
| **Flow self-description** | a flow and Reactor | item types, eligibility, exclusions, session and arena hints |
| **Authority config** | a companion repo and Reactor | roles, step grants, the capability vocabulary, read scope |

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
