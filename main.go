package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/fkautz/gomez/libgomez"
	"github.com/spf13/cobra"
)

var cfgFileToCompile string
var cfgFileToOutput string

func main() {
	cfgFileToCompile = "simple.go"
	cfgFileToOutput = "out.ll"
	cmdCompiler := &cobra.Command{
		Use:   "compile",
		Short: "compile a gomez app",
		Long:  "Gomez is a compiler",
		Run:   runCompiler,
	}
	cmdCompiler.Flags().StringVarP(&cfgFileToCompile, "input", "i", "simple.go", "file to compile")
	cmdCompiler.Flags().StringVarP(&cfgFileToOutput, "output", "o", "simple.ll", "output file")

	rootCmd := &cobra.Command{Use: "App"}
	rootCmd.AddCommand(cmdCompiler)
	rootCmd.Execute()
}

func runCompiler(cmd *cobra.Command, args []string) {
	file, err := os.OpenFile(cfgFileToOutput, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	llvmIR, err := compileGomezToLLVM(cfgFileToCompile)
	log.Println(llvmIR)
	file.WriteString(llvmIR)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func compileGomezToLLVM(filename string) (string, error) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return "", err
	}

	ast.Print(fset, tree)
	return libgomez.GenerateLLVM(fset, tree)
}
