# Documentation Index

This is the map of `docs/`. It is the one file in the root that is not a specification —
everything else there is.

**The rules are defined once, org-wide, in [org/normative.md](org/normative.md)**: which
locations bind and which do not, the tag header, where status lives and where it never does,
one fact one home, the lifecycle of a specification, and what is enforced mechanically. This
index does not restate them.

**This project's status query.** Each root document's tag is a GitHub label, spelled as the
file's basename minus `.md`, and the remaining work for a document is:

> `gh issue list --label <tag> --state open --limit 200`

`--state open` and `--limit 200` are both written out deliberately: `gh issue list` defaults to
a limit of 30, so a specification with more remaining work than that would silently
under-report and read as nearly done.

---

## Specifications

- [gate-contract.md](gate-contract.md) — What a project's gates must satisfy, and what a runner
  may conclude from one.

## Organization-wide corpus — binding

Vendored from [promise-language/org](https://github.com/promise-language/org) at the release
named in [org/stamp.json](org/stamp.json). Never edited here: an issue about one of these
documents is filed against `org` (org/normative.md §7); what this project files locally under
their tags is its own compliance gaps.

- [org/normative.md](org/normative.md) — What makes a document binding, and the one docs
  structure every project holds.
- [org/engineering-guide.md](org/engineering-guide.md) — How code in this organization is
  written, in any language.
- [org/engineering-guide-promise.md](org/engineering-guide-promise.md) — The engineering guide
  applied to Promise source.
- [org/engineering-guide-go.md](org/engineering-guide-go.md) — The engineering guide applied to
  Go source.
- [org/cli-guide.md](org/cli-guide.md) — How every command-line tool behaves at its invocation
  surface.
- [org/stamp.json](org/stamp.json) — The version stamp: the org release these copies came from,
  with per-file hashes.

## Archive — superseded or delivered

- [archive/engineering-guide.md](archive/engineering-guide.md) — base's own engineering guide,
  once the source other repositories vendored. Superseded by
  [org/engineering-guide.md](org/engineering-guide.md), whose home is the org repository.
