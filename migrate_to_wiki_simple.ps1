# TinyBrain Wiki Migration Script (PowerShell) - Simplified
# This script helps migrate documentation from markdown files to GitHub wiki

Write-Host "🧠 TinyBrain Wiki Migration Script" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan

# Check if we're in a git repository
try {
    $repoUrl = git remote get-url origin
    Write-Host "📁 Repository: $repoUrl" -ForegroundColor Green
} catch {
    Write-Host "❌ Error: Not in a git repository" -ForegroundColor Red
    exit 1
}

# Create wiki directory
$wikiDir = "wiki-migration"
if (Test-Path $wikiDir) {
    Remove-Item -Recurse -Force $wikiDir
}
New-Item -ItemType Directory -Path $wikiDir | Out-Null

Write-Host "📝 Creating wiki pages..." -ForegroundColor Yellow

# 1. Home Page
$homeContent = @"
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
"@

$homeContent | Out-File -FilePath "$wikiDir\Home.md" -Encoding UTF8

# 2. Getting Started
$gettingStartedContent = @"
# Getting Started

## Installation

### From Source
``````bash
go install github.com/rainmana/tinybrain/cmd/server@latest
``````

### Docker
``````bash
docker pull rainmana/tinybrain:latest
docker run -p 8080:8080 rainmana/tinybrain
``````

