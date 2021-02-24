package middleware

import (
	"log"
	"net/http"
	"photoshare/models"
	"photoshare/service"
	"photoshare/utility"

	"github.com/gin-gonic/gin"
)

func UserValidate(c *gin.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(http.StatusOK, models.Nologin("token有误，请重新登录"))
		c.Abort()
		return
	}
	log.Println(cookie.Value)
	var id int32
	id, err = utility.GetIdByToken(cookie.Value)
	if err != nil {
		c.JSON(http.StatusOK, models.Nologin(err.Error()))
		c.Abort()
		return
	}
	var user models.User
	user, err = service.GetUserInfoById(id)
	if err != nil {
		c.JSON(http.StatusOK, models.Nologin(err.Error()))
		c.Abort()
		return
	}
	c.Set("user", user)
	c.Next()
}
