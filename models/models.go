package models

import "time"

type Bank struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Branches  []Branch  `gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE" json:"branches,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Branch struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	BankID    uint       `gorm:"not null;index" json:"bank_id"`
	Bank      Bank       `gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE" json:"bank,omitempty"`
	Name      string     `gorm:"not null" json:"name"`
	Address   string     `json:"address"`
	Customers []Customer `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE" json:"customers,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Customer struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	BranchID         uint              `gorm:"not null;index" json:"branch_id"`
	Branch           Branch            `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE" json:"branch,omitempty"`
	Name             string            `gorm:"not null" json:"name"`
	Email            string            `gorm:"uniqueIndex" json:"email"`
	Phone            string            `json:"phone"`
	CustomerAccounts []CustomerAccount `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"customer_accounts,omitempty"`
	Loans            []Loan            `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"loans,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type SavingsAccount struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	Balance          float64           `gorm:"not null;default:0" json:"balance"`
	CustomerAccounts []CustomerAccount `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"customer_accounts,omitempty"`
	Transactions     []Transaction     `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"transactions,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type CustomerAccount struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CustomerID uint           `gorm:"not null;index" json:"customer_id"`
	Customer   Customer       `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
	AccountID  uint           `gorm:"not null;index" json:"account_id"`
	Account    SavingsAccount `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
	HolderRole string         `gorm:"not null;default:'primary_holder'" json:"holder_role"` // primary_holder or joint_holder
	CreatedAt  time.Time      `json:"created_at"`
}

type Transaction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	AccountID uint           `gorm:"not null;index" json:"account_id"`
	Account   SavingsAccount `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
	Type      string         `gorm:"not null" json:"type"` // DEPOSIT or WITHDRAW
	Amount    float64        `gorm:"not null" json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
}

type Loan struct {
	ID                 uint          `gorm:"primaryKey" json:"id"`
	CustomerID         uint          `gorm:"not null;index" json:"customer_id"`
	Customer           Customer      `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
	LoanType           string        `gorm:"not null" json:"loan_type"`
	PrincipalAmount    float64       `gorm:"not null" json:"principal_amount"`
	InterestRate       float64       `gorm:"not null;default:12" json:"interest_rate"` // 12% fixed
	TotalPayableAmount float64       `gorm:"not null" json:"total_payable_amount"`
	PendingAmount      float64       `gorm:"not null" json:"pending_amount"`
	StartDate          time.Time     `gorm:"not null" json:"start_date"`
	EndDate            *time.Time    `json:"end_date,omitempty"`
	Status             string        `gorm:"not null;default:'ACTIVE'" json:"status"` // ACTIVE or CLOSED
	LoanPayments       []LoanPayment `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE" json:"loan_payments,omitempty"`
	CreatedAt          time.Time     `json:"created_at"`
}

type LoanPayment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LoanID      uint      `gorm:"not null;index" json:"loan_id"`
	Loan        Loan      `gorm:"foreignKey:LoanID;constraint:OnDelete:CASCADE" json:"loan,omitempty"`
	Amount      float64   `gorm:"not null" json:"amount"`
	PaymentDate time.Time `gorm:"not null" json:"payment_date"`
	CreatedAt   time.Time `json:"created_at"`
}
