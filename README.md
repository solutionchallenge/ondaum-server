# Ondaum Server

> 📅 This README was written on **May 15, 2025**.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Google Gemini](https://img.shields.io/badge/google%20gemini-8E75B2?style=for-the-badge&logo=google%20gemini&logoColor=white)
![Google Login](https://img.shields.io/badge/google-4285F4?style=for-the-badge&logo=google&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![MySQL](https://img.shields.io/badge/mysql-4479A1.svg?style=for-the-badge&logo=mysql&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)

## 🌍 OVERVIEW
Hello👋 This is team Ondaum. ***Ondaum*** is a pure Korean word, meaning ***'a warm and genuine self'***.

We want to help people around the world live healthier lives by being with Um, an AI professional psychological counseling companion, anytime and anywhere.

Let's start https://ondaum.revimal.me/

## 🛠 SKILLS

### Architecture & Design
- **Architecture Pattern**: Vertical Slice Architecture
- **Design Methodology**: Domain-Driven Design (DDD)
- **Dependency Injection**: [Uber Fx](https://github.com/uber-go/fx) - A dependency injection system for Go applications

### Backend Development
- **Language**: Go
- **Command Line**: [spf13/cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- **Configuration**: [spf13/viper](https://github.com/spf13/viper) - Go configuration with fangs
- **Testing**: [uber-go/gomock](https://github.com/uber-go/mock) - A mocking framework for Go interfaces

### Database & ORM
- **Main Database**: MySQL
- **ORM**: [Bun](https://github.com/uptrace/bun) - A fast and simple ORM for Go
- **Development DB**: [dolthub/go-mysql-server](https://github.com/dolthub/go-mysql-server) - In-memory MySQL server for development and testing

### API & Communication
- **HTTP Framework**: [gofiber/fiber](https://github.com/gofiber/fiber) - Express inspired web framework built on top of Fasthttp
- **API Style**: REST API
- **API Documentation**: [swaggo/swag](https://github.com/swaggo/swag) - Swagger documentation generator for Go
- **Live API Documentation**: [Swagger UI](https://ondaum.revimal.me/api/v1/_sys/swagger)
- **Authentication**: 
  - OAuth 2.0
  - [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - Go implementation of JSON Web Tokens

### AI/ML Integration
- **LLM**: [google.golang.org/genai](https://github.com/googleapis/go-genai) - Official Go client library for Google's Generative AI API

### External Libraries
- [spf13/cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- [spf13/viper](https://github.com/spf13/viper) - Go configuration with fangs
- [dolthub/go-mysql-server](https://github.com/dolthub/go-mysql-server) - In-memory MySQL server for development and testing

## 📁 DIRECTORY

```
.
├── cmd/                # Application running commands
│
├── config/             # Configuration files
│
├── docs/               # Swagger documentations
│
├── internal/           # Private application code
│   ├── domain/         # Domain models and business logic
│   │   ├── chat/       # Chat domain models
│   │   ├── common/     # Common domain models
│   │   ├── diagnosis/  # Diagnosis domain models
│   │   └── user/       # User domain models
│   │
│   ├── handler/        # HTTP request handlers
│   │   ├── future/     # Future-Job handlers
│   │   ├── rest/       # REST-API handlers
│   │   └── websocket/  # Websocket handlers
│   │
│   ├── dependency/     # Dependency injection
│   │
│   └── entrypoint/     # Application entry points
│       └── http/       # HTTP server entrypoint
│
├── migration/          # Database migration files
│   └── sql/            # Migration SQL scripts
│
├── pkg/                # Public library code
│   ├── database/       # Database utilities
│   │   ├── mysql/      # MySQL implementation
│   │   └── memdb/      # In-memory database
│   │
│   ├── future/         # Future utilities
│   │   └── database/   # Database-backed implementation
│   │
│   ├── http/           # HTTP utilities
│   │
│   ├── jwt/            # JWT authentication
│   │
│   ├── llm/            # LLM integration
│   │   └── gemini/     # Google Gemini integration
│   │
│   ├── oauth/          # OAuth integration
│   │   └── google/     # Google OAuth
│   │
│   ├── utils/          # Common utilities
│   │
│   └── websocket/      # WebSocket utilities
│
├── resource/           # Static resources
│   ├── diagnosis/      # Diagnosis resources
│   └── llm/            # LLM resources
│       ├── attachment/ # LLM attachments
│       └── prompt/     # LLM prompts
│
├── test/               # Test files
│   └── mock/           # Test mocks
│
├── .deploy/            # Deployment configurations
│
├── .github/            # GitHub related files
│   └── workflows/      # GitHub Actions workflows
│
├── main.go             # Main application entry
└── go.mod              # Go module definition
```

## 🚀 LAUNCH

```bash
# 1. Install dependencies
go mod download

# 2. Set up environment variables
vi .envrc

# 3. Apply environment variables
source .envrc

# 4. Start the server with local configurations
go run main.go http -n "local"
```

## ⏥ ARCHITECTURE
![SERVER-ARCHITECTURE](https://raw.githubusercontent.com/solutionchallenge/.github/refs/heads/main/assets/images/Ondaum-Server.png)

## 🏎️ PERFORMANCE
![SERVER-PERFORMANCE](https://raw.githubusercontent.com/solutionchallenge/.github/refs/heads/main/assets/images/Ondaum-Performance.png)
_Benchmarked on a GKE Managed Pod (180 mCPU / 256 MiB)_

## 📱 FEATURE
- AI Counseling With Um
- Psychological Assessments
  - International standard tests PHQ-9 / GAD-7 / PSS 
- AI Analysis of Conversation Content
  - Summary and organization of the conversation
  - Sharing feedback and areas for improvement
    
## ✨ VALUE
- Available for consultation anytime, anywhere
- Personalized consultation possible
- Reduced barriers to seeking counseling
- Access to a pre-trained professional psychological counseling AI