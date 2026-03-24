---
name: deep-research
description: In-depth research and analysis on complex topics using Perplexity Deep Research. Use for comprehensive reports, technology evaluations, and multi-faceted questions.
context: fork
agent: general-purpose
allowed-tools:
  - mcp__perplexity__deep_research
  - mcp__perplexity__reason
  - mcp__context7__*
  - mcp__deepwiki__*
  - WebSearch
  - WebFetch
  - Read
  - Write(docs/**)
---

# Deep Research Report

Conduct an exhaustive research report on the following topic:

**Topic:** $ARGUMENTS

## Process

1. Start with `mcp__perplexity__deep_research` for comprehensive analysis
2. Use `mcp__perplexity__reason` to analyze complex sub-questions identified in step 1
3. Cross-reference with `mcp__context7__query-docs` for any library-specific claims
4. Use `mcp__deepwiki__ask_question` for any GitHub repo-specific questions
5. Fetch key source documents with `WebFetch` for detailed verification

## Output Requirements

Produce a structured research report with:

1. **Executive Summary** — Key findings in 3-5 sentences
2. **Background** — Context and why this topic matters
3. **Analysis** — Detailed findings organized by subtopic
4. **Comparison Matrix** — If evaluating alternatives, provide a comparison table
5. **Recommendations** — Specific, actionable recommendations with rationale
6. **Risks and Considerations** — What to watch out for
7. **Sources** — All sources with URLs

Write the report to `docs/research/` with a descriptive filename.
