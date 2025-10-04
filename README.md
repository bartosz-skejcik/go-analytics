# Go Analytics

A lightweight, privacy-focused web analytics backend service built with Go. This project provides a simple yet powerful alternative to traditional analytics platforms, allowing you to track user sessions, page views, and custom events while maintaining user privacy through anonymous tracking.

## Overview

Go Analytics is a self-hosted analytics solution designed to give you full control over your data. It features a RESTful API backend that can be easily integrated into any web application, with support for session tracking, page view monitoring, and custom event logging.

## Key Features

- **Privacy-Focused Anonymous Tracking**: Generates daily-rotating anonymous IDs using IP and user agent hashing
- **Session Management**: Comprehensive session tracking with geolocation, device, and browser information
- **Page View Analytics**: Track page visits with time-on-page metrics and navigation flow
- **Custom Event Logging**: Flexible event system for tracking any custom interactions
- **IP Geolocation**: Automatic location detection (country, region, city) for visitor insights
- **PostgreSQL Storage**: Reliable and scalable data persistence
- **RESTful API**: Clean, well-documented endpoints for easy integration
- **Vercel Deployment Ready**: Configured for serverless deployment on Vercel
- **Docker Support**: Includes Docker Compose setup for local development

## Technology Stack

- **Backend**: Go 1.22.5
- **Web Framework**: Gin (high-performance HTTP framework)
- **Database**: PostgreSQL
- **Deployment**: Vercel (serverless), Docker
- **Libraries**: 
  - `gin-gonic/gin` - HTTP web framework
  - `lib/pq` - PostgreSQL driver
  - `rs/cors` - CORS middleware

## API Endpoints

### Session Tracking
- `POST /api/session` - Record a new user session with page views
- `GET /api/sessions` - Retrieve session data (supports `?distinct=true` for unique visitors)

### Event Tracking
- `POST /api/event` - Log a custom event with arbitrary data
- `GET /api/events` - Retrieve events within a time range

### Analytics Data
- `GET /api/pageviews` - Get page view statistics with optional time filters

## Database Schema

The service automatically creates three main tables:
- **sessions**: Stores user session information including anonymous IDs, timestamps, geolocation, and device data
- **page_views**: Records individual page visits with time spent and navigation order
- **events**: Captures custom events with flexible JSON data storage

## Getting Started

### Local Development with Docker

```bash
# Start PostgreSQL database
docker-compose up -d

# Run the application
go run cmd/main.go
```

### Deployment to Vercel

The project is configured for Vercel serverless deployment with environment variables for PostgreSQL connection.

## Use Cases

- Track website visitor behavior without invasive cookies
- Monitor user flows through your application
- Analyze content engagement with time-on-page metrics
- Log custom business events (button clicks, form submissions, etc.)
- Generate insights about your audience's location and devices

## Future Roadmap

See [TODO.md](TODO.md) for planned features including:
- NPM package for easy frontend integration
- Analytics dashboard with Laravel
- Multi-tenant support with API key authentication

## License

This project is open source and available for personal and commercial use.
