# Logging

> **Tag:** `logging` — remaining work to complete this document: the query named in
> `docs/index.md`.

> **Home:** [promise-language/org](https://github.com/promise-language/org) — this document
> changes here and nowhere else. To change it, file an issue against `org`.

How the development tools and the orchestration system record what they did: one line format, one
place for logs on each machine, one set of bounds, and one path from every machine to a store
where the whole fleet is read at once. A malfunction that crossed a governor, a runner, a flow
binary and a gate is reconstructed from logs after the fact. So every line can be placed on its
host, its arena, its tool and its process, and every process on its parent, without asking the
machine that wrote it.

## Scope

> **A program is bound by what it is for, not by the repository it is built in.**

Bound:

- **Every tool a project's tools module builds** — `make`, `setup`, `verify`, `gate`, `run` and
  every other tool in a project's `bin/`.
- **Every workspace tool**, and every flow binary, including everything the flow SDK does inside
  one.
- **Every program of the orchestration system the organization provides** — the governor, the
  runner, the orchestrator's server and hosting — and every program they start to supervise work.

Not bound:

- **A project's product**: what a project ships to its users. The Promise compiler running on a
  user's machine is a product; the tools that build, test and gate it are bound. A project may
  follow this document for its product. Doing so binds nothing, and following it only in part is
  not a defect under it.

A product runs on machines the organization does not operate, for people who never joined its
fleet. A log home, an identity and a store written onto such a machine would install the
organization's diagnostics on someone else's computer. And a product's logging is a feature its
users judge, which the product's own specification decides.

**The governor's first stage is bound, and ships nothing.** Stage 0 downloads the current stage 1,
runs it and keeps it running. It writes its log like any other bound program, because when stage 1
cannot be fetched or will not stay up, that log is the only account of why. Shipping, and every
other part of the governor that needs changing over time, belongs to stage 1. Stage 0 does not
update itself. It changes only when a person replaces it by hand on each machine, which is costly,
rare, and never reaches every machine. So any version of stage 0 may still be running somewhere,
and whatever it holds has to stay correct for as long as that copy runs. Stage 1 ships stage 0's log
with everything else in the machine's log home ([shipping](#shipping)).

So stage 0's log is the part of this document least able to change. A line type, a field or a
bound that stage 0 writes stays in the fleet until the last hand-installed copy is gone, and an
amendment that changes what stage 0 writes has to be read by the store alongside every earlier
form.

This document does not own:

- the identities a line names, which are [identity](identity.md)'s;
- stdout and stderr, which are [the CLI guide](cli-guide.md#output-modes)'s — a log is neither;
- a record that decides what happens next, such as a journal, a ledger or a step's result, which
  its owner's specification defines;
- what a program records beyond the types below, which its own specification defines under
  [line types](#line-types).

## A log is evidence

> **Nothing reads a log to decide what to do.** A program behaves the same whether its log was
> written, could not be written, or has since been deleted.

A log that one program reads to decide becomes a protocol that nobody versioned, with its reader in
another repository. And the [bounds](#bounds) discard lines by design, so the first time a bound
removed the line a decision depended on, a diagnostic would turn into a wrong answer. What decides
is a journal, a result or a record, with an owner and a contract. What a log holds is what
happened.

> **A failure to log never fails the program.** A write that cannot proceed is reported once on
> stderr and the program goes on. The lines that could not be written are counted, and
> [process end](#process-end) carries the count.

## Where a log is written

> **Every bound process writes its log under `logs/` in its machine's host home**, whose path on
> each platform [identity](identity.md#where-records-are-kept) states: a process in an arena exactly
> as a process outside every arena. That directory is the machine's **log home**. A host has one,
> and so does each guest on it: a container or virtual machine running on the host
> ([identity](identity.md#the-five-identities)).

> **The log home's layout is the contract between every writer and the shipper, and the log home
> holds nothing else:**
>
> ```
> <host home>/logs/
>     .shipper.lock     held by the shipper while it ships
>     .pruned           its modification time records the last pass of the bound
>     <tool>/           one directory per tool, holding that tool's files
> ```
>
> A writer creates `logs/` and its tool's directory when either is missing, readable and writable
> by the account alone. A writer that cannot place its files there writes no log
> ([a log is evidence](#a-log-is-evidence)).

A writer and a shipper are built in different repositories, released on different days, and meet
only in this directory. The path, the three kinds of entry and their names are therefore all
stated here, and neither side is told where the other one is. The directory is private to the
account because a log carries command lines and paths.

> **The log home is a shared location that cannot collide, and every writer reaches it through one
> implementation.** Its path and layout are stated here and in [identity](identity.md#where-records-are-kept),
> every file in it is named by a [unique id](#files), and each language has exactly one
> implementation of this document's writer, which every bound program in that language uses.

Workspace's `docs/tooling.md`, *An agent writes only where it is working*, keeps project and
workspace tools out of shared directories because two arenas writing the same `/tmp/test.csv`
overwrite each other's work, and neither can tell. A log home has none of that failure:

- **Its location is defined.** No writer chooses a path, so no two choose the same one by accident.
- **Its names are unique by construction.** Every name carries a 128-bit id, so no two writers on
  any machine ever write the same file.
- **Its implementation is shared.** A second implementation would drift from the first: a
  different name, a different layout, a bound held differently. And a writer that drifted is
  exactly how a shared location stops being collision-free.

The log home is outside every arena on purpose:

- **An arena is reset, reused and removed by others.** A worktree reset, an arena's setup, and a
  person clearing scratch all touch the checkout, and each is likeliest to run just after a
  failure. A log inside the checkout would go with it, or need the shipper present at every one of
  those moments to rescue it. Outside the arena, creating, cleaning or removing an arena needs no
  shipper and loses nothing.
- **The agent whose work a log records is out of the log's reach.** An agent may write anywhere in
  its checkout, and the guard refuses its writes anywhere else. A log in the log home is evidence
  that agent cannot edit.
- **One directory per machine** is one place for the shipper to watch and one bound to hold,
  however many arenas the machine carries.

> **A log home is never cleared.** What empties it is its bound and its shipper, and nothing else.

The runs worth reading are the ones that went wrong. A clear of any kind removes exactly those,
before anyone has come to read them.

## Files

> **Each process writes its own files, and only its own.** A file is named for the process that
> writes it:
>
> ```
> <log home>/<tool>/<started>-<process id>-<segment>.jsonl
> ```
>
> `<started>` is the moment the operating system says the process started, in UTC, in the basic
> form of ISO 8601 to the microsecond — `YYYYMMDDTHHMMSS.ffffffZ` — with digits the platform does
> not measure written as zeros. `<tool>` is the process's bare name, `<process id>` its drawn
> [id](identity.md#ids) — `p-` and 32 hexadecimal digits, never the operating system's pid — and
> `<segment>` counts from `0`.

```
~/.local/state/promise-language/logs/verify/20260917T081102.123456Z-p-9f86d081884c7d659a2feaa0c55ad015-0.jsonl
```

One writer per file means no locking and no interleaving, and a torn write damages nothing but its
own writer's last line. It also means a crashed process's log is whole up to the crash, instead of
being mixed into a stream it shares with processes that went on.

> **A file's name is unique everywhere, and never changes.**

The process id in the name is the same value every line carries as `process`, every child names as
its parent, and the store keeps each line under. One value names the process wherever it appears,
and nothing is derived from a name or parsed back out of one. The id is 128 random bits, so no two
files anywhere share a name: not on one machine, not in a directory someone gathered every host's
logs into, and not wherever a file is shipped or copied.

The start instant is in the name for the person listing a directory. The basic form sorts in time
order, so a tool's files list in the order its processes started, and it has no `:`, which one
platform refuses in a file name. The tool is a directory rather than part of the name, because tool
names contain `-`, and a name joining several such parts could not be split back into them.

Everything else about a process — its host, guest, arena, pid, parent and command line — is on the
first line of each of its files ([process start](#process-start)), not in the name. A person
starting from a pid finds the file with [`logs`](#reading-a-log), or by reading the first line of
each file in the tool's directory.

> **A file is created exclusively, and never opened for writing again once its process has moved
> on.** A name that already exists is not written to: the process reports the collision on stderr
> and writes no log. Creation that would overwrite is how two processes' evidence becomes one file
> nobody can separate.

> **A segment is created under its final name and is never renamed or moved.** When the next line
> would take a segment past its [bound](#bounds), the writer closes that segment and creates the
> next one, numbered one higher, whose first line is a [process start](#process-start). Segment
> `n` is never written again once segment `n + 1` exists, so the existence of the later segment is
> what says the earlier one is finished.

Rotation by renaming — the live file keeps one name and its full predecessor moves to another — is
avoided for three reasons:

- **It breaks the shipper's coordination.** A reader holding a name and an offset finds different
  bytes under the same name. It must detect that, and if the rename lands mid-read, it loses lines
  or ships them twice.
- **It changes what a finished file looks like.** A rename is the one act a reader cannot tell
  apart from a file still being written, without racing the writer.
- **It fails on a platform that will not rename a file another process holds open**, which a
  shipper or a reader always may.

A file is created where it stays, under the name it keeps, until it is removed. A name, once found,
names the same bytes for as long as the file exists.

## The line

> **A log is JSON Lines: UTF-8, one JSON object per line, each ending in a newline, none spanning
> two.**

> **Every line carries exactly five fields: `time`, `sequence`, `process`, `type` and `data`.**

| Field | Holds |
|---|---|
| `time` | when the line was written, in RFC 3339 UTC with microseconds — `2026-09-17T14:03:07.123456Z` — read from the process's one clock |
| `sequence` | the line's position in its process's log: `0` for the first line, one more for each after, continuing across segments |
| `process` | the process's drawn [id](identity.md#ids), `p-…`: never the operating system's pid, which [process start](#process-start) records as `pid` |
| `type` | the line's [type](#line-types) |
| `data` | an object holding the type's fields, `{}` when it has none |

Five fields, because everything that places a process — host, guest, arena, tool, version, parent
— is the same on every line it writes. It is stated on the [process start](#process-start) line
and joined by `process`. Repeating it on every line would spend most of each line's bytes, against the
bound, on facts that never vary.

`sequence` is what makes loss visible. A gap in a process's sequence is lines that were written and did
not arrive, and that is the one kind of loss a store can report instead of hiding. Together with
`process` it is also what identifies a line, so a line delivered twice is kept once.

> **A field is added to a type, never renamed or repurposed, and absent means unknown** — the rule
> [the CLI guide](cli-guide.md#output-modes) sets for JSON on stdout, for the same reason: the
> reader is in another repository.

> **A line is written with a single write, whole, ending in its newline, and a written byte is
> never changed.** A log changes only by a line appended at its end, or by a whole file removed
> under [the bounds](#bounds) or by its [shipper](#shipping). This binds a machine's log home. What
> the store keeps and deletes is [the store](#the-store)'s.

That is the whole of the coordination between a writer and its shipper. The shipper needs no lock
the writer takes, no signal and no socket: everything up to the last newline is lines that exist,
and anything after it is a line still being written.

> **A line is at most 64 KiB.** A writer cuts the string values in `data` to fit, longest first,
> and marks the line with `"truncated": true` in `data`.

> **A line never carries a credential, a token, or the contents of the environment.** A command
> line is logged, because [the CLI guide](cli-guide.md#explicit-inputs) has a tool name the
> source of a secret, never the secret, and whatever reaches a line anyway is [sanitized](#sanitizing-a-line).

## Sanitizing a line

> **Every line passes through the log sanitizer before it is written, and the sanitizer replaces
> what it detects in place.** It runs inside the one writer implementation, on every string in
> `data` — every field, every argument and every message alike — before the line is cut to its
> bound.

> **The log sanitizer is not a guard.** A disclosure guard, as flow's `docs/disclosure.md` defines
> one, refuses text, returns it to its author, and never modifies what it examines. A log line has
> no author to return it to, and a refused line would be a lost one. So the sanitizer edits
> instead: it replaces what it detects and writes the rest.

> **It runs as the line is written, never afterwards.** The log home never holds what the sanitizer
> removed, and [a written byte is never changed](#the-line) stays true of it.

Sanitizing at shipping instead would leave the unsanitized line on the machine, readable by
anything the account runs, for as long as the bound keeps it. It would also make the store's copy
differ from the machine's.

What the sanitizer replaces, each by a marker naming its category:

| Category | Detected as | Replaced by |
|---|---|---|
| **Credential** | a value in a form credentials take — a token, a key, an authorization header — whether or not it is live | `<credential>` |
| **Email address** | an address in the form `local@domain` | `<email>` |
| **Internal address** | a private network address, or a host name under `.local` or `.internal` | `<address>` |
| **Path** | the rule below | `<path>`, or the path rewritten |

The categories are those of flow's `docs/disclosure.md`, as far as each can be detected in text
without knowing where the text came from. This document owns what the sanitizer does with them:
when it runs, what it writes in their place, how it marks a line, and what happens when it cannot
run.

> **A path is kept only where following it stays inside the account's home directory or the
> arena's checkout.** A path under the home directory is written as `~/` followed by the rest of
> it. A path inside the checkout is written relative to the checkout root. A relative path is
> written as it was given. Each is kept only when following it, through every `..` and every
> symbolic link or junction along the way, never leaves the home directory or the checkout. Any
> other path is replaced by `<path>`. `working_directory` is recorded only when the working
> directory is inside one of the two.

> **A path is written the same way on every platform: `/` between its parts, and `~` for the home
> directory.**

| Platform | Home directory | An absolute path, as free text shows one |
|---|---|---|
| **Windows** | the account's profile folder, as the platform's own call reports it: `C:\Users\jane` | a drive letter and a separator (`C:\`, `C:/`), a UNC path (`\\server\share`), or an extended-length path (`\\?\C:\…`) |
| **Every other platform**, macOS and Linux included | the account's home directory, as the platform's own call reports it | a path beginning with `/` |

On Windows a path is compared without regard to case, `\` and `/` are the same separator, and an
extended-length prefix is removed before comparing. `C:\Users\Jane\prog\forge_1\docs\x.md` and
`c:/users/jane/prog/forge_1/docs/x.md` are both under the home directory, and both are written
`~/prog/forge_1/docs/x.md`. Linux running under Windows is not Windows: its home directory is its
own, and the Windows profile it can reach at `/mnt/c/Users/jane` is outside that home.

One form, because the store reads lines from every platform. A query for a path, or a person
comparing two lines, should not depend on which machine wrote them. A `\` would also have to be
escaped in every JSON string that carried it.

A path outside those places is where a person's name hides: a Windows profile mounted into Linux
(`/mnt/c/Users/jane`), a drive named for its owner, a home in a non-standard place. So does a path
that leaves them by its own steps: `../jane`, followed from a working directory that is the home
directory, names the account next door. What is kept is what diagnosis needs: the checkout, and
the developer tools installed under the home directory.

In a field the writer knows to be a path — `working_directory`, a program, an argument naming a
file — every path, absolute or relative, is judged by this rule. In free text a path is recognised
by its form: an absolute path in its platform's form, as the table shows, or one beginning with
`~`.

> **A line the sanitizer changed says so.** Its `data` carries `sanitized`: the list of categories
> replaced, never what was replaced.

> **A line the sanitizer cannot process is written as a stub.** Its `time`, `sequence`, `process` and `type`
> are kept, and its `data` is replaced whole by `{"sanitizer_failed": true}`.

Nothing unsanitized is ever written. The line's place in its process's sequence still shows, so a
reader sees that something happened there. And the program goes on, as [a failure to
log](#a-log-is-evidence) requires.

> **`truncated`, `sanitized` and `sanitizer_failed` are this document's, and may appear in the
> `data` of any type.**

> **What the sanitizer cannot detect, a log may still hold.** A name in a commit message passed as
> an argument, or a person's name inside a value an author chose to put in a message, can reach a
> line.

Detection is not recognition. No sanitizer can tell a person's name in free text from any other
word, and claiming otherwise would promise a protection nothing delivers. What protects those lines
is who may read them, how long they are kept, that they can be deleted, and that they are never
published ([the store](#the-store)).

## Line types

> **Every type is defined in exactly one specification, which lists its fields.** A type's name is
> its owner and a name, `<owner>.<name>`. The owners `process`, `child` and `log`, and the type
> `message`, are this document's. Any other owner is the repository whose specification defines
> the type: `flow.dispatch`, `governor.stage-restart`, `gate.measure`.

A line whose type no specification defines is a defect in the program that wrote it. So is a
field its type does not list, because an undefined field in a log is a protocol somebody will one
day parse.

> **A line carries facts about work, never the work's content.** No type has a field holding the
> text of a prompt or a response, source, a diff, the body of an issue or a comment, or a document.
> A line about a prompt carries its ids, its sizes, its timing, its cost and its outcome. An
> argument its spawner knows to carry such content — a prompt passed on a command line — is
> written as its length, never its text.

Content already has records of its own, each under its own rules — flow's prompt record, which is
never published, is one — and a log is none of them. A log holding content would be a copy of
those records with none of their rules, sent to a store where more people can read it.

### Process start

The first line of every segment. On segment `0` it is the process's first line. On a later segment
it restates the first one, so each segment can be placed without the others.

| Field | Holds |
|---|---|
| `segment` | the segment this line opens |
| `host` | `{id, label}` of the host: the host home's record, or, in a guest, the host it names |
| `guest` | `{id, label}` of the guest, when the host home's record is a guest's |
| `arena` | `{id, label}` from the checkout's arena record, for a process in an arena |
| `tool` | `{name, version}`: the bare name, and the whole `-version` object |
| `pid` | the process's pid |
| `started` | when the operating system says the process started |
| `boot` | when the operating system says the machine last booted, truncated to the minute: `2026-09-17T06:42:00Z` |
| `ancestors` | `[{pid, started}, …]`, nearest first: the process's ancestry as the operating system observed it ([identity](identity.md#processes)) |
| `args` | the command line, after the program name |
| `working_directory` | the working directory, when [sanitizing](#sanitizing-a-line) keeps it |
| `dropped` | on a later segment, how many segments between the first and this one the bound has removed |

`boot` is how a copied identity record is found. One host id seen from two different boot times at
the same moment is one record on two machines ([identity](identity.md#ids)), and nothing else in a
log could show that.

> **`boot` is truncated to the minute, and never carries seconds or anything finer.**

A boot time to the second or finer is close to unique for a machine, so it would let the same
machine be recognised across labels, ids and records that were meant not to be joined. To the
minute it still tells two machines apart unless both booted in the same minute, and that is all it
is recorded for.

### Process context

The ids of the work a process is serving, written when it learns them. A later line replaces the
fields an earlier one set: `item`, `step`, `step_run`, `invocation` and `session`, each in the form
its owner defines. A child's lines are placed in the same work through its parent, so a child
does not restate its parent's context.

### Child start and child end

Written by a process for every child it starts, bound or not. `child.start` carries the child's
`pid`, `started`, `args` including the program, and `working_directory`. `child.end` carries `pid`, `started`,
the exit `status` or the `signal` that ended it, and `elapsed`, in seconds.

The children that are not bound — `git`, `go test`, an agent's command-line client — are where
much of what fails actually fails, and they write no log of their own. Their parent's two lines are
the only record that they ran. The lines also connect a bound child to its nearest bound ancestor
across any processes in between.

### Message

A line a program's author chose to keep: its `text`, and nothing else.

> **A message is written on purpose.** A program writes a `message` where its author decided a line
> is worth keeping, and nowhere else. Nothing copies stderr, stdout, or a child's output into a
> log.

A captured stream would log whatever anything printed, at a volume nobody chose: progress lines,
a compiler echoing source, a tool repeating a person's words. None of that was reviewed as
something to keep. A line an author wrote on purpose says what it means, and is reviewed with the
code that writes it. What a stream carries belongs to whoever reads the stream ([the CLI
guide](cli-guide.md#output-modes)). A program that runs unattended writes a message for each
thing an operator would need to know afterwards.

> **A message has no level.** It is not `debug`, `info`, `warn` or `error`, and no switch turns a
> class of messages on or off.

A level is a second decision about a line that has already been judged worth writing, and it is
made by whoever wrote the line, before anyone knows what the reader will need. `debug` lines
written only when diagnostics are on are missing from exactly the run that failed. A line
someone might want belongs in the log; a line nobody would want is not written. And what a line
means is carried by its type and its text: a failure is a `process.end` with its status, a
`child.end` with its signal, or a type its owner defines, never a word a reader has to trust.

### Process end

The last line a process writes: the exit `status` it is about to return, its `elapsed` time in
seconds, and `unwritten`, the number of lines it could not write.

> **A process whose log has a start and no end did not end normally.** The store says which of two
> things is true: the process is still running, or it is gone and left no record. The missing
> line is itself the evidence, never a gap to fill in.

### Log pruned

Written by the process that removed files from the log home under [the bounds](#bounds): `files`
and `bytes` removed, and `oldest`, the earliest start instant among the files removed.

## Bounds

[Everything a program writes has a
ceiling](engineering-guide.md#everything-a-program-writes-has-a-ceiling). These are the ceilings
for logs.

> **A segment is at most 8 MiB. A process keeps its first segment and its last three.** When a
> process opens its fifth segment, it removes the oldest segment between the first and the new
> last three, and the new segment's [process start](#process-start) counts the removal in
> `dropped`.

The first segment is kept because it holds when a condition started and what the process was
started as. A long-lived process that malfunctions slowly has, by the time anyone looks, a tail
full of symptoms and a start that dates them.

> **A log home holds at most 1 GiB and 25 000 files, and no file last written more than 14 days
> ago.**

> **Writers hold the bound; the shipper is never needed for it.** Every bound process enforces the
> log home's bound as it starts, before writing its first line, and again each time it opens a new
> segment, in either case only when the last pass was an hour or more ago. The modification time
> of `.pruned` in the log home records the last pass.

A pass removes files in the order they were last written, oldest first, until the log home is
within its bound, and records what it removed as `log.pruned`. Files the store has not yet
acknowledged are removed like any others.

The bound is the writers' because the shipper is the program most likely to be missing when the
bound is needed. A log home fills fastest when nothing is shipping it: the shipper has died, the
store is down, or no shipper was ever installed on the machine. Tying the bound to the shipper would
leave it unbounded in exactly those cases. Running a pass as a new segment opens keeps the bound
held on a machine where nothing new starts and a long-running process keeps writing.

It never removes a file whose process may still be running. It lists the processes running on the
machine once per pass, and a file whose name carries a start instant that no running process
started at is not running. A match could be an unrelated process that started at the same
instant, which only delays that file's removal to a later pass, so the check needs no file to be
opened.

The log home is where logs pile up while the store cannot be reached, which is why its bound
matters most then. What the bound removes during an outage is exactly the loss the outage causes,
and `log.pruned` records it, instead of letting it show up later as an unexplained gap.

The bounds are not configurable. A limit each machine can raise is a limit machines raise until
there is none, and changing one of these numbers is an amendment to this document.

## Reading a log

> **A log is read with `logs`, a workspace tool.** It reads its machine's log home and changes
> nothing in it: it never removes, prunes or rewrites a file. Reading a log therefore never changes
> what the bounds or the shipper see.

It answers, on any machine and in the same way:

- **which processes the log home holds**, with each one's tool, arena, pid and start time, and
  whether it ended and how, is still running, or is gone without an end;
- **one process's log**, found by its pid, by its drawn id, or by its tool and a time;
- **a running process's log as it is written**, read the way a [shipper](#shipping) reads it, up
  to the last newline;
- **a process's ancestors and descendants**, resolved from each process's recorded `ancestors`
  and the [child lines](#child-start-and-child-end) around them.

Its output follows [the CLI guide](cli-guide.md#output-modes): rendered for a person at a
terminal, and the lines themselves, unchanged, when read by a program. Its command surface is
workspace's, in its `docs/tool-contract.md`.

A person looking for a log starts from something the file layout was never designed around: a
pid read from `ps`, an item, an arena, the governor's last hour. A directory listing answers only
the questions a file name was built for. One reader, installed with the tools every checkout
already has, answers all of those questions the same way on every machine. It also leaves the
layout free to serve the two parties that write and ship it.

## Shipping

> **Shipping is a role, and one program on a machine fills it.** The shipper ships its machine's
> log home. It holds `.shipper.lock` in the log home for as long as it ships, and a program that
> cannot take that lock does not ship.

A role, because what shipping requires is the rules in this section, not a particular program. A
program fills it by following them. The program that fills the role is stated as an end state and
changes by amendment:

- **On a host the governor runs on**, the governor's stage 1 is the shipper.
- **In a guest**, the runner serving the guest's arenas is the shipper.
- **A machine where no program fills the role ships nothing.** Its logs stay where they were
  written, within their bounds.

One shipper per log home, because two would send every line twice and race each other to remove
the same file. The lock sits in the log home it guards.

> **The lock is one the operating system releases when its holder dies.** It is an advisory lock on
> the file, held through an open handle, never a pid written into a file. A shipper that crashes
> leaves no lock behind, so the next program to fill the role takes it without judging whether the
> last holder is gone.

A dead shipper costs only shipping. Writers keep writing and keep the log home within its
[bound](#bounds). When a shipper returns, it resumes from the progress it recorded, within the 14
days below. A file the bound removed in the meantime shows up in the store as a gap in its process's
`sequence`, next to the `log.pruned` line that explains it.

> **Nothing stale is shipped.** Before its first delivery, a shipper runs a pass of the
> [bound](#bounds), whether or not the last pass was less than an hour ago. After that it never
> sends a line from a file last written more than 14 days ago. That holds for the file of a process
> still running too, which the bound keeps on the machine but the shipper leaves there.

A machine that has not run for six months comes back with a log home full of what happened before
it stopped. Nobody can still diagnose from any of it. Sending it would fill the store with months of
history at once, and keep for up to 90 more days data about people that the machine itself was about
to delete. What a log home keeps is 14 days, and that is the whole of what may leave it.

> **A shipper ships while the process writes.** A line reaches the store within a minute of being
> written, whenever the store can be reached, and a governor running for a month is read in the
> store as it runs, not when it exits.

The coordination is the file itself ([the line](#the-line)):

- **The shipper reads a file from where it last stopped to the last newline.** Bytes after the last
  newline are a line still being written. They are not read as a line.
- **The shipper keeps its progress in its own state, never in the logs it reads.** For each file
  that is the byte offset and the `sequence` it has had acknowledged. It changes nothing in the log
  home except to remove a file it is done with.
- **A segment is finished** when a later segment of the same process exists, when it holds a
  `process.end`, or when its process is observed gone. A partial line at the end of a finished
  segment is a torn write, and the shipper reports it and never sends it as a line.
- **A shipper opens a log so that its writer can still remove it**, so a file being read never
  holds up the writer's bound on a platform that locks open files.

> **Shipping never stands between a process and its log.** A writer never waits on a shipper, a
> pipe or the network. A store that is down, slow or unreachable changes nothing about what a
> process does or how long it takes.

> **Delivery is at least once.** The shipper sends whole lines, in `sequence` order for each process,
> and resends whatever was not acknowledged. The store identifies a line by `process` and `sequence` and
> keeps one copy of it.

> **A shipper removes a file once the file is finished and the store has acknowledged it to its
> last line.** Until then the file stays where it was written, within the log home's bound.

With the shipper keeping up, the log home holds little beyond the files of running processes and
whatever the store has not yet acknowledged. No other program moves a log, and no lifecycle of an
arena involves the shipper.

> **A shipper reports a process that ended without `process.end`.** When a finished segment has no
> end, and the process's recorded `pid` and `started` no longer run on its machine, the shipper
> tells the store once.

The shipper learns the store's address from the orchestrator its machine is registered with, and
that orchestrator either offers a store or offers none. An orchestrator without one leaves every log
on the machine that wrote it. The form of the request that carries lines to the store, and of the
acknowledgment that comes back, is base's.

## The store

> **The store answers which processes ran where, and what each one wrote, across the whole fleet.**
> It joins every line to its process's [process start](#process-start) by `process`, and every
> process to its parent. A line is then found by host, guest, arena, tool, version, item, step run,
> type and time, without reading any machine again. A process tree reads from the top down:
> a runner, the flow binary it started, the gate that binary ran, the test that gate ran.

> **The store reports what is missing as plainly as what is present:** a gap in a process's `sequence`,
> a torn write, a process gone without an end, and a host id running on two machines at once.

A store that shows only the lines it received makes a fleet with a broken shipper look healthy.

> **The store keeps a line for 90 days from the moment it receives it, and then deletes it.**

Ninety days covers a condition that develops slowly across releases. It is also a period the
development fleet's privacy notice can state as a number. [Everything a program writes has a
ceiling](engineering-guide.md#everything-a-program-writes-has-a-ceiling) holds for the store
as for any other sink. Any bound on its size is declared in its own specification, and that bound
only ever deletes sooner.

> **The projects' maintainers read the store. A contributor reads only the lines of processes that
> ran on hosts registered under their account, and on those hosts' guests and arenas.**

Which hosts are a contributor's is answered by the orchestrator's registration record at the
moment they read ([identity](identity.md#registration)). The store is never told, and never keeps
a copy.

> **The store keeps lines, not connections.** It records no source address of a delivery, and it
> holds ids, never accounts.

> **The store, and the orchestrator that keeps the registration record, run in the United States.**
> They run on hosts in a fleet that hosting manages, on one of the providers hosting supports, in a
> region within the United States. Hosting's host model gives each fleet one location, stated
> where the fleet is defined.

The country is a rule here, because it decides which law a contributor's data is held under and
what the development fleet's privacy notice promises. The provider and the region within the
country are not: they are facts of a deployment that can change without any rule here changing,
and naming them would give a second home to a fact the fleet's definition already holds.

> **Moving the store or the orchestrator to another country is an amendment to this document**,
> and it lands only after the development fleet's privacy notice states the new country.

> **Everything under a host can be exported or deleted.** At a maintainer's request the store
> exports, whole, every line of every process that ran on a host, its guests and its arenas. On a
> separate request it deletes the same set. A deletion leaves a record of what was deleted and
> when, holding no line. Either request covers the host's registration record too
> ([identity](identity.md#registration)).

A person who leaves the fleet can ask for their data and for its removal, and the host id is how
their machines are found. Deleting evidence does not conflict with [a log is
evidence](#a-log-is-evidence): nothing decides on a log, so removing one changes no decision.

> **Nothing in a log is published.** A log names hosts, arenas, paths and command lines, which are
> the disclosure categories of flow's `docs/disclosure.md`. A line quoted anywhere public passes
> the disclosure guard first, like any other text.
