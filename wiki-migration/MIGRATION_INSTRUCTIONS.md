# GitHub Wiki Migration Instructions

## Step 1: Enable GitHub Wiki
1. Go to your GitHub repository: https://github.com/rainmana/tinybrain
2. Click on "Wiki" tab
3. Click "Create the first page"
4. Copy content from `Home.md` as the home page

## Step 2: Create Wiki Pages
Create the following pages in your GitHub wiki:

### Core Pages
- **Getting Started**: Copy from `Getting-Started.md`
- **Core Features**: Copy from `core/ADVANCED_FEATURES.md`
- **Memory Categories**: Copy from `core/ENHANCED_MEMORY_CATEGORIES.md`
- **Security Templates**: Copy from `core/TINYBRAIN_SECURITY_TEMPLATES.md`

### Intelligence & Reconnaissance
- **Intelligence Framework**: Copy from `intelligence/INTELLIGENCE_RECON_FRAMEWORK.md`
- **Intelligence Templates**: Copy from `intelligence/INTELLIGENCE_SECURITY_TEMPLATES.md`
- **MITRE ATT&CK Integration**: Copy from `intelligence/MITRE_ATTACK_INTEGRATION.md`
- **Intelligence Enhancements**: Copy from `intelligence/TINYBRAIN_INTELLIGENCE_ENHANCEMENT_SUMMARY.md`

### Reverse Engineering
- **Reverse Engineering**: Copy from `reverse-engineering/REVERSE_ENGINEERING_FRAMEWORK.md`
- **Insight Mapping**: Copy from `reverse-engineering/INSIGHT_MAPPING_FRAMEWORK.md`

### Security Patterns
- **CWE Security Patterns**: Copy from `security-patterns/CWE_SECURITY_PATTERNS.md`
- **Multi-Language Patterns**: Copy from `security-patterns/MULTI_LANGUAGE_SECURITY_PATTERNS.md`

### Integration & Development
- **AI Assistant Integration**: Copy from `integration/AI_ASSISTANT_INTEGRATION.md`
- **Development Setup**: Copy from `integration/CURSOR_SETUP.md`

### Examples
- **Basic Usage**: Copy from `examples/basic_usage.md`

## Step 3: Update Main README
1. Replace current `README.md` with content from `README_SIMPLIFIED.md`
2. Update repository URLs to match your actual repository
3. Test all links

## Step 4: Clean Up Repository
After migrating to wiki, you can remove these files from the main repository:
- All `*.md` files except `README.md`
- `*.json` files that are documentation
- `examples/` directory (if moved to wiki)

## Step 5: Test
1. Verify all wiki pages load correctly
2. Test all internal links
3. Verify main README links to wiki
4. Test navigation between pages

## Notes
- GitHub wiki uses different markdown syntax for internal links
- Use `[Page Name](Page-Name)` format for internal links
- External links work the same as regular markdown
- Images need to be uploaded to the wiki or linked from the main repository

## Benefits of Wiki Migration
1. **Better Organization**: Structured navigation and categorization
2. **Improved Discoverability**: Easy to find specific information
3. **Cleaner Repository**: Main repo focuses on code, not documentation
4. **Collaborative Editing**: Multiple contributors can edit wiki pages
5. **Version Control**: GitHub tracks wiki changes
6. **Search**: Built-in wiki search functionality
7. **Mobile Friendly**: Better mobile reading experience
