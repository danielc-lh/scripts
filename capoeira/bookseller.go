package capoeira

import (
	"fmt"
	"sync"
	"time"
)

// 1. Define the locations

type Seller struct{}

func (Seller) Name() string { return "Seller" }

type Buyer struct{}

func (Buyer) Name() string { return "Buyer" }

// Helper function for book lookup
func getBook(title string) (price int, deliveryDate time.Time, found bool) {
	switch title {
	case "TAPL":
		return 80, time.Date(2023, 8, 3, 0, 0, 0, 0, time.UTC), true
	case "HoTT":
		return 120, time.Date(2023, 9, 18, 0, 0, 0, 0, time.UTC), true
	default:
		return 0, time.Time{}, false
	}
}

const BUDGET = 100

// 2. Define the BooksellerChoreography struct
type BooksellerChoreography struct {
	Title  Located
	Budget Located
}

func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	default:
		return 0
	}
}

func toDate(v any) time.Time {
	switch val := v.(type) {
	case time.Time:
		return val
	case string:
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			return time.Time{}
		}
		return t
	default:
		return time.Time{}
	}
}

// 3. Provide an implementation of Run(...)
func (c BooksellerChoreography) Run(op ChoreoOp) interface{} {

	// Buyer sends title of book they want to
	titleAtSeller := op.Comm(Buyer{}, Seller{}, c.Title)
	fmt.Printf("Title at seller: %v from op: %v\n", titleAtSeller.Value, op.(ProjectorChoreoOp).Target.Name())
	priceAtSeller := op.Locally(Seller{}, func() interface{} {
		title := titleAtSeller.Value.(string)
		price, _, found := getBook(title)
		fmt.Println("Price: ", price)
		if found {
			return price
		}
		return nil
	})

	priceAtBuyer := op.Comm(Seller{}, Buyer{}, priceAtSeller)
	decisionAtBuyer := op.Locally(Buyer{}, func() interface{} {
		if priceAtBuyer.Value != nil {
			price := toInt(priceAtBuyer.Value)
			fmt.Printf("Buyer: Price is %d\n", price)
			budget := toInt(c.Budget.Value)
			decision := price < budget
			if decision {
				fmt.Printf("The buyer can buy the book, since $%v < $%v\n", price, budget)
			}
			return decision
		}
		fmt.Println("The book does not exist")
		return false
	})
	decision := op.Broadcast(Buyer{}, decisionAtBuyer).(bool)
	if decision {
		deliveryDateAtSeller := op.Locally(Seller{}, func() interface{} {
			title := titleAtSeller.Value.(string)
			_, deliveryDate, _ := getBook(title)
			return deliveryDate
		})
		deliveryDateAtBuyer := op.Comm(Seller{}, Buyer{}, deliveryDateAtSeller)
		op.Locally(Buyer{}, func() interface{} {
			deliveryDate := toDate(deliveryDateAtBuyer.Value)
			fmt.Printf("The book will be delivered on %s\n", deliveryDate.Format(time.RFC3339))
			return nil
		})
	} else {
		op.Locally(Buyer{}, func() interface{} {
			fmt.Println("The buyer cannot buy the book")
			return nil
		})
	}
	return decision
}

// 4. RunBookSellerProtocol: creates transports, projectors, and runs each endpoint
func RunBookSellerProtocol(title string, transport Transport) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Seller endpoint
	go func() {
		defer wg.Done()
		sellerProjector := NewProjector(Seller{}, transport)
		sellerProjector.EppAndRun(
			BooksellerChoreography{
				Title:  sellerProjector.Remote(Buyer{}),
				Budget: sellerProjector.Remote(Buyer{}),
			},
		)
	}()

	// Buyer endpoint
	go func() {
		defer wg.Done()
		buyerProjector := NewProjector(Buyer{}, transport)
		buyerProjector.EppAndRun(
			BooksellerChoreography{
				Title:  buyerProjector.Local(title),
				Budget: buyerProjector.Local(BUDGET),
			},
		)
	}()

	wg.Wait()
}
