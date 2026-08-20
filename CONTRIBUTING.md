# Contributing to base

**base** is part of the **Promise Lang** project, hosted in the `promise-language`
organization and maintained under Promise Lang LLC.

## Contributor License Agreement (CLA) required

Before any pull request can be merged, you must sign the **Promise Lang Contributor License
Agreement**. When you open your first pull request, the CLA Assistant bot will post a link to
sign. You only need to sign once — it covers all future contributions across the project.

- **Individual contributors** sign the Individual CLA.
- **Contributors acting on behalf of an employer** also have their employer sign the
  Corporate CLA.

You retain copyright in your contribution; the CLA grants Promise Lang LLC the rights it
needs to administer, distribute, and sublicense it as part of the project.

## Licensing of contributions

Unless you state otherwise, any contribution you intentionally submit for inclusion is
dual-licensed under the [Apache License 2.0](LICENSE-APACHE) and the [MIT
License](LICENSE-MIT), with no additional terms or conditions. This is core, LLC-covered
code: contributions must **not** introduce code under a copyleft license (GPL, LGPL, AGPL,
EUPL, or similar) or code of uncertain provenance.

## How to contribute

base holds the contracts every BASE participant speaks, plus the reusable machinery built on
them (see the [README](README.md)). The architecture is specified in the Reactor repo — read
[base-engineering.md](https://github.com/promise-language/reactor/blob/main/docs/base-engineering.md)
and [design.md](https://github.com/promise-language/reactor/blob/main/docs/design.md) before
proposing a change here.

1. Open an issue describing the bug or feature, where practical, so the design can be
   discussed before you invest in a PR.
2. **A change to a contract is a change to both sides of a boundary.** Every type here is
   spoken by at least two processes, so say in the issue what else has to move with it, and
   what happens to a peer running the previous version.
3. Keep the generic layer generic. Anything specific to one project belongs in that project's
   repo or its companion BASE repo, never here — see *What does not live here* in the README.
4. The gate manifest and output envelope are **language-neutral by design**. A change that
   makes them satisfiable only with a BASE library, or only in Promise, breaks the one
   deliberate polyglot boundary in the system.
5. Open a pull request and sign the CLA when prompted.

## Commit identity

Commits must carry a `@users.noreply.github.com` author and committer email. The repo's
pre-commit hook enforces this locally; a branch ruleset enforces it server-side. Activate the
hook in a fresh clone with:

```sh
git config core.hooksPath .githooks
```
