package model

// Uncategorized is the category assigned to any transaction that no rule
// matched. Keeping it an explicit constant (rather than the empty string)
// means summaries and exports always show a real bucket.
const Uncategorized = "Uncategorized"

// The default category vocabulary. Rules are free to emit any string, but
// shipping a sensible baseline keeps the default rules.yaml and summaries
// consistent. These intentionally read as everyday budgeting buckets.
const (
	CategoryGroceries     = "Groceries"
	CategoryRent          = "Rent"
	CategoryUtilities     = "Utilities"
	CategorySubscriptions = "Subscriptions"
	CategoryTransport     = "Transport"
	CategoryDiningOut     = "Dining Out"
	CategoryShopping      = "Shopping"
	CategoryHealth        = "Health"
	CategoryIncome        = "Income"
	CategoryCash          = "Cash"
	CategoryFees          = "Fees"
	CategoryInsurance     = "Insurance"
)

// DefaultCategories lists the built-in buckets in display order.
var DefaultCategories = []string{
	CategoryIncome,
	CategoryGroceries,
	CategoryDiningOut,
	CategoryRent,
	CategoryUtilities,
	CategorySubscriptions,
	CategoryTransport,
	CategoryShopping,
	CategoryHealth,
	CategoryInsurance,
	CategoryFees,
	CategoryCash,
	Uncategorized,
}

// CategoryOrDefault returns the given category, falling back to Uncategorized
// when it is empty.
func CategoryOrDefault(c string) string {
	if c == "" {
		return Uncategorized
	}
	return c
}
