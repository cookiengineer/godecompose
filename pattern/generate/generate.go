package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/cookiengineer/godecompose/disasm"
	"github.com/cookiengineer/godecompose/function"
	"github.com/cookiengineer/godecompose/pattern/matcher"
)

type Generator struct {
	matches      []matcher.Match
	instructions []disasm.Instruction
	functions    []*function.Function
	packages     map[string][]*function.Function
	structs      []*function.StructType
}

func New(matches []matcher.Match, instructions []disasm.Instruction) *Generator {
	return &Generator{
		matches:      matches,
		instructions: instructions,
	}
}

func NewForProject(
	matches []matcher.Match,
	instructions []disasm.Instruction,
	funcs []*function.Function,
	pkgs map[string][]*function.Function,
	structs []*function.StructType,
) *Generator {
	return &Generator{
		matches:      matches,
		instructions: instructions,
		functions:    funcs,
		packages:     pkgs,
		structs:      structs,
	}
}

func (g *Generator) Generate() string {
	sort.Slice(g.matches, func(i, j int) bool {
		return g.matches[i].StartAddr < g.matches[j].StartAddr
	})

	var buf strings.Builder
	lastAddr := uint64(0)

	if len(g.instructions) > 0 {
		lastAddr = g.instructions[0].Address
	}

	for _, match := range g.matches {
		if match.StartAddr > lastAddr {
			g.emitRawRange(&buf, lastAddr, match.StartAddr)
		}
		buf.WriteString(g.expandTemplate(match))
		buf.WriteString("\n")
		lastAddr = match.EndAddr
	}

	endAddr := lastAddr
	if len(g.instructions) > 0 {
		last := g.instructions[len(g.instructions)-1]
		endAddr = last.Address + uint64(last.Size)
	}
	if lastAddr < endAddr {
		g.emitRawRange(&buf, lastAddr, endAddr)
	}

	return buf.String()
}

func (g *Generator) WriteProject(dir string, goModule string) error {
	if goModule == "" {
		goModule = "decompiled"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	addrToIdx := make(map[uint64]int, len(g.instructions))
	for i, inst := range g.instructions {
		addrToIdx[inst.Address] = i
	}
	funcInsts := make(map[string][]disasm.Instruction)
	for _, f := range g.functions {
		for addr := f.EntryPoint; addr < f.EndAddr; {
			if idx, ok := addrToIdx[addr]; ok {
				inst := g.instructions[idx]
				funcInsts[f.Name] = append(funcInsts[f.Name], inst)
				addr += uint64(inst.Size)
			} else {
				addr++
			}
		}
	}

	pkgStructs := make(map[string][]*function.StructType)
	for _, st := range g.structs {
		pkgStructs[st.PackagePath] = append(pkgStructs[st.PackagePath], st)
	}

	sort.Slice(g.functions, func(i, j int) bool {
		return g.functions[i].EntryPoint < g.functions[j].EntryPoint
	})

	fset := token.NewFileSet()

	if mainFuncs, ok := g.packages["main"]; ok {
		if err := g.writePackageAST(fset, dir, "", mainFuncs, pkgStructs, funcInsts); err != nil {
			return err
		}
	}

	for pkgPath, funcs := range g.packages {
		if pkgPath == "main" {
			continue
		}
		pkgDir := filepath.Join(dir, pkgPath)
		if err := g.writePackageAST(fset, pkgDir, pkgPath, funcs, pkgStructs, funcInsts); err != nil {
			return err
		}
	}

	modContent := fmt.Sprintf("module %s\n\ngo 1.21\n", goModule)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644); err != nil {
		return fmt.Errorf("write go.mod: %w", err)
	}

	return nil
}

