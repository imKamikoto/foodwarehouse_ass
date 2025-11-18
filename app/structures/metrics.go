package structures

type Metrics struct {
	Received  int
	Discarded int
	Delivered int
}

func (m *Metrics) LogArrival(b Batch) {
	m.Received++
}

func (m *Metrics) LogDiscard(b Batch) {
	m.Discarded++
}

func (m *Metrics) LogDelivery(b Batch) {
	m.Delivered++
}
