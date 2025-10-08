---
layout: default
title: Contributing
nav_order: 17
description: "Contributing guide for TinyBrain"
---

# Contributing

We welcome contributions to TinyBrain! This guide will help you get started with contributing code, documentation, or other improvements.

## Table of Contents
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contribution Process](#contribution-process)
- [Code Standards](#code-standards)
- [Testing Requirements](#testing-requirements)
- [Pull Request Guidelines](#pull-request-guidelines)

## Getting Started

Before contributing, please:

1. **Review the Documentation**: Familiarize yourself with TinyBrain's architecture and features
2. **Check Existing Issues**: Look for existing issues or feature requests
3. **Discuss Major Changes**: Open an issue to discuss significant changes before starting work
4. **Read the Code of Conduct**: Follow our community guidelines

## Development Setup

### Prerequisites

- Go 1.19 or later
- Git
- SQLite3
- Make (optional, for convenience)

### Setup Steps

**1. Fork the Repository**

Fork the TinyBrain repository to your GitHub account.

**2. Clone Your Fork**

```bash
git clone https://github.com/YOUR_USERNAME/tinybrain.git
cd tinybrain
```

**3. Add Upstream Remote**

```bash
git remote add upstream https://github.com/rainmana/tinybrain.git
```

**4. Setup Development Environment**

```bash
# Download dependencies
make dev-setup

# Or manually
go mod download
go mod tidy
```

**5. Verify Setup**

```bash
# Run tests
make test

# Build the project
make build
```

## Contribution Process

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-export-import`
- `fix/database-connection-leak`
- `docs/update-api-examples`

### 2. Make Your Changes

- Write clean, maintainable code
- Follow existing code style and conventions
- Add or update tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run linter
make lint

# Run security checks
make security
```

### 4. Commit Your Changes

Use clear, descriptive commit messages:

```bash
git add .
git commit -m "Add export/import functionality for sessions"
```

**Commit Message Guidelines**:
- Use present tense ("Add feature" not "Added feature")
- Be specific about what changed
- Reference issue numbers when applicable

### 5. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 6. Submit a Pull Request

1. Go to your fork on GitHub
2. Click "New Pull Request"
3. Select your feature branch
4. Fill out the PR template with:
   - Clear description of changes
   - Related issue numbers
   - Testing performed
   - Breaking changes (if any)

## Code Standards

### Go Code Style

- Follow standard Go formatting (use `go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and concise
- Handle errors appropriately

**Example**:
```go
// CreateSession creates a new security assessment session
func (r *Repository) CreateSession(ctx context.Context, session *models.Session) error {
    if session == nil {
        return fmt.Errorf("session cannot be nil")
    }
    
    // Validate session fields
    if err := session.Validate(); err != nil {
        return fmt.Errorf("invalid session: %w", err)
    }
    
    // Insert into database
    _, err := r.db.ExecContext(ctx, query, session.ID, session.Name, ...)
    if err != nil {
        return fmt.Errorf("failed to create session: %w", err)
    }
    
    return nil
}
```

### Documentation Standards

- Update README.md for new features
- Add examples for new MCP tools
- Update API documentation
- Include inline code comments
- Write clear commit messages

### Testing Standards

- Write tests for all new functionality
- Maintain or improve test coverage
- Test both success and error paths
- Use table-driven tests where appropriate
- Clean up test resources (databases, files)

**Example Test**:
```go
func TestCreateSession(t *testing.T) {
    tests := []struct {
        name    string
        session *models.Session
        wantErr bool
    }{
        {
            name: "valid session",
            session: &models.Session{
                ID:   "test-id",
                Name: "Test Session",
            },
            wantErr: false,
        },
        {
            name:    "nil session",
            session: nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := repo.CreateSession(ctx, tt.session)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Testing Requirements

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v ./internal/database -run TestNewDatabase

# Run benchmarks
make bench
```

### Test Coverage

- Aim for 90%+ test coverage (current target)
- All new features must include tests
- Tests should cover edge cases and error conditions
- Integration tests for database operations
- Unit tests for business logic

### Verify Test Coverage

```bash
# Check current coverage
go test -cover ./...

# Generate coverage report
make test
# Open coverage.html in browser
```

## Pull Request Guidelines

### Before Submitting

- [ ] All tests pass
- [ ] Code follows Go formatting standards
- [ ] New features have tests
- [ ] Documentation updated
- [ ] No linting errors
- [ ] Commit messages are clear

### PR Template

**Title**: Clear, descriptive title

**Description**:
- What does this PR do?
- Why is this change needed?
- How was it tested?
- Are there any breaking changes?

**Related Issues**: #123, #456

**Checklist**:
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes
- [ ] All tests passing

### Review Process

1. Maintainer review
2. Address feedback
3. Tests must pass
4. Approval required
5. Squash and merge

## Development Targets

Available Makefile targets for development:

**Building**:
- `make build` - Build the binary
- `make install` - Install to GOPATH/bin

**Testing**:
- `make test` - Run tests with coverage
- `make test-verbose` - Verbose test output
- `make bench` - Run benchmarks

**Quality**:
- `make lint` - Run go vet
- `make fmt` - Format code
- `make security` - Run security checks

**Development**:
- `make dev-setup` - Setup dev environment
- `make run` - Run the server
- `make run-dev` - Run with dev database

**Database**:
- `make db-init` - Initialize database directory
- `make db-reset` - Reset database

**Cleanup**:
- `make clean` - Clean build artifacts

## Getting Help

- **Issues**: Open a GitHub issue for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check the docs/ directory
- **Examples**: See examples/ for usage patterns

## License

By contributing to TinyBrain, you agree that your contributions will be licensed under the MIT License.

## Next Steps

- Review [Architecture](architecture.md) to understand the system design
- Check [API Examples](api/examples.md) for implementation patterns
- See [Roadmap](roadmap.md) for planned features
