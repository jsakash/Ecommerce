package controllers

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jsakash/ecommers/pkg/database"
	"github.com/jsakash/ecommers/pkg/models"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func OtpLog(c *gin.Context) {

	Mob := c.Query("number")
	var (
		accountSid string
		authToken  string
		fromPhone  string
		client     *twilio.RestClient
	)

	result := ChekNumber(Mob)

	if !result {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Mobile number doesnt exist! Please SignUp",
		})
		return
	}

	// Get Twillio credentials from .env file
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")

	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	Mobile := "+91" + Mob
	// Creatin 4 digit OTP
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999-1000) + 1000
	otp := strconv.Itoa(value)
	otpDb := models.Otp{Mobile: Mob, Otp: otp}
	database.DB.Create(&otpDb)

	params := openapi.CreateMessageParams{}
	params.SetTo(Mobile)
	params.SetFrom(fromPhone)
	params.SetBody("Your OTP - " + otp)

	_, err := client.Api.CreateMessage(&params)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "error sending OTP",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "OTP Sent Succesfully",
	})

}

// Checking if the number belongs to any user
func ChekNumber(str string) bool {

	mNumber := str
	var checkOtp models.Users
	database.DB.Raw("SELECT phone_number FROM users WHERE phone_number=?", mNumber).Scan(&checkOtp)
	return checkOtp.Phone_Number == mNumber

}

func ValidateOtp(c *gin.Context) {

	sotp := c.Query("otp")
	var userNum string

	database.DB.Raw("SELECT mobile FROM otps WHERE otp=?", sotp).Scan(&userNum)

	var user models.Users
	//database.DB.Where("users.phone_number = ?", userNum).Find(&user)
	database.DB.First(&user, "phone_number = ?", userNum)

	// Look up request user
	var otp models.Otp
	database.DB.First(&otp, "otp = ?", sotp)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid OTP",
		})

		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// SIgn and get the complete encoded token as a string using the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}
	// Sent it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuthorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "ok",
		"data":    tokenString,
	})

	database.DB.Raw("DELETE FROM otps WHERE mobile=?", userNum).Scan(&otp)

}
