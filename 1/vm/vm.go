package vm

import (
	"fmt"
	"malang/code"
	"malang/compiler"
	"malang/object"
)

const StackSize = 2048

type VM struct {
	// compiler生成的常量和指令
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // 始终指向栈中的下一个空闲槽。栈顶的值是stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// 获取栈顶元素
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	// 取指
	for ip := 0; ip < len(vm.instructions); ip++ {
		// 字节转换为操作码
		op := code.Opcode(vm.instructions[ip])

		// 解码
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			res := leftValue + rightValue
			// 结果压栈
			vm.push(&object.Integer{Value: res})
		}
	}
	return nil
}

// 元素压栈
func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

// 元素弹栈
func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}
