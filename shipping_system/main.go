package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var jwtSecret = []byte("supersecretkey")
var users = make(map[string]*User)

func main() {
	r := gin.New()
	// r.Static("/templates", "./templates")
	server.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message":"Welcome to the JWT Authentication Example"})
	})
	
	viewRoutes := r.Group("/view"){
		viewRoutes.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", nil)
		})

		viewRoutes.POST("/register", func(c *gin.Context) {
			var user User
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
				return	
			}
			
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
	
			user.Password = string(hashedPassword)
			users[user.Username] = &user
	
			c.HTML(http.StatusOK, "register.html",gin.H{"message": "User registered successfully"})
		})
	}

	
	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		
		existingUser, ok := users[user.Username]
		if !ok || bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or Password"})
			return
		}

		// generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": existingUser.ID,
			"username": existingUser.Username,
		})

		jwtToken, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": jwtToken})
	})

	r.Run(":8080")
}