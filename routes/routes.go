package routes

import (
	"github.com/gin-gonic/gin"
	"practice/controllers"
)

func SetRouter() *gin.Engine{
	r:= gin.Default()
	r.POST("/create-tra/",controllers.CreateTrainer)
	r.GET("/get-tra/",controllers.GetTrainers)
	r.POST("/create-stu/", controllers.CreateStudent)
	r.GET("/get-stu/", controllers.GetStudents)
	r.POST("/create-enq/", controllers.CreateEnquiry)
	r.POST("/add-enq/", controllers.AddEnquiry)
	r.POST("/upload/", controllers.UploadFile)
	r.POST("/change/", controllers.ChangePaswd)

	
	return r


}