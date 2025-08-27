package main

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB
var err error
var dbType = "sqlite3" //sqlite3||mysql
var url string

// 初始化数据库连接
func init() {
	//默认sqlite3
	url := "test.db"
	if dbType == "mysql" {
		url = "root:root@tcp(192.168.80.1:3306)/test?charset=utf8mb4&parseTime=True"
	}

	//连接mysql
	db, err = sqlx.Connect(dbType, url)
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	fmt.Println("mysql 连接成功")
}

type User struct {
	Id   uint
	Name string `db:"name"` //起别名映射数据库字段
	Age  int
}

func createTable() {
	// 数据表创建
	sqlCreate := `
	CREATE TABLE IF NOT EXISTS user (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT,
	    age INTEGER
	);
	`
	_, err = db.Exec(sqlCreate)
	if err != nil {
		fmt.Println("表创建失败")
		return
	}
	fmt.Println("sqlite3 user表创建成功")
}
func createMysqlTable() {
	// 数据表创建
	sqlCreate := `CREATE TABLE user (
	    id int4  NOT NULL AUTO_INCREMENT  COMMENT 'id主键',
		name varchar(255) comment '名称',
	    age int4 comment '年龄',
		PRIMARY KEY (id) USING BTREE
	)ENGINE=InnoDB DEFAULT CHARSET=utf8 comment='用户表'
	`
	_, err = db.Exec(sqlCreate)
	if err != nil {
		fmt.Println("表创建失败：", err)
		return
	}
	fmt.Println("mysql user表创建成功")
}

// 基本增删改查练习
func baseSqlUse() {
	// 1.表创建
	// createTable()
	// if dbType == "mysql" {
	// 	createMysqlTable()
	// }else{
	// 	createTable()
	// }

	// 2.数据插入
	// sqlInsert := `
	// insert into user(name, age) values ("张三", 18)
	// `
	// _, err = db.Exec(sqlInsert)
	// if err != nil {
	// 	fmt.Println("用户张三信息插入失败：", err)
	// 	return
	// }
	// fmt.Println("用户张三信息插入成功")

	//2.1数据插入对象映射
	// sqlInsertByMapping := `
	// insert into user(name, age) values (:name, :age)
	// `
	// user := User{
	// 	Name: "wangwu",
	// 	Age:  19,
	// }
	// db.NamedExec(sqlInsertByMapping, &user)
	// db.NamedExec(sqlInsertByMapping, map[string]interface{}{
	// 	"name": "zhaoliu",
	// 	"age":  17,
	// })

	// fmt.Println("wangwu插入成功")

	// 3删除
	// sqlDelete := `
	// 	delete from user where id = ?
	// `
	// db.Exec(sqlDelete, 1)
	// fmt.Println("删除成功")

	// 4.修改
	// sqlUpdate := `
	// update user set age = ? where name = ?
	// `
	// db.Exec(sqlUpdate, 10, "wangwu")
	// fmt.Println("修改成功")

	// 5.查询单条
	// sqlSelectOne := `
	// select id,name, age from user where name = ? and age = ? order by id desc limit 1
	// `
	// user := User{}
	// db.Get(&user, sqlSelectOne, "wangwu", 10)
	// fmt.Println("查询单条成功：", user)

	// 5.1查询多条
	sqlSelectMany := `
	select id,name, age from user where name = ? and age = ? order by id desc
	`
	users := []User{}
	db.Select(&users, sqlSelectMany, "wangwu", 10)
	fmt.Println("查询多条成功：", users)

	// 5.2查询多条映射条件
	sqlQueryMany := `
	select id,name, age from user where age > :age order by id desc
	`
	user := User{
		Name: "lisi",
		Age:  10,
	}
	rows, err := db.NamedQuery(sqlQueryMany, user)
	if err != nil {
		fmt.Println("查询失败：", err)
		return
	}
	for rows.Next() {
		var u User
		rows.StructScan(&u)
		fmt.Println(u)
	}
}

