# dbx

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/chrisfentiman/dbx/actions/workflows/ci.yml/badge.svg)](https://github.com/chrisfentiman/dbx/actions/workflows/ci.yml)

> Read-only Databricks SQL CLI for safe data exploration, with built-in [Claude Code](https://claude.ai/claude-code) integration.

dbx wraps the Databricks SQL API with a two-layer SQL guard that blocks all write operations, a schema whitelist that restricts access to approved data, and compact output optimized for LLM-assisted analysis.

## Install

```bash
curl -sL https://raw.githubusercontent.com/chrisfentiman/dbx/main/install.sh | sh
```

## Quick Start

```bash
dbx up                # scaffold Claude Code skills + check dependencies
dbx config setup      # interactive Databricks connection setup
dbx doctor            # verify everything works
```

## Commands

### Setup
| Command | Description |
|---|---|
| `dbx up` | Scaffold Claude Code skills, rules, CLAUDE.md. Checks for uv and pyright. |
| `dbx config setup` | Interactive Databricks connection setup |
| `dbx config validate` | Test connectivity |
| `dbx doctor` | Verify config, connectivity, version, and binary integrity |
| `dbx update` | Self-update to latest release |

### Explore
| Command | Description |
|---|---|
| `dbx schemas <catalog>` | List allowed schemas |
| `dbx tables <catalog>.<schema>` | List tables in a schema |
| `dbx describe <catalog>.<schema>.<table>` | Show column names, types, nullability |
| `dbx sample <catalog>.<schema>.<table>` | Preview rows |
| `dbx preview <catalog>.<schema>.<table>` | Describe + sample combined |

### Query
| Command | Description |
|---|---|
| `dbx query "SELECT ..."` | Execute read-only SQL |
| `dbx query -f query.sql` | Execute SQL from a file |
| `dbx query "SELECT ..." -o results.csv` | Write output to a file |

### Flags
`--format json|csv|table` · `--limit N` · `--catalog <name>` · `--timeout <duration>`

## Configuration

`dbx config setup` creates `~/.dbx/config.yaml`:

```yaml
host: https://your-workspace.cloud.databricks.com
token: dapi...
warehouse_id: abc123
default_catalog: main
allowed_schemas:
  - my_schema
  - analytics
```

Environment variables override the config file: `DATABRICKS_HOST`, `DATABRICKS_TOKEN`, `DATABRICKS_WAREHOUSE_ID`, `DBX_ALLOWED_SCHEMAS`, `DBX_DEFAULT_CATALOG`.

## Claude Code Integration

`dbx up` scaffolds a project directory with skills for data exploration:

```
.claude/
├── skills/dbx-{analyze,config,explore,query,report}/
└── rules/{code-quality,geospatial,python}.md
CLAUDE.md
```

Files are diff-checked on every `dbx up` — only changed content is written. Run `dbx up` after updating dbx to get the latest skills.

## Security

- **Read-only** — two-layer SQL validation (keyword scan + AST parse) blocks all DML/DDL
- **Schema whitelist** — queries restricted to schemas in your config
- **Checksum verification** — `dbx update` verifies SHA256 before replacing the binary
- **Binary integrity** — `dbx doctor` compares your binary against the official release

## Development

```bash
make build    # compile to ./dbx-cli
make test     # run tests
make lint     # go vet
```

## License

MIT
