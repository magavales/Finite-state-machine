package pkg

import (
	"gonum.org/v1/gonum/stat/combin"
	"strconv"
)

type FSM struct {
	n        int               // длина регистра
	phi      []int             // функция фи
	psi      []int             // функция пси
	graph    *Graph            // граф автомата
	PhiTable map[State][]State // таблица состояний, полученных из функции фи
	PsiTable map[State][]int   // таблица состояний, полученных из функции пси
	Delta    int
	Mu       int
}

// Создание нового экземпляра класса FSM. На вход методу подается длина регистра, функции фи и пси. Метод взвращает
// экземпляр класса FSM
// Входные данные: n int, phi, psi []int
// Выходные данные: *FSM
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

// Получение компонентов связности автомата
// Входные данные: null
// Выходные данные: [][]State
func (f *FSM) GetConnectivityComponents() [][]State {
	return f.graph.GetConnectivityComponents()
}

// Получение компонентов сильной связности автомата
// Входные данные: null
// Выходные данные: [][]State
func (f *FSM) GetStrongConnectivityComponents() [][]State {
	return f.graph.GetStrongConnectivityComponents()
}

// Инициализация класса FSM. В рамках данной функции создаются граф, таблица фи и пси, где находятся состояния
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

// Вычисление полинома Жегалкина. На выход функции подается значение inputX, которое будет иметь значения {0,1}, состояение
// state и функция psi.
// Входные данные: inputX int, state, psi []int
// Выходные данные: int
func ComputePolinom(inputX int, state, psi []int) int {
	zp := []int{1}
	extendedCurrentState := append(state, inputX)

	for i := 1; i < len(extendedCurrentState)+1; i++ {
		gen := combin.NewCombinationGenerator(len(extendedCurrentState), i)
		for gen.Next() {
			var product = 1
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

// Данная функция генерирует значения {0,1}. Генерация выполняется n раз, где n - длина регистра автомата
// Входные данные: n
// Выходные данные: []int = {0,1}
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
