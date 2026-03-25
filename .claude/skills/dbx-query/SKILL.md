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

## Constraints

- SELECT queries ONLY — no INSERT, UPDATE, DELETE, CREATE, DROP, ALTER
- Always fully qualify table names: `catalog.schema.table`
- Always list columns explicitly — no `SELECT *`
- Include `LIMIT` for exploratory queries
- Use CTEs for complex multi-step queries

## Verify-Query-Validate Cycle

### 1. Plan (think before writing SQL)
Before writing a query, reason through:
- What tables do I need?
- Do I know the column names, or do I need to verify with `dbx describe`?
- If joining tables, what is the join path and cardinality?
- Check memory for previously discovered join paths

### 2. Verify (non-optional for unfamiliar tables)
Run `dbx describe` on any table you're not certain about:
```bash
dbx describe <catalog>.<schema>.<table>
```
Review output for exact column names, data types, and NULL constraints.

### 3. Query
Construct and execute:
```bash
dbx query "SELECT col1, col2 FROM catalog.schema.table WHERE condition ORDER BY col1 LIMIT 100"
```
For SQL in a file: `dbx query -f query.sql`
To save output: `dbx query "SELECT ..." --format csv -o results.csv`

### 4. Validate
Before accepting results:
- Is the row count reasonable (> 0 and not suspiciously high)?
- Do sample values look like real data?
- If joining, are there unexpected duplicates (wrong join cardinality)?

### 5. Iterate if needed
If results are unexpected:
1. Check column types with `dbx describe`
2. Sample raw data with `dbx sample`
3. Adjust query and re-execute

## Example: Simple Aggregation

**User asks**: "How many orders per month in 2025?"

**Plan**: Need the orders table. I know from memory that `main.analytics.orders` has `order_id` and `order_date` columns.

**Query**:
```sql
SELECT
  DATE_TRUNC('month', order_date) AS month,
  COUNT(order_id) AS order_count
FROM main.analytics.orders
WHERE order_date >= '2025-01-01'
GROUP BY DATE_TRUNC('month', order_date)
ORDER BY month
```

**Validate**: 12 rows returned (Jan-Dec), counts look reasonable, no NULLs.

## Example: Multi-Table Join

**User asks**: "Show customer names with their total order value"

**Plan**: Need customers and orders. Check memory for join path — `orders.customer_id = customers.customer_id` (many-to-one). Need to aggregate orders.

**Verify**: Run `dbx describe` on both tables to confirm column names.

**Query**:
```sql
SELECT
  c.customer_id,
  c.name,
  SUM(o.order_amount) AS total_value
FROM main.analytics.customers c
JOIN main.analytics.orders o ON c.customer_id = o.customer_id
GROUP BY c.customer_id, c.name
ORDER BY total_value DESC
LIMIT 20
```

**Validate**: 20 rows, values look reasonable, no duplicate customer_ids (GROUP BY working correctly).

## Output Format

### Query
```sql
-- The SQL that was executed
```

### Results Summary
Plain-English interpretation of what the data shows.

### Raw Results
The actual query output.
