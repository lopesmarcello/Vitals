# Vitals ✨

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

### Live Demo
[https://vitals-2osr.onrender.com/](https://vitals-2osr.onrender.com/)

Vitals is a web performance analyzer that provides insights into your website's performance and health. It's a great tool to identify potential issues and optimize your site for a better user experience.

## How it works

Vitals uses a combination of techniques to analyze your website:

- **Browser Analysis:** It uses a headless Chrome browser to load your website and collect performance metrics like loading time, resource usage, and more.
- **Network Analysis:** It checks the network requests made by your website, providing information about the status of each request, the time it took, and the size of the resources.
- **Link Checker:** It crawls your website to find all the links and checks if they are working correctly.

## Getting Started

You can run Vitals locally using Docker.

1. **Clone the repository:**

   ```bash
   git clone https://github.com/lopesmarcello/vitals.git
   ```

2. **Build and run the Docker container:**

   ```bash
   docker-compose up --build -d
   ```

3. **Open your browser and go to `http://localhost:3000`**

## TODO

- [ ] Turn "check links" optional
- [ ] Centralized Configuration
- [ ] Structured Logging
- [ ] Comprehensive Testing
- [ ] Hardened Dockerfile
- [ ] CI/CD Pipeline
- [ ] Security Hardening
- [ ] Observability
- [ ] Code Quality and Refactoring

## Contributing

Contributions are welcome! Please see our [Contributing Guidelines](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
