# Banking System - Complete Implementation

A comprehensive banking system built with **Go**, **Gin**, **GORM**, and **PostgreSQL**.

## 🏦 System Overview

This banking system provides a complete solution for managing:
- Multiple banks and branches
- Customer accounts and profiles
- Savings accounts with deposit/withdrawal operations
- Loans with 12% fixed interest rate
- Transaction history and loan payment tracking
- Real-time interest calculations

---

## ✨ Key Features

### 1. Bank & Branch Management
- Create and manage multiple banks
- Create multiple branches under each bank
- View bank and branch details

### 2. Customer Management
- Register customers
- Link customers to specific branches
- View complete customer profiles

### 3. Savings Accounts
- Open savings accounts (one per customer)
- Deposit and withdraw money
- View account balance
- Track transaction history

### 4. Loan Management
- Take loans with 12% fixed interest rate
- View loan details (principal, pending amount, interest)
- Repay loans (partial or full)
- Automatic loan closure on full repayment

### 5. Interest Calculation
- Real-time interest calculation: `(Pending Amount × 12) / 100`
- View yearly interest based on pending amount

### 6. Transaction Safety
- Database transactions for atomic operations
- Balance validation before withdrawal
- Loan status tracking

---

## 📁 Project Structure

```
banking-system/
├── main.go                    # Server entry point
├── go.mod                     # Go module definition
│
├── config/
│   └── db.go                 # Database initialization & config
│
├── models/
│   └── models.go             # All database models:
│                            #   - Bank
│                            #   - Branch
│                            #   - Customer
│                            #   - SavingsAccount
│                            #   - Transaction
│                            #   - Loan
│                            #   - LoanPayment
│
├── controllers/
│   ├── bank_controller.go    # Bank endpoints
│   ├── branch_controller.go  # Branch endpoints
│   ├── customer_controller.go # Customer endpoints
│   ├── account_controller.go # Account endpoints
│   └── loan_controller.go    # Loan endpoints
│
├── services/
│   └── banking_service.go    # Business logic layer
│
├── routes/
│   └── routes.go             # API route definitions
│
└── Documentation/
    ├── WORKFLOW.md           # Complete workflow guide
    ├── API_DOCUMENTATION.md  # API reference
    ├── FEATURES.md           # Features list
    ├── DATABASE_SCHEMA.md    # DB schema & ER diagram
    ├── SETUP.md              # Installation guide
    └── README.md             # This file
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 12+

### 1. Setup Database

```bash
# Create database
createdb banking_system

# Or using psql
psql -U postgres
CREATE DATABASE banking_system;
```

### 2. Environment Variables

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=banking_system
export DB_SSLMODE=disable
export PORT=8080
```

### 3. Install & Run

```bash
cd "/Users/bipl/Assignment 2"
go mod tidy
go run main.go
```

Expected output:
```
Database connection established successfully
Database migrations completed successfully
Starting Banking System API on port 8080
```

### 4. Test API

```bash
# Health check
curl http://localhost:8080/health

# Create bank
curl -X POST http://localhost:8080/banks \
  -H "Content-Type: application/json" \
  -d '{"name":"HDFC Bank"}'
```

---

## 📚 API Endpoints

### Banks
- `POST /banks` - Create bank
- `GET /banks/:id` - Get bank details
- `GET /banks` - Get all banks

### Branches
- `POST /branches` - Create branch
- `GET /branches/:id` - Get branch details
- `GET /banks/:bank_id/branches` - Get branches by bank

### Customers
- `POST /customers` - Register customer
- `GET /customers/:id` - Get customer details
- `GET /branches/:branch_id/customers` - Get customers by branch

### Accounts
- `POST /accounts/savings` - Open savings account
- `GET /accounts/:id` - Get account details
- `POST /accounts/:id/deposit` - Deposit money
- `POST /accounts/:id/withdraw` - Withdraw money
- `GET /accounts/:id/transactions` - Get transaction history

