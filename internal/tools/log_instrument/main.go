package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type funcContext struct {
	file        *ast.File
	fset        *token.FileSet
	funcName    string
	resultNames []string
}

func main() {
	root := flag.String("root", ".", "root directory to process")
	flag.Parse()

	var goFiles []string
	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() != "." && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if info.Name() == "vendor" || info.Name() == "tmp" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".go" {
			if strings.Contains(path, filepath.Join("internal", "tools", "log_instrument")) {
				return nil
			}
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "walk error: %v\n", err)
		os.Exit(1)
	}

	sort.Strings(goFiles)

	for _, file := range goFiles {
		if err := instrumentFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "failed to instrument %s: %v\n", file, err)
			os.Exit(1)
		}
	}
}

func instrumentFile(path string) error {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	modified := false

	for _, decl := range fileAst.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		if alreadyInstrumented(fn) {
			continue
		}

		ctx := funcContext{
			file: fileAst,
			fset: fset,
		}
		ctx.funcName = buildFunctionName(fn)

		resultNames := ensureNamedResults(fn)
		ctx.resultNames = resultNames

		prependStatements(fn, buildInstrumentationStmts(&ctx, fn))

		processBlock(&ctx, fn.Body)

		modified = true
	}

	if !modified {
		return nil
	}

	astutil.AddImport(fset, fileAst, "go.uber.org/zap")
	astutil.AddImport(fset, fileAst, "time")

	fmt.Println("instrumented", path)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, fileAst); err != nil {
		return fmt.Errorf("format file: %w", err)
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func alreadyInstrumented(fn *ast.FuncDecl) bool {
	if len(fn.Body.List) == 0 {
		return false
	}
	if assign, ok := fn.Body.List[0].(*ast.AssignStmt); ok {
		if len(assign.Lhs) == 1 {
			if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == "__logParams" {
				return true
			}
		}
	}
	if stmt, ok := fn.Body.List[0].(*ast.ExprStmt); ok {
		if call, ok := stmt.X.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if innerCall, ok := sel.X.(*ast.CallExpr); ok {
					if sel2, ok := innerCall.Fun.(*ast.SelectorExpr); ok {
						if id, ok := sel2.X.(*ast.Ident); ok && id.Name == "zap" && sel2.Sel.Name == "L" {
							if sel.Sel.Name == "Info" && len(call.Args) > 0 {
								if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING && lit.Value == "\"function.entry\"" {
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func buildFunctionName(fn *ast.FuncDecl) string {
	if fn.Recv == nil {
		return fn.Name.Name
	}
	if len(fn.Recv.List) == 0 {
		return fn.Name.Name
	}
	recv := fn.Recv.List[0]
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), recv.Type)
	return fmt.Sprintf("%s.%s", buf.String(), fn.Name.Name)
}

func ensureNamedResults(fn *ast.FuncDecl) []string {
	results := fn.Type.Results
	if results == nil || len(results.List) == 0 {
		return nil
	}

	var names []string
	counter := 0
	for _, field := range results.List {
		if len(field.Names) == 0 {
			name := fmt.Sprintf("result%d", counter)
			field.Names = []*ast.Ident{ast.NewIdent(name)}
			counter++
			names = append(names, name)
		} else {
			for _, name := range field.Names {
				names = append(names, name.Name)
			}
		}
	}
	return names
}

func prependStatements(fn *ast.FuncDecl, stmts []ast.Stmt) {
	newList := make([]ast.Stmt, 0, len(stmts)+len(fn.Body.List))
	newList = append(newList, stmts...)
	newList = append(newList, fn.Body.List...)
	fn.Body.List = newList
}

func buildInstrumentationStmts(ctx *funcContext, fn *ast.FuncDecl) []ast.Stmt {
	paramsLiteral := &ast.CompositeLit{
		Type: &ast.MapType{Key: ast.NewIdent("string"), Value: ast.NewIdent("any")},
	}

	addParam := func(name string) {
		if name == "" || name == "_" {
			return
		}
		paramsLiteral.Elts = append(paramsLiteral.Elts, &ast.KeyValueExpr{
			Key:   &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(name)},
			Value: ast.NewIdent(name),
		})
	}

	if fn.Recv != nil {
		for _, field := range fn.Recv.List {
			for _, name := range field.Names {
				addParam(name.Name)
			}
		}
	}

	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			if len(field.Names) == 0 {
				continue
			}
			for _, name := range field.Names {
				addParam(name.Name)
			}
		}
	}

	stmts := []ast.Stmt{}

	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("__logParams")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{paramsLiteral},
	})

	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("__logStart")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.SelectorExpr{X: ast.NewIdent("time"), Sel: ast.NewIdent("Now")}}},
	})

	resultExpr := buildResultExpr(ctx.resultNames)

	deferArgs := []ast.Expr{
		&ast.BasicLit{Kind: token.STRING, Value: "\"function.exit\""},
		buildZapString("func", ctx.funcName),
		buildZapAny("result", resultExpr),
		buildZapDuration("duration", &ast.CallExpr{
			Fun:  &ast.SelectorExpr{X: ast.NewIdent("time"), Sel: ast.NewIdent("Since")},
			Args: []ast.Expr{ast.NewIdent("__logStart")},
		}),
	}

	stmts = append(stmts, &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.FuncLit{
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{X: buildZapCall("Info", deferArgs)},
					},
				},
			},
		},
	})

	entryArgs := []ast.Expr{
		&ast.BasicLit{Kind: token.STRING, Value: "\"function.entry\""},
		buildZapString("func", ctx.funcName),
		buildZapAny("params", ast.NewIdent("__logParams")),
	}

	stmts = append(stmts, &ast.ExprStmt{X: buildZapCall("Info", entryArgs)})

	return stmts
}

