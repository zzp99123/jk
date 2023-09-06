// 接口与类型定义
package main

import (
	"fmt"
	"io"
)

// 接口
type List interface {
	Add(index int, val any)
	Append(val any)
	Delete(index int)
}

// 结构体
type LinkedList struct {
	Name string
	Age  int
}

func (this *LinkedList) Add(index int, val any) {
	this.Name = "123"
	//fmt.Println(this.Name)
}

// 结构体与指针 如果声明了一个指针，但是没有赋值，那么它是 nil
func UserList() {
	l1 := LinkedList{}
	l1ptr := &l1
	var l2 LinkedList = *l1ptr
	l2.Age = 18
	l2.Name = "zzp"
	fmt.Println(l1ptr)
	fmt.Println(l2)
	var l3ptr *LinkedList
	fmt.Println(l3ptr)
}

// 方法接收器
func (this LinkedList) UserName(name string) {
	fmt.Printf("%p \n", &this)
	this.Name = name
	//fmt.Printf("%p \n", this)
}

// 传地址直接改变下面的值
func (this *LinkedList) UserAge(age int) {
	//fmt.Printf("%p \n", &this)
	this.Age = age
}

// 结构体自引用
type node struct {
	next *node
}
type Fish struct {
}
type FakeFish struct {
}

func (f Fish) Swim() {
	fmt.Println("swim")
}
func (f FakeFish) FakeSwim() {
	fmt.Println("FakeSwim")
}

// 衍生类型
func UserFish() {
	f1 := Fish{}
	f1.Swim()
	f2 := FakeFish{}
	f2.FakeSwim()
	//类似转换 只能访问字段和方法 不能改变字段与方法
	f3 := Fish(f2)
	f3.Swim()
}

// 结构体实现接口
type NodeList struct {
	head *node
}

func (u NodeList) Add(idx int, val any) {
	panic("implement me")
}

func (u NodeList) Append(val any) {
	panic("implement me")
}
func (u NodeList) Delete(idx int) {
	panic("implement me")
}

// 组合
type Outer struct {
	Inner
}
type Outer1 struct {
	*Inner
}
type Inner struct {
}
type Outer2 struct {
	io.Closer
}

func (o Outer) Name() string {
	return "Outer"
}
func (i Inner) SayHello() {
	fmt.Println("hello" + i.Name())
}
func (i Inner) Name() string {
	return "Inner"
}
func UserOuter() {
	o := Outer{}
	o.SayHello()
}

func main() {
	//结构体初始化
	//l := &LinkedList{
	//	Name: "zzp",
	//}
	//l := LinkedList{
	//	Name: "zzp",
	//}
	//var l LinkedList
	//l.Add(1, "123")
	//fmt.Println(l)
	//结构体字段初始化
	//l := LinkedList{
	//	Name: "zxn",
	//	Age:  18,
	//}
	//fmt.Println(l)
	//UserList()
	//l := LinkedList{
	//	Name: "zzp",
	//	Age:  18,
	//}
	//l.UserName("ok")
	//l.UserAge(23)
	//fmt.Println(l)
	//l := &LinkedList{
	//	Name: "zzp",
	//	Age:  18,
	//}
	//l.UserName("ok")
	//l.UserAge(23)
	//fmt.Println(l)
	//UserFish()
	//n := NodeList{}
	//n.Add(1, "123")
	o := Outer{}
	//a := Outer1{}
	i := Inner{}
	UserOuter()
	o.Name()
	i.SayHello()
	//i.Name()

}
