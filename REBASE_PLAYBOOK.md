# Rebasing the fork + its dependency forks onto upstream

This repo (`openrm/krakend-ce`) and five dependency forks --
`openrm/krakend-jose`, `openrm/krakend-cel`, `openrm/krakend-opencensus`,
`openrm/krakend-martian`, `openrm/krakend-bloomd` -- carry a small set of
custom patches on top of the real `krakend/*` (or, for bloomd, an original
openrm) upstream. Periodically each needs to be rebased onto current
upstream so security fixes keep flowing without losing the custom patches.
This is the procedure, written down so it doesn't have to be re-derived.

## 1. Per-fork: pick the right rebase target

**Always check for a major-version-line trap before touching anything.**
Compare `upstream/master`'s `go.mod` module path against the latest tag on
the SAME major-version line krakend-ce actually depends on:

```
git fetch upstream --tags
head -1 <(git show upstream/master:go.mod)          # e.g. module .../v3
for t in $(git tag -l 'v2.*' | sort -V | tail -5); do
  echo "$t: $(git show $t:go.mod | head -1)"        # e.g. v2.12.3: module .../v2
done
grep <module-name> /path/to/krakend-ce/go.mod        # confirm which line krakend-ce pins
```

If `upstream/master` has moved to a newer major version than what
krakend-ce depends on, **rebase onto the latest tag on the old line**, not
onto `upstream/master`. Blindly rebasing onto `upstream/master` in that
case replays commits that don't apply cleanly (different module system,
different APIs) and wastes time resolving conflicts in code you don't even
want. This has bitten every rebase pass so far: krakend-martian and
krakend-jose both had `upstream/master` already on `/v3` while krakend-ce
needs `/v2`. krakend-opencensus was the one exception -- its upstream never
split, so `upstream/master` was the right target.

## 2. Before rebasing, check whether the work already exists on `origin`

More than once this session, `origin/master` already had a commit
implementing the exact same fix, done independently (a parallel session,
or a teammate) -- with a different commit hash, sometimes on a different
base entirely. **Diff actual file content, not commit hashes or messages**
-- hashes differ for byte-identical content across independent rebases,
and similar-sounding commit messages can hide a fix built on a stale base:

```
diff <(git show origin/master:path/to/file.go) path/to/file.go
git merge-base --is-ancestor <the-tag-you-picked-in-step-1> origin/master
```

- If `origin/master`'s version is content-identical (or a strict
  superset) AND already based on the correct target from step 1: don't
  re-push your own redraft, just verify `origin/master` builds/tests and
  move on. (This happened for krakend-cel and krakend-opencensus.)
- If `origin/master`'s version *looks* similar but is NOT based on the
  correct target (the ancestor check fails): it's on a stale, separate
  lineage -- your properly-targeted rebase is still the right one to push,
  even though it means discarding a commit that looks superficially
  similar. (This happened for krakend-martian: `origin/master` had an
  identically-named commit built on go 1.17 and the old `krakendio` module
  path, years behind the `v2.3.0` tag that was actually the right target.)

## 3. Doing the rebase

```
git merge-base master <target>                      # compute fresh -- don't reuse
                                                      # a merge-base computed against
                                                      # a DIFFERENT target (see below)
git checkout -b rebase-onto-<target> master          # branch from YOUR branch,
                                                      # not from <target> itself
git rebase --onto <target> <merge-base>
```

Two mechanical traps hit repeatedly:
- **Recompute the merge-base for the actual target.** A merge-base
  computed against `upstream/master` (wrong target, per step 1) is *not*
  the same commit as the merge-base against the correct tag. Reusing the
  old value replays commits that are already part of the target, which
  looks like an unrelated conflict (e.g. a dependency-bump commit you
  never touched) rather than the real problem.
- **Branch from your own branch, not from the target.** `git checkout -b
  foo <target>` followed by `git rebase --onto <target> <merge-base>`
  rebases `<merge-base>..<target>` -- i.e. the target's own commits --
  onto itself, not your commits. Branch from `master` (or whatever your
  working branch is), *then* rebase `--onto <target>`.

