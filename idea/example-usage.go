package main

import (
	"log"
	"net/http"
)

// Calculator implements CalculatorHandler
type Calculator struct {
	cc *CalculatorClient
	ci *ClientInfo
}

func (ch *CalculatorHandler) Add(a int, b int) (int, error) {
	return a + b, nil
}

func (ch *CalculatorHandler) Subtract(a int, b int) (int, error) {
	return b - a, nil
}

func (ch *CalculatorHandler) Clear() {
	// do nothing
}

func main() {
	calculatorServer := NewCalculatorServer(func(cc *CalculatorClient, ci *ClientInfo) *CalculatorHandler {
		return &Calculator{
			cc: cc,
			ci: ci,
		}
	})

	http.Handle("/", calculatorServer)
	err := http.ListenAndServe(":1324", nil)
	if err != nil {
		log.Fatal(err)
	}
}
