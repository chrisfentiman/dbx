---
name: dbx-analyze
description: End-to-end data exploration combining schema discovery, sampling, querying, and analysis. Use for open-ended questions like "what does our lead data look like?" or "explore the transactions table" or "how are these tables related?"
allowed-tools:
  - Bash(dbx *)
argument-hint: "[question or topic to explore]"
---

# Databricks Data Analysis

Conduct an end-to-end data exploration to answer a question or understand a dataset.

## Objective

$ARGUMENTS

## Analysis Protocol

### Phase 1: Discovery
Identify relevant tables for this analysis:
```bash
dbx schemas <catalog>
dbx tables <catalog>.<schema>
```
Scan table names and descriptions for relevance to the objective.

### Phase 2: Profiling
For each relevant table, build a data profile:
```bash
dbx describe <catalog>.<schema>.<table>
dbx sample <catalog>.<schema>.<table> --limit 10 --format table
dbx query "SELECT COUNT(*) as row_count FROM <catalog>.<schema>.<table>"
dbx query "SELECT <col>, COUNT(*) as cnt FROM <catalog>.<schema>.<table> GROUP BY <col> ORDER BY cnt DESC LIMIT 20"
```

Profile should include:
- Row count and key cardinality
- Column types and nullability
- Value distributions for categorical columns
- Min/max/avg for numeric columns
- Date ranges for temporal columns

### Phase 3: Relationship Mapping
Identify how tables connect:
- Look for matching column names across tables
- Check foreign key patterns (e.g., agent_id in multiple tables)
- Test joins to verify relationships:
```bash
dbx query "SELECT COUNT(*) FROM <table1> a JOIN <table2> b ON a.key = b.key"
```

### Phase 4: Analysis
Based on the profile, write targeted queries to answer the objective:
- Aggregations, joins, window functions as needed
- Always iterate: run query, interpret, refine

### Phase 5: Synthesis
Combine findings into a coherent narrative.

## Output Format

### Data Landscape
What tables exist and how they relate (include a relationship diagram if multiple tables).

### Key Metrics
Important numbers discovered during profiling.

### Analysis Findings
Detailed answers to the exploration objective.

### Data Quality Notes
Any issues found: nulls, duplicates, inconsistencies, missing data.

### Follow-Up Questions
What additional analysis would be valuable.
