package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
	//data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x04, 0x0a, 0x0b}
	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e} //3-2=1
	vm := NewVM(data)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop().(int)
	fmt.Printf("answer is %d", result)
	assert.Equal(t, 1, result)
}

func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	value := s.Pop()
	assert.Equal(t, value, 1)

	value = s.Pop()
	assert.Equal(t, value, 2)
}
