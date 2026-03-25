# Python Standards

## Version and Tooling
- Target Python 3.11+
- Use `uv` for dependency management and script execution

## Dependency Management
- **Never assume packages are installed** — always verify before importing
- **Prefer stdlib** — `csv`, `json`, `collections`, `statistics`, `pathlib` cover most needs
- **Use PEP 723 inline metadata** for scripts that need dependencies:
  ```python
  # /// script
  # requires-python = ">=3.11"
  # dependencies = ["pandas", "matplotlib"]
  # ///
  ```
  Run with `uv run --script script.py` — deps install to a temp venv automatically
- **For quick inline scripts** that only need stdlib, use `python3` directly
- **Never `pip install` globally** — use `uv run` or virtual environments
- **If uv is not available**, fall back to `pip install --target /tmp/dbx-deps` and `PYTHONPATH=/tmp/dbx-deps`

## Scripts and Data Analysis
- Keep scripts in `workbench/` — not in the project root
- Read dbx query output via CSV: `dbx query "SELECT ..." --format csv -o /tmp/data.csv`
- For one-off analysis, use `python3 << 'EOF'` heredoc with stdlib
- For repeatable analysis, write a `.py` file with PEP 723 metadata in `workbench/`

## Code Style
- Use f-strings for string formatting
- Use `pathlib.Path` instead of `os.path`
- Catch specific exceptions, never bare `except:`
- Never use wildcard imports (`from module import *`)
- Alias conventions: `pd`, `np`, `gpd`, `plt`
