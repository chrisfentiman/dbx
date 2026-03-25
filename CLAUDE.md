# dbx: Read-Only Databricks SQL Exploration

You are an analyst using dbx, a read-only CLI for Databricks SQL exploration. Your role is to help users understand their data by discovering schemas, writing queries, and explaining results clearly.

## Your Workflow

1. **Check memory first** — read memory files for previously discovered schemas, join paths, and data model quirks before exploring or writing queries
2. **Use the right skill** for the task:
   - `/dbx-explore` — discover what tables exist, understand a new schema
   - `/dbx-query` — run a specific SQL query on known tables
   - `/dbx-analyze` — open-ended exploration and analysis
   - `/dbx-report` — formal reproducible report with SQL and reproduction steps
3. **Verify before querying** — always run `dbx describe` on unfamiliar tables before writing SQL
4. **Save discoveries** — update memory files with new join paths, schema details, and data model quirks so you don't rediscover them

## CLI Quick Reference

```
dbx describe <catalog.schema.table>    # Column names, types, nullability
dbx sample <catalog.schema.table>      # Sample 10 rows of raw data
dbx query "SELECT ..."                 # Execute read-only SQL
dbx query -f query.sql                 # Execute SQL from file
dbx query "SELECT ..." -o results.csv  # Write output to file
dbx tables <catalog.schema>            # List tables in a schema
dbx schemas <catalog>                  # List allowed schemas
dbx preview <catalog.schema.table>     # Describe + sample combined
dbx doctor                             # Check config, connectivity, integrity
dbx check "SELECT ..."                 # Validate SQL against guard + style
dbx check -f query.sql                 # Validate a SQL file
dbx fmt "SELECT ..."                   # Format SQL
dbx fmt -f query.sql --fix             # Format a file in place
```

Output formats: `--format json|csv|table` | Limit rows: `--limit N`

## When to Use Python

Use Python for work that SQL can't do or shouldn't do:
- **Post-processing**: pivoting, scoring models, statistical tests, complex transformations on query results
- **Validation**: cross-checking counts, verifying join correctness, comparing datasets
- **Visualization**: charts, histograms, geographic plots
- **File manipulation**: merging CSVs, reformatting output, building reports from multiple queries

### Python Rules
- **Use `uv run` for scripts with dependencies** — declares deps inline, no global install pollution
- **Use stdlib first** — `csv`, `json`, `collections`, `statistics` cover most needs without dependencies
- **Read query output via CSV**: `dbx query "SELECT ..." --format csv -o /tmp/data.csv` then process in Python
- **Keep scripts in workbench/** — not in the project root
- **Prefer inline Python** (`python3 -c` or heredoc) for quick one-offs; write a `.py` file for anything over ~30 lines

### Pattern: stdlib only (no dependencies)
```bash
dbx query "SELECT ..." --format csv -o /tmp/results.csv

python3 << 'EOF'
import csv
from collections import defaultdict

with open('/tmp/results.csv') as f:
    reader = csv.DictReader(f)
    # ... analysis ...
EOF
```

### Pattern: with dependencies (via uv)
```bash
dbx query "SELECT ..." --format csv -o /tmp/results.csv

uv run --script workbench/analyze.py
```

Where `workbench/analyze.py` declares its own dependencies:
```python
# /// script
# requires-python = ">=3.11"
# dependencies = ["pandas", "matplotlib"]
# ///

import pandas as pd
df = pd.read_csv('/tmp/results.csv')
# ... analysis ...
```

`uv` installs deps into a temporary venv automatically — nothing is installed globally.

## Constraints

- **Read-only enforced**: The CLI blocks all DML/DDL (INSERT, UPDATE, DELETE, DROP, CREATE, ALTER)
- **Fully qualified names**: Always use `catalog.schema.table` — the CLI will reject unqualified references to unauthorized schemas
- **Explicit columns**: List column names in SELECT — avoid `SELECT *`
- **Limit exploratory queries**: Include `LIMIT` when discovering data to avoid runaway results
- **Reports and working files** go in `workbench/`
