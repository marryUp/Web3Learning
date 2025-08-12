package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// 1.获取数组中只出现一次的数字
func singleNumber(nums []int) int {
	fmt.Println("1.入参为：", nums)
	//	记录出现次数的map
	numCountMap := make(map[int]int)

	//遍历数组并记录数字出现的次数
	for _, num := range nums {
		value, ok := numCountMap[num]
		if ok {
			value++
			numCountMap[num] = value
		} else {
			numCountMap[num] = 1
		}
	}
	// 打印map
	fmt.Println(numCountMap)

	// 找出map中只出现一次的数字
	for key, value := range numCountMap {
		if value == 1 {
			fmt.Println("只出现一次的数字为：", key)
			return key
		} else {
			fmt.Printf("%d出现的次数为：%d \n", key, value)
		}
	}
	return 0
}

// 2.回文数
func isPalindrome(x int) bool {
	fmt.Println("2.入参为：", x)
	// 判断是否为负数
	if x < 0 {
		return false
	}

	//临时数字字符串数组
	strArr := make([]string, 0)

	//转为字符串遍历
	xstring := strconv.Itoa(x)
	for _, value := range xstring {
		strArr = append(strArr, string(value))
	}

	//判断传来的数字是否为0结尾
	if len(strArr) > 1 && (string(xstring[len(xstring)-1]) == "0") {
		return false
	}

	//倒序字符串
	strLast := ""
	for i := len(strArr) - 1; i >= 0; i-- {
		strLast += string(strArr[i])
	}
	//判斷正序和倒序字符串是否相等
	if strLast == xstring {
		return true
	}
	return false
}

// 3.给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
func isValid(s string) bool {

	fmt.Printf("3.入参为:%s,长度为:%d\n", s, len(s))
	//判断字符串长度
	if len(s)%2 != 0 {
		fmt.Println("长度不能为奇数！")
		return false
	}

	stack := []rune{}
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		switch char {
		case '(', '{', '[':
			//如果是左括号，则压入栈
			stack = append(stack, char)
		case ')', '}', ']':
			//如果是右边括号，则判断栈里面最上面那个是否和当前元素是一对，如果是一对就弹出栈
			if len(stack) == 0 || stack[len(stack)-1] != pairs[char] {
				return false
			}
			//弹出栈
			stack = stack[:len(stack)-1]
		}
	}

	//如果栈里面没有元素则证明一对一对的都弹完了
	return len(stack) == 0
}

// 4.编写一个函数来查找字符串数组中的最长公共前缀。如果不存在公共前缀，返回空字符串。例如：strs = ["flower","flow","flight"]
func longestCommonPrefix(strs []string) string {
	fmt.Println("4.传入参数为：", strs)
	//判断要素个数
	if len(strs) == 1 {
		return strs[0]
	}

	//临时存放第一个元素的每个字符
	strMap := make(map[int]string)
	//最小相同个数
	longestNumMin := 0
	for in, str := range strs {
		//每个元素重置计数
		longestNum := 0
		//遍历元素里面的每个字符
		for index, value := range str {
			// 如果map为空则将字符串放入map
			if in == 0 {
				strMap[index] = string(value)
			} else {
				v, ok := strMap[index]
				// 如果不为空则判断当前字符是否存在当前map中
				if ok && string(v) == string(value) {
					longestNum++
				} else {
					// 如果不在map中则停止遍历当前元素
					break
				}

			}

		}
		fmt.Println("longestNum:", longestNum)
		//判断如果longestNum小于longestNumMin则重新赋值
		if in > 1 && longestNum > longestNumMin {

		} else {
			longestNumMin = longestNum
		}

	}

	fmt.Println("最小相同个数：", longestNumMin)
	lastStr := ""
	//循环遍历map拼接相同的字符串
	for i := 0; i < longestNumMin; i++ {
		lastStr += strMap[i]
	}

	return lastStr
}

// 5.将大整数加 1，并返回结果的数字数组。
func plusOne(digits []int) []int {
	fmt.Println("5.传入参数为：", digits)

	len := len(digits)

	//遍历数字数组
	for i := len - 1; i >= 0; i-- {
		//将小于9的数字+1
		if digits[i] < 9 {
			digits[i]++
			fmt.Println("1.返回打印：", digits)
			return digits
		}
		//等于9的就重置为0
		digits[i] = 0

	}

	//处理特殊999等进位的,将第一位置为1,将其他位数的0追加在后面
	digits = append([]int{1}, digits...)
	fmt.Println("2.返回特殊结果，打印：", digits)
	return digits
}

