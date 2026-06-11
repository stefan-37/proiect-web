package handler

import(
	"backend/service"
	"github.com/gin-gonic/gin"
)

func UserSignUp(c *gin.Context) {
	service.UserSignUp(c, database)
}

func UserDelete(c *gin.Context) {
	service.UserDelete(c, database)
}

func UserUpdate(c *gin.Context) {
	service.UserUpdate(c, database)
}

func UserGet(c *gin.Context) {
	service.UserGet(c, database)
}

func UserSubscribe(c *gin.Context) {
	service.UserSubscribe(c, database)
}

func GetUserSubscriptions(c *gin.Context) {
	service.GetUserSubscriptions(c, database)
}

func GetClasses(c *gin.Context) {
	service.GetClasses(c, database)
}

func BookClass(c *gin.Context) {
	service.BookClass(c, database)
}

func UserLogout(c *gin.Context) {
	service.UserLogout(c)
}

func UserResetPassword(c *gin.Context) {
	service.UserResetPassword(c, database)
}

func GetReservations(c *gin.Context) {
	service.GetReservations(c, database)
}

func DeleteReservation(c *gin.Context) {
	service.DeleteReservation(c, database)
}

func ForgotPassword(c *gin.Context) {
	service.ForgotPassword(c, database)
}

func GetUserPayments(c *gin.Context) {
	service.GetUserPayments(c, database)
}

func GetUserReceipt(c *gin.Context) {
	service.GetUserReceipt(c, database)
}

func GetBookingReceipt(c *gin.Context) {
	service.GetBookingReceipt(c, database)
}