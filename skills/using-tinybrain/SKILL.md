---
name: using-tinybrain
description: >-
  Use TinyBrain (the security-focused memory MCP server) to capture findings, build a
  knowledge graph of an engagement, and recall it across sessions. Invoke this skill
  whenever the TinyBrain MCP tools are available AND you are doing multi-step security
  work — code review, vulnerability analysis, CTFs, threat modeling, intelligence
  analysis, malware/binary analysis, or any task where findings accumulate and you'll
  want them later. Also use it whenever the user mentions TinyBrain, "store this finding",
  "remember this for the assessment", "what did we find", or asks you to track an
  investigation. Reach for it proactively — the whole point of TinyBrain is that you
  record as you go rather than scrambling to reconstruct context at the end.
---

# Using TinyBrain

TinyBrain is **persistent, structured memory for security work**. Treat it as an
extension of your own working memory: a place to write findings down the moment you
notice them, link them into a graph, and read them back in a later turn or a later
session instead of re-deriving them.

The failure mode this skill exists to prevent: doing a whole investigation in your head,
then either forgetting the details after a context compaction or dumping an unstructured
wall of text at the end. The fix is simple and habitual — **capture as you go, connect
what relates, and retrieve before you re-investigate.**

This skill assumes the TinyBrain MCP tools are connected. If they are not, say so rather
than guessing; the tools may surface with a server prefix (e.g. `mcp__tinybrain__create_session`)
— refer to them by their logical names below regardless of prefix.

## The core loop

Most security work with TinyBrain follows the same rhythm. You won't do every step every
time, but keep the whole loop in mind:

1. **Orient** — before starting, check whether memory already exists for this work.
   Call `list_sessions` and, if you find a relevant one, `get_context_summary` to pull
   back what you already know. Resume, don't restart.
2. **Open a session** — `create_session` with the right `task_type` (and, for intel work,
   `intelligence_type`/`classification`/`threat_level`). One session = one engagement,
   challenge, or investigation. The session is the container everything else hangs off.
3. **Capture** — as findings appear, `create_memory` each one with a real category,
   calibrated priority/confidence, and tags. Store *while you work*, not at the end.
4. **Connect** — when two memories relate (evidence supports a finding, one vuln enables
   another, a mitigation addresses a weakness), `create_relationship`. This is what turns
   a pile of notes into a knowledge graph you can traverse.
5. **Checkpoint** — at meaningful transitions, `create_context_snapshot` so a future
   session (or you, after a compaction) can restore the working state. Track multi-stage
   work with `create_task_progress` / `update_task_progress`.
6. **Retrieve** — before re-investigating anything, search memory first. `search_memories`
   for keyword/FTS, `semantic_search` for "I meant the same idea in other words",
   `get_related_entries` to walk the graph from a known node.

## Capture well — the part that matters most

A memory is only as useful as it is findable and trustworthy later. Five fields carry that
weight; spend a moment getting them right rather than dumping defaults.

- **title** — specific and searchable. "SQLi in `/login` username param" beats "found a bug".
  Future-you searches by these words.
- **category** — the controlled vocabulary that lets you filter ("show me every `exploit`").
  Pick the most specific fit: `finding`, `vulnerability`, `exploit`, `payload`, `technique`,
  `tool`, `reference`, `context`, `hypothesis`, `evidence`, `recommendation`, `note`, plus
  intelligence categories (`osint`, `ioc`, `threat_actor`, ...). Full list and meanings in
  `references/memory-hygiene.md`.
- **priority (0–10)** — triage signal. Calibrate honestly: a critical, exploitable finding
  is 9–10; a passing observation is 1–2. If everything is priority 7, nothing is. This drives
  `check_high_priority_memories` and ranking.
- **confidence (0.0–1.0)** — how sure you are. A verified, reproduced finding is ~0.95; an
  untested hunch is ~0.3. Storing a low-confidence `hypothesis` is good — it's a TODO with a
  paper trail — but label it honestly so you don't later mistake a guess for a fact.
- **tags** — cross-cutting handles that category can't express (`auth`, `idor`, the target
  host, the CWE). These make tag-search and later correlation work.

When a finding maps to MITRE ATT&CK, pass `mitre_tactic` / `mitre_technique` (e.g. `T1190`);
TinyBrain stores them as searchable tags so you can later pull "everything tagged Initial
Access". See the domain skills for when this pays off.

**Why store as you go:** the value compounds only if the record exists before you need it.
A finding you "remember to log later" is one a compaction can erase first. Logging it the
moment you confirm it costs one tool call and buys durability.

## Connect — build the graph, not a list

Relationships are what make TinyBrain more than a notepad. When you create a memory that
relates to an existing one, link them with `create_relationship` using the type that names
the actual relationship: `evidence → supports → vulnerability`, `vulnerability → exploits`,
`mitigation → mitigates → weakness`, `finding → depends_on → finding`. The full vocabulary
(`depends_on`, `causes`, `mitigates`, `exploits`, `references`, `contradicts`, `supports`,
`related_to`, `parent_of`, `child_of`) is in `references/memory-hygiene.md`.

The payoff comes at retrieval: from any node, `get_related_entries` walks the graph — so
from a vulnerability you can reach its evidence, its exploit, and the recommendation that
closes it, without remembering their titles.

## Retrieve — never re-derive what you already found

Before re-investigating a host, re-reading a function, or re-deriving a conclusion, **search
first.** Pick the retrieval mode that fits:

- `search_memories` — keyword / full-text. Best when you know roughly what words are in it.
- `semantic_search` — conceptual similarity. Best when you remember the idea, not the wording.
- `get_related_entries` — graph traversal from a known memory.
- `get_context_summary` — a session-level digest; call it when resuming work to reload the
  whole picture at once.
- `find_similar_memories` / `check_duplicates` — before storing, to avoid logging the same
  finding twice.

## Orient at the start of security work

When the user opens (or returns to) a security task and TinyBrain is connected, take a beat
to orient before diving in: `list_sessions` to see if this work already has a home, and
`get_context_summary` on the relevant session to reload findings. This is the single highest
-leverage habit — it's what makes memory feel continuous instead of amnesiac.

## Anti-patterns to avoid

- **Hoarding without retrieving.** Storing is half the loop; if you never search, you've just
  made write-only notes. Search before you re-investigate.
- **Flat priority/confidence.** Defaulting everything to 5 / 0.5 throws away the triage signal
  the schema exists to capture.
- **Orphan memories.** A finding with no relationships is hard to navigate to later. Link it.
- **End-of-session data dumps.** Storing everything at the end defeats the durability purpose
  and usually loses detail. Capture in the moment.
- **Re-deriving silently.** If you find yourself re-reading the same code or re-scanning the
  same host, you skipped a search.

## Going deeper

- `references/memory-hygiene.md` — full category, priority, confidence, and relationship
  vocabularies with calibration guidance.
- `references/tool-reference.md` — all TinyBrain tools grouped by purpose, with when to use
  each.
- For specific engagement types, see the companion skills: `tinybrain-code-review`,
  `tinybrain-ctf`, and `tinybrain-threat-intel`. They build on this loop with workflow
  tailored to each domain.
