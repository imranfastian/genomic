# 🧬 Genomic Secure API — Golang Microservice

A secure, containerized microservice built in Go to demonstrate core backend skills:

- RESTful APIs
- JWT Authentication
- PostgreSQL integration
- Docker & Kubernetes-ready
- Follows DevOps best practices

## 🚀 Features

✅ Secure login with JWT  
✅ Protected genomic data access  
✅ PostgreSQL for structured data  
✅ Docker-based dev workflow  
✅ Kubernetes manifests for deployment  
✅ Follows security-first and cloud-native patterns

---

## 🧱 API Endpoints

| Method | Endpoint                        | Description                            |
| ------ | ------------------------------- | -------------------------------------- |
| POST   | `http://localhost:8080/login`   | Returns JWT token                      |
| GET    | `http://localhost:8080/genomes` | Returns genome records (JWT-protected) |

---

## 🐳 Local Development (Docker Compose)

```bash
cp .env.example .env
docker compose --env-file .env up --build
```
