// 泛型语法
package main

import "fmt"

func add[T int | float32](a T, b T) T {
	return a + b
}

// 泛型切片
type slice[T int | float64 | string] []T

// 我们把需要用到的类型参数，提前在[]里进行定义，然后在后面实际的变量类型中进行使用，必须要先定义，后使用
type SliceInt []int
type SliceFloat []float64
type SliceString []string

// 泛型map变量
// a := map[string]string 正常写法
type Map[key int | string, value int | float64] map[key]value
type map1 map[int]int
type map2 map[int]float64
type map3 map[string]int
type map4 map[string]float64

// 泛型结构体
type Struct[T int | uint | string] struct {
	Name T
	Age  T
	sex  float64
}

func main() {
	res := add(100, 200)
	println(res)
	//泛型变量实例化
	type MyInt int
	var int1 MyInt = 3
	fmt.Println(int1)
	//切片
	type slice[T int | float32 | string] []T
	//var MySlice slice[int] = []int{1, 2, 3}
	MySlice := slice[int]{1, 2, 3}
	MySliceString := slice[string]{"123", "345", "789"}
	fmt.Println(MySlice, MySliceString)
	//map
	type Map1[key int | string, value string | uint] map[key]value
	//var MyMap Map1[int, string] = map[int]string{
	//	123: 345,
	//	1:   2,
	//}
	MyMap := Map1[int, string]{
		123: "123",
		456: "456",
	}
	fmt.Println(MyMap)
	//结构体
	//var s Struct[int]
	//s.Age = 123
	s := Struct[string]{
		Name: "123",
	}
	fmt.Println(s)
	//泛型嵌套
	type MyStruct[a int | string, p map[a]string] struct {
		Name    string
		content a
		job     p
	}
	m := MyStruct[int, map[int]string]{
		Name:    "zzp",
		content: 123,
		job:     map[int]string{1: "abc"},
	}
	fmt.Println(m)

	type slice1[T int | float32] []T
	type ok[c int | float32, v slice1[c]] struct {
		Name  c
		title v
	}
	//s = ok[int, slice1[int]]{
	//	Name:  123,
	//	title: []int{123, 2},
	//}
	//m := slice1[int]{1, 23, 3}
	//复杂嵌套 Slice2其实就是继承和实现了Slice1，也就是说Slice2的类型参数约束的取值范围，必须是在Slice1的取值范围里
	type slice5[T int | string | float32] []T
	type slice6[T int | string] slice5[T]
	n := slice5[int]{1, 2, 3}
	o := slice6[string]{"zzp", "abc"}
	fmt.Println("pppppppppppppp", n, o)

	//因为它是单一递归继承的，只会检查它的上一级的取值范围是否覆盖
	type Slice1[T bool | float64 | string | int] []T
	type Slice2[T bool | float64 | string] Slice1[T]
	type Slice3[T bool | int] Slice2[T]

}
