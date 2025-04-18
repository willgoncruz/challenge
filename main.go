package main

import (
	css "challenge/client"
	"challenge/ledger"
	"challenge/scheduler"
	"flag"
	"log"
	"sync"
	"time"
)

var (
	endpoint = flag.String("endpoint", "https://api.cloudkitchens.com", "Problem server endpoint")
	auth     = flag.String("auth", "", "Authentication token (required)")
	name     = flag.String("name", "", "Problem name. Leave blank (optional)")
	seed     = flag.Int64("seed", 0, "Problem seed (random if zero)")

	rate = flag.Duration("rate", 500*time.Millisecond, "Inverse order rate")
	min  = flag.Duration("min", 4*time.Second, "Minimum pickup time")
	max  = flag.Duration("max", 8*time.Second, "Maximum pickup time")
)

func main() {
	flag.Parse()

	client := css.NewClient(*endpoint, *auth)
	id, orders, err := client.New(*name, *seed)
	if err != nil {
		log.Fatalf("Failed to fetch test problem: %v", err)
	}

	// ------ Simulation harness logic goes here using rate, min and max ------
	// var actions []model.Action
	// for _, order := range orders {
	// 	log.Printf("Received: %+v", order)

	// 	// actions = append(actions, model.Action{Timestamp: time.Now().UnixMicro(), ID: order.ID, Action: model.Place})
	// 	ledger.Audit(order, model.Place)
	// 	time.Sleep(*rate)
	// }
	// ------------------------------------------------------------------------

	// -------------- Schedulers --------------------------
	wg := &sync.WaitGroup{}
	wg.Add(1) // Wait for schedulers to finish

	// Async process the place and pickup orders
	go scheduler.NewPlaceScheduler(rate, scheduler.NewPickupScheduler(min, max)).Process(orders, wg)
	// go scheduler.NewPickupScheduler(min, max).Process(orders, wg)

	wg.Wait()
	// ------------------------------------------------------------------------

	result, err := client.Solve(id, *rate, *min, *max, ledger.Retrieve())
	if err != nil {
		log.Fatalf("Failed to submit test solution: %v", err)
	}
	log.Printf("Test result: %v", result)
}
