package finance

import "github.com/Suke2004/atlas-go/internal/db"

// ProjectCostAttribution represents how infrastructure cost is split across projects.
type ProjectCostAttribution struct {
	ProjectName string  `json:"project_name"`
	MonthlyCost float64 `json:"monthly_cost"`
	Percentage  int     `json:"percentage"`
}

// FinanceSummary is the aggregated financial view for the dashboard.
type FinanceSummary struct {
	TotalIncome   float64                  `json:"total_income"`
	TotalExpenses float64                  `json:"total_expenses"`
	NetSavings    float64                  `json:"net_savings"`
	SavingsRate   int                      `json:"savings_rate"`
	ProjectCosts  []ProjectCostAttribution `json:"project_costs"`
	TopExpenses   []db.Transaction         `json:"top_expenses"`
}
