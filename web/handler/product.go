package handler

import (
	"akuntansi/model"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewProduct(c *gin.Context) {

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/product/product_new.html"))
	session := sessions.Default(c)
	userID := session.Get("id").(int32)
	if err := tmpl.Execute(c.Writer, gin.H{"userID": userID, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
		fmt.Println(err)
	}

}

func NewProductProductStock(c *gin.Context) {

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/product/product_new_stock.html"))
	session := sessions.Default(c)
	userID := session.Get("id").(int32)
	if err := tmpl.Execute(c.Writer, gin.H{"userID": userID, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
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
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	dateParsed, err := time.ParseInLocation("2006-01-02", c.PostForm("created_at"), jakartaLocation)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp := dateParsed.Unix()
	inputTutor.CreatedAt = timestamp
	inputTutor.TotalStock = 0
	if inputTutor.Stock != nil {
		inputTutor.TotalStock = *inputTutor.Stock
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

func CreateProductStock(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputProduct
	var product model.Product
	session := sessions.Default(c)
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product")
		return
	}
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	dateParsed, err := time.ParseInLocation("2006-01-02", c.PostForm("created_at"), jakartaLocation)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	db.Where("id=?", c.PostForm("product_id")).First(&product)

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp := dateParsed.Unix()
	inputTutor.CreatedAt = timestamp
	inputTutor.Name = product.Name
	inputTutor.PriceMember = product.PriceMember
	inputTutor.PriceNonmember = product.PriceNonmember
	inputTutor.Pv = product.Pv
	inputTutor.TotalStock = product.TotalStock
	if inputTutor.Stock != nil {
		inputTutor.TotalStock = *inputTutor.Stock + product.TotalStock
	}
	db.Debug().Table("product").Where("name=?", product.Name).Updates(map[string]interface{}{
		"total_stock": inputTutor.TotalStock,
	})
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
	var product model.Product
	session := sessions.Default(c)
	id := c.Param("id")
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product")
		return
	}
	if c.PostForm("created_at") != "" {
		jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			session.Set("error", err.Error())
			session.Save()
			c.Redirect(http.StatusSeeOther, "/product-user")
			return
		}
		fmt.Println(c.PostForm("created_at"))
		dateParsed, err := time.ParseInLocation("2006-01-02", c.PostForm("created_at"), jakartaLocation)
		if err != nil {
			session.Set("error", err.Error())
			session.Save()
			c.Redirect(http.StatusSeeOther, "/product-user")
			return
		}

		// Mengonversi time.Time ke timestamp UNIX (int64)
		timestamp := dateParsed.Unix()
		inputTutor.CreatedAt = timestamp
	}
	db.Where("id=?", id).First(&product)
	if inputTutor.Stock != nil && (*inputTutor.Stock > product.Stock) {
		fmt.Println("hg +")
		db.Debug().Table("product").Where("name=?", product.Name).Updates(map[string]interface{}{
			"total_stock": product.TotalStock + (*inputTutor.Stock - product.Stock),
		})
	}
	if inputTutor.Stock != nil && (*inputTutor.Stock < product.Stock) {
		fmt.Println("hg")
		fmt.Println(product.Name)
		db.Debug().Table("product").Where("name=?", product.Name).Updates(map[string]interface{}{
			"total_stock": product.TotalStock - (product.Stock-*inputTutor.Stock),
		})
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
	userID := session.Get("id").(int32)
	c.HTML(http.StatusOK, "product.html", gin.H{"userID": userID, "Error": session.Get("error"), "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func GetDataProduct(c *gin.Context, db *gorm.DB) {
	month, _ := strconv.Atoi(c.PostForm("month"))
	day, _ := strconv.Atoi(c.PostForm("day"))
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
		Where(searchQuery, queryParams...)

	if day != 0 && month != 0 && year != 0 {
		query = query.Where("DAY(FROM_UNIXTIME(created_at)) = ?", day).
			Where("MONTH(FROM_UNIXTIME(created_at)) = ?", month).
			Where("YEAR(FROM_UNIXTIME(created_at)) = ?", year)
	} else if day != 0 && month != 0 {
		query = query.Where("DAY(FROM_UNIXTIME(created_at)) = ?", day).
			Where("MONTH(FROM_UNIXTIME(created_at)) = ?", month)
	} else if month != 0 && year != 0 {
		query = query.Where("MONTH(FROM_UNIXTIME(created_at)) = ?", month).
			Where("YEAR(FROM_UNIXTIME(created_at)) = ?", year)
	} else if day != 0 {
		query = query.Where("DAY(FROM_UNIXTIME(created_at)) = ?", day)
	} else if month != 0 {
		query = query.Where("MONTH(FROM_UNIXTIME(created_at)) = ?", month)
	} else if year != 0 {
		query = query.Where("YEAR(FROM_UNIXTIME(created_at)) = ?", year)
	}

	query = query.Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&productusers)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	for i := range productusers {
		jakartaLocation, err1 := time.LoadLocation("Asia/Jakarta")
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}
		if productusers[i].CreatedAt != 0 {
			date := time.Unix(int64(productusers[i].CreatedAt), 0)
			date = date.In(jakartaLocation)
			productusers[i].CreatedAt_t = date.Format("2006-01-02 15:04")
		}
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
	var package2 model.Product
	id := c.Query("id")

	db.Where("id=?", id).First(&package2)
	db.Debug().Table("product").Where("name=?", package2.Name).Updates(map[string]interface{}{
		"total_stock": package2.TotalStock - package2.Stock,
	})
	err := db.Where("id=?", id).Delete(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product")
		return
	}

	c.Redirect(http.StatusSeeOther, "/product")
}
