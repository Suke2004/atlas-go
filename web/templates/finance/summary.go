package finance

import "github.com/Suke2004/atlas-go/internal/db"

type ProjectCostAttribution struct {
	ProjectName string  `json:"project_name"`
	MonthlyCost float64 `json:"monthly_cost"`
	Percentage  int     `json:"percentage"`
}

type FinanceSummary struct {
	TotalIncome  float64                  `json:"total_income"`
	TotalExpenses float64                  `json:"total_expenses"`
	NetSavings   float64                  `json:"net_savings"`
	SavingsRate  int                      `json:"savings_rate"`
	ProjectCosts []ProjectCostAttribution `json:"project_costs"`
	TopExpenses  []db.Transaction         `json:"top_expenses"`
}
