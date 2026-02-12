package services

import (
	"banking-system/config"
	"banking-system/models"
	"errors"
	"time"
)

type BankService struct{}

func NewBankService() *BankService {
	return &BankService{}
}

func (bs *BankService) CreateBank(name string) (*models.Bank, error) {
	bank := models.Bank{
		Name: name,
	}

	if result := config.GetDB().Create(&bank); result.Error != nil {
		return nil, result.Error
	}

	return &bank, nil
}

func (bs *BankService) GetBankByID(id uint) (*models.Bank, error) {
	var bank models.Bank
	result := config.GetDB().Preload("Branches").First(&bank, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bank, nil
}

type BranchService struct{}

func NewBranchService() *BranchService {
	return &BranchService{}
}

func (bs *BranchService) CreateBranch(bankID uint, name, address string) (*models.Branch, error) {
	var bank models.Bank
	if result := config.GetDB().First(&bank, bankID); result.Error != nil {
		return nil, errors.New("bank not found")
	}

	branch := models.Branch{
		BankID:  bankID,
		Name:    name,
		Address: address,
	}

	if result := config.GetDB().Create(&branch); result.Error != nil {
		return nil, result.Error
	}

	return &branch, nil
}

type CustomerService struct{}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (cs *CustomerService) RegisterCustomer(branchID uint, name, email, phone string) (*models.Customer, error) {
	var branch models.Branch
	if result := config.GetDB().First(&branch, branchID); result.Error != nil {
		return nil, errors.New("branch not found")
	}

	customer := models.Customer{
		BranchID: branchID,
		Name:     name,
		Email:    email,
		Phone:    phone,
	}

	if result := config.GetDB().Create(&customer); result.Error != nil {
		return nil, result.Error
	}

	return &customer, nil
}

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (as *AccountService) OpenSavingsAccount(customerID uint, holderRole string) (*models.SavingsAccount, error) {
	var customer models.Customer
	if result := config.GetDB().First(&customer, customerID); result.Error != nil {
		return nil, errors.New("customer not found")
	}
	if holderRole == "" {
		holderRole = "primary_holder"
	}
	tx := config.GetDB().Begin()
	account := models.SavingsAccount{
		Balance: 0,
	}
	if result := tx.Create(&account); result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	customerAccount := models.CustomerAccount{
		CustomerID: customerID,
		AccountID:  account.ID,
		HolderRole: holderRole,
	}
	if result := tx.Create(&customerAccount); result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &account, nil
}
func (as *AccountService) AddAccountHolder(accountID, customerID uint, holderRole string) (*models.CustomerAccount, error) {

	var account models.SavingsAccount
	if result := config.GetDB().First(&account, accountID); result.Error != nil {
		return nil, errors.New("account not found")
	}
	var customer models.Customer
	if result := config.GetDB().First(&customer, customerID); result.Error != nil {
		return nil, errors.New("customer not found")
	}
	var existingLink models.CustomerAccount
	if result := config.GetDB().Where("customer_id = ? AND account_id = ?", customerID, accountID).First(&existingLink); result.RowsAffected > 0 {
		return nil, errors.New("customer is already linked to this account")
	}
	customerAccount := models.CustomerAccount{
		CustomerID: customerID,
		AccountID:  accountID,
		HolderRole: holderRole,
	}
	if result := config.GetDB().Create(&customerAccount); result.Error != nil {
		return nil, result.Error
	}
	return &customerAccount, nil
}
func (as *AccountService) GetAccountBalance(accountID uint) (float64, error) {
	var account models.SavingsAccount
	if result := config.GetDB().First(&account, accountID); result.Error != nil {
		return 0, result.Error
	}
	return account.Balance, nil
}

func (as *AccountService) GetTransactionHistory(accountID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := config.GetDB().Where("account_id = ?", accountID).Order("created_at DESC").Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (as *AccountService) GetAccountHolders(accountID uint) ([]models.CustomerAccount, error) {
	var holders []models.CustomerAccount
	result := config.GetDB().Where("account_id = ?", accountID).Preload("Customer").Find(&holders)
	if result.Error != nil {
		return nil, result.Error
	}
	return holders, nil
}

type LoanService struct{}

func NewLoanService() *LoanService {
	return &LoanService{}
}

func (ls *LoanService) CreateLoan(customerID uint, loanType string, principalAmount float64) (*models.Loan, error) {

	var customer models.Customer
	if result := config.GetDB().First(&customer, customerID); result.Error != nil {
		return nil, errors.New("customer not found")
	}
	interestRate := 12.0
	totalPayableAmount := principalAmount + (principalAmount * interestRate / 100.0)

	loan := models.Loan{
		CustomerID:         customerID,
		LoanType:           loanType,
		PrincipalAmount:    principalAmount,
		InterestRate:       interestRate,
		TotalPayableAmount: totalPayableAmount,
		PendingAmount:      totalPayableAmount,
		StartDate:          time.Now(),
		Status:             "ACTIVE",
	}

	if result := config.GetDB().Create(&loan); result.Error != nil {
		return nil, result.Error
	}

	return &loan, nil
}

func (ls *LoanService) GetLoanByID(loanID uint) (*models.Loan, error) {
	var loan models.Loan
	result := config.GetDB().Preload("LoanPayments").First(&loan, loanID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &loan, nil
}

func (ls *LoanService) GetCustomerLoans(customerID uint) ([]models.Loan, error) {
	var loans []models.Loan
	result := config.GetDB().Where("customer_id = ?", customerID).Preload("LoanPayments").Find(&loans)
	if result.Error != nil {
		return nil, result.Error
	}
	return loans, nil
}

func (ls *LoanService) RepayLoan(loanID uint, amount float64) error {
	tx := config.GetDB().Begin()

	var loan models.Loan
	if result := tx.First(&loan, loanID); result.Error != nil {
		tx.Rollback()
		return errors.New("loan not found")
	}

	if loan.Status == "CLOSED" {
		tx.Rollback()
		return errors.New("loan is already closed")
	}

	if amount > loan.PendingAmount {
		tx.Rollback()
		return errors.New("repayment amount exceeds pending amount")
	}

	loan.PendingAmount -= amount
	if loan.PendingAmount == 0 {
		loan.Status = "CLOSED"
	}

	if result := tx.Save(&loan); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	payment := models.LoanPayment{
		LoanID:      loanID,
		Amount:      amount,
		PaymentDate: time.Now(),
	}

	if result := tx.Create(&payment); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (ls *LoanService) CalculateYearlyInterest(loanID uint) (float64, error) {
	var loan models.Loan
	if result := config.GetDB().First(&loan, loanID); result.Error != nil {
		return 0, result.Error
	}
	interest := (loan.PendingAmount * loan.InterestRate) / 100.0
	return interest, nil
}

func (ls *LoanService) GetLoanPaymentHistory(loanID uint) ([]models.LoanPayment, error) {
	var payments []models.LoanPayment
	result := config.GetDB().Where("loan_id = ?", loanID).Order("payment_date DESC").Find(&payments)
	if result.Error != nil {
		return nil, result.Error
	}
	return payments, nil
}
