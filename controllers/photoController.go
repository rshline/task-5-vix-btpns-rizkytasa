package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/helpers/auth"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/helpers/messages"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/helpers/util"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/models"
	"gorm.io/gorm"
)

// GET: /photos
func GetPhoto(c *gin.Context) {
	//Create list photo
	photos := []models.Photo{}

	//Set database
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Debug().Model(&models.Photo{}).Limit(100).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	if len(photos) > 0 {
		for i := range photos {
			user := models.User{}
			err := db.Model(&models.User{}).Where("id = ?", photos[i].UserID).Take(&user).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": err.Error(),
					"data":    nil,
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"message": "Data retrieved successfully",
		"data": photos,
	})
}

// POST: /photos
func AddPhoto(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	//Get token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"status":  "Error",
			"message": "token not found",
			"data":    nil,
		})
		return
	}

	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	var user_has_login models.User
	err = db.Debug().Where("email = ?", email).First(&user_has_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User not found",
			"data":    nil,
		})
		return
	}

	input := models.Photo{}
	if err := util.ReadRequest(&input, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}	

	input.Init()
	input.UserID = user_has_login.ID
	err = input.Validate("upload") //Validate photo
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	var old_photo models.Photo
	err = db.Debug().Model(&models.Photo{}).Where("user_id = ?", user_has_login.ID).Find(&old_photo).Error
	if err != nil {
		if err.Error() == "Data not found" {
			err = db.Debug().Create(&input).Error //Create photo to database
			if err != nil {
				formattedError := messages.ErrorMessage(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": formattedError.Error(),
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  "Success",
				"message": "Photo uploaded successfully!",
				"data":    input,
			})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	input.ID = old_photo.ID
	err = db.Debug().Model(&old_photo).Updates(&input).Error
	if err != nil {
		formattedError := messages.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo changed successfully!",
		"data":    input,
	}) //Return response

}


// PUT: /photos
func UpdatePhoto(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"status":  "Error",
			"message": "token not found",
			"data":    nil,
		})
		return
	}

	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	var user_has_login models.User
	err = db.Debug().Where("email = ?", email).First(&user_has_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User not found",
			"data":    nil,
		})
		return
	}

	input := models.Photo{}
	if err := util.ReadRequest(&input, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}	

	err = input.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	if user_has_login.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User unauthorized to edit.",
			"data":    nil,
		})
		return
	}

	err = db.Model(&photo).Updates(&input).Error
	if err != nil {
		formattedError := messages.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo updated successfully!",
		"data":    photo,
	})
}


// DELETE: /photos
func DeletePhoto(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"status":  "Error",
			"message": "token not found",
			"data":    nil,
		})
		return
	}

	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	var user_has_login models.User
	if err := db.Debug().Where("email = ?", email).First(&user_has_login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User not found",
			"data":    nil})
		return
	}

	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	if user_has_login.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User not authorized!",
			"data":    nil,
		})
		return
	}

	err = db.Debug().Delete(&photo).Error
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
		"message": "Photo deleted successfully!",
		"data":    nil})
}