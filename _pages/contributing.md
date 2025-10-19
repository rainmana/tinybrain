---
layout: default
title: Contributing
permalink: /contributing/
---

# Contributing

We welcome contributions to TinyBrain! This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git
- Docker (optional)
- Make

### Development Setup
1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/tinybrain.git
   cd tinybrain
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the project:
   ```bash
   make build
   ```
5. Run tests:
   ```bash
   make test
   ```

## Contribution Guidelines

### Code Style
- Follow Go coding standards
- Use `gofmt` to format code
- Use `golint` for linting
- Write comprehensive tests
- Document public APIs

### Commit Messages
Use conventional commit format:
```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build process or auxiliary tool changes

Examples:
```
feat(intelligence): add OSINT data collection capabilities
fix(memory): resolve memory leak in search functionality
docs(api): update API documentation for new endpoints
```

### Pull Request Process
1. Create a feature branch from `main`
2. Make your changes
3. Add tests for new functionality
4. Update documentation if needed
5. Run tests and ensure they pass
6. Submit a pull request

### Pull Request Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] New tests added for new functionality
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or breaking changes documented)
```

## Development Areas

### Core Features
- Memory management
- Session handling
- Search functionality
- Relationship mapping
- Context snapshots

### Intelligence Features
- OSINT capabilities
- HUMINT integration
- SIGINT analysis
- MITRE ATT&CK mapping
- Threat intelligence

### Reverse Engineering
- Malware analysis
- Binary analysis
- Vulnerability research
- Protocol analysis
- Tool integration

### Security Patterns
- CWE integration
- OWASP patterns
- Multi-language support
- Authorization templates
- Security datasets

### Integration
- MCP protocol
- AI assistant integration
- API development
- Cloud integration
- CI/CD support

## Testing

### Unit Tests
```bash
go test ./internal/...
```

### Integration Tests
```bash
go test -tags=integration ./...
```

### Test Coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Performance Tests
```bash
go test -bench=. ./...
```

## Documentation

### Code Documentation
- Document all public functions
- Use Go doc comments
- Include examples where appropriate
- Keep documentation up to date

### API Documentation
- Update OpenAPI specifications
- Document new endpoints
- Include request/response examples
- Update error codes and messages

### User Documentation
- Update README files
- Add usage examples
- Document configuration options
- Update troubleshooting guides

## Security

### Security Guidelines
- Follow secure coding practices
- Validate all inputs
- Use parameterized queries
- Implement proper authentication
- Follow principle of least privilege

### Vulnerability Reporting
Report security vulnerabilities privately:
1. Email: security@tinybrain.dev
2. Use PGP encryption if possible
3. Include detailed reproduction steps
4. Allow reasonable time for response

### Security Review Process
- All code changes require security review
- Security team reviews all pull requests
- Automated security scanning in CI/CD
- Regular security audits

## Performance

### Performance Guidelines
- Optimize for memory usage
- Minimize database queries
- Use efficient algorithms
- Profile performance-critical code
- Monitor resource usage

### Performance Testing
```bash
go test -bench=. -benchmem ./...
```

### Performance Monitoring
- Monitor memory usage
- Track database performance
- Measure response times
- Profile CPU usage

## Release Process

### Version Numbering
We use semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

### Release Checklist
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Version numbers updated
- [ ] CHANGELOG updated
- [ ] Release notes prepared
- [ ] Security review completed

### Release Steps
1. Update version numbers
2. Update CHANGELOG
3. Create release branch
4. Run full test suite
5. Create release tag
6. Build and test release artifacts
7. Publish release

## Community

### Communication
- GitHub Issues: Bug reports and feature requests
- GitHub Discussions: General discussion and questions
- Discord: Real-time chat and collaboration
- Email: security@tinybrain.dev for security issues

### Code of Conduct
We follow the Contributor Covenant Code of Conduct:
- Be respectful and inclusive
- Focus on constructive feedback
- Respect different viewpoints
- Accept responsibility for mistakes
- Help create a welcoming environment

### Recognition
Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation
- Community highlights

## Getting Help

### Documentation
- [Getting Started](getting-started/)
- [API Reference](api-reference/)
- [Core Features](core-features/)
- [Integration Guide](integration/)

### Support Channels
- GitHub Issues: Technical issues and bugs
- GitHub Discussions: Questions and ideas
- Discord: Real-time help and discussion
- Email: support@tinybrain.dev

### Mentorship
- New contributor mentorship program
- Code review guidance
- Architecture discussions
- Best practices sharing

## License

By contributing to TinyBrain, you agree that your contributions will be licensed under the MIT License.

## Thank You

Thank you for considering contributing to TinyBrain! Your contributions help make security tools more accessible and effective for the community.
