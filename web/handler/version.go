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
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDataSelectVersion(c *gin.Context, db *gorm.DB) {
	versionId := c.Query("q")
	var version []model.VersionSelect
	err := db.
		Where("name like ?", "%"+versionId+"%").
		Order("id desc").
		Find(&version).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan version", http.StatusOK, version)
	c.JSON(http.StatusOK, formatter)
}

func IndexVersion(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "version.html", gin.H{"Info": session.Get("info"), "Error": session.Get("error"), "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL")})
}

func NewVersion(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/version/version_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditVersion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var version model.Version

	err := db.Where("id=?", id).First(&version).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/version/version_edit.html"))
	if err := tmpl.Execute(c.Writer, gin.H{"Data": version, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}

}

func ViewVersion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var version model.Version

	err := db.Where("id=?", id).First(&version).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "version_view.html", gin.H{"Data": version, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")})
}

func CreateVersion(c *gin.Context, db *gorm.DB) {
	var inputVersion model.Version
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int32(userData["user_id"].(float64))
	err := c.ShouldBind(&inputVersion)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/version")
		return
	}
	inputVersion.CreatedBy = &userID
	err = db.Debug().Create(&inputVersion).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/version")
		return
	}

	c.Redirect(http.StatusFound, "/version/view/"+strconv.Itoa(int(inputVersion.ID)))
}

func UpdateVersion(c *gin.Context, db *gorm.DB) {
	var inputVersion model.Version
	id := c.Param("id")
	session := sessions.Default(c)
	err := c.ShouldBind(&inputVersion)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/version")
		return
	}

	err = db.Debug().Model(&inputVersion).Where("id=?", id).Updates(&inputVersion).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/version")
		return
	}
	db.Debug().Table("version").Where("id=?", id).Updates(map[string]interface{}{
		"updated_by": userID,
	})
	c.Redirect(http.StatusFound, "/version/view/"+id)
}

func GetDataVersion(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameVersion(orderColumnIdx)

	var totalRecords int64
	var versions []model.Version

	searchQuery, queryParams := buildSearchQueryVersion(searchValue)
	query := db.Model(&model.Version{}).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&versions)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            versions,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameVersion(idx int) string {
	columnsMapping := map[int]string{
		2: "name",
		3: "type",
		4: "is_active",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryVersion(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		if strings.ToLower(searchValue) == "active" {
			searchQuery := "is_active=1"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "inactive" {
			searchQuery := "is_active=0"
			var params []interface{}
			return searchQuery, params
		}
		searchQuery := "name LIKE ? OR type LIKE ? OR is_active LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func GetDataSection(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	id := c.Param("id")

	fmt.Println(orderColumnIdx)
	orderColumn := getColumnNameSection(orderColumnIdx)

	var totalRecords int64
	var versions []model.Section

	searchQuery, queryParams := buildSearchQuerySection(searchValue)
	query := db.Model(&model.Section{}).
		Preload("Version").
		Where("version_id=?", id).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&versions)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            versions,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameSection(idx int) string {
	columnsMapping := map[int]string{
		2: "sequence",
		3: "name",
		4: "minimum_score",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "sequence"
	}
	return colName
}

func buildSearchQuerySection(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "version_id LIKE ? OR sequence LIKE ? OR name LIKE ? OR minimum_score LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}
