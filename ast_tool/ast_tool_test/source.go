package ast_tool_test

import "C"
import (
	"encoding/json"
	"fmt"
)

//export Dic
type Dic struct {
	Msg    string
	Cookie string
}

/*
1
2
3
4
*/
// GetName
//export GetName
func GetName(pyStr *C.char) *C.char {
	fmt.Println(1111, C.GoString(pyStr))
	var s Dic
	err := json.Unmarshal([]byte(C.GoString(pyStr)), &s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
	return C.CString(s.Msg)
}
