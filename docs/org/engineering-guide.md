# Engineering Guide

> **Tag:** `engineering-guide` — remaining work to complete this document: the query named in
> `docs/index.md`.

> **Home:** [promise-language/org](https://github.com/promise-language/org) — this document
> changes here and nowhere else. To change it, file an issue against `org`.

How code in this organization is written, in any language — naming, shape, testing, visibility,
effects, and what to do when the platform is in the way. The language-specific form of these rules
lives in the per-language guides — [`engineering-guide-promise.md`](engineering-guide-promise.md),
[`engineering-guide-go.md`](engineering-guide-go.md) — which apply this document to one language
and never contradict it. Its blockquotes are read as [normative.md](normative.md#end-state) says:
each is an invariant, and the prose under it is why.

## Why this is in the tree

A change is made inside a materialized worktree of one repository and nothing else. A rule that
lives in another repo is not in the agent's context at the moment it has to be followed, so
referencing it is the same as not having it.

**So every repository holds a copy, byte-identical to the source, and machinery keeps it that
way.** The copy is tracked content that arrives by sync and is never hand-edited: an edit to it
is refused, and the refusal points at the home repository, which is where a rule changes. The
rules of the copy — where it sits, how it is listed, how a gap against it is filed — are
[normative.md](normative.md)'s, and this guide relies on them without restating them. What is
true only of one project lives in that project's own documents, which cite this one — never in
edits to the copy.

## One obvious way

> **For every task the project offers exactly one obvious way to do it.** A second way is a defect,
> not a convenience.

Two ways to do the same thing is two things to learn, two things to search for, two behaviours to
keep equivalent, and a fork in every transcript — and the two never stay equivalent, because
nothing tells anyone when they diverge. This is the design-level form of [define once](#define-once):
aliases, fallbacks, alternate spellings, "it also works if you…", and convenience wrappers that
shadow the real path all fail it. When a second way is found, one of the two is removed.

## Naming

- **Full English words.** `print_line`, not `println`. `execute`, not `exec`.
- **The abbreviation dictionary below is the closed list of exceptions, and its entries are
  mandatory.** Where a mapping exists, the abbreviation is the correct form, not a tolerated one;
  where none exists, the full word is. It applies in every language and on every surface a name
  appears — code, flags, wire fields.
- **Proper names are verbatim.** Technologies, formats, and algorithms keep their own spelling —
  `base64`, `utf8`, `json`, `sha256`, `url` — with only the casing following the language's
  convention.
- **A quantity is a number and its unit, and neither is an abbreviation.** `5h`, `7d`, `10mb` are
  how a quantity is written, not a clipped word, and the unit is one a reader needs no knowledge
  of the project to read — `-pace-7d` answers "how long" to anyone. Where the quantity's type has
  a grammar, the name spells it as the grammar does: a duration's is
  [flag form](cli-guide.md#flag-form)'s. What the full-word rule forbids is truncating a word —
  `-h` for `-help` loses letters only context restores — and `7d` loses nothing.
- **Do not prefix a member with its type.** `response.status`, not `response.status_code`. The
  caller already has the context.

**The abbreviation dictionary.** A reader can always expand a full word but must have memorized
each abbreviation, so an entry earns its place only when the short form is so universally
recognized that it reads *better* than the word — `repo` is what `git` and `gh` call it; nobody
says *identifier*. An entry covers its plural (`args`, `repos`). The list grows by an issue
against this document, one word at a time, and it is not a licence: a word that is not in it is
spelled in full, however common its clipping. Two extensions are sanctioned, and neither removes
an entry. A language guide may add mappings particular to that language. And a project that ships
a **released product** — a language and its standard library, a published module set — may add
mappings for that product's public surface, in the project's own specification of that surface:
the names a public reader meets are governed by a table the product's authors own and its users
can read, not by one written for the organization's tooling. Everywhere else — this repository,
every project's `bin/`, every internal surface — the table above and the language guide's
additions are the whole list.

| Abbreviation | Word |
|---|---|
| `abs` | absolute |
| `arg` | argument |
| `attr` | attribute |
| `config` | configuration |
| `dest` | destination |
| `dir` | directory |
| `env` | environment |
| `func` | function |
| `hex` | hexadecimal |
| `id` | identifier |
| `info` | information |
| `init` | initialize |
| `len` | length |
| `max` | maximum |
| `millis` | milliseconds |
| `min` | minimum |
| `pos` | position |
| `prev` | previous |
| `repo` | repository |
| `src` | source |
| `stderr` | standard error |
| `stdin` | standard input |
| `stdout` | standard output |
| `sync` | synchronize |
| `var` | variable |

## Types and shape

- **A set that can be closed is closed.** An operating system is `linux` / `darwin` / `windows`,
  not a string — a string admits `macOS`, and the failure is silent: a lookup that misses falls
  through to a default and runs the wrong thing. A closed set refuses the value at the boundary and
  names what was allowed.

  **Closed is also the only reversible choice.** Opening a set later accepts values that were
  previously refused, which breaks nothing. Closing one later refuses values already written down,
  which breaks everything holding them. So the question is never "might this need to grow?" but
  "must it be open *today*?" What stays open is what a project invents and the system only carries
  — tags, metric names. What closes is anything the system must *interpret*.
- **Absence lives in the type, never in a sentinel.** A value that is sometimes not there is
  expressed with the language's optional form — not an empty string, not a zero, not a `has_x`
  bool beside the value it guards. Every sentinel spelling admits a state that means nothing; the
  typed form deletes those states and puts the check where the compiler can insist on it. An empty
  collection is not absence: no tags is a real empty list.
- **Identities are types, never bare `string` or `int`.** A function taking three `string`s takes
  them in an order nothing checks, and an item id assigned to a project id is a bug the compiler
  should have caught. Identity types are immutable value types.
- **A quantity is the standard library's type for it.** A timeout is a `Duration`, not `"30m"`; a
  moment is a `DateTime`, not an epoch `int`. The typed form parses once at the boundary and is a
  number everywhere after. If the quantity has no type yet, ask for one upstream rather than
  passing a string.

## Comments and documentation

- **Every public declaration carries API documentation**, in the language's doc form. It is the
  surface that tooling and agents read; describe behaviour, not the signature.
- **No decorative banners.** `// ── Section ──` carries no meaning, costs tokens, and rots.
- **Default to no comments.** Names carry meaning. Comment the *why* when it is non-obvious — a
  hidden constraint, a subtle invariant, a workaround and the issue it waits on.
- **Documentation is part of the change, not after it.** A change that makes a document wrong is
  incomplete; the docs corpus is load-bearing, and its cross-references are claims meant to be
  checked.
- **Published text cites repository-relative paths, never absolute ones.** An absolute path in a
  report, an issue, or a summary names a person and a machine no reader shares — unusable to
  everyone who reads it, before any disclosure rule says a word.

> **No `TODO` comments. Fix it now, or file it and reference the issue.**

A `TODO` is a backlog entry filed somewhere with no backlog semantics: nothing lists them, nothing
ranks them, nothing closes them, and nobody ever sweeps them. Stale `TODO`s train every reader to
skip all of them, which destroys the few that meant something. **The project's issue tracking
system is the single source of truth for undone work** — whichever system the project uses; a
comment may *reference* an issue, and must never *be* the record.

> **No plans, phases, or task lists in documentation either. A document defines the destination; it
> never reports the distance travelled.**

The distinction is **normative versus record**, not present versus future. *"A lease names its
holder as `(host, pid, start time)`"* is a claim about the design, equally true before any code
exists and after all of it does. *"Phase 2: lease ledger — done"* decays the moment anything moves.
The test: would this sentence need editing when work happens, even though nothing about the design
changed? If yes, it is a record, and it belongs where records live. A document points at where
status lives rather than reporting it.

## Finish it or file what is left

> **A change is done when it satisfies the normative document. Until then, the difference is an
> issue — filed before the change lands, not after.**

Landing something partial is allowed and often right. Landing it *silently* is not: an
implementation that is mostly there reads exactly like one that is finished, and the missing part
is met later as a **defect** — someone spends the diagnosis budget of a bug to rediscover a gap
that was understood perfectly well at the time.

- **The normative document defines done, not the diff.** The test is whether the document is now
  true. Where it is not, that delta is the issue, and the document stays as written.
- **File before landing.** *"I will file it after"* is the promise a `TODO` makes, and it decays
  the same way.
- **The issue says what is missing, not that something is missing** — written while it is still in
  your head, because that is the only moment it is cheap.
- **A gap found in someone else's work is filed the same way**, whether or not you are the one to
  close it.

## Do not work around the platform

> **When the language, compiler, runtime, or tooling is in the way, file it upstream and stop.**

Do not restructure code to dodge a parser bug, add redundant casts to sidestep a type-checker gap,
or reimplement something because a feature is missing. Every absorbed workaround is a defect the
platform never learns about: a gap is a platform request, not a local problem.

- **File against the owning project** with what broke, the workaround if one is unavoidable, and
  the priority.
- **An unavoidable workaround is commented and linked to its issue.** An undocumented workaround is
  indistinguishable from a mistake six months later.
- **Reuse the standard library.** A local utility that duplicates a stdlib concept is the same
  defect as a workaround, with a longer half-life.

## Define once

> **Each piece of state lives in one place. Each behaviour is implemented once.**

- **Never two copies that must be kept in sync.** Two copies are two things that can disagree, and
  nothing tells you when they have.
- **A genuine performance duplicate is an explicit cache**, with one clear owner and defined
  invalidation — never a second source of truth.
- **Implement a behaviour once and reuse it.** The minor variations in a copied implementation are
  where the two copies start disagreeing.
- **Wire types are one shared module used by both sides** — not hand-kept-in-sync copies.
- **The exception is deliberate vendoring**, as with this document — byte-identical, hash-checked,
  with its source named.

## Plan first

> **A plan names the files and functions it will change and what each change does — and it plans
> the smallest change that resolves the item.** A plan that could have been written without
> opening the repository is not a plan.

- **A bug is reproduced before it is fixed** — a minimal failing test or command. A fix for a
  bug nobody reproduced fixes a guess.
- **The boundary is part of the plan.** Say what is deliberately not being done, and why: the
  smallest change is a decision, and unstated it reads as an omission.
- **Planning and implementing are different acts.** A planner that starts implementing ends up
  describing work it has already half-done, and the plan stops being reviewable as a statement
  of intent.

> **A reviewed plan is followed, or renegotiated — never silently substituted.** Discovering
> mid-implementation that the plan is wrong is a finding: say so and stop. The plan was
> reviewed; the substitution was not.

## Keep a change to its subject

**Do not refactor or add unrelated features while implementing something else.** A change that
carries passengers is a change nobody can review, and a revert takes the passengers with it.
Improvements found along the way are their own item, filed and linked, not folded in.

An unrelated *fix* folded in silently is the worst passenger. Once the same problem is repaired
on the mainline — likeliest exactly when it is also filed, or visible to someone else — the two
are one fix spelled twice; git cannot know they are equivalent, so they conflict on every rebase
from then on, with nothing in the history to say they were the same repair. A breakage you trip
over is routed by the pre-existing-failure rule in [Sharing a mainline](#sharing-a-mainline):
fixed in place and declared, or filed and left alone — never both, and never quietly.

## Sharing a mainline

The mainline moves while a change is being made. Three rules keep the two from corrupting each
other:

- **A pre-existing failure is fixed in place or filed — never ignored, and never both.** A
  failure that reproduces on the base commit, without your changes, still has to be routed.
  A small, clearly-right repair rides along with this change, **named in it** so the reviewer
  sees a deliberate passenger — every filed item carries the fixed cost of a whole resolution,
  and filing every one-line breakage buys process with nothing. Anything larger, riskier, or
  plausibly already being fixed by someone else is filed so the breakage is recorded — and left
  alone. Doing both is how one repair gets spelled twice.
- **A duplicate fix is resolved toward the mainline.** When a conflict exists because the
  mainline already landed the repair this change carries, keeping both sides spells one fix
  twice: take the mainline's version, drop this one, and record what was dropped. Every other
  conflict integrates both sides — the intent of the change and the incoming work — never
  whichever side makes the markers go away.
- **After a rebase, the full gate runs again.** A rebase can introduce semantic conflicts git
  never flags: both sides apply cleanly, and the result is still wrong.

## Time is not a coordinate

> **Production code carries no literal time constants.** A timeout, a retry delay, a poll interval
> is a named parameter declared at one boundary — never a magic number at the point of use.

A literal `30` deep in a call is a decision nobody can find, tune, or test: it cannot be shortened
for a test, lengthened for a slow platform, or even located when it fires. Naming it in one place
is [define once](#define-once) applied to durations — and forces the question of whether the
constant should exist at all, because most stand in for a signal that was never built.

> **Nothing synchronizes on time — in tests or anywhere else.** Order concurrent operations
> explicitly: a channel, a ready handshake, an awaited task — never by sleeping and assuming an
> event happened.

A sleep standing in for a happens-before edge is load-sensitive: the window that passes on an idle
laptop collapses on a loaded runner, and no amount of lengthening fixes it. This is not a ban on
testing timing *behaviour*; it is a ban on using a clock where a signal belongs.

> **Nothing retries to make a failure go away.** Code does not re-attempt an operation in the
> hope that it succeeds, until it succeeds, or to let something else finish first. A retry
> implements a contract or it does not exist: it is admissible only where a specification names
> the failure transient and the operation safe to repeat, and the code cites that contract at the
> retry.

A retry is the time rule with the clock taken out: it assumes the failure is transient without
establishing it, and it is tuned by the same reflex — when three attempts do not clear it, someone
writes five, and nothing distinguishes 3 from 5 from 100000. It is worse than a sleep. A sleep
delays a result; a retry changes it, converting a reproducible defect into an intermittent one,
and the first failure — the one carrying the diagnostic — is the one the loop discards. A retry is
almost always arguable by analogy and almost never justified by evidence; the named contract is
how the two are told apart. A failure nobody understood is reported, with its evidence, and filed
— never re-attempted.

> **A test never reads the wall clock.** Every instant a test depends on is injected and pinned;
> a test whose outcome changes with the date it runs on is asserting the date.

> **Code that takes a clock takes all of its time from it.** Once a component accepts an injected
> clock, every reading of time inside it — stamping, expiry, pruning, backoff — goes through that
> clock. A direct read of the system clock beside an injected one is a second clock, and the two
> agree only in production, where nobody is looking.

A wall-clock dependency is the environment dependency ([testing](#testing)) that arrives on a
schedule: it passes on the day it is written, so review and the gate approve it, and it fails
later on a tree nobody touched, where bisection finds every commit equally guilty. Pinning the
test's clock is not the fix on its own — a mixed clock makes the code look testable, and the tests
it passes are valid only within a window around the day they were written. The fix is one source
of time, and the test controls it. This bans an outcome that depends on the calendar, not
measuring real elapsed time where that is the behaviour under test.

## Everything a program writes has a ceiling

> **A sink that only grows is a leak with a slower clock.** A log, a journal, a spool, a cache, a
> retained buffer, a table of past runs — anything appended to across a process's life is
> declared with a bound where it is created, the way a timeout is a named parameter declared at
> one boundary.

"It is only a log" is the reasoning that produces a gigabyte of one repeated line. The damage is
not the disk. An unbounded sink is why a malfunction runs unnoticed: the process that
malfunctions is the one that writes fastest, and past a certain size the evidence is a file
nobody opens — the pathology funds its own concealment.

- **The bound belongs to the writer** — not to an operator's cron job, and not to a reader that
  trims on the way in. A sink is unbounded until the thing appending to it says otherwise.
- **Cap by rotating and keeping segments, never by truncating to the tail.** What is worth
  reading in a log is usually when a condition *started*, and a cap that keeps only the most
  recent bytes discards exactly that.
- **For a slow-onset condition the origin is worth as much as the tail.** Rotation that drops the
  oldest segment first still loses the line that dates the whole thing; where the first record
  matters, the bound keeps the first segment alongside the last ones. A bound that destroys the
  diagnosis is not a bound worth having.

## Look for the silent classes

Correctness bugs announce themselves. These do not, so they are a standing obligation in **any code
a change touches** — not only in the lines it adds:

| | |
|---|---|
| **Memory leaks** | zero tolerance. Every allocation and resource needs a release path, and every leak is a regression rather than a pre-existing condition |
| **Lifetime errors** | double free, use after free, missing scope cleanup |
| **Concurrency races** | lock ordering, park/wake, channel close |
| **Resource waste** | handles, connections, and processes opened and never accounted for |
| **Unbounded sinks** | a log, journal, or cache with no ceiling — [everything a program writes has a ceiling](#everything-a-program-writes-has-a-ceiling) |

**Anything found here is filed at critical priority**, whether or not the current change caused it.
These are the classes that survive review, pass tests, and surface in production as something
unrelated.

## Visibility

- **Nothing is public that does not have to be.** Module-only, or narrower, by default.
- Public is a commitment: it is what other modules compile against and what wire compatibility
  constrains. Reaching for it early is how a private detail becomes a permanent one.

## Never weaken the gate

> **When a check fails, fix the cause.** Never skip, disable, delete, or mark as
> expected-to-fail a failing check to get past it — and never suppress the class of failure it
> guards as the fix.

A check that is genuinely wrong is a finding, not an obstacle: say so and stop, and the change
to the check lands as its own reviewed change. Weakened in passing, a gate keeps its name and
loses its meaning, and every later change inherits a hole exactly where something already went
wrong once.

## The tree holds deliberate source

- **Nothing lands that was not meant to.** Build artifacts, scratch files, and stray binaries
  are deleted before the change is done — not left unstaged, where the next add finds them
  again.
- **Never `.gitignore` something to get past a check.** An ignored file survives locally, keeps
  local verification passing, and breaks the mainline that never receives it — the tree that was
  verified is no longer the tree that lands.

## Testing

- **Every behavioural change is tested** — the change itself, its edge cases, and its error paths.
  A test that covers only the happy path documents the feature and verifies almost nothing.
- **A test must fail when the behaviour regresses — that is the only question that matters.**
  Would this test fail if the change were reverted? A test that passes either way is worse than
  no test, because it reads as coverage.
- **Do not pad.** A test that restates the implementation, or asserts what the type system
  already guarantees, costs review time forever and catches nothing.
- **Audit resource invariants, not just line coverage.** Every allocating type gets a cleanup
  test; concurrency code gets a stress test. A covered line that leaks is covered and wrong.
- **Code that cannot be tested is not finished — restructure it so that it can be.** An
  untestable remainder that survives the attempt is filed, naming what resists testing and what
  would have to change: a reason is accountable, a bare list is a handoff to nobody.
- **Co-locate tests with the code they test.** A separate tree is for cross-cutting integration
  tests only.
- **Tests never rely on the environment they happen to run in** — the host's locale, a
  developer's `PATH`, ambient credentials, the working directory, the current date or time
  ([time is not a coordinate](#time-is-not-a-coordinate)). A test states its whole world or
  builds it.
- **Zero leaks, and the check never gets suppressed.** There is no annotation for tolerating one.
- **Every wire contract has a conformance suite**, and every implementation passes the same one.
- Synchronization, retry, and clock rules for tests are in
  [Time is not a coordinate](#time-is-not-a-coordinate).

## No hidden effects

- **Every effect is visible at the call site.** If a function can fail, its signature says so; if
  a value is consumed or mutated, the call site shows it.
- **Self-contained by default.** Avoid global state, implicit initialization, and ambient context.
  A reader — human or agent — should be able to read a file top to bottom and know what it does.
- **Do not depend on environment details unless the task is the environment.** Code that changes
  behaviour on what it finds around it — an environment variable, a config file it was never told
  about, the shape of the host — has effects no call site shows.

> **An environment variable is never an input.** No normal use of any tool or module requires
> setting one; the sanctioned uses are the two named in the CLI guide — the tool's own debug
> diagnostics, and guard-enforced containment markers.

An environment variable is inherited and hidden: it selects behaviour invisibly, leaks into every
child process, and makes two identical command lines do different things. It is the exact opposite
of [one obvious way](#one-obvious-way). A capability reachable only by setting a variable is a
missing flag, and the fix is the flag.

## Ask rather than guess

> **A decision you cannot make from the item, the code, and the documents is asked for — never
> guessed, and never worked around.**

The ask carries the decision needed, the evidence it rests on, and a recommendation: a reader
cannot choose between options without seeing what they are choosing about. Ask only what you
genuinely cannot decide — a question is for a missing decision, not for something unread.

## Evidence over assertion

- **A claim of "already done", "cannot be done", or "not needed" carries proof** — the commit,
  the code, the reproduction. Without it, the claim is indistinguishable from giving up, and
  producing the work stays the expected outcome.
- **A concrete failing input beats a general worry.** In review, name the input and the wrong
  output it produces; a worry without one is a question, not a finding.
- **Nothing is manufactured to look thorough.** Finding nothing is a legitimate result, stated
  plainly; a padded finding spends everyone's attention on nothing.
- **Report what happened, not what was intended.** A summary describes the change that exists —
  which is not always the change that was planned.

## Prompts point here

A flow's step prompts are the natural place for these rules to be repeated, and repeating them
there is the same defect as everything else in [Define once](#define-once): a prompt that restates
a rule is a second copy that drifts, and the one in the prompt wins because it is the one the agent
read.

> **A step prompt cites this document. It does not carry its own copy of the rules.**

That works precisely because the guide is in the tree — the agent resolving an item can read it.
Prompts stay about the *task*: what this step is for, what its artifact is, what it must not do.
And most of what a prompt would otherwise repeat should be a mechanical bound instead: a rule that
could be a bound and is written as a sentence is a rule that will eventually be ignored exactly
once.

## Enforcement

> **Every resolution includes an engineering review step that judges the change against this
> guide** — the same way the normative-docs step judges it against the project's specifications.

A guide consulted only when someone remembers is advisory everywhere. The review step makes it a
gate on every change: the reviewer's question is "which rules here does this change break", asked
while the change is still cheap to fix.

A rule with a mechanical check behind it is enforced; a rule without one is upheld by the review
step, and the gap is an open item carrying this document's tag. **When a check is built, it is
built beside this guide** — these rules govern every repository, so a lint that enforces one is
written once and provisioned, not written once per project.
