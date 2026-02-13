package capoeira

import "fmt"

// Location represents a participant in a choreography.
type Location interface {
	Name() string
}

// Located represents a value located at a specific location.
type Located struct {
	Value    interface{}
	Location Location
}

// MultiplyLocated represents a value located at multiple locations.
type MultiplyLocated struct {
	Values map[string]interface{} // location name -> value
}

func NewMultiplyLocated() MultiplyLocated {
	return MultiplyLocated{Values: make(map[string]interface{})}
}

func (ml *MultiplyLocated) Add(location Location, value interface{}) {
	ml.Values[location.Name()] = value
}

func (ml *MultiplyLocated) Get(location Location) interface{} {
	return ml.Values[location.Name()]
}

// ChoreoOp provides methods for choreographic operations.
type ChoreoOp interface {
	Locally(location Location, computation func() interface{}) Located
	Comm(sender, receiver Location, data Located) Located
	Broadcast(sender Location, data Located) interface{}
	Multicast(sender Location, destinations []Location, data Located) MultiplyLocated
}

// Choreography is an interface for choreography logic.
type Choreography interface {
	Run(op ChoreoOp) interface{}
}

// Projector performs end-point projection and runs a choreography.
type Projector struct {
	Target    Location
	Transport Transport
}

func NewProjector(target Location, transport Transport) *Projector {
	return &Projector{
		Target:    target,
		Transport: transport,
	}
}

func (p *Projector) Local(value interface{}) Located {
	return Located{Value: value, Location: p.Target}
}

func (p *Projector) Remote(location Location) Located {
	return Located{Value: nil, Location: location}
}

// ProjectorChoreoOp implements ChoreoOp for a specific target and transport.
type ProjectorChoreoOp struct {
	Target    Location
	Transport Transport
}

func (op ProjectorChoreoOp) Locally(location Location, computation func() interface{}) Located {
	if location.Name() == op.Target.Name() {
		return Located{Value: computation(), Location: location}
	}
	return Located{Value: nil, Location: location}
}

func (op ProjectorChoreoOp) Comm(sender, receiver Location, data Located) Located {
	if sender.Name() == op.Target.Name() && sender.Name() == receiver.Name() {
		return Located{Value: data.Value, Location: receiver}
	}
	if sender.Name() == op.Target.Name() {
		// Send via transport
		if t, ok := op.Transport.(Transport); ok {
			fmt.Printf("Sending from %s to %s. data: %+v\n", sender.Name(), receiver.Name(), data.Value)
			t.Send(sender.Name(), receiver.Name(), data.Value)
		}
		return Located{Value: data.Value, Location: receiver}
	} else if receiver.Name() == op.Target.Name() {
		// Receive via transport
		if t, ok := op.Transport.(Transport); ok {
			fmt.Printf("Receiving from %s at %s\n", sender.Name(), receiver.Name())
			val := t.Receive(sender.Name(), receiver.Name())
			fmt.Printf("Received val: %+v\n", val)
			return Located{Value: val, Location: receiver}
		}
		return Located{Value: nil, Location: receiver}
	}
	return Located{Value: nil, Location: receiver}
}

func (op ProjectorChoreoOp) Broadcast(sender Location, data Located) interface{} {
	if sender.Name() == op.Target.Name() {
		if t, ok := op.Transport.(Transport); ok {
			for _, dest := range t.Locations() {
				if dest != sender.Name() {
					t.Send(sender.Name(), dest, data.Value)
				}
			}
		}
		return data.Value
	}
	if t, ok := op.Transport.(Transport); ok {
		return t.Receive(sender.Name(), op.Target.Name())
	}
	return data.Value
}

func (op ProjectorChoreoOp) Multicast(sender Location, destinations []Location, data Located) MultiplyLocated {
	ml := NewMultiplyLocated()
	if sender.Name() == op.Target.Name() {
		if t, ok := op.Transport.(Transport); ok {
			for _, dest := range destinations {
				if dest.Name() != sender.Name() {
					t.Send(sender.Name(), dest.Name(), data.Value)
				}
			}
		}
		for _, dest := range destinations {
			ml.Add(dest, data.Value)
		}
	} else {
		if t, ok := op.Transport.(Transport); ok {
			for _, dest := range destinations {
				if dest.Name() == op.Target.Name() {
					val := t.Receive(sender.Name(), dest.Name())
					ml.Add(dest, val)
				} else {
					ml.Add(dest, nil)
				}
			}
		}
	}
	return ml
}

// EppAndRun performs end-point projection to run a choreography for the target location.
func (p *Projector) EppAndRun(choreo Choreography) interface{} {
	op := ProjectorChoreoOp{
		Target:    p.Target,
		Transport: p.Transport,
	}
	return choreo.Run(op)
}