func (g *Generator) writePackageAST(fset *token.FileSet, dir string, pkgPath string, funcs []*function.Function, pkgStructs map[string][]*function.StructType, funcInsts map[string][]disasm.Instruction) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create package dir %s: %w", dir, err)
	}

	pkgName := "main"
	if pkgPath != "" {
		parts := strings.Split(pkgPath, "/")
		pkgName = parts[len(parts)-1]
	}

	fileName := pkgName + ".go"
	if pkgName == "main" {
		fileName = "main.go"
	}

	file := &ast.File{
		Name: ast.NewIdent(pkgName),
	}

	usedNames := make(map[string]bool)

	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].EntryPoint < funcs[j].EntryPoint
	})

	lookupPkg := pkgPath
	if lookupPkg == "" {
		lookupPkg = pkgName
	}
	if structs, ok := pkgStructs[lookupPkg]; ok {
		for _, st := range structs {
			if st.Name != "" {
				decl, name := g.buildStructDecl(st, usedNames)
				if decl != nil {
					file.Decls = append(file.Decls, decl)
					usedNames[name] = true
				}
			}
		}
	}

	for _, f := range funcs {
		insts := funcInsts[f.Name]
		if len(insts) == 0 {
			continue
		}

		blocks := disasm.BuildControlFlowGraph(insts, []uint64{f.EntryPoint})
		f.Blocks = blocks

		funcDecl, name := g.buildFunctionDecl(f, insts, usedNames)
		if funcDecl != nil {
			file.Decls = append(file.Decls, funcDecl)
			usedNames[name] = true
		}
	}

	f, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if err := printer.Fprint(f, fset, file); err != nil {
		return fmt.Errorf("print AST: %w", err)
	}

	return nil
}

func (g *Generator) buildStructDecl(st *function.StructType, usedNames map[string]bool) (ast.Decl, string) {
	name := deduplicateName(sanitizeFuncName(st.Name), usedNames)
	fields := function.InferStructFields(st)
	var fieldList []*ast.Field
	if len(fields) > 0 {
		for _, fld := range fields {
			ft := fld.Type
			if ft == "" {
				ft = "interface{}"
			}
			fieldList = append(fieldList, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fld.Name)},
				Type:  ast.NewIdent(ft),
				Comment: &ast.CommentGroup{List: []*ast.Comment{{
					Text: fmt.Sprintf("// offset %s (%d refs)", fld.Offset, fld.Count),
				}}},
			})
		}
	} else {
		fieldList = []*ast.Field{{
			Type: ast.NewIdent("interface{}"),
			Comment: &ast.CommentGroup{List: []*ast.Comment{{
				Text: "// fields unknown",
			}}},
		}}
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{
			Name: ast.NewIdent(name),
			Type: &ast.StructType{Fields: &ast.FieldList{List: fieldList}},
		}},
	}, name
}

func (g *Generator) buildFunctionDecl(f *function.Function, insts []disasm.Instruction, usedNames map[string]bool) (*ast.FuncDecl, string) {
	blocks := disasm.BuildControlFlowGraph(insts, []uint64{f.EntryPoint})
	if len(blocks) == 0 {
		return nil, ""
	}

	structure := disasm.StructureControlFlow(f.Name, blocks)

	var funcMatches []matcher.Match
	for _, m := range g.matches {
		if m.StartAddr >= f.EntryPoint && m.EndAddr <= f.EndAddr {
			funcMatches = append(funcMatches, m)
		}
	}
	sort.Slice(funcMatches, func(i, j int) bool {
		return funcMatches[i].StartAddr < funcMatches[j].StartAddr
	})

	sig := function.ReconstructSignature(f)
	baseName := sanitizeFuncName(sig.Name)
	name := deduplicateName(baseName, usedNames)

	body := g.buildFunctionBody(f, insts, blocks, structure, funcMatches)

	var recv *ast.FieldList
	if sig.Receiver != "" {
		recvName := strings.ToLower(sig.Receiver[:1])
		recvType := ast.NewIdent(sig.Receiver)
		recv = &ast.FieldList{List: []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent(recvName)},
			Type:  recvType,
		}}}
		if sig.IsPointer {
			recv.List[0].Type = &ast.StarExpr{X: recvType}
		}
	}

	var params []*ast.Field
	for _, a := range sig.Args {
		typ := a.Type
		if typ == "" {
			typ = "interface{}"
		}
		params = append(params, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(a.Name)},
			Type:  ast.NewIdent(typ),
		})
	}

	var results []*ast.Field
	for _, r := range sig.Returns {
		typ := r.Type
		if typ == "" {
			typ = "interface{}"
		}
		results = append(results, &ast.Field{Type: ast.NewIdent(typ)})
	}
	var resultList *ast.FieldList
	if len(results) > 0 {
		resultList = &ast.FieldList{List: results}
	}

	if name == "main" {
		params = nil
		resultList = nil
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent(name),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: resultList,
		},
		Body: &ast.BlockStmt{List: body},
		Recv: recv,
	}, name
}

