package capoeira

import (
	"fmt"
	"testing"
	"time"
)

func TestParkingChoreography(t *testing.T) {

	fmt.Println("Starting Parking Protocol Example")
	transport := NewChannelTransport([]string{Ticketer{}.Name(), ParkingAuthority{}.Name(), Printer{}.Name()})
	fmt.Println("\n----------------------------------------")
	fmt.Println("Running Parking Protocol with Local Channel Transport")
	ticketChan := RunParkingProtocol(transport)

	// Wait for the ticket channel to receive a parking space
	select {
	case space := <-ticketChan:
		fmt.Printf("Received ticket for space %d\n", space.number)
		if space.number != 3 {
			t.Errorf("Expected ticket for space 3 but got %d", space.number)
		}
		close(ticketChan)
	case <-time.After(5 * time.Second):
		close(ticketChan)
		t.Error("Expected ticket but got none")
	}
}
