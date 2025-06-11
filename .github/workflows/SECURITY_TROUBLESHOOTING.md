# Security Scanning Troubleshooting Guide

## Issue: "Resource not accessible by integration" when uploading SARIF results

### Root Cause
The error occurs when GitHub Actions workflows attempt to upload security scan results (SARIF format) to GitHub's Code Scanning feature but lack the necessary permissions.

### Solutions Applied

1. **Created Security Scan Workflow** (`security-scan.yaml`)
   - Added proper permissions: `security-events: write`
   - Configured Trivy scanning for both filesystem and Docker images
   - Set up proper SARIF upload with categorization

2. **Required Permissions for SARIF Upload**
   ```yaml
   permissions:
     contents: read
     security-events: write  # Required for uploading SARIF results
     actions: read          # Required for private repos
     id-token: write        # Required for OIDC token
   ```

### Repository Settings to Verify

1. **Enable GitHub Advanced Security** (if using private repo):
   - Go to Settings → Security & analysis
   - Enable "GitHub Advanced Security"
   - Enable "Code scanning"

2. **Check Branch Protection Rules**:
   - Ensure security scans don't block necessary workflows
   - Configure status checks appropriately

3. **Token Permissions**:
   - Verify GITHUB_TOKEN has sufficient permissions
   - For organization repos, check organization security policies

### Manual Resolution Steps

If you're still seeing the error:

1. **Check Repository Settings**:
   ```bash
   # Verify Advanced Security is enabled for your repo
   gh api repos/:owner/:repo --jq '.security_and_analysis'
   ```

2. **Re-run Failed Workflow**:
   - Go to Actions tab
   - Find the failed workflow
   - Click "Re-run all jobs"

3. **Test Security Scanning**:
   ```bash
   # Test Trivy locally first
   docker run --rm -v $(pwd):/workspace aquasec/trivy fs /workspace --format sarif
   ```

### Expected Outcome

After implementing these changes:
- ✅ Security scans will run on push/PR/schedule
- ✅ SARIF results will upload to GitHub Security tab
- ✅ Vulnerabilities will be visible in the Security overview
- ✅ No more "Resource not accessible" errors

### Monitoring

The security workflow will now run:
- On every push to main/dev/develop branches
- On pull requests to main
- Daily at 2 AM UTC (scheduled scan)

Results will appear in:
- **Security tab** → Code scanning alerts
- **Pull request checks** (if scan finds issues)
- **Actions tab** (workflow execution logs)