### Loans
- `POST /loans` - Take loan
- `GET /loans/:id` - Get loan details
- `GET /customers/:customer_id/loans` - Get customer loans
- `POST /loans/:id/repay` - Repay loan
- `GET /loans/:id/interest` - Calculate loan interest
- `GET /loans/:id/payments` - Get loan payment history

---

## 💡 Usage Examples

### Complete Workflow Example

```bash
# 1. Create Bank
BANK=$(curl -s -X POST http://localhost:8080/banks \
  -H "Content-Type: application/json" \
  -d '{"name":"State Bank"}')
BANK_ID=$(echo $BANK | jq .id)

# 2. Create Branch
BRANCH=$(curl -s -X POST http://localhost:8080/branches \
  -H "Content-Type: application/json" \
  -d "{\"bank_id\":$BANK_ID,\"name\":\"Mumbai\",\"address\":\"Fort\"}")
BRANCH_ID=$(echo $BRANCH | jq .id)

# 3. Register Customer
CUSTOMER=$(curl -s -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d "{\"branch_id\":$BRANCH_ID,\"name\":\"Rajesh\",\"email\":\"rajesh@bank.com\",\"phone\":\"9876543210\"}")
CUSTOMER_ID=$(echo $CUSTOMER | jq .id)

# 4. Open Account
ACCOUNT=$(curl -s -X POST http://localhost:8080/accounts/savings \
  -H "Content-Type: application/json" \
  -d "{\"customer_id\":$CUSTOMER_ID}")
ACCOUNT_ID=$(echo $ACCOUNT | jq .id)

# 5. Deposit Money
curl -s -X POST http://localhost:8080/accounts/$ACCOUNT_ID/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount":50000}' | jq .

# 6. Take Loan
LOAN=$(curl -s -X POST http://localhost:8080/loans \
  -H "Content-Type: application/json" \
  -d "{\"customer_id\":$CUSTOMER_ID,\"principal_amount\":100000}")
LOAN_ID=$(echo $LOAN | jq .id)

# 7. Check Interest
curl -s http://localhost:8080/loans/$LOAN_ID/interest | jq .

# 8. Repay Loan
curl -s -X POST http://localhost:8080/loans/$LOAN_ID/repay \
  -H "Content-Type: application/json" \
  -d '{"amount":20000}' | jq .

# 9. View Account
curl -s http://localhost:8080/accounts/$ACCOUNT_ID | jq .
```

---

## 🔒 Data Integrity

All critical operations use database transactions:

```
Deposit:
  1. Validate account exists
  2. Update balance
  3. Create transaction record
  4. Commit atomically

Withdraw:
  1. Validate balance >= amount
  2. Deduct from balance
  3. Create transaction record
  4. Commit atomically

Loan Repayment:
  1. Validate pending amount >= amount
  2. Deduct from pending
  3. Create payment record
  4. Close loan if pending = 0
  5. Commit atomically
```

---

## 📊 Interest Calculation

```
Formula: Interest = (Pending Amount × 12) / 100

Example:
- Principal: ₹100,000
- Interest Rate: 12% p.a.
- After ₹20,000 repayment:
  Pending: ₹80,000
  Interest: (80,000 × 12) / 100 = ₹9,600
```

---

## 🛠 Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.21+ |
| **Web Framework** | Gin Gonic |
| **ORM** | GORM |
| **Database** | PostgreSQL |
| **Server Port** | 8080 |

---

## 📖 Documentation

- **[WORKFLOW.md](WORKFLOW.md)** - Complete system workflow and architecture
- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Detailed API reference with examples
- **[FEATURES.md](FEATURES.md)** - Complete features list
- **[DATABASE_SCHEMA.md](DATABASE_SCHEMA.md)** - ER diagram and SQL schema
- **[SETUP.md](SETUP.md)** - Installation and setup guide

---

## 🧪 Testing

### Manual Testing with cURL

All API endpoints can be tested using cURL commands provided in [API_DOCUMENTATION.md](API_DOCUMENTATION.md).

### Automated Testing

