![Approval Workflow Dashboard](screenshots/dashboard.png)

# Approval Workflow System

A full-stack approval workflow system built with **Go**, **PostgreSQL**, **Docker**, and a **React-based admin panel**.

This project demonstrates how environment-based changes (e.g. feature flags, configuration updates) can be requested, reviewed, approved or rejected through a clean API and a modern web interface.

---

## 🚀 Features

### Backend (Go)
- RESTful API built with `chi`
- PostgreSQL database with migrations
- Transaction-safe approval & rejection logic
- Approval state machine:
  - `PENDING`
  - `APPROVED`
  - `REJECTED`
- Action logging (`CREATED`, `APPROVED`, `REJECTED`)
- Clear separation of layers:
  - handlers
  - store / repository
  - domain types
- Dockerized PostgreSQL setup

### Frontend (React)
- Modern admin panel built with React (Vite)
- List all approval requests
- View request details and change payload
- Create new requests
- Approve / Reject pending requests
- View request logs (audit trail)
- Clean and dark-themed UI

---

## 🧱 Tech Stack

**Backend**
- Go
- chi router
- PostgreSQL
- Docker & Docker Compose

**Frontend**
- React (Vite)
- JavaScript
- Fetch API

---

## 📦 Project Structure

approval-workflow-system/
├── cmd/api # Application entry point
├── internal
│ ├── approval # Business logic & store
│ ├── server # HTTP handlers & routes
│ ├── db # Database connection
│ └── types # Domain models
├── migrations # SQL migrations
├── scripts # Migration runner
├── frontend # React admin panel
├── docker-compose.yml
└── README.md

yaml
Kodu kopyala

---

## ⚙️ How to Run Locally

### 1️⃣ Start PostgreSQL with Docker

```bash
docker compose up -d
Wait until the database is ready.

2️⃣ Run database migrations
bash
Kodu kopyala
go run ./scripts/migrate.go
3️⃣ Start the backend API
bash
Kodu kopyala
go run ./cmd/api
API will be available at:

arduino
Kodu kopyala
http://localhost:8080
4️⃣ Start the frontend panel
bash
Kodu kopyala
cd frontend
npm install
npm run dev
Frontend will be available at:

arduino
Kodu kopyala
http://localhost:5173
🔌 API Endpoints (Summary)
Method	Endpoint	Description
GET	/requests	List all requests
POST	/requests	Create a new request
GET	/requests/{id}	Get request details
POST	/requests/{id}/approve	Approve a request
POST	/requests/{id}/reject	Reject a request
GET	/requests/{id}/logs	Get request logs

🧪 Example Workflow
Create a request from dev → staging

Request status starts as PENDING

Reviewer approves or rejects the request

Status changes to APPROVED or REJECTED

All actions are logged and visible in the UI

🎯 Why This Project?
This project was built to:

Practice real-world backend architecture in Go

Implement transactional business rules

Design a clean approval workflow

Build a usable admin panel without relying on Postman

Simulate production-like systems used in modern companies

📌 Future Improvements
Authentication & role-based access

Deployment (Render / Railway / Fly.io)

Request filtering & search

Diff viewer for change payloads

Real-time updates (SSE / WebSocket)