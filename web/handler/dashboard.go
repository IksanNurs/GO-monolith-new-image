package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) {
	// session := sessions.Default(c)
	// email := session.Get("email").(string)
	var countRequest int64
	var countExam int64
	var countOrder int64
	err := db.Debug().
		Table("product").
		Count(&countRequest).
		Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = db.Debug().
		Table("product_user").
		Count(&countExam).
		Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = db.Debug().
		Table("user").
		Count(&countOrder).
		Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{"countOrder": countOrder, "countExam": countExam, "countRequest": countRequest, "AuthURL": os.Getenv("AUTH_ADMIN_URL")})
}