### Pre-built Binaries
Download from [Releases](https://github.com/rainmana/tinybrain/releases)

## Basic Usage

### 1. Start the Server
``````bash
tinybrain-server --config config.json
``````

### 2. Create a Session
``````bash
curl -X POST http://localhost:8080/sessions -H "Content-Type: application/json" -d '{"name": "Security Assessment", "task_type": "penetration_test", "intelligence_type": "osint"}'
``````

### 3. Store Intelligence
``````bash
curl -X POST http://localhost:8080/memory -H "Content-Type: application/json" -d '{"session_id": "session-id", "title": "OSINT Finding", "content": "Social media analysis reveals...", "category": "intelligence", "intelligence_type": "osint", "threat_level": "medium"}'
``````

## Configuration

See [Configuration](Configuration) for detailed setup options.

## Next Steps

- [Core Features](Core-Features)
- [Intelligence & Reconnaissance](Intelligence-&-Reconnaissance)
- [API Reference](API-Reference)
"@

$gettingStartedContent | Out-File -FilePath "$wikiDir\Getting-Started.md" -Encoding UTF8

# 3. Core Features
$coreFeaturesContent = @"
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
"@

$coreFeaturesContent | Out-File -FilePath "$wikiDir\Core-Features.md" -Encoding UTF8

Write-Host "📋 Copying documentation files..." -ForegroundColor Yellow

# Create subdirectories
$subdirs = @("intelligence", "reverse-engineering", "security-patterns", "integration", "core", "examples")
foreach ($dir in $subdirs) {
    New-Item -ItemType Directory -Path "$wikiDir\$dir" | Out-Null
}

# Copy existing documentation files
$fileMappings = @{
    "INTELLIGENCE_RECON_FRAMEWORK.md" = "intelligence\"
    "INTELLIGENCE_SECURITY_TEMPLATES.md" = "intelligence\"
    "MITRE_ATTACK_INTEGRATION.md" = "intelligence\"
    "TINYBRAIN_INTELLIGENCE_ENHANCEMENT_SUMMARY.md" = "intelligence\"
    "REVERSE_ENGINEERING_FRAMEWORK.md" = "reverse-engineering\"
    "INSIGHT_MAPPING_FRAMEWORK.md" = "reverse-engineering\"
    "CWE_SECURITY_PATTERNS.md" = "security-patterns\"
    "CWE_TINYBRAIN_INTEGRATION.md" = "security-patterns\"
    "MULTI_LANGUAGE_SECURITY_PATTERNS.md" = "security-patterns\"
    "ENHANCED_LANGUAGE_LIBRARY_PATTERNS.md" = "security-patterns\"
    "ENHANCED_AUTHORIZATION_TEMPLATES.md" = "security-patterns\"
    "SECURITY_CODE_REVIEW_DATASET.md" = "security-patterns\"
    "AI_ASSISTANT_INTEGRATION.md" = "integration\"
    "INTEGRATION_TEST_RESULTS.md" = "integration\"
    "CURSOR_SETUP.md" = "integration\"
    "ADVANCED_FEATURES.md" = "core\"
    "ENHANCED_MEMORY_CATEGORIES.md" = "core\"
    "TINYBRAIN_SECURITY_TEMPLATES.md" = "core\"
}

foreach ($file in $fileMappings.Keys) {
    if (Test-Path $file) {
        $destDir = "$wikiDir\$($fileMappings[$file])"
        Copy-Item $file $destDir
        Write-Host "✅ Copied $file to $destDir" -ForegroundColor Green
    }
}

# Copy examples
if (Test-Path "examples\basic_usage.md") {
    Copy-Item "examples\basic_usage.md" "$wikiDir\examples\"
    Write-Host "✅ Copied examples\basic_usage.md" -ForegroundColor Green
}

Write-Host "✅ Documentation files copied to $wikiDir\" -ForegroundColor Green

# Create migration instructions
$migrationInstructions = @"
# Wiki Migration Instructions

## Step 1: Enable GitHub Wiki
1. Go to your GitHub repository
2. Click on "Wiki" tab
3. Click "Create the first page"
4. Copy content from Home.md as the home page

## Step 2: Create Wiki Pages
Create the following pages in your GitHub wiki:

### Core Pages
- **Getting Started**: Copy from Getting-Started.md
- **Core Features**: Copy from Core-Features.md
- **Configuration**: Copy from config.example.json (format as markdown)

### Intelligence & Reconnaissance
- **Intelligence Framework**: Copy from intelligence\INTELLIGENCE_RECON_FRAMEWORK.md
- **Intelligence Templates**: Copy from intelligence\INTELLIGENCE_SECURITY_TEMPLATES.md
- **MITRE ATT&CK Integration**: Copy from intelligence\MITRE_ATTACK_INTEGRATION.md
- **Intelligence Enhancements**: Copy from intelligence\TINYBRAIN_INTELLIGENCE_ENHANCEMENT_SUMMARY.md

### Reverse Engineering
- **Reverse Engineering**: Copy from reverse-engineering\REVERSE_ENGINEERING_FRAMEWORK.md
- **Insight Mapping**: Copy from reverse-engineering\INSIGHT_MAPPING_FRAMEWORK.md

### Security Patterns
- **CWE Security Patterns**: Copy from security-patterns\CWE_SECURITY_PATTERNS.md
- **CWE Integration**: Copy from security-patterns\CWE_TINYBRAIN_INTEGRATION.md
- **Multi-Language Patterns**: Copy from security-patterns\MULTI_LANGUAGE_SECURITY_PATTERNS.md
- **Language Library Patterns**: Copy from security-patterns\ENHANCED_LANGUAGE_LIBRARY_PATTERNS.md
- **Authorization Templates**: Copy from security-patterns\ENHANCED_AUTHORIZATION_TEMPLATES.md
- **Security Datasets**: Copy from security-patterns\SECURITY_CODE_REVIEW_DATASET.md

### Integration & Development
- **AI Assistant Integration**: Copy from integration\AI_ASSISTANT_INTEGRATION.md
- **Test Results**: Copy from integration\INTEGRATION_TEST_RESULTS.md
- **Development Setup**: Copy from integration\CURSOR_SETUP.md

### Core Features (Additional)
- **Advanced Features**: Copy from core\ADVANCED_FEATURES.md
- **Memory Categories**: Copy from core\ENHANCED_MEMORY_CATEGORIES.md
- **Security Templates**: Copy from core\TINYBRAIN_SECURITY_TEMPLATES.md

### Examples
- **Basic Usage**: Copy from examples\basic_usage.md

## Step 3: Update Main README
1. Replace current README.md with content from README_SIMPLIFIED.md
2. Update repository URLs to match your actual repository
3. Test all links

## Step 4: Clean Up Repository
After migrating to wiki, you can remove these files from the main repository:
- All *.md files except README.md
- *.json files that are documentation
- examples\ directory (if moved to wiki)

## Step 5: Test
1. Verify all wiki pages load correctly
2. Test all internal links
3. Verify main README links to wiki
4. Test navigation between pages

## Notes
- GitHub wiki uses different markdown syntax for internal links
- Use [Page Name](Page-Name) format for internal links
- External links work the same as regular markdown
- Images need to be uploaded to the wiki or linked from the main repository
"@

$migrationInstructions | Out-File -FilePath "$wikiDir\MIGRATION_INSTRUCTIONS.md" -Encoding UTF8

Write-Host "📋 Migration instructions created in $wikiDir\MIGRATION_INSTRUCTIONS.md" -ForegroundColor Green

Write-Host ""
Write-Host "🎉 Wiki migration preparation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Review the files in $wikiDir\" -ForegroundColor White
Write-Host "2. Follow the instructions in $wikiDir\MIGRATION_INSTRUCTIONS.md" -ForegroundColor White
Write-Host "3. Enable GitHub wiki for your repository" -ForegroundColor White
Write-Host "4. Create wiki pages with the provided content" -ForegroundColor White
Write-Host "5. Update your main README" -ForegroundColor White
Write-Host "6. Clean up the repository (optional)" -ForegroundColor White
Write-Host ""
Write-Host "📁 All files are ready in: $wikiDir\" -ForegroundColor Cyan
