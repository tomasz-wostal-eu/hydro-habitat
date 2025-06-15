# Security Configuration and Workflow Guide

## Overview

This repository is configured with comprehensive security scanning and dependency management to ensure the safety and integrity of the Hydro Habitat application.

## Security Features Implemented

### 1. Automated Security Scanning

**Trivy Security Scanner** (`security-scan.yaml`)
- **Filesystem scanning**: Scans source code for vulnerabilities and secrets
- **Docker image scanning**: Scans built container images for vulnerabilities
- **SARIF upload**: Results are uploaded to GitHub Security tab
- **Schedule**: Runs on push, PR, and daily at 2 AM UTC

**Scan Coverage:**
- Go dependencies (backend)
- NPM dependencies (frontend)
- Docker base images
- Secret detection in source code

### 2. Dependency Management

**Dependabot** (`.github/dependabot.yml`)
- **Automated updates**: Weekly dependency updates
- **Multi-ecosystem**: Go modules, NPM packages, Docker, GitHub Actions
- **Security-focused**: Prioritizes security updates
- **Controlled**: Limited PR count to avoid overwhelming

**Dependency Security Updates** (`dependency-security.yaml`)
- **Weekly automated**: Runs security-focused dependency updates
- **Vulnerability checking**: Uses `govulncheck` for Go and `npm audit` for NPM
- **Automated PRs**: Creates pull requests with security updates
- **Test verification**: Ensures updates don't break functionality

### 3. Code Quality & Linting

**Existing Linting** (Status: ✅ Complete)
- **Backend**: `golangci-lint` with zero issues
- **Frontend**: `eslint` with zero issues
- **Test compliance**: All 74 tests passing

## Repository Settings Required

### 1. Enable GitHub Advanced Security

For private repositories, enable these features in Settings → Security & analysis:

```yaml
Required Settings:
✅ Dependency graph
✅ Dependabot alerts
✅ Dependabot security updates
✅ Code scanning
✅ Secret scanning
```

### 2. Branch Protection Rules

Recommended protection for `main` branch:

```yaml
Branch Protection:
✅ Require status checks to pass
✅ Require up-to-date branches
✅ Include security scan results
✅ Restrict pushes to main
```

## Workflow Permissions

All security workflows include proper permissions:

```yaml
permissions:
  contents: read          # Read repository contents
  security-events: write  # Upload SARIF results
  actions: read          # Required for private repos
  id-token: write        # Required for OIDC token
```

## Monitoring & Alerts

### Security Dashboard Locations

1. **Security Tab**: View all vulnerability alerts and scan results
2. **Actions Tab**: Monitor workflow execution and logs
3. **Pull Requests**: Automatic security check status
4. **Dependabot Tab**: Review dependency update PRs

### Alert Severity Levels

- **CRITICAL**: Immediate action required
- **HIGH**: Address within 1 week
- **MEDIUM**: Address within 1 month
- **LOW**: Address during regular maintenance

## Manual Security Checks

### Local Security Scanning

```bash
# Scan backend with Trivy
docker run --rm -v $(pwd):/workspace aquasec/trivy fs /workspace/backend

# Scan frontend with Trivy
docker run --rm -v $(pwd):/workspace aquasec/trivy fs /workspace/frontend

# Go vulnerability check
cd backend && govulncheck ./...

# NPM security audit
cd frontend && npm audit
```

### Update Dependencies Manually

```bash
# Update Go dependencies
cd backend && go get -u all && go mod tidy

# Update NPM dependencies
cd frontend && npm update && npm audit fix
```

## Incident Response

### When Security Issues Are Found

1. **Review the alert** in GitHub Security tab
2. **Assess severity** and impact on your application
3. **Update dependencies** using automated PRs or manual updates
4. **Test thoroughly** after applying security fixes
5. **Deploy updates** as soon as possible for critical issues

### Emergency Security Updates

For critical vulnerabilities:

1. Create hotfix branch: `git checkout -b hotfix/security-fix`
2. Apply the security patch
3. Run all tests: `make test` (if available)
4. Create emergency PR with fast-track review
5. Deploy immediately after merge

## Best Practices

### Development Workflow

- ✅ Never commit secrets or API keys
- ✅ Review Dependabot PRs promptly
- ✅ Monitor security alerts weekly
- ✅ Keep dependencies up to date
- ✅ Run security scans before major releases

### Container Security

- ✅ Use minimal base images (Alpine, Distroless)
- ✅ Regular base image updates
- ✅ Multi-stage builds to reduce attack surface
- ✅ Non-root user in containers

## Troubleshooting

### SARIF Upload Issues

If you encounter "Resource not accessible by integration":

1. **Check permissions**: Ensure workflow has `security-events: write`
2. **Enable Advanced Security**: Required for private repos
3. **Verify token scope**: GITHUB_TOKEN needs sufficient permissions
4. **Re-run workflow**: Sometimes transient GitHub API issues

### Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| Dependabot PRs failing | Check if tests pass, review dependency conflicts |
| Security scan timeouts | Reduce scan scope or increase timeout values |
| False positive alerts | Review and dismiss non-applicable vulnerabilities |
| High vulnerability count | Prioritize by severity, update critical issues first |

## Support & Resources

- **GitHub Security Documentation**: https://docs.github.com/en/code-security
- **Trivy Documentation**: https://trivy.dev/
- **Dependabot Documentation**: https://docs.github.com/en/code-security/dependabot
- **Go Security**: https://go.dev/security/
- **NPM Security**: https://docs.npmjs.com/auditing-package-dependencies-for-security-vulnerabilities
