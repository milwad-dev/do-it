# Do‑It 📝

Do-It is a to-do application built with Golang and a REST API. It offers a fast and scalable solution for task management, enabling users to easily add, update, and delete tasks. The app uses Golang’s robust performance alongside a well-structured REST API to provide a seamless and efficient experience for managing daily tasks.

### 🔥 Features

- ✅ Full CRUD for tasks
- 📝 Swagger (OpenAPI 3) UI for API testing
- ⚙️ Dockerized environment
- 🧪 Testable structure with modular design
- ♻️ Clean Restful design
- 🌐 `.env` support for easy configuration
- 🪶 Lightweight and dependency minimal

### 🚀 Getting Started

Prerequisites
- Go ≥1.22 (module support)
- Docker & Docker Compose (optional, for containerized runs)
- Make sure to have a `.env` file—see `.env.example` for reference.

#### Run Locally

```shell
git clone https://github.com/milwad-dev/do-it.git
cd do-it
cp .env.example .env           # configure DB and API settings
go mod download                # installs dependencies
go run main.go                 # starts server on :8000 by default
```

#### Run with Docker Compose

```shell
docker-compose up --build
./wait-for-db.sh              # waits for DB container to be ready
```

### 📐 API Documentation

Access Swagger UI at:

```shell
http://localhost:8000/api/swagger/
```

Inspect interactive docs and Swagger JSON at:

```shell
http://localhost:8000/api/swagger/doc.json
```

### 🧩 Project Structure

```perl
├── internal/           # core business logic and handlers
├── http-requests/      # HTTP client utilities
├── docs/               # auto-generated Swagger docs (via swaggo)
├── Dockerfile
├── docker-compose.yml
├── wait-for-db.sh      # DB readiness helper
└── main.go             # app entry point
```

### ⚙️ Environment Variables

Customize using .env file:

```ini
APP_PORT=8000

DB_HOST=...
DB_PORT=...
DB_USER=...
DB_PASS=...
DB_NAME=...
```

### 🛠 Contributing

Contributions are welcome—follow these steps:

1. Fork the repo
2. Create a new branch: `git checkout -b feature/YourFeature`
3. Make changes & add tests
4. Submit a Pull Request

### ✅ License

This project is open‑source under the MIT License. See the full text in [LICENSE](https://github.com/milwad-dev/do-it/blob/master/LICENSE). 

### 💡 Related Projects by Milwad Khosravi

- [ToWork](https://github.com/milwad-dev/towork-backend) (Laravel + REST API) – ToDo app backend in Laravel 
- [Go Shop](https://github.com/milwad-dev/go-shop) – Another REST‑based to‑do application in Go 

### 👨‍💻 Author

Created with ❤️ by [Milwad Khosravi](https://github.com/milwad-dev)

### 🔭 What's Next

1. Add Grafana
2. Add Elasticsearch
3. Use Goroutines
4. Use Kafka
5. Add scheduler for reminder tasks
