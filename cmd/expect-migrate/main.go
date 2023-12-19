// Package main contains a cli application that rewrites go source code containing test code using
// github.com/halimath/expect-go with a dot import and matchers from v0.1.0 to the same test code using
// github.com/halimath/expect v0.3.0 with regular imports. It handles most of the matchers (for some a
// manual work is needed), handels ExpectThat and EnsureThat but does not wrap multiple expectations into
// a single call to expect.That. Thus, this tool can reduce the amount of manual work needed when upgrading
// expect to v0.3.0 but some remainings are still needed in order to have readable, maintainable test code.
//
// Usage:
//   - either: cat source_test.go | expect1to3 > source_test.go
//   - or: expect1to3 source_test.go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strconv"
)

func main() {
	var err error

	if len(os.Args) == 1 {
		err = migrate("stdin.go", os.Stdin, os.Stdout)
	} else {
		err = migrateFile(os.Args[1])
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func migrateFile(filename string) error {
	renamedFilename := filename + "~"
	if err := os.Rename(filename, renamedFilename); err != nil {
		return fmt.Errorf("failed to rename input file: %v", err)
	}

	in, err := os.Open(renamedFilename)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer in.Close()

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open output file: %v", err)
	}
	defer out.Close()

	if err := migrate(filename, in, out); err != nil {
		return err
	}

	return nil
}

func migrate(filename string, src io.Reader, out io.Writer) error {
	// Parse the source file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return fmt.Errorf("failed to parse source file: %v", err)
	}

	// Rewrite the AST to update existing calls to ExpectThat as well as
	// change the existing import.
	ast.Walk(&expectationRewriteVisitor{}, file)

	// Rewrite existing import statement
	ast.Inspect(file, rewriteImports)

	// Insert another import for expect/is
	addImportToIsPackage(file)

	// Render the updated AST back to go source code.
	if err := printer.Fprint(out, fset, file); err != nil {
		return fmt.Errorf("failed to write source code: %v", err)
	}

	return nil
}

var matcherNameTranslationTable = map[string]string{
	"Equal":     "EqualTo",
	"DeepEqual": "DeepEqualTo",
}

func translateMatcherFuncName(n string) string {
	// TODO: What about the matcher Len, Nil, NotNil?

	if translated, ok := matcherNameTranslationTable[n]; ok {
		return translated
	}

	return n
}

var chainingMethodNames = []string{
	"Is",
	"Has",
	"And",
	"Matches",
}

func addImportToIsPackage(file *ast.File) {
	var isImportFound bool
	for _, imp := range file.Imports {
		if imp.Path.Value == strconv.Quote("github.com/halimath/expect/is") {
			isImportFound = true
			break
		}
	}

	if isImportFound {
		return
	}

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			importSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote("github.com/halimath/expect/is")}}
			genDecl.Specs = append(genDecl.Specs, importSpec)
			break
		}
	}

}

func rewriteImports(node ast.Node) bool {
	if imp, ok := node.(*ast.ImportSpec); ok {
		if imp.Path.Value == "\"github.com/halimath/expect-go\"" {
			imp.Path.Value = "\"github.com/halimath/expect\""
			imp.Name = nil
		}

		return false
	}

	return true
}

type blockFrame struct {
	block            *ast.BlockStmt
	lastExpectation  *ast.CallExpr
	currentStmtIndex int
	stmtsToRemove    []int
}

type expectationRewriteVisitor struct {
	blocks Stack[*blockFrame]
}

func (v *expectationRewriteVisitor) removeMarkedStatments(currentBlock *blockFrame) {
	stmts := make([]ast.Stmt, 0, len(currentBlock.block.List))

	for i := range currentBlock.block.List {
		if !contains(currentBlock.stmtsToRemove, i) {
			stmts = append(stmts, currentBlock.block.List[i])
		}
	}

	currentBlock.block.List = stmts
}

func (v *expectationRewriteVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	currentBlock, _ := v.blocks.Peek()

	if b, ok := node.(*ast.BlockStmt); ok {
		v.blocks.Push(&blockFrame{
			block:            b,
			currentStmtIndex: -1,
		})

		for _, n := range b.List {
			ast.Walk(v, n)
		}

		blockStatement, _ := v.blocks.Pop()
		v.removeMarkedStatments(blockStatement)

		return nil
	}

	if _, ok := node.(ast.Stmt); ok {
		if currentBlock != nil {
			currentBlock.currentStmtIndex++
		}
	}

	exprStmt, ok := node.(*ast.ExprStmt)
	if !ok {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	call, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}
		return v
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	if !contains(chainingMethodNames, sel.Sel.Name) {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	if len(call.Args) != 1 {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	matcher, ok := call.Args[0].(*ast.CallExpr)
	if !ok {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	expectThatCall, ok := sel.X.(*ast.CallExpr)
	if !ok {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	expectThat, ok := expectThatCall.Fun.(*ast.Ident)
	if !ok {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	if expectThat.Name != "ExpectThat" && expectThat.Name != "EnsureThat" {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return v
	}

	failNow := expectThat.Name == "EnsureThat"

	if len(expectThatCall.Args) != 2 {
		if currentBlock != nil {
			currentBlock.lastExpectation = nil
		}

		return nil
	}

	got := expectThatCall.Args[1]

	newMatcherArgs := []ast.Expr{got}
	newMatcherArgs = append(newMatcherArgs, matcher.Args...)
	matcher.Args = newMatcherArgs

	var matcherFun ast.Expr = &ast.SelectorExpr{
		X: &ast.Ident{
			Name: "is",
		},
		Sel: &ast.Ident{
			Name: translateMatcherFuncName(matcher.Fun.(*ast.Ident).Name),
		},
	}

	matcher.Fun = matcherFun

	if failNow {
		failNowCall := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "expect"},
				Sel: &ast.Ident{Name: "FailNow"},
			},
			Args: []ast.Expr{matcher},
		}
		matcher = failNowCall
	}

	if currentBlock.lastExpectation != nil {
		// The previous expression was a call to expect.That
		// Append the expectation that call
		currentBlock.lastExpectation.Args = append(currentBlock.lastExpectation.Args, matcher)

		// Mark the statement for removal
		currentBlock.stmtsToRemove = append(currentBlock.stmtsToRemove, currentBlock.currentStmtIndex)
	} else {
		call.Fun = &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "expect",
			},
			Sel: &ast.Ident{
				Name: "That",
			},
		}

		call.Args = expectThatCall.Args[:1]

		call.Args = append(call.Args, matcher)
		currentBlock.lastExpectation = call
	}

	return nil
}

type Stack[T any] []T

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Stack[T]) Pop() (v T, ok bool) {
	if len(*s) == 0 {
		return
	}

	v = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	ok = true
	return
}

func (s *Stack[T]) Peek() (v T, ok bool) {
	if len(*s) == 0 {
		return
	}

	v = (*s)[len(*s)-1]
	ok = true
	return
}

func contains[S ~[]E, E comparable](s S, v E) bool {
	for i := range s {
		if s[i] == v {
			return true
		}
	}

	return false
}
