package handler

import (
	"akuntansi/helper"
	"akuntansi/model"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDataSelectTutor(c *gin.Context, db *gorm.DB) {
	tutorId := c.Query("q")
	var tutor []model.TutorSelect
	err := db.
		Where("name like ?", "%"+tutorId+"%").
		Find(&tutor).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan tutor", http.StatusOK, tutor)
	c.JSON(http.StatusOK, formatter)
}

func IndexTutor(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "tutor.html", gin.H{"Info": session.Get("info"), "Error": session.Get("error"), "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL")})
}

func NewTutor(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/tutor/tutor_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string)}); err != nil {
		fmt.Println(err)
	}

}

func EditTutor(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var tutor model.Tutor
	var user model.User

	err := db.Where("id=?", id).First(&tutor).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	db.Where("id=?", tutor.UserID).First(&user)

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/tutor/tutor_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"UserID": tutor.UserID, "EmailData": user.Email, "Data": tutor, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string)}); err != nil {
		fmt.Println(err)
	}

}

func CreateTutor(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputTutor
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/tutor")
		return
	}
	inputTutor.CreatedBy = int32(userID)
	err = db.Debug().Create(&inputTutor).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/tutor")
		return
	}

	c.Redirect(http.StatusFound, "/tutor")
}

func UpdateTutor(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputTutor
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	id := c.Param("id")
	session := sessions.Default(c)
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/tutor")
		return
	}

	err = db.Debug().Model(&inputTutor).Where("id=?", id).Updates(&inputTutor).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/tutor")
		return
	}
	db.Debug().Table("tutor").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	})
	c.Redirect(http.StatusFound, "/tutor")
}

func GetDataTutor(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameTutor(orderColumnIdx)

	var totalRecords int64
	var tutors []model.Tutor

	searchQuery, queryParams := buildSearchQueryTutor(searchValue)
	query := db.Model(&model.Tutor{}).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&tutors)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            tutors,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameTutor(idx int) string {
	columnsMapping := map[int]string{
		2: "name",
		3: "is_active",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryTutor(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "name LIKE ? OR is_active LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}
