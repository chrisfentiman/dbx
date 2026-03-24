---
name: dbx-query
description: Execute read-only SQL queries against Databricks. Use when you need to run a specific SQL query, aggregate data, join tables, or answer analytical questions about Databricks data.
allowed-tools:
  - Bash(dbx *)
argument-hint: "[SQL or question about the data]"
---

# Databricks SQL Query

Execute a read-only SQL query against Databricks and interpret the results.

## Task

$ARGUMENTS

## Query Protocol

### Step 1: Understand the Schema
Before writing SQL, verify table structure if you haven't already:
```bash
dbx describe <catalog>.<schema>.<table>
```

### Step 2: Write and Execute SQL
Construct a SQL query that answers the question. Rules:
- SELECT only (no INSERT, UPDATE, DELETE, DROP, CREATE, ALTER)
- Always fully qualify table names: catalog.schema.table
- Use explicit column names (avoid SELECT *)
- Add ORDER BY and LIMIT for large result sets
- Use CTEs for complex queries (readability)

Execute:
```bash
dbx query "SELECT col1, col2 FROM catalog.schema.table WHERE condition ORDER BY col1 LIMIT 100"
```

### Step 3: Iterate if Needed
If results are unexpected:
1. Check column types with dbx describe
2. Sample raw data with dbx sample
3. Adjust query and re-execute

### Step 4: Interpret Results
- Summarize the key findings from the data
- Call out any anomalies or unexpected patterns
- If the data suggests follow-up questions, note them

## Output Format

### Query
```sql
-- The SQL that was executed
```

### Results Summary
Plain-English interpretation of what the data shows.

### Raw Results
The actual query output.
