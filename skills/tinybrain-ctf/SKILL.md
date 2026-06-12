---
name: tinybrain-ctf
description: >-
  Work a CTF challenge or other authorized security-learning exercise while tracking
  reconnaissance, hypotheses, payloads, and the path to the flag in TinyBrain — so you never
  lose a working payload, repeat a dead end, or forget what you already enumerated. Use this
  skill when the TinyBrain MCP tools are available AND you are solving a CTF, working a
  practice/retired box (HTB, picoCTF, a crackme, a deliberately-vulnerable app like Juice
  Shop or DVWA), or the user says "let's do this CTF", "solve this challenge", "capture the
  flag", or "work this box". Builds on the using-tinybrain skill.
---

# CTF and security-learning workflow with TinyBrain

CTFs and practice boxes are authorized, sandboxed learning. They're also exactly the kind of
multi-stage, easy-to-lose-track-of work TinyBrain is built for: lots of enumeration, many
hypotheses, payloads you'll want again, and a winning path worth remembering for the writeup.

Read the `using-tinybrain` skill first for the core loop. This skill adds CTF-specific habits.

## Workflow

1. **One session per challenge/box.** `create_session` with `task_type: "penetration_test"`
   (boxes) or `"general"` (puzzle-style challenges), named for the challenge. If you're
   resuming, `get_context_summary` first to reload what you've enumerated and tried.

2. **Log reconnaissance as you go.** Open ports, services, versions, endpoints, interesting
   files — each a `create_memory`, `category: "context"` or `finding`, tagged with the host
   and service. This is the map you'll keep returning to; storing it once beats re-scanning.

3. **Track hypotheses as first-class TODOs.** "Login form might be SQLi", "that uploader
   probably allows a webshell" — store as `category: "hypothesis"` at low-to-mid confidence.
   As you test them, `update_memory` to bump confidence and flip toward `vulnerability` (it
   worked) or note it as a dead end (it didn't). This is what stops you from re-trying the same
   idea an hour later.

4. **Never lose a working payload.** The moment a payload, request, or command does something
   useful, store it `category: "payload"` (or `exploit`) with the exact text in `content` and
   high confidence. Link it `exploits` → the vulnerability it leverages. CTF flags are often
   one tweaked payload away from the last one — keep them.

5. **Capture the chain.** When step B depends on what step A gave you (creds from one service
   unlocking another, an LFI enabling log poisoning), link the memories with `depends_on` or
   `causes`. The chain of links *is* your eventual writeup of how the box fell.

6. **Store the flag and the path.** When you get a flag, store it (`category: "evidence"`,
   high priority) and make sure it's linked back through the exploit chain that produced it.

## Before you enumerate again, search

The classic CTF time-sink is re-enumerating something you already mapped. Before re-scanning a
host or re-poking an endpoint, `search_memories` by host/service tag. `get_related_entries`
from your current foothold to see what you already linked to it.

## Track stage progress for multi-stage boxes

For boxes with clear phases (recon → foothold → privesc → root), use
`create_task_progress` / `update_task_progress` so resuming is instant: you can see you're at
"privesc, 60%, trying the SUID binary" without re-reading every note.

## Writeup at the end

Ask TinyBrain to reconstruct the solve: `get_context_summary` for the arc, then traverse the
flag's relationships backward (`get_related_entries`) to recover recon → hypothesis → payload
→ exploit → flag in order. A clean writeup falls out of a well-linked graph.

## Keep it in the authorized lane

This skill is for CTFs, retired/practice boxes, and deliberately-vulnerable training apps —
targets you're explicitly permitted to attack. If a request shifts toward a real,
non-consented target, stop and say so; that's outside what this skill (or TinyBrain) is for.
