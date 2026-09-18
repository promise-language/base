# CLI Tools

> **Tag:** `cli-guide` — remaining work to complete this document: the query named in
> `docs/index.md`.

> **Home:** [promise-language/org](https://github.com/promise-language/org) — this document
> changes here and nowhere else. To change it, file an issue against `org`.

How every command-line tool in the organization behaves at its invocation surface: how it takes
parameters, how it reports, and how it refuses. Its blockquotes are read as
[normative.md](normative.md#end-state) says: each is an invariant, and the prose under it is why.

## Scope

This document governs the invocation surface of every command-line tool an organization project
ships — the `bin/` tools and any binary a project's contributors or agents run by hand. It does not
own:

- **Which tools a project must have** — workspace's tool contract, its `docs/tool-contract.md`,
  owns that.
- **The gate envelope and the `--envelope` protocol** — base's gate contract, its
  `docs/gate-contract.md`, owns that.
- **How tools are built, provisioned, and kept fresh** — workspace's tooling document, its
  `docs/tooling.md`, owns that.

## Explicit inputs

> **A tool reads no environment variable to decide what it does.** Every parameter arrives as an
> explicit argument on the command line.

An environment variable is an argument with no audit trail: it is invisible in the invocation, it
leaks across process boundaries the caller never considered, and two invocations that look identical
behave differently. A reader of a command line must be able to know everything the tool was told.
The one thing a command line does not show is a configured default
([configuration](#configuration)), and a tool that has one can be asked for it.

Two narrow uses are sanctioned, and both are outside the tool's outcome:

- **Debugging the tool itself.** A variable that turns on the tool's own diagnostics may exist. It
  must never change what the tool does — only what it reports about its own execution.
- **Guard-enforced containment markers.** A variable a parent sets so that guards deny classes of
  action in the entire subprocess tree is containment, not argument transport: it is read by the
  guard, not by the tool, and it can only narrow what is possible, never select behaviour.

> **A parameter another program owns is named by its source, never read from its variable.** A
> credential, a profile, a provider's account: the tool takes a flag naming the source — the
> profile, the chain, the file — the platform that owns the variable resolves it, and the result
> names what it resolved to.

A provider's SDK reads its own variables; that is the provider's contract, and a tool that refused
it would be refusing the platform rather than protecting an audit trail. What the tool may not do
is let the variable decide silently. `-credentials production` on the command line says which
source was used, the identity in the result says what it became, and two invocations that read
alike acted alike or say why not. Where the source is a machine fact it is configurable
([configuration](#configuration)) — a default the command line can be asked for, never a variable
the command line cannot see. A variable the platform defines for locating its own fixed places —
the user's home, the executable path — is read by the platform's own call and decides nothing
about what the tool does.

## Flag form

> **A flag is a name, dash-separated when multiword, prefixed with `-` or `--`.** Its words are
> the engineering guide's ([naming](engineering-guide.md#naming)): full English words, the
> abbreviation dictionary's entries, and a quantity with its unit. The two prefixes are the same
> flag: tools normalize the prefix once, then match the name exactly.

`--my-long-flag` and `-my-long-flag` are one flag. `--myLongFlag`, `--my_long_flag`, and a
clipped word are not flags at all — they are unknown input ([fail closed](#fail-closed)).

> **A name — a flag's or a subcommand's — is lowercase ASCII letters `a`–`z` and digits `0`–`9`,
> with `-` as the only separator within it.** Nothing else: no uppercase, no underscores, no
> dots, no characters outside ASCII. The one character that may join names is the `:` of a
> compound subcommand, below.

Case is the cheapest way to mint an accidental alias — `--force` and `--Force` are one name to a
person and two to a matcher — and anything beyond lowercase ASCII is a name that types
differently across keyboards, shells, and platforms. A closed alphabet keeps every name exactly
as greppable, quotable, and portable as the one canonical spelling requires.

> **A subcommand may be compound: two or more names joined by `:`, `<concept>:<instance>` being
> the common shape.** The compound is one name — addressed exactly, declared in the tool's
> specification, never a prefix to search under.

A gate is `tested:root`: a concept every project shares and an instance only the project knows,
and a finer instance adds a segment. The instance is not a value — `tested -instance root` makes
every consumer rejoin what one name declared, one flow asked for, and one judge named — and it is
not a second path segment, which would leave `root` meaning nothing on its own. The alphabet
stays computable: a name, then any number of `:` each followed by a name; an empty segment is
not a name.

> **One name per flag. No aliases, no fallbacks.** There is exactly one canonical way to pass each
> parameter, and that way is the one help text, error messages, and documentation use.

A second name for the same thing is a second thing to search for, a second thing to deny in a
guard, and a fork in every transcript. The words a name may use are the engineering guide's, and
a flag has no exception of its own.

> **A value flag declares a type, and takes its value as the next argument or attached with
> `=`.** String, integer, duration, path, enumeration, or a list of one of those — boolean is not
> a value type — and the type is part of the flag's definition, printed by `-help`, and a value
> that does not satisfy it is a usage error naming the flag, the value, and the expected type.

`-timeout 30s` and `-timeout=30s` are one parameter, as the two prefixes are one flag: the tool
normalizes the spelling once and matches the value against the type.

> **A duration is one positive integer and one unit: `ms`, `s`, `m`, `h`, or `d`.** No fraction,
> no sign, no compound — `1h30m` is `90m` — and no other unit, on the command line, in the
> parameter file, and in any field a tool writes.

A type named without its grammar is one grammar per tool: `1h30m`, `90m`, `1.5h`, `5400s`, and
`PT1H30M` are each defensible, a reader of two tools' help cannot tell which each takes, and a
parameter file carrying a duration means different things to different readers of one field. One
grammar is one parser in every language and one value that travels a command line, a file, and a
wire field unchanged. The units are the ones a reader needs no knowledge of the tool to read.

> **A list is `list of <type>`: comma-separated on the command line, an array in the parameter
> file, every element checked as its type.** An empty element is a usage error.

A list declared as a string is a parser per tool: `-help` says `string`, which is true and useless;
one tool drops an empty element where the next refuses it and a third takes a trailing comma; and
the parameter file has nothing to map an array to. Naming the element type keeps the check where
the type is.

> **A boolean flag has a spelling for each state that differs from its default — `-my-flag`
> asserts it, `-no-my-flag` denies it — and a spelling that could only restate the default does
> not exist.** Neither takes a value; `-my-flag=true` and `-my-flag=false` are rejected with a
> pointer to the spelling that exists.

A boolean whose default is off has `-my-flag` alone; one whose default is on has `-no-my-flag`
alone; one whose default is decided at run time — from what the tool finds — has both.
`-no-dry-run` on a tool that does not dry-run by default is a second spelling of saying nothing:
it can be passed to no effect, it cannot appear in a transcript that means anything, and it is one
more name to match and to document for a behaviour that already happens. The denial is a word
rather than `=false` because a decision hidden inside a value is one a reader scanning for the
flag's name misreads — and it exists at all only when there is something to deny. A default that
changes gains or loses its spelling in the same change.

## One order

> **The command path comes first, complete and uninterrupted; every flag follows it; every
> positional argument follows the flags.**

```
tool <subcommand…> <flags…> <positional args…>
```

`tool sync -json origin` — never `tool -json sync`, and never `tool sync origin -json`. There is
no "global flag position": a flag that every command supports, `-json` included, is still written
after the command path, because it modifies the command being run, and until the path is complete
there is no such command. Tools where a flag works in one position and silently not in another
are the standing failure this rule exists to delete — two positions is two spellings of the same
invocation, and the second one is an alias.

`-help` and `-version` obey the same rule: `tool -help` is the root command's help — the full
surface, per [help and version](#help-and-version) — and `tool sync -help` is `sync`'s.

> **A flag appearing after the first positional argument is an error, never a positional.** A parser
> that stops at the first non-flag and hands the rest through untouched has silently reinterpreted
> the invocation; refusing it is the [fail-closed](#fail-closed) rule applied to position.

The operator who typed `tool sync origin -json` wanted JSON output; a tool that instead passes
`-json` to the backend as a name has done something no one asked, without a word.

> **`--` ends the flags.** Everything after it is positional, verbatim — and a positional that
> begins with `-` is accepted only after it.

Without the marker, a value like a file named `-report` is indistinguishable from a flag, and
guessing is worse than either answer. With it, the boundary between flags and arguments is
explicit exactly where it would otherwise be ambiguous.

## No general switches

> **No flag answers questions the tool has not asked yet.** Blanket switches — `-yes`, `-force`,
> "assume yes to everything" — are not permitted. **Overrides are named and independent: one
> flag per thing being overridden, never one flag meaning "ignore whatever refused".**

Every override names the one condition it overrides, so consenting to one risk never consents to
another — and the thing overridden is a *condition*, never a *phase*. A flag that skips a
verification step is named for one thing and switches off everything the step checks, including
whatever it grows to check later, in every script that already passes it: that is a flag
answering questions the tool has not asked yet, arriving under a name that looks specific. A tool
that would need `-yes` is a tool that asks questions interactively; it should instead refuse with
a typed, named condition and the specific flag that overrides it. An override that takes effect
is named in the result, so a transcript shows not only that consent was given but what it was
spent on.

## Output modes

> **Output is human-readable when stdout is a terminal and JSON when it is not.** Every tool
> supports `-json` and `-human` to force the mode regardless of piping. Passing both is a usage
> error.

The mode is decided by stdout only — never stderr, never an environment variable — and it is decided
for every command the tool has, `-help` and `-version` included ([help and
version](#help-and-version)). A rule with one exception has to be known, where a rule without one
can be assumed, and the caller most likely to forget the exception is the script that pipes every
command alike.

> **A command whose output a named contract fixes has one mode.** Its specification names the
> contract; on that command `-json` and `-human` are unknown input, and the refusal names the
> contract ([fail closed](#fail-closed)).

A gate prints one envelope; a judge prints one verdict. The shape is the contract's, not this
document's, and a caller that asked for it is reading it. Honouring `-human` there destroys what
the caller is parsing; accepting it and printing the envelope anyway is the silently dropped flag
this document forbids; so the flag is refused, and the refusal says which contract owns the
output. The script that pipes every command alike loses nothing — what it reads is the contract's
form either way — and the only reader the rule touches is the one who typed a mode the command
does not have.

> **Stdout carries the result and nothing else. Progress and narration go to stderr.**

That is what makes `tool > out.json` and `tool -json 2>/dev/null` both behave. JSON on stdout is a
stable interface: fields are added, never renamed or repurposed, and absent means unknown rather
than zero.

## Help and version

> **Every tool supports `-help`**: it prints the subcommands and, per subcommand, every flag with
> its type and a one-line description — or simply every flag, when the tool has no subcommands.
> It exits 0, and it is the only place the full flag list appears. A tool that is not fit to act
> answers neither this nor `-version` ([exit codes](#exit-codes)).

A flag whose default is configurable ([configuration](#configuration)) also shows the value in force
and the file it came from. In JSON, `-help` is the surface as data — the same facts, in one object
every tool shares:

```json
{
  "project": "sometool",
  "description": "…",
  "flags": [{"name": "json", "type": "boolean", "description": "…"}],
  "commands": [
    {
      "name": "sync",
      "description": "…",
      "flags": [{"name": "timeout", "type": "duration", "description": "…"}]
    }
  ]
}
```

The top-level `description` and `flags` are the root command's — the whole tool's, when it has no
subcommands, and then `commands` is empty. A flag's `type` is the one [flag form](#flag-form) has it
declare — a list adds `element`, the type of its elements — and a configurable flag adds `value`
and `source`.

> **A command set the tool does not author is described, never listed.** Where the commands are
> the project's — the gates it answers — `-help` describes them as one class, names the
> invocation that enumerates them, and lists none of them.

An enumeration in two places is one fact with two homes, free to drift: `bin/gate --list` is the
gates' one home, and help says so. In JSON the class is one entry under `commands` whose `name` is
the pattern, `<gate>`, and whose `enumerated-by` is the invocation. A tool knows at definition
time which of its commands it authors, so a checker can tell the class from a tool that left
commands out.

> **Every tool supports `-version`**: it identifies the binary as either a comprehensible release
> version or the commit hash it was built from. It exits 0.

A binary that cannot say what it is cannot be the subject of a bug report — and `-version` is the
command most likely to be read by a machine, since it is how one program identifies another. In
human mode it is one line, `<project> <text>`. In JSON it is one object, so no consumer parses a
version out of prose and two tools never disagree on where the fields are:

```json
{
  "project": "sometool",
  "text": "v0.12.0",
  "major": 0,
  "minor": 12,
  "patch": 0
}
```

- **`project`** — the project the binary is built from, as its repository is named. Never folded
  into another field.
- **`text`** — the version as the human line shows it, without the project name: a release
  version, or the commit hash.
- **`major`, `minor`, `patch`** — present when the version is semantic, so no consumer parses
  `text`; `prerelease` and `build` join them when the version carries those parts. A version that is
  not semantic omits them rather than zeroing them: absent means unknown ([output
  modes](#output-modes)), and a `0` would be compared as a number.
- **`commit`** — the hash the binary was built from, when known.

> **`-help` and `-version` do nothing else.** An invocation carrying either prints and exits:
> nothing is done, written, or contacted. Alongside them the tool accepts the `-json` and `-human`
> of [output modes](#output-modes) and nothing more; any other flag, or a positional, is a usage
> error ([fail closed](#fail-closed)).

They are the two flags an operator types to learn what a tool is before running it, and a help
request that runs the gate, or a version query that rewires the hooks, has done work nobody
asked for under the one flag that promised none. Whatever else is on the line is refused rather
than ignored: the operator who typed `tool sync -help origin` may have meant the help and may
have meant the command, and guessing is worse than either answer.

> **A tool whose every action is a subcommand orients when invoked bare.** `tool` alone prints
> its `-version` line, then a brief help: its most common commands with their one-line
> descriptions, and how to reach the full surface (`tool -help`). It exits 0 and does nothing
> else. In JSON it is the `-version` object with `commands` added — the entries `-help` gives,
> for the common commands only.

A bare `bin/issue` is a person finding their footing. The version comes first because "which build
is this" is the first thing a bug report needs; the brief help answers "what do I type next" without
the wall of definitions [fail closed](#fail-closed) refuses to dump, and the full list stays one
flag away. Which commands are common is the tool's own specification's choice, and the list is short
by design. A tool with a root action has no such courtesy — invoked bare, it runs.

## Fail closed

> **An unknown flag is an error that names it** — the bad flag, the closest existing flag when one
> is close, and how to get the supported list (`-help`). Nothing is silently ignored.

The error does not dump the flag list: a wall of definitions buries the one fact the operator
needs. The pointer to `-help` is the list, one step away.

> **Invocation errors are rejected before any action.** The tool exits 2 having claimed nothing,
> written nothing, and contacted nothing.

The same rule covers unknown subcommands, missing required parameters, values failing their type,
misplaced flags ([one order](#one-order)), and contradictory parameters. Validation is exhaustive:
every problem with the invocation is reported, not just the first.

## Invocation from a file

> **`-json-input /path/to/args.json` supplies parameters from a JSON file whose keys map exactly
> to the tool's flags**, plus `"args"`, an array carrying the positional arguments.
> An unknown key in the file is an unknown flag ([fail closed](#fail-closed)).

The file is a transport for the same closed parameter set, not a second configuration system: no key
exists in the file that does not exist as a flag, and each is spelled exactly as the flag is. A
boolean flag appears under its own name with the value `true` — `"dry-run": true`, or `"no-fetch":
true` — and `false` is a usage error, because a spelling that says nothing does not exist on the
command line either ([flag form](#flag-form)). The name is deliberately not `-file` or `-input` —
bare words a tool's own domain will want for its actual inputs; `-json-input` names the mechanism,
so it collides with nothing a tool processes.

> **A parameter set both in the file and on the command line is a usage error.** There is no
> precedence between the two, because precedence is a fallback ([flag form](#flag-form)).

## Subcommands

> **The command set is closed in both directions**: no command is added outside the tool's
> definition, and no command answers to a name not in the set. Each command has exactly one name.

A set the tool does not author — the gates a project answers — is closed by the definition that
lists it, and the tool's specification names that definition. Its members are compound names
([flag form](#flag-form)), and `bin/gate tested:root --envelope` is the [one order](#one-order):
the path `gate tested:root`, complete, then the flag the runner appends.

> **Addressing is exact.** An identifier the user types resolves to exactly what it names, never to
> something that merely contains or resembles it.

## Exit codes

> **`0` — did what was asked**, including when there was nothing to do. An empty result is not an
> error. **`1` — could not complete**, or stopped on a condition a human must clear. **`2` — the
> invocation itself was malformed**, and nothing was done. **`3` — the tool refused to act**,
> because it is not fit to: built from source other than the tree beside it, not built by the
> project's builder, its own precondition unmet. Nothing was done, and the refusal names the
> condition and what clears it.

The four answer two questions. Was the subject examined: `0` and `1` say yes, `2` and `3` say no.
Whose repair is it: `1` the subject's, `2` the invocation's, `3` the installation's. A stale tool
that exits `1` reports the subject as bad, and its caller names a repair that is not the repair —
the one thing the message will not say is "rebuild the tools".

> **A tool that is not fit to act refuses every invocation**, `-help` and `-version` included,
> before it reads the command line. The one tool that never refuses on this ground is the builder
> the refusal names: it needs nothing pre-built, so the way out is always open.

Help and version are not an exemption here but the reason for the rule. What a stale binary would
print is the surface it was built with, and that is exactly what is out of date: it would describe
flags that have moved and name a version that is not what the tree holds, confidently and in the
one place an operator goes to learn what a tool is. A binary that cannot be trusted to measure
cannot be trusted to describe itself either. Refusing before the command line is read is what
keeps the answer from depending on which flag was typed.

> **A refusal is the result.** It goes to stdout in the mode in force ([output
> modes](#output-modes)): in JSON one object, `{"condition": "…", "recovery": "…"}`, naming the
> condition and what clears it; in human mode one line saying the same. A caller reads the status
> to know it received no answer, and the object to know why.

A sentinel at the head of stderr is prose, not a protocol: the first rewording breaks every
matcher, and the matchers live in other repositories than the tool.

## Configuration

> **A tool has no configuration file.** The answer to "should this be configurable" is no, and a
> tool that reads a file anyway names in its own specification each key it reads and why. A tool
> reading a key its specification does not name has a defect, not a convention.

A configuration file is the [explicit inputs](#explicit-inputs) hazard in another container: an
input the command line does not show, so two invocations that look identical behave differently.
What lets a tool carry one at all is that the specification naming its keys is a reviewed document —
so the set of configured tools and of configurable keys is knowable, and the decision to make
something configurable is taken where decisions are reviewed, not where a parameter got tedious to
type.

> **Only a machine fact is configurable.** A value belongs in configuration when it is a property
> of the *machine* — true for every invocation on that host, and different on the next — never a
> property of the invocation.

The root a fleet of checkouts sits under is a machine fact; a timeout, a target, a mode is not.
This is the test that keeps configuration from becoming a second spelling of the flag set, which
[invocation from a file](#invocation-from-a-file) already refuses.

> **A configured value is a default, not an argument.** It stands where the built-in default
> stood, and a flag overrides it exactly as a flag overrides any default. Every configurable
> value has a flag; nothing is reachable through the file alone.

This is not the precedence [invocation from a file](#invocation-from-a-file) refuses. The file and
the command line are not two sources of one argument; they are a default and an override, the two
things every flag with a default already has.

> **One location, one form.** Configuration lives at `~/.config/<tool>/config.json` — on Windows,
> under the user's roaming application-data folder as `<tool>\config.json` — and its keys are the
> flags they default, spelled as [invocation from a file](#invocation-from-a-file) spells them. A
> cache lives under `~/.cache/<tool>/` and state under `~/.local/state/<tool>/` — on Windows, under
> the user's local application-data folder as `<tool>\cache\` and `<tool>\state\`; neither is
> configuration.

The paths are the XDG defaults on every platform but Windows, because a command-line tool is used
from a shell and that is where its neighbours keep theirs. The `XDG_*` variables that would relocate
them are not read ([explicit inputs](#explicit-inputs)): the location is the same for every
invocation on the host, which is what makes it a machine fact rather than an ambient one. A file
that does not parse, or names a key that is not a configurable flag, is refused before any action
([fail closed](#fail-closed)).

A **released product** whose layout this rule does not fit — a language toolchain with module
caches, build outputs, and per-project state is more than one directory of each — takes its own
layout under one condition: a specification in that project names every location the product
writes, what each holds, and which of the three kinds it is. The exception is the specification,
not the product; a location no specification names lands at the default above.

> **A configured tool can be asked what it was told.** `-help` shows, for every configurable
> flag, the value in force and the file it came from ([help and version](#help-and-version)).

That is what keeps the file from being the invisible input [explicit inputs](#explicit-inputs)
refuses: the command line does not show it, but one command does.
