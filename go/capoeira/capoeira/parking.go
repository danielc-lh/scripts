package capoeira

import (
	"fmt"
	"sync"
	"time"
)

type Garage struct {
	spaces []ParkingSpace
}

type ParkingSpace struct {
	number    int
	occupied  bool
	occupant  string
	startTime time.Time
	duration  time.Duration
}

// Locations
type Ticketer struct{}

func (t Ticketer) Name() string { return "ticketer" }

type ParkingAuthority struct{}

func (t ParkingAuthority) Name() string { return "parking_authority" }

type Printer struct{}

func (t Printer) Name() string { return "printer" }

// choreography
type TicketingChoreography struct{}

func getGarageState() interface{} {
	return Garage{
		spaces: []ParkingSpace{
			{number: 1, occupied: false}, // empty
			{number: 2, occupied: true, occupant: "Alice", startTime: time.Now(), duration: time.Hour},                       // full + paid
			{number: 3, occupied: true, occupant: "Bob", startTime: time.Now().Add(-2 * time.Hour), duration: 1 * time.Hour}, // overdue
		},
	}
}

// run

func (t TicketingChoreography) Run(op ChoreoOp) interface{} {
	garageAtTicketer := op.Locally(Ticketer{}, func() interface{} {
		return getGarageState()
	})

	// this is the key!! everyone needs to know about the garage
	// since we need to know how many spots we will need to try and make decisions for,
	// send decisions for, and receive decisions for.
	garage := op.Broadcast(Ticketer{}, garageAtTicketer).(Garage)

	for _, space := range garage.spaces {
		// Check if the space is occupied and if the duration has expired
		decisionAtPA := op.Locally(
			ParkingAuthority{}, func() interface{} {
				decisionAtPA := space.occupied && space.startTime.Add(space.duration).Before(time.Now())
				return decisionAtPA
			})
		decision := op.Broadcast(ParkingAuthority{}, decisionAtPA).(bool)
		fmt.Printf("Space %d decision: %v\n", space.number, decision)

		if decision {
			// the space is expired, so send it to the printer
			spaceAtTicketer := op.Comm(ParkingAuthority{}, Printer{}, Located{Value: space, Location: ParkingAuthority{}})
			return op.Locally(Printer{}, func() interface{} {
				printer := Printer{}
				if spaceAtTicketer.Location != printer {
					fmt.Printf("wrong location! Got %v\n", spaceAtTicketer.Location)
				}
				space, ok := spaceAtTicketer.Value.(ParkingSpace)
				if !ok {
					fmt.Printf("failed to cast to ParkingSpace\n")
					return nil
				}
				fmt.Printf("Printing ticket for space %d occupied by %s\n", space.number, space.occupant)
				return space
			})
		}
	}
	return nil
}

// creates transports, projectors, and runs each endpoint
// returns a channel of parking spaces that should be ticketed (i.e. the expired ones)
func RunParkingProtocol(transport Transport) chan ParkingSpace {
	var wg sync.WaitGroup
	wg.Add(3)

	toTicket := make(chan ParkingSpace, 1)

	// Ticketer endpoint
	go func() {
		defer wg.Done()
		ticketerProjector := NewProjector(Ticketer{}, transport)
		ticketerProjector.EppAndRun(TicketingChoreography{})
	}()

	// ParkingAuthority endpoint
	go func() {
		defer wg.Done()
		printerProjector := NewProjector(ParkingAuthority{}, transport)
		printerProjector.EppAndRun(
			TicketingChoreography{},
		)
	}()

	// Printer endpoint
	go func() {
		defer wg.Done()
		printerProjector := NewProjector(Printer{}, transport)
		space := printerProjector.EppAndRun(
			TicketingChoreography{},
		)
		// Send the space to the result channel
		toTicket <- space.(Located).Value.(ParkingSpace)
	}()

	wg.Wait()
	return toTicket
}
