# Python Standards

## Version and Tooling
- Target Python 3.11+ (use modern syntax: `match`, `type` aliases, `ExceptionGroup`)
- Use `uv` for dependency management and script execution
- Use `ruff` for linting and formatting (not black, isort, flake8 separately)
- Use `mypy` for type checking with strict mode
- Use `pytest` for testing (not unittest)
- Use `pyproject.toml` for project configuration (not setup.py/setup.cfg)

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
- **Never `pip install` globally** — use `uv run`, `uv pip install --python`, or virtual environments
- **If uv is not available**, fall back to `pip install --target /tmp/dbx-deps` and `PYTHONPATH=/tmp/dbx-deps`

## Scripts and Data Analysis
- Keep scripts in `workbench/` — not in the project root
- Read dbx query output via CSV: `dbx query "SELECT ..." --format csv -o /tmp/data.csv`
- For one-off analysis, use `python3 << 'EOF'` heredoc with stdlib
- For repeatable analysis, write a `.py` file with PEP 723 metadata in `workbench/`

## Code Style
- Type hints on all public function signatures
- Docstrings on all public functions and classes (Google style)
- Use `pathlib.Path` instead of `os.path`
- Use f-strings for string formatting
- Use `dataclasses` or `pydantic` for data structures, not plain dicts
- Prefer `enum.Enum` for finite sets of choices
- Use `logging` module, never `print()` for production code

## Imports
- Group imports: stdlib, third-party, local
- Use absolute imports for cross-module references
- Never use wildcard imports (`from module import *`)
- Alias conventions: `pd`, `np`, `gpd`, `plt`

## Error Handling
- Catch specific exceptions, never bare `except:`
- Use `raise ... from e` for exception chaining
- Validate at system boundaries (API inputs, file reads, external data)
- Don't over-validate internal function calls

## Testing
- Test files mirror source structure: `src/foo/bar.py` -> `tests/foo/test_bar.py`
- Use fixtures for shared setup
- Parametrize tests with `@pytest.mark.parametrize`
- Aim for meaningful assertions, not coverage numbers
- Name tests: `test_<function>_<scenario>_<expected_outcome>`
