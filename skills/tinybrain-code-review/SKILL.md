---
name: tinybrain-code-review
description: >-
  Run a security-focused code review while recording findings in TinyBrain so they form a
  navigable, deduplicated, standards-mapped knowledge base instead of a one-off report.
  Use this skill when the TinyBrain MCP tools are available AND you are reviewing code for
  security issues, doing a secure code audit, triaging a vulnerability in a codebase, or the
  user asks you to "review this for security", "audit this code", "find vulnerabilities in",
  or map findings to CWE/CVE/OWASP. Builds on the using-tinybrain skill.
---

# Security code review with TinyBrain

This is defensive, authorized review work: finding weaknesses in code so they can be fixed.
TinyBrain turns a review from a throwaway pass into durable, structured knowledge — every
finding categorized, linked to its evidence and its fix, mapped to CWE/CVE, and reusable in
the next review of the same codebase.

Read the `using-tinybrain` skill first for the core capture/connect/retrieve loop. This skill
layers a review-specific workflow on top.

## Workflow

1. **Orient and open a session.** Check `list_sessions` for prior reviews of this codebase;
   if one exists, `get_context_summary` to recall known issues before re-reading anything.
   Otherwise `create_session` with `task_type: "security_review"`, naming it for the codebase
   or PR under review.

2. **Record findings as you read.** Each issue is a `create_memory` with:
   - `category: "vulnerability"` for a confirmed weakness, `hypothesis` for something that
     looks suspicious but you haven't confirmed, `finding` for a noteworthy observation.
   - `title` naming the issue and location: "SQLi in `UserRepo.findByName` — string-concat query".
   - `content` with the vulnerable code, the data flow (source → sink), and why it's exploitable.
   - `priority` from impact × exploitability; `confidence` from how sure you are it's real.
   - `tags`: the file, the component, and the **CWE id** (e.g. `CWE-89`) — this is what lets
     you map to CVEs later and correlate across the review.
   - `source`: the file path and line range.

3. **Capture evidence separately and link it.** Store the proof — the exact sink, a PoC
   request, a failing test — as a `category: "evidence"` memory, then
   `create_relationship` with `supports` from the evidence to the vulnerability. Evidence that
   travels with the finding is what makes a report credible and a fix verifiable.

4. **Record the fix as a recommendation and link it.** When you propose remediation, store it
   as `category: "recommendation"` and link it `mitigates` → the vulnerability. Now the
   finding, its proof, and its fix are one traversable cluster.

5. **Map to standards.** For confirmed weaknesses, use `map_to_cve` with the CWE id to pull
   related CVEs, and `map_to_compliance` with the relevant standard (e.g. OWASP, PCI-DSS) to
   show coverage and gaps. Store notable results as `reference` memories linked to the finding.

6. **Correlate before reporting.** Run `analyze_risk_correlation` on the session to surface
   chains where individually-medium issues combine into something worse (e.g. an IDOR plus a
   verbose error becomes account enumeration). `check_high_priority_memories` gives you the
   lede for the writeup.

## Before you re-read code, search

If you're about to re-examine a function or revisit a pattern, `search_memories` (by file,
CWE, or component tag) first — you may have already reviewed it. `check_duplicates` before
storing a finding in a large review so the same issue across call sites doesn't get logged
five times; if it's genuinely the same root cause at multiple sites, store once and note the
sites in `content`, or link the instances with `related_to`.

## Producing the report

When asked for the writeup, retrieve rather than recall: `get_context_summary` for the
overview, `check_high_priority_memories` for the criticals, and walk each finding's
relationships (`get_related_entries`) to assemble finding → evidence → recommendation as a
unit. The memory graph *is* the report's structure.

## Calibration reminders specific to review

- A reachable, exploitable sink with attacker-controlled input is high priority (8–10). An
  issue gated behind auth or requiring an unlikely precondition is lower — say so in `content`.
- Don't inflate confidence on findings you haven't traced end-to-end. A `hypothesis` at 0.4
  ("this looks like it might be injectable, haven't confirmed the source is tainted") is an
  honest, useful record — and a reminder to go confirm it.
