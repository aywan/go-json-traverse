package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type state int

const (
	unknownState state = iota
	objectState
	arrayState
	valueState
)

type stateStackType struct {
	stack   []state
	current state
}

func (s *stateStackType) Push(new state) {
	s.stack = append(s.stack, s.current)
	s.current = new
}

func (s *stateStackType) Pop() state {
	if len(s.stack) == 0 {
		current := s.current
		s.current = unknownState
		return current
	}
	current := s.current
	s.current = s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return current
}

func (s *stateStackType) Top() state {
	return s.current
}

func main() {
	in, err := os.Open("in.json")
	if err != nil {
		log.Fatal(err)
	}

	stack := stateStackType{make([]state, 0, 100), valueState}

	d := json.NewDecoder(in)
	for {
		token, err := d.Token()
		if err != nil {
			break
		}

		if token == nil {
			if stack.Top() != valueState {
				log.Fatal("not value state")
			}
			fmt.Println("val null")
			stack.Pop()
			if stack.Top() == arrayState {
				stack.Push(valueState)
			}
			continue
		}

		switch token.(type) {
		case json.Delim:
			delim := token.(json.Delim)
			switch delim {
			case '{':
				fmt.Println("object")
				stack.Push(objectState)

			case '[':
				fmt.Println("array")
				stack.Push(arrayState)
				stack.Push(valueState)

			case '}':
				if stack.Pop() != objectState {
					log.Fatal("not object state")
				}
				if stack.Top() == valueState {
					stack.Pop()
				}
				if stack.Top() == arrayState {
					stack.Push(valueState)
				}
				fmt.Println("end object")

			case ']':
				if stack.Pop() != valueState {
					log.Fatal("not value state")
				}
				if stack.Pop() != arrayState {
					log.Fatal("not array state")
				}
				if stack.Top() == valueState {
					stack.Pop()
				}
				fmt.Println("end array")
			}
		case string:
			tokenStr := token.(string)
			switch stack.Top() {
			case objectState:
				fmt.Println("key", tokenStr)
				stack.Push(valueState)

			case valueState:
				fmt.Println("val string", tokenStr)
				stack.Pop()
				if stack.Top() == arrayState {
					stack.Push(valueState)
				}

			default:
				log.Fatal("unknown string state")
			}

		case float64:
			tokenFloat := token.(float64)
			if stack.Top() != valueState {
				log.Fatal("not float64 state")
			}
			fmt.Println("val number", tokenFloat)
			stack.Pop()
			if stack.Top() == arrayState {
				stack.Push(valueState)
			}

		case bool:
			tokenBool := token.(bool)
			if stack.Top() != valueState {
				log.Fatal("not value state")
			}
			fmt.Println("val bool", tokenBool)
			stack.Pop()
			if stack.Top() == arrayState {
				stack.Push(valueState)
			}

		default:
			log.Fatal("unknown type")
		}
	}

}
