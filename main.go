package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	etag "github.com/pablor21/echo-etag/v4"

	"go-store/templates"
	"go-store/types"
)

func main() {
	// TODO: Fill in your products here with name -> price as the key -> value pair.
	products := map[string]float64{
		"Lily":  199.99,
		"Lotus": 2212.50,
		"Rose":  95.29,
	}
	e := echo.New()
	e.Use(etag.Etag())

	e.Static("assets", "./assets")

	e.GET("/store", func(ctx echo.Context) error {
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(products)))
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
		purchaseInfo := types.PurchaseInfo{
			FName:        ctx.FormValue("fName"),
			LName:        ctx.FormValue("lName"),
			Email:        ctx.FormValue("email"),
			ProductName:  ctx.FormValue("productName"),
			ProductPrice: float64(products[ctx.FormValue("productName")]),
			Quantity:     q,
			RoundUp:      r,
			SubTotal:     float64(q) * float64(products[ctx.FormValue("productName")]),
			TotalTax:     float64(q)*float64(products[ctx.FormValue("productName")])*0.0875 + float64(q)*float64(products[ctx.FormValue("productName")]),
			TotalRound:   math.Ceil(float64(q)*float64(products[ctx.FormValue("productName")])*0.0875 + float64(q)*float64(products[ctx.FormValue("productName")])),
		}
		return Render(ctx, http.StatusOK, templates.Base(templates.PurchaseConfirmation(purchaseInfo)))
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
