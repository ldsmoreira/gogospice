package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type Element struct {
	label    string
	nodeA    string
	nodeB    string
	nodeC    string
	nodeD    string
	waveform string
	device   string
	value    int
}

func main() {
	content, err := ioutil.ReadFile("netlist2.net")
	if err != nil {
		log.Fatal(err)
	}

	var elementArray []Element
	const nElementsConst = 10
	slices := strings.Split(string(content), "\n")

	var slicenet [7]string
	var slicematrix [nElementsConst][len(slicenet)]string

	for index, value := range slices {
		for index2, value2 := range strings.Split(value, " ") {
			slicematrix[index][index2] = value2
		}
	}

	for _, value := range slicematrix {

		if value[0] == "" {
			continue
		}

		switch string(value[0][0]) {

		case "I":

			current, _ := strconv.Atoi(value[4])
			element := Element{label: value[0], nodeA: value[1], nodeB: value[2],
				waveform: value[3], device: "current_source", value: current}

			elementArray = append(elementArray, element)

		case "R":

			resistence, _ := strconv.Atoi(value[3])
			element := Element{label: value[0], nodeA: value[1], nodeB: value[2],
				device: "resistor", value: resistence}
			elementArray = append(elementArray, element)

		case "G":

			transcondutancy, _ := strconv.Atoi(value[5])
			element := Element{label: value[0], nodeA: value[1],
				nodeB: value[2], nodeC: value[3],
				nodeD: value[4], device: "controled_current_source",
				value: transcondutancy}

			elementArray = append(elementArray, element)

		default:

			fmt.Printf("%s.\n", "Not found")
		}

	}

	var matrix_dimension int = 0

	for _, value := range elementArray {

		nodeA, _ := strconv.Atoi(value.nodeA)
		nodeB, _ := strconv.Atoi(value.nodeB)
		nodeC, _ := strconv.Atoi(value.nodeC)
		nodeD, _ := strconv.Atoi(value.nodeD)

		if nodeA >= matrix_dimension && value.nodeA != "" {
			matrix_dimension = nodeA
		}
		if nodeB >= matrix_dimension && value.nodeB != "" {
			matrix_dimension = nodeB
		}
		if nodeC >= matrix_dimension && value.nodeC != "" {
			matrix_dimension = nodeC
		}
		if nodeD >= matrix_dimension && value.nodeD != "" {
			matrix_dimension = nodeD
		}
	}

	matrix_dimension = matrix_dimension + 1

	condutance_matrix := mat.NewDense(matrix_dimension, matrix_dimension, nil)
	current_vector := mat.NewDense(1, matrix_dimension, nil)

	for _, element := range elementArray {

		switch element.device {

		case "current_source":

			nodeA, _ := strconv.Atoi(element.nodeA)
			nodeB, _ := strconv.Atoi(element.nodeB)
			current_vector.Set(0, nodeA, current_vector.At(0, nodeA)-float64(element.value))
			current_vector.Set(0, nodeB, current_vector.At(0, nodeB)+float64(element.value))

		case "resistor":

			nodeA, _ := strconv.Atoi(element.nodeA)
			nodeB, _ := strconv.Atoi(element.nodeB)
			condutance_matrix.Set(nodeA, nodeA, condutance_matrix.At(nodeA, nodeA)+1/float64(element.value))
			condutance_matrix.Set(nodeB, nodeB, condutance_matrix.At(nodeB, nodeB)+1/float64(element.value))
			condutance_matrix.Set(nodeA, nodeB, condutance_matrix.At(nodeA, nodeB)-1/float64(element.value))
			condutance_matrix.Set(nodeB, nodeA, condutance_matrix.At(nodeB, nodeA)-1/float64(element.value))

		case "controled_current_source":
			nodeA, _ := strconv.Atoi(element.nodeA)
			nodeB, _ := strconv.Atoi(element.nodeB)
			nodeC, _ := strconv.Atoi(element.nodeC)
			nodeD, _ := strconv.Atoi(element.nodeD)
			condutance_matrix.Set(nodeA, nodeC, condutance_matrix.At(nodeA, nodeC)+float64(element.value))
			condutance_matrix.Set(nodeB, nodeD, condutance_matrix.At(nodeB, nodeD)+float64(element.value))
			condutance_matrix.Set(nodeB, nodeC, condutance_matrix.At(nodeB, nodeC)-float64(element.value))
			condutance_matrix.Set(nodeA, nodeD, condutance_matrix.At(nodeA, nodeD)-float64(element.value))
		}
	}

	// var condutance_matrix_sized mat.Dense

	condutance_matrix_sized := condutance_matrix.Slice(1, matrix_dimension, 1, matrix_dimension)
	current_vector_sized := current_vector.Slice(0, 1, 1, matrix_dimension)

	var condutance_matrix_sized_INV mat.Dense
	err = condutance_matrix_sized_INV.Inverse(condutance_matrix_sized)
	if err != nil {
		log.Fatalf("condutance_matrix is not invertible: %v", err)
	}

	var node_tension_vector mat.Dense
	node_tension_vector.Mul(&condutance_matrix_sized_INV, current_vector_sized.T())

	fc := mat.Formatted(&node_tension_vector, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fc)
}
