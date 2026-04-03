# Contributing to ZarishSphere GitHub Bot

## Development Setup

### Prerequisites
- Go 1.21 or later
- GitHub Personal Access Token with repo permissions

### Setup
```bash
git clone https://github.com/zarishsphere/zs-ops-github-bot.git
cd zs-ops-github-bot
go mod download
make build
```

### Testing
```bash
make test
make test-coverage
```

## Code Standards

### Go Code
- Follow standard Go formatting (`go fmt`)
- Use `golangci-lint` for code quality
- Write tests for all new functionality
- Use meaningful variable and function names

### Commits
- Use conventional commit format
- Keep commits focused and atomic
- Write clear commit messages

### Pull Requests
- Create feature branches from `main`
- Ensure CI passes before requesting review
- Provide clear description of changes
- Update documentation as needed

## Security
- Never commit secrets or tokens
- Use environment variables for configuration
- Follow principle of least privilege
- Report security issues privately

## License
By contributing, you agree to license your contributions under the same license as the project (Apache 2.0).