package handler

import (
	"akuntansi/helper"
	"akuntansi/model"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDataSelectMeeting(c *gin.Context, db *gorm.DB) {
	meetingName := c.Query("q")
	var tutor []model.MeetingSelect
	err := db.
		Where("name like ?", "%"+meetingName+"%").
		Find(&tutor).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan meeting", http.StatusOK, tutor)
	c.JSON(http.StatusOK, formatter)
}

func NewMeeting(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundle model.Bundle

	err := db.Where("id=?", id).First(&bundle).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/meeting/meeting_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": bundle, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditMeeting(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var meeting model.Meeting
	query := db.Model(&model.Meeting{}).
		Preload("Bundle").
		Preload("Tutor").
		Where("id=?", id).
		First(&meeting)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	date := time.Unix(meeting.StartedAt, 0)
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle error jika gagal memuat zona waktu
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	date = date.In(jakartaLocation)
	meeting.StartedAtTe = date.Format("2006-01-02 15:04")

	date1 := time.Unix(meeting.EndedAt, 0)
	date1 = date1.In(jakartaLocation)
	meeting.EndedAtTe = date1.Format("2006-01-02 15:04")

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/meeting/meeting_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": meeting, "NilToEmptyString": NilToZeroValue, "TutorID": meeting.Tutor.ID, "TutorName": meeting.Tutor.Name, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func ViewMeeting(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var meeting model.Meeting
	query := db.Model(&model.Meeting{}).
		Preload("Bundle").
		Preload("Tutor").
		Where("id=?", id).
		First(&meeting)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	date := time.Unix(meeting.StartedAt, 0)
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle error jika gagal memuat zona waktu
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	date = date.In(jakartaLocation)
	meeting.StartedAtTe = date.Format("2006-01-02 15:04")

	date1 := time.Unix(meeting.EndedAt, 0)
	date1 = date1.In(jakartaLocation)
	meeting.EndedAtTe = date1.Format("2006-01-02 15:04")

	c.HTML(http.StatusOK, "meeting_view.html", gin.H{"NilToEmptyString": NilToZeroValue, "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host), "JWTToken": session.Get("token").(string), "Data": meeting, "TutorID": meeting.Tutor.ID, "TutorName": meeting.Tutor.Name, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")})
}

func CreateMeeting(c *gin.Context, db *gorm.DB) {
	var inputMeeting model.InputMeeting
	session := sessions.Default(c)
	startedatt := c.PostForm("started_at_t")
	endedatt := c.PostForm("ended_at_t")
	err := c.ShouldBind(&inputMeeting)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}
	dateParsed, err := time.ParseInLocation("2006-01-02T15:04", startedatt, jakartaLocation)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp := dateParsed.Unix()
	inputMeeting.StartedAt = timestamp

	dateParsed1, err1 := time.ParseInLocation("2006-01-02T15:04", endedatt, jakartaLocation)
	if err1 != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp1 := dateParsed1.Unix()
	inputMeeting.EndedAt = timestamp1

	err = db.Debug().Create(&inputMeeting).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}

	c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
}

func UpdateMeeting(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var inputMeeting model.InputMeeting
	session := sessions.Default(c)
	startedatt := c.PostForm("started_at_t")
	endedatt := c.PostForm("ended_at_t")
	err := c.ShouldBind(&inputMeeting)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}
	dateParsed, err := time.ParseInLocation("2006-01-02T15:04", startedatt, jakartaLocation)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp := dateParsed.Unix()
	inputMeeting.StartedAt = timestamp

	dateParsed1, err1 := time.ParseInLocation("2006-01-02T15:04", endedatt, jakartaLocation)
	if err1 != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp1 := dateParsed1.Unix()
	inputMeeting.EndedAt = timestamp1

	err = db.Debug().Model(&inputMeeting).Where("id=?", id).Updates(&inputMeeting).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
		return
	}

	c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputMeeting.BundleID)))
}

func ActionDeleteMeeting(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var meeting model.Meeting
	id := c.Query("id")
	bd_id := c.Query("bd_id")

	err := db.Where("id=?", id).Delete(&meeting).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
		return
	}

	c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
}
