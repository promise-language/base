# Identity

> **Tag:** `identity` — remaining work to complete this document: the query named in
> `docs/index.md`.

> **Home:** [promise-language/org](https://github.com/promise-language/org) — this document
> changes here and nowhere else. To change it, file an issue against `org`.

> **This document identifies the machines and workspaces the Promise language projects are
> developed on, and nothing else.** It covers the development tools and the orchestration system
> that the projects' contributors and maintainers use. It does not apply to the Promise compiler,
> to programs built with it, or to any other product a project ships.

> **An identity here names a machine or something running on one, and never a person:**
>
> - a **host**: a physical or virtual machine;
> - a **guest**: a container or virtual machine running on a host. It is a machine, not a
>   visitor;
> - an **arena**: a checkout of a project that work is done in;
> - a **tool**: one build of one program;
> - a **process**: one run of a tool.

An identity here is seen by more people than its owner. The projects' maintainers see every host,
guest and arena in the fleet. A contributor sees only their own: the hosts registered under their
account, and what runs there. So no identity holds anything personal:

- **An id is random**, and says nothing about a machine or who owns it.
- **A host's label is approved by the person whose machine it is**, before anyone else sees it
  ([human labels](#human-labels)). A machine's own host name is often its owner's —
  `janes-macbook-pro` — and it becomes the label only if its owner chooses it.
- **An account name, an email address, a home directory, a hardware serial or an address is never
  part of any identity**, and nothing here derives one from them.

What names the place work runs in and the thing that runs it: the host, the guest, the arena, the
tool and the process. Each is named by an id that is unique without anyone coordinating it. A
host, a guest and an arena also carry a human label: a short, meaningful handle a person can read
where the id alone is random digits. Each is created one way, by one party, at one moment.
[Logging](logging.md) is the first reader — a log line nobody can place is a line nobody can use —
but leases, exclusion scopes and attribution name the same five things, and name them with the
same ids.

## Scope

This document owns what a host, a guest, an arena, a tool and a process are; the form of each
one's id and human label; who creates each, where each is kept, and what each survives. It binds
every program [logging](logging.md#scope) binds, and every document that names one of the five.

It does not own the identities of work — a project, an item, a step, a step run, a claim, an
account — which are the orchestrators' and base's wire types'. It does not own how a host or an
arena is registered with an orchestrator, or what makes a registration trusted: that is the
credential the orchestrator issues, and [registration](#registration) says why the id here is not
it. It does not own what a log line carries, which is [logging](logging.md)'s.

## The five identities

| Identity | What it names | Created by | Kept in | Survives | Ends with |
|---|---|---|---|---|---|
| **Host** | a physical or virtual machine, with the one account on it that runs the development tools | provisioning, with a person confirming its label | the host home's record | reboots, relabelling, reprovisioning, every tool upgrade | the host home's removal |
| **Guest** | a container or virtual machine on a host, which arenas live in | whatever creates the guest | the guest's own host home | the guest's restarts | the guest's destruction |
| **Arena** | a checkout units of work are leased to | provisioning, the first time it provisions the checkout | `.workspace/arena.json` in the checkout | reprovisioning, moving the directory, every clear of `.home/` | the checkout's removal |
| **Tool** | one build of one program, as it was invoked | the build | the binary | — | — |
| **Process** | one execution of a tool | the process, as it starts | its own memory, and its log | nothing | its exit |

> **One operating-system account on a machine runs the development tools and the orchestration
> system, and the host home is that account's.** Provisioning refuses to create a host identity
> on a machine where it can see those tools running under another account.

One account keeps one machine one host. A second account's tools could not see the first account's
host home, so they would create a second identity for the same machine, keep a second log home and
need a second shipper. And the account is how a contributor's own hosts are told apart from
everyone else's: two accounts on one machine would let one person read another's lines through the
host they share.

**A container or virtual machine whose host nothing here can see is a host, not a guest.** A cloud
platform's container is one: a guest is recorded against its parent, and a parent nobody can name
is not one.

**A fresh clone is a new arena**, even of the same project into the same directory. The arena is
the checkout, not the path it sits at, and a new clone starts no unit of work where the last one
stopped.

## Ids

> **An id is 128 random bits from the operating system's secure random source, written as a kind
> letter, a dash, and 32 lowercase hexadecimal digits.** The kinds are `h` a host, `g` a guest, `a`
> an arena and `p` a process: `h-9f86d081884c7d659a2feaa0c55ad015`.

Random, because every other source fails somewhere this system runs. A derived id, from a host
name, a machine id, a hardware address or a path, is copied perfectly by a cloned image, collides
across hosts that share a naming habit, and discloses what it was derived from. A minted id needs
a server to mint it, and a laptop running `bin/verify`, a standalone checkout and a cloud
container that has not registered yet have none. Randomness needs nothing but the machine it runs
on, and it gives uniqueness to parties who never talk to each other.

128 bits, because processes are counted across the whole fleet for as long as the store keeps
them. At 64 bits two processes sharing an id stop being a curiosity once there are some billions
of them; at 128 bits they never meet. Nothing is gained past that, and every log line carries one.

The kind letter is there so an id read out of context says what it names. An arena's id in a host
field is then a malformed value that is refused, not a lookup that quietly finds nothing, which is
[identities are types](engineering-guide.md#types-and-shape) applied to the text form.

> **An id never changes and is never reused.** Relabelling, moving and reprovisioning keep it, and
> nothing creates a second one for a thing that has one. A thing whose record is lost is a new
> thing.

> **An id means nothing.** No part of it is derived from what it names, so nothing can be learned
> from one and nothing in one goes stale.

> **An identity record is never part of an image.** A host image, a guest image, a snapshot kept
> for cloning and a template are all built without one, and the record is created where the
> machine first runs.

An image that carries a record makes every machine started from it the same host. Logs from all of
them join into one account of one machine that never existed, and the file whose only job is to
tell machines apart is what merged them. [Logging](logging.md#process-start) records the boot
time of every process, so the store can see one record running on two machines and report it.

## Human labels

> **A human label is for a person to read, never to compare.** It can be edited at any time, it is
> not unique, and nothing is addressed by it. Every record that names something carries the id,
> and the label only beside it for display.

A label is not a person's name. It is a short handle for a host, a guest or an arena, so that a view
of the fleet reads `sofia` and `forge_1` where it would otherwise read two strings of 32 random
hexadecimal digits.

> **A host's, a guest's or an arena's label is 1 to 64 characters from lowercase ASCII letters,
> digits and `_`.**

A derivation lowercases what it derives from, replaces any other character with `_`, and cuts the
result to 64. A label a person gives that is not in this form is refused, and the refusal states the
form. The label is not changed silently into something the person did not type.

Lowercase only, because a label is read, typed and searched for by people, and `Sofia`, `sofia` and
`SOFIA` would be three spellings of one machine to every reader but a case-insensitive filesystem.
One case removes the question.

`_` is the only punctuation a label may hold, so that `-`, `.` and `/` stay free to separate labels
wherever they are joined. `sofia/arena_130/forge_1` splits back into a host, a guest and an arena
without escaping or guessing, because none of the three can contain the character that joins
them. A label built from letters, digits and `_` is also one word to every terminal, editor and
search, so a double-click selects all of it.

A tool has no label. It is known by its executable's own name, fixed by the project that builds it
([tools](#tools)), which is not in this form and is not held to it.

| Identity | Human label, as derived | Confirmed |
|---|---|---|
| **Host** | the machine's host name up to its first dot, lowercased | by a person, when the host is created |
| **Guest** | the name its creator gave the container or virtual machine | no |
| **Arena** | the checkout directory's own name, prefixed with the repository name of the checkout's origin and `_` unless the directory's name already begins with it: `forge_1`, `promise_mac_3` | no |
| **Tool** | none: a tool is known by the bare name it was invoked as | — |
| **Process** | none: a process is known by its tool and its id | — |

> **A host's label is confirmed by the person who creates the host.** At a terminal, provisioning
> shows the derived label and asks for it to be accepted or replaced. Without a terminal, the label
> is given on the command line by whoever creates the host on a person's behalf. A host whose label
> was neither confirmed nor given is not created: provisioning refuses and names the flag that
> supplies it.

Only a host's label is asked for. A host's derived label comes from whatever a DNS server, a cloud
vendor or an installer chose — `localhost`, `ip-10-0-3-17`, `DESKTOP-7Q2M4KD` — or from its owner,
as `janes-macbook-pro`. It is the label the maintainers see, so the person whose machine it is decides
it before anyone else sees it. It is also the label a person scans for first in a view of
the whole fleet. A host is created once per machine, so asking
costs one question per machine. A guest's label and an arena's are taken from choices someone
already made — the container's name, the directory the clone went into — so they already mean
something. They are also created by the dozen, unattended, where a question would have nobody to
answer it.

Relabelling changes the record's label and nothing else. Lines already written keep the label that was
current when they were written, and the id joins them to the lines written after.

## Where records are kept

> **The host home is one directory per machine, at one path per platform.** It holds the machine's
> identity record and its [log home](logging.md#where-a-log-is-written), and nothing that belongs
> to a single tool.

| Platform | Host home |
|---|---|
| Windows | `promise-language\state\` in the account's local application-data folder |
| Every other platform, macOS and Linux included | `.local/state/promise-language/` in the account's home directory |

The account's home directory and its local application-data folder are the ones the platform's
own call reports. No variable relocates the host home, `XDG_STATE_HOME` included: it is the same
place for every program on the machine, which is what lets a program that writes there and a
program that reads there find each other without being told. An account the platform reports no
home for has no host home, and a program running under it treats every record as missing.

These are the locations [the CLI guide](cli-guide.md#configuration) gives a tool's state,
with `promise-language` where a tool's name would be: the host home is state, but the host's and
not any one tool's. A tool's own state stays where the CLI guide puts it.

The host home's record is `identity.json`. On a host:

```json
{"kind": "host", "id": "h-9f86d081884c7d659a2feaa0c55ad015", "label": "sofia"}
```

In a guest it names the guest, and the host the guest runs on:

```json
{
  "kind": "guest",
  "id": "g-2c26b46b68ffc68ff99b453c1d304134",
  "label": "arena_130",
  "host": {"id": "h-9f86d081884c7d659a2feaa0c55ad015", "label": "sofia"}
}
```

The host's label is copied at the moment the guest is created, for display, and can fall behind a
relabelling. The id is what joins the two records.

An arena's record is `.workspace/arena.json` in its checkout:

```json
{"id": "a-7d865e959b2466918c9863afca942d0f", "label": "forge_1"}
```

> **A record is written only by its creator, and read by everything.** Every program may read the
> host home's record and its own checkout's arena record. Only provisioning writes a host's or an
> arena's record, and only a guest's creator writes a guest's. A program that finds no record
> reports the thing as unknown, and never creates an id for it.

> **A project or workspace tool running in an arena writes nothing outside its checkout.** Tools
> run concurrently in many arenas on one machine, and two of them writing the same path outside
> their checkouts overwrite each other's work with neither able to tell (workspace's
> `docs/tooling.md`, *An agent writes only where it is working*).

> **Infrastructure those tools share may write to a shared location outside every checkout, and
> only where all three hold:**
>
> - **a binding norm defines it**: the location, its layout, and every write made to it;
> - **concurrent writes are safe by construction**: no two writers ever write the same file, whatever
>   else runs at the same moment;
> - **one reusable library per language implements the writing**, and every program that writes
>   there uses it. Writing there any other way is a defect.
>
> The [log home](logging.md#where-a-log-is-written) is such a location. The host home's record is
> not: provisioning alone writes it.

So a tool in an arena could not create a host record even if it would. And a program that created
an id wherever it found none would give the host a second identity, which the next such program
would replace with a third.

## Creating an identity

Every creator runs the same procedure, implemented once:

1. **Read the record.** If it exists and parses, stop: the thing already has an identity, so
   creating one is idempotent. Relabelling is a separate act.
2. **A record that exists but does not parse is refused, never replaced.** Overwriting it would
   give the thing a second identity while lines written under the first are still on their way to
   the store. The refusal names the file.
3. **Draw the id.**
4. **Derive the label**, and for a host, have it confirmed.
5. **Write the record** to a temporary file in the same directory, then rename it into place, so
   no reader ever sees half a record.

Who runs it:

- **A host's identity** is created by workspace's provisioning the first time it acts on the host.
  Hosting, for a host it creates, supplies the label through the machine's first boot, where
  provisioning creates the record.
- **A guest's identity** is created by the guest's creator — the arena provider that starts the
  container or virtual machine, or hosting. Only the creator knows the parent, so it writes the
  record into the guest before any arena there runs, copying the parent from its own host home's
  record.
- **An arena's identity** is created by workspace's provisioning when it first provisions the
  checkout.

## Tools

> **A tool is named by the bare name it runs as, and by its `-version` object.**

The bare name is the name it dispatched as — `issue`, never `bin/issue` and never
`workspace-darwin-arm64` — because one image installed under six names is six different tools to
a person reading what ran. The version is the object [the CLI guide](cli-guide.md#help-and-version)
has every tool print, taken whole: nothing here adds a second spelling of a version.

## Processes

> **A process draws its id as it starts, before it writes anything.**

> **A process is also named by what the operating system observes: its pid and the moment it
> started.** The drawn id joins what a process wrote. The observed pair is what can be checked
> against a running machine, to learn whether that process is still alive: a pid alone is reused,
> and a drawn id cannot be asked of the kernel.

Both are needed because they answer different questions. The drawn id is unique across the fleet
and is known from the process's first line, but only the process's own word vouches for it. The
observed pair can be verified at any moment on its host, but it is unique only there and only
until the host reboots. A lease holder is the observed pair on its host.

> **A process records its ancestry as the operating system observes it.** As it starts, it reads
> its parent's pid and that process's start time, then the same for the parent's parent, and so on.
> It stops only when the chain ends: after eight levels, at the machine's first process, at an
> ancestor it cannot read, or at an ancestor that started after the process below it, which is a
> pid already reused. It records every ancestor it read, nearest first, on its [process
> start](logging.md#process-start).

The chain is recorded whole, never cut short at the first ancestor that has a log of its own.
Stopping there would mean reading the log home before writing a first line, so what a process
records would depend on other writers and on what the bound had already removed. An observation
should say the same thing whatever else is running. The cost is about 45 bytes an ancestor, once
per segment. The top of the chain is worth having too: it says whether a process ran under a
runner, an agent's session, a terminal or the system's service manager, even when nothing up there
keeps a log.

> **A process's parent is its nearest ancestor that left a record.** A reader of the log — `logs`,
> or the store — takes the first recorded ancestor that either has a process start of its own or
> was recorded starting as another process's [child](logging.md#child-start-and-child-end), on the
> same machine. The process that wrote that line is the parent.

Nothing is passed from parent to child, so nothing has to cross the processes in between that know
nothing about lineage. Nothing can arrive stale or wrong either. A value handed down can be exported
by a shell, copied into an unrelated process, or inherited by a process that was never the sender's
child. An observation is what the operating system says, and nothing between the two processes can
change it. It also keeps an environment variable out of every tool's inputs, which is where [the
CLI guide](cli-guide.md#explicit-inputs) wants none.

A chain rather than only the parent, because two processes that record nothing can sit in a row:
`bin/verify` runs `go test`, which runs a test binary, which runs `bin/gate`. `bin/verify` records
`go test` as its child, but nothing records the test binary. From `bin/gate`, only the chain reaches
back to a process that left a record.

An ancestor that exited before its descendant started cannot be observed, because the operating
system has already given the orphan another parent, so lineage stops there. Lineage matters for work
one process waits on — a runner and its flow binary, a flow binary and its gate, the governor's two
stages — and in those the ancestor is alive when the descendant records it.

## Registration

> **Registering a host, a guest or an arena with an orchestrator presents the id this document
> defines. The orchestrator does not mint another.**

Whether a registration is trusted is a separate question, and the orchestrator answers it with a
credential it issues and rotates, never with the id. An id is a name, and a name every program on
the machine can read is not a secret. A copied record is caught where a copied credential is
caught, by the orchestrator refusing a credential it has superseded. The id stays what every log
line, lease and exclusion on that host is joined by, whether or not any orchestrator ever sees it.

> **Registration is where an identity first leaves its machine.** Creating an identity sends
> nothing anywhere. Registering joins a host's id, and its guests' and arenas' ids, to the account
> that registered them. That join is kept in the orchestrator's registration record, which the
> projects' maintainers read, and only they do.

The join is what turns a random id into a named person's machine, so it is kept in one record and
copied nowhere. The store holds ids, never accounts ([logging](logging.md#the-store)).

> **A registration record is discarded once its host has not phoned home for 90 days.** A host
> phones home with any request made with its credential: registering, renewing a session token, or
> shipping. The joins of the host's guests and arenas are discarded with it.

Ninety days is how long the store keeps a line, so a registration record never outlives the last
line it could place. A host that last phoned home 90 days ago has shipped nothing since.

> **A request to export or delete everything under a host covers its registration record**, along
> with the lines the store keeps ([logging](logging.md#the-store)).

> **A person is shown where the development fleet's privacy notice is at the two moments a person
> is present for this:** when provisioning asks them to confirm a host's label, and when a machine
> is first registered with an orchestrator that offers a store.
