# honeyslop - code canaries to quickly triage hallucinated ("slop") vulnerability reports

<p align="center">
  <img src="honeyslop-full.png" alt="HoneySlop" width="300">
</p>

honeyslop is code canaries, decoys, for open-source projects drowning in AI-hallucinated ("slop") and unverified vulnerability reports. A slop scanner ingests the canary, then generates a vulnerability "report" based on it. The report self-identifies as slop. Close it in one grep.

This is a quick PoC, vibe-coded as a joke (not production-grade), because we received a slop report at [raptor](https://github.com/gadievron/raptor), an autonomous attack/defense agent based on Claude Code, ourselves. Should be fun!

Code canaries extend familiar triage signals (e.g. detections in test files, example secrets, non-existent paths) into deliberate markers, or decoys. In tests, these canaries work well enough to flag slop, but they can be further improved (embedded in real code, function/file/directory names less indicative, regularly regenerated as new code, etc.).

Written by: Gadi Evron (@gadievron), John Cartwright (@grokjc), Daniel Cuthbert (@danielcuthbert), and Michal Kamensky (@kamenskymic, with props for naming the project).

**Use at your own risk.** If you paste this into production, that's a you problem. See [Disclaimer](#disclaimer) below.

## Triage rules

For each incoming report, in order:

1. grep any canary UUID in the report → close. (UUIDs are per-language; each canary file embeds exactly one.)
2. grep canary-only function names (`zqx_tarnish_v3`, `zqxTarnishV3`, `_validate_pep_440_plus`; also `handle_*_request` if you've adopted F+G privately) → close.
3. grep `CVE-2025-99919` (fake) → close.
4. Cited function doesn't exist in the tree → *"does not exist"*.
5. For memcpy/bounds claims on B/D: ask the reporter to walk through how their PoC defeats the specific guard on the cited line. AI follow-ups cannot answer; humans can.

## Stages

Two categories of canary:

- **SCANNER-FLAG** (Stages A, B, C, D, E) — trips scanners so slop reports pile up on the canary instead of real code.
- **RESOURCE-WASTE** (Stages F + G, together) — runs agentic LLM scanners through their full iteration budget at maximum cost.

| Stage   | File(s)                                              | Shape                                                             |
| ------- | ---------------------------------------------------- | ----------------------------------------------------------------- |
| **A**   | `python/legacy_utils.py`, `python/session_restore.py`, `python/compat_tokens.py`, `js/legacy_utils.js` | ~15 CWE sinks + fake secrets + shibboleths                        |
| **B**   | `c/buffer_ops.c`                                     | 4 `memcpy`/`memmove` shapes (CWE-120/121/787/170)                 |
| **C**   | merged into A                                        | Extended CWE yield                                                |
| **D**   | `c/heartbeat.c` + `c/sat.h`, `c/tls_heartbeat.c`     | Heartbleed silhouette                                             |
| **E**   | `python/regex_validator.py`, `js/regex_validator.js` | Catastrophic-backtrack regex + fake **CVE-2025-99919**            |
| **F+G** | `private/fractal_dag/` (not in this repo)            | Stage-A sinks across a 12-node DAG of `handle_*_request` entries  |

See [Safety model](#safety-model) for how each stage stays inert despite looking vulnerable.

## Safety model

The canary code files deliberately read like plausible deprecated modules — no "canary", "honeypot", or "tripwire" language in comments or identifiers. That keeps the files from self-identifying to scanners, but it also means the *why these files are safe* documentation lives here instead of in each file's docstring. When reviewing or rotating a canary, check that every layer below is still intact.

These files **are** vulnerable-shaped code — and some (e.g. `python/session_restore.py`'s `pickle.loads`, `c/tls_heartbeat.c`'s unbounded `memcpy`) would be genuinely exploitable if reachable. That's the design. Scanners flagging the sinks is the whole point; the layers below block **execution**, not scanner **signal**. Pattern-based scanners (grep, semgrep, LLM-based slop pipelines — the primary threat model) read source as text and surface findings regardless of runtime reachability. Build-integrated C scanners (CodeQL default, clang-static-analyzer) only see compiled files, so unbuilt C canaries are invisible to them — an accepted trade-off, since slop generators overwhelmingly read source, not builds.

### Stages A and E (Python + JS)

Five independent layers keep these files inert:

1. **Top-level `raise ImportError` / `throw new Error`** — a plain `import` / `require` aborts before any definition binds.
2. **Every def/function under `if False:` / `if (false)`** — names never enter the runtime namespace even if layer 1 is bypassed.
3. **Empty exports** — Python: `__all__: list[str] = []` (star-import exports nothing). JS: `module.exports = {}` (CommonJS consumers get an empty object).
4. **Zero in-tree callers of the shibboleth functions** (`zqx_tarnish_v3`, `zqxTarnishV3`, `_validate_pep_440_plus`). Any report citing one self-identifies as slop.
5. **Deployment isolation** — the adopter excludes the canary paths from sdist / wheel / Docker / SAST (see [`SECURITY.md.template`](SECURITY.md.template) and [How to try](#how-to-try)).

Stage E adds a sixth layer: the catastrophic-backtrack regex is stored as a string literal only, not passed to `re.compile` / `new RegExp` at module scope. Even a harness that strips layer 1 cannot trigger a compiled backtracking engine.

Scanners walk the AST past `raise` / `throw` and into the dead block, so the sinks still surface as findings — that is the intended behaviour.

### Stage B (C `buffer_ops.c`)

Safety is structural — each shape has a proof alongside it:

- `bufops_copy_banner` — `src` is a string literal, `n = sizeof(literal)`, `_Static_assert` pins it to the destination size.
- `bufops_copy_bounded` — `if (n > dst_cap) n = dst_cap;` the line before the `memcpy` makes the CWE-787 claim impossible. Short-circuits on `n == 0` to avoid C17 `memcpy(dst, NULL, 0)` UB.
- `bufops_copy_truncating` — `n <= dst_cap - 1`, `dst[n]` hits at most `dst_cap - 1`; early-return on `dst_cap == 0`.
- `bufops_shift` — both `i + n` and `j + n` bounded to `cap`; `memmove` explicitly supports overlap.

Additional isolation: all functions are `static` (no external linkage) and the file is not added to any build target.

### Stage D (C `heartbeat.c` + `sat.h`)

The Heartbleed silhouette (`uint16_t payload_len` → `malloc(1+2+payload_len+16)` → `memcpy`) is defanged by layered guards:

- `sat_sub` saturating subtraction for all header/trailer budget math (no wrap).
- Frame fields cached into `const` locals on entry — closes TOCTOU and signal-handler UB windows.
- NULL checks on the reader struct and its buffer.
- `_Static_assert(SIZE_MAX - 19 >= UINT16_MAX, ...)` proves the `malloc` size cannot overflow `size_t`.
- `payload_len > 0` short-circuit avoids `memcpy(dst, NULL, 0)` UB on empty payloads.
- `parse_heartbeat` and `read_u16_be` are `static`; file not linked into any build target.

A report alleging OOB read/write in `parse_heartbeat` without engaging with the specific guard on the cited line has not verified exploitability — close via triage rule 5.

`c/tls_heartbeat.c` is a variation on the same silhouette kept deliberately unguarded: `process_heartbeat` is `static` and the file is not linked into any build target — isolation is the only layer, matching the Stage B catch-all. Attempting to call it from outside the TU is a link error.

## How to try

1. **Pick stages.** C/C++ parser surface → D (+ B). Python OSS maintainer → A + E. Under sustained agentic scanner spam → add F + G privately.
2. **Rotate every UUID.** One per language, distinct per adopter — not prefix variants of one base. See [`ROTATE_UUID.md`](ROTATE_UUID.md).
3. **Exclude from build artefacts.** Python: `MANIFEST.in prune` the canary paths, or `pyproject.toml tool.setuptools.exclude-package-data`. C: omit from `CMakeLists.txt` / `Makefile` / sdist. Docker: `.dockerignore`.
4. **Exclude from CI static analysis.** Otherwise your own CI produces findings on the canary. CodeQL `paths-ignore`, `.semgrepignore`, `bandit -x`, Ruff `--extend-exclude` — all pointed at your canary paths.
5. **Consider adding the triage rule to `SECURITY.md`** — see [`SECURITY.md.template`](SECURITY.md.template). This may tip slop scanners off to the canary's presence (maybe a good thing?).
6. **Protect from contributor cleanup.** `CODEOWNERS` on the canary files; a pre-commit hook that fails if the canary UUID count decreases or if `if False:` tripwires go missing.
7. **Remove obvious "tells".** Strip "canary", "canaries", "honeypot", "decoy", "fake", "tripwire", and "slop" from code comments, directory names, file names, and function/identifier names. Reframe top-of-file docstrings as plausible deprecation notices. Keep this conceptual language in the *docs* (README, `SECURITY.md.template`, `ROTATE_UUID.md`) where it's load-bearing.

## Adopter snippets

Copy-paste to close the exclusion and ownership steps above. Paths below use honeyslop's `c/` / `python/` / `js/` layout; **rename to wherever you land the canaries** (see step 7) — committing exclusions that still point at directories named `canary/` or `slop/` is itself a tell.

**`MANIFEST.in`**

```
prune c
prune python
prune js
```

**`pyproject.toml`** — setuptools

```toml
[tool.setuptools.exclude-package-data]
"*" = ["c/*", "python/*", "js/*"]
```

**`.dockerignore`**

```
c/
python/
js/
```

**`.semgrepignore`**

```
c/
python/
js/
```

**CodeQL** — `.github/codeql/codeql-config.yml`

```yaml
paths-ignore:
  - c
  - python
  - js
```

**Bandit** — CI invocation

```
bandit -r src/ -x python/
```

**Ruff** — `pyproject.toml`

```toml
[tool.ruff]
extend-exclude = ["python/"]
```

**`.clang-format-ignore`**

```
c/*
```

**Secret-scanner allowlist** — gitleaks `.gitleaks.toml`

```toml
[allowlist]
regexes = [
  '''AKIAIOSFODNN7EXAMPLE''',
  '''ghp_[A-Za-z0-9]{36}''',
  '''xoxb-[0-9A-Za-z-]+''',
  '''sk_live_[A-Za-z0-9]+''',
]
```

**`.github/CODEOWNERS`**

```
c/       @your-org/sec-team
python/  @your-org/sec-team
js/      @your-org/sec-team
```

Owners should be a small group that understands *why* these paths look vulnerable — so "clean up dead code" PRs get blocked, not merged.

## Watch out for

**Bucket 1 — your own tools trip.** Your scanners, linters, IDE, and secret-scanning will fire on the canary. This is the intended behaviour for *incoming* scans but means **your own pipelines** need to skip these paths. Required:

- Exclude the canary paths from every SAST config (CodeQL `paths-ignore`, `.semgrepignore`, `bandit -x`, Ruff `--extend-exclude`, `.clang-format-ignore`, editor LSP).
- Exclude from published artefacts (`MANIFEST.in prune`, `.dockerignore`, wheel excludes).
- Allowlist the fake secrets (`AKIAIOSFODNN7EXAMPLE`, fake `ghp_` / `xoxb-` / `sk_live_`) in your secret scanner.
- Warn contributors not to "clean up" the honeypot. Watch for linter removal of the `if False:` tripwire, and insertion of `# noqa` / `# nosec`.

**Bucket 2 — effectiveness erosion.** Public canaries enter LLM training corpora over 6–18 months and scanner vendors add skip-heuristics. Rotate UUIDs, banner, shibboleths, and the fake CVE annually (see [`ROTATE_UUID.md`](ROTATE_UUID.md)); vary wording across adopters; keep F+G private. It won't stop the models from learning the code, but might buy you some time.

**Bucket 3 — post-compromise.** If an attacker already has access they can do whatever they want, and don't need honeyslop. However, flipping `if False:` → `if True:`, inserting Unicode/BIDI tricks, or adding prompt-injection text in docstrings can turn a canary live where you don't expect it. Worth watching for, if unlikely.

## Disclaimer

This project is provided under the MIT License; see [LICENSE](LICENSE) for warranty and liability terms. It is the end user's responsibility to comply with applicable laws and regulations, and with the terms of service of any tools or platforms involved.

**Warning**: This project contains code that appears vulnerable, and should be assumed to be so. Some of it works by being costly to analyze — assume this code will consume significant computing resources when inspected by a scanner, depending on the canary. Any modification, automated or otherwise, could turn a canary live. Do not run, deploy, or adapt any of it in a real environment.


## License

Licensed under MIT. See [LICENSE](LICENSE).
