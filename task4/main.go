package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/ml_test/task4_model/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

var db *gorm.DB
var err error
var logger *logrus.Logger

// 1.数据库连接初始化
func init() {
	//sqlite3数据库连接
	db, err = gorm.Open(sqlite.Open("blogdb.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败")
	}
	fmt.Println("数据库连接成功。。。")
	// 创建同步表
	createTables()
	// 声明日志初始化
	logger = initLogger()
}

// 2.创建相关表
func createTables() {
	db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
}

// 3.用戶注册和认证
func Register(c *gin.Context) {
	var user model.User
	//接收参数并赋值到user
	err := c.ShouldBindJSON(&user)
	//信息绑定判断
	dealBindInfo(err, c)

	userName := user.Username
	//查询当前名称是否已经存在
	userList := []model.User{}

	db.Debug().Unscoped().Find(&userList, "username = ?", userName)

	fmt.Println("查询到的用户有：", userList)
	// 如果已经存在，则返回报错
	if len(userList) > 0 {
		logger.Info("当前用户名称已存在：", userName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "当前用户名称已存在：" + userName})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Info("密码加密失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user.Password = string(hashedPassword)

	//创建用户
	err = db.Omit("created_at", "updated_at", "deleted_at").Create(&user).Error
	if err != nil {
		logger.Info("用户创建失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})

}

// 3登录和认证
func Login(c *gin.Context) {
	var user model.User

	//信息绑定结构体
	err = c.ShouldBindJSON(&user)
	//信息绑定判断
	dealBindInfo(err, c)

	//查出库中的当前用户信息
	var userDb model.User
	err = db.Debug().Unscoped().Where("username = ?", user.Username).First(&userDb).Error
	if err != nil {
		logger.Info("用户名或者密码有误")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或者密码有误"})
		return
	}

	// 两个密码加密后对比?
	err = bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(user.Password))
	if err != nil {
		logger.Info("用户名或者密码有误!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或者密码有误!"})
		return
	}

	//生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userDb.Id,
		"username": userDb.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	fmt.Println("生成的tokenStr:", token)

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		logger.Info("生成获取token失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取token失败"})
		return
	}
	//将token返回给前端
	c.JSON(http.StatusCreated, gin.H{"token": tokenString})

}

// 通过token获取用户信息
func GetUserFromToken(tokenString string) (*model.User, error) {

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	fmt.Println("kaishi第一步。。。。。err:", err)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token无效")
	}
	fmt.Println("kaishi第二步。。。。。")
	// 获取 claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Info("无法获取claims")
		return nil, fmt.Errorf("无法获取claims")
	}
	fmt.Println("claims解析token信息为：", claims)

	// 获取用户ID
	userIdFloat, ok := claims["id"].(float64)
	if !ok {
		logger.Info("token中没有用户id")
		return nil, fmt.Errorf("token中没有用户id")
	}
	userId := uint(userIdFloat)

	// 查询数据库
	var user model.User
	if err := db.Unscoped().First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 4.1文章创建
func createPost(c *gin.Context) {

	//验证是否是否登录
	token := c.GetHeader("Authorization")
	//判断token
	dealTokenEmpty(token, c)
	user, err := GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// c.JSON(http.StatusBadRequest, gin.H{"error": "请先登录认证"})
		return
	}
	userId := user.Id
	var post model.Post
	err = c.ShouldBindJSON(&post)
	//信息绑定判断
	dealBindInfo(err, c)

	if post.Title == "" || post.Content == "" {
		logger.Info("文章的标题和内容必填")
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章的标题和内容必填"})
		return
	}

	post.UserId = userId
	fmt.Println("要创建的文章信息为：", post)
	err = db.Debug().Create(&model.Post{
		Title:   post.Title,
		Content: post.Content,
		UserId:  post.UserId,
	}).Error
	if err != nil {
		fmt.Println("创建报错：", err)
		logger.Info("文章信息创建失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章信息创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": "ok"})

}

// 4.2文章的获取

func getPosts(c *gin.Context) {

	postTitle := c.Query("postTitle")

	posts := []model.Post{}
	if postTitle != "" {
		db.Debug().Unscoped().Find(&posts, "title = ?", postTitle)
	} else {
		db.Debug().Find(&posts)
	}
	fmt.Println("获取到的文章有：", posts)

	c.JSON(http.StatusCreated, gin.H{"data": posts})

}

// 4.3文章的修改
func updatePost(c *gin.Context) {

	//验证是否登录
	token := c.GetHeader("Authorization")
	//判断token
	dealTokenEmpty(token, c)

	//通过token获取当前登录用户信息
	user, err := GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// c.JSON(http.StatusBadRequest, gin.H{"error": "请先登录认证"})
		return
	}
	userId := user.Id
	var post model.Post
	err = c.ShouldBindJSON(&post)
	//信息绑定判断
	dealBindInfo(err, c)

	fmt.Println("绑定信息为：", post)
	if post.Id == 0 {
		logger.Info("主键id不能为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "主键id不能为空"})
		return
	}

	// 判断只能修改自己的文章信息
	postDb := model.Post{}
	db.Debug().Unscoped().Find(&postDb, "id = ?", post.Id)
	if postDb.Id == 0 || postDb.UserId != userId {
		logger.Info("只能修改自己的文章信息")
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能修改自己的文章信息"})
		return
	}

	db.Debug().Unscoped().Model(&postDb).Updates(&post)

	c.JSON(http.StatusCreated, gin.H{"code": "ok"})

}

