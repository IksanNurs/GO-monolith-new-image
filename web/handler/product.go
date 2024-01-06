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

func NewProduct(c *gin.Context) {

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/product/product_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
		fmt.Println(err)
	}

}

func EditProduct(c *gin.Context, db *gorm.DB) {
	id := c.Query("id")
	session := sessions.Default(c)
	var package1 model.Product

	err := db.Where("id=?", id).First(&package1).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/product/product_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": package1, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func CreateProduct(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputProduct
	session := sessions.Default(c)
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product")
		return
	}
	err = db.Debug().Create(&inputTutor).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product")
		return
	}

	c.Redirect(http.StatusFound, "/product")
}

func UpdateProduct(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputProduct
	session := sessions.Default(c)
	id := c.Param("id")
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product")
		return
	}
	err = db.Debug().Model(&inputTutor).Where("id=?", id).Updates(&inputTutor).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/product")
		return
	}

	c.Redirect(http.StatusFound, "/product")
}

func IndexProduct(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "product.html", gin.H{"Error": session.Get("error"), "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func GetDataProduct(c *gin.Context, db *gorm.DB) {
	month, _ := strconv.Atoi(c.PostForm("month"))
	year, _ := strconv.Atoi(c.PostForm("year"))
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	orderColumn := getColumnProduct(orderColumnIdx)

	var totalRecords int64
	var productusers []model.Product
	searchQuery, queryParams := buildSearchQueryProduct(searchValue)
	query := db.Debug().Model(&model.Product{}).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&productusers)
	if month != 0 && year != 0 {
		query = query.Where("YEAR(FROM_UNIXTIME(created_at)) = ?", year).
			Where("MONTH(FROM_UNIXTIME(created_at)) = ?", month)
	} else {
		if month != 0 {
			query = query.Where("MONTH(FROM_UNIXTIME(created_at)) = ?", month)
		}
		if year != 0 {
			query = query.Where("YEAR(FROM_UNIXTIME(created_at)) = ?", year)
		}
	}

	query.Count(&totalRecords).
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

func getColumnProduct(idx int) string {
	columnsMapping := map[int]string{
		2: "name",
		3: "stock",
		4: "pv",
		5: "price_member",
		6: "price_nonmember",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryProduct(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "name LIKE ? OR stock LIKE ? OR pv LIKE ? OR price_member LIKE ? OR price_nonmember LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func DeleteProduct(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var package1 model.Product
	id := c.Query("id")

	err := db.Where("id=?", id).Delete(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product")
		return
	}

	c.Redirect(http.StatusSeeOther, "/product")
}
