package core

import (
	"fmt"
)

type Instruction byte

//can't use 0-9 because those are reserved for operation
const (
	InstrPushInt  Instruction = 0x0a // 10
	InstrAdd      Instruction = 0x0b // 11
	InstrPushByte Instruction = 0x0c //12
	InstrPack     Instruction = 0x0d //13
	InstrSub      Instruction = 0x0e // 14
)

type Stack struct {
	data []any //stack can have anything from int, string, byte to slice of byte
	sp   int   //stack pointer
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

func (s *Stack) Push(v any) {
	s.data[s.sp] = v
	s.sp++
}

//TODO: pop function logic needs to be checked
func (s *Stack) Pop() any {
	// value := s.data[s.sp-1]
	// s.sp--

	value := s.data[0]
	s.data = append(s.data[:0], s.data[1:]...)
	s.sp--

	return value
}

type VM struct {
	data  []byte
	ip    int //instruction counter
	stack *Stack
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: NewStack(128),
	}
}

func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])

		//execute the instruction
		if err := vm.Exec(instr); err != nil {
			return err
		}
		vm.ip++
		fmt.Println(instr)

		//considering each instruction of one byte
		if vm.ip > len(vm.data)-1 {
			break
		}
	}
	return nil
}

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1]))

	case InstrPushByte:
		vm.stack.Push(byte(vm.data[vm.ip-1]))

	case InstrPack:
		n := vm.stack.Pop().(int) //length of pack
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte)
		}
		vm.stack.Push(b)

	case InstrSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a - b)

	case InstrAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a + b)
	}

	return nil
}