func (g *Generator) buildFunctionBody(f *function.Function, insts []disasm.Instruction, blocks []*disasm.BasicBlock, structure *disasm.StructuredFunc, funcMatches []matcher.Match) []ast.Stmt {
	var body []ast.Stmt

	addrToLabel := make(map[uint64]string)
	blockByAddr := make(map[uint64]*disasm.BasicBlock)
	addrToBlockIdx := make(map[uint64]int)
	labelNum := 0
	for i, b := range blocks {
		blockByAddr[b.StartAddr] = b
		addrToBlockIdx[b.StartAddr] = i
		for _, inst := range b.Instructions {
			if inst.IsBranch && inst.BranchTarget != 0 && !inst.IsCall {
				if _, ok := addrToLabel[inst.BranchTarget]; !ok {
					addrToLabel[inst.BranchTarget] = fmt.Sprintf("L%d", labelNum)
					labelNum++
				}
			}
		}
	}

	loopHeads := make(map[int]bool)
	loopExits := make(map[int]bool)
	if structure != nil {
		for _, loop := range structure.Loops {
			loopHeads[loop[0]] = true
			loopExits[loop[1]] = true
		}
	}

	emittedBlocks := make(map[int]bool)
	usedLabels := make(map[string]bool)
	collectLabelUses(blocks, structure, &funcMatches, addrToLabel, addrToBlockIdx, usedLabels)

	for i := 0; i < len(blocks); i++ {
		if emittedBlocks[i] {
			continue
		}

		block := blocks[i]
		blockInsts := block.Instructions
		if len(blockInsts) == 0 {
			continue
		}

		if loopHeads[i] {
			stmts := g.buildLoop(blocks, i, structure, &funcMatches, addrToLabel, blockByAddr, addrToBlockIdx, emittedBlocks)
			body = append(body, stmts...)
			continue
		}

		if label, ok := addrToLabel[block.StartAddr]; ok && usedLabels[label] {
			body = append(body, &ast.LabeledStmt{
				Label: ast.NewIdent(label),
				Stmt:  &ast.EmptyStmt{},
			})
		}

		lastInst := blockInsts[len(blockInsts)-1]
		isCond := lastInst.IsConditional && len(block.Successors) > 1

		if isCond {
			stmts := g.buildCondBlock(blocks, block, i, structure, &funcMatches, addrToLabel, blockByAddr, addrToBlockIdx, emittedBlocks)
			body = append(body, stmts...)
		} else if !emittedBlocks[i] {
			stmts := g.buildPlainBlock(block, blockInsts, funcMatches)
			body = append(body, stmts...)
			emittedBlocks[i] = true
		}
	}

	return body
}

