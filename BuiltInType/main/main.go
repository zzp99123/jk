// go内置类型 数组 切片 map
package main

import "fmt"

// 数组
func Array() {
	a := [3]int{1, 2, 3}
	fmt.Printf("a:%v,len:%d,cap:%d\n", a, len(a), cap(a))
	a1 := [3]int{9, 8}
	fmt.Printf("a:%v,len:%d,cap:%d\n", a1, len(a1), cap(a1))
	var a2 [3]int
	fmt.Printf("a:%v,len:%d,cap:%d\n", a2, len(a2), cap(a2))
	fmt.Printf("a[1]:%d", a[1])
}

// 切片
func Slice() {
	s := []int{1, 23, 4}
	fmt.Printf("a:%v,len:%d,cap:%d\n", s, len(s), cap(s))
	//5表示初始化5个参数，6表示初始化容量为6
	s1 := make([]int, 5, 6)
	fmt.Printf("a:%v,len:%d,cap:%d\n", s1, len(s1), cap(s1))
	//s1 = append(s1, 7)
	//fmt.Printf("a:%v,len:%d,cap:%d\n", s1, len(s1), cap(s1))
	s2 := append(s1, 7)
	fmt.Printf("a:%v,len:%d,cap:%d\n", s2, len(s2), cap(s2))
	s2 = append(s2, 8)
	fmt.Printf("a:%v,len:%d,cap:%d\n", s2, len(s2), cap(s2))
	s3 := make([]int, 4)
	fmt.Printf("a:%v,len:%d,cap:%d\n", s3, len(s3), cap(s3))
	s3 = append(s3, 5)
	fmt.Printf("a:%v,len:%d,cap:%d\n", s3, len(s3), cap(s3))

	fmt.Printf("s[2]:%d", s[2])
}

// 子切片 arr[:end]=arr[0:end] arr[start:] = arr[start:len(arr)]
func SubSlice() {
	s := []int{1, 2, 3, 4, 5}
	s1 := s[1:]
	fmt.Printf("a:%v,len:%d,cap:%d\n", s, len(s), cap(s))
	fmt.Printf("a:%v,len:%d,cap:%d\n", s1, len(s1), cap(s1))
}

//核心：共享数组。
//子切片和切片究竟会不会互相影响，就抓住一点：它们是不是还共享数组？
//什么意思？ 就是如果它们结构没有变化，那肯定是共享的；但是结构变化了，就可能不是共享了。 什么情况下结构会发生变化？扩容了。
//所以，切片与子切片，切片作为参数传递到别的方法、结构体里面，任何情况下你要判断是否内存共享，那么 就一个点：有没有扩容

func ShareSlice() {
	s := []int{1, 2, 3, 4, 5}
	s1 := s[2:]
	fmt.Printf("a:%v,len:%d,cap:%d\n", s, len(s), cap(s))    //[1,2,3,4,5]
	fmt.Printf("a:%v,len:%d,cap:%d\n", s1, len(s1), cap(s1)) //[3,4,5]
	s1[0] = 99
	fmt.Printf("a:%v,len:%d,cap:%d\n", s, len(s), cap(s))    //[1,2,99,4,5]
	fmt.Printf("a:%v,len:%d,cap:%d\n", s1, len(s1), cap(s1)) //[99,4,5]
	s1 = append(s1, 199)
	fmt.Printf("a:%v,len:%d,cap:%d\n", s, len(s), cap(s))    //[1,2,99,4,5]
	fmt.Printf("a:%v,len:%d,cap:%d\n", s1, len(s1), cap(s1)) //[99,4,5,199]
	//append了以后在改变s1的值以后 s的值不变
	s1[1] = 1999
	fmt.Printf("a:%v,len:%d,cap:%d\n", s, len(s), cap(s))    //[1,2,99,4,5]
	fmt.Printf("a:%v,len:%d,cap:%d\n", s1, len(s1), cap(s1)) //[99,1999,5,199]

}

func main() {
	//Array()
	//Slice()
	//SubSlice()
	//ShareSlice()
	//切片
	m := map[string]string{
		"key": "1",
	}
	m = make(map[string]string, 4)
	m["key"] = "123"
	fmt.Println(m)
}
