# Security Policy

## 🛡️ Our Commitment to Security

The Shantilly project takes security seriously. We are committed to maintaining a secure codebase and protecting our users.

## 🔒 Reporting Security Vulnerabilities

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities by:

1. **Email:** Send detailed information to: [security@shantilly.dev](mailto:security@shantilly.dev)
2. **Response Time:** We will respond within 48 hours
3. **Coordination:** We'll work with you to understand and address the vulnerability

### What to Include

When reporting a vulnerability, please include:

- **Description:** Detailed description of the vulnerability
- **Impact:** Potential impact if exploited
- **Reproduction:** Steps to reproduce the vulnerability
- **Environment:** Operating system, Shantilly version, configuration
- **Proof of Concept:** Code or screenshots demonstrating the vulnerability

## 🛠️ Supported Versions

We actively support the following versions of Shantilly with security updates:

| Version | Supported          |
| ------- | ------------------ |
| v1.x.x  | ✅                |
| v0.x.x  | ❌ (End of Life)  |

## 🔐 Security Features

### Built-in Security

- **No External Dependencies for Core:** Shantilly core functionality doesn't require external network access
- **Input Validation:** All YAML configurations are validated before processing
- **Sandboxed Execution:** Script execution is contained and monitored
- **No Persistent Storage:** Shantilly doesn't store sensitive data

### Best Practices for Users

1. **YAML Configuration:**
   - Validate YAML syntax before processing
   - Don't include sensitive data in YAML files
   - Use file permissions to protect configuration files

2. **Script Execution:**
   - Only execute scripts from trusted sources
   - Review scripts before execution
   - Use principle of least privilege for file access

3. **Environment:**
   - Keep Shantilly updated to the latest version
   - Use secure terminals that support modern encryption
   - Monitor system logs for unusual activity

## 🚨 Vulnerability Response Process

### 1. **Initial Response** (Within 48 hours)
- Acknowledge receipt of the report
- Assign a tracking number
- Begin investigation

### 2. **Investigation** (1-7 days)
- Reproduce and verify the vulnerability
- Assess the scope and impact
- Develop a fix or mitigation

### 3. **Fix Development** (1-14 days)
- Develop and test the fix
- Ensure the fix doesn't introduce new vulnerabilities
- Prepare documentation

### 4. **Disclosure** (After fix is ready)
- Coordinate with the reporter on disclosure timing
- Release security update
- Publish security advisory
- Credit the reporter (unless they prefer anonymity)

## 🛡️ Security Measures

### Code Security

- **Static Analysis:** Automated security scanning with CodeQL
- **Dependency Scanning:** Automated vulnerability scanning of dependencies
- **Secret Scanning:** Detection of accidentally committed secrets
- **Regular Audits:** Periodic security reviews of the codebase

### Infrastructure Security

- **Protected Branches:** Main branch requires approvals and status checks
- **Signed Commits:** All commits must be GPG signed
- **Dependency Updates:** Automated security updates via Dependabot
- **Security Policies:** Enforced through GitHub security features

## 📋 Security Checklist

For Contributors:

- [ ] All code changes reviewed by at least one other developer
- [ ] No secrets (API keys, passwords, tokens) in code
- [ ] Input validation for all user-provided data
- [ ] Error messages don't leak sensitive information
- [ ] Tests cover security-related functionality
- [ ] Documentation includes security considerations

## 🔗 Related Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://golang.org/doc/faq#security)
- [GitHub Security Features](https://docs.github.com/en/code-security)
- [CVE Database](https://cve.mitre.org/)

## 📞 Contact

For security-related questions or concerns:

- **Security Email:** [security@shantilly.dev](mailto:security@shantilly.dev)
- **GitHub Security Advisories:** [Security Advisories](https://github.com/helton-godoy/shantilly/security/advisories)

## 🎯 Recognition

We believe in recognizing security researchers who help improve Shantilly's security. Depending on the nature of the disclosure, we may:

- Credit researchers in security advisories
- Include researchers in release notes
- Offer private acknowledgment in our README
- Provide security researcher rewards for significant discoveries

**Note:** All security communications will be treated confidentially until a coordinated disclosure is agreed upon.

---

*This security policy is based on industry best practices and will be updated as needed to reflect the evolving security landscape.*
