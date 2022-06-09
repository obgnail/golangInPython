package ast_tool

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/obgnail/golangInPython/utils"
	log "github.com/sirupsen/logrus"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	ExportFlag = "//export"
)

var (
	ExportRegex *regexp.Regexp = regexp.MustCompile(`//export\s*?(\w+)`)
)

type ASTField struct {
	Name string
	Type string
	Tag  string
}

func NewASTField(name, Typ, tag string) *ASTField {
	return &ASTField{Name: name, Type: Typ, Tag: tag}
}

func joinFieldName(fields []*ASTField) string {
	res := make([]string, 0)
	for _, field := range fields {
		res = append(res, field.Name+" "+field.Type)
	}
	return strings.Join(res, ",")
}

// ASTFunction AST函数对象
type ASTFunction struct {
	Package       string      // 函数所在包
	Name          string      // 函数名字
	Exported      bool        // 是否向外公开
	CGOExportName string      // cgo导出的名字
	Recv          []*ASTField // 函数接收者
	Params        []*ASTField // 函数参数
	Results       []*ASTField // 函数返回值
}

// Println Debug Only
func (f *ASTFunction) Println() {
	if f == nil {
		return
	}
	fmt.Print("所在包:" + f.Package)
	fmt.Print(",是否公开:", f.Exported)
	fmt.Print(",函数名字:" + f.Name)
	fmt.Print(",CGO函数名称:", f.CGOExportName)
	fmt.Print(",函数接收者:")
	fmt.Print(joinFieldName(f.Recv))
	fmt.Print(",函数参数:")
	fmt.Print(joinFieldName(f.Params))
	fmt.Print(",函数返回值:")
	fmt.Print(joinFieldName(f.Results))
	fmt.Println()
}

// 1.数组或者切片类型 √
// 2.用户定义的类型或基本数据类型 √
// 3.选择表达式  √
// 4.指针表达式 √
// 5.映射类型 √
// 6.函数类型 √
// 7.管道类型 √
// 8.匿名结构体 ×
func exprToTypeStringRecursively(expr ast.Expr) string {
	if arr, ok := expr.(*ast.ArrayType); ok {
		if arr.Len == nil {
			return "[]" + exprToTypeStringRecursively(arr.Elt)
		} else if lit, ok := arr.Len.(*ast.BasicLit); ok {
			return fmt.Sprintf("[%s]%s", lit.Value, exprToTypeStringRecursively(arr.Elt))
		} else {
			// TODO 完备性检查
			log.Fatalf("no such expr: %v", expr)
		}
	}
	if _, ok := expr.(*ast.InterfaceType); ok {
		return "interface{}"
	}
	if indent, ok := expr.(*ast.Ident); ok {
		return indent.Name
	} else if selExpr, ok := expr.(*ast.SelectorExpr); ok {
		return exprToTypeStringRecursively(selExpr.X) + "." + exprToTypeStringRecursively(selExpr.Sel)
	} else if star, ok := expr.(*ast.StarExpr); ok {
		return "*" + exprToTypeStringRecursively(star.X)
	} else if mapType, ok := expr.(*ast.MapType); ok {
		return fmt.Sprintf("map[%s]%s", exprToTypeStringRecursively(mapType.Key), exprToTypeStringRecursively(mapType.Value))
	} else if funcType, ok := expr.(*ast.FuncType); ok {
		params := parseFieldList(funcType.Params)
		results := parseFieldList(funcType.Results)
		return fmt.Sprintf("func(%s)(%s)", joinFieldName(params), joinFieldName(results))
	} else if chanType, ok := expr.(*ast.ChanType); ok {
		if chanType.Dir == ast.SEND {
			return "chan <- " + exprToTypeStringRecursively(chanType.Value)
		} else if chanType.Dir == ast.RECV {
			return "<- chan " + exprToTypeStringRecursively(chanType.Value)
		} else {
			return "chan " + exprToTypeStringRecursively(chanType.Value)
		}
	} else if chanType, ok := expr.(*ast.StructType); ok {
		// 不考虑匿名结构体类型
		log.Warnf("do not support StructType: %v", chanType)
	}
	// TODO 完备性检查
	log.Fatal("err expr: %v", expr)
	return ""
}

