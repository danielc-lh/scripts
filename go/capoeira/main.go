package main

import (
	"fmt"

	"github.com/danielc-lh/scripts/capoeira"
)

func main() {
	// fmt.Println("Starting Bookseller Protocol Example")
	// localTransport := capoeira.NewChannelTransport(capoeira.Seller{}.Name(), capoeira.Buyer{}.Name())
	// httpTransport := capoeira.NewHTTPTransport([]string{capoeira.Seller{}.Name(), capoeira.Buyer{}.Name()})

	// fmt.Println("\n----------------------------------------\n")
	// fmt.Println("Running Bookseller Protocol with Local Channel Transport \n")
	// capoeira.RunBookSellerProtocol("TAPL", localTransport)
	// fmt.Println("\n----------------------------------------\n")
	// fmt.Println("\nRunning Bookseller Protocol with HTTP Transport \n")
	// capoeira.RunBookSellerProtocol("HoTT", httpTransport)

	fmt.Println("Starting Parking Protocol Example")
	transport := capoeira.NewChannelTransport([]string{capoeira.Ticketer{}.Name(), capoeira.ParkingAuthority{}.Name(), capoeira.Printer{}.Name()})
	fmt.Println("\n----------------------------------------")
	fmt.Println("Running Parking Protocol with Local Channel Transport")
	capoeira.RunParkingProtocol(transport)
}
