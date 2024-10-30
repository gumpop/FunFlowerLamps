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
	TimeStamp    string
}

type Customer struct {
	ID    int
	FName string
	LName string
	Email string
}

type Product struct {
	ID      int
	Name    string
	Image   string
	Price   float64
	InStock int
}

type Order struct {
	ID           int
	CustID       int
	ProdID       int
	ProductPrice float64
	Quantity     int64
	RoundUp      bool
	SubTotal     float64
	Tax          float64
	Round        float64
	TimeStamp    int64
}

type CustomerResults struct {
	Customers []Customer
	Num       int
	Customer2 Customer
	Customer3 string
	Customer4 Customer
	Customer5 string
	Customer6 Customer
}

type OrderResults struct {
	Num      int
	Orders   []PurchaseInfo
	AfterNum int
}

type ProductResults struct {
	Products []Product
	Num1     int
	Num2     int
}
