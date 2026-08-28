# The gate contract

**Normative.** This document defines what a project's gates must satisfy, and what a runner may conclude from one. Every statement is a requirement. Where the code does not satisfy one, an issue is open against it, tagged **`gate-contract`** — this document's filename without its extension. That is the convention for every normative document here, so the gaps against one are a label query rather than a reading of the tracker, and a document and its open deltas cannot drift apart under renaming. Nothing here records progress, status, or history.

This is the layer's one contract that is deliberately **not** an SDK interface. A gate measures the tree it comes from, so it is built from that tree and written in whatever the tree is written in. A project satisfies this document by running a program and printing JSON — no BASE library, no Promise, no code generation. That is what makes it the one place another language legitimately enters the system, and it is why this contract is stated here in prose rather than only in the Promise types that happen to implement it. A gate author reading `gate/*.pr` to learn the contract is reading an implementation of it.

Two artifacts carry the whole contract: a **manifest**, which a project prints to declare what gates it has, and an **envelope**, which one gate prints to report what it measured.

## What a gate is

A **gate measures something and reports the measurements.** Coverage, size, duration, failure counts — numbers, with the type each was measured in.

**A gate never modifies what it measures.** It may write elsewhere — a build cache, a report — but the subject it reports on is exactly as it found it. Measuring and then repairing is not a gate: a measurement that changes its subject cannot be repeated, because the second run answers a question about a different thing.

**The subject is the tracked tree.** That is what makes the rule checkable rather than an appeal to intent: a runner records the tracked state before spawning and compares it after, and paths the project ignores are outside the subject, which is where a build cache and a report legitimately go. A gate that leaves a tracked file changed has broken the contract, whatever it printed.

**A gate does not decide.** The thresholds a measurement is judged against — a cap, a ratcheting baseline — are not the gate's. A gate carrying its own thresholds can be made to pass by editing the gate, and when the subject is a change written by an agent, the agent can edit it. So `test_failures: 3` is not a verdict; it is a pass or a failure depending on state the gate does not have and must not be given.

