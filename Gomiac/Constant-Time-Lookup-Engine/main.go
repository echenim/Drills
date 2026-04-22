package main

import "fmt"

type FastLookupEngine struct {
	data  []string
	index map[string]int
}

// new function to initialize the data structure
func NewFastLookupEngine() *FastLookupEngine {
	return &FastLookupEngine{
		data:  []string{},
		index: make(map[string]int),
	}
}

func (e *FastLookupEngine) Put(value string) {
	_, exist := e.index[value]
	if !exist {
		e.data = append(e.data, value)
		e.index[value] = len(e.data) - 1
	}
}

func (e *FastLookupEngine) Check(value string) {
	if _, exist := e.index[value]; !exist {
		fmt.Printf("No, %v does not exist", value)
		return
	}
	fmt.Printf("Yes, %v exist", value)
}

func (e *FastLookupEngine) Remove(value string) {
	if k, exist := e.index[value]; exist {
		// Swap the element to remove with the last element
		lastindex := len(e.data) - 1
		e.data[k], e.data[lastindex] = e.data[lastindex], e.data[k]
		e.data = e.data[:lastindex]
		delete(e.index, value)
		if k != lastindex {
			e.index[e.data[k]] = k
		}

	}
}

func main() {
	engine := NewFastLookupEngine()
	engine.Put("apple")
	engine.Put("banana")
	engine.Put("cherry")

	fmt.Printf("\n Initial state: %v \n", engine.data)

	engine.Check("banana")
	engine.Remove("banana")
	fmt.Printf("\n current state: %v \n", engine.data)
	engine.Check("banana")
	engine.Check("apple")
	engine.Check("cherry")
	engine.Check("banana")
	fmt.Printf("\n current state: %v \n", engine.data)
}