func buildResultExpr(resultNames []string) ast.Expr {
	if len(resultNames) == 0 {
		return ast.NewIdent("nil")
	}
	if len(resultNames) == 1 {
		return ast.NewIdent(resultNames[0])
	}
	lit := &ast.CompositeLit{
		Type: &ast.MapType{Key: ast.NewIdent("string"), Value: ast.NewIdent("any")},
	}
	for _, name := range resultNames {
		lit.Elts = append(lit.Elts, &ast.KeyValueExpr{
			Key:   &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(name)},
			Value: ast.NewIdent(name),
		})
	}
	return lit
}

func buildZapCall(method string, args []ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.CallExpr{Fun: &ast.SelectorExpr{X: ast.NewIdent("zap"), Sel: ast.NewIdent("L")}},
			Sel: ast.NewIdent(method),
		},
		Args: args,
	}
}

func buildZapString(key, value string) ast.Expr {
	return &ast.CallExpr{
		Fun:  &ast.SelectorExpr{X: ast.NewIdent("zap"), Sel: ast.NewIdent("String")},
		Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(key)}, &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(value)}},
	}
}

func buildZapAny(key string, value ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun:  &ast.SelectorExpr{X: ast.NewIdent("zap"), Sel: ast.NewIdent("Any")},
		Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(key)}, value},
	}
}

func buildZapDuration(key string, value ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun:  &ast.SelectorExpr{X: ast.NewIdent("zap"), Sel: ast.NewIdent("Duration")},
		Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(key)}, value},
	}
}

func processBlock(ctx *funcContext, block *ast.BlockStmt) {
	if block == nil {
		return
	}

	newList := make([]ast.Stmt, 0, len(block.List))
	for _, stmt := range block.List {
		switch s := stmt.(type) {
		case *ast.ReturnStmt:
			newList = append(newList, rewriteReturn(ctx, s)...)
			continue
		case *ast.IfStmt:
			processBlock(ctx, s.Body)
			insertErrorLog(ctx, s)
			processElse(ctx, s.Else)
		case *ast.ForStmt:
			processBlock(ctx, s.Body)
		case *ast.RangeStmt:
			processBlock(ctx, s.Body)
		case *ast.SwitchStmt:
			processSwitch(ctx, s)
		case *ast.TypeSwitchStmt:
			processTypeSwitch(ctx, s)
		case *ast.SelectStmt:
			processSelect(ctx, s)
		}
		newList = append(newList, stmt)
	}
	block.List = newList
}

func processElse(ctx *funcContext, stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		processBlock(ctx, s)
	case *ast.IfStmt:
		processBlock(ctx, s.Body)
		insertErrorLog(ctx, s)
		processElse(ctx, s.Else)
	}
}

