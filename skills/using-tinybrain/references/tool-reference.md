# TinyBrain tool reference

All TinyBrain MCP tools, grouped by purpose. Tools may appear with a server prefix
(`mcp__tinybrain__<name>`); the logical names are below. Several tools have aliases that
behave identically — both names are noted.

## Sessions

| Tool | Required args | Use |
|---|---|---|
| `create_session` | `name`, `task_type` | Start an engagement. Optional: `description`, `intelligence_type`, `classification`, `threat_level`, `metadata`. |
| `get_session` | `session_id` | Fetch one session. |
| `list_sessions` | — | List sessions; optional `task_type`, `status`, `limit`. Call this to orient. |

## Memory: create / read / update / delete

| Tool | Required args | Use |
|---|---|---|
| `create_memory` (alias `store_memory`) | `session_id`, `title`, `content`, `category` | Store a finding. Optional: `priority`, `confidence`, `tags`, `source`, `content_type`, `intelligence_type`, `threat_level`, `mitre_tactic`, `mitre_technique`. |
| `get_memory` | `memory_id` | Fetch one memory (also bumps its access tracking). |
| `update_memory` | `memory_id` | Update fields: `title`, `content`, `category`, `priority`, `confidence`, `tags`, `source`. |
| `delete_memory` | `memory_id` | Delete one memory. |
| `get_detailed_memory_info` | `memory_id` | Memory plus access/debug metadata. |

## Search & retrieval

| Tool | Required args | Use |
|---|---|---|
| `search_memories` (alias `search_memory`) | `query` | Keyword/FTS search. Optional `session_id`, `category`, `categories`, `tags`, `search_type`, `min_priority`, `min_confidence`, `limit`. |
| `semantic_search` | `query`, `session_id` | Conceptual similarity search. |
| `find_similar_memories` | `session_id`, `content` | Find memories similar to a block of content. |
| `check_duplicates` | `session_id`, `title`, `content` | Pre-store dedup check. |
| `get_context_summary` | `session_id` | Session-level digest; call when resuming. |

## Relationships

| Tool | Required args | Use |
|---|---|---|
| `create_relationship` | `source_memory_id`, `target_memory_id`, `relationship_type` | Link two memories. Optional `strength`, `description`. |
| `get_related_entries` (alias `get_related_memories`) | `memory_id` | Walk the graph from a node. Optional `relationship_type`, `limit`. |

## Context snapshots

| Tool | Required args | Use |
|---|---|---|
| `create_context_snapshot` | `session_id`, `name` | Save working state at a checkpoint. Optional `description`, `context_data`. |
| `get_context_snapshot` | `snapshot_id` | Restore a snapshot. |
| `list_context_snapshots` | `session_id` | List a session's snapshots. |

## Task progress

| Tool | Required args | Use |
|---|---|---|
| `create_task_progress` | `session_id`, `task_name`, `stage`, `status` | Start tracking a multi-stage task. Optional `progress_percentage`, `notes`. |
| `update_task_progress` | `session_id`+`task_name` (or `task_id`), `stage`, `status` | Advance a task. |
| `get_task_progress` | `task_id` | Fetch one task. |
| `list_task_progress` | `session_id` | List a session's tasks; optional `status`. |

## Batch & lifecycle

| Tool | Required args | Use |
|---|---|---|
| `batch_create_memories` | `session_id`, `memory_requests` | Store many at once. |
| `batch_update_memories` | `memory_updates` | Update many at once. |
| `batch_delete_memories` | `memory_ids` | Delete many at once. |
| `cleanup_old_memories` | `max_age_days` | Prune by age (supports dry-run). |
| `cleanup_low_priority_memories` | `max_priority`, `max_confidence` | Prune low-value memories. |
| `cleanup_unused_memories` | `max_unused_days` | Prune memories never accessed. |

## Templates & analysis

| Tool | Required args | Use |
|---|---|---|
| `get_security_templates` | — | Predefined memory templates. |
| `create_memory_from_template` | `session_id`, `template_name` | Instantiate a template. |
| `analyze_risk_correlation` | `session_id` | Correlate vulnerabilities into risk chains. |
| `map_to_cve` | `session_id`, `cwe_id` | Map a CWE to related CVEs. |
| `map_to_compliance` | `session_id`, `standard` | Map findings to a compliance standard. |

## Security knowledge hub

| Tool | Required args | Use |
|---|---|---|
| `query_nvd` | — | Query the NVD CVE dataset. |
| `query_attack` | — | Query MITRE ATT&CK. |
| `query_owasp` | — | Query OWASP testing data. |
| `download_security_data` | — | Populate the knowledge-hub datasets. |
| `get_security_data_summary` | — | Show what hub data is loaded. |

## Notifications & ops

| Tool | Required args | Use |
|---|---|---|
| `get_notifications` | `session_id` | Alerts (high-priority finds, duplicates, cleanups). |
| `mark_notification_read` | `notification_id` | Dismiss a notification. |
| `check_high_priority_memories` | `session_id` | Surface high-priority items. |
| `check_duplicate_memories` | `session_id` | Surface likely duplicates. |
| `get_memory_stats` / `get_database_stats` / `get_system_diagnostics` | — | Stats and health. |
| `health_check` | — | Liveness/DB check. |
| `export_session_data` / `import_session_data` | `session_id` / `import_data` | Portability. |

## Embeddings (semantic plumbing)

| Tool | Required args | Use |
|---|---|---|
| `generate_embedding` | `text` | Generate an embedding for text. |
| `calculate_similarity` | `embedding1`, `embedding2` | Cosine similarity between two embeddings. |

Note: the embedding implementation is currently lightweight; prefer `search_memories` and
`semantic_search` for day-to-day retrieval.
