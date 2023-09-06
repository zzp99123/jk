// 实现删除切片特定下标元素的方法。
package main

import "fmt"

func Delate(slice []int, index int) []int {
	for i := 0; i < len(slice); i++ {
		if i == index {
			slice = append(slice[:i], slice[i+1:]...)
			//slice = slice[i+1:]
		}
	}
	//fmt.Println(slice)
	return slice
}
func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	//index := 3
	//for i := 0; i < len(s); i++ {
	//	if i == index {
	//		s = s[:i]
	//	}
	//}
	//fmt.Println(s)
	//s = s[:3]
	//fmt.Println(s)
	//q := []int{}
	res := Delate(s, 3)
	fmt.Println(res)
}
