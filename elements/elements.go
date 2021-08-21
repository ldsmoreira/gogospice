package elements

type CurrentSource struct {
	label    string
	nodeA    int
	nodeB    int
	waveform string
	device   string
	value    float64
}
type ControledCurrentSource struct {
	label  string
	nodeA  int
	nodeB  int
	nodeC  int
	nodeD  int
	device string
	value  float64
}
type Resistor struct {
	label  string
	nodeA  int
	nodeB  int
	device string
	value  float64
}
