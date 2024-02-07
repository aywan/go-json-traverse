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
	stack []state
}

func (s *stateStackType) Push(new state) {
	s.stack = append(s.stack, new)
}

func (s *stateStackType) Pop() state {
	current := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return current
}

func main() {
	in, err := os.Open("in.json")
	if err != nil {
		log.Fatal(err)
	}

	stateStack := stateStackType{make([]state, 0, 100)}
	curState := valueState

	d := json.NewDecoder(in)
	for {
		token, err := d.Token()
		if err != nil {
			break
		}

		switch token.(type) {
		case json.Delim:
			delim := token.(json.Delim)
			switch delim {
			case '{':
				stateStack.Push(curState)
				curState = objectState

			case '[':
				stateStack.Push(curState)
				stateStack.Push(arrayState)
				curState = valueState

			case '}':
				if curState != objectState {
					log.Fatal("not object state")
				}
				curState = stateStack.Pop()
				if curState == valueState {
					curState = stateStack.Pop()
				}

			case ']':
				if curState != valueState {
					log.Fatal("not value state")
				}
				curState = stateStack.Pop()
				if curState != arrayState {
					log.Fatal("not array state")
				}
				curState = stateStack.Pop()
				if curState == valueState {
					curState = stateStack.Pop()
				}
			}
		case string:
			tokenStr := token.(string)
			switch curState {
			case objectState:
				fmt.Println("key", tokenStr)
				stateStack.Push(curState)
				curState = valueState

			case valueState:
				fmt.Println("val string", tokenStr)
				curState = stateStack.Pop()
				if curState == arrayState {
					stateStack.Push(arrayState)
					curState = valueState
				}

			default:
				log.Fatal("unknown string state")
			}

		case float64:
			tokenFloat := token.(float64)
			if curState != valueState {
				log.Fatal("not float64 state")
			}
			fmt.Println("flo string", tokenFloat)
			curState = stateStack.Pop()
			if curState == arrayState {
				stateStack.Push(arrayState)
				curState = valueState
			}
		}
	}

}
