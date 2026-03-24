---
name: dbx-config
description: Check Databricks CLI configuration and connectivity. Use to verify dbx setup before running queries.
disable-model-invocation: true
allowed-tools:
  - Bash(dbx *)
---

# Databricks Config Check

Verify the dbx CLI is properly configured and can connect to Databricks.

## Steps

1. Show current configuration:
```bash
dbx config
```

2. Validate connectivity:
```bash
dbx config validate
```

Report the configuration status and whether the connection is working.
