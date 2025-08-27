package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 1.编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
func addTen(num *int) {
	fmt.Println("1.传入参数为：", num)
	*num += 10
}

// 2.实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2
func sliceEleValueDouble(nums *[]int) {
	fmt.Println("2.传入参数为：", nums)
	for index, value := range *nums {
		(*nums)[index] = value * 2
	}
}

// 3.使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
func startGoroutineDouble() {
	go func() { //打印从1到10的奇数
		for i := 1; i <= 10; i++ {
			time.Sleep(time.Second * 1)
			if i%2 == 1 {
				fmt.Println("打印奇数协程：", i)
			}
		}
	}()

	go func() {
		for i := 2; i <= 10; i++ {
			time.Sleep(time.Second * 1)
			if i%2 == 0 {
				fmt.Println("打印偶数协程：", i)
			}
		}
	}()

	select {}
}

type TaskInfo struct {
	id     int
	idLock sync.Mutex
	chnl   chan string
}

func (this *TaskInfo) numAdd() {
	this.idLock.Lock()
	this.id++
	this.idLock.Unlock()

}

// 4.设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。(没搞明白)
func taskStart(task *TaskInfo) {

	wb := sync.WaitGroup{}

	var tasks = []func(){
		func() {
			start := time.Now()
			for i := 0; i < 100000; i++ {
				task.numAdd()
			}
			timeSpend := time.Since(start)
			fmt.Println("任务1执行数字：", task.id)
			fmt.Println("任务1开始执行，耗时：", timeSpend)
			wb.Done()
		},
		func() {
			start := time.Now()
			for i := 0; i < 100000; i++ {
				task.numAdd()
			}
			timeSpend := time.Since(start)
			fmt.Println("任务2执行数字：", task.id)
			fmt.Println("任务2开始执行，耗时：", timeSpend)
			wb.Done()
		},
		func() {
			start := time.Now()
			for i := 0; i < 100000; i++ {
				task.numAdd()
			}
			timeSpend := time.Since(start)
			fmt.Println("任务3执行数字：", task.id)
			fmt.Println("任务3开始执行，耗时：", timeSpend)
			wb.Done()
		},
	}

	wb.Add(len(tasks))
	for _, fn := range tasks {
		go fn()

	}
	wb.Wait()
	fmt.Println("所有任务执行完毕")

	// fmt.Println("协程开启。。。-->>", task.id)
	// task.numAdd()
	// time.Sleep(time.Second * 3)

	// task.chnl <- "stop"

	// stopTask(task)
}

// 停止协程
func stopTask(task *TaskInfo) {
	fmt.Println("进入关闭协程")
	stopStr := <-task.chnl
	fmt.Println("0停止符号：", stopStr)
	switch stopStr {
	case "start":
		fmt.Println("1停止符号：", stopStr)
	case "stop":
		fmt.Println("2停止符号：", stopStr)
		fmt.Printf("%d协程已经关闭。。。。\n", task.id)
		close(task.chnl)
	default:
		fmt.Println("3停止符号：", stopStr)
	}

}

// 5.定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
type Shape interface {
	Area()
	Perimeter()
}
type Rectangle struct {
	msg string
}
type Circle struct {
	msg string
}

func (this *Rectangle) Area() {
	fmt.Println("调用Rectangle的Area方法，其中msg为：", this.msg)
}

func (this *Rectangle) Perimeter() {
	fmt.Println("调用Rectangle的Perimeter方法，其中msg为：", this.msg)
}
func (this *Circle) Area() {
	fmt.Println("调用Circle的Area方法，其中msg为：", this.msg)
}
func (this *Circle) Perimeter() {
	fmt.Println("调用Circle的Perimeter方法，其中msg为：", this.msg)
}

// 6.使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
type Person struct {
	Name string
	Age  uint
}

type Employee struct {
	EmployeeID string
	person     Person
}

func (this *Employee) PrintInfo() {
	fmt.Println("员工编号：", this.EmployeeID)
	fmt.Println("员工姓名：", this.person.Name)
	fmt.Println("员工年龄：", this.person.Age)
}

// 7.编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
func normalChannel() {
	chnl := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			chnl <- i
			fmt.Println("协程写入数字：", i)
			time.Sleep(time.Second * 1)
		}
		close(chnl)

	}()

	go func() {
		for {
			msg := <-chnl
			fmt.Println("监听到数数字为：", msg)
			time.Sleep(time.Second * 1)
		}
	}()
	select {}
}

// 8.实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
func softChannel() {
	chnl := make(chan int, 5)

	go func() {
		for i := 1; i <= 100; i++ {
			chnl <- i
			fmt.Println("协程写入数字：", i)
			time.Sleep(time.Second * 1)
		}
		close(chnl)

	}()

	go func() {
		for {
			msg := <-chnl
			fmt.Println("监听到数数字为：", msg)
			time.Sleep(time.Second * 1)
		}
	}()

	select {}
}

// 9.编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func syncNumAdd() {
	id := 0
	syncLock := sync.Mutex{}
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("第%d个协程启动\n", i+1)
			for j := 0; j < 1000; j++ {
				syncLock.Lock()
				id++
				syncLock.Unlock()
			}
		}()
	}
	time.Sleep(time.Second * 1)
	fmt.Println("最终累计计算值为：", id)
}

// 10.编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func unLockNumAdd() {
	var id int64 = 0

	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("第%d个协程启动\n", i+1)
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&id, 1)
			}
		}()
	}
	time.Sleep(time.Second * 1)
	fmt.Println("无锁最终累计计算值为：", id)
}

func main() {

	//1.编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
	// num := 10
	// addTen(&num)
	// fmt.Println("1.加上10之后的结果：", num)

	//2.实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2
	// nums := []int{1, 2, 3, 5}
	// sliceEleValueDouble(&nums)
	// fmt.Println("2.每个要素2倍之后的结果：", nums)

	//3.使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	// startGoroutineDouble()

	//4.设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。(不会)
	task := &TaskInfo{
		id:   0,
		chnl: make(chan string, 10),
	}

	taskStart(task)

	// for i := 0; i < 10; i++ {
	// 	go func() {
	// 		taskStart(task)
	// 	}()
	// }
	// 定义函数数组

	// time.Sleep(time.Second)

	//5.定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
	//创建结构体
	// rectangle := &Rectangle{
	// 	msg: "结构体Rectangle",
	// }
	// rectangle.Area()
	// rectangle.Perimeter()
	// circle := &Circle{
	// 	msg: "结构体Circle",
	// }
	// circle.Area()
	// circle.Perimeter()

	//6.使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
	// employee := Employee{
	// 	EmployeeID: "12312",
	// 	person: Person{
	// 		Name: "张三",
	// 		Age:  28,
	// 	},
	// }
	// employee.PrintInfo()

	//7.编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	// normalChannel()

	//8.实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
	// softChannel()

	//9.编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	// syncNumAdd()

	//10.使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	// unLockNumAdd()
}