func rewriteReturn(ctx *funcContext, ret *ast.ReturnStmt) []ast.Stmt {
	if len(ret.Results) == 0 {
		return []ast.Stmt{ret}
	}

	// Handle tuple-returning expressions (e.g., return foo()) when multiple results are expected.
	if len(ret.Results) == 1 && len(ctx.resultNames) > 1 {
		assign := &ast.AssignStmt{
			Lhs: make([]ast.Expr, len(ctx.resultNames)),
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ret.Results[0]},
		}
		for i, name := range ctx.resultNames {
			assign.Lhs[i] = ast.NewIdent(name)
		}
		return []ast.Stmt{assign, &ast.ReturnStmt{}}
	}

	stmts := make([]ast.Stmt, 0, len(ret.Results)+1)
	for i, expr := range ret.Results {
		if i >= len(ctx.resultNames) {
			continue
		}
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(ctx.resultNames[i])},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{expr},
		})
	}
	stmts = append(stmts, &ast.ReturnStmt{})
	return stmts
}

func processSwitch(ctx *funcContext, s *ast.SwitchStmt) {
	if s.Body == nil {
		return
	}
	for _, stmt := range s.Body.List {
		if clause, ok := stmt.(*ast.CaseClause); ok {
			block := &ast.BlockStmt{List: clause.Body}
			processBlock(ctx, block)
			clause.Body = block.List
		}
	}
}

func processTypeSwitch(ctx *funcContext, s *ast.TypeSwitchStmt) {
	if s.Body == nil {
		return
	}
	for _, stmt := range s.Body.List {
		if clause, ok := stmt.(*ast.CaseClause); ok {
			block := &ast.BlockStmt{List: clause.Body}
			processBlock(ctx, block)
			clause.Body = block.List
		}
	}
}

func processSelect(ctx *funcContext, s *ast.SelectStmt) {
	if s.Body == nil {
		return
	}
	for _, stmt := range s.Body.List {
		if comm, ok := stmt.(*ast.CommClause); ok {
			block := &ast.BlockStmt{List: comm.Body}
			processBlock(ctx, block)
			comm.Body = block.List
		}
	}
}

func insertErrorLog(ctx *funcContext, ifStmt *ast.IfStmt) {
	errIdent := detectErrorIdent(ifStmt.Cond)
	if errIdent == nil {
		return
	}

	if len(ifStmt.Body.List) > 0 {
		if existing, ok := ifStmt.Body.List[0].(*ast.ExprStmt); ok {
			if call, ok := existing.X.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "Error" {
						if callLit, ok := call.Args[0].(*ast.BasicLit); ok && callLit.Kind == token.STRING && callLit.Value == "\"function.error\"" {
							return
						}
					}
				}
			}
		}
	}

	errArgs := []ast.Expr{
		&ast.BasicLit{Kind: token.STRING, Value: "\"function.error\""},
		buildZapString("func", ctx.funcName),
		buildZapError(errIdent.Name),
		buildZapAny("params", ast.NewIdent("__logParams")),
	}

	logStmt := &ast.ExprStmt{X: buildZapCall("Error", errArgs)}
	ifStmt.Body.List = append([]ast.Stmt{logStmt}, ifStmt.Body.List...)
}

func detectErrorIdent(expr ast.Expr) *ast.Ident {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op != token.NEQ && e.Op != token.EQL {
			return nil
		}
		if id, ok := e.X.(*ast.Ident); ok {
			if isNilIdent(e.Y) && e.Op == token.NEQ && isErrorName(id.Name) {
				return id
			}
		}
		if id, ok := e.Y.(*ast.Ident); ok {
			if isNilIdent(e.X) && e.Op == token.NEQ && isErrorName(id.Name) {
				return id
			}
		}
	}
	return nil
}

func isNilIdent(expr ast.Expr) bool {
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name == "nil"
	}
	return false
}

func buildZapError(name string) ast.Expr {
	return &ast.CallExpr{
		Fun:  &ast.SelectorExpr{X: ast.NewIdent("zap"), Sel: ast.NewIdent("Error")},
		Args: []ast.Expr{ast.NewIdent(name)},
	}
}

func isErrorName(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "err")
}
