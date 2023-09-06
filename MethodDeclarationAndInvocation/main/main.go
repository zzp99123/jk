// 方法声明与调用
package main

import "fmt"

func Func0(name string) string {
	return "hello world"
}
func Func1(a, b, c int, str string) (string, error) {
	return "", nil
}
func Func2(a int, b int) (str string, err error) {
	str = "hello"
	return
}
func Func3(a int, b int) (str string, err error) {
	res := "hello"
	return res, err
}

// 方法作为变量
func Func4() {
	myFunc := Func3
	_, _ = myFunc(5, 6)
}

// 局部方法
func Func5() {
	fn := func(name string) string {
		return "hello" + name
	}
	ok := fn("abc")
	fmt.Println(ok)
}

// 方法作为返回值
func Func6() func(name string) string {
	return func(name string) string {
		return "hello" + name
	}
}

// 匿名方法
func Func7() {
	fn := func() string {
		return "hello world"
	}()
	fmt.Println(fn)
}

// 闭包 :闭包如果使用不当可能会引起内存泄露的问题。即一个对象被闭包引用的话，它是不会被垃圾回收的。 记住这个结论，你后面面试用得上
func Closure(name string) func() string {
	return func() string {
		return "hello" + name
	}
}

// 不定参数
func Indefinite(name string, alias ...string) {
	if len(alias) > 0 {
		fmt.Println(alias[0])
	}
}

// defer类似于栈 先进后出，也就是，先定义的后执行，后定义的先执行
func Defer() {
	defer func() {
		fmt.Println("第一个defer")
	}()
	defer func() {
		fmt.Println("第二个defer")
	}()
}

// defer与闭包
func DeferClosure() {
	i := 0
	defer func() {
		print(i)
	}()
	i = 1
}
func DeferClosureVal() {
	i := 0
	defer func(val int) {
		print(val)
	}(i)
	i = 1
}

// defer 修改返回值 如果是带名字的返回值，那么可以修改这个返回值，否则不能修改
func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}
func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}

type MySturt struct {
	Name string
}

func DeferReturnV2() *MySturt {
	a := &MySturt{Name: "jerry"}
	defer func() {
		a.Name = "tom"
	}()
	return a
}

// 9还在循环中 到10跳出循环进入defer中
func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

// 先进后出
func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}
func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}
func main() {
	//n := Func0("abc")
	//fmt.Println(n)
	//a, c := Func1(1, 2, 3, "abc")
	//fmt.Println(a, c)
	//a, err := Func2(1, 2)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(a)
	//n, err := Func3(2, 3)
	//fmt.Println(n, err)
	//n, err := Func2(1, 2)
	//println(n, err)
	////忽略返回值
	//_, _ = Func2(1, 2)
	//n, _ = Func2(1, 2)
	////str是新变量 需要:=
	//str, _ := Func2(1, 2)
	//println(str)
	//Func5()
	//fn := Func6()
	//fmt.Println(fn("zzp"))
	//Func7()
	//Indefinite("zzp", "daming", "xixue")
	//Defer()
	//DeferClosure()
	//DeferClosureVal()
	//res := DeferReturn()
	//res := DeferReturnV1()
	//res := DeferReturnV2()
	//fmt.Println(res)
	//DeferClosureLoopV1()
	//DeferClosureLoopV2()
	DeferClosureLoopV3()
}