**Resolving conflicts**: when a conflict is two features touching the same
spot (e.g. upstream added header-propagation support right where the
fork's own refresh-token feature also hooks in), check whether the same
package has a sibling file that merged cleanly (e.g. `gin/jose.go` merged
fine while `mux/jose.go` conflicted, for the exact same feature) -- it
shows you the correct combined shape to replicate, rather than guessing.

## 4. Verify before pushing

```
go build ./... && go vet ./... && go test ./...
```

If `go vet`/`go test` fails, **check whether the exact same failure
already exists on the pre-rebase branch** before assuming the rebase
caused it:

```
git stash && git checkout master -- . && go vet ./...   # or just checkout the old branch
```

Twice this session, a `go vet` failure (a function signature had grown a
new required parameter, and call sites in `_test.go` files were never
updated) turned out to pre-date the rebase entirely -- worth fixing while
there since it blocks CI either way, but not the rebase's fault, and
should be described as such in the commit message.

## 5. Pushing: check reachability before any force-push

Every fork here has diverged from `origin/<branch>` (that's the whole
point), so pushing the rebased result is always a force-push. **Never
force-push without checking whether the current tip would become
unreachable:**

```
git fetch origin
TIP=$(git rev-parse origin/master)
git branch -r --contains $TIP
git tag --contains $TIP
```

- If `$TIP` shows up on some *other* ref already (a dependabot branch, an
  old release tag -- this was true for every fork in this session, since
  dependabot branches are cut from master and inherit its history): safe
  to force-push directly.
- If not: create and push an archive branch pointing at `$TIP` **first**,
  then force-push the rebased branch:
  ```
  git branch archive/pre-<description>-rebase $TIP
  git push origin archive/pre-<description>-rebase
  git push --force-with-lease origin master
  ```

## 6. After pushing: point krakend-ce's go.mod at the real commit

Once a fork is pushed, pin a real version in krakend-ce's `go.mod` with a
plain `replace` directive -- **no module renaming, no local or relative
filesystem paths**:

```
go get github.com/<upstream-org>/krakend-X/v2@github.com/openrm/krakend-X/v2@<commit-sha>
```

More precisely, since the fork keeps the exact same module identity the
upstream package has (`github.com/krakend/krakend-X/v2`) even though it's
hosted at a different URL (`github.com/openrm/krakend-X`):

```
go get github.com/openrm/krakend-X/v2@<commit-sha>   # fails, but the error message
                                                       # prints the resolved pseudo-version --
                                                       # copy it from there
go mod edit -replace="github.com/krakend/krakend-X/v2=github.com/openrm/krakend-X/v2@<pseudo-version>"
go mod tidy
```

**This works with zero changes to the fork itself.** Go's `replace`
directive only requires that the *replacement target's own declared
module path* matches the **left-hand side** (the path being replaced) --
not the URL it's actually fetched from. Since every fork here still
declares itself as `github.com/krakend/krakend-X/v2` (never renamed), this
resolves cleanly. The only real mistake to avoid is dropping the `/v2`
suffix on the right-hand side import path (`github.com/openrm/krakend-X`
without `/v2` looks like a *different, incompatible* major version to Go
and fails with a confusing "post-v1 module path" error) -- include it and
this just works.

**Don't rename the fork's module identity to "fix" this.** An earlier pass
this session did exactly that (`github.com/krakend/krakend-X/v2` ->
`github.com/openrm/krakend-X/v2`, touching every self-import in the fork
plus every importer in krakend-ce and krakend-bloomd) based on a
mis-syntaxed replace attempt that looked like a hard Go limitation but
wasn't. It was reverted once the real cause (the missing `/v2` suffix) was
found and verified with an isolated test build. Renaming adds real,
permanent complexity -- more files to touch on every future rebase, extra
review surface -- for a problem the plain `replace` directive already
solves for free.

**No sibling checkouts, no relative paths, no local filesystem paths
needed anywhere.** Since every replace target resolves over the network
(a real GitHub commit, not a local directory), the Dockerfile and CI need
nothing beyond the plain single-repo pattern (`COPY . /app`, one
`actions/checkout`) -- verified by building the image from this repo
alone, with no sibling directories present at all, and Trivy-scanning the
result at 0 vulnerabilities.

## Net effect, this pass (2026-09-26 to 2026-09-29)

krakend-jose rebased onto `v2.12.3` (not `upstream/master`, which had
moved to `/v3`); krakend-cel and krakend-opencensus were already correctly
based once their own commits were checked -- no rebase needed, just
verification; krakend-martian rebased onto `v2.3.0` (not `upstream/master`,
same `/v3` trap), discarding a superficially-similar but badly-stale
`origin/master` commit; krakend-bloomd needed no rebase (it's not a
namespace-preserving fork). krakend-ce itself was already on the latest
available tag (`v2.13.11`) -- no rebase needed there this pass.

All five forks' known CVEs fixed (mostly `golang.org/x/crypto`, plus a
handful of `grpc`/`otel`/`cel-go` ones) and verified at 0 via `trivy fs`/
`trivy image` against the real built binary, not just source inference.

Every dependency, including all five forks, now resolves via a plain
version-pinned `require`/`replace` in krakend-ce's `go.mod` -- no local or
relative filesystem paths, no module renaming.
