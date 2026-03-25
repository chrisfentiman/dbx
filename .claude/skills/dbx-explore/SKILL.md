---
name: dbx-explore
description: Explore Databricks catalogs, schemas, tables, and columns. Use when discovering what data is available, understanding table structures, or finding relevant tables for analysis.
allowed-tools:
  - Bash(dbx *)
argument-hint: "[catalog or schema or table name]"
---

# Databricks Schema Explorer

Systematically explore Databricks to discover and understand available data.

## Target

$ARGUMENTS

## Constraints

- Check memory first — do not rediscover schemas already documented
- Always describe tables before making assumptions about columns
- Save all discoveries to memory after exploring

## Exploration Protocol

### 1. Check Memory
Before exploring, check if this schema or table has been previously documented in memory files. If already known, summarize from memory instead of re-running commands.

### 2. Discover Structure
Start broad, then drill down:
```bash
dbx catalogs                        # What catalogs exist?
dbx schemas <catalog>               # What schemas are available?
dbx tables <catalog>.<schema>       # What tables are in this schema?
```

### 3. Describe Tables
For each relevant table, get the full column definition:
```bash
dbx describe <catalog>.<schema>.<table>
```
Note: column names, data types, NULL constraints.

### 4. Sample Data
Preview actual data to understand content and quality:
```bash
dbx sample <catalog>.<schema>.<table>
```
Look at 5-10 rows. Check for: NULL patterns, date formats, ID formats, unexpected values.

### 5. Identify Relationships
After describing 2-3 tables, reason about joins:
- Do any ID columns appear in multiple tables?
- Are there naming conventions? (e.g., `table_id` pattern)
- What is the likely cardinality? (one-to-many, many-to-one)

### 6. Save to Memory
Update memory files with:
- Schema inventory (tables, row counts, key columns)
- Verified join paths with cardinality
- Column details for important tables
- Any data model quirks or edge cases

## Example: Exploring a New Schema

**User**: "What's in the analytics schema?"

**Step 1**: Check memory — analytics not previously explored.

**Step 2**: List tables:
```
dbx tables main.analytics
→ customers, orders, products, reviews, shipments
```

**Step 3**: Describe key tables:
```
dbx describe main.analytics.customers
→ customer_id (INT), email (VARCHAR), created_date (DATE)

dbx describe main.analytics.orders
→ order_id (INT), customer_id (INT), order_date (DATE), amount (DECIMAL)
```

**Step 4**: Sample to verify:
```
dbx sample main.analytics.customers
→ Real data, no NULLs in customer_id, email looks valid
```

**Step 5**: Identify joins:
- `orders.customer_id` → `customers.customer_id` (many-to-one)

**Step 6**: Save to memory — schema inventory, join path, column details.

## Output Format

### Available Data
List of catalogs/schemas discovered

### Table Inventory
| Table | Key Columns | Notes |
|-------|------------|-------|

### Relationships
Documented join paths with cardinality

### Recommendations
Which tables are most relevant to the user's question and why
