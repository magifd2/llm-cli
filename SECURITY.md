# Security Policy

## Reporting a Vulnerability

We take the security of `llm-cli` seriously. If you discover a security vulnerability within this project, we encourage you to report it to us as quickly as possible. We are committed to addressing all legitimate vulnerabilities in a timely manner.

**Please do NOT open a public GitHub issue for security vulnerabilities.**

To report a vulnerability, please refer to the security contact information associated with the project's GitHub repository. You can usually find this on the repository's main page or in the project's `README.md`.

### What to Include in Your Report

To help us quickly understand and resolve the issue, please include the following information in your report:

*   **Description of the vulnerability**: A clear and concise description of the vulnerability.
*   **Steps to reproduce**: Detailed steps on how to reproduce the vulnerability. This should include any specific configurations, commands, or inputs required.
*   **Impact**: Explain the potential impact of the vulnerability (e.g., data exposure, unauthorized access, denial of service).
*   **Affected versions**: Specify which versions of `llm-cli` are affected.
*   **Proof of concept (optional but helpful)**: If possible, provide a proof-of-concept (PoC) code or a demonstration that illustrates the vulnerability.
*   **Your contact information (optional)**: If you wish to be credited for your discovery, please include your name or handle.

### Our Commitment

Upon receiving a vulnerability report, we will:

1.  Acknowledge receipt of your report within **7 business days**.
2.  Investigate the reported vulnerability promptly.
3.  Keep you informed of our progress and any remediation plans.
4.  Notify you when the vulnerability has been resolved.

## Security Principles

`llm-cli` adheres to the following security principles:

*   **Security First**: Security is the highest priority, overriding all other considerations such as functionality or performance.
*   **Input Validation**: Never trust user input, including environment variables. All external inputs must be validated and sanitized to prevent injection attacks.
*   **Sensitive Information Handling**: Sensitive information (API keys, credentials) must never be hardcoded or stored in insecure locations.
*   **Dependency Security**: All code, dependencies, and configurations must be reviewed for potential security vulnerabilities before being committed. Regularly scan project dependencies for known vulnerabilities.
*   **Secure Development Lifecycle**: Incorporate threat modeling, security-focused code reviews, and safe testing practices throughout the development process.
*   **Secure Configuration Storage**: Configuration files containing sensitive information (e.g., API keys) must be stored with restrictive file and directory permissions (e.g., `0600` for files, `0700` for directories) to prevent unauthorized access.
