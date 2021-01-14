package main

import (
	"log"
	//"encoding/json"
	"net/http"
	"os"
	"strconv"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	//"github.com/gin-contrib/sessions/cookie"
)

var session_key = "Email"

// Timeline モデルの宣言
type Timeline struct {
	gorm.Model
	Type string  `form:"type"`
	Year int `form:"yearmonth"`
	Month int
	Content string `form:"content" binding:"required"`
}

// User モデルの宣言
type User struct {
	gorm.Model
	Email string `form:"email" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
}
// PasswordEncrypt パスワードをhash化
func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CompareHashAndPassword hashと非hashパスワード比較
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func gormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	DBMS := "mysql"
	USER := os.Getenv("MYSQL_USER")
	PASS := os.Getenv("MYSQL_PASSWORD")
	DBNAME := os.Getenv("MYSQL_DATABASE")
	CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// DBの初期化
func dbInit() {
	db := gormConnect()
	// コネクション解放
	defer db.Close()
	db.AutoMigrate(&Timeline{}) //構造体に基づいてテーブルを作成
	db.AutoMigrate(&User{})
}

// ユーザー登録処理
func createUser(email string, password string) []error {
	passwordEncrypt, _ := PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&User{Email: email, Password: passwordEncrypt}).GetErrors(); err != nil {
		return err
	}
	return nil

}

// ユーザーを一件取得
func getUser(email string) User {
	db := gormConnect()
	var user User
	db.First(&user, "email = ?", email)
	db.Close()
	return user
}

// つぶやき登録処理
func createTimeline(content string) {
	db := gormConnect()
	defer db.Close()
	// Insert処理
	db.Create(&Timeline{Content: content})
}

// つぶやき更新
func updateTimeline(id int, tweetText string) {
	db := gormConnect()
	var tline Timeline
	db.First(&tline, id)
	tline.Content = tweetText
	db.Save(&tline)
	db.Close()
}

// つぶやき全件取得
func getAllTimelines() []Timeline {
	db := gormConnect()

	defer db.Close()
	var tlines []Timeline
	// FindでDB名を指定して取得した後、orderで登録順に並び替え
	db.Order("created_at desc").Find(&tlines)
	return tlines
}

// つぶやき一件取得
func getTimeline(id int) Timeline {
	db := gormConnect()
	var tline Timeline
	db.First(&tline, id)
	db.Close()
	return tline
}

// つぶやき削除
func deleteTimeline(id int) {
	db := gormConnect()
	var tline Timeline
	db.First(&tline, id)
	db.Delete(&tline)
	db.Close()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	//store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob("views/templates/*.html")

	dbInit()

	router.Static("/css", "views/static/css")
	router.Static("/images", "views/static/images")

	// 一覧
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get(session_key) != nil {
			tlines := getAllTimelines()
			c.HTML(200, "index.html", gin.H{"tweets": tlines})
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
		
	})

	// ユーザーログイン画面
	router.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get(session_key) != nil {
			c.Redirect(http.StatusFound, "/")
		} else {
			c.HTML(200, "login.html", gin.H{})
		}
	})

	// ユーザーログイン
	router.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)

		// DBから取得したユーザー
		var user = getUser(c.PostForm("email"))
		if len(user.Email) == 0 {
			// 登録を行います。
			email := c.PostForm("email")
			password := c.PostForm("password")
			// 登録ユーザーが重複していた場合にはじく処理
			if err := createUser(email, password); err != nil {
				c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
			}
		} 
        
		//パスワード(Hash)
		dbPassword := getUser(c.PostForm("email")).Password
		log.Println(dbPassword)
		// フォームから取得したユーザーパスワード
		formPassword := c.PostForm("password")

		// ユーザーパスワードの比較
		if err := CompareHashAndPassword(dbPassword, formPassword); err != nil {
			log.Println("ログインできませんでした")
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"err": err})
			c.Abort()
		} else {
			//text, _ := json.Marshal(user)
			session.Set(session_key, user.Email)
			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"err": err})
				return
			}
			log.Println("ログインできました")
			c.Redirect(http.StatusFound, "/")
		}
	})

	router.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete(session_key)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		c.Redirect(http.StatusFound, "/login")
	})

	a := router.Group("/")
	a.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		log.Println(session.Get(session_key))
		if session.Get(session_key) == nil {
			c.Redirect(200, "/login")
		}
	})

	// つぶやき登録
	router.POST("/new", func(c *gin.Context) {
		var form Timeline
		// バリデーション処理
		if err := c.Bind(&form); err != nil {
			tlines := getAllTimelines()
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"tweets": tlines, "err": err})
			c.Abort()
		} else {
			content := c.PostForm("content")
			createTimeline(content)
			c.Redirect(302, "/")
		}
	})

	// つぶやき詳細
	router.GET("/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		// パラメータから受け取った値をint化
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		tline := getTimeline(id)
		c.HTML(200, "detail.html", gin.H{"tweet": tline})
	})

	// 更新
	router.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		tline := c.PostForm("tweet")
		updateTimeline(id, tline)
		c.Redirect(302, "/")
	})

	// 削除確認
	router.GET("/delete_check/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		tline := getTimeline(id)
		c.HTML(200, "delete.html", gin.H{"tweet": tline})
	})

	// 削除
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		deleteTimeline(id)
		c.Redirect(302, "/")

	})

	router.Run()
}
