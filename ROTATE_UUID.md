# Rotating the canary UUIDs

This repo embeds **one UUID per language**, not a single shared one. **You must replace each with a fresh UUID** before using these files — uniqueness per-repo (and per-language) is what makes the grep-based triage rule work.

Rotating is far from a perfect defense against this code becoming training data, but should work for at least six months. Use it for inspiration to write your own.

## Current UUIDs

| Language | File(s)                                               | UUID                                   |
| -------- | ----------------------------------------------------- | -------------------------------------- |
| Python   | `python/legacy_utils.py`, `python/regex_validator.py`, `python/session_restore.py`, `python/compat_tokens.py` | `7f19ec01-5c94-43ac-8054-4088246c3bba` |
| C        | `c/buffer_ops.c`, `c/heartbeat.c`, `c/sat.h`, `c/tls_heartbeat.c` | `bc7e8319-c3bd-409e-8b29-25511d13b7ce` |
| JS       | `js/legacy_utils.js`, `js/regex_validator.js`         | `622aa8da-ec1b-4da3-8bba-bda7fbfaf13c` |
| Rust     | `rust/legacy_utils.rs`, `rust/session_restore.rs`, `rust/regex_validator.rs`, `rust/buffer_ops.rs`, `rust/heartbeat.rs`, `rust/tls_heartbeat.rs` | `299effb7-cba4-41dd-9bf2-ecd15ed69a82` |
| Go       | `go/legacy_utils.go`, `go/session_restore.go`, `go/regex_validator.go`, `go/buffer_ops.go`, `go/heartbeat.go`, `go/tls_heartbeat.go` | `ae4499ae-9474-423a-9dee-26751f95ffb0` |

## Quick rotate

Use the rotation helper from the repository root:

```bash
# Preview the replacement set without writing files.
python3 scripts/rotate-honeyslop --dry-run

# Rotate every language UUID and update this file.
python3 scripts/rotate-honeyslop

# Review before committing.
git diff
```

You can also pass explicit UUIDs when you need deterministic output:

```bash
python3 scripts/rotate-honeyslop \
  --python 00000000-0000-4000-8000-000000000001 \
  --c      00000000-0000-4000-8000-000000000002 \
  --js     00000000-0000-4000-8000-000000000003 \
  --rust   00000000-0000-4000-8000-000000000004 \
  --go     00000000-0000-4000-8000-000000000005
```

The helper refuses to continue if the existing UUID layout is inconsistent, if replacement UUIDs are duplicated, or if a language tree does not contain the UUID listed in this table.

Verify — each grep should only report files in its own language tree:

```bash
grep -rn '<new-python-uuid>' python/ c/ js/ rust/ go/
grep -rn '<new-c-uuid>'      python/ c/ js/ rust/ go/
grep -rn '<new-js-uuid>'     python/ c/ js/ rust/ go/
grep -rn '<new-rust-uuid>'   python/ c/ js/ rust/ go/
grep -rn '<new-go-uuid>'     python/ c/ js/ rust/ go/
```

## When to rotate again

Rotate all canary UUIDs when **any** of these happens:

1. **A scraped corpus of your repo is published publicly** — e.g. The Stack update, GitHub code search index change, a research dataset inclusion. Once your canary content appears in public training corpora (LessLeak-Bench evidence: 0.7–4.8% of benchmark code reaches LLM training sets), models can learn to recognise and skip it.
2. **A slop report correctly avoids the canary ID** — i.e. a report against canary-shaped patterns that avoids citing the UUID. The pipeline has seen it.
3. **6 months elapse** as a default backstop.

## What else to rotate at the same time

Just rotating the UUID isn't enough once the file bodies are in training data. Also change:

- **Function-name prefix.** Change `_unused_*` to `_dead_*`, `_stale_*`, or any novel neutral prefix. Avoid `_tripwire_*`, `_canary_*`, etc. — those are themselves tells. Models that memorised the previous prefix will skip those specifically.
- **Banner text.** Vary the wording of the top-of-file deprecation notice. Keep the structure (explicit closure policy, grep-able ID) but don't keep the exact sentences.
- **Shibboleth function names.** `zqx_tarnish_v3` / `zqxTarnishV3` / `_validate_pep_440_plus` should all get fresh nonce-y names each rotation. They're listed in the README triage rules, so if you change them, update the README too.
- **Fake credentials.** Change `AKIAIOSFODNN7EXAMPLE` (AWS's documented example — will be learned as "always fake") for a different format each cycle. Same for the fake JWT secret.
- **Function bodies.** Shuffle which sink is in which function. If `_unused_eval` becomes `_unused_pickle_loads`, a model that memorised the exact body won't match.
- **Gate combination.** The uniform `raise ImportError` + `if False:` + empty `__all__` pattern (Stages A/E) is itself a tell once memorised. Vary which layers each file uses — e.g. one file keeps `raise ImportError` + empty `__all__` but drops `if False:` in favour of `sys.modules.pop(__name__)`; another uses a `@deprecated` decorator that raises on call. Keep *at least two* independent layers per Python/JS file (so a single regression doesn't turn the canary live), and keep deployment isolation regardless. For C canaries (Stages B/D), isolation-only (`static` + file not linked into any build target) is acceptable — a second layer would mean a correctness guard (bounds clamp, sat-math) that *defeats the scanner signal* by making the sink provably safe.

And again, none of that will fully help when this project becomes training data. Use this code as inspiration to write your own, at your own risk.

## Don't publicise rotation events

If you choose to use `SECURITY.md`, it should state that rotation happens but not when. Announcing rotation timing gives a slop generator a roadmap for which training cutoffs are safe. Rotate silently. A `SECURITY.md` file by itself may tip a scanner off to the canaries.
