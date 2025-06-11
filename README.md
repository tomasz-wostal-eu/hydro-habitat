[![CI, Release and Publish](https://github.com/tomasz-wostal-eu/hydro-habitat/actions/workflows/release.yaml/badge.svg)](https://github.com/tomasz-wostal-eu/hydro-habitat/actions/workflows/release.yaml)

# Hydro Habitat 🐠

A modern aquarium tank management system built with Go backend and React frontend. Hydro Habitat helps aquarium enthusiasts track and manage their tank inventory with features for monitoring water types, locations, and detailed notes.

## 🚀 Features

- **Tank Management**: Create, edit, delete, and view aquarium tanks
- **Water Type Tracking**: Support for tap, RO (Reverse Osmosis), and RO/DI water
- **Location Management**: Track tank placement by room and rack location
- **Inventory System**: Assign and track inventory numbers
- **Notes & Documentation**: Add detailed notes for each tank
- **Responsive UI**: Modern, mobile-friendly interface built with React and Tailwind CSS
- **REST API**: Full RESTful API with Swagger documentation
- **Database**: PostgreSQL with UUID primary keys and ENUM types

## 🏗️ Architecture

### Backend (Go)
- **Framework**: Gin web framework
- **Database**: PostgreSQL with sqlx for database operations
- **Documentation**: Swagger/OpenAPI integration
- **Testing**: Comprehensive unit and integration tests
- **Architecture**: Clean architecture with separate domain, store, and API layers

### Frontend (React)
- **Framework**: React 18 with Vite
- **Styling**: Tailwind CSS with Catppuccin color scheme
- **HTTP Client**: Axios for API communication
- **Testing**: Vitest with React Testing Library
- **Icons**: Lucide React icons

### Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose for development and testing
- **Database**: PostgreSQL 16 Alpine
- **Web Server**: Caddy for frontend serving
- **Load Testing**: K6 performance tests

## 📦 Quick Start

### Prerequisites
- Docker and Docker Compose
- Make (optional, for convenience commands)

### 1. Clone the Repository
```bash
git clone <repository-url>
cd hydro_habitat
```

### 2. Start the Application
```bash
# Using Make (recommended)
make up

# Or using Docker Compose directly
docker-compose up --build -d
```

### 3. Access the Application
- **Frontend**: http://localhost (port 80)
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html

### 4. View Logs
```bash
make logs
# Or: docker-compose logs -f
```

## 🛠️ Available Commands

### Application Management
```bash
make up               # Start all services in the background
make down             # Stop and remove all services
make logs             # Display logs for all services
```

### Testing
```bash
make test-unit        # Run backend unit tests
make test-integration # Run backend integration tests
```

### Code Quality
```bash
make lint             # Run linting for both backend and frontend
make lint-backend     # Run Go linting checks
make lint-frontend    # Run ESLint and Prettier checks
```

### Security Scanning
```bash
make security         # Run security scanning for backend
make security-backend # Run Go security scanning with gosec
```

> **Note**: All `make` commands use `docker-compose` for local development. Our GitHub Actions workflows use `docker compose` (space instead of hyphen) for compatibility with newer CI environments.

### 5. Stop the Application
```bash
make down
# Or: docker-compose down --volumes --remove-orphans
```

## 🧪 Testing

### Backend Tests

#### Unit Tests (No Database Required)
```bash
make test-unit
# Or: docker-compose run --rm backend-unit-test
```

#### Integration Tests (With Database)
```bash
make test-integration
# Or: docker-compose run --rm backend-integration-test
```

#### Running Tests Locally
```bash
cd backend
go test ./...                    # All tests
go test -v ./api ./config        # Specific packages
go test -tags=integration ./...  # Integration tests only
```

**Note**: Integration tests require an active database connection. When running locally, you need to provide the proper database credentials via environment variables (e.g., `DATABASE_URL`). If not configured, integration tests will be skipped automatically. Use `make test-integration` to run integration tests with the proper database configuration.

### Frontend Tests
```bash
# Run all frontend tests
docker-compose run --rm frontend-test

# Or run locally in frontend directory
cd frontend
npm install
npm run test:run     # Run once
npm run test         # Watch mode
npm run test:ui      # UI mode
npm run test:coverage # With coverage
```

### Load Testing
```bash
# Install K6 first, then run smoke tests
cd k6-tests
k6 run smoke-test.js
```

## ✨ Code Quality & Linting

### Running All Linting
```bash
make lint  # Runs both backend and frontend linting
```

### Backend Linting (Go)
```bash
make lint-backend
# Or: docker-compose run --rm backend-lint
```

The backend uses **golangci-lint** with the following enabled linters:
- `govet` - Go vet examiner
- `errcheck` - Check for unchecked errors
- `staticcheck` - Static analysis checks
- `unused` - Check for unused code
- `gosimple` - Simplify code suggestions
- `ineffassign` - Detect ineffectual assignments
- `typecheck` - Type checking
- `goimports` - Import formatting
- `misspell` - Spell checking
- `gocyclo` - Cyclomatic complexity

### Frontend Linting (JavaScript/React)
```bash
make lint-frontend
# Or: docker-compose run --rm frontend-lint
```

The frontend uses **ESLint** with React-specific rules:
- React hooks linting
- React refresh plugin
- React-specific best practices
- Code formatting with Prettier integration

### Local Development Linting
```bash
# Backend (requires Go and golangci-lint installed)
cd backend
golangci-lint run ./...

# Frontend (requires Node.js)
cd frontend
npm run lint        # Run ESLint
npm run lint:fix    # Auto-fix ESLint issues
npm run format      # Format with Prettier
npm run format:check # Check Prettier formatting
```

## 🔒 Security Scanning

### Backend Security Scanning
```bash
make security         # Runs backend security scanning
make security-backend # Runs Go security scanning with gosec
# Or: docker-compose run --rm backend-security
```

The backend security scanning uses **gosec** (Go Security Checker) to identify:
- SQL injection vulnerabilities
- Command injection risks
- Path traversal issues
- Weak cryptographic practices
- Unsafe file operations
- Hardcoded credentials
- TLS configuration issues

### Manual Security Checks
```bash
# Scan with Trivy (requires Trivy installation)
trivy fs ./backend
trivy fs ./frontend

# Go vulnerability database check
cd backend && govulncheck ./...

# NPM security audit
cd frontend && npm audit
```

### Automated Security Features
- **GitHub Security Scanning**: Trivy scans run automatically on push/PR
- **Dependabot**: Automated dependency security updates
- **SARIF Upload**: Security results integrated with GitHub Security tab
- **Container Scanning**: Docker images scanned for vulnerabilities

See [`SECURITY.md`](SECURITY.md) for detailed security configuration and incident response procedures.

## 🗄️ Database Schema

The application uses PostgreSQL with the following main entities:

### Tanks Table
```sql
CREATE TABLE tanks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    room VARCHAR(255),
    rack_location VARCHAR(255),
    inventory_number VARCHAR(255),
    volume_liters INTEGER NOT NULL,
    water water_type NOT NULL DEFAULT 'tap',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Water Types
- `tap` - Tap water
- `ro` - Reverse Osmosis water
- `rodi` - RO/DI (Reverse Osmosis/Deionized) water

## 🛠️ Development

### Project Structure
```
hydro_habitat/
├── backend/                 # Go backend service
│   ├── api/                # HTTP handlers and routing
│   ├── config/             # Configuration management
│   ├── domain/             # Domain models and business logic
│   ├── store/              # Data access layer
│   ├── migrations/         # Database migrations
│   └── docs/               # Swagger documentation
├── frontend/               # React frontend application
│   └── src/
│       ├── test/           # Frontend tests
│       └── App.jsx         # Main application component
├── k6-tests/               # Load testing scripts
├── docker-compose.yaml     # Container orchestration
└── Makefile               # Build automation
```

### API Endpoints

#### Tanks
- `GET /api/v1/tanks` - List all tanks
- `POST /api/v1/tanks` - Create a new tank
- `GET /api/v1/tanks/{id}` - Get tank by ID
- `PUT /api/v1/tanks/{id}` - Update tank
- `DELETE /api/v1/tanks/{id}` - Delete tank

#### Health Check
- `GET /health` - Health check endpoint

### Environment Variables
```bash
# Database Configuration
POSTGRES_USER=hydro           # Default: hydro
POSTGRES_PASSWORD=hydrosecret # Default: hydrosecret
POSTGRES_DB=hydro_habitat     # Default: hydro_habitat

# Backend Configuration
DATABASE_URL=postgres://...   # Auto-configured in Docker
GIN_MODE=release             # Gin framework mode
```

### Local Development

#### Backend Development
```bash
cd backend
go mod download
go run main.go
```

#### Frontend Development
```bash
cd frontend
npm install
npm run dev  # Starts Vite dev server on http://localhost:5173
```

## 🐳 Docker Services

The application consists of several Docker services:

### Production Services
- **postgres**: PostgreSQL database
- **backend**: Go API server (production build)
- **frontend**: React application served by Caddy

### Test Services
- **backend-unit-test**: Runs Go unit tests
- **backend-integration-test**: Runs Go integration tests with database
- **frontend-test**: Runs React/Vitest tests

### Code Quality Services
- **backend-lint**: Runs golangci-lint for Go code quality checks
- **frontend-lint**: Runs ESLint for React/JavaScript code quality checks

### Security Services
- **backend-security**: Runs gosec security scanner for Go code vulnerabilities

## 📊 Monitoring and Health Checks

- Database health checks ensure PostgreSQL is ready before starting dependent services
- Backend exposes `/health` endpoint for monitoring
- All services include restart policies for resilience

## 🔧 Configuration

### Docker Compose Profiles
The application uses a single `docker-compose.yaml` file with different services for different purposes:
- Production services start automatically with `docker-compose up`
- Test services are run on-demand with `docker-compose run`

### Multi-stage Docker Builds
- **Backend**: Uses multi-stage builds with separate test and production stages
- **Frontend**: Optimized production builds with Caddy web server

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Write tests for new features
- Follow existing code style and patterns
- Run linting before committing (`make lint`)
- Run security scanning before major releases (`make security`)
- Ensure all tests pass before submitting PR (`make test-unit test-integration`)
- Update documentation as needed
- Review security alerts and dependencies regularly

### Pull Request Process
Our PR workflow automatically runs comprehensive quality checks:

#### **Automated Checks Include:**
- **Code Quality**: Backend linting (golangci-lint) and frontend linting (ESLint)
- **Security Scanning**: Vulnerability detection with gosec and Trivy
- **Testing**: Unit tests, integration tests, and full application testing
- **Build Verification**: Ensures all Docker services build successfully
- **API Integration**: Verifies backend endpoints and frontend connectivity

#### **Smart Change Detection:**
- Only runs relevant checks based on changed files
- Backend changes trigger Go-specific quality checks
- Frontend changes trigger React/JavaScript quality checks
- Workflow changes are automatically detected

#### **Required Status Checks:**
- All linting must pass (0 issues)
- All tests must pass (unit + integration)
- Security scans must show no critical vulnerabilities
- Application must start successfully and respond to health checks

PRs cannot be merged until all automated checks pass. You can run the same checks locally using the available `make` commands.

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# View database logs
docker-compose logs postgres

# Reset database
make down && make up
```

#### Frontend Build Issues
```bash
# Clear node modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
```

#### Port Conflicts
If you encounter port conflicts:
- Frontend (port 80): Change in docker-compose.yaml under frontend service
- Backend (port 8080): Change in docker-compose.yaml under backend service
- Database (port 5432): Change in docker-compose.yaml under postgres service

### Useful Commands
```bash
# View all running containers
docker-compose ps

# Execute commands in running containers
docker-compose exec backend sh
docker-compose exec postgres psql -U hydro -d hydro_habitat

# Remove all data (including database)
docker-compose down --volumes --remove-orphans
docker system prune -f
```
---

Built with ❤️ for the aquarium community
