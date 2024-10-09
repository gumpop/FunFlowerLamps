package types

// TODO: If you choose to use a struct rather than individual parameters to your view, you might flesh this one out:
type PurchaseInfo struct {
	FName        string
	LName        string
	Email        string
	ProductName  string
	ProductPrice float64
	Quantity     int64
	RoundUp      bool
	SubTotal     float64
	TotalTax     float64
	TotalRound   float64
}
