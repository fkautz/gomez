package gomez

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strconv"
	"strings"
)

func GenerateLLVM(fset *token.FileSet, tree ast.Node) (string, error) {
	g := &gomezGenerator{
		importTable:   &SymbolTable{},
		symbolTable:   &SymbolTable{},
		functionTable: &SymbolTable{},
		output:        bytes.NewBufferString(""),
	}
	g.fset = fset
	a, b := g.walk(tree)
	fmt.Println("complete")
	fmt.Println(a)
	fmt.Println(b)
	return g.output.String(), nil
}

type gomezGenerator struct {
	fset          *token.FileSet
	output        *bytes.Buffer
	importTable   *SymbolTable
	symbolTable   *SymbolTable
	functionTable *SymbolTable
	functionType  string
	counter       int
}

func (g *gomezGenerator) walk(node ast.Node) (string, error) {
	if node != nil {
		log.Println("  start: ", node.Pos())
		log.Println("  end: ", node.End())
	}
	switch typedNode := node.(type) {
	case *ast.File:
		{
			log.Println("package name: ", typedNode.Name)
			ast.Print(g.fset, typedNode.Name)
			// Initialize symbol tables
			g.importTable.PushFrame()
			g.symbolTable.PushFrame()
			g.functionTable.PushFrame()

			// TODO Populate imports
			for name, obj := range typedNode.Scope.Objects {
				if obj.Kind == ast.Pkg {
					// TODO add type to symbol
					g.importTable.AddSymbol(name, nil, "")
				}
			}

			// Populate variable table
			for name, obj := range typedNode.Scope.Objects {
				if obj.Kind == ast.Var {
					// TODO add type to symbol
					fmt.Println("FOUND GLOBAL VARIABLE")
					g.symbolTable.AddSymbol(name, []string{"int32"}, "@"+name)
				}
			}

			// Populate function table
			for name, obj := range typedNode.Scope.Objects {
				if obj.Kind == ast.Fun {
					// TODO add type to function symbol
					g.functionTable.AddSymbol(name, nil, "")
					fmt.Println("In Symbol Table: ")
					for k, v := range g.symbolTable.frames[0].variables {
						fmt.Println(k, v)
					}
				}
			}

			// Emit functions
			for _, decl := range typedNode.Decls {
				switch typedDecl := decl.(type) {
				case *ast.FuncDecl:
					{
						g.counter = 0
						g.importTable.PushFrame()
						g.symbolTable.PushFrame()

						if _, err := g.walk(decl); err != nil {
							return "", err
						}

						g.importTable.PopFrame()
						g.symbolTable.PopFrame()
					}
				case *ast.GenDecl:
					{
						result, err := g.walk(typedDecl)
						if err != nil {
							return "", err
						}
						result = strings.Replace(result, "=", "= global", 1)
						g.output.WriteString("@" + result + "\n")
					}
				}
			}
			g.importTable.PopFrame()
			g.symbolTable.PopFrame()
			g.functionTable.PopFrame()
			return "", nil
		}
	case *ast.Ident:
		{
			switch typedNode.Obj.Kind {
			case ast.Fun:
				{
					return "@" + typedNode.Name, nil
				}
			case ast.Var:
				{
					variable := "%" + strconv.Itoa(g.counter)
					g.counter = g.counter + 1
					_, _, internalName, _ := g.symbolTable.FindVariable(typedNode.Name)
					// TODO handle not ok
					g.output.WriteString("  " + variable + " = load i32* " + internalName + ", align 4\n")
					return variable, nil
				}
			default:
				{
					return "", errors.New("Ident not implemented")
				}
			}
		}
	case *ast.GenDecl:
		{
			name := ""
			value := ""
			var err error
			fmt.Println("GenDecl: ", typedNode)
			switch typedSpec := typedNode.Specs[0].(type) {
			case *ast.ValueSpec:
				{
					name = typedSpec.Names[0].Name
					value, err = g.walk(typedSpec.Values[0])
					if err != nil {
						return "", err
					}
				}
			}
			fmt.Println("var name: " + name)
			//			g.variables.AddVariable("")
			return name + " = i32 " + value, nil
		}
	case *ast.FuncDecl:
		{
			functionName := typedNode.Name.Name

			if typedNode.Type.Results != nil && typedNode.Type.Results.List != nil {
				switch len(typedNode.Type.Results.List) {
				case 0:
					g.functionType = "void"
				case 1:
					g.functionType = "i32"
				default:
					{
						names := make([]string, 0)
						for _, field := range typedNode.Type.Results.List {
							for _, name := range field.Names {
								names = append(names, name.Name)
							}
						}
						g.functionType = ""
						g.functionType = "i32 " + strings.Join(names, ", i32 ")
					}
				}
			} else {
				g.functionType = "void"
			}

			log.Println("Function Name: " + functionName)
			g.output.WriteString("define " + g.functionType + " @" + functionName + "(")
			fmt.Println("recv")
			fmt.Println(typedNode.Type.Params)
			fmt.Println(typedNode)
			allocations := make([]string, 0)
			if len(typedNode.Type.Params.List) > 0 {
				for i, input := range typedNode.Type.Params.List {
					fmt.Println("i", i)
					for j, recv := range input.Names {
						fmt.Println("j", j)
						if !(i == 0 && j == 0) {
							g.output.WriteString(", ")
						}
						allocations = append(allocations, recv.Name)
						g.output.WriteString("i32 %_" + recv.Name)
					}
				}
			}

			g.output.WriteString(") {\nentry:\n")
			for _, alloc := range allocations {
				g.output.WriteString("  %" + alloc + " = alloca i32, align 4\n")
			}

			for _, alloc := range allocations {
				// TODO handle not ok
				g.output.WriteString("  store i32 %_" + alloc + ", i32* %" + alloc + ", align 4\n")
				g.symbolTable.AddSymbol(alloc, nil, "%"+alloc)
			}
			g.walk(typedNode.Body)
			if g.functionType == "void" {
				g.output.WriteString("  ret void\n")
			}
			g.functionType = ""
			g.output.WriteString("}\n\n")
			return "", nil
		}
	case *ast.BlockStmt:
		{
			for _, statement := range typedNode.List {
				fmt.Println(statement)
				g.walk(statement)
			}
			return "", nil
		}
	case *ast.ReturnStmt:
		{
			result := "void"
			var err error
			if len(typedNode.Results) != 0 {
				result, err = g.walk(typedNode.Results[0])
				if err != nil {
					return "", err
				}
				g.output.WriteString("  ret i32 " + result + "\n")
			} else {
				g.output.WriteString("  ret void\n")
			}
			return "", nil
		}
	case *ast.AssignStmt:
		{
			fmt.Println("Assign Statement")
			fmt.Println(typedNode)
			left := typedNode.Lhs[0].(*ast.Ident).Name
			right, err := g.walk(typedNode.Rhs[0])
			if err != nil {
				return "", err
			}
			_, _, internalName, ok := g.symbolTable.FindVariable(left)
			if ok != nil {
				left = g.symbolTable.AddSymbol(left, nil, "%"+left)
				g.output.WriteString("  %" + left + " = alloca i32, align 4\n")
				internalName = "%" + left
			}
			//			g.output.WriteString(";" + left + "\n")
			g.output.WriteString("  store i32 " + right + ", i32* " + internalName + ", align 4\n")
			return "", nil
		}
	case *ast.ExprStmt:
		{
			return g.walk(typedNode.X)
		}
	case *ast.BasicLit:
		{
			return typedNode.Value, nil
		}
	case *ast.BinaryExpr:
		{
			left, err := g.walk(typedNode.X)
			if err != nil {
				return "", err
			}
			right, err := g.walk(typedNode.Y)
			if err != nil {
				return "", err
			}
			op := typedNode.Op.String()
			result := "%" + strconv.Itoa(g.counter)
			g.counter++
			switch op {
			case "+":
				{
					g.output.WriteString("  " + result + " = add nsw i32 " + left + ", " + right + "\n")
				}
			case "-":
				{
					g.output.WriteString("  " + result + " = sub nsw i32 " + left + ", " + right + "\n")
				}
			case "*":
				{
					g.output.WriteString("  " + result + " = mul nsw i32 " + left + ", " + right + "\n")
				}
			case "/":
				{
					g.output.WriteString("  " + result + " = sdiv i32 " + left + ", " + right + "\n")
				}
			case "%":
				{
					g.output.WriteString("  " + result + " = srem i32 " + left + ", " + right + "\n")
				}
			case "<":
				{
					g.output.WriteString("  " + result + " = icmp slt i32 " + left + ", " + right + "\n")
				}
			case ">":
				{
					g.output.WriteString("  " + result + " = icmp sgt i32 " + right + ", " + left + "\n")
				}
			}
			return result, nil
		}
	case *ast.CallExpr:
		{
			funcName, err := g.walk(typedNode.Fun)
			fmt.Println("calling: " + funcName)
			if err != nil {
				return "", err
			}
			args := make([]string, 0)
			for _, expr := range typedNode.Args {
				result, err := g.walk(expr)
				if err != nil {
					return "", err
				}
				args = append(args, result)
			}
			for i, s := range args {
				args[i] = "i32 " + s
			}
			callArgs := strings.Join(args, ", ")
			result := "%" + strconv.Itoa(g.counter)
			g.counter++
			g.output.WriteString("  " + result + " = call i32 " + funcName + "(" + callArgs + ") \n")
			return result, nil
		}
	case *ast.IfStmt:
		{
			condition, err := g.walk(typedNode.Cond)
			if err != nil {
				return "", err
			}

			output := g.output

			g.output = bytes.NewBufferString("")

			// true
			trueLabel := strconv.Itoa(g.counter)
			g.counter++
			if _, err = g.walk(typedNode.Body); err != nil {
				return "", err
			}

			trueOutput := g.output

			g.output = bytes.NewBufferString("")
			falseLabel := strconv.Itoa(g.counter)
			g.counter++
			if typedNode.Else != nil {
				if _, err = g.walk(typedNode.Else); err != nil {
					return "", err
				}
			}

			falseOutput := g.output

			joinLabel := strconv.Itoa(g.counter)
			g.counter++

			g.output = output
			g.output.WriteString("  br i1 " + condition + ", label %" + trueLabel + ", label %" + falseLabel + "\n")

			// true branch
			g.output.WriteString("; <label>:" + trueLabel + "\n")
			g.output.WriteString(trueOutput.String())
			g.output.WriteString("  br label %" + joinLabel + "\n")

			// false branch
			g.output.WriteString("; <label>:" + falseLabel + "\n")
			g.output.WriteString(falseOutput.String())
			g.output.WriteString("  br label %" + joinLabel + "\n")

			// join label
			g.output.WriteString("; <label>:" + joinLabel + "\n")
			return "", nil
		}
	case *ast.ForStmt:
		{
			g.output.WriteString("; for statement\n")
			// create block
			g.symbolTable.PushFrame()
			// initialize
			g.output.WriteString("; init statement\n")
			init, err := g.walk(typedNode.Init)
			if err != nil {
				return "", err
			}
			g.output.WriteString(init)

			// label condition
			conditionLabel := strconv.Itoa(g.counter)
			g.counter++
			g.output.WriteString("; condition statement\n")
			g.output.WriteString("  br label %" + conditionLabel + "\n")
			g.output.WriteString("; <label>:" + conditionLabel + "\n")

			// if condition is false, jump to label end
			condition, err := g.walk(typedNode.Cond)
			if err != nil {
				return "", err
			}
			// TODO REMOVE PLACEHOLDERS

			// swap writer
			mainWriter := g.output
			g.output = bytes.NewBufferString("")

			// body
			g.output.WriteString("; for body\n")
			bodyLabel := strconv.Itoa(g.counter)
			g.counter++
			g.output.WriteString("; <label>:" + bodyLabel + "\n")
			_, err = g.walk(typedNode.Body)
			if err != nil {
				return "", err
			}

			//			// post
			g.output.WriteString("; for post\n")
			if typedNode.Post != nil {
				_, err = g.walk(typedNode.Post)
				if err != nil {
					return "", err
				}
			}

			// jump to start label
			g.output.WriteString("; for jump to condition\n")
			g.output.WriteString("  br label %" + conditionLabel + "\n")
			// label end
			g.output.WriteString("; for end label\n")
			endLabel := strconv.Itoa(g.counter)
			g.counter++
			g.output.WriteString("; <label>:" + endLabel + "\n")

			// write branch instruction to main buffer
			mainWriter.WriteString("; branch in for condition\n")
			mainWriter.WriteString("  br i1 " + condition + ", label %" + bodyLabel + ", label %" + endLabel + "\n")

			// flush temporary buffer to main buffer
			mainWriter.WriteString(g.output.String())

			// reset g.output to main buffer
			g.output = mainWriter

			// drop block
			g.symbolTable.PopFrame()
			return "", nil
		}
	case *ast.IncDecStmt:
		{
			// TODO
			log.Printf("case: %T\n", typedNode)
			return "", errors.New("Not Implemented")
		}
	default:
		{
			log.Printf("case: %T\n", typedNode)
			return "", errors.New("Not Implemented")
		}
	}
}
