#  Banking System Backend

A scalable banking system backend built using **Go**, **Gin**, **GORM**, and **PostgreSQL**.

This project simulates real-world banking operations including account management, transactions, and loan processing with guaranteed data consistency using database transactions.

---

##  Key Features

- Multi-bank & multi-branch architecture
- Customer onboarding & profile management
- Savings accounts (deposit & withdrawal)
- Immutable transaction history
- Loan system with fixed **12% annual interest**
- Total payable amount calculated at loan creation
- Automatic loan closure after full repayment
- RESTful API design with proper HTTP status codes
- Atomic operations using database transactions

---

## 🛠 Tech Stack

- **Go (Golang)**
- **Gin** – HTTP Web Framework
- **GORM** – ORM Library
- **PostgreSQL** – Relational Database

---

## 📂 Project Structure

banking-system/
├── cmd/main.go
├── config/
├── controllers/
├── models/
├── routes/
├── services/
└── schema.sql

---

## ⚡ Getting Started

### Install Dependencies
```bash
go mod download
Run the Application
go run cmd/main.go
Server runs at:
http://localhost:8080
🔌 API Overview
Method	Endpoint	Description
POST	/banks	Create bank
POST	/customers	Register customer
POST	/accounts/savings	Open savings account
POST	/accounts/{id}/deposit	Deposit money
POST	/accounts/{id}/withdraw	Withdraw money
POST	/loans	Create loan
POST	/loans/{id}/repay	Repay loan
