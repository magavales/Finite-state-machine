package pkg

import (
	"gonum.org/v1/gonum/stat/combin"
	"strconv"
)

type FSM struct {
	n        int
	phi      []int
	psi      []int
	graph    *Graph
	PhiTable map[State][]State
	PsiTable map[State][]int
}

// NewFSM
// Создание нового экземпляра класса FSM
// /*
func NewFSM(n int, phi, psi []int) *FSM {
	return &FSM{
		n:        n,
		phi:      phi,
		psi:      psi,
		graph:    NewGraph(),
		PhiTable: make(map[State][]State),
		PsiTable: make(map[State][]int),
	}
}

// GetConnectivityComponents
// Получение компонентов связности автомата
// /*
func (f *FSM) GetConnectivityComponents() [][]State {
	return f.graph.GetConnectivityComponents()
}

// GetStrongConnectivityComponents
// Получение компонентов сильной связности автомата
// /*
func (f *FSM) GetStrongConnectivityComponents() [][]State {
	return f.graph.GetStrongConnectivityComponents()
}

// Init
// Инициализация класса FSM
// /*
func (f *FSM) Init() {
	var (
		zpPhi int
		zpPsi int
	)
	for state := range generateBinaryCombinations(f.n) {
		node := NewState(state)
		for x := range 2 {
			zpPhi = ComputePolinom(x, state, f.phi)
			zpPsi = ComputePolinom(x, state, f.psi)

			newState := NewState(append(state[1:], zpPhi))
			weight := strconv.Itoa(x) + strconv.Itoa(zpPsi)

			f.PhiTable[node] = append(f.PhiTable[node], newState)
			f.PsiTable[node] = append(f.PsiTable[node], zpPsi)

			f.graph.AddEdge(node, newState, weight)
		}
	}
}

// ComputePolinom
// Вычисление полинома Жегалкина
// /*
func ComputePolinom(inputX int, state, psi []int) int {
	zp := []int{1}
	extendedCurrentState := append(state, inputX)

	for i := 1; i < len(extendedCurrentState)+1; i++ {
		gen := combin.NewCombinationGenerator(len(extendedCurrentState), i)
		for gen.Next() {
			var product int = 1
			for _, idx := range gen.Combination(nil) {
				product *= extendedCurrentState[idx]
			}
			zp = append(zp, product)
		}
	}

	if len(zp) != len(psi) {
		panic("Wtf happen! Why zp != coeffs")
	}

	var count = 0
	for i := range len(psi) {
		zp[i] *= psi[i]
		if zp[i] == 1 {
			count++
		}
	}
	return count % 2
}

/*
Генерация значений булевой функции - {0,1}
*/
func generateBinaryCombinations(n int) <-chan []int {
	result := make(chan []int)

	go func() {
		defer close(result)

		var generateCombinationsHelper func([]int, int)
		generateCombinationsHelper = func(currentCombination []int, index int) {
			if index == n {
				combination := make([]int, n)
				copy(combination, currentCombination)
				result <- combination
				return
			}
			currentCombination[index] = 0
			generateCombinationsHelper(currentCombination, index+1)
			currentCombination[index] = 1
			generateCombinationsHelper(currentCombination, index+1)
		}

		initialCombination := make([]int, n)
		generateCombinationsHelper(initialCombination, 0)
	}()

	return result
}
