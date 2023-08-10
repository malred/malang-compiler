// evaluator/builtins.go
package evaluator

import (
	// "fmt"
	"malang/object"
)

var builtins = map[string]*object.Builtin{
	// 传入字符串或数组,求长度
	"len": object.GetBuiltinByName("len"),
	// 输出内容
	"puts": object.GetBuiltinByName("puts"),
	// 传入数组,返回第一个元素
	"first": object.GetBuiltinByName("first"),
	// 传入数组,返回最后一个元素
	"last": object.GetBuiltinByName("last"),
	// 传入数组,返回除了第一个元素以外的所有元素,返回的是新分配的数组
	"rest": object.GetBuiltinByName("rest"),
	// 向数组末尾添加新元素,返回一个新数组
	"push": object.GetBuiltinByName("push"),
	// todo: 文件读写 网络编程 数据库(用原生的"database/sql")
}
