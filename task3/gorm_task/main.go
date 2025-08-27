package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 定义用户结构体
type User struct {
	Id      uint `gorm:"primaryKey"`
	Name    string
	PostNum int    `gorm:"default:0"`
	Posts   []Post `gorm:"foreignKey:UserId"` // 一对多关系
}

// 定义文章结构体
type Post struct {
	Id            uint `gorm:"primaryKey"`
	UserId        uint // 外键字段
	Title         string
	CommentStatus string    `gorm:"default:'无评论'"`     // 文章的评论状态字段，默认值为"无评论"
	CommentNum    int       `gorm:"default:0"`         // 文章的评论数量统计字段
	Comments      []Comment `gorm:"foreignKey:PostId"` // 一对多关系
}

// 定义评论结构体
type Comment struct {
	Id      uint `gorm:"primaryKey"`
	PostId  uint // 外键字段
	Content string
}

func ConnectDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		fmt.Println("SQlite连接失败：", err)
		return nil
	}
	return db
}

// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	// 创建文章前，更新用户的文章数量
	// tx.Model(&User{}).Where("id = ?", post.UserId).Update("post_num", gorm.Expr("post_num + ?", 1))
	err = tx.Model(&User{}).Where("id = ?", post.UserId).Update("post_num", gorm.Expr("post_num + ?", 1)).Error
	if err != nil {
		fmt.Println("更新用户文章数量失败：", err)
	}
	return
}

// 为 Comment 模型添加一个钩子函数，在评论创建后自动更新文章的评论数量统计字段。
func (comment *Comment) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("评论创建之后：", comment.Id)
	fmt.Println("评论创建之后：", comment.PostId)
	return tx.Debug().Model(&Post{}).Where("id = ?", comment.PostId).Update("comment_num", gorm.Expr("comment_num + ?", 1)).Update("comment_status", "").Error
}

// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (comment *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("文章ID:", comment.PostId)
	// 查询文章的评论数量
	var count int64
	err = tx.Model(&Comment{}).Where("post_id = ?", comment.PostId).Count(&count).Error
	if err != nil {
		fmt.Println("查询文章评论数量失败：", err)
		return err
	}
	fmt.Println("文章ID:", comment.PostId, "当前评论数量:", count)
	if count == 0 {
		// 如果评论数量为0，更新文章的评论状态为"无评论"
		err = tx.Model(&Post{}).Where("id = ?", comment.PostId).UpdateColumns(map[string]interface{}{"comment_status": "无评论", "comment_num": 0}).Error
		if err != nil {
			fmt.Println("更新文章评论状态失败：", err)
			return err
		}
		fmt.Println("当前文章暂无评论:", comment.PostId)
	} else {
		tx.Debug().Model(&Post{}).Where("id = ?", comment.PostId).UpdateColumn("comment_num", count)
	}
	return
}

func main() {
	db := ConnectDb()
	if db == nil {
		fmt.Print("数据库连接失败，程序退出")
	}
	fmt.Println("SQlite连接成功")
	// 1.1创建表
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 1.2插入数据
	// db.Create(&User{Name: "ML", Posts: []Post{
	// 	{Title: "First Post", Comments: []Comment{{Content: "Great post!"}, {Content: "Thanks for sharing!"}}},
	// 	{Title: "Second Post", Comments: []Comment{{Content: "Interesting read!"}, {Content: "I learned a lot!"}}},
	// }})

	// db.Create(&User{Name: "Ml", PostNum: 1, Posts: []Post{
	// 	{

	// 		CommentNum:    1,
	// 		CommentStatus: "",
	// 		Comments: []Comment{
	// 			{
	// 				PostId:  1,
	// 				Content: "第一条评论，初始",
	// 			},
	// 		},
	// 	},
	// },
	// })
	// db.Create(&Comment{PostId: 6, Content: "第三条评论"})
	// fmt.Println("插入评论数据成功")

	// 2.查询数据
	// 	2.1编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	// users := []User{}
	// err := db.Preload("Posts").Preload("Posts.Comments").Find(&users, "name = ?", "ML").Error
	// if err != nil {
	// 	fmt.Println("查询用户及其文章和评论失败：", err)
	// 	return
	// }
	// fmt.Println("查询到的用户记录为：", users)
	// for _, user := range users {
	// 	fmt.Printf("用户ID: %d, 用户名: %s\n", user.Id, user.Name)
	// 	for _, post := range user.Posts {
	// 		fmt.Printf("  文章ID: %d, 标题: %s\n", post.Id, post.Title)
	// 		for _, comment := range post.Comments {
	// 			fmt.Printf("    评论ID: %d, 内容: %s\n", comment.Id, comment.Content)
	// 		}
	// 	}
	// }

	// 2.2编写Go代码，使用Gorm查询评论数量最多的文章信息。
	// var maxPost = Post{}
	// db.Debug().Preload("Comments").Order("(select count(*) from comments where comments.post_id = posts.id) DESC").First(&maxPost)
	// fmt.Println("查询到的文章记录为：", maxPost)

	//3.钩子函数
	// 	3.1为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	// BeforeCreate 方法已经定义在 Post 结构体中的钩子函数，下面为测试方法
	// db.Create(&Post{UserId: 1, Title: "第3个文章 Post"})
	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
	// db.Create(&Comment{PostId: 3, Content: "这是你的评论4"})

	// comment := &Comment{
	// 	Id:     11,
	// 	PostId: 3,
	// }
	// db.Delete(&comment)

}
