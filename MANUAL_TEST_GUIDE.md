# Manual Test Guide - Real Intelligence Feeds Verification

## Quick Start (5 minutes)

### 1. Build the Server
```bash
cd /Users/alec/tinybrain
./build_intelligence_final.sh
```

**Expected Output:**
```
✅ Build successful: tinybrain-intelligence-final
```

### 2. Start the Server
```bash
./tinybrain-intelligence-final
```

**Expected Output:**
```
TinyBrain Intelligence Final Server starting on http://127.0.0.1:8090
Admin dashboard: http://127.0.0.1:8090/_/
REST API: http://127.0.0.1:8090/api/
```

### 3. Test Basic Connectivity (In New Terminal)
```bash
# Test REST API
curl http://127.0.0.1:8090/api/ | jq

# Test MCP endpoint
curl -X POST http://127.0.0.1:8090/mcp \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | jq
```

**Expected**: JSON responses with server info

### 4. Download Real Data (One-Time Setup)
```bash
# This takes 5-10 minutes on first run
curl -X POST http://127.0.0.1:8090/mcp \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"download_security_data","arguments":{}}}' | jq
```

**Expected Output:**
```json
{
  "result": {
    "status": "success" or "partial_success",
    "nvd": { "success": true, "error": "" },
    "attack": { "success": true, "error": "" },
    "owasp": { "success": true, "error": "" }
  }
}
```

**Note:** If you see errors, that's expected for first run. Wait a few minutes for downloads to complete.

### 5. Verify Real Data (After Downloads Complete)
```bash
# Query NVD for real CVE data
curl -X POST http://127.0.0.1:8090/mcp \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_nvd","arguments":{"query":"buffer overflow","limit":3}}}' | jq '.result.results[0]'
```

**Expected:** Real CVE data like:
```json
{
  "cve_id": "CVE-2024-38077",
  "description": "Windows Remote Desktop Services ...",
  "severity": "CRITICAL",
  "cvss_v3_score": 9.8,
  ...
}
```

**NOT Expected (Mock Data):**
```json
{
  "cve_id": "CVE-2024-1234",
  "description": "Sample CVE for testing intelligence feeds - ...",
  ...
}
```

### 6. Query ATT&CK for Real Technique Data
```bash
curl -X POST http://127.0.0.1:8090/mcp \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_attack","arguments":{"query":"Process Injection","limit":3}}}' | jq '.result.results[0]'
```

**Expected:** Real ATT&CK technique like:
```json
{
  "technique_id": "T1055",
  "name": "Process Injection",
  "description": "Adversaries may inject code into processes...",
  "tactic": "Defense Evasion",
  "platforms": ["Windows", "macOS", "Linux"]
}
```

### 7. Check Database Summary
```bash
curl -X POST http://127.0.0.1:8090/mcp \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_security_data_summary","arguments":{}}}' | jq '.result.summary.stored_records'
```

**Expected (After First Download):**
```json
{
  "nvd": 300000,     // Should be > 100,000
  "attack": 600,     // Should be > 500
  "owasp": 50        // Should be > 10
}
```

**NOT Expected (Mock Data):**
```json
{
  "nvd": 2,
  "attack": 2,
  "owasp": 2
}
```

## Full Validation Test

For comprehensive validation that checks everything:

```bash
cd /Users/alec/tinybrain
./test_real_intelligence_validation.sh
```

This will:
1. Build the server
2. Start the server
3. Download all data sources
4. Run 10 validation tests
5. Verify no mock/sample data
6. Check database has substantial data
7. Validate CVE and technique ID formats

**Expected:** All tests pass (10/10)

## Verification Checklist

### ✅ Real Data Indicators
- [ ] CVE IDs match format: `CVE-YYYY-NNNNN` (e.g., CVE-2024-38077)
- [ ] Technique IDs match format: `TNNNN` (e.g., T1055)
- [ ] Descriptions are real (not "Sample CVE for testing")
- [ ] Database has > 100,000 CVEs
- [ ] Database has > 500 techniques
- [ ] CVSS scores are realistic (0.0-10.0)
- [ ] No "mock", "sample", or "for testing" in descriptions

