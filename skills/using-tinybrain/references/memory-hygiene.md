# Memory hygiene: the controlled vocabularies

These are the values TinyBrain's schema accepts and what each one means. Using them
consistently is what makes later filtering, ranking, and correlation work. When in doubt,
pick the most specific value that fits.

## Session task types

Pass one to `create_session` as `task_type`. It frames the whole engagement.

| task_type | Use for |
|---|---|
| `security_review` | General security review of a system or design |
| `penetration_test` | Authorized penetration test of an in-scope target |
| `exploit_dev` | Developing or validating an exploit (authorized/research) |
| `vulnerability_analysis` | Analyzing a specific vulnerability or class |
| `threat_modeling` | Modeling threats against a system (defensive) |
| `incident_response` | Investigating/responding to an incident |
| `intelligence_analysis` | OSINT/threat-intel collection and analysis |
| `reverse_engineering` | Reversing a binary, protocol, or format |
| `malware_analysis` | Analyzing malware behavior/structure |
| `general` | Anything that doesn't fit a specific type |

For `intelligence_analysis` and related work, also set `intelligence_type` (one of `osint`,
`humint`, `sigint`, `geoint`, `masint`, `techint`, `finint`, `cybint`), and optionally
`classification` and `threat_level` (`low`/`medium`/`high`/`critical`). These are stored in
session metadata.

## Memory categories

Pass one to `create_memory` as `category`. Pick the most specific.

**Core security categories**

| category | Meaning |
|---|---|
| `finding` | A confirmed observation of interest |
| `vulnerability` | A specific weakness that could be exploited |
| `exploit` | A working exploit or exploitation method |
| `payload` | A concrete payload (input, shellcode, request) |
| `technique` | A method or TTP (how something is done) |
| `tool` | A tool used or worth using |
| `reference` | External reference material (link, doc, CVE) |
| `context` | Background/environmental context |
| `hypothesis` | An unverified idea worth testing (a TODO with a trail) |
| `evidence` | Proof supporting a finding (request/response, screenshot, log) |
| `recommendation` | A remediation or next-step recommendation |
| `note` | A general note that doesn't fit elsewhere |

**Intelligence categories** — `intelligence`, `osint`, `humint`, `sigint`, `geoint`,
`masint`, `techint`, `finint`, `cybint`

**Intelligence objects** — `threat_actor`, `attack_campaign`, `ioc`, `ttp`, `pattern`,
`correlation`

**Reconnaissance** — `target_analysis`, `infrastructure_mapping`, `vulnerability_assessment`

**Analysis** — `malware_analysis`, `binary_analysis`, `vulnerability_research`

## Priority scale (0–10)

A triage signal. The point is *spread* — if everything is a 7, the field is useless.

| Range | Meaning |
|---|---|
| 9–10 | Critical: exploitable now, or directly mission-relevant |
| 7–8 | High: serious, likely exploitable or important |
| 4–6 | Medium: worth attention, not urgent |
| 1–3 | Low: minor observation, context, or housekeeping |
| 0 | Negligible / placeholder |

Drives `check_high_priority_memories`, `cleanup_low_priority_memories`, and ranking.

## Confidence scale (0.0–1.0)

How sure you are this is true. Label honestly — a confident guess and a verified fact are
different things, and conflating them is how investigations go wrong.

| Range | Meaning |
|---|---|
| 0.9–1.0 | Verified / reproduced |
| 0.7–0.9 | Strong evidence, not fully confirmed |
| 0.4–0.7 | Plausible, partial evidence |
| 0.1–0.4 | Hunch / hypothesis, untested |

Storing a low-confidence `hypothesis` is encouraged — it's a tracked lead. Just don't let a
0.3 later read as a fact.

## Relationship types

Pass one to `create_relationship` as `relationship_type`. Name the *actual* relationship and
mind the direction (source → target).

| relationship_type | Reads as (source → target) | Typical use |
|---|---|---|
| `supports` | source supports target | evidence → finding |
| `exploits` | source exploits target | exploit → vulnerability |
| `mitigates` | source mitigates target | recommendation → vulnerability |
| `depends_on` | source depends on target | finding → precondition |
| `causes` | source causes target | misconfig → exposure |
| `references` | source references target | note → external reference |
| `contradicts` | source contradicts target | evidence → hypothesis |
| `related_to` | loosely associated | catch-all when nothing more specific fits |
| `parent_of` / `child_of` | hierarchy | campaign → sub-incident |

A relationship also takes an optional `strength` (0.0–1.0) and `description` — use the
description to say *why* they're linked in a few words; it's invaluable when traversing later.

## Tagging

Tags are free-form, cross-cutting handles category can't express: the target host, an
affected component, a CWE id, an attack class (`idor`, `ssrf`), a campaign name. Good tags
are the difference between "I know it's in here somewhere" and finding it in one search.

The intelligence parameters on `create_memory` (`intelligence_type`, `threat_level`,
`mitre_tactic`, `mitre_technique`) are stored as namespaced tags (`intel:osint`,
`threat:high`, `mitre_technique:T1190`), so you can later search for, say, everything tagged
with a given ATT&CK technique.
