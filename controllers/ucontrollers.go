package controllers

import (
	"net/http"
	"practice/db"
	"practice/models"
	"time"

	"github.com/disintegration/imaging"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func istTOutc(t time.Time) time.Time {
	ist, _ := time.LoadLocation("Asia/Kolkata")
	return t.In(ist).UTC()
}
func CreateTrainer(c *gin.Context) {
	var trainers models.Trainer
	if err := c.ShouldBindJSON(&trainers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hPass, err := bcrypt.GenerateFromPassword([]byte(trainers.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}
	trainers.Password = string(hPass)
	db.DB.Create(&trainers)
	c.JSON(http.StatusOK, gin.H{"message": " trainer added successsfully"})
}

func CreateStudent(c *gin.Context) {
	var students models.Student
	if err := c.ShouldBindJSON(&students); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hPass, err := bcrypt.GenerateFromPassword([]byte(students.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}
	students.Password = string(hPass)
	db.DB.Create(&students)
	tnow := istTOutc(time.Now())
	// students.CreatedUTC=istTOutc(time.Now())
	students.CreatedUTC = tnow.AddDate(0, 0, 10)
	c.JSON(http.StatusOK, gin.H{"message": " student added successsfully"})
}

// func GetStudents(c *gin.Context){
// 	var students []models.Student
// 	if err := db.DB.Find(&students).Error; err!=nil{
// 		c.JSON(http.StatusNotFound, gin.H{"error":"Student not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, students)

// 	var res []struct{
// 		Id string `json:"id"`
// 		Name string `json:"name"`
// 		Email string `json:"email"`
// 		Branch string `json:"branch"`
// 		Batches struct{
// 			Sub int  `json:"sub"`
// 			Batchno string   `json:"batchno"`
//         Trainer struct{
// 			Name string `json:"name"`
// 			Branch  struct{
// 				Id string `json:"id"`
// 				Name string `json:"name"`
// 			} `json:"branch"`
// 		} `json:"trainer"`

// 	}  `json:"batches"`
//     // if err := c.ShouldBindJSON(&res); err != nil {
//     //     c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//     //     return
//     // }

//     db.DB.Debug("select t.branch as tbranch, s.batches as sb from students as s left join trainers as t on  t.name in (select b->>'trainer' from jsonb_array_elements(s.batches) as b)")
// }

// type StudentBatchResponse struct {

//     Name      string         `json:"name"`
//     Branch    datatypes.JSON `json:"branch"`
//     Batches   datatypes.JSON `json:"batches"`
//     TrainerBranch datatypes.JSON `json:"trainer_branch"`
// }

// func GetStudents(c *gin.Context) { ----all details will be printed nd only trainer name will be thr

//     type StudentBatchResponse struct {
//         Id      uint           `json:"id"`
//         Name    string         `json:"name"`
//         Branch  string         `json:"branch"`
//         Batches datatypes.JSON `json:"batches"`
//     }

//     var result []StudentBatchResponse

//     query :=
//         `SELECT s.id,s.name,s.branch,(
// 		SELECT json_agg(jsonb_set(b,'{trainer}',to_jsonb(t)))
//         FROM jsonb_array_elements(s.batches) AS b
//         LEFT JOIN trainers t
//         ON trim(lower(t.name)) = trim(lower(b->>'trainer'))) AS batches FROM students s`

//     if err := db.DB.Raw(query).Scan(&result).Error; err != nil {
//         fmt.Println("DB ERROR:", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//	    c.JSON(http.StatusOK, result)
//	}
type StudentBatchResponse struct {
	Id      uint           `json:"id"`
	Name    string         `json:"name"`
	Branch  string         `json:"branch"`
	Batches datatypes.JSON `json:"batches"`
}

// func GetStudents(c *gin.Context) { ---------all details of trainer will be printed
//     type StudentBatchResponse struct {
//         Id      uint           `json:"id"`
//         Name    string         `json:"name"`
//         Branch  string         `json:"branch"`
//         Batches datatypes.JSON `json:"batches"`
//     }

//     var students []StudentBatchResponse
//     if err := db.DB.Raw(`SELECT id, name, branch, batches FROM students`).Scan(&students).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     var trainers []struct {
//         Id    uint           `json:"id"`
//         Name  string         `json:"name"`
//         Email datatypes.JSON `json:"email"`
//         Phone datatypes.JSON `json:"phone"`
//         Branch datatypes.JSON `json:"branch"`
//     }
//     db.DB.Raw(`SELECT id, name, email, phone, branch FROM trainers`).Scan(&trainers)

//     for i, student := range students {
//         batches := student.Batches
//         for _, t := range trainers {
//             batches = datatypes.JSON(bytes.ReplaceAll(batches, []byte(`"`+t.Name+`"`), []byte(t.ToJSON())))
//         }
//         students[i].Batches = batches
//     }

//     c.JSON(http.StatusOK, students)
// }

func GetStudents(c *gin.Context) {

	var result []StudentBatchResponse

	query :=
		`SELECT s.id,s.name,s.branch,(SELECT json_agg(jsonb_set(b,'{trainer}',to_jsonb(
        json_build_object('id', t.id,'name', t.name,'email', t.email,'phone', t.phone,'branch', t.branch,
        'location', t.location))))
        FROM jsonb_array_elements(s.batches) AS b
        LEFT JOIN trainers t
        ON trim(lower(t.name)) = trim(lower(b->>'trainer'))) AS batches
        FROM students s`

	if err := db.DB.Raw(query).Scan(&result).Error; err != nil {
		fmt.Println("DB ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetTrainers(c *gin.Context) {
	var trainers []models.Trainer
	if err := db.DB.Find(&trainers).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trainer not found"})
		return
	}
	c.JSON(http.StatusNotFound, trainers)
}

// func UploadFile(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "file not found"})
// 		return
// 	}

// 	openedfile, err := file.Open()
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "cannot open file"})
// 		return
// 	}
// 	defer openedfile.Close()

// 	img, err := imaging.Decode(openedfile)
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "failed to decode image"})
// 		return
// 	}
// 	thumb := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)

// 	origBuf, thumbBuf := new(bytes.Buffer), new(bytes.Buffer)
// 	imaging.Encode(origBuf, img, imaging.JPEG)
// 	imaging.Encode(thumbBuf, thumb, imaging.JPEG)

// 	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=lion dbname=practiceDB sslmode=disable")
// 	if err != nil {
// 		fmt.Println("SQL Open error:", err)
// 		c.JSON(500, gin.H{"error": "database connection failed"})
// 		return
// 	}

// 	if err := db.Ping(); err != nil {
// 		fmt.Println("DB Ping error:", err)
// 		c.JSON(500, gin.H{"error": "cannot reach database"})
// 		return
// 	}

// 	var id int
// 	db.QueryRow("INSERT INTO images (filename, original, thumbnail) VALUES ($1,$2,$3) RETURNING id",
// 		file.Filename, origBuf.Bytes(), thumbBuf.Bytes()).Scan(&id)

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "uploaded successfully", "id": id})

