package ast_tool

import (
	"fmt"
	"go/ast"
	"go/token"
)

type ASTStruct struct {
	Package       string // 所在包
	Name          string // 名字
	Exported      bool   // 是否向外公开
	FieldList     []*ASTField
}

// Println Debug Only
func (s *ASTStruct) Println() {
	if s == nil {
		return
	}
	fmt.Print("所在包:" + s.Package)
	fmt.Print(",是否公开:", s.Exported)
	fmt.Print(",函数名字:" + s.Name)
	fmt.Print(",属性列表:")
	fmt.Print(joinFieldName(s.FieldList))
	fmt.Println()
}

func CreateASTStructFromASTNode(node ast.Node, pkg string) *ASTStruct {
	st, ok := node.(*ast.GenDecl)
	if !ok {
		return nil
	}
	if st.Tok != token.TYPE {
		return nil
	}
	typeSpec, ok := st.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	typeNode, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	astStruct := &ASTStruct{
		Package:       pkg,
		Name:          typeSpec.Name.Name,
		Exported:      typeSpec.Name.IsExported(),
		FieldList:     parseFieldList(typeNode.Fields),
	}
	return astStruct
}
