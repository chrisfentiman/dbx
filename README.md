# dbx

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/chrisfentiman/dbx/actions/workflows/ci.yml/badge.svg)](https://github.com/chrisfentiman/dbx/actions/workflows/ci.yml)

> Read-only Databricks SQL CLI for safe data exploration, with built-in [Claude Code](https://claude.ai/claude-code) integration.

dbx wraps the Databricks SQL API with a security-first approach: a two-layer SQL guard blocks all write operations, a schema whitelist restricts access to approved data, and compact JSON output is optimized for LLM-assisted analysis. Run `dbx up` in any directory to scaffold Claude Code skills for schema discovery, querying, and report generation.

## Table of Contents

- [Install](#install)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Configuration](#configuration)
- [Claude Code Integration](#claude-code-integration)
- [Security](#security)
- [Development](#development)
- [License](#license)

## Install

### One-line install (recommended)

```bash
curl -sL https://raw.githubusercontent.com/chrisfentiman/dbx/main/install.sh | sh
```

### Manual download

Download from [Releases](https://github.com/chrisfentiman/dbx/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/chrisfentiman/dbx/releases/latest/download/dbx-darwin-arm64.tar.gz | tar xz
sudo mv dbx-darwin-arm64 /usr/local/bin/dbx

# macOS (Intel)
curl -L https://github.com/chrisfentiman/dbx/releases/latest/download/dbx-darwin-amd64.tar.gz | tar xz
sudo mv dbx-darwin-amd64 /usr/local/bin/dbx

# Linux (x86_64)
curl -L https://github.com/chrisfentiman/dbx/releases/latest/download/dbx-linux-amd64.tar.gz | tar xz
sudo mv dbx-linux-amd64 /usr/local/bin/dbx
```

### From source

```bash
git clone https://github.com/chrisfentiman/dbx.git
cd dbx
make build
mv dbx-cli /usr/local/bin/dbx
```

## Quick Start

```bash
dbx up                # scaffold Claude Code skills in current directory
dbx config setup      # interactive Databricks connection setup
dbx config validate   # test connectivity
dbx schemas main      # list schemas
dbx query "SELECT * FROM main.my_schema.my_table LIMIT 10"
```

## Commands

| Command | Description |
|---|---|
| `dbx up [--force]` | Scaffold current directory with skills and config |
| `dbx config setup` | Interactive connection setup |
| `dbx config validate` | Test connectivity |
| `dbx config` | Show current config |
| `dbx catalogs` | List allowed catalogs |
| `dbx schemas <catalog>` | List allowed schemas |
| `dbx tables <catalog>.<schema>` | List tables |
| `dbx describe <table>` | Column details |
| `dbx sample <table>` | Sample rows |
| `dbx preview <table>` | Describe + sample |
| `dbx query "SQL"` | Execute read-only SQL |
| `dbx version` | Print version |

**Flags:** `--format json|csv|table` `--limit N` `--catalog <name>` `--timeout <duration>`

## Configuration

Stored in `~/.dbx/config.yaml`:

```yaml
host: https://your-workspace.cloud.databricks.com
token: dapi...
warehouse_id: abc123
default_catalog: main
allowed_schemas:
  - silver_analytics
  - gold_reporting
default_format: json
query_timeout: 30s
```

Environment variables override the config file:

| Variable | Description |
|---|---|
| `DATABRICKS_HOST` | Workspace URL |
| `DATABRICKS_TOKEN` | Personal access token |
| `DATABRICKS_WAREHOUSE_ID` | SQL warehouse ID |
| `DBX_ALLOWED_SCHEMAS` | Comma-separated schema whitelist |
| `DBX_DEFAULT_CATALOG` | Default catalog |

## Claude Code Integration

`dbx up` scaffolds a project directory for Claude Code:

```
.claude/
├── skills/
│   ├── dbx-analyze/     # open-ended data exploration
│   ├── dbx-config/      # connection verification
│   ├── dbx-explore/     # schema and table discovery
│   ├── dbx-query/       # SQL query execution
│   └── dbx-report/      # reproducible report generation
└── rules/               # code quality standards
```

- **Skills (`dbx-*`)** are always updated on `dbx up`
- **Rules** are only written if they don't already exist
- **`CLAUDE.md`** is skipped if it exists (use `--force` to overwrite)

## Security

- **Read-only enforcement** — two-layer SQL validation (keyword scan + AST parse) blocks all DML/DDL
- **Schema whitelist** — queries can only reference schemas in your config
- **Credentials isolated** — stored in `~/.dbx/config.yaml`, never in the project directory

## Development

```bash
make build    # compile to ./dbx-cli
make test     # run tests
make lint     # go vet
make clean    # remove binary
```

## License

MIT