package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

// Colors for terminal output
const (
	colorGreen  = "\033[92m"
	colorRed    = "\033[91m"
	colorYellow = "\033[93m"
	colorBlue   = "\033[94m"
	colorReset  = "\033[0m"
)

// API results tracker
var apiResults = struct {
	passed []string
	failed []string
}{
	passed: []string{},
	failed: []string{},
}

// Response structures
type BankResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BranchResponse struct {
	ID      uint   `json:"id"`
	BankID  uint   `json:"bank_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type CustomerResponse struct {
	ID       uint   `json:"id"`
	BranchID uint   `json:"branch_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type AccountResponse struct {
	ID      uint    `json:"id"`
	Balance float64 `json:"balance"`
}

type LoanResponse struct {
	ID                 uint    `json:"id"`
	CustomerID         uint    `json:"customer_id"`
	LoanType           string  `json:"loan_type"`
	PrincipalAmount    float64 `json:"principal_amount"`
	TotalPayableAmount float64 `json:"total_payable_amount"`
	PendingAmount      float64 `json:"pending_amount"`
	Status             string  `json:"status"`
}

type AccountDetail struct {
	AccountNo  uint    `json:"account_no"`
	BranchName string  `json:"branch_name"`
	BankName   string  `json:"bank_name"`
	Balance    float64 `json:"balance"`
}

type CustomerWithAccounts struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Accounts []AccountDetail `json:"accounts"`
}

// Logging functions
func logSuccess(apiName, message string) {
	fmt.Printf("%s✓ %s - SUCCESS%s %s\n", colorGreen, apiName, colorReset, message)
	apiResults.passed = append(apiResults.passed, apiName)
}

func logFailure(apiName, errorMsg string) {
	fmt.Printf("%s✗ %s - FAILED%s %s\n", colorRed, apiName, colorReset, errorMsg)
	apiResults.failed = append(apiResults.failed, apiName)
}

func logInfo(message string) {
	fmt.Printf("%sℹ %s%s\n", colorBlue, message, colorReset)
}

func logWarning(message string) {
	fmt.Printf("%s⚠ %s%s\n", colorYellow, message, colorReset)
}

func printSeparator() {
	fmt.Printf("%s%s%s\n", colorYellow, "================================================================================", colorReset)
}

func logRequest(method, endpoint string, payload interface{}) {
	fmt.Printf("\n%s→ REQUEST:%s %s %s\n", colorBlue, colorReset, method, endpoint)
	if payload != nil {
		jsonData, _ := json.MarshalIndent(payload, "  ", "  ")
		fmt.Printf("  %sPayload:%s\n  %s\n", colorYellow, colorReset, string(jsonData))
	}
}

func logResponse(statusCode int, response interface{}) {
	color := colorGreen
	if statusCode >= 400 {
		color = colorRed
	}
	fmt.Printf("%s← RESPONSE:%s Status: %d\n", color, colorReset, statusCode)
	if response != nil {
		jsonData, _ := json.MarshalIndent(response, "  ", "  ")
		fmt.Printf("  %sBody:%s\n  %s\n", colorYellow, colorReset, string(jsonData))
	}
}

