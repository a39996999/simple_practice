package mail

import (
	"chatroom/model"
	"chatroom/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type EmailRequest struct {
	Email string `json:"Email"`
}

func SendMailToken(c *gin.Context) {
	email := EmailRequest{}
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if email.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid email",
		})
		return
	}
	mailToken, err := utils.GenerateToken()
	host_ip := os.Getenv("local_host") + ":" + os.Getenv("server_port")
	verifyUrl := "http://" + host_ip + "/v1/user/verifymail/" + mailToken
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	isVerifyEmail, err := model.RecordSendMail(email.Email, mailToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	} else if isVerifyEmail != false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email was verify",
		})
		return
	}
	message := []byte("To: " + email.Email + "\r\n" +
		"Subject: Verify mail\r\n" + "\r\n" +
		verifyUrl + "\r\n")
	err = utils.SendMail(email.Email, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "send mail sucessfully",
		})
	}
}

func VerifyMailCode(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "token invalid",
		})
		return
	}
	err := model.VerifyMail(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "verify mail sucessfully",
		})
	}
}
