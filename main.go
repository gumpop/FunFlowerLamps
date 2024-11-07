package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	etag "github.com/pablor21/echo-etag/v4"

	dbs "go-store/db"
	"go-store/templates"
	"go-store/types"

	"github.com/go-sql-driver/mysql"
)

var conn *sql.DB

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "mmackay",
		Passwd: "Apple2",
		DBName: "mmackay",
	}
	// Get a database handle.
	var err error
	conn, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	e := echo.New()
	e.Use(etag.Etag())

	e.Static("assets", "./assets")

	e.GET("/store", func(ctx echo.Context) error {
		var products []types.Product
		products, err = dbs.GetAllProducts(conn)
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(products)))
	})

	e.GET("/order_entry", func(ctx echo.Context) error {
		var products []types.Product
		products, err = dbs.GetAllProducts(conn)
		return Render(ctx, http.StatusOK, templates.Base(templates.OrderEntry(products)))
	})

	e.GET("/get_prod_stock", func(ctx echo.Context) error {
		var stock int
		var id int64
		id, err := strconv.ParseInt(ctx.QueryParam("product"), 10, 64)
		if err != nil {
			//eating it NOMNOMNOM
		}
		stock, err = dbs.GetInStock(int(id), conn)

		return ctx.String(http.StatusOK, fmt.Sprintf("%d", stock))
	})

	e.GET("/get_customers", func(ctx echo.Context) error {
		var pattern url.Values
		pattern = ctx.QueryParams()
		var customers []types.Customer
		if pattern["fName"] != nil {
			customers = dbs.GetCustomersByFirst(ctx.QueryParam("fName"), conn)
		} else {
			customers = dbs.GetCustomersByLast(ctx.QueryParam("lName"), conn)
		}
		return Render(ctx, http.StatusOK, templates.CustTable(customers))
	})

	e.POST("/add_order", func(ctx echo.Context) error {

		q, err := strconv.ParseInt(ctx.FormValue("quantity"), 10, 64)
		if err != nil {
			//eating it NOMNOMNOM
		}
		r := false
		pd, err := strconv.ParseInt(ctx.FormValue("productName"), 10, 64)
		if err != nil {
			//eating it NOMNOMNOM
		}
		ts, err := strconv.ParseInt(ctx.FormValue("timestamp"), 10, 64)
		if err != nil {
			//eating it NOMNOMNOM
		}
		var cust types.Customer
		cust, err = dbs.GetCustomerByEmail(ctx.FormValue("email"), conn)
		var message string
		if err != nil {
			message = "Welcome new customer! Thank you for your order, " + ctx.FormValue("fName") + " " + ctx.FormValue("lName")
			dbs.AddCustomer(ctx.FormValue("email"), ctx.FormValue("fName"), ctx.FormValue("lName"), conn)
			cust, err = dbs.GetCustomerByEmail(ctx.FormValue("email"), conn)
		} else {
			message = "Welcome Back! Thank you for your order, " + ctx.FormValue("fName") + " " + ctx.FormValue("lName")
		}
		var prod types.Product
		prod, err = dbs.GetProductByID(int(pd), conn)
		var add bool
		add, err = dbs.CheckOrder(cust.ID, prod.ID, ts, conn)
		if add {
			dbs.AddOrder(prod.ID, cust.ID, int(q), prod.Price, prod.Price*float64(q)*0.0875, math.Ceil(prod.Price*float64(q)*1.0875)-prod.Price, ts, conn)
			dbs.SellProduct(int(q), prod.ID, conn)
		}

		var purchaseInfo types.PurchaseInfo
		purchaseInfo.Email = cust.Email
		purchaseInfo.FName = cust.FName
		purchaseInfo.LName = cust.LName
		purchaseInfo.ProductName = prod.Name
		purchaseInfo.ProductPrice = prod.Price
		purchaseInfo.Quantity = q
		purchaseInfo.RoundUp = r
		purchaseInfo.SubTotal = prod.Price * float64(q)
		purchaseInfo.TimeStamp = time.Now().String()
		purchaseInfo.TotalTax = prod.Price * float64(q) * 1.0875
		purchaseInfo.TotalRound = math.Ceil(prod.Price * float64(q) * 1.0875)

		return Render(ctx, http.StatusOK, templates.PurchaseConfirmation(message, purchaseInfo))
	})

	e.POST("/purchase", func(ctx echo.Context) error {

		q, err := strconv.ParseInt(ctx.FormValue("quantity"), 10, 64)
		if err != nil {
			//eating it NOMNOMNOM
		}
		r, err := strconv.ParseBool(ctx.FormValue("rup"))
		if err != nil {
			//eating it NOMNOMNOM
		}
		pd, err := strconv.ParseInt(ctx.FormValue("productName"), 10, 64)
		if err != nil {
			//eating it NOMNOMNOM
		}
		ts, err := strconv.ParseInt(ctx.FormValue("timestamp"), 10, 64)
		if err != nil {
			//eating it NOMNOMNOM
		}
		var cust types.Customer
		cust, err = dbs.GetCustomerByEmail(ctx.FormValue("email"), conn)
		var message string
		if err != nil {
			message = "Welcome new customer! Thank you for your order, " + ctx.FormValue("fName") + " " + ctx.FormValue("lName")
			dbs.AddCustomer(ctx.FormValue("email"), ctx.FormValue("fName"), ctx.FormValue("lName"), conn)
			cust, err = dbs.GetCustomerByEmail(ctx.FormValue("email"), conn)
		} else {
			message = "Welcome Back! Thank you for your order, " + ctx.FormValue("fName") + " " + ctx.FormValue("lName")
		}
		var prod types.Product
		prod, err = dbs.GetProductByID(int(pd), conn)
		var add bool
		add, err = dbs.CheckOrder(cust.ID, prod.ID, ts, conn)
		if add {
			dbs.AddOrder(prod.ID, cust.ID, int(q), prod.Price, prod.Price*float64(q)*0.0875, math.Ceil(prod.Price*float64(q)*1.0875)-prod.Price, ts, conn)
			dbs.SellProduct(int(q), prod.ID, conn)
		}

		var purchaseInfo types.PurchaseInfo
		purchaseInfo.Email = cust.Email
		purchaseInfo.FName = cust.FName
		purchaseInfo.LName = cust.LName
		purchaseInfo.ProductName = prod.Name
		purchaseInfo.ProductPrice = prod.Price
		purchaseInfo.Quantity = q
		purchaseInfo.RoundUp = r
		purchaseInfo.SubTotal = prod.Price * float64(q)
		purchaseInfo.TimeStamp = time.Now().String()
		purchaseInfo.TotalTax = prod.Price * float64(q) * 1.0875
		purchaseInfo.TotalRound = math.Ceil(prod.Price * float64(q) * 1.0875)

		return Render(ctx, http.StatusOK, templates.Base(templates.PurchaseConfirmation(message, purchaseInfo)))
	})

	e.GET("/dbQueries", func(ctx echo.Context) error {
		var customerResults types.CustomerResults
		customerResults.Customers, err = dbs.GetAllCustomers(conn)
		customerResults.Num, err = dbs.GetNumberCustomers(conn)
		customerResults.Customer2, err = dbs.GetCustomerById(2, conn)
		_, err = dbs.GetCustomerById(3, conn)
		if err != nil {
			customerResults.Customer3 = "Customer does not exist"
		}
		customerResults.Customer4, err = dbs.GetCustomerByEmail("mmackay@mines.edu", conn)
		if err != nil {
			fmt.Printf("%v", err)
		}
		_, err = dbs.GetCustomerByEmail("dduck@mines.edu", conn)
		if err != nil {
			customerResults.Customer5 = "Customer does not exist"
		}
		err = dbs.AddCustomer("dduck@mines.edu", "donald", "duck", conn)
		customerResults.Customer6, err = dbs.GetCustomerByEmail("dduck@mines.edu", conn)

		var orderResults types.OrderResults
		orderResults.Num, _ = dbs.GetNumberOrders(conn)
		err = dbs.AddOrder(1, 1, 2, 64.5, 2.5, 0.5, time.Now().Unix(), conn)
		orderResults.Orders, _ = dbs.GetAllOrders(conn)
		orderResults.AfterNum, _ = dbs.GetNumberOrders(conn)

		var productResults types.ProductResults
		productResults.Products, _ = dbs.GetAllProducts(conn)
		dbs.SellProduct(2, 1, conn)
		productResults.Num1, _ = dbs.GetInStock(1, conn)
		dbs.SellProduct(10, 1, conn)
		productResults.Num2, _ = dbs.GetInStock(1, conn)

		return Render(ctx, http.StatusOK, templates.Base(templates.Queries(customerResults, orderResults, productResults)))

	})

	e.GET("/admin", func(ctx echo.Context) error {
		var customerResults []types.Customer
		var orderResults []types.PurchaseInfo
		var productResults []types.Product

		customerResults, err := dbs.GetAllCustomers(conn)
		if err != nil {
			fmt.Printf("%v", err)
		}

		orderResults, errr := dbs.GetAllOrders(conn)
		if errr != nil {
			fmt.Printf("%v", errr)
		}

		productResults, errrr := dbs.GetAllProducts(conn)
		if errrr != nil {
			fmt.Printf("%v", errrr)
		}

		return Render(ctx, http.StatusOK, templates.Base(templates.Admin(customerResults, orderResults, productResults)))

	})

	e.Logger.Fatal(e.Start(":8000"))
}

// INFO: This is a simplified render method that replaces `echo`'s with a custom
// one. This should simplify rendering out of an echo route.
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
