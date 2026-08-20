# Maintainers — repo administration

## Author-identity enforcement

This repo is public, so GitHub branch rulesets restricting commit author and committer email
are available. Apply [ruleset-author-identity.json](.github/ruleset-author-identity.json)
with:

```sh
gh api -X POST /repos/promise-language/base/rulesets \
  --input .github/ruleset-author-identity.json
```

The ruleset rejects any push whose commits do not have **both** author email and committer
email matching `11466501+djabi@users.noreply.github.com`. It applies to all branches via
`~ALL`.

Verify:

```sh
gh api /repos/promise-language/base/rulesets | jq '.[] | {name, target, enforcement}'
```

Remove (using an id from the listing above):

```sh
gh api -X DELETE /repos/promise-language/base/rulesets/<id>
```

## Local enforcement

The server-side ruleset is the authority. Locally, identity is enforced in three layers so a
bad commit is caught before it is ever pushed.

### 1. Pre-commit hook (primary gate)

[`.githooks/pre-commit`](.githooks/pre-commit) rejects any commit whose **author or
committer** email is not a `@users.noreply.github.com` address. It reads the impending
identity with `git var GIT_AUTHOR_IDENT` / `GIT_COMMITTER_IDENT`, so it catches both
`user.email` config and `GIT_*_EMAIL` environment overrides.

The hook is not active in a fresh clone until `core.hooksPath` is pointed at it:

```sh
git config core.hooksPath .githooks
git config core.hooksPath          # -> .githooks
```

Once this repo has its own dev tooling, the bootstrap will do that automatically and the hook
will exec `bin/precommit` — the full commit gate — instead of doing the identity check
inline. The hook already defers to `bin/precommit` when it exists, so that transition needs
no change here.

### 2. Correct identity in local config

The repo's local `.git/config` sets `user.name` and `user.email` to the correct noreply
identity, so a plain `git commit` in this worktree uses them and passes the hook.

### 3. useConfigOnly (belt-and-suspenders)

Refuses commits if `user.email` is somehow unset, rather than falling back to a default:

```sh
git config user.useConfigOnly true
```

## CLA

[`.github/workflows/cla.yml`](.github/workflows/cla.yml) runs the CLA Assistant against the
org-wide signature file in the private `promise-language/cla-signatures` repo, so one
signature covers every Promise Lang repo. It needs the `CLA_SIGNATURES_TOKEN` secret — a PAT
with Contents read/write on that repo — to be available to this repository.
