package varis

import "sync"

type synapse struct {
	weight    float64
	in        chan float64
	out       chan float64
	cache     float64
	inNeuron  *Neuron
	outNeuron *Neuron
}

func (syn *synapse) live() {
	for {
		syn.cache = <-syn.in
		outputValue := syn.cache * syn.weight
		syn.out <- outputValue
	}
}

func ConnectNeurons(in *Neuron, out *Neuron, weight float64) {
	syn := &synapse{
		weight:    weight,
		in:        make(chan float64),
		out:       make(chan float64),
		inNeuron:  in,
		outNeuron: out,
	}

	in.conn.addOutputSynapse(syn)
	out.conn.addInputSynapse(syn)

	go syn.live()
}

type connection struct {
	inSynapses  []*synapse
	outSynapses []*synapse
}

func (c *connection) addOutputSynapse(syn *synapse) {
	c.outSynapses = append(c.outSynapses, syn)
}

func (c *connection) addInputSynapse(syn *synapse) {
	c.inSynapses = append(c.inSynapses, syn)
}

func (c *connection) collectSignals() []float64 {

	inputCount := len(c.inSynapses)
	inputSignals := make([]float64, inputCount)

	wg := sync.WaitGroup{}
	wg.Add(inputCount)

	for i := range inputSignals {
		go func(index int) {
			inputSignals[index] = <-c.inSynapses[index].out
			wg.Done()
		}(i)
	}

	wg.Wait()
	return inputSignals
}

func (c *connection) broadcastSignals(value float64) {
	for _, o := range c.outSynapses {
		o.in <- value
	}
}

func (c *connection) changeWeight(delta float64) {
	for _, s := range c.inSynapses {
		s.weight += s.cache * delta
	}
}
