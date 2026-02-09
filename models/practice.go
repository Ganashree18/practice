package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Trainer struct{
    gorm.Model
	Name string `json:"name"`
	Email datatypes.JSON `json:"email"`
	Phone  datatypes.JSON  `json:"phone"`
	Branch  datatypes.JSON  `json:"branch"`
	Password  string `json:"password"`
	Location  datatypes.JSON  `json:"location"`
}
  

type Student struct{
	gorm.Model
	Name string `json:"name"`
	Email datatypes.JSON `json:"email"`
	Phone  datatypes.JSON  `json:"phone"`
	Password  string `json:"password"`
	Branch string `json:"branch"`
	Address  datatypes.JSON  `json:"address"`
	Batches  datatypes.JSON  `json:"batches"`
	Images string `json:"images"`
	CreatedUTC time.Time  `json:"createdutc"`
	
}

type Enquiry struct{
	gorm.Model
	Name string `json:"name"`
	Email datatypes.JSON `json:"email"`
	Contact  datatypes.JSON  `json:"contact"`
	College  string  `json:"college"`
	Yop  string `json:"yop"`
	Degree string `json:"degree"`
	EnquiredCourse  datatypes.JSON  `json:"enquiredcourse"`
	Referral datatypes.JSON `json:"referral"`

}

	
	


// func (s *Student) BeforeCreate(tx *gorm.DB) (err error) {
// 	now := time.Now().UTC()
// 	s.CreatedUTC = now
// 	s.UpdatedUTC = now
// 	return
// }

// func (s *Student) BeforeUpdate(tx *gorm.DB) (err error) {
// 	s.UpdatedUTC = time.Now().UTC()
// 	return
// }

// func (s *Student) BeforeDelete(tx *gorm.DB) (err error) {
// 	now := time.Now().UTC()
// 	tx.Model(s).Update("DeletedUTC", now)
// 	s.DeletedUTC = &now
// 	return
// }



