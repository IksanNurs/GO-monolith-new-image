package handler

import (
	"akuntansi/model"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func ClearSessionInfoAndError(c *gin.Context) {
	// Hapus sesi di sini
	session := sessions.Default(c)
	session.Delete("info")
	session.Delete("error")
	session.Delete("success")
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Session cleared"})
}


func IndexLogin(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("id") != nil {
		c.Redirect(http.StatusSeeOther, "/dashboard")
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/auth/login.html"))
	if err := tmpl.Execute(c.Writer, gin.H{"AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
		fmt.Println(err)
	}
}

func OauthToken(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	user := model.User{}
	if session.Get("id") == nil {
		trimEmail := strings.Trim(c.PostForm("email"), " ")
		trimPassword := strings.Trim(c.PostForm("password"), " ")
		err := db.Where("email=? AND password=?", trimEmail, trimPassword).First(&user).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error_title": "Gagal Login", "error_text": "email atau password tidak sesuai"})
			return
		}
		re := strings.Replace(os.Getenv("ALLOWED_USERS"), " ", "", -1)
		spl := strings.Split(re, ",")
		// for _, h := range spl {
		// 	tr := strings.Trim(h, " ")
		// 	if tr == "2" {
		// 		fmt.Println(" 2" + "-" + hh)
		// 	}
		// }

		if !slices.Contains(spl, strconv.Itoa(int(user.ID))) {
			c.JSON(http.StatusBadRequest, gin.H{"error_title": "Gagal Login", "error_text": "anda tidak memiliki akses, hubungi admin"})
			return
		}
		session.Set("id", user.ID)
		session.Save()
	}

	c.JSON(http.StatusOK, nil)
}
