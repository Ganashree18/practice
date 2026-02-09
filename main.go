package main 
import (
	"log"
	"practice/db"
	"practice/models"
	"practice/routes"
)


func main(){
	db.ConnectPostgres()

	err:= db.DB.AutoMigrate(&models.Trainer{})
	err= db.DB.AutoMigrate(&models.Student{})
	err= db.DB.AutoMigrate(&models.Enquiry{})
	if err != nil{
		log.Fatal("migration unsuccessful", err)
	}
	r:= routes.SetRouter()
	r.Run(":8080")

}









		
	
	

