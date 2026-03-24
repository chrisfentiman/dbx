---
name: quick-search
description: Quick search for a specific question or fact check using Perplexity. Use for fast lookups, version checks, or quick answers.
allowed-tools:
  - mcp__perplexity__search
  - WebSearch
---

# Quick Search

Get a fast answer to a specific question.

**Question:** $ARGUMENTS

Use `mcp__perplexity__search` to find the answer. Be concise — return the direct answer with source attribution. If the question requires more depth, suggest using `/research` instead.
