Banking System – API List

- Health
GET /health – Check API status

- Banks
POST /banks – Create bank
GET /banks – Get all banks
GET /banks/:id – Get bank by ID

- Branches
POST /branches – Create branch
GET /branches/:id – Get branch by ID
GET /banks/:bank_id/branches – Get branches of a bank

- Customers
POST /customers – Register customer
GET /customers/:id – Get customer details
GET /branches/:branch_id/customers – Get customers by branch
GET /customers/:customer_id/loans – Get customer loans

- Savings Accounts
POST /accounts/savings – Open savings account
POST /accounts/:account_id/holders – Add joint holder
GET /accounts/:id – Get account details
POST /accounts/:id/deposit – Deposit money
POST /accounts/:id/withdraw – Withdraw money
GET /accounts/:id/transactions – Get transaction history

- Loans
POST /loans – Take loan
GET /loans/:id – Get loan details
POST /loans/:id/repay – Repay loan
GET /loans/:id/interest – Calculate loan interest
GET /loans/:id/payments – Get loan payment history



