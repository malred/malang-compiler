package compiler

import (
	"fmt"
	"malang/ast"
	"malang/code"
	"malang/lexer"
	"malang/object"
	"malang/parser"
	"testing"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

// expectedInstructions是[][]byte，需要转为[]byte
func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

func testIntegerObject(expected int64, actual object.Object) error {
	res, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if res.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", res.Value, expected)
	}

	return nil
}

func testConstants(
	t *testing.T,
	expected []interface{},
	actual []object.Object,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want %d", len(actual), len(expected))
	}

	// 遍历比较常量池
	for i, constant := range expected {
		// 根据不同类型，进行测试
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func testInstructions(
	expected []code.Instructions,
	actual code.Instructions,
) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot=%q", i, concatted, actual)
		}
	}

	return nil
}

func TestIntegerArithmetic(t *testing.T) {
	ts := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
			},
		},
	}

	runCompilerTests(t, ts)
}

func runCompilerTests(t *testing.T, ts []compilerTestCase) {
	t.Helper()

	// 编译器生成字节码
	for _, tt := range ts {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compiler(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		// 将传给虚拟机的指令
		bytecode := compiler.Bytecode()

		// 测试字节码是否正确
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		// 测试常量池是否正确
		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
	}
}
