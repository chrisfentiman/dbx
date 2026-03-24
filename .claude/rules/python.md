# Python Standards

## Version and Tooling
- Target Python 3.11+ (use modern syntax: `match`, `type` aliases, `ExceptionGroup`)
- Use `ruff` for linting and formatting (not black, isort, flake8 separately)
- Use `mypy` for type checking with strict mode
- Use `pytest` for testing (not unittest)
- Use `pyproject.toml` for project configuration (not setup.py/setup.cfg)

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
