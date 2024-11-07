package database

import (
	"database/sql"
	"fmt"
	"go-store/types"
	"time"
)

func GetAllCustomers(db *sql.DB) ([]types.Customer, error) {
	var Customers []types.Customer
	rows, err := db.Query("SELECT * FROM customers")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var customer types.Customer
		if err := rows.Scan(&customer.ID, &customer.FName, &customer.LName, &customer.Email); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		Customers = append(Customers, customer)
	}

	return Customers, nil
}

func GetAllProducts(db *sql.DB) ([]types.Product, error) {

	var theProducts []types.Product
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product types.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Image, &product.Price, &product.InStock); err != nil {
			return nil, fmt.Errorf(" %v", err)
		}
		theProducts = append(theProducts, product)
	}

	return theProducts, nil
}

func GetAllOrders(db *sql.DB) ([]types.PurchaseInfo, error) {
	var Orders []types.PurchaseInfo
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("order : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var order types.Order
		if err := rows.Scan(&order.ID, &order.ProdID, &order.CustID, &order.Quantity, &order.ProductPrice, &order.Tax, &order.Round, &order.TimeStamp); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		Orders = append(Orders, convToPurchase(order, db))
	}

	return Orders, nil
}

// conv order to purchase info
func convToPurchase(order types.Order, conn *sql.DB) types.PurchaseInfo {
	var purchase types.PurchaseInfo
	var cust types.Customer
	var prod types.Product

	cust, _ = GetCustomerById(order.CustID, conn)
	prod, _ = GetProductByID(order.ProdID, conn)

	purchase.FName = cust.FName
	purchase.LName = cust.LName
	purchase.Email = cust.Email
	purchase.ProductName = prod.Name
	purchase.ProductPrice = prod.Price
	purchase.Quantity = order.Quantity
	purchase.SubTotal = order.SubTotal
	purchase.TotalTax = order.SubTotal + order.Tax
	purchase.TotalRound = order.SubTotal + order.Tax + order.Round
	if order.Round == 0 {
		purchase.RoundUp = false
	} else {
		purchase.RoundUp = true
	}
	purchase.TimeStamp = time.UnixMilli(int64(order.TimeStamp)).String()

	return purchase
}

func GetNumberOrders(db *sql.DB) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM orders")
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	var num int
	rows.Next()
	if err := rows.Scan(&num); err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	return num, nil
}
func GetNumberCustomers(db *sql.DB) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM customers")
	if err != nil {
		return -1, fmt.Errorf("%v", err)
	}
	defer rows.Close()
	var num int
	rows.Next()
	if err := rows.Scan(&num); err != nil {
		return -1, fmt.Errorf("%v", err)
	}

	return num, nil
}

func GetCustomerById(id int, db *sql.DB) (types.Customer, error) {
	rows, _ := db.Prepare("SELECT * FROM customers WHERE id=?")
	defer rows.Close()
	var customer types.Customer
	err := rows.QueryRow(id).Scan(&customer.ID, &customer.FName, &customer.LName, &customer.Email)
	if err != nil {
		return customer, fmt.Errorf("%v", err)
	}

	return customer, nil
}

func GetCustomerByEmail(email string, db *sql.DB) (types.Customer, error) {
	rows, _ := db.Prepare("SELECT * FROM customers WHERE email=?")
	defer rows.Close()

	var customer types.Customer
	err := rows.QueryRow(email).Scan(&customer.ID, &customer.FName, &customer.LName, &customer.Email)
	if err != nil {
		return customer, fmt.Errorf("%v", err)
	}

	return customer, nil
}

func GetCustomersByFirst(pattern string, db *sql.DB) []types.Customer {
	stmt, _ := db.Prepare("SELECT * FROM customers WHERE FirstName LIKE ?")
	rows, _ := stmt.Query(fmt.Sprintf("%%%s%%", pattern))

	var customer types.Customer
	var customers []types.Customer
	for rows.Next() {
		_ = rows.Scan(&customer.ID, &customer.FName, &customer.LName, &customer.Email)
		customers = append(customers, customer)
	}

	return customers
}

func GetCustomersByLast(pattern string, db *sql.DB) []types.Customer {
	stmt, _ := db.Prepare("SELECT * FROM customers WHERE LastName LIKE ?")
	rows, _ := stmt.Query(fmt.Sprintf("%%%s%%", pattern))

	var customer types.Customer
	var customers []types.Customer
	for rows.Next() {
		_ = rows.Scan(&customer.ID, &customer.FName, &customer.LName, &customer.Email)
		customers = append(customers, customer)
	}

	return customers
}

func AddCustomer(email string, fName string, lName string, db *sql.DB) error {
	rows, _ := db.Prepare("INSERT INTO customers (FirstName, LastName, email) VALUES (?, ?, ?)")
	defer rows.Close()

	_, err := rows.Exec(fName, lName, email)

	return err
}

func AddOrder(pID int, cID int, num int, price float64, tax float64, donation float64, timestamp int64, db *sql.DB) error {
	rows, _ := db.Prepare("INSERT INTO orders (product_id, customer_id, quantity, price, tax, donation, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer rows.Close()

	_, err := rows.Exec(pID, cID, num, price, tax, donation, timestamp)

	return err
}

func SellProduct(num int, id int, db *sql.DB) error {
	var stock int
	stock, err := GetInStock(id, db)

	if num > stock {
		rows, _ := db.Prepare("UPDATE products SET in_stock = 0 WHERE id=?")
		defer rows.Close()

		_, err2 := rows.Exec(id)
		return err2
	} else {
		var newStock = stock - num
		rows, _ := db.Prepare("UPDATE products SET in_stock = ? WHERE id=?")
		defer rows.Close()

		_, err2 := rows.Exec(newStock, id)
		return err2
	}

	return err
}

func GetInStock(id int, db *sql.DB) (int, error) {
	var stock int
	rows, _ := db.Prepare("SELECT in_stock FROM products WHERE id=?")
	defer rows.Close()

	err := rows.QueryRow(id).Scan(&stock)
	return stock, err
}

func GetInStockByName(name string, db *sql.DB) (int, error) {
	var stock int
	rows, _ := db.Prepare("SELECT in_stock FROM products WHERE product_name=?")
	defer rows.Close()

	err := rows.QueryRow(name).Scan(&stock)
	return stock, err
}

func GetProductByID(id int, db *sql.DB) (types.Product, error) {
	rows, _ := db.Prepare("SELECT * FROM products WHERE id=?")
	defer rows.Close()
	var product types.Product
	err := rows.QueryRow(id).Scan(&product.ID, &product.Name, &product.Image, &product.Price, &product.InStock)
	if err != nil {
		return product, fmt.Errorf("%v", err)
	}

	return product, nil
}

func CheckOrder(cid int, pid int, time int64, db *sql.DB) (bool, error) {
	rows, _ := db.Prepare("SELECT COUNT(*) FROM orders WHERE product_id=? AND customer_id=? AND timestamp=?")
	defer rows.Close()
	var num int
	err := rows.QueryRow(pid, cid, time).Scan(&num)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	if num == 0 {
		return true, nil
	}
	return false, nil
}