func parseFieldList(fList *ast.FieldList) []*ASTField {
	if fList == nil {
		return nil
	}

	dst := make([]*ASTField, 0)
	list := fList.List
	for i := 0; i < len(list); i++ {
		names := list[i].Names
		typeStr := exprToTypeStringRecursively(list[i].Type)
		fieldTag := ""
		if list[i].Tag != nil {
			fieldTag, _ = strconv.Unquote(list[i].Tag.Value)
		}
		for j := 0; j < len(names); j++ {
			dst = append(dst, NewASTField(names[j].Name, typeStr, fieldTag))
		}
		if len(names) == 0 {
			dst = append(dst, NewASTField("", typeStr, fieldTag))
		}
	}
	return dst
}

func parseCGOExportName(fDoc *ast.CommentGroup) string {
	if fDoc == nil {
		return ""
	}
	for _, doc := range fDoc.List {
		if strings.HasPrefix(doc.Text, ExportFlag) {
			exportName := ExportRegex.FindStringSubmatch(doc.Text)
			if len(exportName) == 2 {
				return exportName[1]
			}
		}
	}
	return ""
}

func CreateASTFunctionFromASTNode(node ast.Node, pkg string) *ASTFunction {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}
	astFunction := &ASTFunction{
		Package:       pkg,
		Name:          fn.Name.Name,
		Exported:      fn.Name.IsExported(),
		CGOExportName: parseCGOExportName(fn.Doc),
		Params:        parseFieldList(fn.Type.Params),
		Results:       parseFieldList(fn.Type.Results),
		Recv:          parseFieldList(fn.Recv),
	}
	return astFunction
}

func CreateASTFunction(Path string) (functions []*ASTFunction, err error) {
	Path, err = utils.Abs(Path)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var searchFiles []string
	if !utils.IsDir(Path) {
		searchFiles = append(searchFiles, Path)
	} else {
		files, err := utils.GetFilesInDir(Path)
		if err != nil {
			return nil, errors.Trace(err)
		}
		for _, file := range files {
			absPath := filepath.Join(Path, file)
			if suffix := path.Ext(file); !utils.IsDir(absPath) && suffix == ".go" {
				searchFiles = append(searchFiles, absPath)
			}
		}
	}

	for _, file := range searchFiles {
		fns, err := CreateASTFunctionFromFile(file)
		if err != nil {
			return nil, errors.Trace(err)
		}
		functions = append(functions, fns...)
	}
	return
}

func CreateASTFunctionFromFile(filePath string) (functions []*ASTFunction, err error) {
	filePath, err = utils.Abs(filePath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	rawData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", string(rawData), parser.ParseComments)
	if err != nil {
		return nil, errors.Trace(err)
	}

	//ast.Print(fileSet, file)

	pkg := ""
	ast.Inspect(file, func(node ast.Node) bool {
		if pk, ok := node.(*ast.Ident); ok && pkg == "" {
			pkg = pk.Name
		}
		if fn := CreateASTFunctionFromASTNode(node, pkg); fn != nil {
			functions = append(functions, fn)
		}
		return true
	})
	return
}

// deprecated
func CreateASTFromFile(targetPath string) (functions []*ASTFunction, structs []*ASTStruct, err error) {
	targetPath, err = utils.Abs(targetPath)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	rawData, err := ioutil.ReadFile(targetPath)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", string(rawData), parser.ParseComments)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	//ast.Print(fileSet, file)

	pkg := ""
	ast.Inspect(file, func(node ast.Node) bool {
		if pk, ok := node.(*ast.Ident); ok && pkg == "" {
			pkg = pk.Name
		}
		if fn := CreateASTFunctionFromASTNode(node, pkg); fn != nil {
			functions = append(functions, fn)
		}
		if st := CreateASTStructFromASTNode(node, pkg); st != nil {
			structs = append(structs, st)
		}
		return true
	})
	return
}