// 事务的使用
func useTransaction() (err error) {
	// 开启事务
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println("开启事务失败")
		return
	}

	// 最后关闭事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("事务回滚")
			panic(r)
		} else if err != nil {
			tx.Rollback()
			fmt.Println("事务回滚")
		} else {
			tx.Commit()
			fmt.Println("事务提交")
		}
	}()

	// 处理事务
	//2.1数据插入对象映射
	sqlInsertByMapping := `
	insert into user(name, age) values (:name, :age)
	`
	user := User{
		Name: "wangwu11",
		Age:  11,
	}
	rs, err := tx.NamedExec(sqlInsertByMapping, &user)
	if err != nil {
		// tx.Rollback()
		// fmt.Print("插入失败，事务回滚：", err)
		return err
	}
	n, err := rs.RowsAffected()
	if n < 1 {
		return errors.New("Insert error")
	}
	// 更新操作
	sqlUpdate := `
		update user set age = :age where id = :id 
	`
	userUp := User{
		Id:   3,
		Name: "wangwu22",
		Age:  22,
	}

	ures, err := tx.NamedExec(sqlUpdate, &userUp)
	if err != nil {
		// tx.Rollback()
		// fmt.Print("更新失败，事务回滚：", err)
		return err
	}
	num, err := ures.RowsAffected()
	if num != 1 {
		fmt.Println("更新数量为0,返回事务报错")
		return errors.New("Update error")
	}
	return nil

}

// --------------------------上面是练习，下面是作业----------------------------------
// 1员工结构体
type Employee struct {
	Id         int
	Name       string
	Department string
	Salary     float32
}

// 1sqlite3上创建员工信息表
func createEmployeeTable() {
	sqlCreate := `
	create table if not exists employees(
	 
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT,
	    department text,
		salary float
	
	)
	`
	_, err = db.Exec(sqlCreate)
	if err != nil {
		fmt.Println("创建员工信息表失败")
		return
	}
	fmt.Println("员工信息表创建成功")

}

// 2book结构体
type Book struct {
	Id     uint
	Title  string
	Author string
	Price  float32
}

// 2创建books表，包含字段 id 、 title 、 author 、 price 。
func createBookTable() {
	createSql := `
		create table if not exists books(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	    title TEXT,
	    author text,
		price float
		)
	`
	_, err = db.Exec(createSql)
	if err != nil {
		fmt.Println("创建books表失败：", err)
		return
	}
	fmt.Println("创建books成功。。。")
}

func main() {

	defer db.Close()
	//基本增删改查练习
	// baseSqlUse()

	// 事务的使用
	// useTransaction()

	//--------------------------上面是练习，下面是作业----------------------------------
	//0.默认连接数据库，创建员工信息表
	// 1.createEmployeeTable()
	// 新增员工信息
	// insertSql := `
	// insert into employees(name, department, salary) values (:name,:department,:salary)
	// `
	// employeeSlice := []Employee{
	// 	{
	// 		Name:       "zhangsan",
	// 		Department: "技术部",
	// 		Salary:     8000,
	// 	},
	// 	{
	// 		Name:       "lisi",
	// 		Department: "技术部",
	// 		Salary:     8800,
	// 	},
	// 	{
	// 		Name:       "wangwu",
	// 		Department: "技术部",
	// 		Salary:     10000,
	// 	},
	// }

	// //批量插入
	// for emp := range employeeSlice {

	// 	db.NamedExec(insertSql, emp)
	// }

	//1.1编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	employeeNew := Employee{
		Department: "技术部",
	}
	queryListSql := `
		select * from employees where department = :department
	`
	rows, err := db.NamedQuery(queryListSql, employeeNew)
	if err != nil {
		fmt.Println("查询技术部员工信息报错：", err)
		return

	}

	for rows.Next() {
		var emp = Employee{}
		rows.StructScan(&emp)
		fmt.Println("查询到技术部员工有：", emp)
	}

	//1.2编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	selectMaxSalarySql := `
		select * from employees order by salary desc limit 1
	`
	var maxSalaryEmployee = Employee{}
	db.Get(&maxSalaryEmployee, selectMaxSalarySql)
	fmt.Println("最高工资员工信息为：", maxSalaryEmployee)

	//2假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	// 要求 ：
	// 2.1定义一个 Book 结构体，包含与 books 表对应的字段。
	// createBookTable()

	// 2.2编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	//先录入数据
	// booksInsert := []Book{
	// 	{
	// 		Title:  "钢铁是怎样练成的",
	// 		Author: "ml",
	// 		Price:  88.8,
	// 	}, {
	// 		Title:  "穷爸爸富爸爸",
	// 		Author: "zhangsan",
	// 		Price:  38.18,
	// 	},
	// }
	// booksInsertSql := `
	// insert into books(title,author,price) values(:title,:author,:price)
	// `
	// db.NamedExec(booksInsertSql, booksInsert)

	// 再查询数据
	books := []Book{}

	selectSql := `
	select * from books where price > ?
	`
	db.Select(&books, selectSql, 50)

	fmt.Println("大于50的书籍有：", books)

}
