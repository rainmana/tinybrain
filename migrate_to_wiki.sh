#!/bin/bash

# TinyBrain Wiki Migration Script
# This script helps migrate documentation from markdown files to GitHub wiki

set -e

echo "🧠 TinyBrain Wiki Migration Script"
echo "=================================="

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# Get repository URL
REPO_URL=$(git remote get-url origin)
echo "📁 Repository: $REPO_URL"

# Extract repository name and owner
if [[ $REPO_URL =~ github\.com[:/]([^/]+)/([^/]+)\.git ]]; then
    OWNER="${BASH_REMATCH[1]}"
    REPO="${BASH_REMATCH[2]}"
    echo "👤 Owner: $OWNER"
    echo "📦 Repository: $REPO"
else
    echo "❌ Error: Could not extract repository information"
    exit 1
fi

# Create wiki directory
WIKI_DIR="wiki-migration"
mkdir -p "$WIKI_DIR"

echo "📝 Creating wiki pages..."

# 1. Home Page
cat > "$WIKI_DIR/Home.md" << 'EOF'
# TinyBrain - Security-Focused LLM Memory Storage

TinyBrain is a Model Context Protocol (MCP) server designed for security professionals, penetration testers, and AI assistants working on offensive security tasks.

## Quick Start
[Getting Started](Getting-Started)

## Key Features
- **Intelligence Gathering**: OSINT, HUMINT, SIGINT, and more
- **Reverse Engineering**: Malware analysis, binary analysis, vulnerability research
- **MITRE ATT&CK Integration**: Complete framework support
- **Security Patterns**: CWE, OWASP, and multi-language patterns
- **Memory Management**: 30+ memory categories for security data

## Documentation
- [Getting Started](Getting-Started)
- [Core Features](Core-Features)
- [Intelligence & Reconnaissance](Intelligence-&-Reconnaissance)
- [Reverse Engineering](Reverse-Engineering)
- [Security Patterns](Security-Patterns)
- [Integration & Development](Integration-&-Development)
- [API Reference](API-Reference)
- [Contributing](Contributing)
EOF

# 2. Getting Started
cat > "$WIKI_DIR/Getting-Started.md" << 'EOF'
# Getting Started

## Installation

### From Source
```bash
go install github.com/rainmana/tinybrain/cmd/server@latest
```

### Docker
```bash
docker pull rainmana/tinybrain:latest
docker run -p 8080:8080 rainmana/tinybrain
```

