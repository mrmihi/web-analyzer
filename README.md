# Scraper Service

A Go-based web scraping service that analyzes webpages and provides detailed information about their structure and content.

## API Documentation

The API provides endpoints for analyzing webpages and retrieving metrics.

### Main Endpoints:
- `POST /api/v1/analyze`: Analyze a webpage by providing a URL
- `GET /api/v1/system/metrics`: Get Prometheus metrics
- `GET /api/v1/example`: Example endpoint

## Prerequisites

- [Go 1.24+](https://golang.org/dl/) - The Go programming language
- [Google Chrome](https://www.google.com/chrome/) - Required for headless browser automation
- [Docker](https://www.docker.com/) (optional) - For running the complete setup with monitoring and logging
- [Node.js](https://nodejs.org/) (optional) - For commitlint and other development tools

## Technologies Used

### Backend
- [Go 1.24](https://golang.org/) - Programming language
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [go-rod](https://github.com/go-rod/rod) - Browser automation library for web scraping
- [Zap](https://github.com/uber-go/zap) - Structured logging
- [Viper](https://github.com/spf13/viper) - Configuration management
- [OpenTelemetry](https://opentelemetry.io/) - Distributed tracing

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
   git clone https://github.com/yourusername/scraper.git
   cd scraper
   ```

### Development

- Run the development server with hot reloading:
   ```bash
   make dev
   ```

### Docker Setup

- Start the complete setup with Docker Compose:
   ```bash
   make sandbox
   ```

- Stop all services:
   ```bash
   make teardown
   ```

- Access the services:
  - Client: [http://localhost:5173](http://localhost:5173)
  - Scraper API: [http://localhost:8080](http://localhost:8080)
  - Grafana: [http://localhost:3000](http://localhost:3000)

## Usage

### Analyzing a Webpage

Send a POST request to `/api/v1/` with a JSON body containing the URL to analyze:

```bash
curl -X POST http://localhost:8080/api/v1/analyze\
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
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

- `common/` - Common utilities and error handling
- `config/` - Configuration management
- `consts/` - Service constants
- `dto/` - Data Transfer Objects
- `handlers/` - HTTP handlers and controllers
- `infra/` - Infrastructure configuration
  - `docker-compose.yml` - Docker Compose configuration
- `internal/` - Internal packages
  - `logger/` - Logging utilities
  - `scraper/` - Web scraping functionality (analyzer, utils)
- `middleware/` - Middleware functions
- `routes/` - HTTP routes and server setup
- `services/` - Business logic services

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
- **Solution**: Implemented a page pool for parallel processing and added timeout handling to ensure the API remains responsive.

### Challenge 2: Link Analysis
- **Problem**: Analyzing all links on a webpage could lead to excessive resource usage for pages with many links.
- **Solution**: Limited the number of links to analyze to 50 and implemented concurrent processing with goroutines.


## Possible Improvements

1. **Caching**: Implement response caching to improve performance for frequently requested URLs.
2. **Rate Limiting**: Enhance the rate limiting middleware to prevent abuse.
3. **Authentication**: Add authentication for API access.
4. **More Analysis Features**: Add more webpage analysis features such as:
   - Image analysis
   - SEO metrics
   - Performance metrics
   - Accessibility checks
5. **API Documentation**: Add Swagger/OpenAPI documentation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