// 4.4文章的删除
func deletePost(c *gin.Context) {
	//验证是否登录
	token := c.GetHeader("Authorization")
	//判断token
	dealTokenEmpty(token, c)
	//通过token获取当前登录用户信息
	user, err := GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// c.JSON(http.StatusBadRequest, gin.H{"error": "请先登录认证"})
		return
	}
	userId := user.Id
	var post model.Post
	err = c.ShouldBindJSON(&post)
	//信息绑定判断
	dealBindInfo(err, c)

	fmt.Println("绑定信息为：", post)
	if post.Id == 0 {
		logger.Info("主键id不能为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "主键id不能为空"})
		return
	}

	// 判断只能修改自己的文章信息
	postDb := model.Post{}
	db.Debug().Unscoped().Find(&postDb, "id = ?", post.Id)
	if postDb.Id == 0 || postDb.UserId != userId {
		logger.Info("只能删除自己的文章信息")
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能删除自己的文章信息"})
		return
	}

	db.Debug().Model(&postDb).Where("id = ?", post.Id).Delete(&post)
	fmt.Println("删除成功：", post.Id)
	c.JSON(http.StatusCreated, gin.H{"code": "ok"})

}

// 5.1添加评论
func addComment(c *gin.Context) {
	// 判断是否认证过
	b, user := haveValid(c)
	if !b {
		logger.Info("暂未登录认证")
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先注册并登录认证"})
		return
	}
	userId := user.Id

	var comment model.Comment
	err = c.ShouldBindJSON(&comment)
	//信息绑定判断
	dealBindInfo(err, c)

	if comment.PostId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要评论的文章"})
		return
	}
	if comment.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
		return
	}
	comment.UserId = userId

	// 添加评论
	db.Debug().Create(&comment)

	c.JSON(http.StatusCreated, gin.H{"code": "ok"})
}

// 5.1查看评论
func listComment(c *gin.Context) {

	// comment.UserId = userId
	postId := c.Query("postId")
	if postId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要评论的文章"})
		return
	}
	var comments []model.Comment
	// 添加评论
	db.Debug().Find(&comments, "post_id = ? ", postId)

	c.JSON(http.StatusCreated, gin.H{"data": comments})
}

// 6.1 返回报错各种封装
// 登录认证校验
func dealTokenEmpty(token string, c *gin.Context) {
	if token == "" {
		logger.Info("暂未登录认证")
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先登录认证"})
		return
	}
}

// 实体绑定校验
func dealBindInfo(err error, c *gin.Context) {
	if err != nil {
		logger.Info("信息绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": "信息绑定失败"})
		return
	}
}

// 判断用户是否认证过
func haveValid(c *gin.Context) (bool, *model.User) {
	//验证是否登录
	fmt.Println("获取到的参数为：", c)
	token := c.GetHeader("Authorization")
	fmt.Println("获取到的 token:", token)
	if token == "" {
		return false, nil
	}
	//通过token获取当前登录用户信息
	user, err := GetUserFromToken(token)
	if err != nil {
		return false, nil
	}
	if user.Id != 0 {
		return true, user
	}
	return false, nil
}

//6.2使用日志库

func initLogger() *logrus.Logger {
	logger := logrus.New()

	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)

	// 设置输出格式为 JSON
	logger.SetFormatter(&logrus.JSONFormatter{})

	// 设置输出到文件和控制台
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Info("无法打开日志文件，使用默认输出")
	}

	return logger
}

func main() {

	fmt.Println("task4 启动。。。。")

	//2.数据库设计与模型定义
	//init()//默认执行

	// 启动监听服务
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	defer func() {
		logger.Info("博客后台启动成功。。。")
		fmt.Println("博客后台启动成功。。。")
		router.Run() // 监听并在 0.0.0.0:8080 上启动服务
	}()

	//3.1注册接口
	router.POST("/register", func(c *gin.Context) {
		Register(c)
	})

	// 3.2登录和认证接口
	router.POST("/login", func(c *gin.Context) {
		Login(c)
	})

	/**4文章管理功能
		实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
	实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
	实现文章的更新功能，只有文章的作者才能更新自己的文章。
	实现文章的删除功能，只有文章的作者才能删除自己的文章。
	*/

	//4.1实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。

	router.POST("/post/create", func(c *gin.Context) {
		createPost(c)
	})

	//4.2实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
	router.GET("/post/getPosts", func(c *gin.Context) {
		getPosts(c)
	})
	//4.3实现文章的更新功能，只有文章的作者才能更新自己的文章。
	router.POST("/post/updatePost", func(c *gin.Context) {
		updatePost(c)
	})
	//4.4实现文章的删除功能，只有文章的作者才能删除自己的文章。
	router.POST("/post/deletePost", func(c *gin.Context) {
		deletePost(c)
	})

	/**5.评论功能
	实现评论的创建功能，已认证的用户可以对文章发表评论。
	实现评论的读取功能，支持获取某篇文章的所有评论列表。
	*/
	//5.1实现评论的创建功能，已认证的用户可以对文章发表评论。
	router.POST("/comment/addComment", func(c *gin.Context) {
		addComment(c)
	})
	//5.2实现评论的读取功能，支持获取某篇文章的所有评论列表。
	router.GET("/comment/listComment", func(c *gin.Context) {
		listComment(c)
	})
}
