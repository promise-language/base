# Normative Documents

> **Tag:** `normative` — remaining work to complete this document: the query named in
> `docs/index.md`.

> **Home:** [promise-language/org](https://github.com/promise-language/org) — this document
> changes here and nowhere else. To change it, file an issue against `org`.

What makes a document in a managed project binding, what one must contain, and the rules that
keep two of them from ever disagreeing. Every managed project holds the same documentation
structure, so a reader who has learned one tree has learned them all. This document's subject is
the documents themselves; it is one of them, and every rule below applies to it.

## Location

There is no configuration file, and no marker decides what a file is. The directory a file sits
in determines it, and the [header](#header) line only restates what the directory already
decided:

| Location | What a file there is | Binding? |
|----------|----------------------|----------|
| `docs/` root | A **specification**: what the project *should* be — the intended end state. | **Yes.** Work that contradicts one stops and is resolved — the document amended, the item adjusted, or the item rejected — never shipped as a quiet deviation. |
| `docs/org/` | The **organization-wide corpus**, vendored from its home repository. | **Yes.** Never edited in place: it changes at its home, and reaches the project by sync. |
| `docs/proposals/` | An end state that has **not been ratified** — a draft, an RFC, a direction under discussion. | No. |
| `docs/archive/` | An end state that has been **superseded or delivered** — kept for history. | No. |
| `docs/research/` | Background analysis feeding a decision — an assessment, not a design. | No. |

Those five rows are the whole vocabulary. A project may lack `research/` or have an empty
`archive/`; it may not invent a sixth location or assign one of these a different meaning. The
corpus's home repository is the one tree with no `docs/org/`: its root *is* the corpus, and it
is held to every rule here exactly as a project is.

**`docs/index.md` is the map and the one file in the root that is not a specification.** Every
tracked file under `docs/` is listed in it, wherever it lives — the section an entry sits under
is where its binding status is written down. `docs/org/` is listed once, as the directory, by
way of the **stamp** it carries: the stamp names the release the copies came from and every
member, so the corpus's map is the corpus's own, and a sync adds and removes members without
editing a file the project owns. The index also names the project's status query
([header](#header)), which is the one per-project fact this shared document cannot carry.

## Header

Every Markdown document under `docs/` — the index excepted — opens with its title and, beneath
it, the header lines its location requires. The location decides; the lines restate the decision
for the reader who arrived by a link rather than by the directory, and they are checked against
the location ([mechanical checks](#mechanical-checks)) so the two cannot drift. One form per
location — the root, `proposals/`, `archive/`, `research/`, in that order:

```markdown
> **Tag:** `<basename>` — remaining work to complete this document: the query named in `docs/index.md`.
> **Proposal.** Not normative: an end state under discussion, binding nothing until ratified.
> **Archived.** Not normative: superseded or delivered, kept for history.
> **Research.** Not normative: an assessment feeding a decision, not a design.
```

The tag is always the file's basename minus `.md`, so the vocabulary is the directory listing
and nothing duplicates it. The query is the project's own — a `gh issue list --label`
invocation on a GitHub-tracked project — and its exact spelling is stated once, in
`docs/index.md`, which is why the tag line points there instead of carrying it: a query repeated
in every header is a copy in every header.

**A proposal carries both lines**, the proposal line first and the tag line under it. A document
that describes an intended end state owes work toward that end state — the rules it has not yet
decided ([end state](#end-state)) — whether or not it binds anything yet, and the tag is what
lists that work. The label is made when the proposal is, so the decisions taken while it was a
draft stay under the same query after [ratification](#lifecycle), and ratification stops being
the moment a document first becomes something one can ask a question about. `archive/` and
`research/` carry no tag: one has no remaining work by definition, and the other is an
assessment rather than an end state to complete. A document under `docs/org/` carries the header
it has at its home — the tag line and the home line — byte for byte, like the rest of it.

A document that other repositories are held to — every specification in the corpus's home —
also carries, after its tag line, its **home line**: the repository it changes in, and how to
ask for the change. A reader who meets the document anywhere but its home is the reader that
line exists for ([reconciliation](#reconciliation)), so it is written in the source, worded to
be true wherever the file is read, and never added by whatever delivers the file. One form:

```markdown
> **Home:** [<owner>/<repository>](<repository url>) — this document changes here and nowhere else. To change it, file an issue against `<repository>`.
```

## Sections

> **A section is addressed by its slug, never by a number.** Headings below the title are
> unnumbered. A heading is a short name in letters, digits, and spaces, unique within its
> document, and its slug — the heading lowercased, spaces as dashes — is the section's identity.
> A reference, in a document, an item, or a comment, is `<file>#<slug>`:
> `normative.md#reconciliation`.

A number is a position, and a position moves. Inserting a section renumbers everything after it,
and every reference written against the old numbering still resolves — to the wrong section,
silently. A slug moves only when its heading is renamed, and then every reference breaks loudly:
in the tree the link check names each one, and outside it the old slug is a string to search for.
Renaming a heading is therefore an amendment, and the change that renames it repairs every
reference in the tree.

The alphabet is closed for the same reason a flag's is: the renderer that makes an anchor strips
punctuation, so a heading carrying any has a slug a reader cannot predict from the text. With
letters, digits, and spaces alone, the slug is computable by a reader and by the check alike.
Inside a document a reference is a link, and its text is the heading or the phrase the sentence
needs; a bare `#slug` is never shown to a reader as prose.

## End state

**A specification describes what the project should be, never how far along it is.** No status
sections, no progress notes, no phasing, no "currently", "not yet", or "implemented", and no
inline markers naming an item — a status section arriving one sentence at a time.

The practical test: **a specification reads identically the day before and the day after the
work that implements it.** A sentence that would have to change when an item closes is status,
and does not belong.

Where status lives instead, and all three are queries or single homes rather than prose in a
specification:

- **Per document**: the open items carrying the document's tag *are* its status section, always
  current.
- **For the project**: the README's Status section is the one sanctioned home for "how far along
  is this" prose. Nothing else records progress.
- **Undecided design**: a question whose answer would change what the specification says is an
  item carrying the document's tag, like every other piece of remaining work on it
  ([reconciliation](#reconciliation)) — never a section of prose. A specification that says a
  thing is undecided is read once and built around, and nothing ever comes back to it; an item is
  listed by the query the tag line points at, and somebody owns pushing it. So a document plus
  its open items is the end state, exactly as an implementation plus its gap items is. Undecided
  *time* is not a question at all — work merely not done yet is an ordinary item under the tag.
  This holds for a proposal exactly as for a specification: it is an end state under discussion,
  and what is still being discussed is its remaining work, under its own tag
  ([header](#header)).

**A rule stated as a blockquote is an invariant**, and the prose under it is why. If the rule
and the reasoning ever disagree, the rule is what the implementation must satisfy.

## One home

**A fact is specified in exactly one document.** Two specifications must never define the same
thing, and no specification may claim authority over another. *Supersedes*, *takes precedence
over*, *overrides*, *this document wins*, *the authoritative version is* — a document that needs
one of these is proof that a fact has two homes, and the remedy is always to delete the
duplicate and cross-reference, never to rank the copies.

Split a shared subject by giving each document a different **kind** of statement: the *model* (what
a thing means, what invariant it preserves) in the document that owns the concept; the *contract*
(its surface, parameters, interactions) in the reference for that surface. Neither restates the
other, and each links to the other once.

**A document is created only where no existing one is its home**, and the change that creates it
says which documents were considered and why the fact belongs in neither. A subject that nobody
owns is the reason to write a document; a subject someone owns is the reason not to, because a
second document about it is the duplicate this section exists to prevent, arriving at the one
moment nobody is looking for one.

**A fact whose home is source code stays there.** Prose says where to look, not what it will
say. And a fact whose home is the org corpus stays there: a project document cites `docs/org/`,
it does not restate it.

## Links

Link to the document that owns a fact. If a passage must be edited whenever its target changes,
it is a copy however it is worded — a paraphrase and a quotation drift identically. A copy is
sanctioned only where a machine checks it: the vendored `docs/org/`, byte-identical, verified
against its stamp, and named as a copy by its own header; and the [header](#header) line, which
restates what the directory decided and is checked against it.

A document in another repository is cited by that repository's name and the document's, and a
section by its slug ([sections](#sections)) — never by a link into a branch, whose target moves
without any change in the citing tree, and never by a section number. A citation that must be
re-read whenever the other repository moves is the drift above arriving from the far end.

## Lifecycle

Three transitions, each one reviewed change:

- **Ratification.** A design begins in `docs/proposals/`, unbound, untagged, freely rewritten —
  and **written as the specification it would become**. The [end-state](#end-state) voice,
  [one home](#one-home) per fact, and [links](#links) all apply to it,
  because what makes it a proposal is where it sits, not how it is written. Ratifying it is one
  act: `git mv` into the root, delete the proposal line, move its index entry — the move *is*
  the decision. Nothing else changes: the tag line it already carried stays, its label stays, and
  the items under that label stay. The test mirrors the end state's: **a proposal reads
  identically the day before and the day after it is ratified**, but for its location, the
  proposal line it no longer carries, its index entry, and the paths of the relative links the
  move re-rooted. Questions still open under its tag do not block the move, any more than they
  block a specification; what blocks it is text that would need rewriting first, and a question
  whose answer would rewrite the document is exactly that. A proposal that would need rewriting first is not ready, and a rewrite folded into
  the ratification is a change nobody can diff against what was proposed.
- **Amendment.** An ordinary reviewed diff, landing **before or with** the change that
  implements it, never after: a specification trailing its implementation has stopped describing
  the end state and started reporting history.
- **Retirement.** `git mv` into `docs/archive/`, replace the tag line with the archive line, and
  move its index entry. Content stays; location — and so authority — changes. The label outlives
  the document: deleting it would strip the closed items that record the work.

  **Completion never retires a specification.** A document whose implementation is finished has
  an empty tag query — its healthiest state, not its end. It stays in the root, where it keeps
  the next change from quietly undoing the work and gives every future reconciliation pass its
  measure; retiring it on delivery would turn "implemented" back into "unspecified", the
  [reconciliation](#reconciliation) gap it took the work to close. A specification retires only
  when it stops describing the intended end state: a ratified replacement supersedes it, its
  subject is removed from the project, or the direction is abandoned — each a decision about the
  design, never a report that work finished. The one document that *is* delivered is one that
  was inherently one-shot — a staged migration, a bounded sequence — and it archives when it
  completes; that is the "delivered" in the [location](#location) table's archive row, and it is
  the exception, not the pattern.

## Reconciliation

> **Every gap between a specification and the implementation is covered by an open item carrying
> that document's tag.**

That invariant is what makes the tag query a *complete* status section. Three kinds of gap —
unbuilt, divergent, unspecified — are recorded the same way, and a gap in the *document* is
still a gap: a rule that is wrong, unclear, or missing is covered by an item under the same tag
as a rule that is unmet. After a ratification or amendment, walk the document against the
implementation and file the items that close every gap, as its own change.

**A gap closes from the side that is wrong, never by moving the other side to meet it.** A
compliance gap — the implementation short of the rule — closes by changing the implementation,
and the document stays as written. A definition gap — the rule short of what is intended —
closes by amending the document ([lifecycle](#lifecycle)). Closing the item records that the gap
is gone; an item may not be closed while its gap remains. Relaxing a rule to fit what was built,
or building to a rule known to be wrong, is the quiet deviation the [location](#location) table
forbids, spelled from the other end.

**An item carries what a reader needs to check it in a minute**: the rule — the section, quoted,
not paraphrased — what fails under it, with evidence the reader can confirm, and what would
close it: the change to the implementation, or the replacement text, or the statement that the
item is asking for one.

**One item may cover several gaps, and carries every tag it answers to.** The invariant above is
coverage, not arithmetic: what it forbids is a gap no item names, never gaps that share one. An
item is sized by the change that closes it — one subject, one resolution — because every item
costs a plan, a review and a full run of the gates whether it closes a sentence or a section, and
an item per document is that cost multiplied by the corpus, with an ordering between the pieces
that nobody chose. An item that says only that a section should be clearer cannot be
closed by anything checkable.

> **An issue about a document is filed where the document originates.** A defect in a rule, a
> change request, an amendment proposal — these go to the repository the document's home line
> names, never to a project holding a copy.

Documents live in many repositories; their meaning must not. What a project files locally under
a document's tag is only its **own** gaps against it — the compliance work its tree owes. So a
shared document's tag names two different queries in two places, and both are complete: in the
home repository, the remaining work on the *definition*, together with the home's own compliance
gaps — one tag, because both are remaining work on that document, and a second tag would split
one status into two queries neither of which is whole; in each project, that project's remaining
work toward it. An issue filed in the wrong place is invisible to the query that should have
listed it — the same defect as a fact with two homes. The vendored copies make the wrong place
the *natural* place: a reader meets the document in the project's tree and files where they
stand. Such an issue is transferred to the home repository, not worked where it landed — and the
home line is what tells both the reader and the transferrer where that is.

When the corpus is amended, this same pass runs in every project against the delta: the home
repository's release is what starts it, and each project's tag queries are where it lands.

## Mechanical checks

- Every relative link in every tracked Markdown file resolves to a tracked file, and a link
  carrying a fragment resolves to a heading in it — and **a failing link is repaired by fixing
  the reference, never by removing the link.** Repoint it at the document that now owns the
  fact, or delete the reference together with the claim it supports. A link turned into plain
  text leaves the reference exactly as wrong as it was and makes the check report a pass: a
  repaired reference and a silenced one must never look alike.
- Every tracked file under `docs/` is listed in `docs/index.md` ([location](#location)) — and
  **a missing index is a failure, not a pass**: nothing to report and nothing was checked must
  never look alike.
- The line under every title is the one its location requires ([header](#header)), a tag names
  the file it sits in, and a document that carries a home line carries it in the one form.
- Every heading below a title is a name in the closed alphabet, unique in its document
  ([sections](#sections)).
- `docs/org/` is refused at the edit by the guard and verified against its stamp by the
  integration gate.

Everything else here is upheld by review, and the gaps against this document are items carrying
its tag.
