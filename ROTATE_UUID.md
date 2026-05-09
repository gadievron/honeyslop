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

## Quick rotate

```bash
# Generate fresh UUIDs (one per language)
NEW_PY=$(uuidgen | tr 'A-Z' 'a-z')
NEW_C=$(uuidgen  | tr 'A-Z' 'a-z')
NEW_JS=$(uuidgen | tr 'A-Z' 'a-z')
NEW_RS=$(uuidgen | tr 'A-Z' 'a-z')

OLD_PY="7f19ec01-5c94-43ac-8054-4088246c3bba"
OLD_C="bc7e8319-c3bd-409e-8b29-25511d13b7ce"
OLD_JS="622aa8da-ec1b-4da3-8bba-bda7fbfaf13c"
OLD_RS="299effb7-cba4-41dd-9bf2-ecd15ed69a82"

# Replace in each language tree (review the diff before committing).
# `xargs -r` (GNU) skips the sed call if grep finds nothing — otherwise
# sed would hang reading from stdin.
grep -rl "$OLD_PY" python/ | xargs -r sed -i.bak "s/$OLD_PY/$NEW_PY/g"
grep -rl "$OLD_C"  c/      | xargs -r sed -i.bak "s/$OLD_C/$NEW_C/g"
grep -rl "$OLD_JS" js/     | xargs -r sed -i.bak "s/$OLD_JS/$NEW_JS/g"
grep -rl "$OLD_RS" rust/   | xargs -r sed -i.bak "s/$OLD_RS/$NEW_RS/g"
find python c js rust -name '*.bak' -delete

# Don't forget the table at the top of this file
sed -i.bak "s/$OLD_PY/$NEW_PY/g; s/$OLD_C/$NEW_C/g; s/$OLD_JS/$NEW_JS/g; s/$OLD_RS/$NEW_RS/g" ROTATE_UUID.md
rm -f ROTATE_UUID.md.bak

# Verify — each grep should only report files in its own language tree
grep -rn "$NEW_PY" python/ c/ js/ rust/
grep -rn "$NEW_C"  python/ c/ js/ rust/
grep -rn "$NEW_JS" python/ c/ js/ rust/
grep -rn "$NEW_RS" python/ c/ js/ rust/
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
