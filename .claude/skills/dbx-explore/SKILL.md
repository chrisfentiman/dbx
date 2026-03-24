---
name: dbx-explore
description: Explore Databricks catalogs, schemas, tables, and columns. Use when discovering what data is available, understanding table structures, or finding relevant tables for analysis.
allowed-tools:
  - Bash(dbx *)
argument-hint: "[catalog or schema or table name]"
---

# Databricks Schema Explorer

Explore the Databricks Unity Catalog to discover and understand available data.

## Target

$ARGUMENTS

## Exploration Protocol

### Step 1: Discover Available Schemas
If no specific table is mentioned, start broad:
```bash
dbx catalogs
```
Then drill into a catalog:
```bash
dbx schemas <catalog>
```

### Step 2: Find Relevant Tables
Once you know the schema:
```bash
dbx tables <catalog>.<schema>
```

### Step 3: Understand Table Structure
For each relevant table:
```bash
dbx describe <catalog>.<schema>.<table>
```

### Step 4: Preview Data
Get a small sample to understand the data:
```bash
dbx sample <catalog>.<schema>.<table> --limit 5 --format table
```

## Output Format

Summarize findings as:

### Available Data
- List of relevant catalogs/schemas discovered

### Table Inventory
For each relevant table:
| Table | Description | Key Columns | Row Sample |
|-------|------------|-------------|------------|

### Relationships
Note any foreign key patterns or join opportunities between tables (matching column names, ID references, etc.).

### Recommendations
Suggest which tables are most relevant to the user's question and why.