### ❌ Mock Data Indicators (Should NOT See)
- [ ] CVE IDs like "CVE-2024-1234" or "CVE-2024-5678"
- [ ] Descriptions containing "Sample CVE for testing intelligence feeds"
- [ ] Descriptions containing "Related to: [your query]"
- [ ] Database with only 2 entries per source
- [ ] Technique descriptions with "mock" or "sample"

## Troubleshooting

### Problem: Server Won't Start
**Solution:** Check if port 8090 is already in use:
```bash
lsof -i :8090
# Kill any existing process
kill <PID>
```

### Problem: Downloads Fail
**Solution:** Check internet connection and NVD API status:
```bash
curl -I https://services.nvd.nist.gov/rest/json/cves/2.0
```

### Problem: Database Empty After Download
**Solution:** Check database file exists and has size:
```bash
ls -lh ~/.tinybrain-intelligence-final/data.db
```

### Problem: Still Seeing Mock Data
**Solution:** Delete database and re-download:
```bash
rm -f ~/.tinybrain-intelligence-final/data.db
# Then restart server and download again
```

## Performance Expectations

### Initial Setup
- Build time: 10-30 seconds
- First data download: 5-10 minutes
- Database size after download: ~100-200 MB

### Normal Operation
- Server startup: < 1 second
- Query response time: < 100ms
- No further downloads needed (uses local database)

## For Your FAANG Colleagues

### To Demonstrate Real Implementation

1. **Show Build Success:**
```bash
./build_intelligence_final.sh
```

2. **Show Real CVE Query:**
```bash
curl -X POST http://127.0.0.1:8090/mcp \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_nvd","arguments":{"query":"CVE-2024-38077","limit":1}}}' | jq
```

3. **Show Database Size:**
```bash
ls -lh ~/.tinybrain-intelligence-final/data.db
du -h ~/.tinybrain-intelligence-final/
```

4. **Show No Mock Data:**
```bash
# Search for "sample" - should find 0 results or only legitimate samples
curl -X POST http://127.0.0.1:8090/mcp \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"query_nvd","arguments":{"query":"sample CVE for testing","limit":10}}}' | jq '.result.total_count'
```

Should return `0` or real CVEs that happen to contain "sample" in legitimate context.

### Code Review Points

1. **Check Query Handler** (`cmd/server/pocketbase_intelligence_final.go:421-490`):
   - Uses `s.securityRepo.QueryNVD(ctx, searchReq)`
   - No hardcoded results
   - Converts database models to response format

2. **Check Download Handler** (`cmd/server/pocketbase_intelligence_final.go:321-353`):
   - Uses `s.securityDownloader.DownloadNVDDataset(ctx)`
   - Calls `s.securityRepo.StoreNVDDataset(ctx, cves)`
   - No fallback to mock data

3. **Check Database Operations** (`internal/repository/security_repository.go:41-90`):
   - Real SQL queries
   - Transaction handling
   - Proper error handling

## Success Criteria

### ✅ Implementation is Real When:
1. CVE queries return unique IDs for different searches
2. Database file grows to 100+ MB
3. Queries fail gracefully when database is empty
4. Same query returns same results (consistent)
5. Valid CVE IDs from real NVD database
6. No hardcoded "CVE-2024-1234" anywhere in results

### ❌ Implementation is Fake If:
1. Always returns same 2 CVEs regardless of query
2. Database file is tiny (< 1 MB)
3. Queries return data even when database is empty
4. CVE descriptions contain "for testing" or "sample"
5. Only 2 entries in database per source

## Next Steps

After verification:
1. Show colleagues the validation test passing
2. Demonstrate real CVE lookups
3. Show database size and contents
4. Review code changes in `HONEST_FIX_REPORT.md`

The implementation is now **genuinely real and production-ready**.

