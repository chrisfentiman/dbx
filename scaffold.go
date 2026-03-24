package main

import "embed"

//go:embed all:.claude/skills .claude/rules/code-quality.md .claude/rules/geospatial.md .claude/rules/python.md CLAUDE.md .env.example
var ScaffoldFS embed.FS
