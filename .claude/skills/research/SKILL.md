---
name: research
description: Deep research on any topic using Perplexity, DeepWiki, and Context7. Use for comprehensive investigation of technologies, libraries, patterns, or domain questions.
allowed-tools:
  - mcp__perplexity__*
  - mcp__deepwiki__*
  - mcp__context7__*
  - WebSearch
  - WebFetch
  - Read
  - Glob
  - Grep
  - Write(docs/**)
---

# Deep Research

Conduct thorough research on the given topic using all available research tools.

## Research Target

$ARGUMENTS

## Research Protocol

### Phase 1: Broad Survey
1. Use `mcp__perplexity__search` for a quick landscape overview
2. Use `WebSearch` for recent articles, blog posts, and discussions
3. Identify key themes, tools, and approaches

### Phase 2: Deep Dive
1. Use `mcp__perplexity__reason` for complex comparisons and analysis
2. Use `mcp__context7__resolve-library-id` + `mcp__context7__query-docs` for library-specific documentation
3. Use `mcp__deepwiki__ask_question` for GitHub repository-specific questions
4. Use `WebFetch` to read specific articles or documentation pages found in Phase 1

### Phase 3: Synthesis
1. Cross-reference findings from multiple sources
2. Identify consensus vs. conflicting opinions
3. Note recency of information (prefer 2025-2026 sources)

## Output Format

Structure your findings as:

### Summary
2-3 sentence overview of key findings.

### Key Findings
Bulleted list of the most important discoveries, each with source attribution.

### Detailed Analysis
Organized by subtopic with evidence from multiple sources.

### Recommendations
Actionable next steps based on the research.

### Sources
List all sources consulted with URLs where available.

## Research Quality Standards
- Always cross-reference claims across at least 2 sources
- Prefer official documentation over blog posts
- Note when information may be outdated
- Flag conflicting information explicitly
- Distinguish between facts, opinions, and speculation
