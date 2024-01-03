package main

import (
	"container/list"
	"fmt"
	"math/rand"
)

// 递归
func recur(n int) int {
	if n == 1 {
		return 1
	}
	res := recur(n - 1)
	return n + res
}

// 尾递归
func tailRecur(n, res int) int {
	if n == 0 {
		return res
	}
	return tailRecur(n-1, res+n)
}

// 递归树 斐波那契数列
func fib(n int) int {
	if n == 1 || n == 2 {
		return n - 1
	}
	res := fib(n-1) + fib(n-2)
	return res
}

// 显示栈模拟来调用栈的行为从而将递归转化为迭代形式
func forLoopRecur(n int) int {
	//List列表是一种非连续存储的容器，由多个节点组成，节点通过一些变量记录彼此之间的关系。 列表有多种实现方法，如单链表、双链表等。
	stack := list.New()
	res := 0
	//递 递归调用
	for i := n; i > 0; i-- {
		stack.PushBack(i)
	}
	//归:返回结果
	for stack.Len() != 0 {
		res += stack.Back().Value.(int)
		stack.Remove(stack.Back())
	}
	return res
}

// 时间复杂度 冒泡排序
func bubbleSort(nums []int) int {
	n := 0
	for i := len(nums) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if nums[j] > nums[j+1] {
				tumps := nums[j]
				nums[j] = nums[j+1]
				nums[j+1] = tumps
			}
		}
		n += 3
	}
	return n
}

// 指阶数（细胞分裂）
func exponential(n int) int {
	c, b := 0, 1
	for i := 0; i < n; i++ { // 1 2 3
		for j := 0; j < b; j++ {
			c++ // 1 2 3
		}
		b *= 2 //2 4 8
	}
	return c // 2 * n-1
}

// 指数阶 其递归地一分为二，经过n次分裂后停止
func expRecur(n int) int {
	if n == 1 {
		return 1
	}
	return expRecur(n-1) + expRecur(n-1) + 1
}

// 对数阶（循环实现） 每轮缩减到一半
func logarithmic(n int) int {
	c := 0
	for n > 1 {
		n = n / 2
		c++
	}
	return c
}

// 对数阶（递归实现）
func logRecur(n int) int {
	if n <= 1 {
		return 1
	}
	return logRecur(n/2) + 1
}

// 线性对数阶 二叉树的每一次的操作总数都为n 数共有log2n + 1层 主流排序算法的时间复杂度通常为 例如快速排序、归并排序、堆排序等。
func linearLogRecur(n int) int {
	if n <= 1 {
		return 1
	}
	c := linearLogRecur(n/2) + linearLogRecur(n/2)
	for i := 0; i < n; i++ {
		c++
	}
	return c
}

// 阶乘阶（递归实现）
func factorialRecur(n int) int {
	if n == 0 {
		return 1
	}
	c := 0
	for i := 0; i < n; i++ {
		c += factorialRecur(n - 1)
	}
	return c
}
func randomNumbers(n int) []int {
	c := make([]int, n)
	for i := 0; i < n; i++ {
		c[i] = i + 1
	}
	//打乱切片
	rand.Shuffle(len(c), func(i, j int) {
		c[i], c[j] = c[j], c[i]
	})
	return c

}
func findOne(c []int) int {
	for i := 0; i < len(c); i++ {
		if c[i] == 1 {
			return i
		}
	}
	return -1
}

func main() {
	//res := recur(4)
	//res := tailRecur(3, 5)
	//res := fib(5)
	//res := forLoopRecur(5)
	//n := []int{8, 45, 2, 74, 9, 1, 2, 4}
	//res := bubbleSort(n)
	res := factorialRecur(4) // 4 *3 *2
	fmt.Println(res)
}
