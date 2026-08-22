# wire-consumer

A separate module that imports `wire` over a relative path and uses it the way a flow would. The
module's own tests prove it compiles; this proves someone else can consume it.

```sh
promise test tests/wire-consumer            # the batch tests
promise test tests/wire-consumer/main.pr    # the snapshot test, which must be named
```

**The snapshot test in `main.pr` is not found by the directory scan**, and that is
[promise-language/promise#37](https://github.com/promise-language/promise/issues/37): a snapshot
test runs only in a file *not* named `*_test.pr`, while discovery only globs `*_test.pr`, so no
filename satisfies both. The batch tests in `consumer_test.pr` do run under the plain command, so
the module is covered either way — the snapshot is kept for the one thing it adds, that a real
`main()` links and executes against the module.

`use` is module-scoped rather than per-file, so it is declared once. It lives in **`main.pr`**
deliberately: the directory scan compiles the whole module either way, but a single-file run of
`main.pr` needs the import in that file. `promise test consumer_test.pr` on its own therefore does
not resolve `wire` — run the directory instead, which is the supported path. Whether that scoping
is right is under review upstream.