func (g *Generator) buildCondBlock(blocks []*disasm.BasicBlock, block *disasm.BasicBlock, blockIdx int, structure *disasm.StructuredFunc, funcMatches *[]matcher.Match, addrToLabel map[uint64]string, blockByAddr map[uint64]*disasm.BasicBlock, addrToBlockIdx map[uint64]int, emittedBlocks map[int]bool) []ast.Stmt {
	var stmts []ast.Stmt
	insts := block.Instructions

	var blockMatches []matcher.Match
	for _, m := range *funcMatches {
		if m.StartAddr >= block.StartAddr && m.EndAddr <= block.EndAddr {
			blockMatches = append(blockMatches, m)
		}
	}

	condExpr := ast.Expr(ast.NewIdent("true"))
	lastMatchIdx := -1
	if len(blockMatches) > 0 {
		lastIdx := len(blockMatches) - 1
		lastMatch := blockMatches[lastIdx]
		_, cond := g.parseCondFromMatch(lastMatch)
		if cond != nil {
			condExpr = cond
			lastMatchIdx = lastIdx
		}
	}

	var body []ast.Stmt

	lastAddr := block.StartAddr
	for i, match := range blockMatches {
		if i == lastMatchIdx {
			continue
		}
		if match.StartAddr > lastAddr {
			body = append(body, g.emitFunctionRangeAST(lastAddr, match.StartAddr, insts)...)
		}
		body = append(body, g.expandTemplateAST(match)...)
		lastAddr = match.EndAddr
	}
	if lastAddr < block.EndAddr {
		body = append(body, g.emitFunctionRangeAST(lastAddr, block.EndAddr, insts)...)
	}

	emittedBlocks[blockIdx] = true

	ifStmt := &ast.IfStmt{
		Cond: condExpr,
		Body: &ast.BlockStmt{List: body},
	}

	if len(block.Successors) > 1 {
		elseTarget := block.Successors[1]
		elseBlock := blockByAddr[elseTarget.StartAddr]
		if elseBlock != nil && !emittedBlocks[addrLookup(blocks, elseBlock.StartAddr)] {
			elseIdx := addrLookup(blocks, elseBlock.StartAddr)
			mergeBlock := findMergeBlock(structure, blockIdx)

			if elseIdx >= 0 && isSimpleBlockContent(elseBlock, mergeBlock) {
				elseStmts := g.buildPlainBlock(elseBlock, elseBlock.Instructions, *funcMatches)
				ifStmt.Else = &ast.BlockStmt{List: elseStmts}
				emittedBlocks[elseIdx] = true
				if mergeBlock >= 0 && mergeBlock < len(blocks) {
					emittedBlocks[mergeBlock] = false
				}
			} else if label, ok := addrToLabel[elseTarget.StartAddr]; ok {
				ifStmt.Else = &ast.BlockStmt{List: []ast.Stmt{
					&ast.BranchStmt{Tok: token.GOTO, Label: ast.NewIdent(label)},
				}}
			}
		}
	}

	stmts = append(stmts, ifStmt)
	return stmts
}

func (g *Generator) buildLoop(blocks []*disasm.BasicBlock, headIdx int, structure *disasm.StructuredFunc, funcMatches *[]matcher.Match, addrToLabel map[uint64]string, blockByAddr map[uint64]*disasm.BasicBlock, addrToBlockIdx map[uint64]int, emittedBlocks map[int]bool) []ast.Stmt {
	var stmts []ast.Stmt
	head := blocks[headIdx]
	insts := head.Instructions
	if len(insts) == 0 {
		return stmts
	}

	var blockMatches []matcher.Match
	for _, m := range *funcMatches {
		if m.StartAddr >= head.StartAddr && m.EndAddr <= head.EndAddr {
			blockMatches = append(blockMatches, m)
		}
	}

	lastInst := insts[len(insts)-1]
	isCond := lastInst.IsConditional && len(head.Successors) > 1

	var condExpr ast.Expr = nil
	if isCond {
		condExpr = ast.NewIdent("true")
		if len(blockMatches) > 0 {
			_, cond := g.parseCondFromMatch(blockMatches[len(blockMatches)-1])
			if cond != nil {
				condExpr = cond
			}
		}
	}

	var body []ast.Stmt

	lastAddr := head.StartAddr
	for _, match := range blockMatches {
		if match.StartAddr > lastAddr {
			body = append(body, g.emitFunctionRangeAST(lastAddr, match.StartAddr, insts)...)
		}
		body = append(body, g.expandTemplateAST(match)...)
		lastAddr = match.EndAddr
	}
	if lastAddr < head.EndAddr {
		body = append(body, g.emitFunctionRangeAST(lastAddr, head.EndAddr, insts)...)
	}

	emittedBlocks[headIdx] = true

	exitIdx := -1
	if isCond && len(head.Successors) > 1 {
		exitTarget := head.Successors[1]
		exitIdx = addrToBlockIdx[exitTarget.StartAddr]
	} else if !isCond && len(head.Successors) > 0 {
		nextTarget := head.Successors[0]
		exitIdx = addrToBlockIdx[nextTarget.StartAddr]
	}

	for i := headIdx + 1; i < len(blocks) && (exitIdx < 0 || i < exitIdx); i++ {
		if emittedBlocks[i] {
			continue
		}
		b := blocks[i]
		if len(b.Instructions) == 0 {
			continue
		}

		var bMatches []matcher.Match
		for _, m := range *funcMatches {
			if m.StartAddr >= b.StartAddr && m.EndAddr <= b.EndAddr {
				bMatches = append(bMatches, m)
			}
		}

		lastBAddr := b.StartAddr
		for _, match := range bMatches {
			if match.StartAddr > lastBAddr {
				body = append(body, g.emitFunctionRangeAST(lastBAddr, match.StartAddr, b.Instructions)...)
			}
			body = append(body, g.expandTemplateAST(match)...)
			lastBAddr = match.EndAddr
		}
		if lastBAddr < b.EndAddr {
			body = append(body, g.emitFunctionRangeAST(lastBAddr, b.EndAddr, b.Instructions)...)
		}
		emittedBlocks[i] = true
	}

	forStmt := &ast.ForStmt{
		Cond: condExpr,
		Body: &ast.BlockStmt{List: body},
	}
	stmts = append(stmts, forStmt)
	return stmts
}