### Pre-built Binaries
Download from [Releases](https://github.com/rainmana/tinybrain/releases)

## Basic Usage

### 1. Start the Server
```bash
tinybrain-server --config config.json
```

### 2. Create a Session
```bash
curl -X POST http://localhost:8080/sessions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Security Assessment",
    "task_type": "penetration_test",
    "intelligence_type": "osint"
  }'
```

### 3. Store Intelligence
```bash
curl -X POST http://localhost:8080/memory \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-id",
    "title": "OSINT Finding",
    "content": "Social media analysis reveals...",
    "category": "intelligence",
    "intelligence_type": "osint",
    "threat_level": "medium"
  }'
```

## Configuration

See [Configuration](Configuration) for detailed setup options.

## Next Steps

- [Core Features](Core-Features)
- [Intelligence & Reconnaissance](Intelligence-&-Reconnaissance)
- [API Reference](API-Reference)
EOF

# 3. Core Features
cat > "$WIKI_DIR/Core-Features.md" << 'EOF'
# Core Features

## Memory Management
- 30+ memory categories
- Intelligent search and retrieval
- Relationship mapping
- Context snapshots

## Session Management
- Multi-session support
- Task progress tracking
- Status management
- Metadata storage

## Search Capabilities
- Full-text search
- Semantic search
- Category filtering
- Priority-based search

## Advanced Features
- Pattern recognition
- Correlation analysis
- Insight mapping
- Knowledge graphs

## Security Templates
- Pre-configured templates
- OWASP integration
- CWE patterns
- Multi-language support
EOF

# Copy existing documentation files to wiki directory
echo "📋 Copying documentation files..."

# Intelligence & Reconnaissance
mkdir -p "$WIKI_DIR/intelligence"
cp "INTELLIGENCE_RECON_FRAMEWORK.md" "$WIKI_DIR/intelligence/"
cp "INTELLIGENCE_SECURITY_TEMPLATES.md" "$WIKI_DIR/intelligence/"
cp "MITRE_ATTACK_INTEGRATION.md" "$WIKI_DIR/intelligence/"
cp "TINYBRAIN_INTELLIGENCE_ENHANCEMENT_SUMMARY.md" "$WIKI_DIR/intelligence/"

# Reverse Engineering
mkdir -p "$WIKI_DIR/reverse-engineering"
cp "REVERSE_ENGINEERING_FRAMEWORK.md" "$WIKI_DIR/reverse-engineering/"
cp "INSIGHT_MAPPING_FRAMEWORK.md" "$WIKI_DIR/reverse-engineering/"

# Security Patterns
mkdir -p "$WIKI_DIR/security-patterns"
cp "CWE_SECURITY_PATTERNS.md" "$WIKI_DIR/security-patterns/"
cp "CWE_TINYBRAIN_INTEGRATION.md" "$WIKI_DIR/security-patterns/"
cp "MULTI_LANGUAGE_SECURITY_PATTERNS.md" "$WIKI_DIR/security-patterns/"
cp "ENHANCED_LANGUAGE_LIBRARY_PATTERNS.md" "$WIKI_DIR/security-patterns/"
cp "ENHANCED_AUTHORIZATION_TEMPLATES.md" "$WIKI_DIR/security-patterns/"
cp "SECURITY_CODE_REVIEW_DATASET.md" "$WIKI_DIR/security-patterns/"

# Integration & Development
mkdir -p "$WIKI_DIR/integration"
cp "AI_ASSISTANT_INTEGRATION.md" "$WIKI_DIR/integration/"
cp "INTEGRATION_TEST_RESULTS.md" "$WIKI_DIR/integration/"
cp "CURSOR_SETUP.md" "$WIKI_DIR/integration/"

# Core Features
mkdir -p "$WIKI_DIR/core"
cp "ADVANCED_FEATURES.md" "$WIKI_DIR/core/"
cp "ENHANCED_MEMORY_CATEGORIES.md" "$WIKI_DIR/core/"
cp "TINYBRAIN_SECURITY_TEMPLATES.md" "$WIKI_DIR/core/"

# Examples
mkdir -p "$WIKI_DIR/examples"
cp "examples/basic_usage.md" "$WIKI_DIR/examples/"

echo "✅ Documentation files copied to $WIKI_DIR/"

# Create migration instructions
cat > "$WIKI_DIR/MIGRATION_INSTRUCTIONS.md" << EOF
# Wiki Migration Instructions

## Step 1: Enable GitHub Wiki
1. Go to https://github.com/$OWNER/$REPO
2. Click on "Wiki" tab
3. Click "Create the first page"
4. Copy content from \`Home.md\` as the home page

## Step 2: Create Wiki Pages
Create the following pages in your GitHub wiki:

### Core Pages
- **Getting Started**: Copy from \`Getting-Started.md\`
- **Core Features**: Copy from \`Core-Features.md\`
- **Configuration**: Copy from \`config.example.json\` (format as markdown)

### Intelligence & Reconnaissance
- **Intelligence Framework**: Copy from \`intelligence/INTELLIGENCE_RECON_FRAMEWORK.md\`
- **Intelligence Templates**: Copy from \`intelligence/INTELLIGENCE_SECURITY_TEMPLATES.md\`
- **MITRE ATT&CK Integration**: Copy from \`intelligence/MITRE_ATTACK_INTEGRATION.md\`
- **Intelligence Enhancements**: Copy from \`intelligence/TINYBRAIN_INTELLIGENCE_ENHANCEMENT_SUMMARY.md\`

### Reverse Engineering
- **Reverse Engineering**: Copy from \`reverse-engineering/REVERSE_ENGINEERING_FRAMEWORK.md\`
- **Insight Mapping**: Copy from \`reverse-engineering/INSIGHT_MAPPING_FRAMEWORK.md\`

### Security Patterns
- **CWE Security Patterns**: Copy from \`security-patterns/CWE_SECURITY_PATTERNS.md\`
- **CWE Integration**: Copy from \`security-patterns/CWE_TINYBRAIN_INTEGRATION.md\`
- **Multi-Language Patterns**: Copy from \`security-patterns/MULTI_LANGUAGE_SECURITY_PATTERNS.md\`
- **Language Library Patterns**: Copy from \`security-patterns/ENHANCED_LANGUAGE_LIBRARY_PATTERNS.md\`
- **Authorization Templates**: Copy from \`security-patterns/ENHANCED_AUTHORIZATION_TEMPLATES.md\`
- **Security Datasets**: Copy from \`security-patterns/SECURITY_CODE_REVIEW_DATASET.md\`

### Integration & Development
- **AI Assistant Integration**: Copy from \`integration/AI_ASSISTANT_INTEGRATION.md\`
- **Test Results**: Copy from \`integration/INTEGRATION_TEST_RESULTS.md\`
- **Development Setup**: Copy from \`integration/CURSOR_SETUP.md\`

### Core Features (Additional)
- **Advanced Features**: Copy from \`core/ADVANCED_FEATURES.md\`
- **Memory Categories**: Copy from \`core/ENHANCED_MEMORY_CATEGORIES.md\`
- **Security Templates**: Copy from \`core/TINYBRAIN_SECURITY_TEMPLATES.md\`

### Examples
- **Basic Usage**: Copy from \`examples/basic_usage.md\`

## Step 3: Update Main README
1. Replace current \`README.md\` with content from \`README_SIMPLIFIED.md\`
2. Update repository URLs to match your actual repository
3. Test all links

## Step 4: Clean Up Repository
After migrating to wiki, you can remove these files from the main repository:
- All \`*.md\` files except \`README.md\`
- \`*.json\` files that are documentation
- \`examples/\` directory (if moved to wiki)

## Step 5: Test
1. Verify all wiki pages load correctly
2. Test all internal links
3. Verify main README links to wiki
4. Test navigation between pages

## Notes
- GitHub wiki uses different markdown syntax for internal links
- Use \`[Page Name](Page-Name)\` format for internal links
- External links work the same as regular markdown
- Images need to be uploaded to the wiki or linked from the main repository
EOF

echo "📋 Migration instructions created in $WIKI_DIR/MIGRATION_INSTRUCTIONS.md"

# Create a script to update the main README
cat > "$WIKI_DIR/update_readme.sh" << 'EOF'
#!/bin/bash

# Update main README with simplified version
echo "Updating main README..."

# Backup current README
cp README.md README_backup.md

# Replace with simplified version
cp README_SIMPLIFIED.md README.md

echo "✅ README updated"
echo "📋 Backup saved as README_backup.md"
EOF

chmod +x "$WIKI_DIR/update_readme.sh"

echo "🔧 README update script created: $WIKI_DIR/update_readme.sh"

# Create a cleanup script
cat > "$WIKI_DIR/cleanup_repo.sh" << 'EOF'
#!/bin/bash

# Clean up repository after wiki migration
echo "🧹 Cleaning up repository..."

# List of files to remove (after confirming wiki migration)
FILES_TO_REMOVE=(
    "ADVANCED_FEATURES.md"
    "AI_ASSISTANT_INTEGRATION.md"
    "CURSOR_SETUP.md"
    "CWE_SECURITY_PATTERNS.md"
    "CWE_TINYBRAIN_INTEGRATION.md"
    "ENHANCED_AUTHORIZATION_TEMPLATES.md"
    "ENHANCED_LANGUAGE_LIBRARY_PATTERNS.md"
    "ENHANCED_MEMORY_CATEGORIES.md"
    "INSIGHT_MAPPING_FRAMEWORK.md"
    "INTELLIGENCE_RECON_FRAMEWORK.md"
    "INTELLIGENCE_SECURITY_TEMPLATES.md"
    "INTEGRATION_TEST_RESULTS.md"
    "MITRE_ATTACK_INTEGRATION.md"
    "MULTI_LANGUAGE_SECURITY_PATTERNS.md"
    "REVERSE_ENGINEERING_FRAMEWORK.md"
    "SECURITY_CODE_REVIEW_DATASET.md"
    "TINYBRAIN_INTELLIGENCE_ENHANCEMENT_SUMMARY.md"
    "TINYBRAIN_SECURITY_TEMPLATES.md"
    "GITHUB_WIKI_MIGRATION_PLAN.md"
    "README_SIMPLIFIED.md"
    "test_intelligence_enhancements.go"
    "test_enhancements.ps1"
    "test_enhancements_simple.ps1"
    "TEST_RESULTS_INTELLIGENCE_ENHANCEMENTS.md"
)

echo "Files that will be removed:"
for file in "${FILES_TO_REMOVE[@]}"; do
    if [ -f "$file" ]; then
        echo "  - $file"
    fi
done

read -p "Do you want to remove these files? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    for file in "${FILES_TO_REMOVE[@]}"; do
        if [ -f "$file" ]; then
            rm "$file"
            echo "✅ Removed $file"
        fi
    done
    echo "🧹 Cleanup complete"
else
    echo "❌ Cleanup cancelled"
fi
EOF

chmod +x "$WIKI_DIR/cleanup_repo.sh"

echo "🧹 Cleanup script created: $WIKI_DIR/cleanup_repo.sh"

echo ""
echo "🎉 Wiki migration preparation complete!"
echo ""
echo "Next steps:"
echo "1. Review the files in $WIKI_DIR/"
echo "2. Follow the instructions in $WIKI_DIR/MIGRATION_INSTRUCTIONS.md"
echo "3. Enable GitHub wiki for your repository"
echo "4. Create wiki pages with the provided content"
echo "5. Update your main README"
echo "6. Clean up the repository (optional)"
echo ""
echo "📁 All files are ready in: $WIKI_DIR/"
