# Data Explorer

This project provides controlled SQL access to Databricks through a Go CLI (`dbx`) wrapped by Claude Code skills.

## Using the Data Explorer

### Always Check Memory First

Before exploring schemas, writing joins, or looking up table relationships, **read your memory files** for previously discovered join paths, ID mappings, table inventories, and data model quirks. Do NOT rediscover what's already known — it wastes context and time.

### When to Use Each Skill

| Skill | When to Use |
|---|---|
| `/dbx-config` | Verify setup before first query of a session |
| `/dbx-explore` | Discovering what tables exist, understanding a new schema you haven't seen before |
| `/dbx-query` | Running a specific SQL query when you already know the tables and joins |
| `/dbx-analyze` | Open-ended exploration: "what does this data look like?", "how are these tables related?" |
| `/dbx-report` | Producing a formal reproducible report with SQL, results, and reproduction instructions |

### Skill Selection Guidelines

- If the question is about **known tables with known joins** (check memory first), use `/dbx-query` directly
- If the question involves **tables you haven't seen before**, use `/dbx-explore` first, then `/dbx-query`
- If the user wants a **saved report**, always finish with `/dbx-report`
- If you need to **map relationships** between unfamiliar tables, use `/dbx-analyze`

### Workflow for Data Questions

1. **Check memory** — look for known join paths and table info
2. **Use the right skill** — don't manually run `dbx` commands outside of skills unless debugging
3. **Save new discoveries** — when you find a new join path, table relationship, or ID mapping, update your memory

## Updating Memory

When you discover something new about the data model:

- A new join path between tables
- A new table or schema
- An ID mapping or naming convention
- A data model quirk (e.g., "user_id changes format between systems")
- A project/program and which tables/filters apply to it (e.g., `form_type = 'xyz'`)

**Save it to memory immediately.** Keep it structured — add to the appropriate section rather than appending randomly.

## CLI Reference

```bash
dbx catalogs                              # List allowed catalogs
dbx schemas <catalog>                      # List allowed schemas
dbx tables <catalog>.<schema>              # List tables in a schema
dbx describe <catalog>.<schema>.<table>    # Column details
dbx sample <catalog>.<schema>.<table>      # Sample rows (default 10)
dbx query "SELECT ..."                     # Execute read-only SQL
dbx preview <catalog>.<schema>.<table>     # Describe + sample combined
dbx config                                 # Show config
dbx config validate                        # Test connectivity
dbx config setup                           # Interactive setup
```

Flags: `--format json|table|csv`, `--limit N`, `--catalog <name>`, `--timeout <seconds>`

## Important Constraints

- **Read-only**: The SQL guard blocks all DML/DDL (INSERT, UPDATE, DELETE, DROP, CREATE, ALTER)
- **Schema whitelist**: Queries are restricted to the schemas configured in `~/.dbx/config.yaml`
- **JSON output is compact**: No pretty-printing — optimized for LLM token efficiency
- Reports go in `workbench/`
