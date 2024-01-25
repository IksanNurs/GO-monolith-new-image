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

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EditUser(c *gin.Context, db *gorm.DB) {
	id := c.Query("id")
	session := sessions.Default(c)
	var package1 model.User

	err := db.Where("id=?", id).First(&package1).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/user/users_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": package1, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func NewUser(c *gin.Context) {

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/user/users_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
		fmt.Println(err)
	}

}

func CreateUser(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputUser
	session := sessions.Default(c)
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}
	err = db.Debug().Create(&inputTutor).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	c.Redirect(http.StatusFound, "/user")
}

func UpdateUser(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputUser
	session := sessions.Default(c)
	id := c.Param("id")
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}
	err = db.Debug().Model(&inputTutor).Where("id=?", id).Updates(&inputTutor).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/user")
		return
	}

	c.Redirect(http.StatusFound, "/user")
}

func IndexUser(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "users.html", gin.H{"Error": session.Get("error"), "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func ActionGetAllUserCPNS(c *gin.Context, db *gorm.DB) {
	var User []model.UserSelect
	params := c.Query("q")
	err := db.
		Where("name like ?", "%"+params+"%").
		Order("id desc").
		Limit(20).
		Find(&User).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan user", http.StatusOK, gin.H{"user": User})
	c.JSON(http.StatusOK, formatter)
}

func ActionGetAllUserCPNSByID(c *gin.Context, db *gorm.DB) {
	var User model.User
	params := c.Param("id")
	err := db.
		Where("id=?", params).
		First(&User).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan user", http.StatusOK, gin.H{"user": User})
	c.JSON(http.StatusOK, formatter)
}

func GetDataUser(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	orderColumn := getColumnUser(orderColumnIdx)

	var totalRecords int64
	var productusers []model.User
	searchQuery, queryParams := buildSearchQueryUser(searchValue)
	query := db.Debug().Model(&model.User{}).
		Where(searchQuery, queryParams...).
		Where("id!=?",1).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&productusers)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            productusers,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnUser(idx int) string {
	columnsMapping := map[int]string{
		2: "email",
		3: "name",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryUser(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "email LIKE ? OR name LIKE ? OR id LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}
func DeleteUser(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var package1 model.User
	id := c.Query("id")

	err := db.Where("id=?", id).Delete(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/user")
}