func (g *Generator) buildPlainBlock(block *disasm.BasicBlock, insts []disasm.Instruction, funcMatches []matcher.Match) []ast.Stmt {
	var stmts []ast.Stmt

	var blockMatches []matcher.Match
	for _, m := range funcMatches {
		if m.StartAddr >= block.StartAddr && m.EndAddr <= block.EndAddr {
			blockMatches = append(blockMatches, m)
		}
	}

	lastInst := insts[len(insts)-1]
	var nonCondMatches []matcher.Match
	for _, m := range blockMatches {
		if m.StartAddr >= lastInst.Address && m.StartAddr <= lastInst.Address+uint64(lastInst.Size) {
			continue
		}
		nonCondMatches = append(nonCondMatches, m)
	}

	lastAddr := block.StartAddr
	for _, match := range nonCondMatches {
		if match.StartAddr > lastAddr {
			stmts = append(stmts, g.emitFunctionRangeAST(lastAddr, match.StartAddr, insts)...)
		}
		stmts = append(stmts, g.expandTemplateAST(match)...)
		lastAddr = match.EndAddr
	}
	if lastAddr < block.EndAddr {
		stmts = append(stmts, g.emitFunctionRangeAST(lastAddr, block.EndAddr, insts)...)
	}

	return stmts
}

func (g *Generator) parseCondFromMatch(match matcher.Match) ([]ast.Stmt, ast.Expr) {
	template := g.expandTemplate(match)
	template = strings.TrimSpace(template)

	cond := g.parseCond(template)

	if strings.HasPrefix(template, "if ") {
		if idx := strings.Index(template, "{"); idx >= 0 {
			bodyStr := template[idx:]
			stmts := g.parseTemplateStmts(bodyStr)
			return stmts, cond
		}
	}

	return nil, cond
}

func (g *Generator) parseCond(template string) ast.Expr {
	if strings.HasPrefix(template, "if ") {
		condStr := template[3:]
		if idx := strings.Index(condStr, "{"); idx >= 0 {
			condStr = strings.TrimSpace(condStr[:idx])
		}
		condStr = strings.TrimSpace(condStr)
		if expr, err := parser.ParseExpr(condStr); err == nil {
			return expr
		}
		return parseFallbackCond(condStr)
	}
	return ast.NewIdent("true")
}

func parseFallbackCond(condStr string) ast.Expr {
	condStr = strings.TrimSpace(condStr)
	if condStr == "" || condStr == "_" || condStr == "true" {
		return ast.NewIdent("true")
	}

	for _, op := range []string{"==", "!=", "<=", ">=", "<", ">"} {
		if idx := strings.Index(condStr, op); idx >= 0 {
			left := strings.TrimSpace(condStr[:idx])
			right := strings.TrimSpace(condStr[idx+len(op):])
			var tok token.Token
			switch op {
			case "==":
				tok = token.EQL
			case "!=":
				tok = token.NEQ
			case "<":
				tok = token.LSS
			case ">":
				tok = token.GTR
			case "<=":
				tok = token.LEQ
			case ">=":
				tok = token.GEQ
			}
			return &ast.BinaryExpr{
				X:  parseIdentOrLiteral(left),
				Op: tok,
				Y:  parseIdentOrLiteral(right),
			}
		}
	}

	return parseIdentOrLiteral(condStr)
}

