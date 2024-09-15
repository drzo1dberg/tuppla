package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/drzo1dberg/tuppla/go-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/gin-contrib/cors"
)

var DB *gorm.DB
var jwtSecret = []byte("6a7f4d1d9a8b67f32c90dfe4b9a6e9d6f8a7c9d2e6f9a9b2d9f4c7b2d1e9f7a4")

func main () {
   initDatabase()
   router := gin.Default()
   router.Use(cors.Default())

   router.POST("/api/register", func(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Password konnte nicht gehast werden"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	if err := DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Benutzer konnte nicht angelegt werden"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Registrierung erfolgreich"})
   })

   router.POST("/api/login", func(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültige Anmeldedaten"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil{
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültige Anmeldedaten"})
		return	
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 	user.ID,
		"exp": 		time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Token konnte nicht erstellt werden"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})

   })

   router.GET("/api/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello from Go Backend with Gin!!",
			})
		})
		protected := router.Group("/api")
		protected.Use(AuthMiddleware())
		
		protected.GET("/posts", GetPosts)
		protected.POST("/posts", CreatePost)
		
   router.Run(":8080")
}

func initDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Datenbankverbindung fehlgeschlagen", err)
	}

	DB.AutoMigrate(&models.User{}, &models.Post{})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Keine Berechtigung"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			  }
			  return jwtSecret, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx.Set("user_id", uint(claims["user_id"].(float64)))
			  } else {
				if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültiges Token"})
				ctx.Abort()
				}
				return
			  }
			  ctx.Next()
	}
}
func GetPosts(c *gin.Context) {
	var posts []models.Post
	DB.Preload("User").Order("created_at desc").Find(&posts)
	c.JSON(http.StatusOK, gin.H{"posts": posts})
  }
  
  func CreatePost(c *gin.Context) {
	var input struct {
	  Content string `json:"content" binding:"required"`
	}
  
	if err := c.ShouldBindJSON(&input); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
  
	userID, _ := c.Get("user_id")
	post := models.Post{
	  UserID:  userID.(uint),
	  Content: input.Content,
	}
  
	DB.Create(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Beitrag erstellt"})
}