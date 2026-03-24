# Code Quality Standards

## General Principles
- Write code that is correct first, clear second, fast third
- Prefer simple, obvious solutions over clever ones
- Every function should do one thing well
- Don't add abstractions until you need them in 3+ places
- Delete dead code — don't comment it out

## Git Workflow
- Commit messages: imperative mood, explain WHY not WHAT
- Keep commits focused on a single logical change
- Run tests before committing
- Never commit secrets, credentials, or API keys

## Dependencies
- Pin exact versions in production (`==`), use ranges in libraries (`>=`)
- Prefer well-maintained, widely-used libraries
- Evaluate new dependencies critically: is the dependency worth the cost?
- Keep dependencies minimal — stdlib is often enough

## Documentation
- Code should be self-documenting through clear naming
- Comments explain WHY, not WHAT
- Docstrings on public interfaces (functions, classes, modules)
- Keep README current with setup instructions and examples

## Security
- Never hardcode secrets — use environment variables or secret managers
- Validate and sanitize all external input
- Use parameterized queries for all database operations
- Keep dependencies updated for security patches
- Use HTTPS for all external API calls

## Project Organization
- Flat is better than nested (don't over-nest directories)
- Group by feature/domain, not by type (don't put all models in one folder)
- Keep test files close to source files
- Configuration in one place (`pyproject.toml` or `.env`)
