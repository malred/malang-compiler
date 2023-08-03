// repl/repl.go
package repl
/*
#include <stdio.h>

void pf(){
	printf("hello malang!\n");
}
*/
import "C"
import (
	"bufio"
	"fmt"
	"io"
	"malang/compiler"
	"malang/evaluator"
	"malang/lexer"
	"malang/object"
	"malang/parser"
	"malang/util"
	"malang/vm"
)

const PROMPT = ">> "

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`
const MALRED_LOGO_IMG = `
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO\OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO  ,OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO     [OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO        \OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO          ,OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO             [OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO/\OOOOOOOOOOOOOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^  ,OOOOOOOOOOOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^     ,OOOOOOOOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^        \OOOOOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^          ,OOOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^           =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^           =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^           =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^ [OOOOOOOOO^           =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^    \OOOOOO^           =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^      ,OOOO^           =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^       =OOO^           =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^       =OOO^           =OOOO               OOOO                /OOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^       =OOO^       ]\  =OOOO               OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^       =OOO^    ]O  \\/OOOO              ]O/.               ,OOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^       =OOO^,/O[[[[[[[[OOOOO          ]O/[[[[[[[[[[[[[[[[[[[OOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^       =OOOO\]]]]]]]]]]]/OOO      ]OO]]]]]]]]]]]]]]]]]]]]]]OOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^      ]O[                [O ,/O[[                  .    =OOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO^   ]OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO\//                                         ,/O/   OO OOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO      ,OOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOO^             ,/OOOOOOO^              ,/OOOOOOOO            OOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^           OOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOO       ]OOOOOOOOOOOOOOOOOO  ]OOOOOOOOOOOOOOOOO^           OOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO        ,]/OOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOO[   OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^      OOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOO      =OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO      =OOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOO^      =OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO       OOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOO/         ,OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO/         ,OOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOO/             \OOOOOOOOOOOOOOOOOOOOOOOOOOOOOO            ,OOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOO^              \OOOOOOOOOOOOOOOOOOOOOOOOOO              /OOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOO^   ,]/        ,\OOOOOOOOOOOOOOOOOOOO[         /]]    OOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO          ,[OOOOOOOOOOO/[           /OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^                                  OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO                                   =OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO^        ]]               ,]]        OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO\    ,OOOOOOO^        OOOOOOO\    ,/OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO      =OOOOOOOOO\OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO\     ,OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO`
const MALRED_LOGO = `
 __  __    __    __    ____  ____  ____  
(  \/  )  /__\  (  )  (  _ \( ___)(  _ \ 
 )    (  /(__)\  )(__  )   / )__)  )(_) )
(_/\/\_)(__)(__)(____)(_)\_)(____)(____/ 

`
const ERROR_LOGO = `
 ____  ____  ____  _____  ____ 
( ___)(  _ \(  _ \(  _  )(  _ \
 )__)  )   / )   / )(_)(  )   /
(____)(_)\_)(_)\_)(_____)(_)\_)

`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	io.WriteString(out, MALRED_LOGO)
	// 加载标准库
	std := util.LoadStd()
	l := lexer.New(std)
	p := parser.New(l)
	program := p.ParseProgram()
	evaluator.Eval(program, env)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func StartVM(in io.Reader, out io.Writer) {
	C.pf() 
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// 求值不同
		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithState(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MALRED_LOGO_IMG)
	io.WriteString(out, ERROR_LOGO)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func ReadAndEval(input string) {
	std := util.LoadStd()
	l := lexer.New(std + input)
	p := parser.New(l)
	env := object.NewEnvironment()
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		fmt.Println(MALRED_LOGO_IMG)
		fmt.Println(ERROR_LOGO)
		fmt.Println("Woops! We ran into some monkey business here!\n")
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg + "\n")
		}
	}
	evaluator.Eval(program, env)
	// fmt.Printf(">> %v\n", evaluator.Eval(pro, env).Inspect())
}
