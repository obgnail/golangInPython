package converter

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/obgnail/golangInPython/ast_tool"
	"github.com/obgnail/golangInPython/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const pythonTemplatePath = "converter/python.tmpl"

type RenderFunction struct {
	Name          string
	ArgTypes      string
	ResType       string
	DefArgs       string
	DefRes        string
	CallArgs      string
	NeedReturn    bool
	ResNeedDecode bool
}

func NewRenderFunction(name string, params, results []*ast_tool.ASTField) RenderFunction {
	var ctypeList []string
	var defArgsList []string
	var callArgsList []string
	for _, param := range params {
		types, ok := TypeConverterMap[param.Type]
		if !ok {
			log.Fatalf("do not support such type: %s", param.Type)
		}
		ctypeList = append(ctypeList, types[0])
		defArgsList = append(defArgsList, fmt.Sprintf("%s: %s", param.Name, types[1]))

		pythonArg := param.Name
		if types[1] == "str" {
			pythonArg = fmt.Sprintf(`%s.encode("utf-8")`, pythonArg)
		}
		callArgsList = append(callArgsList, pythonArg)
	}

	resultsLen := len(results)
	resultCType := "null"
	resultPythonType := "None"

	if resultsLen > 1 {
		log.Fatalf("clang function does not support multi result: %v", results)
	} else if resultsLen == 1 {
		result := results[0]
		types, ok := TypeConverterMap[result.Type]
		if !ok {
			log.Fatalf("do not support such type: %s", result.Type)
		}
		resultCType, resultPythonType = types[0], types[1]
	}

	fn := RenderFunction{
		Name:          name,
		ArgTypes:      fmt.Sprintf("[%s]", strings.Join(ctypeList, ", ")),
		ResType:       resultCType,
		DefArgs:       strings.Join(defArgsList, ", "),
		DefRes:        resultPythonType,
		CallArgs:      strings.Join(callArgsList, ", "),
		NeedReturn:    resultsLen == 1,
		ResNeedDecode: resultPythonType == "str",
	}
	return fn
}

type Converter struct {
	sourceDir string
	targetDir string
	goDir     string
	Functions []RenderFunction
}

func NewConverter(sourceDir, targetDir, goDir string, source []*ast_tool.ASTFunction) (c *Converter, err error) {
	if sourceDir, err = utils.Abs(sourceDir); err != nil {
		return nil, errors.Trace(err)
	}
	if targetDir, err = utils.Abs(targetDir); err != nil {
		return nil, errors.Trace(err)
	}

	c = &Converter{sourceDir: sourceDir, targetDir: targetDir, goDir: goDir}
	for _, s := range source {
		if len(s.CGOExportName) != 0 {
			c.Functions = append(c.Functions, NewRenderFunction(s.Name, s.Params, s.Results))
		}
	}
	return
}

func (c *Converter) Compile() error {
	targetFile := filepath.Join(c.targetDir, "main.so")
	cmd := exec.Command(c.goDir, "build", "-buildmode=c-shared", "-o", targetFile)
	cmd.Dir = c.sourceDir
	if _, err := cmd.Output(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (c *Converter) GenPythonFile() (err error) {
	var f *os.File
	file := filepath.Join(c.targetDir, "dll.py")
	f, err = os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()
	tpl := template.Must(template.ParseFiles(pythonTemplatePath))
	if err := tpl.Execute(f, c.Functions); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (c *Converter) Convert() error {
	if len(c.Functions) == 0 {
		return nil
	}
	if err := c.Compile(); err != nil {
		return errors.Trace(err)
	}
	if err := c.GenPythonFile(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func Convert(source, target, goDir string) error {
	functions, err := ast_tool.CreateASTFunction(source)
	if err != nil {
		return errors.Trace(err)
	}
	c, err := NewConverter(source, target, goDir, functions)
	if err != nil {
		return errors.Trace(err)
	}
	if err := c.Convert(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

var (
	// CGO: [ctypes, Python]
	TypeConverterMap = map[string][2]string{
		"*C.char":     {"c_char_p", "str"},
		"*C.schar":    {"c_char_p", "str"},
		"*C.uchar":    {"c_ubyte", "int"},
		"C.short":     {"c_short", "int"},
		"C.ushort":    {"c_ushort", "int"},
		"C.int":       {"c_int", "int"},
		"C.uint":      {"c_uint", "int"},
		"C.long":      {"c_long", "int"},
		"C.ulong":     {"c_ulong", "int"},
		"C.longlong":  {"c_longlong", "int"},
		"C.ulonglong": {"c_ulonglong", "int"},
		"C.float":     {"c_float", "float"},
		"C.double":    {"c_double", "float"},
		"C.size_t":    {"c_size_t", "int"},
	}
)
