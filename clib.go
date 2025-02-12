package clib

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>

// 定义一个函数指针类型，与 libmylib.so 中的 my_function 类型匹配
typedef int (*myFuncType)(int, int);

// 包装函数：传入函数指针和参数，调用该函数指针指向的函数
int callMyFunction(void *f, int arg1, int arg2) {
    return ((myFuncType)f)(arg1, arg2);
}
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func FuncClib() {
	// 指定动态库路径，可以是相对路径或绝对路径
	libPath := C.CString("./libmyfunc.so")
	defer C.free(unsafe.Pointer(libPath))

	// 使用 RTLD_LAZY 模式加载动态库
	handle := C.dlopen(libPath, C.RTLD_LAZY)
	if handle == nil {
		errStr := C.dlerror()
		fmt.Println("dlopen error:", C.GoString(errStr))
		return
	}
	// 程序退出前关闭动态库
	defer C.dlclose(handle)

	// 查找动态库中名为 "add" 的符号
	symName := C.CString("add")
	defer C.free(unsafe.Pointer(symName))
	symbol := C.dlsym(handle, symName)
	fmt.Println("symbol type: ", reflect.TypeOf(symbol).String())
	if symbol == nil {
		errStr := C.dlerror()
		fmt.Println("dlsym error:", C.GoString(errStr))
		return
	}

	// 调用包装函数，通过从 dlsym 得到的函数指针调用 my_function
	arg1 := C.int(42)
	arg2 := C.int(11)

	result := C.callMyFunction(symbol, arg1, arg2)
	fmt.Printf("Result of my_function(%d, %d) = %d\n", int(arg1), int(arg2), int(result))
}
