package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// Course เป็นโครงสร้างข้อมูลคอร์ส
type Course struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Instructor string  `json:"instructor"`
}

func main() {
	// กำหนด DSN สำหรับเชื่อมต่อฐานข้อมูล MySQL
	dsn := "root:123456789@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"


	// เชื่อมต่อกับฐานข้อมูล MySQL
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("เชื่อมต่อฐานข้อมูลไม่สำเร็จ: " + err.Error())
	}

	// สร้าง Table (ถ้ายังไม่มี)
	err = db.AutoMigrate(&Course{})
	if err != nil {
		panic("ไม่สามารถสร้างตารางในฐานข้อมูล: " + err.Error())
	}

	// สร้างเซิร์ฟเวอร์ Gin
	r := gin.Default()

	// กำหนด Middleware และ CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
	}))

	// สร้างเส้นทาง API
	v1 := r.Group("/api/v1")
	{
		v1.GET("/courses", GetCourses)         // เส้นทางสำหรับดึงข้อมูลคอร์สทั้งหมด
		v1.GET("/courses/:id", GetCourseByID)  // เส้นทางสำหรับดึงข้อมูลคอร์สตาม ID
		v1.POST("/courses", CreateCourse)      // เส้นทางสำหรับสร้างคอร์ส
	}

	// รันเซิร์ฟเวอร์ที่พอร์ต 8080
	r.Run(":8080")
}

// GetCourses เป็น handler ที่ให้ข้อมูลทั้งหมดของคอร์ส
func GetCourses(c *gin.Context) {
	var courses []Course
	// ดึงข้อมูลทั้งหมดจากฐานข้อมูล
	if err := db.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลคอร์สได้"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

// GetCourseByID เป็น handler ที่ให้ข้อมูลของคอร์สตาม ID ที่ระบุ
func GetCourseByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var course Course
	// ดึงข้อมูลคอร์สตาม ID จากฐานข้อมูล
	if err := db.Where("id = ?", id).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบคอร์ส"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

// CreateCourse เป็น handler ที่ให้สร้างคอร์สใหม่
func CreateCourse(c *gin.Context) {
	var newCourse Course

	// ดึงข้อมูล JSON จาก request body
	if err := c.ShouldBindJSON(&newCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
		return
	}

	// บันทึกข้อมูลลงในฐานข้อมูล
	if err := db.Create(&newCourse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างคอร์สได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newCourse})
}
