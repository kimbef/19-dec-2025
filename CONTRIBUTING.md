# Contributing to CI/CD Pipeline Application

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to this project.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style Guidelines](#code-style-guidelines)
- [Testing Requirements](#testing-requirements)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [CI/CD Integration Points](#cicd-integration-points)

## 📜 Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow:

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards other community members

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed:
- Go 1.20 or higher
- Docker and Docker Compose
- Make
- Git
- golangci-lint (install with `make install-tools`)

### Setting Up Your Development Environment

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/19-dec-2025.git
   cd 19-dec-2025
   ```

2. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/kimbef/19-dec-2025.git
   ```

3. **Install dependencies**
   ```bash
   make deps
   make install-tools
   ```

4. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

5. **Start development database**
   ```bash
   docker-compose up -d db
   ```

6. **Run the application**
   ```bash
   make run
   ```

## 🔄 Development Workflow

### Creating a Feature Branch

```bash
# Update your main branch
git checkout main
git pull upstream main

# Create a feature branch
git checkout -b feature/your-feature-name
```

### Making Changes

1. **Write code** following our [Code Style Guidelines](#code-style-guidelines)
2. **Write tests** for your changes
3. **Run tests locally**
   ```bash
   make test
   make test-coverage
   ```
4. **Run linter**
   ```bash
   make lint
   make fmt
   ```
5. **Ensure code builds**
   ```bash
   make build
   ```

### Keeping Your Branch Updated

```bash
# Fetch latest changes
git fetch upstream

# Rebase your branch
git rebase upstream/main
```

## 🎨 Code Style Guidelines

### Go Conventions

Follow standard Go conventions and best practices:

1. **Formatting**
   - Use `gofmt` for formatting (run `make fmt`)
   - Lines should be ≤120 characters
   - Use tabs for indentation

2. **Naming**
   - Use camelCase for variables and functions
   - Use PascalCase for exported types and functions
   - Use ALL_CAPS for constants
   - Package names should be lowercase, single-word

3. **Error Handling**
   - Always check and handle errors
   - Wrap errors with context using `fmt.Errorf`
   - Use custom error types for domain-specific errors

4. **Comments**
   - Document all exported functions, types, and constants
   - Use complete sentences in comments
   - Start comments with the name of the thing being described

Example:
```go
// Handler contains dependencies for API handlers.
// It provides methods for handling HTTP requests and responses.
type Handler struct {
    db  *database.DB
    log *logger.Logger
}

// NewHandler creates a new API handler with the given dependencies.
// It returns a properly initialized Handler instance.
func NewHandler(db *database.DB, log *logger.Logger) *Handler {
    return &Handler{
        db:  db,
        log: log,
    }
}
```

5. **Project Structure**
   - `cmd/`: Application entry points
   - `internal/`: Private application code
   - `pkg/`: Public library code
   - Keep packages small and focused
   - Avoid circular dependencies

## ✅ Testing Requirements

### Writing Tests

1. **Test Coverage**
   - Aim for >80% code coverage
   - Test both success and failure cases
   - Include edge cases and boundary conditions

2. **Test Organization**
   - Place tests in `*_test.go` files
   - Use table-driven tests for multiple scenarios
   - Name tests descriptively: `TestFunctionName_Scenario`

Example:
```go
func TestCreatePipeline_Success(t *testing.T) {
    // Arrange
    // ... setup code

    // Act
    result, err := CreatePipeline(req)

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if result.Name != expected.Name {
        t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
    }
}
```

3. **Integration Tests**
   - Use test database for integration tests
   - Clean up test data after each test
   - Use meaningful test data

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/api/...

# Run specific test
go test -v -run TestHealthCheck ./internal/api/
```

## 📝 Commit Message Guidelines

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

### Examples

```
feat(api): add pipeline creation endpoint

Implement POST /api/v1/pipelines endpoint with validation
and database integration.

Closes #123
```

```
fix(database): resolve connection pool exhaustion

Adjust max connections and idle timeout to prevent
pool exhaustion under high load.
```

```
docs(readme): update API documentation

Add examples for all endpoints and clarify environment
variable configuration.
```

## 🔀 Pull Request Process

### Before Submitting

1. **Update your branch**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks**
   ```bash
   make all  # Runs fmt, vet, test, and build
   make lint
   ```

3. **Update documentation** if needed
   - README.md for user-facing changes
   - Code comments for implementation details
   - API documentation for endpoint changes

### Submitting a Pull Request

1. **Push your branch**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create PR** on GitHub with:
   - Clear title describing the change
   - Description explaining what and why
   - Reference to related issues
   - Screenshots for UI changes (if applicable)

3. **PR Template**
   ```markdown
   ## Description
   Brief description of changes

   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update

   ## Testing
   - [ ] Tests added/updated
   - [ ] All tests pass
   - [ ] Linter passes

   ## Related Issues
   Closes #XXX
   ```

### Review Process

1. **Automated checks** must pass:
   - Tests across Go 1.20, 1.21, 1.22
   - Linting with golangci-lint
   - Code coverage maintained
   - Security scan passes

2. **Code review** by maintainers:
   - At least one approval required
   - Address review comments
   - Keep discussion professional and constructive

3. **Merge**:
   - Squash and merge for feature branches
   - Maintainers will merge after approval

## 🔧 CI/CD Integration Points

### Understanding the CI Pipeline

The GitHub Actions workflow (`.github/workflows/ci.yml`) runs on every push and PR:

1. **Test Stage**
   - Runs tests on Go 1.20, 1.21, 1.22
   - Requires PostgreSQL service
   - Uploads coverage to Codecov

2. **Lint Stage**
   - Runs golangci-lint
   - Checks code formatting
   - Verifies go vet passes

3. **Build Stage**
   - Compiles application
   - Uploads build artifacts

4. **Docker Stage**
   - Builds Docker image
   - Pushes to registry (main branch only)

5. **Security Stage**
   - Scans for vulnerabilities
   - Reports to GitHub Security

### Local CI Simulation

Simulate CI checks locally:

```bash
# Run full CI-like checks
make fmt
make vet
make lint
make test
make build
make docker-build
```

### Adding CI/CD Steps

When adding new CI/CD steps:

1. Test locally first
2. Update `.github/workflows/ci.yml`
3. Document in CONTRIBUTING.md
4. Update README.md if user-facing

## 🐛 Reporting Bugs

### Before Reporting

1. Check existing issues
2. Ensure you're using the latest version
3. Try to reproduce with minimal setup

### Bug Report Template

```markdown
**Description**
Clear description of the bug

**To Reproduce**
Steps to reproduce:
1. ...
2. ...

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: 
- Go version:
- Application version:

**Logs**
Relevant log output
```

## 💡 Feature Requests

### Feature Request Template

```markdown
**Problem Statement**
What problem does this solve?

**Proposed Solution**
How should this work?

**Alternatives Considered**
Other approaches considered

**Additional Context**
Any other relevant information
```

## 📞 Getting Help

- **Documentation**: Check README.md first
- **Issues**: Search existing issues
- **Discussions**: Use GitHub Discussions for questions
- **Contact**: Reach out to maintainers

## 🎉 Recognition

Contributors are recognized in:
- GitHub contributors page
- Release notes
- Special thanks in major releases

Thank you for contributing! 🙏
