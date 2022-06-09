package ast_tool_test_test

import (
	"github.com/obgnail/golangInPython/ast_tool"
	"testing"
)

func assertEqual(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("Not Equal. %d %d", a, b)
	}
}

func TestCreateASTFromFile(t *testing.T) {
	functions, structs, err := ast_tool.CreateASTFromFile("./source.go")
	assertEqual(t, err, nil)
	assertEqual(t, len(functions), 1)
	assertEqual(t, len(structs), 1)

	st := structs[0]
	assertEqual(t, st.Name, "Dic")
	assertEqual(t, st.Package, "ast_tool_test")
	assertEqual(t, st.FieldList[0].Type, "string")
	assertEqual(t, st.FieldList[1].Name, "Cookie")

	fn := functions[0]
	assertEqual(t, fn.Name, "GetName")
	assertEqual(t, fn.Package, "ast_tool_test")
	assertEqual(t, fn.CGOExportName, "GetName")
	assertEqual(t, fn.Params[0].Name, "pyStr")
	assertEqual(t, fn.Params[0].Type, "*C.char")
	assertEqual(t, len(fn.Results), 0)

	for _, f := range functions {
		f.Println()
	}
	for _, f := range structs {
		f.Println()
	}
}

func TestCreateASTFunction(t *testing.T) {
	functions, err := ast_tool.CreateASTFunction("/Users/heyingliang/go/src/github.com/obgnail/golangInPython")
	assertEqual(t, err, nil)

	for _, f := range functions {
		f.Println()
	}
}
