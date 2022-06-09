package main

import "C"
import (
	"encoding/json"
	"fmt"
)

type Dic struct {
	Msg    string
	Cookie string
}

//export noArgNoRes
func noArgNoRes() {
	fmt.Println("123")
}

//export oneArgNoRes
func oneArgNoRes(b C.int) {
	fmt.Println("printInt:", b)
}

//export twoArgsNoRes
func twoArgsNoRes(a, b C.int) {
	fmt.Println(a + b)
}

//export noArgOneRes
func noArgOneRes() *C.char {
	return C.CString("123")
}

//export oneArgOneRes
func oneArgOneRes(a C.int) *C.char {
	fmt.Println(a)
	return C.CString("123")
}

//export twoArgsOneRes
func twoArgsOneRes(a, b C.int) C.int {
	return a + b
}

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

//export number_add
func number_add(a, b C.int) C.int {
	return a + b
}

//export addFloat
func addFloat(a, b C.float) C.float {
	return a + b
}

//export addDouble
func addDouble(a, b C.double) C.double {
	return a + b
}

func main() {}
