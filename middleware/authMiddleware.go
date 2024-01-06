package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("id") == nil {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
		userID := session.Get("id").(int32)
		re := strings.Replace(os.Getenv("ALLOWED_USERS"), " ", "", -1)
		spl := strings.Split(re, ",")
		// for _, h := range spl {
		// 	tr := strings.Trim(h, " ")
		// 	if tr == "2" {
		// 		fmt.Println(" 2" + "-" + hh)
		// 	}
		// }

		if !slices.Contains(spl, strconv.Itoa(int(userID))) {
			session.Clear()
			session.Save()
			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		c.Set("userData", nil)
		c.Next()
	}
}
