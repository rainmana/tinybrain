---
name: tinybrain-threat-intel
description: >-
  Collect, structure, and correlate threat intelligence in TinyBrain — indicators, threat
  actors, campaigns, and TTPs — mapped to MITRE ATT&CK so the picture builds across sources
  and sessions instead of living in scattered notes. Use this skill when the TinyBrain MCP
  tools are available AND you are doing intelligence analysis, threat modeling, OSINT
  collection, IOC tracking, actor/campaign profiling, or the user mentions threat intel,
  MITRE ATT&CK mapping, IOCs, "track this threat actor", or "model the threats against".
  Builds on the using-tinybrain skill.
---

# Threat intelligence & threat modeling with TinyBrain

Intelligence work is the canonical "many sources, one evolving picture" problem — exactly
what TinyBrain's categories, relationships, and ATT&CK tagging are designed for. The goal is a
living graph: indicators tied to the actors and campaigns they belong to, behaviors mapped to
ATT&CK, and correlations that emerge as you add sources.

Read the `using-tinybrain` skill first for the core loop. This skill adds intel-specific
structure.

## Workflow

1. **Open an intelligence session.** `create_session` with `task_type: "intelligence_analysis"`,
   set `intelligence_type` (`osint`, `cybint`, `sigint`, ...), and `classification` /
   `threat_level` as appropriate. For defensive modeling of a system, use
   `task_type: "threat_modeling"` instead.

2. **Use the right object categories.** Intel has its own vocabulary — use it so the graph
   stays queryable:
   - `ioc` — indicators (hashes, IPs, domains, URLs). Tag with type and the actor/campaign.
   - `threat_actor` — actor profiles (capabilities, motivation, known targeting).
   - `attack_campaign` — campaigns (timeline, victims, attribution confidence).
   - `ttp` / `technique` — observed behaviors. Tag these with their ATT&CK ids.
   - `correlation` — a link you've inferred across sources.
   Plus the collection categories (`osint`, `humint`, ...) for raw collected material.

3. **Map behaviors to MITRE ATT&CK.** When you record a TTP, pass `mitre_tactic` and
   `mitre_technique` (e.g. tactic `TA0001` Initial Access, technique `T1566` Phishing).
   TinyBrain stores these as searchable tags, so you can later pull "every technique this
   actor uses" or "everything mapped to Initial Access" across the whole graph. Use
   `query_attack` to look up technique details to enrich a record.

4. **Build the actor/campaign graph.** Link aggressively — this is where intel value lives:
   - `ioc` → `related_to` → `threat_actor` (this indicator belongs to this actor)
   - `technique` → `child_of` → `attack_campaign` (this behavior was seen in this campaign)
   - `attack_campaign` → `related_to` → `threat_actor` (attribution)
   - `evidence` → `supports` → any attribution or assessment claim
   Use the relationship `description` to record *why* you linked them and your confidence in
   the link — attribution is probabilistic, and the reasoning matters as much as the edge.

5. **Calibrate confidence like an analyst.** Attribution and correlation are rarely certain.
   A confirmed indicator is high confidence; an inferred actor link is mid; a single-source
   attribution claim is low. Store the uncertain ones — labeled honestly — so the assessment
   shows its work.

6. **Correlate and assess.** `analyze_risk_correlation` and `search_memories` by ATT&CK tag or
   actor surface patterns across sources. `get_related_entries` from an actor walks out to its
   indicators, campaigns, and techniques for a full profile.

## Enrich from the knowledge hub

`query_attack`, `query_nvd`, and `query_owasp` pull authoritative data you can attach to
records as `reference` memories. `map_to_cve` connects a CWE to CVEs an actor might exploit.
`get_security_data_summary` shows what hub data is loaded (run `download_security_data` if it's
empty).

## Producing an intelligence product

For a briefing or assessment, retrieve and traverse rather than recall: `get_context_summary`
for the overview, then walk the graph from the central actor/campaign to assemble indicators,
techniques (with ATT&CK mapping), and the evidence behind each assessment. Surface confidence
levels explicitly — a good intel product distinguishes what's confirmed from what's inferred.
