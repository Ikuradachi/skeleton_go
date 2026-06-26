# <center> 🚀 Go API Skeleton </center>

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![Echo](https://img.shields.io/badge/Echo-v4-00ADD8?style=for-the-badge)
![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)

**Template Modern untuk REST API dengan Golang**

Skeleton project yang powerful dan scalable untuk membangun REST API dengan Go. Dilengkapi dengan best practices, clean architecture, dan tools modern untuk development yang cepat.

</div>

---

## 📋 Deskripsi

Go API Skeleton adalah template production-ready untuk membangun REST API menggunakan Golang. Project ini mengimplementasikan clean architecture dengan separation of concerns yang jelas, memudahkan development, testing, dan maintenance aplikasi Anda.

## ✨ Fitur Unggulan

### 🏗️ Clean Architecture
Struktur project yang terorganisir dengan separation of concerns yang jelas. Handler, Service, dan Repository layer yang terpisah untuk maintainability maksimal.

### ⚡ Echo Framework
Menggunakan Echo, framework Go yang high-performance dan minimalist. Dilengkapi dengan middleware untuk CORS, logging, dan recovery.

### 🗄️ Database Ready
Support untuk PostgreSQL dan MySQL dengan sqlx. Konfigurasi database yang mudah dan connection pooling yang optimal.

### 📝 YAML Configuration
Konfigurasi aplikasi menggunakan Viper dengan format YAML. Mudah dikelola dan mendukung multiple environments.

### 📦 Modular Structure
Feature-based folder structure yang memudahkan scaling. Setiap feature memiliki entity, handler, service, dan repository sendiri.

### 🛡️ Middleware Ready
Pre-configured middleware untuk CORS, logging, dan panic recovery. Siap untuk ditambahkan authentication dan authorization.

## 🏛️ Arsitektur

Project ini mengikuti **Clean Architecture** dengan struktur sebagai berikut:

```
skeleton/
├── cmd/api/              # Entry point aplikasi
│   └── main.go          # Main application
├── internal/            # Private application code
│   ├── api/            # Router & middleware
│   │   ├── router.go
│   │   └── middleware.go
│   ├── features/       # Feature modules
│   │   └── auth/       # Contoh: Auth feature
│   │       ├── entity.go
│   │       ├── handler.go
│   │       ├── service.go
│   │       └── repository.go
│   └── shared/         # Shared utilities
│       ├── request.go
│       └── response.go
├── pkg/                # Public reusable packages
│   ├── config/        # Configuration management
│   ├── logger/        # Logging utilities
│   ├── migration/     # Database migrations
│   └── utils/         # Helper functions
├── app/               # Static files & frontend
├── docs/              # Documentation
├── test/              # Test files
├── Dockerfile         # Docker configuration
├── docker-compose.yaml # Docker compose configuration
├── go.mod
└── go.sum
```

### Separation of Concerns

Setiap layer memiliki tanggung jawab yang jelas:

- **Handler:** Menangani HTTP request/response
- **Service:** Business logic dan validasi
- **Repository:** Database operations
- **Entity:** Data models dan structures

## 🛠️ Technology Stack

| Technology | Deskripsi |
|-----------|-----------|
| **Go** | Modern Go version |
| **Echo** | High performance web framework |
| **sqlx** | Database toolkit |
| **PostgreSQL** | Primary database (optional) |
| **MySQL** | Alternative database (optional) |
| **Viper** | Configuration management |
| **Docker** | Containerization |

## 🐳 Docker Deployment

Aplikasi dapat dijalankan secara containerized menggunakan Docker dan Docker Compose.

### Docker Compose (Direkomendasikan)
Jalankan perintah berikut untuk membangun dan menjalankan service:
```bash
docker compose up -d --build
```

### Docker CLI
Jika ingin membangun dan menjalankan container secara manual:

1. Build docker image:
   ```bash
   docker build -t skeleton-go .
   ```
2. Jalankan container dengan melakukan mount konfigurasi:
   ```bash
   docker run -d \
     -p 38600:38600 \
     --name skeleton-go \
     -v $(pwd)/pkg/config/config.yaml:/app/pkg/config/config.yaml \
     skeleton-go
   ```

<div align="center">

**Built with ❤️ using Go**

🌐 [Live Demo](https://ikuradachi.github.io/skeleton_go/) | 📦 [GitHub Repository](https://github.com/Ikuradachi/skeleton_go)

⭐ Star this repository if you find it helpful!

</div>

