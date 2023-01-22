package controllers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/app"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/helpers/auth"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/helpers/hash"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/helpers/messages"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/helpers/util"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// POST: /users/register
func Register(c *gin.Context)  {

	// connect db
	db := c.MustGet("db").(*gorm.DB)

	// get input
	var input app.RegisterInput

	err := util.ReadRequest(&input, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Email = input.Email
	u.Password = input.Password

	_, err = u.Init()
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": "Error",
			"message": err.Error(),
			"data": nil,
		})
		return
	}

	// validate input
	err = u.ValidateInput("userinput") 
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// hash password
	hashedPassword, err := hash.GenerateHash(u.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	u.Password = string(hashedPassword)

	// create new user
	err = db.Debug().Create(&u).Error
	if err != nil {
		formattedError := messages.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	//response
	data := app.RegisterInput{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User registered succesfully!",
		"data":    data,
	})
}

// POST: /users/login
func Login(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var input app.LoginInput
	if err := util.ReadRequest(&input, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	u := models.User{}
	u.Email = input.Email
	u.Password = input.Password

	_,err := u.InitLogin()
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	err = u.ValidateInput("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	var user_data app.LoginData
	err = db.Debug().Table("users").Select("*").Joins("LEFT JOIN photos ON photos.user_id = users.id").
		Where("users.email = ?", u.Email).Find(&user_data).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "Error",
				"message": u.Email + " has not been registered!",
			})
		return
	}

	err = hash.VerifyPassword(user_data.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		formattedError := messages.ErrorMessage(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	token, err := auth.GenerateJWT(user_data.Email, user_data.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	data := app.LoginData{
		ID: user_data.ID, 
		Username: user_data.Username, 
		Email: user_data.Email, 
		Token: token,
	}

	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"status":  "Success",
		"message": "Login successfully!",
		"data":    data,
	})
}

// PUT: /users/:userId (Update User)
func UpdateUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	// get user data
	var user models.User
	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User not found",
			"data":    nil,
		})
		return
	}

	var input app.RegisterInput
	if err := util.ReadRequest(&input, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{			
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	new_data := models.User{}
	new_data.Username = input.Username
	new_data.Email = input.Email
	new_data.Password = input.Password

	err = new_data.ValidateInput("userinput") 
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	hashedPassword, err := hash.GenerateHash(new_data.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	new_data.Password = string(hashedPassword)

	err = db.Debug().Model(&user).Updates(&new_data).Error
	if err != nil {
		formattedError := messages.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	data := app.RegisterInput{
		ID:        new_data.ID,
		Username:  new_data.Username,
		Email:     new_data.Email,
		CreatedAt: new_data.CreatedAt,
		UpdatedAt: new_data.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User updated succesfully!",
		"data":    data,
	})
}

// DELETE: /users/:userId (Delete User)
func DeleteUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User not found",
			"data":    nil,
		})
		return
	}

	err = db.Debug().Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User deleted succesfully",
		"data":    nil,
	})
}