# Security Policy

## Supported Versions

Currently, only the latest version of budge is actively supported with security updates.

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < latest| :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in budge, please report it privately:

1. **Email**: Create an issue on GitHub with a generic title like "Security concern" and request private disclosure
2. **Provide details**: Include a description of the vulnerability, steps to reproduce, and potential impact
3. **Response time**: We aim to respond within 48 hours

### What to Include

- Type of vulnerability (e.g., credential exposure, connection security)
- Steps to reproduce
- Potential impact
- Suggested fix (if you have one)

## Security Considerations

### Configuration File Security

budge stores IMAP credentials in `~/.config/budge/config.yaml`. To protect your credentials:

```bash
chmod 600 ~/.config/budge/config.yaml
```

budge will warn you if the permissions are not set correctly.

### Credential Storage

- **Never commit** your `config.yaml` to version control
- Use app-specific passwords (not your main account password)
- The repository includes a `.gitignore` to prevent accidental commits

### IMAP Connection Security

- Always use TLS (`tls: true` on port 993) or STARTTLS for email connections
- Avoid plain-text IMAP connections unless on localhost (e.g., ProtonMail Bridge)

## Disclosure Policy

Once a vulnerability is fixed:
1. We'll create a GitHub Security Advisory
2. Release a patched version
3. Publicly disclose after users have had time to update

Thank you for helping keep budge secure!
