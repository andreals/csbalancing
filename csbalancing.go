package csbalancing

import (
	"errors"
	"sort"
	"sync"
)

// Entity is a struct to manipulate CSs and Customers data
type Entity struct {
	ID    int
	Score int
}

type byScore []Entity

func (e byScore) Len() int           { return len(e) }
func (e byScore) Less(i, j int) bool { return e[i].Score < e[j].Score }
func (e byScore) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// csQuantity is a struct to manipulate the quantity of customers a CS has
type csQuantity struct {
	ID       int
	Quantity int
}

type byQuantity []csQuantity

func (p byQuantity) Len() int           { return len(p) }
func (p byQuantity) Less(i, j int) bool { return p[i].Quantity < p[j].Quantity }
func (p byQuantity) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// CustomerSuccessBalancing is a function to balance the number of CS customers, and return the CS that has the most customers
func CustomerSuccessBalancing(customerSuccess []Entity, customers []Entity, customerSuccessAway []int) int {

	var (
		wg     sync.WaitGroup
		csQtys []csQuantity
	)

	// Remove CSs aways from the slice
	for i := 0; i < len(customerSuccess); i++ {
		if isCSAway(customerSuccess[i].ID, customerSuccessAway) {
			customerSuccess = append(customerSuccess[:i], customerSuccess[i+1:]...)
			i--
		}
	}

	// Sort entities by Score ascending
	sort.Sort(byScore(customerSuccess))
	sort.Sort(byScore(customers))

	for _, customer := range customers {

		wg.Add(1)
		go func(customer Entity) {

			defer wg.Done()
			for _, cs := range customerSuccess {

				// Checks if the customer's score is less than or equal to the CS score
				if customer.Score <= cs.Score {

					// check if CS already has a costumer, if so, increment, if not add
					idx, err := getCSQtyIndex(csQtys, cs.ID)
					if err != nil {
						csQtys = append(csQtys, csQuantity{ID: cs.ID, Quantity: 1})
						return
					}

					csQtys[idx].Quantity++
					return
				}
			}

		}(customer)
	}

	wg.Wait()

	// If there were no CS, it returns 0
	if len(csQtys) < 1 {
		return 0
	}

	// Sort CSs by Quantity descending
	sort.Sort(sort.Reverse(byQuantity(csQtys)))

	// Check if there was a tie
	if len(csQtys) > 1 && csQtys[0].Quantity == csQtys[1].Quantity {
		return 0
	}

	// Returns the first position of CS (the one with the highest number of customers)
	return csQtys[0].ID
}

// getCSQtyIndex returns the index of the slice that has the customer ID entered or an error if it does not exist
func getCSQtyIndex(customerSuccessQuantitys []csQuantity, customerSuccessID int) (int, error) {
	for i, csQty := range customerSuccessQuantitys {
		if csQty.ID == customerSuccessID {
			return i, nil
		}
	}
	return 0, errors.New("index not found")
}

// isCSAway returns if received CS ID is away
func isCSAway(customerSuccessID int, customerSuccessAway []int) bool {
	for _, idx := range customerSuccessAway {
		if customerSuccessID == idx {
			return true
		}
	}
	return false
}