Create a bash script `test_banking.sh`:

```bash
#!/bin/bash

API="http://localhost:8080"

# Test health
echo "Testing health..."
curl -s $API/health | jq .

# Create bank
echo -e "\nCreating bank..."
curl -s -X POST $API/banks \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Bank"}' | jq .

# Continue with other tests...
```

---

## 🔧 Configuration

Environment variables:

```
DB_HOST=localhost          # PostgreSQL host
DB_PORT=5432              # PostgreSQL port
DB_USER=postgres          # Database user
DB_PASSWORD=postgres      # Database password
DB_NAME=banking_system    # Database name
DB_SSLMODE=disable        # SSL mode
PORT=8080                 # API server port
```

---

## 📝 Error Handling

The API returns proper HTTP status codes:

- `200 OK` - Request successful
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error response format:
```json
{
    "error": "Descriptive error message"
}
```

---

## 🚦 Running the Server

### Development Mode

```bash
go run main.go
```

### Production Build

```bash
go build -o banking-system

# Run binary
./banking-system
```

### With Hot Reload (Optional)

Install `air`:
```bash
go install github.com/cosmtrek/air@latest
air
```

---

## 📋 Database Verification

Check tables after migrations:

```sql
psql -U postgres -d banking_system

\dt              -- List tables
\d banks         -- Describe table structure
SELECT * FROM banks;
SELECT COUNT(*) FROM customers;
```

---

## 🎯 Use Cases

### Scenario 1: Regular Customer Operations
1. Customer opens account
2. Deposits monthly salary
3. Withdraws for expenses
4. Views transaction history

### Scenario 2: Loan Management
1. Customer takes ₹1,00,000 loan
2. Views yearly interest: ₹12,000
3. Makes partial repayments
4. Loan automatically closes on full repayment

### Scenario 3: Multi-Branch Bank
1. Main HQ creates multiple branches
2. Each branch registers customers
3. Central reporting across branches

---

## 🔐 Security Features

- ✅ Input validation on all endpoints
- ✅ SQL injection prevention (GORM)
- ✅ Balance validation for withdrawals
- ✅ Loan status validation
- ✅ Database constraint enforcement

**Future Enhancements**:
- JWT authentication
- Role-based access control
- Rate limiting
- API key management

---

## 📈 Performance

### Optimizations Included
- Database indexes on foreign keys
- Indexes on frequently searched fields
- Connection pooling via GORM
- Atomic transactions to prevent race conditions

### Scalability
- Designed for horizontal scaling
- Stateless API design
- Database-driven consistency

---

## 🐛 Troubleshooting

### Database Connection Error
```bash
# Verify PostgreSQL is running
brew services list

# Check database exists
createdb banking_system
```

### Port Already in Use
```bash
# Find and kill process
lsof -ti:8080 | xargs kill -9

# Or use different port
PORT=8081 go run main.go
```

### Go Module Issues
```bash
go mod tidy
go mod download
```

---

## 📦 Deployment

### Docker (Optional)

Create `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o banking-system
EXPOSE 8080
CMD ["./banking-system"]
```

Build and run:
```bash
docker build -t banking-system .
docker run -p 8080:8080 -e DB_HOST=host.docker.internal banking-system
```

---

## 🤝 Contributing

This is an educational project. Feel free to:
- Extend features
- Improve documentation
- Add authentication
- Implement advanced features

---

## 📄 License

This project is for educational purposes.

---

## 📞 Support

For documentation and examples, refer to:
- [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - API reference
- [WORKFLOW.md](WORKFLOW.md) - System workflow
- [SETUP.md](SETUP.md) - Installation steps

---

## ✅ Checklist

- ✅ Database models created
- ✅ All controllers implemented
- ✅ Service layer built
- ✅ Routes configured
- ✅ Database migrations working
- ✅ Transaction safety ensured
- ✅ Interest calculation implemented
- ✅ Error handling in place
- ✅ Documentation complete
- ✅ API endpoints tested

---

**Happy Banking! 🏦**