**The boundary is who may write a threshold, not where it sits.** Both a cap and a baseline live in the tree — see [The baseline](#the-baseline) — and that is what makes a verdict reproducible rather than what puts it at risk. What the artifact rule forbids is the party under judgement being able to move them.

**A gate does not report whether it finished.** A process killed for memory, or truncated mid-write by a full disk, is not alive to say so, and one that exits cleanly having measured nothing can state something false. What became of a run is the account of whoever spawned it.

**A gate is reproducible**: the same subject gives the same measurement to anyone who runs it, anywhere. That promise has two halves, and they fail identically.

- **The measurement half** is why a gate is a program and not a shell script. A script inherits whatever the environment hands it — user configuration, path differences, shell dialects — so two hosts disagree about a textually identical gate for reasons that are not about the subject.
- **The judgement half has the same shape.** A verdict is a measurement judged against a threshold, so it is reproducible only if the threshold is too. A threshold that can differ between two runs over the same subject — a per-host cache, or state on a server that moves on its own schedule — makes two hosts disagree about a textually identical gate, again for a reason that is not about the subject.

Same failure one layer up, and the same fix in both halves: **what a verdict depends on must be a function of the subject.** A gate that is a program rather than a script does not inherit the host; a threshold versioned with the tree does not vary by who is asking or when. That is why the baseline lives in the tree rather than with an orchestrator — checking out a commit from a month ago and judging it against today's baseline answers a question about neither.

## The manifest

A project declares its gates by printing a manifest. Nothing is registered by hand anywhere else: a second copy of "which gates exist" is a copy that goes stale silently, and the manifest is the only one.

Every value in a manifest is the **project's** to state. Reactor owns only what it does with them — where a gate runs, when, and what its measurements are judged against.

A manifest carries a `schema_version`, the `project` it claims to belong to, the `targets` the product is built for, an optional `preflight` command run once after a fresh checkout, and its `gates`.

**The project field is a claim, not a fact.** Reactor knows the project from the worktree's origin; a manifest naming a different one is refused rather than believed.

**Targets are open, and closed per manifest.** A project declares its own set — `wasm32` and `linux/amd64` are both legitimate and there is no universal shape to check — and Reactor never interprets a target, only matches and groups on it. A gate may speak for a target only if the manifest declares it.

### One gate

| Field | Meaning |
|---|---|
| `name` | Stable id, unique within the manifest. Keys metric history and baselines, so renaming a gate discards its history |
| `command` | The exec line. See below |
| `fix` | A deterministic remediation line, optional. **Never run by a gate runner** — it modifies, so nothing that must measure may invoke it |
| `timeout` | How long the gate's work may take. Bounds work, never queue wait: the clock starts at spawn |
| `blocks` | The transitions this gate is a precondition for. Empty blocks nothing |
| `schedule` | Monitor cadence. Absent means a pure precondition, run only when the transition it blocks is attempted |
| `target` | Which declared target this gate's measurements speak for. Absent means it speaks only for the host that ran it |
| `host_os`, `host_arch` | Where it may run. An empty filter means any, on both axes |
| `serialized_by` | Exclusions it must hold alone, each `scope:leaf` |
| `tags` | Free-form labels, used to select subsets |
| `metrics` | What it measures, and the terms each is judged on |

**What a project omits is absent, not empty.** No `fix`, no `target`, no `schedule` are each a real absence with a meaning. An empty `host_arch`, `serialized_by` or `tags` list is a real empty list — no filter, rather than a missing one.

### What a metric declares

A metric declares its `name` — the key it appears under in an envelope — its `type` (`int`, `float` or `bool`), its `direction` (`up` or `down`, which way is better, and the only way a ratchet ever moves), its `mode` (`enforced` gates a transition, `informational` is recorded and never blocks), and optionally a `cap`: an absolute bound no history may relax.

**`bool` measures a property, not a collapsed count.** Plenty of gates report a yes or no — whether the tree builds for a target, whether every source file carries the licence header, whether a vendored copy still matches upstream. Encoded as `0` and `1` those say the wrong thing: they invite comparison and arithmetic that mean nothing, and nothing distinguishes them from a count that happens to be small.

The rule that keeps `bool` from spreading is the same one that made groups worth having. **Use it when the subject genuinely has two states, never to summarise something that has more.** `tests_pass: true` is strictly worse than `test_failures: 3`: it discards the magnitude that lets a ratchet move by degrees and a regression be located, and it cannot distinguish a tree that got slightly worse from one that collapsed. A bool that replaced a count has thrown away the reason the number was measured.

**`false` is worse than `true`**, which is what keeps `direction` meaning the same thing for every type — `up` for a property that should become true and stay true, `down` for one that should become false. A bool ratchet is the strictest kind there is: once a project builds for a target, a later run saying it does not is a regression with no room to argue about degree. A `cap` on a bool fixes the value from the first run rather than from history — "true from the outset", as against the ratchet's "never worse than it has been".

**The type is declared even though the envelope states it too.** The two are a claim and its check: a run reporting a float where the manifest says int is a mismatch to name, not a widening to absorb — the envelope is refused and the run broke the contract, rather than the value being narrowed to fit. Changing a metric's type mid-history is a discontinuity rather than a conversion — the old values measured something else — so a changed type is flagged for review at manifest refresh rather than compared across. Absorbed silently it would move a ratchet that by construction never moves back, leaving a floor no later run could tell was wrong.

### Closed vocabularies

Four sets are closed, and a value outside one is refused at the boundary rather than missed at a lookup. A `macOS` written where `darwin` was meant does not error on its own — it misses a lookup and runs the wrong command, or nothing.

| Set | Values |
|---|---|
| `host_os` | `linux`, `darwin`, `windows`, `any` |
| `host_arch` | `amd64`, `arm64` |
| `blocks` | `commit`, `branch.create`, `push:branch:own`, `push:origin`, `pr.create`, `pr.merge` |
| exclusion scope | `project`, `host`, `arena`, `global` |

**Closed is also the only reversible choice.** Opening a set later accepts values previously refused, which breaks nothing; closing one later refuses values already written down, which breaks everything holding them — and the wire evolves additively, so there is no version in which that becomes safe.

A gate **blocks** a transition; it cannot invent one. The set above is the VCS half of the capability vocabulary, and an entry outside it names something no authority can grant.

An exclusion's **scope** is closed because the scope is what makes its leaf resolvable — `project:` resolves against the item's project, `host:` against the machine. The leaf is opaque on purpose, so a project can invent `project:migration` without a change to this layer.

### Rules a manifest must satisfy

A manifest that breaks any of these is refused whole, rather than accepted with the offending gate dropped:

- Gate names are unique within the manifest.
- Every gate declares a positive timeout.
- No gate is inert. A gate declaring neither a transition to block nor a cadence to run on can never run, and a manifest carrying one is rejected rather than left to sit.
- No gate speaks for a target the manifest does not declare.

## The exec line

A command is one line: a program and its arguments, optionally dispatched per host OS by an override that wins over the default on the OS it names.

**An exec line is exec'd, never interpreted.** `bin/gate test --wasm` becomes the program `bin/gate` with the arguments `test` and `--wasm`. No shell: no globbing, no pipes, no `&&`, no redirection, and no quoting grammar. The line splits on whitespace, and that is the whole grammar.

The reason is portability, not preference. This must work under PowerShell as well as `cmd.exe` and the POSIX shells, and they disagree about quoting, escaping and operators — so a shell-interpreted line would mean different things on different hosts. Per-OS dispatch cannot rescue it, because a project would have to know which shell the runner happened to invoke rather than only which OS it was on. It also keeps a boundary the design depends on: a shell is capability the arena grants, and a gate comes from the tree.

**A gate that needs shell features writes a script and names it**, which puts the interpreter in the project's own file where it is chosen explicitly.

**The runner appends `--envelope`.** The manifest declares only the gate's own line; the flag is protocol, not project configuration, so it has one spelling everywhere and a runner adds it without being told to. It is appended last, after whatever the project declared, because that is the only rule that works without parsing the line: `bin/gate test --wasm` is run as `bin/gate test --wasm --envelope`.

**A gate prints an envelope only when it was given the flag.** Any other invocation is a person or an agent at a terminal, and must produce nothing on stdout and a non-zero exit. Silence and failure are what keep the human path from becoming a second channel — a bare invocation that printed measurements and exited `0` would be read as a pass by the first script that wrapped it, which is the ambiguity the three parties exist to remove.

What it prints on stderr is what the caller needed: the gate's name, what it measures, and the command that runs it.

```
$ bin/gate test
test — measures test_count, test_failures, excluded_count

  bin/run test
```

## The envelope

A gate prints **one JSON object on stdout**, and that is its entire structured output. Human-readable progress belongs on stderr, which nothing parses.

**Progress reaches the reader as it is written.** Gates run long — a test suite, a full build — and a reader watching one needs to see it working rather than a silence that is indistinguishable from a hang. So stderr passes through unprocessed and unbuffered: the runner does not collect it, reformat it, prefix it or hold it until the end. The gate is writing to the reader, and the runner is not in the middle.

**The gate is given the reader's own stream, not a pipe the runner copies.** That is the requirement, not one way of meeting it — "forwarded faithfully" is a plausible reading of the paragraph above that satisfies every word of it and defeats it. A runner that pipes stderr and copies each line onward has changed the one thing that matters: the gate is now writing to a pipe rather than to a terminal, and most runtimes switch from line to block buffering when their output is not a terminal. Ten minutes of progress arrives as one block at exit — precisely the silence the rule exists to prevent, produced by an implementation that reads as correct.

**And it fails silently.** A gate whose progress arrives all at once still prints a valid envelope and still gets a correct verdict, so nothing errors and nobody investigates. The only symptom is a reader watching a gate that appears to have hung, which is indistinguishable from a gate that is simply slow.

An orchestrator that needs progress in a durable record is the one case that must pipe, and it accepts the loss of liveness knowingly rather than discovering it.

**Only then does the runner speak.** The structured output — the verdict, the measurements against the terms they were judged on — is produced after the envelope has been read. The gate's progress and the runner's summary never interleave, because they do not overlap in time.

**Progress does not extend a deadline.** A gate that prints forever still reaches its timeout: the timeout bounds work, and a gate that is talking is not thereby making progress. Resetting a deadline on output would turn a wedged-but-chatty gate into one that runs until something else kills it.

The bound in the other direction is on what is *held*, not on what is emitted. stdout is read into memory and parsed, so it is capped; stderr is never held by the runner at all, so there is nothing to cap.

One gate run reports **one target**. That is invariant, not a convenience.

An envelope carries a `schema_version`, the `target` its measurements speak for, its `metrics`, any per-part `groups`, and an `incomplete_reason` when it has one.

**Measurements carry the type they were measured in.** A bare `5` is an int and a float alike, and JSON cannot tell them apart — so a reader could not reconstruct what a gate wrote without the manifest open beside it, and a value disagreeing with what the manifest declared would be a widening nobody could see rather than a mismatch someone can name. Integers stay integers: a test count encoded as a float is a count nobody can compare exactly. `bool` is the one type JSON describes on its own — `true` is not a number in any reading — and it is declared anyway, so that every measurement is checked against the manifest by the same rule rather than one type being trusted because the wire happened to be unambiguous about it.

**Groups say where.** A total says a ratchet moved; a group says which package, file or suite it moved in, which is the difference between a bisect and a glance. A group carries the same metric names the manifest declared, scoped to one part of the run.

**Completeness is the absence of a reason.** A run that measured less than usual — a suite skipped for a missing tool, a subset deliberately selected — carries the reason it did. A run that measured everything carries nothing. There is no separate completeness flag, because a flag is a second field a reader could find disagreeing with the reason sitting beside it. An incomplete run with nothing to say is the one state an envelope cannot mean, and is refused.

**A baseline is never moved from an incomplete run.** Honest numbers that understate the tree are indistinguishable from a regression unless the run says so, and ratcheting on one lowers the floor for a reason that is not about the code.

**There is no marker for "the envelope is whole."** A gate killed mid-write emits malformed JSON, which already fails to parse; a flag meaning that would say nothing the parser does not.

**An envelope that arrived over the wire has been through none of the sending side's checks.** A consumer validates after decoding, so a gate written in another language is held to exactly the rule a gate written in Promise is.

## Running a gate

**A gate is never invoked directly.** A caller asks a runner to run it and reads what the runner reports; spawning the process is the runner's job. Three parties with three jobs — **the gate measures, the runner observes, verify judges** — and the exit code a caller sees belongs to one of the last two, never to the gate.

**The gate's own exit code is not consulted.** The states that matter most are ones a gate cannot report, because it is not alive in them: killed at the declared timeout, killed for memory, truncated mid-write on a full disk. And the substitution fails in the safe-looking direction too — a gate that exits `0` having printed nothing has stated something false, and a caller reading that number believes it.

This is also what keeps the exit code from being a second channel. The gate has exactly one output; everything else a caller learns is the runner's account of what it watched. **Two channels that could disagree never arise, so there is no rule for when they do.**

**The runner comes from outside the tree.** That is the boundary the thresholds sit behind and it is there for the same reason: a gate an agent can edit must not decide whether that agent's change passed, and a runner taken from the worktree could be edited to skip. Being remote is not what makes a runner trusted — one that never speaks to a server and only keeps tabs on what it started satisfies this exactly as well. **The property is whose artifact it is.**

### What the runner reports

One of five outcomes, and they are the full set:

| Outcome | Means | Retry |
|---|---|---|
| **measured** | The process completed and printed a valid envelope. Whether the numbers are acceptable is a different question | — |
| **timed out** | Killed at the manifest's declared timeout | maybe |
| **could not start** | The program named by the exec line is absent or not executable. Nothing ran | never |
| **died** | Killed by a signal, or exited without printing a readable envelope | yes |
| **broke the contract** | The process printed something that is not an envelope, disagreed with what the manifest declared, or modified the subject it measured | never |

**One of these is decided before a process exists.** `could not start` is what the spawn itself reports — the exec syscall refuses a program that is absent or not executable — so it is settled with no output to read and no exit code to have. The other four are decided by what the process did. (A missing program surfacing as exit status `127` is a shell reporting it, which means a shell was between the runner and the program; the exec line forbids that, so a runner that sees `127` has a different defect than the one it looks like.)

The last three outcomes are separated because they are owned by three different people. `died` is a condition of the host and a retry is the correct response. `broke the contract` is a defect in the gate's own code, and `could not start` is a defect in what the manifest declared or what the arena delivered — both recur on every retry forever, and a retry loop pointed at one reads as a flaky host for as long as anyone lets it run. An orchestrator that cannot tell them apart attributes the failure to the wrong repository.

### When a gate breaks the contract

**A violation is not a failing measurement, and is never recorded as one.** A gate that broke the contract has not reported that its subject is bad; it has failed to report anything. Conflating the two is the expensive mistake, and it is expensive in both directions — read as a failure it blames a change for a defect in the gate, and read as a pass it lets a transition through on a measurement nobody made.

Three things follow, and they are requirements:

- **The run yields no measurements, and no baseline moves.** A violation that fed a ratchet would move a floor for a reason that is not about the code, and a ratchet by construction never moves back — so a single violation absorbed silently leaves a floor no later honest run can meet, and none of them can tell why.
- **A transition the gate blocks is not allowed, and not because the subject failed.** The gate is a precondition and it did not speak, so the precondition is unmet. Reporting that as a failing gate points the next person at the change instead of at the gate.
- **It recurs.** `broke the contract` and `could not start` are both defects in the project's own artifacts, so a retry reproduces them exactly. Retrying is a loop, not a recovery.

**A modified worktree is spent.** The remaining gates selected for that transition must not run in it: they would measure a tree that no one proposed and no one reviewed, and their numbers would be honest about the wrong thing. The arena is re-materialized or the transition is abandoned; there is no partial recovery, because nothing downstream can tell which gates ran before the modification and which ran after.

## Running one gate by hand

A person or an agent wanting one gate's result runs **`bin/run <gate>`**, and that is the same path the runner takes rather than a parallel one. Running a single gate is common and is not a lesser case: it is faster than every gate that blocks a transition, and it is what someone iterating on one failure actually wants.

`bin/run` is set up by the project's `./make` alongside `verify` and the rest of the tooling, so it is generated rather than authored — a project does not write a launcher per gate, and does not write this one either.

**It reaches a verdict, and prints it for a human.** Rendering the envelope is the judging layer's job because it is the only layer that can do it: it holds the caps, the directions and the baselines, so it can put a number beside the terms it was judged on. A gate could only ever print the left-hand column.

```
$ bin/run format
unformatted_files    3   cap 0          ✗
```

That is also why a gate keeps exactly one output mode. A gate that pretty-printed when it thought a human was watching would have two, and one of them would not parse.

## The baseline

A cap is declared in the manifest and changes only when a person edits it. A **baseline** is derived: the best a metric has been, which a passing run ratchets in the declared direction and which never moves back.

**A baseline lives in the tree it measures, versioned with it.** That is what makes a verdict reproducible, and the reasoning is the same one that made a gate a program: what a verdict depends on must be a function of the subject. A baseline held by an orchestrator moves on its own schedule, so checking out last month's commit and judging it against today's baseline answers a question about neither tree — and adding thirty tests today would retroactively fail every tree that came before. Versioned with the tree, a commit carries the terms it was judged on, and `bin/run` reaches a verdict offline, on any machine, for any commit.

**Only the runner writes it.** A baseline moves when a complete, passing run ratchets it, and by no other route. It never moves from an incomplete run: honest numbers that understate the tree would lower a floor for a reason that is not about the code, and a ratchet by construction never moves back.

**A resolution's diff may not contain the baseline.** This is the artifact rule at the one place where the threshold and the subject share a tree: an agent that can move a baseline can pass itself, and a lowered floor is invisible afterwards because no later run can tell it was wrong. A change authored by a resolution that touches the baseline is refused — not merged and then flagged, because by then the floor has moved.

The rule is about the **author**, not the file. A person lowering a baseline deliberately, in a reviewed change, is doing something legitimate that the ratchet has no other way to express — a metric that got worse for a reason the project accepts. Forbidding that outright would make the ratchet a floor nobody could ever lower, which is a different failure. What is forbidden is the party under judgement moving it as a side effect of being judged.

## Where the verdict is made

The verdict is computed from an envelope's measurements, against the caps in the manifest and the baselines in the tree. It exists in neither the envelope nor the exit code, and it is reached outside the tree by something the tree cannot edit.

**Verify is a selector over the manifest, not an entry in it.** The set of gates that must be green before a transition is allowed on a host of this shape is derived from what the gates declare they block, filtered by where they may run and by any tags the caller narrowed to. Nothing declares "the verify set" separately, so the check a developer runs and the check that refuses the push are read out of one declaration and cannot drift.

**A verdict rests only on a gate's measurements.** Nothing permitted to modify its subject may be the thing a decision reads, however convenient its output. A repairing tool asked *"did this pass"* answers about a state that did not exist when the question was asked, and the tree it describes is not the tree anyone proposed — so a pass means the repairs were applied, which is a different claim from the one a reviewer, a bisect or a later rebuild will check.

That is why `fix` is never run by a gate runner, and the rule belongs here rather than only there: **it constrains the caller, not the tool.** A repairing command that repairs is behaving correctly, so no check on it can catch this — the defect is entirely in what was asked, and it looks like a green result rather than like an error.

**Eligibility is a filter, never a verdict.** A gate that may not run on this host has not passed and has not failed; it has not spoken.