func parseIdentOrLiteral(s string) ast.Expr {
	s = strings.TrimSpace(s)
	if s == "nil" {
		return ast.NewIdent("nil")
	}
	if s == "true" || s == "false" {
		return ast.NewIdent(s)
	}
	if n, err := strconv.ParseInt(s, 0, 64); err == nil {
		return &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", n)}
	}
	return ast.NewIdent(s)
}

func (g *Generator) parseTemplateStmts(bodyStr string) []ast.Stmt {
	bodyStr = strings.TrimSpace(bodyStr)
	if strings.HasPrefix(bodyStr, "{") {
		bodyStr = bodyStr[1:]
	}
	if strings.HasSuffix(bodyStr, "}") {
		bodyStr = bodyStr[:len(bodyStr)-1]
	}
	bodyStr = strings.TrimSpace(bodyStr)

	if bodyStr == "" {
		return nil
	}

	if stmts, err := parseGoStmts(bodyStr); err == nil {
		return stmts
	}

	if strings.Contains(bodyStr, "return ") || strings.Contains(bodyStr, "panic(") || strings.Contains(bodyStr, "goto ") {
		stmts, _ := parseGoStmts(bodyStr + "\n")
		return stmts
	}

	return []ast.Stmt{
		&ast.EmptyStmt{},
	}
}