// 6.删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	fmt.Println("6.传入参数为：", nums)
	//判断数组长度
	if len(nums) < 1 {
		return 0
	}
	//声明前一个元素的index
	i := 0
	for j := 1; j < len(nums); j++ {
		//如果后面元素不等于前面元素，计数i++,并将前面一个数改成刚才的后面那个数
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j]
		}
	}
	//返回不重复的总个数
	return i + 1
}

// 7.合并区间,按起点排序是关键！
func merge(intervals [][]int) [][]int {
	fmt.Println("7.传入参数为：", intervals)
	if len(intervals) < 2 {
		return intervals
	}
	// 按起点排序,网上摘抄的
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	resArr := [][]int{}
	//前一个对比数组
	firstArr := intervals[0]
	//后一个对比数组
	secondArr := intervals[1]

	// flag := false

	for j := 1; j < len(intervals); j++ {

		secondArr = intervals[j]
		//如果相交则加入新的数组
		mergeResArr := mergeArr(firstArr, secondArr)
		if mergeResArr != nil {
			// 如果为空直接放入
			if len(resArr) == 0 {
				resArr = append(resArr, mergeResArr)
			} else {
				// 如果不为空则和上一个接着对比
				newArr := mergeArr(resArr[len(resArr)-1], mergeResArr)
				if newArr != nil {
					resArr[len(resArr)-1] = newArr
				} else {
					resArr = append(resArr, mergeResArr)
				}
			}

			//将合并后的结果作为下一个的起始数组
			firstArr = mergeResArr

		} else {
			//如果结果数组中为空，则将两个都放进去
			if len(resArr) == 0 {
				resArr = append(resArr, firstArr)
				resArr = append(resArr, secondArr)
			} else {
				//如果上两个数组合并成功，则不需要放进去
				if !reflect.DeepEqual(resArr[len(resArr)-1], firstArr) {
					resArr = append(resArr, firstArr)
				}
				resArr = append(resArr, secondArr)

			}
			//后一个数组覆盖前一个
			firstArr = secondArr
		}

	}

	fmt.Println("7.返回结果为：", resArr)
	return resArr

}

// 7.1合并两个数组
func mergeArr(arr1 []int, arr2 []int) []int {

	nums1 := []int{}
	nums2 := []int{}
	if arr1[0] > arr2[0] {
		nums1 = arr2
		nums2 = arr1
	} else {
		nums1 = arr1
		nums2 = arr2
	}
	resNums := []int{}

	if nums1[0] <= nums2[0] && nums2[0] <= nums1[1] {
		resNums = append(resNums, nums1[0])
		if nums2[1] < nums1[1] {
			resNums = append(resNums, nums1[1])
		} else {
			resNums = append(resNums, nums2[1])
		}
	} else {
		return nil
	}

	return resNums
}

// 8.两数之和 给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。(简单)
func twoSum(nums []int, target int) []int {
	fmt.Println("8.传入参数为：", nums, target)
	resArr := []int{}
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if (nums[i] + nums[j]) == target {
				resArr = append(resArr, i, j)
			}
		}
	}
	fmt.Println("8.返回结果为：", resArr)
	return resArr
}

func main() {
	// 1.只出现一次的数字
	// nums := []int{1, 2, 3, 2, 1}
	// singleNum := singleNumber(nums)
	// fmt.Println("得到的运行结果,只出现一次的数字为：", singleNum)

	//2.回文数
	// fmt.Println(isPalindrome(22))

	//3.给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效(需要重复掌握)
	// fmt.Println(isValid("()[]{}"))

	//4.编写一个函数来查找字符串数组中的最长公共前缀。如果不存在公共前缀，返回空字符串。
	// strs := []string{"aaa", "aa", "aaa"}
	// fmt.Println(longestCommonPrefix(strs))

	//5.将大整数加 1，并返回结果的数字数组。(需要重复掌握)
	// digits := []int{9, 9, 9}
	// fmt.Println("5.输出结果：", plusOne(digits))

	//6.删除有序数组中的重复项(需要重复掌握)
	// digits := []int{9, 9, 2}
	// fmt.Println("6.输出结果：", removeDuplicates(digits))

	//7.合并区间,数组排序为重点
	// intervals := [][]int{
	// 	{1, 4},
	// 	{0, 4},
	// 	{0, 10},
	// }
	// merge(intervals)

	//8.两数之和 给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
	testArr := []int{3, 4, 6, 7}
	twoSum(testArr, 9)
}
