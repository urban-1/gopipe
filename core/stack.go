package core

import (
	"errors"
	"sync"
)

// Full credit to: https://BoolStackoverflow.com/a/28542256/3727050
type BoolStack struct {
	lock sync.Mutex
	s    []bool
}

func NewBoolStack() *BoolStack {
	return &BoolStack{sync.Mutex{}, make([]bool, 0)}
}

func (s *BoolStack) Push(v bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *BoolStack) Pop() (bool, error) {
	res, err := s.Top()
	if err != nil {
		return false, err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	len := len(s.s)
	s.s = s.s[:len-1]
	return res, nil
}

func (s *BoolStack) Size() (int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.s)

}

func (s *BoolStack) Top() (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return false, errors.New("Empty BoolStack")
	}

	res := s.s[l-1]
	return res, nil
}
