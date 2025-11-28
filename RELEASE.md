# TinyBrain Release Guide

## Creating a New Release

### Prerequisites
- Ensure all changes are committed and pushed to `main` branch
- All tests pass: `go test -v ./cmd/server`
- Build succeeds: `go build -o bin/server ./cmd/server`

### Release Process

#### 1. Determine Version Number
Follow [Semantic Versioning](https://semver.org/):
- **Patch** (v1.2.3 → v1.2.4): Bug fixes, minor changes
- **Minor** (v1.2.3 → v1.3.0): New features, backwards compatible
- **Major** (v1.2.3 → v2.0.0): Breaking changes

#### 2. Create and Push Tag

```bash
# Set your version
VERSION="v1.2.2"

# Create annotated tag
git tag -a $VERSION -m "Release $VERSION"

# Push tag to GitHub
git push origin $VERSION
```

#### 3. Monitor GitHub Actions
- Go to: https://github.com/rainmana/tinybrain/actions
- Watch the "Release" workflow run
- Builds typically take 2-3 minutes

#### 4. Verify Release
- Check: https://github.com/rainmana/tinybrain/releases
- Verify all binaries are present:
  - `tinybrain_*_Darwin_arm64.tar.gz` (macOS Apple Silicon)
  - `tinybrain_*_Darwin_x86_64.tar.gz` (macOS Intel)
  - `tinybrain_*_Linux_x86_64.tar.gz` (Linux x86_64)
  - `tinybrain_*_Linux_arm64.tar.gz` (Linux ARM64)
  - `tinybrain_*_Windows_x86_64.zip` (Windows)
  - `checksums.txt`

#### 5. Test `go install`

```bash
# Test installing the new version
go install github.com/rainmana/tinybrain/cmd/server@$VERSION

# Verify it works
server serve --help
```

### Troubleshooting

#### Release Workflow Fails
- Check GitHub Actions logs for errors
- Common issues:
  - Tests failing: Fix tests and push fix
  - Build errors: Check go.mod and dependencies
  - Permission issues: Verify GitHub token permissions

#### `go install` Still Shows 404
- Wait 5-10 minutes for Go proxy to update
- Try with `GOPROXY=direct go install ...`
- Clear Go module cache: `go clean -modcache`

#### Binaries Don't Work
- Verify CGO_ENABLED=0 in .goreleaser.yml
- Check architecture matches target platform
- Ensure all dependencies are included

### Updating Release Notes

Edit the release on GitHub to add:
- Detailed changelog
- Breaking changes (if any)
- Migration guides
- Screenshots or examples

### Rolling Back a Release

```bash
# Delete the tag locally
git tag -d $VERSION

# Delete the tag remotely
git push origin :refs/tags/$VERSION

# Delete the release on GitHub (via web interface)
```

## Automated Release Process

The GitHub Actions workflow automatically:
1. ✅ Runs tests
2. ✅ Builds binaries for all platforms
3. ✅ Creates archives (.tar.gz for Unix, .zip for Windows)
4. ✅ Generates checksums
5. ✅ Creates GitHub release
6. ✅ Uploads all artifacts
7. ✅ Generates changelog from commits

## Local Testing (Optional)

Install GoReleaser locally to test before pushing tags:

```bash
# Install GoReleaser (macOS)
brew install goreleaser

# Test release build without publishing
goreleaser release --snapshot --clean

# Check generated binaries
ls -lh dist/
```

## Commit Message Format

For better changelogs, use conventional commits:
- `feat: add new feature` → Features section
- `fix: resolve bug` → Bug Fixes section
- `docs: update readme` → Excluded from changelog
- `chore: update dependencies` → Excluded from changelog