// HTTP helper function
func makeRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// Test functions
func testCreateBank(name string) *uint {
	payload := map[string]string{"name": name}
	logRequest("POST", "/banks", payload)

	resp, err := makeRequest("POST", baseURL+"/banks", payload)
	if err != nil {
		logFailure("POST /banks", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data BankResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 201 {
		logSuccess("POST /banks", fmt.Sprintf("Created bank: %s (ID: %d)", data.Name, data.ID))
		return &data.ID
	}

	logFailure("POST /banks", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testGetBank(bankID uint) *BankResponse {
	endpoint := fmt.Sprintf("/banks/%d", bankID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /banks/:id", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data BankResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /banks/:id", fmt.Sprintf("Retrieved bank: %s", data.Name))
		return &data
	}

	logFailure("GET /banks/:id", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testCreateBranch(bankID uint, name, address string) *uint {
	payload := map[string]interface{}{
		"bank_id": bankID,
		"name":    name,
		"address": address,
	}
	logRequest("POST", "/branches", payload)

	resp, err := makeRequest("POST", baseURL+"/branches", payload)
	if err != nil {
		logFailure("POST /branches", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data BranchResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 201 {
		logSuccess("POST /branches", fmt.Sprintf("Created branch: %s (ID: %d)", data.Name, data.ID))
		return &data.ID
	}

	logFailure("POST /branches", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testGetBranch(branchID uint) *BranchResponse {
	endpoint := fmt.Sprintf("/branches/%d", branchID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /branches/:id", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data BranchResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /branches/:id", fmt.Sprintf("Retrieved branch: %s", data.Name))
		return &data
	}

	logFailure("GET /branches/:id", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testCreateCustomer(branchID uint, name, email, phone string) *uint {
	payload := map[string]interface{}{
		"branch_id": branchID,
		"name":      name,
		"email":     email,
		"phone":     phone,
	}
	logRequest("POST", "/customers", payload)

	resp, err := makeRequest("POST", baseURL+"/customers", payload)
	if err != nil {
		logFailure("POST /customers", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data CustomerResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 201 {
		logSuccess("POST /customers", fmt.Sprintf("Created customer: %s (ID: %d)", data.Name, data.ID))
		return &data.ID
	}

	logFailure("POST /customers", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testGetCustomer(customerID uint) *CustomerResponse {
	endpoint := fmt.Sprintf("/customers/%d", customerID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /customers/:id", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data CustomerResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /customers/:id", fmt.Sprintf("Retrieved customer: %s", data.Name))
		return &data
	}

	logFailure("GET /customers/:id", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testListCustomersWithAccounts() []CustomerWithAccounts {
	logRequest("GET", "/customers", nil)

	resp, err := makeRequest("GET", baseURL+"/customers", nil)
	if err != nil {
		logFailure("GET /customers", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data []CustomerWithAccounts
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /customers", fmt.Sprintf("Retrieved %d customers with accounts", len(data)))
		return data
	}

	logFailure("GET /customers", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testCreateAccount(customerID uint, holderRole string) *uint {
	payload := map[string]interface{}{
		"customer_id": customerID,
		"holder_role": holderRole,
	}
	logRequest("POST", "/accounts/savings", payload)

	resp, err := makeRequest("POST", baseURL+"/accounts/savings", payload)
	if err != nil {
		logFailure("POST /accounts/savings", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	logResponse(resp.StatusCode, result)

	if resp.StatusCode == 201 {
		// Extract account from nested structure
		if account, ok := result["account"].(map[string]interface{}); ok {
			if id, ok := account["id"].(float64); ok {
				accountID := uint(id)
				balance := 0.0
				if bal, ok := account["balance"].(float64); ok {
					balance = bal
				}
				logSuccess("POST /accounts/savings", fmt.Sprintf("Created account (ID: %d) with balance: %.2f", accountID, balance))
				return &accountID
			}
		}
		logFailure("POST /accounts/savings", "Invalid response structure")
		return nil
	}

	logFailure("POST /accounts/savings", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testGetAccount(accountID uint) *AccountResponse {
	endpoint := fmt.Sprintf("/accounts/%d", accountID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /accounts/:id", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data AccountResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /accounts/:id", fmt.Sprintf("Retrieved account with balance: %.2f", data.Balance))
		return &data
	}

	logFailure("GET /accounts/:id", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testCreateLoan(customerID uint, loanType string, principalAmount float64) *uint {
	payload := map[string]interface{}{
		"customer_id":      customerID,
		"loan_type":        loanType,
		"principal_amount": principalAmount,
	}
	logRequest("POST", "/loans", payload)

	resp, err := makeRequest("POST", baseURL+"/loans", payload)
	if err != nil {
		logFailure("POST /loans", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data LoanResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 201 {
		logSuccess("POST /loans", fmt.Sprintf("Created %s loan (ID: %d) - Principal: %.2f, Total: %.2f",
			data.LoanType, data.ID, data.PrincipalAmount, data.TotalPayableAmount))
		return &data.ID
	}

	logFailure("POST /loans", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testGetLoan(loanID uint) *LoanResponse {
	endpoint := fmt.Sprintf("/loans/%d", loanID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /loans/:id", err.Error())
		return nil
	}
	defer resp.Body.Close()

	var data LoanResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /loans/:id", fmt.Sprintf("Retrieved loan - Pending: %.2f, Status: %s",
			data.PendingAmount, data.Status))
		return &data
	}

	logFailure("GET /loans/:id", fmt.Sprintf("Status: %d", resp.StatusCode))
	return nil
}

func testGetAllBanks() {
	logRequest("GET", "/banks", nil)

	resp, err := makeRequest("GET", baseURL+"/banks", nil)
	if err != nil {
		logFailure("GET /banks", err.Error())
		return
	}
	defer resp.Body.Close()

	var data []BankResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /banks", fmt.Sprintf("Retrieved %d banks", len(data)))
		return
	}

	logFailure("GET /banks", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testGetBranchesByBank(bankID uint) {
	endpoint := fmt.Sprintf("/banks/%d/branches", bankID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /banks/:bank_id/branches", err.Error())
		return
	}
	defer resp.Body.Close()

	var data []BranchResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /banks/:bank_id/branches", fmt.Sprintf("Retrieved %d branches", len(data)))
		return
	}

	logFailure("GET /banks/:bank_id/branches", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testGetCustomersByBranch(branchID uint) {
	endpoint := fmt.Sprintf("/branches/%d/customers", branchID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /branches/:branch_id/customers", err.Error())
		return
	}
	defer resp.Body.Close()

	var data []CustomerResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /branches/:branch_id/customers", fmt.Sprintf("Retrieved %d customers", len(data)))
		return
	}

	logFailure("GET /branches/:branch_id/customers", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testAddAccountHolder(accountID, customerID uint, holderRole string) {
	payload := map[string]interface{}{
		"customer_id": customerID,
		"holder_role": holderRole,
	}
	endpoint := fmt.Sprintf("/accounts/%d/holders", accountID)
	logRequest("POST", endpoint, payload)

	resp, err := makeRequest("POST", baseURL+endpoint, payload)
	if err != nil {
		logFailure("POST /accounts/:account_id/holders", err.Error())
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 201 {
		logSuccess("POST /accounts/:account_id/holders", fmt.Sprintf("Added holder to account %d", accountID))
		return
	}

	logFailure("POST /accounts/:account_id/holders", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testDeposit(accountID uint, amount float64) {
	payload := map[string]interface{}{
		"amount": amount,
	}
	endpoint := fmt.Sprintf("/accounts/%d/deposit", accountID)
	logRequest("POST", endpoint, payload)

	resp, err := makeRequest("POST", baseURL+endpoint, payload)
	if err != nil {
		logFailure("POST /accounts/:account_id/deposit", err.Error())
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("POST /accounts/:account_id/deposit", fmt.Sprintf("Deposited %.2f to account %d", amount, accountID))
		return
	}

	logFailure("POST /accounts/:account_id/deposit", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testWithdraw(accountID uint, amount float64) {
	payload := map[string]interface{}{
		"amount": amount,
	}
	endpoint := fmt.Sprintf("/accounts/%d/withdraw", accountID)
	logRequest("POST", endpoint, payload)

	resp, err := makeRequest("POST", baseURL+endpoint, payload)
	if err != nil {
		logFailure("POST /accounts/:account_id/withdraw", err.Error())
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("POST /accounts/:account_id/withdraw", fmt.Sprintf("Withdrew %.2f from account %d", amount, accountID))
		return
	}

	logFailure("POST /accounts/:account_id/withdraw", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testGetTransactions(accountID uint) {
	endpoint := fmt.Sprintf("/accounts/%d/transactions", accountID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /accounts/:account_id/transactions", err.Error())
		return
	}
	defer resp.Body.Close()

	var data []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /accounts/:account_id/transactions", fmt.Sprintf("Retrieved %d transactions", len(data)))
		return
	}

	logFailure("GET /accounts/:account_id/transactions", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testGetCustomerLoans(customerID uint) {
	endpoint := fmt.Sprintf("/customers/%d/loans", customerID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /customers/:customer_id/loans", err.Error())
		return
	}
	defer resp.Body.Close()

	var data []LoanResponse
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /customers/:customer_id/loans", fmt.Sprintf("Retrieved %d loans", len(data)))
		return
	}

	logFailure("GET /customers/:customer_id/loans", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testRepayLoan(loanID uint, amount float64) {
	payload := map[string]interface{}{
		"amount": amount,
	}
	endpoint := fmt.Sprintf("/loans/%d/repay", loanID)
	logRequest("POST", endpoint, payload)

	resp, err := makeRequest("POST", baseURL+endpoint, payload)
	if err != nil {
		logFailure("POST /loans/:id/repay", err.Error())
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("POST /loans/:id/repay", fmt.Sprintf("Repaid %.2f to loan %d", amount, loanID))
		return
	}

	logFailure("POST /loans/:id/repay", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testGetLoanInterest(loanID uint) {
	endpoint := fmt.Sprintf("/loans/%d/interest", loanID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /loans/:id/interest", err.Error())
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /loans/:id/interest", fmt.Sprintf("Retrieved interest for loan %d", loanID))
		return
	}

	logFailure("GET /loans/:id/interest", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func testGetLoanPayments(loanID uint) {
	endpoint := fmt.Sprintf("/loans/%d/payments", loanID)
	logRequest("GET", endpoint, nil)

	resp, err := makeRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logFailure("GET /loans/:id/payments", err.Error())
		return
	}
	defer resp.Body.Close()

	var data []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	logResponse(resp.StatusCode, data)

	if resp.StatusCode == 200 {
		logSuccess("GET /loans/:id/payments", fmt.Sprintf("Retrieved %d payments", len(data)))
		return
	}

	logFailure("GET /loans/:id/payments", fmt.Sprintf("Status: %d", resp.StatusCode))
}

func main() {
	printSeparator()
	logInfo(fmt.Sprintf("Starting Banking System API Tests - %s", time.Now().Format("2006-01-02 15:04:05")))
	printSeparator()

	// Test 1: Create Banks
	logInfo("\n[STEP 1] Creating Banks...")
	bank1ID := testCreateBank("State Bank of India")
	bank2ID := testCreateBank("HDFC Bank")

	if bank1ID == nil || bank2ID == nil {
		logWarning("Cannot proceed without banks")
		return
	}

	// Test 2: Get Bank Details
	logInfo("\n[STEP 2] Retrieving Bank Details...")
	testGetBank(*bank1ID)
	testGetBank(*bank2ID)

	// Test 3: Create Branches
	logInfo("\n[STEP 3] Creating Branches...")
	branch1ID := testCreateBranch(*bank1ID, "Mumbai Main Branch", "123 Marine Drive, Mumbai")
	branch2ID := testCreateBranch(*bank1ID, "Delhi Branch", "456 Connaught Place, Delhi")
	branch3ID := testCreateBranch(*bank2ID, "Bangalore Branch", "789 MG Road, Bangalore")

	if branch1ID == nil || branch2ID == nil || branch3ID == nil {
		logWarning("Cannot proceed without branches")
		return
	}

	// Test 4: Get Branch Details
	logInfo("\n[STEP 4] Retrieving Branch Details...")
	testGetBranch(*branch1ID)
	testGetBranch(*branch2ID)

	// Test 5: Create Customers
	logInfo("\n[STEP 5] Creating Customers...")
	timestamp := time.Now().Format("20060102150405")
	customer1ID := testCreateCustomer(*branch1ID, "Rajesh Kumar", fmt.Sprintf("rajesh.%s@example.com", timestamp), "9876543210")
	time.Sleep(100 * time.Millisecond)
	timestamp2 := time.Now().Format("20060102150405")
	customer2ID := testCreateCustomer(*branch2ID, "Priya Sharma", fmt.Sprintf("priya.%s@example.com", timestamp2), "9876543211")
	time.Sleep(100 * time.Millisecond)
	timestamp3 := time.Now().Format("20060102150405")
	customer3ID := testCreateCustomer(*branch3ID, "Amit Patel", fmt.Sprintf("amit.%s@example.com", timestamp3), "9876543212")

	if customer1ID == nil || customer2ID == nil || customer3ID == nil {
		logWarning("Cannot proceed without customers")
		return
	}

	// Test 6: Get Customer Details
	logInfo("\n[STEP 6] Retrieving Customer Details...")
	testGetCustomer(*customer1ID)
	testGetCustomer(*customer2ID)

	// Test 7: Create Accounts
	logInfo("\n[STEP 7] Creating Savings Accounts...")
	account1ID := testCreateAccount(*customer1ID, "primary_holder")
	account2ID := testCreateAccount(*customer1ID, "primary_holder") // Second account for customer 1
	account3ID := testCreateAccount(*customer2ID, "primary_holder")
	account4ID := testCreateAccount(*customer3ID, "primary_holder")
	account5ID := testCreateAccount(*customer3ID, "primary_holder") // Second account for customer 3

	if account1ID == nil || account3ID == nil {
		logWarning("Some accounts failed to create")
	}

	// Test 8: Get Account Details
	logInfo("\n[STEP 8] Retrieving Account Details...")
	if account1ID != nil {
		testGetAccount(*account1ID)
	}
	if account2ID != nil {
		testGetAccount(*account2ID)
	}
	if account3ID != nil {
		testGetAccount(*account3ID)
	}
	if account4ID != nil {
		testGetAccount(*account4ID)
	}
	if account5ID != nil {
		testGetAccount(*account5ID)
	}

	// Test 9: Create Loans
	logInfo("\n[STEP 9] Creating Loans...")
	loan1ID := testCreateLoan(*customer1ID, "Home Loan", 500000)
	loan2ID := testCreateLoan(*customer2ID, "Personal Loan", 100000)
	loan3ID := testCreateLoan(*customer3ID, "Car Loan", 300000)
	loan4ID := testCreateLoan(*customer1ID, "Education Loan", 200000) // Second loan for customer 1

	// Test 10: Get Loan Details
	logInfo("\n[STEP 10] Retrieving Loan Details...")
	if loan1ID != nil {
		testGetLoan(*loan1ID)
	}
	if loan2ID != nil {
		testGetLoan(*loan2ID)
	}
	if loan3ID != nil {
		testGetLoan(*loan3ID)
	}
	if loan4ID != nil {
		testGetLoan(*loan4ID)
	}

	// Test 11: List All Customers with Accounts
	logInfo("\n[STEP 11] Listing All Customers with Account Details...")
	customersData := testListCustomersWithAccounts()

	if customersData != nil {
		printSeparator()
		logInfo("CUSTOMER ACCOUNT SUMMARY:")
		printSeparator()
		for _, customer := range customersData {
			fmt.Printf("\n%sCustomer: %s%s\n", colorBlue, customer.Name, colorReset)
			fmt.Printf("  Email: %s\n", customer.Email)
			fmt.Printf("  Phone: %s\n", customer.Phone)
			fmt.Printf("  Number of Accounts: %d\n", len(customer.Accounts))

			if len(customer.Accounts) > 0 {
				fmt.Printf("  %sAccounts:%s\n", colorGreen, colorReset)
				for idx, account := range customer.Accounts {
					fmt.Printf("    %d. Account No: %d\n", idx+1, account.AccountNo)
					fmt.Printf("       Bank: %s\n", account.BankName)
					fmt.Printf("       Branch: %s\n", account.BranchName)
					fmt.Printf("       Balance: ₹%.2f\n", account.Balance)
				}
			} else {
				fmt.Printf("  %sNo accounts found%s\n", colorYellow, colorReset)
			}
		}
	}

	// Test 12: Get All Banks
	logInfo("\n[STEP 12] Retrieving All Banks...")
	testGetAllBanks()

	// Test 13: Get Branches by Bank
	logInfo("\n[STEP 13] Retrieving Branches by Bank...")
	if bank1ID != nil {
		testGetBranchesByBank(*bank1ID)
	}

	// Test 14: Get Customers by Branch
	logInfo("\n[STEP 14] Retrieving Customers by Branch...")
	if branch1ID != nil {
		testGetCustomersByBranch(*branch1ID)
	}

	// Test 15: Add Account Holder (Joint Account)
	logInfo("\n[STEP 15] Adding Joint Account Holder...")
	if account1ID != nil && customer2ID != nil {
		testAddAccountHolder(*account1ID, *customer2ID, "joint_holder")
	}

	// Test 16: Deposit Money
	logInfo("\n[STEP 16] Testing Deposit Operations...")
	if account1ID != nil {
		testDeposit(*account1ID, 10000.00)
	}
	if account3ID != nil {
		testDeposit(*account3ID, 25000.00)
	}

	// Test 17: Withdraw Money
	logInfo("\n[STEP 17] Testing Withdrawal Operations...")
	if account1ID != nil {
		testWithdraw(*account1ID, 2000.00)
	}

	// Test 18: Get Transaction History
	logInfo("\n[STEP 18] Retrieving Transaction History...")
	if account1ID != nil {
		testGetTransactions(*account1ID)
	}
	if account3ID != nil {
		testGetTransactions(*account3ID)
	}

	// Test 19: Get Customer Loans
	logInfo("\n[STEP 19] Retrieving Customer Loans...")
	if customer1ID != nil {
		testGetCustomerLoans(*customer1ID)
	}

	// Test 20: Repay Loan
	logInfo("\n[STEP 20] Testing Loan Repayment...")
	if loan1ID != nil {
		testRepayLoan(*loan1ID, 50000.00)
	}
	if loan2ID != nil {
		testRepayLoan(*loan2ID, 112000.00) // Full repayment
	}

	// Test 21: Get Loan Interest
	logInfo("\n[STEP 21] Retrieving Loan Interest...")
	if loan1ID != nil {
		testGetLoanInterest(*loan1ID)
	}
	if loan3ID != nil {
		testGetLoanInterest(*loan3ID)
	}

	// Test 22: Get Loan Payment History
	logInfo("\n[STEP 22] Retrieving Loan Payment History...")
	if loan1ID != nil {
		testGetLoanPayments(*loan1ID)
	}
	if loan2ID != nil {
		testGetLoanPayments(*loan2ID)
	}

	// Final Summary
	printSeparator()
	logInfo("TEST SUMMARY:")
	printSeparator()

	totalTests := len(apiResults.passed) + len(apiResults.failed)
	passRate := 0.0
	if totalTests > 0 {
		passRate = float64(len(apiResults.passed)) / float64(totalTests) * 100
	}

	fmt.Printf("\n%s✓ Passed: %d/%d%s\n", colorGreen, len(apiResults.passed), totalTests, colorReset)
	fmt.Printf("%s✗ Failed: %d/%d%s\n", colorRed, len(apiResults.failed), totalTests, colorReset)
	fmt.Printf("%sPass Rate: %.1f%%%s\n\n", colorBlue, passRate, colorReset)

	if len(apiResults.passed) > 0 {
		fmt.Printf("%sAPIs that worked:%s\n", colorGreen, colorReset)
		for _, api := range apiResults.passed {
			fmt.Printf("  ✓ %s\n", api)
		}
	}

	if len(apiResults.failed) > 0 {
		fmt.Printf("\n%sAPIs that failed:%s\n", colorRed, colorReset)
		for _, api := range apiResults.failed {
			fmt.Printf("  ✗ %s\n", api)
		}
	}

	printSeparator()
}