// }

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
		return
	}

	path := "C:\\Users\\pc\\Desktop\\holdmyfiles\\"
	err = c.SaveUploadedFile(file, path+file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't save the file"})
		return
	}
	originalPath := path + file.Filename
	img, err := imaging.Open(originalPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image"})
		return
	}

	thumb := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)
	thumbPath := path + "thumb_" + file.Filename

	err = imaging.Save(thumb, thumbPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't save thumbnail"})
		return
	}

	db.DB.Raw("update students set images = ? where id = 5", path+file.Filename).Scan(&models.Student{})

	c.JSON(http.StatusOK, gin.H{
		"message":   "Uploaded successfully",
		"FileName":  path + file.Filename,
		"thumbnail": "thumb_" + file.Filename,
	})

}

func ChangePaswd(c *gin.Context) {
	var ip struct {
		Id       int    `json:"id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&ip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	hPass, err := bcrypt.GenerateFromPassword([]byte(ip.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}
	ip.Password = string(hPass)

	db.DB.Debug().Raw("update students set password = ? where id = ?", ip.Password, ip.Id).Scan(&models.Student{})
	c.JSON(http.StatusOK, gin.H{"message": "Enquiry registered successfully"})

}

func CreateEnquiry(c *gin.Context) {
	var enquiries models.Enquiry
	if err := c.ShouldBindJSON(&enquiries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&enquiries)
	c.JSON(http.StatusOK, gin.H{"message": "Enquiry registered successfully"})
}

// 	var ip struct{
// 		Id int   `json:"id"`
// 		Name string    `json:"name"`
// 		Yop string   `json:"yop"`
// 		Degree string   `json:"degree"`
// 		EnquiredCourse struct{
// 			Cid int  `json:"cid"`
// 			Cname string   `json:"cname"`
// 		}  `json:"enquiredcourse"`
// 		Contact struct{
// 			Code int  `json:"code"`
// 			Phone string   `json:"phone"`
// 		}  `json:"contact"`

//		Email struct{
//			Email string    `json:"email"`
//		}  `json:"email"`
//		Referral struct{
//			Email string    `json:"email"`
//		}  `json:"referral"`
//		}
//	}
//
// // type Email struct{
// // 	Email string `json:"email"`
// // }
func AddEnquiry(c *gin.Context) {
	var ip models.Enquiry
	if err := c.ShouldBindJSON(&ip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enquiry := models.Enquiry{
		Name:           ip.Name,
		Email:          ip.Email,
		Contact:        ip.Contact,
		College:        ip.College,
		Yop:            ip.Yop,
		Degree:         ip.Degree,
		EnquiredCourse: ip.EnquiredCourse,
		Referral:       ip.Referral,
	}

	type course []struct {
		Cid   uint   `json:"cid"`
		CName string `json:"cname"`
	}
	type email struct {
		Email uint `json:"email"`
	}

	var ecourse course
	json.Unmarshal(enquiry.EnquiredCourse, &ecourse)
	for _, j := range ecourse {
		fmt.Println(j)
	}
	fmt.Println(ecourse)
	marshalcourse, _ := json.Marshal(ecourse)

	db.DB.Debug().Raw("update enquiries set , enquiredcourse = ?, email = ?, phone = ? where id = ?", marshalcourse, ip.Id).Scan(&models.Enquiry{})

	var eemail []email
	json.Unmarshal(enquiry.Email, &eemail)
	for _, j := range eemail {
		fmt.Println(j)

	}
	marshalemail, _ := json.Marshal(eemail)

	db.DB.Debug().Raw("update enquiries set  , enquiredcourse = ?, email = ?, phone = ? where id = ?", marshalcourse, marshalemail, ip.Id).Scan(&models.Enquiry{})

	if err := db.DB.Create(&enquiry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enquiry"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Enquiry added successfully", "data": enquiry})
}

// func AddEnquiryy(c *gin.Context) {
// 	user, _ := c.Get("user")
// 	// userid := user.(models.GeneralUser).ID
// 	type Branch struct {
// 		Id   uint   `json:"id"`
// 		Name string `json:"name"`;
// 	}
// 	type Role struct {
// 		Id   int    `json:"id"`
// 		Name string `json:"name"`
// 	}

// 	type InputAuthor struct {
// 		Id       uint           `json:"id"`
// 		Name     string         `json:"name"`
// 		Roles    datatypes.JSON `json:"roles"`
// 		Branches datatypes.JSON `json:"branches"`
// 	}
// 	var auth InputAuthor
// 	type EnquiryAdd struct {
// 		ID        uint           `json:"id"`
// 		Name      string         `gorm:"size:255;not null" json:"name"`
// 		Email     datatypes.JSON `json:"email"`
// 		Contact   datatypes.JSON `json:"contact"`
// 		Enquiry   datatypes.JSON `json:"enquiry"`
// 		Course    datatypes.JSON `json:"course"`
// 		Records   datatypes.JSON `json:"records"`
// 		Scourse   datatypes.JSON `json:"scourse"`
// 		Education datatypes.JSON `json:"education"`
// 		Refers    datatypes.JSON `json:"refers"`
// 		Author    InputAuthor    `json:"author"`
// 		CreatedAt time.Time      `json:"created"`
// 		UpdatedAt time.Time      `json:"updated"`
// 		Edited    bool           `json:"edited"`
// 		Rid       int            `json:"rid"`
// 		Regular   bool           `json:"regular"`
// 		Special   bool           `json:"special"`
// 	}
// 	ue := user.(models.Enquiry).Email
// 	ec := user.(models.Enquiry).EnquiredCourse
// // 	auth.Id = user.(models.GeneralUser).ID
// // 	auth.Name = user.(models.GeneralUser).Name
// // 	auth.Branches = user.(models.GeneralUser).Branches
// // 	auth.Roles = user.(models.GeneralUser).Roles

// 	var usercourse []struct {
// 		Cid int  `json:"cid"`
// 		Cname string   `json:"cname"`
// //
// 	}
// 	var useremail []struct {

// 		Email string   `json:"email"`
// //
// 	}
// // 	var roles []Role
// 	json.Unmarshal(ue, &useremail)
// 	json.Unmarshal(ec, &usercourse)
// // 	counselor := true
// // 	cf := false
// // 	fc := false
// // 	collector := true
// // 	for _, r1 := range roles {
// // 		if r1.Name == "counselor" {
// // 			cf = true
// // 		}
// // 		if r1.Name == "feecollector" {
// // 			fc = true
// // 		}
// // 	}
// // 	if !cf {
// // 		counselor = false
// // 	}
// // 	if !fc {
// // 		collector = false
// // 	}
// 	enquiry := models.Enquiry{}
// 	var input EnquiryAdd
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fmt.Print((enquiry))

// // 	var cons []struct {
// // 		Code     string `json:"code"`
// // 		Number   string `json:"number"`
// // 		ReadOnly bool   `json:"readonly"`
// // 	}

// // 	var mails []struct {
// // 		Email string `json:"email"`
// // 	}

// // 	var dataNotTaken struct {
// // 		DataNotTaken bool `json:"data_not_taken"`
// // 		Experienced  bool `json:"experienced"`
// // 	}

// 	var referrals []struct {
// 		Ref_code   string `json:"ref_code"`
// 		Ref_phone string `json:"ref_phone"`
// 		Ref_name   string `json:"ref_name"`

// 	}
// 	var refers []struct {
// 		Ref_code   string `json:"ref_code"`
// 		Ref_phone string `json:"ref_phone"`
// 		Ref_name   string `json:"ref_name"`

// 	}
// // 	// var cons []Contacts
// // 	// var mails []Emails

// // 	json.Unmarshal(input.Enquiry, &dataNotTaken)
// // 	json.Unmarshal(input.Contact, &cons)
// // 	json.Unmarshal(input.Email, &mails)
// 	json.Unmarshal(input.Refers, &referrals)
// // 	currentTime1 := time.Now()
// // 	var cqs string
// // 	var crs string
// // 	for _, element := range cons {
// // 		// UpdateNoEnquiryContact(element.Number, user.(models.GeneralUser).ID)
// // 		if len(cqs) > 0 {
// // 			cqs = cqs + " OR "
// // 		}
// // 		if len(crs) > 0 {
// // 			crs = crs + " OR "
// // 		}
// // 		crs = crs + "name = '" + element.Number + "'"
// // 		cqs = cqs + "contact @@ '$[*].number == \"" + element.Number + "\"'"
// // 	}
// // 	for _, element := range mails {
// // 		// UpdateNoEnquiryEmail(element.Email, user.(models.GeneralUser).ID)
// // 		if len(cqs) > 0 {
// // 			cqs = cqs + " OR "
// // 		}
// // 		cqs = cqs + "email @@ '$[*].email == \"" + element.Email + "\"'"
// // 	}
// // 	database.DB.Raw("select * from enquiries where " + cqs).Scan(&enquiry)
// // 	type Course struct {
// // 		Name string `json:"name"`
// // 	}
// // 	type Mode_class struct {
// // 		Id   int    `json:"id"`
// // 		Name string `json:"name"`
// // 	}

// 	type Courses struct {
// 		Agreed_date   string         `json:"agreed_date"`
// 		Author        InputAuthor    `json:"author"`
// 		Branch        []Branch       `json:"branch"`
// 		Course        []Course       `json:"course"`
// 		Mode_class    Mode_class     `json:"mode_class"`
// 		Date          string         `json:"date"`
// 		Modeofenquiry datatypes.JSON `json:"modeofenquiry"`
// 		Trainer       datatypes.JSON `json:"trainer"`
// 		Walkindate    string         `json:"walkindate"`
// 		Conversion    int            `json:"conversion"`
// 		Tats          bool           `json:"tats"`
// 	}

// // 	type Scourses struct {
// // 		Name struct {
// // 			Id   int    `json:"id"`
// // 			Name string `json:"name"`
// // 		} `json:"name"`
// // 		Author      InputAuthor    `json:"author"`
// // 		Oauthor     InputAuthor    `json:"oauthor"`
// // 		Branch      []Branch       `json:"branch"`
// // 		Mode_class  Mode_class     `json:"mode_class"`
// // 		Followup    datatypes.JSON `json:"followup"`
// // 		Date        string         `json:"date"`
// // 		Othercourse bool           `json:"othercourse"`
// // 		OtherCourse string         `json:"other_course"`
// // 		Archived    bool           `json:"archived"`
// // 		Comment     string         `json:"comment"`
// // 	}

// // 	type AllCourses struct {
// // 		Agreed_date   string         `json:"agreed_date"`
// // 		Author        InputAuthor    `json:"author"`
// // 		Oauthor       InputAuthor    `json:"oauthor"`
// // 		Branch        Branch         `json:"branch"`
// // 		Course        Course         `json:"course"`
// // 		Mode_class    Mode_class     `json:"mode_class"`
// // 		Details       datatypes.JSON `json:"details"`
// // 		Followup      datatypes.JSON `json:"followup"`
// // 		Date          string         `json:"date"`
// // 		Modeofenquiry datatypes.JSON `json:"modeofenquiry"`
// // 		Trainer       datatypes.JSON `json:"trainer"`
// // 		Walkindate    string         `json:"walkindate"`
// // 		Conversion    int            `json:"conversion"`
// // 		Tats          bool           `json:"tats"`
// // 	}

// // 	var coursesinput []Courses
// // 	var reEnqCourses []AllCourses
// // 	var reSpcEnqCourses []Scourses
// // 	var courses []AllCourses
// // 	var ocourses []AllCourses
// // 	var p1courses []AllCourses
// // 	// var p2courses []AllCourses
// // 	// var poCourses []AllCourses
// // 	var scoursesinput []Scourses
// // 	var scourses []Scourses
// // 	var education Education
// // 	var dnt struct {
// // 		DataNotTaken bool `json:"data_not_taken"`
// // 		Experienced  bool `json:"experienced"`
// // 	}

// // 	var emails []struct {
// // 		Email string `json:"email"`
// // 	}

// // 	var contacts []struct {
// // 		Code     string `json:"code"`
// // 		Number   string `json:"number"`
// // 		ReadOnly bool   `json:"readonly"`
// // 	}
// // 	type EducationAuthor struct {
// // 		Id        int            `json:"id"`
// // 		Name      string         `json:"name"`
// // 		Roles     datatypes.JSON `json:"roles"`
// // 		Branches  datatypes.JSON `json:"branches"`
// // 		Education Education      `json:"education"`
// // 	}
// // 	var authorEducation []EducationAuthor

// // 	type EnquiryBasic struct {
// // 		Comment         string         `json:"comment"`
// // 		Time_slot       datatypes.JSON `json:"time_slot"`
// // 		Experienced     bool           `json:"experienced"`
// // 		Class_timing    datatypes.JSON `json:"class_timing"`
// // 		Flexi_timing    bool           `json:"flexi_timing"`
// // 		Enquiree_name   string         `json:"enquiree_name"`
// // 		Data_not_taken  bool           `json:"data_not_taken"`
// // 		Enq_for_someone bool           `json:"enq_for_someone"`
// // 		Mode_of_enquiry datatypes.JSON `json:"mode_of_enquiry"`
// // 	}
// // 	var inputBasic EnquiryBasic
// // 	var existingBasic EnquiryBasic
// // 	var message string

// 	// json.Unmarshal(enquiry.Courses, &courses)
// // 	json.Unmarshal(input.Course, &coursesinput)
// // 	json.Unmarshal(enquiry.Scourses, &scourses)
// // 	json.Unmarshal(input.Scourse, &scoursesinput)
// // 	json.Unmarshal(enquiry.Ocourses, &ocourses)
// // 	json.Unmarshal(enquiry.Pcourses, &p1courses)
// // 	// json.Unmarshal(enquiry.Oncourses, &p2courses)
// // 	json.Unmarshal(enquiry.Email, &emails)
// // 	json.Unmarshal(enquiry.Contact, &contacts)
// // 	json.Unmarshal(enquiry.Refers, &refers)
// // 	json.Unmarshal(enquiry.Enquiry, &dnt)
// // 	json.Unmarshal(input.Education, &education)
// // 	json.Unmarshal(education.Author, &authorEducation)
// // 	json.Unmarshal(enquiry.Enquiry, &existingBasic)
// // 	json.Unmarshal(input.Enquiry, &inputBasic)
// // 	json.Unmarshal(enquiry.Rcourses, &reEnqCourses)
// // 	json.Unmarshal(enquiry.Rscourses, &reSpcEnqCourses)
// // 	if enquiry.ID > 0 {
// // 		message = "Repeat Enquiry"
// // 		type Status struct {
// // 			Converted bool   `json:"converted"`
// // 			Name      string `json:"name"`
// // 			Date      string `json:"date"`
// // 		}
// // 		var status Status
// // 		var reEnuiry = enquiry.ReEnquiry
// // 		json.Unmarshal(enquiry.Status, &status)

// // 		if (len(status.Name) > 0 && (input.Edited && (input.Regular || input.Special))) || (len(status.Name) > 0 && !input.Edited) {
// // 			reEnqCourses = append(reEnqCourses, courses...)
// // 			// reEnqCourses = append(reEnqCourses, ocourses...)
// // 			// reEnqCourses = append(reEnqCourses, p2courses...)
// // 			// reEnqCourses = append(reEnqCourses, p1courses...)
// // 			reSpcEnqCourses = append(reSpcEnqCourses, scourses...)
// // 			scourses = make([]Scourses, 0)
// // 			// p2courses = make([]AllCourses, 0)
// // 			// p1courses = make([]AllCourses, 0)
// // 			courses = make([]AllCourses, 0)
// // 			// ocourses = make([]AllCourses, 0)
// // 		}
// // 		for _, e := range coursesinput {
// // 			for _, e2 := range e.Course {
// // 				for _, e1 := range e.Branch {
// // 					bu := false
// // 					for _, b1 := range userbranch {
// // 						if e1.Id == b1.Id {
// // 							bu = true
// // 						}
// // 					}
// // 					var sc AllCourses
// // 					if bu && counselor {
// // 						sc.Author = auth
// // 					} else {
// // 						sc.Oauthor = auth
// // 					}
// // 					sc.Branch = e1
// // 					sc.Course = e2
// // 					sc.Mode_class = e.Mode_class
// // 					sc.Agreed_date = e.Agreed_date
// // 					sc.Date = currentTime1.Format("2006-01-02 15:04:05")
// // 					sc.Modeofenquiry = e.Modeofenquiry
// // 					sc.Walkindate = e.Walkindate
// // 					sc.Trainer = e.Trainer
// // 					// if strings.ToLower(e.Mode_class.Name) == "online" {
// // 					// 	p2found := false
// // 					// 	for _, p2 := range p2courses {
// // 					// 		if p2.Course.Name == sc.Course.Name {
// // 					// 			p2found = true
// // 					// 		}
// // 					// 	}
// // 					// 	if !p2found {
// // 					// 		p2courses = append(p2courses, sc)
// // 					// 	}
// // 					// 	if len(status.Name) > 0 {
// // 					// 		reEnuiry = true
// // 					// 		message = "Re-Enquiry"
// // 					// 	}

// // 					// }
// // 					cfound := false
// // 					if bu && counselor {
// // 						for _, c := range courses {
// // 							// if c.Course.Name == sc.Course.Name && c.Branch.Id == sc.Branch.Id {
// // 							if c.Course.Name == sc.Course.Name && c.Author.Id == auth.Id {
// // 								cfound = true
// // 							}
// // 						}
// // 						if !cfound {
// // 							courses = append(courses, sc)
// // 							// go Walkintat(sc.Modeofenquiry, userid, sc.Date, sc.Branch.Name, sc.Branch.Id, sc.Walkindate)
// // 						}
// // 						if len(status.Name) > 0 {
// // 							reEnuiry = true
// // 							message = "Re-Enquiry"
// // 						}
// // 					} else {
// // 						ofound := false
// // 						for _, c := range ocourses {
// // 							if c.Course.Name == sc.Course.Name && c.Branch.Id == sc.Branch.Id {
// // 								ofound = true
// // 							}
// // 						}
// // 						for _, c := range courses {
// // 							if c.Course.Name == sc.Course.Name && c.Branch.Id == sc.Branch.Id {
// // 								cfound = true
// // 							}
// // 						}
// // 						if !ofound && !cfound {
// // 							sc.Trainer = nil
// // 							ocourses = append(ocourses, sc)
// // 						}
// // 						if len(status.Name) > 0 {
// // 							reEnuiry = true
// // 							message = "Re-Enquiry"
// // 						}
// // 					}
// // 				}

// // 			}
// // 		}
// // 		for _, ele := range scoursesinput {
// // 			var sc Scourses
// // 			sc.Name = ele.Name
// // 			sc.Oauthor = auth
// // 			sc.Branch = ele.Branch
// // 			sc.Mode_class = ele.Mode_class
// // 			sc.Date = time.Now().Format("2006-01-02 15:04:05")
// // 			sc.OtherCourse = ele.OtherCourse
// // 			sc.Othercourse = ele.Othercourse
// // 			sc.Archived = ele.Archived
// // 			scfound := false
// // 			for _, sc := range scourses {
// // 				if sc.Name == ele.Name && sc.OtherCourse == ele.OtherCourse {
// // 					scfound = true
// // 				}
// // 			}
// // 			if !scfound {
// // 				scourses = append(scourses, sc)
// // 			}
// // 			if len(status.Name) > 0 {
// // 				reEnuiry = true
// // 				message = "Re-Enquiry"
// // 			}
// // 		}
// // 		for _, ele := range cons {
// // 			exists := 0
// // 			database.DB.Model(&models.Enquiry{}).Select("id").Where("contact @@ '$.number == \"" + fmt.Sprint(ele.Number) + "\"'").Find(&exists)
// // 			if exists == 0 {
// // 				contacts = append(contacts, ele)
// // 			}
// // 		}
// // 		for _, ele := range mails {
// // 			exists := 0
// // 			database.DB.Model(&models.Enquiry{}).Select("id").Where("email @@ '$.email == \"" + fmt.Sprint(ele.Email) + "\"'").Find(&exists)
// // 			if exists == 0 {
// // 				emails = append(emails, ele)
// // 			}
// // 		}
// // 		for _, ele := range referrals {
// // 			exists := 0
// // 			database.DB.Model(&models.Referrals{}).Select("id").Where("number like ?", ele.Referee_number).Find(&exists)
// // 			if exists == 0 {
// // 				refers = append(refers, ele)
// // 			}
// // 		}
// // 		addEducation := true
// // 		var authEdu []byte
// // 		if input.Edited {
// // 			type IDS struct {
// // 				ID int `json:"id"`
// // 			}
// // 			var highestDegree IDS
// // 			var name IDS
// // 			var streamObject IDS
// // 			var collegeObject IDS
// // 			var inputEducation Education
// // 			var existingEducation Education
// // 			json.Unmarshal(enquiry.Education, &existingEducation)
// // 			json.Unmarshal(existingEducation.HighestDegree, &highestDegree)
// // 			json.Unmarshal(existingEducation.Name, &name)
// // 			json.Unmarshal(existingEducation.CollegeObject, &collegeObject)
// // 			json.Unmarshal(existingEducation.StreamObject, &streamObject)
// // 			var inputName IDS
// // 			var inputStreamObject IDS
// // 			var inputCollegeObject IDS
// // 			var inputHighestDegree IDS
// // 			json.Unmarshal(input.Education, &inputEducation)
// // 			json.Unmarshal(inputEducation.HighestDegree, &inputHighestDegree)
// // 			json.Unmarshal(inputEducation.Name, &inputName)
// // 			json.Unmarshal(inputEducation.CollegeObject, &inputCollegeObject)
// // 			json.Unmarshal(inputEducation.StreamObject, &inputStreamObject)
// // 			if name.ID != inputName.ID || collegeObject.ID != inputCollegeObject.ID || streamObject.ID != inputStreamObject.ID || existingEducation.Percentage != inputEducation.Percentage || existingEducation.YOP != inputEducation.YOP || existingEducation.Otherdegree != inputEducation.Otherdegree || existingEducation.Othercollege != inputEducation.Othercollege || existingEducation.Otherstream != inputEducation.Otherstream || highestDegree != inputHighestDegree {
// // 				addEducation = false
// // 			}
// // 		}

// // 		json.Unmarshal(enquiry.Education, &education)
// // 		json.Unmarshal(education.Author, &authorEducation)
// // 		var xedu Education
// // 		json.Unmarshal(input.Education, &xedu)
// // 		xedu.Author = nil
// // 		auth := EducationAuthor{int(user.(models.GeneralUser).ID), user.(models.GeneralUser).Name, user.(models.GeneralUser).Roles, user.(models.GeneralUser).Branches, xedu}
// // 		if !addEducation || !input.Edited || reEnuiry {
// // 			authorEducation = append(authorEducation, auth)
// // 			authEdu, _ = json.Marshal(authorEducation)
// // 			education.Author = authEdu
// // 		}
// // 		if !counselor && !collector && !dataNotTaken.DataNotTaken {
// // 			json.Unmarshal(input.Education, &education)
// // 		}
// // 		var reSpEnq any
// // 		var cs any
// // 		var othc any
// // 		var sc any
// // 		var em any
// // 		var cn any
// // 		var re any
// // 		var reEnq any

// // 		update := models.Enquiry{}
// // 		if len(courses) > 0 {
// // 			cs, _ = json.Marshal(&courses)
// // 		} else {
// // 			cs = nil
// // 		}
// // 		if len(ocourses) > 0 {
// // 			othc, _ = json.Marshal(&ocourses)
// // 		} else {
// // 			othc = nil
// // 		}
// // 		// var oc, _ = json.Marshal(&p2courses)
// // 		if len(scourses) > 0 {
// // 			sc, _ = json.Marshal(&scourses)
// // 		} else {
// // 			sc = nil
// // 		}
// // 		if len(emails) > 0 {
// // 			em, _ = json.Marshal(&emails)
// // 		} else {
// // 			em = nil
// // 		}
// // 		if len(contacts) > 0 {
// // 			cn, _ = json.Marshal(&contacts)
// // 		} else {
// // 			cn = nil
// // 		}
// // 		if len(refers) > 0 {
// // 			re, _ = json.Marshal(&refers)
// // 		} else {
// // 			re = nil
// // 		}
// // 		edu, _ := json.Marshal(&education)
// // 		if len(reEnqCourses) > 0 {
// // 			reEnq, _ = json.Marshal(&reEnqCourses)
// // 		} else {
// // 			reEnq = nil
// // 		}

// // 		if len(reSpcEnqCourses) > 0 {
// // 			reSpEnq, _ = json.Marshal(reSpcEnqCourses)
// // 		} else {
// // 			reSpEnq = nil
// // 		}

// // 		// if status.Name != "converted" {
// // 		// 	database.DB.Raw("update enquiries set status = null where id = ?", enquiry.ID).Scan(&update)
// // 		// }
// // 		exp := false
// // 		dataNot := false
// // 		if dataNotTaken.Experienced {
// // 			exp = true
// // 		}
// // 		if dnt.DataNotTaken && !dataNotTaken.DataNotTaken {
// // 			dataNot = false
// // 			json.Unmarshal(input.Education, &education)
// // 			education.Author = authEdu
// // 			edu, _ = json.Marshal(&education)
// // 		}
// // 		if dataNotTaken.DataNotTaken && dnt.Experienced {
// // 			dataNot = true
// // 			json.Unmarshal(input.Education, &education)
// // 			education.Author = authEdu
// // 			edu, _ = json.Marshal(&education)
// // 		}
// // 		if dnt.Experienced && !dataNotTaken.Experienced {
// // 			exp = false
// // 			json.Unmarshal(input.Education, &education)
// // 			education.Author = authEdu
// // 			edu, _ = json.Marshal(&education)
// // 		}
// // 		if dnt.DataNotTaken && dataNotTaken.DataNotTaken {
// // 			dataNot = true
// // 			json.Unmarshal(input.Education, &education)
// // 			education.Author = authEdu
// // 			edu, _ = json.Marshal(&education)
// // 			// json.Unmarshal(input.Enquiry, &existingBasic)
// // 		}
// // 		if input.Edited {
// // 			json.Unmarshal(input.Enquiry, &existingBasic)
// // 			enquiry.Name = input.Name
// // 			message = "Edit Enquiry"
// // 		}
// // 		existingBasic.Data_not_taken = dataNot
// // 		existingBasic.Experienced = exp
// // 		var statusa any
// // 		if reEnuiry && !input.Edited {
// // 			statusa = nil
// // 		} else {
// // 			statusa = enquiry.Status
// // 		}

// // 		extBasic, _ := json.Marshal(existingBasic)
// // 		database.DB.Raw("update enquiries set name = ?, email = ?, contact = ?, courses = ?, ocourses = ?, rcourses = ?, scourses = ?, refers = ?, updated_at = ?, education = ?, enquiry = enquiry || ?, re_enquiry = ?, status = ?, rscourses = ? where id = ?", enquiry.Name, em, cn, cs, othc, reEnq, sc, re, currentTime1, edu, extBasic, reEnuiry, statusa, reSpEnq, enquiry.ID).Scan(&update)
// // 		// go updateNoEnquiryForm(enquiry.ID, userid, "1st func")
// // 		go updateRefs(enquiry.ID)
// // 		go updateRecs(enquiry.ID)
// // 		go addCFNEnquiry(enquiry.ID, userid, true)
// // 		go Savehistory(int(enquiry.ID), message, user)
// // 		go Tagleadenqstu(int(enquiry.ID), enquiry.Name, input.Contact, input.Email, "enquiry")
// // 		go TagEnquiryToWebisteEnquiries(&contacts, int(enquiry.ID), &emails)
// // 		// // CounselorStatistics(int(enquiry.ID))
// // 		c.JSON(http.StatusOK, gin.H{"message": "Added enquiry " + enquiry.Name + "  " + fmt.Sprint(enquiry.ID)})
// // 		return
// // 	}
// // 	incubation := false
// // 	for _, e := range coursesinput {
// // 		for _, e2 := range e.Course {
// // 			for _, e1 := range e.Branch {
// // 				bu := false
// // 				for _, b1 := range userbranch {
// // 					if e1.Id == b1.Id {
// // 						bu = true
// // 					}
// // 				}
// // 				var sc AllCourses
// // 				if bu && counselor {
// // 					sc.Author = auth
// // 				} else {
// // 					sc.Oauthor = auth
// // 				}
// // 				if e1.Id == 178 || e1.Id == 105 {
// // 					incubation = true
// // 				}
// // 				sc.Branch = e1
// // 				sc.Course = e2
// // 				sc.Mode_class = e.Mode_class
// // 				sc.Agreed_date = e.Agreed_date
// // 				sc.Date = currentTime1.Format("2006-01-02 15:04:05")
// // 				sc.Modeofenquiry = e.Modeofenquiry
// // 				sc.Walkindate = e.Walkindate
// // 				sc.Trainer = e.Trainer
// // 				// if strings.ToLower(e.Mode_class.Name) == "online" {
// // 				// 	p2found := false
// // 				// 	for _, p2 := range p2courses {
// // 				// 		// if p2.Course.Name == sc.Course.Name && p2.Branch.Id == sc.Branch.Id {
// // 				// 		if p2.Course.Name == sc.Course.Name && p2.Author.Id == auth.Id {
// // 				// 			p2found = true
// // 				// 		}
// // 				// 	}
// // 				// 	if !p2found {
// // 				// 		p2courses = append(p2courses, sc)
// // 				// 	}

// // 				// }
// // 				if bu && counselor {
// // 					cfound := false
// // 					for _, c := range courses {
// // 						if c.Course.Name == sc.Course.Name {
// // 							cfound = true
// // 						}
// // 					}
// // 					if !cfound {
// // 						courses = append(courses, sc)
// // 						// Walkintat(sc.Modeofenquiry, userid, sc.Date, sc.Branch.Name, sc.Branch.Id, sc.Walkindate)
// // 					}
// // 				} else {
// // 					ofound := false
// // 					for _, c := range ocourses {
// // 						// if c.Course.Name == sc.Course.Name && c.Branch.Id == sc.Branch.Id {
// // 						if c.Course.Name == sc.Course.Name && c.Author.Id == auth.Id {
// // 							ofound = true
// // 						}
// // 					}
// // 					if !ofound {
// // 						sc.Trainer = nil
// // 						ocourses = append(ocourses, sc)
// // 					}
// // 				}
// // 			}

// // 		}
// // 	}

// // 	for _, ele := range scoursesinput {
// // 		var sc Scourses
// // 		sc.Name = ele.Name
// // 		sc.Oauthor = ele.Oauthor
// // 		sc.Branch = ele.Branch
// // 		sc.Mode_class = ele.Mode_class
// // 		sc.Date = time.Now().Format("2006-01-02 15:04:05")
// // 		sc.Oauthor = auth
// // 		sc.OtherCourse = ele.OtherCourse
// // 		sc.Othercourse = ele.Othercourse
// // 		sc.Archived = false
// // 		scourses = append(scourses, sc)
// // 	}

// // 	newEnquiry := EducationAuthor{int(user.(models.GeneralUser).ID), user.(models.GeneralUser).Name, user.(models.GeneralUser).Roles, user.(models.GeneralUser).Branches, education}
// // 	authorEducation = append(authorEducation, newEnquiry)
// // 	var authEdu, _ = json.Marshal(authorEducation)
// // 	education.Author = authEdu

// // 	var cs datatypes.JSON
// // 	var othc datatypes.JSON
// // 	var sc datatypes.JSON
// // 	if len(courses) > 0 {
// // 		cs, _ = json.Marshal(&courses)
// // 	} else {
// // 		cs = nil
// // 	}
// // 	if len(ocourses) > 0 {
// // 		othc, _ = json.Marshal(&ocourses)
// // 	} else {
// // 		othc = nil
// // 	}
// // 	if len(scourses) > 0 {
// // 		sc, _ = json.Marshal(&scourses)
// // 	} else {
// // 		sc = nil
// // 	}
// // 	if len(mails) <= 0 {
// // 		input.Email = nil
// // 	}
// // 	if len(referrals) <= 0 {
// // 		input.Refers = nil
// // 	}
// // 	var incub Incubationstruct
// // 	incmarshal, _ := json.Marshal(incub)
// // 	// var oc, _ = json.Marshal(&p2courses)
// // 	// var sc, _ = json.Marshal(&scourses)
// // 	var edu, _ = json.Marshal(&education)
// // 	r := models.Enquiry{}
// // 	r.Name = input.Name
// // 	r.Email = input.Email
// // 	r.Contact = input.Contact
// // 	r.Enquiry = input.Enquiry
// // 	r.Courses = cs
// // 	r.Ocourses = othc
// // 	// r.Oncourses = oc
// // 	r.Scourses = sc
// // 	r.Education = edu
// // 	if incubation == true {
// // 		r.Incubation = incmarshal
// // 	}
// // 	// r.Records = rs
// // 	r.Refers = input.Refers
// // 	r.Author = user.(models.GeneralUser).ID
// // 	err := database.DB.Create(&r).Error

// // 	type enquiryBasic struct {
// // 		Comment         string         `json:"comment"`
// // 		Time_slot       datatypes.JSON `json:"time_slot"`
// // 		Experienced     bool           `json:"experienced"`
// // 		Class_timing    datatypes.JSON `json:"class_timing"`
// // 		Flexi_timing    bool           `json:"flexi_timing"`
// // 		Enquiree_name   string         `json:"enquiree_name"`
// // 		Data_not_taken  bool           `json:"data_not_taken"`
// // 		Enq_for_someone bool           `json:"enq_for_someone"`
// // 		Mode_of_enquiry datatypes.JSON `json:"mode_of_enquiry"`
// // 	}
// // 	var enqenquiry enquiryBasic
// // 	json.Unmarshal(input.Enquiry, &enqenquiry)

// // 	if enqenquiry.Comment != "incubation" {
// // 		if input.Rid > 0 {
// // 			go updateNoEnquiryForm(r.ID, userid)
// // 		}
// // 		go updateRecs(r.ID)
// // 		go updateRefs(r.ID)
// // 		go addCFNEnquiry(r.ID, userid, true)
// // 		Tagenquiry(r.ID, r.Contact, r.Email)
// // 	}
// // 	go Tagleadenqstu(int(r.ID), r.Name, input.Contact, input.Email, "enquiry")
// // 	// go Updatestudentconverted(r)
// // 	if err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// // 		return
// // 	}
// // 	json.Unmarshal(input.Contact, &contacts)
// // 	json.Unmarshal(input.Email, &emails)
// // 	go TagEnquiryToWebisteEnquiries(&contacts, int(r.ID), &emails)
// // 	// Tatsupdatee(r.ID)
// // 	// Walkintat(r.Courses, r.Author, r.CreatedAt)
// // 	c.JSON(http.StatusOK, gin.H{"message": "Added enquiry " + r.Name + " " + fmt.Sprint(r.ID)})

// }

func done() {
	fmt.Println("Haiiii")
}