func parseGoStmts(src string) ([]ast.Stmt, error) {
	src = strings.TrimSpace(src)
	if src == "" {
		return nil, nil
	}
	file, err := parser.ParseFile(token.NewFileSet(), "", "package _; func _() { "+src+" }", parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(file.Decls) > 0 {
		if fd, ok := file.Decls[0].(*ast.FuncDecl); ok && fd.Body != nil {
			return fd.Body.List, nil
		}
	}
	return nil, fmt.Errorf("parse failed")
}

func (g *Generator) expandTemplateAST(match matcher.Match) []ast.Stmt {
	template := g.expandTemplate(match)
	if template == "" {
		return nil
	}

	template = strings.TrimSpace(template)

	if stmts, err := parseGoStmts(template); err == nil {
		return filterUndefinedGotos(stmts)
	}

	return nil
}

func filterUndefinedGotos(stmts []ast.Stmt) []ast.Stmt {
	var result []ast.Stmt
	for _, s := range stmts {
		switch st := s.(type) {
		case *ast.BranchStmt:
			if st.Tok == token.GOTO {
				continue
			}
		case *ast.IfStmt:
			if st.Body != nil {
				st.Body.List = filterUndefinedGotos(st.Body.List)
			}
			if st.Else != nil {
				if elseBlock, ok := st.Else.(*ast.BlockStmt); ok {
					elseBlock.List = filterUndefinedGotos(elseBlock.List)
				}
			}
		case *ast.BlockStmt:
			st.List = filterUndefinedGotos(st.List)
		case *ast.ForStmt:
			if st.Body != nil {
				st.Body.List = filterUndefinedGotos(st.Body.List)
			}
		}
		result = append(result, s)
	}
	return result
}

func (g *Generator) emitFunctionRangeAST(start, end uint64, insts []disasm.Instruction) []ast.Stmt {
	return nil
}

func (g *Generator) emitRawRange(buf *strings.Builder, start, end uint64) {
	var rangeInsts []disasm.Instruction
	noiseOps := map[string]bool{"NOP": true, "NOPL": true, "NOPW": true, "INT": true, "INT3": true, "DATA16": true}
	for _, inst := range g.instructions {
		if inst.Address >= start && inst.Address < end {
			if noiseOps[inst.Opcode] {
				continue
			}
			rangeInsts = append(rangeInsts, inst)
		}
	}
	if len(rangeInsts) == 0 {
		return
	}

	buf.WriteString("\t// unresolved asm\n")
	for _, inst := range rangeInsts {
		buf.WriteString(fmt.Sprintf("\t// %016x: %s\n", inst.Address, inst.IntelSyntax))
	}
}

func (g *Generator) expandTemplate(match matcher.Match) string {
	template := match.Pattern.GenTemplate
	if template == "" {
		return fmt.Sprintf("// matched %s", match.Pattern.Name)
	}

	result := template
	for _, b := range match.Bindings {
		placeholder := "$" + b.CaptureVar
		value := b.Value
		if b.Alias != "" {
			value = b.Alias
		}
		if strings.HasPrefix(value, "0x") {
			value = "L_" + strings.TrimPrefix(value, "0x")
		}
		result = strings.ReplaceAll(result, placeholder, value)
	}

	for _, b := range match.Pattern.Bindings {
		placeholder := "$" + b.CaptureVar
		if _, ok := match.Bindings[b.CaptureVar]; !ok {
			alias := b.Alias
			if strings.HasPrefix(alias, "0x") {
				alias = "L_" + strings.TrimPrefix(alias, "0x")
			}
			result = strings.ReplaceAll(result, placeholder, alias)
		}
	}

	return result
}

func addrLookup(blocks []*disasm.BasicBlock, addr uint64) int {
	for i, b := range blocks {
		if b.StartAddr == addr {
			return i
		}
	}
	return -1
}

func isSimpleBlockContent(block *disasm.BasicBlock, mergeIdx int) bool {
	if block == nil || len(block.Instructions) == 0 {
		return false
	}
	last := block.Instructions[len(block.Instructions)-1]
	if last.IsBranch && !last.IsConditional {
		return true
	}
	if last.IsConditional {
		return false
	}
	return len(block.Instructions) <= 5
}

func findMergeBlock(structure *disasm.StructuredFunc, condIdx int) int {
	if structure == nil || structure.PdomTree == nil {
		return -1
	}
	if condIdx < 0 || condIdx >= len(structure.PdomTree) {
		return -1
	}
	return structure.PdomTree[condIdx]
}

func collectLabelUses(blocks []*disasm.BasicBlock, structure *disasm.StructuredFunc, funcMatches *[]matcher.Match, addrToLabel map[uint64]string, addrToBlockIdx map[uint64]int, usedLabels map[string]bool) {
	for i, block := range blocks {
		insts := block.Instructions
		if len(insts) == 0 {
			continue
		}
		lastInst := insts[len(insts)-1]
		isCond := lastInst.IsConditional && len(block.Successors) > 1
		if !isCond {
			continue
		}
		if len(block.Successors) > 1 {
			elseTarget := block.Successors[1]
			if label, ok := addrToLabel[elseTarget.StartAddr]; ok {
				elseIdx, _ := addrToBlockIdx[elseTarget.StartAddr]
				if elseIdx >= 0 && elseIdx < len(blocks) {
					eBlock := blocks[elseIdx]
					mergeBlock := findMergeBlock(structure, i)
					if len(eBlock.Instructions) > 0 {
						eLast := eBlock.Instructions[len(eBlock.Instructions)-1]
						nextIsCond := eLast.IsConditional && len(eBlock.Successors) > 1
						if nextIsCond || !isSimpleBlockContent(eBlock, mergeBlock) {
							usedLabels[label] = true
						}
					} else if !isSimpleBlockContent(eBlock, mergeBlock) {
						usedLabels[label] = true
					}
				} else {
					usedLabels[label] = true
				}
			}
		}
	}
}

func deduplicateName(base string, usedNames map[string]bool) string {
	name := base
	for i := 1; usedNames[name]; i++ {
		name = fmt.Sprintf("%s_%d", base, i)
	}
	return name
}

func sanitizeFuncName(name string) string {
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "/", "_")

	if name == "" || (len(name) > 0 && name[0] >= '0' && name[0] <= '9') {
		name = "_" + name
	}

	if name == "init" || name == "main_init" {
		return "_init"
	}
	goKeywords := map[string]bool{
		"break": true, "case": true, "chan": true, "const": true, "continue": true,
		"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
		"func": true, "go": true, "goto": true, "if": true, "import": true,
		"interface": true, "map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true, "var": true,
	}
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		if goKeywords[p] {
			parts[i] = "_" + p
		}
	}
	return strings.Join(parts, "_")
}
