package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"SaaS/initializers"
	"SaaS/models"
	"SaaS/services"
)

func SignUp(c *gin.Context) {
	var user models.User

	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения тела запроса"})
		return
	}
	if !services.EmailVerification(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "введён некорректный адрес электронной почты"})
		return
	}
	if len(user.Password_hash) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "минимальная длина пароля - 6 символов"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password_hash), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка хэширования пароля"})
		return
	}
	user.Password_hash = string(hash)

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		if result.Error.Error() == "constraint failed: UNIQUE constraint failed: users.email (2067)" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь с таким адресом электронной почты уже зарегистрирован"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка регистрации"})
		return
	}
	message := fmt.Sprintf("пользователь %s успешно зарегистрирован", user.Name)
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

func SignIn(c *gin.Context) {
	var userToCheck models.User
	var user models.User

	if c.Bind(&userToCheck) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения тела запроса"})
		return
	}

	initializers.DB.First(&user, "email = ?", userToCheck.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь с таким адресом электронной почты не найден"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(userToCheck.Password_hash)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "введён неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка создания токена"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", tokenString, 3600*24*7, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "успешный вход",
	})
}

func SignOut(c *gin.Context) {
	_, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователя не существует"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", "", -1, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "успешный логаут",
	})
}

func Me(c *gin.Context) {
	user := services.GetUser(c)
	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"name":  user.Name,
	})
}
