# Scraper Service

A Go-based web scraping service that analyzes webpages and provides detailed information about their structure and content.

## API Documentation

The API provides endpoints for analyzing webpages and retrieving metrics.

### Main Endpoints:
- `GET /api/v1/analyze`: Analyze a webpage by providing a URL
- `GET /api/v1/system/metrics`: Get Prometheus metrics

## Prerequisites

- [Go 1.24.3](https://golang.org/dl/) - The Go programming language
- [Google Chrome](https://www.google.com/chrome/) - Required for headless browser automation
- [Docker](https://www.docker.com/) (optional) - For running the complete setup with monitoring and logging
- [Node.js](https://nodejs.org/) (optional) - For commitlint and other development tools

## Technologies Used

### Backend
- [Go 1.24.3](https://golang.org/) - Programming language
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [go-rod](https://github.com/go-rod/rod) - Browser automation library for web scraping
- [Zap](https://github.com/uber-go/zap) - Structured logging
- [Viper](https://github.com/spf13/viper) - Configuration management
- [OpenTelemetry](https://opentelemetry.io/) - Distributed tracing (WIP)
- [Gin-contrib/cache](https://github.com/gin-contrib/cache) - Response caching

### DevOps & Monitoring
- [Docker](https://www.docker.com/) - Containerization
- [Prometheus](https://prometheus.io/) - Metrics collection
- [Grafana](https://grafana.com/) - Metrics and logs visualization
- [Air](https://github.com/air-verse/air) - Live reload for development

### Frontend
- The frontend client is available in a separate repository and is included in the Docker Compose setup

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mrmihi/web-analyzer.git
   cd web-analyzer
   ```

### Development

- Run the development server with hot reloading:
   ```bash
   make dev
   ```

### Docker Setup

- Start the complete setup with Docker Compose depending on your OS:
   ```bash
   make sandbox-linux
   make sandbox-windows
   ```

- Stop all services:
   ```bash
   make teardown-linux
   make teardown-windows
   ```
- Access the services:
  - Client: [http://localhost:5173](http://localhost:5173)
  - Scraper API: [http://localhost:8080](http://localhost:8080)
  - Grafana: [http://localhost:3000](http://localhost:3000)

## Usage

### Analyzing a Webpage

Send a GET request to `/api/v1/analyze/` with a URL query parameter:

```bash
curl --location 'http://localhost:8080/api/v1/analyze/?url=https://mrmihi.dev'
```

Example response:

```json
{
  "html_version": "HTML 5",
  "title": "Example Domain",
  "headings": {
    "h1": 1,
    "h2": 0,
    "h3": 0,
    "h4": 0,
    "h5": 0,
    "h6": 0
  },
  "internal_links": 1,
  "external_links": 1,
  "inaccessible_links": 0,
  "login_form": false
}
```

## Project Structure

- `cmd/` - HTTP Server initialization
- `common/` - Common utilities and error handling
- `config/` - Configuration management
- `dto/` - Data Transfer Objects
- `handlers/` - HTTP handlers and controllers
- `infra/` - Infrastructure configuration
  - `docker-compose.yml` - Docker Compose configuration
- `internal/` - Internal packages
  - `logger/` - Logging utilities
  - `scraper/` - Web scraping functionality
- `middleware/` - Middleware functions
- `routes/` - HTTP routes and server setup
- `services/` - Business logic services
- `tests/` - Tests cases

## Main Features

1. **Webpage Analysis**
   - HTML version detection
   - Page title extraction
   - Heading counts (h1-h6)
   - Internal and external link counting
   - Login form detection

2. **Monitoring and Observability**
   - Prometheus metrics
   - Grafana dashboards
   - Structured logging with Zap
   - Distributed tracing with OpenTelemetry (WIP)

3. **Robust Error Handling**
   - Request validation
   - Timeout handling
   - Graceful shutdown

## Challenges and Solutions

### Challenge 1: Headless Browser Automation
- **Problem**: Controlling a headless browser for web scraping can be resource-intensive and prone to timeouts.
- **Solution**: Implemented resource blocking for non-essential content (images, stylesheets, etc.) and added a 120-second timeout for analysis operations to ensure the API remains responsive.

### Challenge 2: Link Analysis
- **Problem**: Analyzing all links on a webpage could lead to excessive resource usage for pages with many links.
- **Solution**: Implemented concurrent processing of links using goroutines to improve performance while maintaining accuracy.


## Possible Improvements

1. **Authentication**: Add authentication for API access.

2. **More Analysis Features**: Add more webpage analysis features such as:
   - Image analysis
   - SEO metrics
   - Performance metrics
   - Accessibility checks

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
