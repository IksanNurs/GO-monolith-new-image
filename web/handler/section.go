package handler

import (
	"akuntansi/model"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewSection(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	id := c.Param("id")
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/section/section_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "ID": id, "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditSection(c *gin.Context, db *gorm.DB) {
	var maxSequence int64
	dataSequence := []int{}
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var section model.Section

	err := db.Where("id=?", id).First(&section).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	result := db.Debug().Table("section").Select("Count(sequence) as count_sequence").Where("version_id=?", section.VersionID).Count(&maxSequence)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	for i := 1; i <= int(maxSequence); i++ {
		dataSequence = append(dataSequence, i)
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/section/section_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Sequence": dataSequence, "Data": section, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}

}

func CreateSection(c *gin.Context, db *gorm.DB) {
	var inputSection model.InputSection
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	err := c.ShouldBind(&inputSection)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/version/view/"+strconv.Itoa(int(inputSection.VersionID)))
		return
	}

	var maxSequence int64
	result := db.Debug().Table("section").Select("Count(sequence) as count_sequence").Where("version_id=?", inputSection.VersionID).Count(&maxSequence)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	inputSection.Sequence = int32(maxSequence) + 1
	inputSection.CreatedBy = int32(userID)
	err = db.Debug().Create(&inputSection).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/version/view/"+strconv.Itoa(int(inputSection.VersionID)))
		return
	}

	c.Redirect(http.StatusFound, "/version/view/"+strconv.Itoa(int(inputSection.VersionID)))
}

func UpdateSection(c *gin.Context, db *gorm.DB) {
	var inputSection model.InputSection
	var sequenceSection model.SequenceSection
	id := c.Param("id")
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	session := sessions.Default(c)
	err := c.ShouldBind(&inputSection)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/version/view/"+strconv.Itoa(int(inputSection.VersionID)))
		return
	}

	old, _ := strconv.Atoi(c.PostForm("old_sequence"))

	if inputSection.Sequence != int32(old) {
		db.Where("sequence=? AND version_id=?", inputSection.Sequence, inputSection.VersionID).First(&sequenceSection)
		db.Debug().Table("section").Where("id=?", sequenceSection.ID).Updates(map[string]interface{}{
			"sequence": old,
		})

	}
	err = db.Debug().Model(&inputSection).Where("id=?", id).Updates(&inputSection).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/version/view/"+strconv.Itoa(int(inputSection.VersionID)))
		return
	}
	db.Debug().Table("section").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	})
	c.Redirect(http.StatusFound, "/version/view/"+strconv.Itoa(int(inputSection.VersionID)))
}
