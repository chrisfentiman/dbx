---
name: dbx-report
description: Generate a reproducible data report with SQL queries, results, and step-by-step reproduction instructions. Use when the user wants a formal report they can re-run themselves.
allowed-tools:
  - Bash(dbx *)
  - Write
argument-hint: "[report topic or question]"
---

# Databricks Data Report

Generate a reproducible data report answering a specific question.

## Report Topic

$ARGUMENTS

## Report Protocol

### Step 1: Plan Queries
Determine what queries are needed to answer the question. List them before executing.

### Step 2: Execute and Capture
Run each query, capturing both the SQL and results:
```bash
dbx query "<sql>" --format json
```

### Step 3: Analyze
Interpret each result set and build the narrative.

### Step 4: Write Report
Write a markdown report to docs/reports/ with this structure:

# Report: <Title>
**Generated**: <date>
**Data Source**: Databricks
**Warehouse**: (from dbx config)

## Executive Summary
2-3 sentence overview of findings.

## Methodology
Describe the approach and any assumptions.

## Queries and Results

### Query 1: <description>
The exact SQL used, the results summary, and interpretation.

### Query 2: <description>
...

## Findings
Detailed analysis and conclusions.

## How to Reproduce
Step-by-step instructions for the user to run these queries themselves:
1. Ensure dbx CLI is configured (run dbx config validate)
2. Run each query above in sequence
3. Expected data freshness: results may vary if underlying data has been updated

Save the report to: docs/reports/<topic-slug>-<date>.md
