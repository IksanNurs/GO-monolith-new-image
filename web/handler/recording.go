package handler

import (
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

func NewRecording(c *gin.Context, db *gorm.DB) {
	id := c.Query("id")
	bd_id := c.Query("bd_id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var recording model.Recording
	// err := db.Debug().
	// 	Preload("Meeting").
	// 	Where("meeting_id=?", id).
	// 	First(&recording).Error

	// if err != nil {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/recording/recording_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": recording, "bundleID": bd_id, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditRecording(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var recording model.Recording
	query := db.Model(&model.Recording{}).
		Preload("Meeting").
		Where("id=?", id).
		First(&recording)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/recording/recording_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": recording, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func CreateRecording(c *gin.Context, db *gorm.DB) {
	var inputRecording model.InputRecording
	session := sessions.Default(c)
	// bundle_id := c.Param("id")
	err := c.ShouldBind(&inputRecording)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/meeting/view/"+strconv.Itoa(int(inputRecording.MeetingID)))
		return
	}

	err = db.Debug().Create(&inputRecording).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/meeting/view/"+strconv.Itoa(int(inputRecording.MeetingID)))
		return
	}

	c.Redirect(http.StatusFound, "/bundle/meeting/view/"+strconv.Itoa(int(inputRecording.MeetingID)))
}

func UpdateRecording(c *gin.Context, db *gorm.DB) {
	var inputRecording model.InputRecording
	session := sessions.Default(c)
	id := c.Param("id")
	err := c.ShouldBind(&inputRecording)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/meeting/view/"+strconv.Itoa(int(inputRecording.MeetingID)))
		return
	}

	err = db.Debug().Model(&inputRecording).Where("id=?", id).Updates(&inputRecording).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/meeting/view/"+strconv.Itoa(int(inputRecording.MeetingID)))
		return
	}

	c.Redirect(http.StatusFound, "/bundle/meeting/view/"+strconv.Itoa(int(inputRecording.MeetingID)))
}

func GetDataRecording(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	meetingID, _ := strconv.Atoi(c.Param("id"))
	orderColumn := getColumnNameRecording(orderColumnIdx)

	var totalRecords int64
	var recordings []model.Recording

	searchQuery, queryParams := buildSearchQueryRecording(searchValue)
	query := db.Debug().Model(&model.Recording{}).
		Preload("Meeting").
		Where(searchQuery, queryParams...).
		Where("meeting_id=?", meetingID).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&recordings)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            recordings,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameRecording(idx int) string {
	columnsMapping := map[int]string{
		2: "meeting_id",
		3: "name",
		4: "file",
		5: "url",
		6: "embeb",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryRecording(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "meeting_id LIKE ? OR name LIKE ? OR file LIKE ? OR url LIKE ? OR embeb LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func ActionDeleteRecording(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var recording model.Recording
	id := c.Query("id")
	me_id := c.Query("me_id")

	err := db.Where("id=?", id).Delete(&recording).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/bundle/meeting/view/"+me_id)
		return
	}

	c.Redirect(http.StatusSeeOther, "/bundle/meeting/view/"+me_id)
}
