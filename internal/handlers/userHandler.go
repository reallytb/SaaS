package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"SaaS/internal/dto"
	"SaaS/internal/initializers"
	"SaaS/internal/models"
	"SaaS/pkg/utils"
)

// SignUp godoc
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} dto.SuccessResponse "Пользователь зарегистрирован"
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос или пользователь уже существует"
// @Router /signup [post]
// @Example request
//
//	{
//	    "email": "user@example.com",
//	    "password_hash": "password123",
//	    "name": "Иван Иванов"
//	}
//
// @Example response 201
//
//	{
//	    "success": true,
//	    "message": "пользователь Иван Иванов успешно зарегистрирован"
//	}
//
// @Example response 400
//
//	{
//	    "success": false,
//	    "error": "пользователь с таким адресом электронной почты уже зарегистрирован"
//	}
func SignUp(c *gin.Context) {
	var user models.User

	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка чтения тела запроса",
		})
		return
	}
	if !utils.EmailVerification(user.Email) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "введён некорректный адрес электронной почты",
		})
		return
	}
	if len(user.PasswordHash) < 6 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "минимальная длина пароля - 6 символов",
		})
		return
	}
	if user.Name == "" || len(user.Name) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "имя пользователя не может быть пустым",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка хэширования пароля",
		})
		return
	}
	user.PasswordHash = string(hash)

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		if result.Error.Error() == "constraint failed: UNIQUE constraint failed: users.email (2067)" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Error:   "пользователь с таким адресом электронной почты уже зарегистрирован",
			})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка регистрации",
		})
		return
	}
	message := fmt.Sprintf("пользователь %s успешно зарегистрирован", user.Name)
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}

// SignIn godoc
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает JWT токен в cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.SuccessResponse "Успешный вход"
// @Failure 400 {object} dto.ErrorResponse "Неверные учетные данные"
// @Router /signin [post]
// @Example request
//
//	{
//	    "email": "user@example.com",
//	    "password_hash": "password123"
//	}
//
// @Example response 200
//
//		{
//		"success": true,
//	    "message": "успешный вход"
//	    "data": {
//	        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	    }
//		}
//
// @Example response 400
//
//	{
//	    "success": false,
//	    "error": "введён неверный пароль"
//	}
//
// @Security BearerAuth
func SignIn(c *gin.Context) {
	var userToCheck models.User
	var user models.User

	if c.Bind(&userToCheck) != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка чтения тела запроса",
		})
		return
	}

	initializers.DB.First(&user, "email = ?", userToCheck.Email)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "пользователь с таким адресом электронной почты не найден",
		})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userToCheck.PasswordHash)) != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "введён неверный пароль",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка создания токена",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", tokenString, 3600*24, "", "", true, true)

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "успешная авторизация",
		Data: gin.H{
			"token": tokenString,
		},
	})
}

// SignOut godoc
// @Summary Выход из системы
// @Description Разлогинивает пользователя и удаляет куки авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse "успешный логаут"
// @Failure 404 {object} dto.ErrorResponse "пользователь не найден"
// @Router /signout [post]
// @Example response 200
//
//			{
//			"success": true,
//		    "message": "успешный логаут"
//	     }
//
// @Example response 404
//
//	{
//	    "success": false,
//	    "error": "пользователь не найден"
//	}
//
// @Security BearerAuth
func SignOut(c *gin.Context) {
	_, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "пользователь не найден",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", "", -1, "", "", true, true)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "успешный логаут",
	})
}

// Me godoc
// @Summary Получить информацию о текущем пользователе
// @Description Возвращает данные авторизованного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse "Данные пользователя"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Router /me [get]
// @Example response 200
//
//	{
//	    "success": true,
//	    "data": {
//	        "email": "user@example.com",
//	        "name": "Иван Иванов",
//	    }
//	}
func Me(c *gin.Context) {
	user := utils.GetUser(c)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data: gin.H{
			"Email": user.Email,
			"Name":  user.Name,
		},
	})
}
