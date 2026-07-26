package finance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	financetemplates "github.com/Suke2004/atlas-go/web/templates/finance"
)

type TransactionInput struct {
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"` // "income" | "expense"
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
}

type FinanceSummary = financetemplates.FinanceSummary
type ProjectCostAttribution = financetemplates.ProjectCostAttribution

type Service interface {
	CreateTransaction(ctx context.Context, userID int64, input TransactionInput) (db.Transaction, error)
	ListTransactions(ctx context.Context, userID int64) ([]db.Transaction, error)
	GetFinanceSummary(ctx context.Context, userID int64) (FinanceSummary, error)
	DeleteTransaction(ctx context.Context, userID, transactionID int64) error
}

type service struct {
	repo         Repository
	projectsRepo projects.Repository
}

func NewService(repo Repository, projRepo projects.Repository) Service {
	return &service{
		repo:         repo,
		projectsRepo: projRepo,
	}
}

func (s *service) CreateTransaction(ctx context.Context, userID int64, input TransactionInput) (db.Transaction, error) {
	if input.Amount <= 0 {
		return db.Transaction{}, fmt.Errorf("amount must be greater than zero")
	}

	if input.TransactionDate == "" {
		input.TransactionDate = time.Now().Format("2006-01-02")
	}

	return s.repo.CreateTransaction(ctx, db.CreateTransactionParams{
		UserID:          userID,
		Amount:          input.Amount,
		Type:            input.Type,
		Category:        strings.TrimSpace(input.Category),
		Description:     strings.TrimSpace(input.Description),
		TransactionDate: input.TransactionDate,
	})
}

func (s *service) ListTransactions(ctx context.Context, userID int64) ([]db.Transaction, error) {
	return s.repo.ListTransactions(ctx, userID)
}

func (s *service) GetFinanceSummary(ctx context.Context, userID int64) (FinanceSummary, error) {
	summaryRow, err := s.repo.GetFinanceSummary(ctx, userID)
	if err != nil {
		return FinanceSummary{}, err
	}

	allTxns, _ := s.repo.ListTransactions(ctx, userID)

	netSavings := summaryRow.TotalIncome - summaryRow.TotalExpenses
	var savingsRate int
	if summaryRow.TotalIncome > 0 {
		savingsRate = int((netSavings / summaryRow.TotalIncome) * 100)
	}

	// Calculate Project Cost Attribution USP
	var projCosts []ProjectCostAttribution
	if s.projectsRepo != nil {
		projList, err := s.projectsRepo.ListProjects(ctx, userID)
		if err == nil && len(projList) > 0 {
			costPerProj := summaryRow.TotalExpenses / float64(len(projList))
			for _, p := range projList {
				pct := 0
				if summaryRow.TotalExpenses > 0 {
					pct = int((costPerProj / summaryRow.TotalExpenses) * 100)
				}
				projCosts = append(projCosts, ProjectCostAttribution{
					ProjectName: p.Name,
					MonthlyCost: costPerProj,
					Percentage:  pct,
				})
			}
		}
	}

	return FinanceSummary{
		TotalIncome:   summaryRow.TotalIncome,
		TotalExpenses: summaryRow.TotalExpenses,
		NetSavings:    netSavings,
		SavingsRate:   savingsRate,
		ProjectCosts:  projCosts,
		TopExpenses:   allTxns,
	}, nil
}

func (s *service) DeleteTransaction(ctx context.Context, userID, transactionID int64) error {
	return s.repo.DeleteTransaction(ctx, db.DeleteTransactionParams{ID: transactionID, UserID: userID})
}
