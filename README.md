
# Approval Workflow System

This project is a Git-native approval workflow service designed to manage feature flag changes between environments 
(e.g. Dev → Staging → Production).

## Tech Stack

- Backend: Go
- Database: PostgreSQL
- Frontend: React (TypeScript)
- Containerization: Docker, docker-compose

## Goals (v1)

- User authentication (JWT)
- Role-based access control (RBAC)
- Approval requests between environments
- Approval & reject endpoints with audit logs
- Simple React dashboard to list and manage requests
