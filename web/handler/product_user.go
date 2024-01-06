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

func EditProductUser(c *gin.Context, db *gorm.DB) {
	id := c.Query("id")
	session := sessions.Default(c)
	var package1 model.ProductUser

	err := db.Preload("User").Preload("Product").Where("id=?", id).First(&package1).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/product_user/productuser_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": package1, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func GetDataSelectProduct(c *gin.Context, db *gorm.DB) {
	versionId := c.Query("q")
	var version []model.ProductSelect
	err := db.
		Where("name like ?", "%"+versionId+"%").
		Order("id desc").
		Find(&version).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan product select", http.StatusOK, version)
	c.JSON(http.StatusOK, formatter)
}

func NewProductUser(c *gin.Context) {

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/product_user/productuser_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
		fmt.Println(err)
	}

}

func CreateProductUser(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputProductUser
	var product model.Product
	IsPrice, _ := strconv.Atoi(c.PostForm("is_price"))
	session := sessions.Default(c)
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	err = db.Where("id=?", inputTutor.ProductID).First(&product).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	if product.Stock != 0 {
		db.Debug().Table("product").Where("id=?", inputTutor.ProductID).Updates(map[string]interface{}{
			"stock": product.Stock - *inputTutor.Quantity,
		})
	} else {
		session.Set("error", "product "+product.Name+" tidak memiliki persedian stok")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	if IsPrice == 0 {
		inputTutor.PriceMember = 0
		inputTutor.PriceNonmember = product.PriceNonmember
	} else {
		inputTutor.PriceMember = product.PriceMember
		inputTutor.PriceNonmember = 0
	}
	err = db.Debug().Create(&inputTutor).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}

	c.Redirect(http.StatusFound, "/product-user")
}

func UpdateProductUser(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputProductUser
	var product model.Product
	var productuser model.ProductUser
	id := c.Param("id")
	IsPrice, _ := strconv.Atoi(c.PostForm("is_price"))
	session := sessions.Default(c)
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	err = db.Where("id=?", inputTutor.ProductID).First(&product).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	err = db.Where("id=?", id).First(&productuser).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	if product.Stock != 0 {
		k := 0
		if *inputTutor.Quantity < productuser.Quantity {
			k = int(productuser.Quantity - *inputTutor.Quantity)
			k += int(product.Stock)
		} else {
			k = int(*inputTutor.Quantity - productuser.Quantity)
			k = int(product.Stock) - k
		}
		fmt.Println(k)
		if k < 0 {
			session.Set("error", "product "+product.Name+" tidak memiliki stok yang cukup")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/product-user")
			return
		}
		db.Debug().Table("product").Where("id=?", inputTutor.ProductID).Updates(map[string]interface{}{
			"stock": k,
		})
	} else {
		session.Set("error", "product "+product.Name+" tidak memiliki persedian stok")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}
	if IsPrice == 0 {
		inputTutor.PriceMember = 0
		inputTutor.PriceNonmember = product.PriceNonmember
	} else {
		inputTutor.PriceMember = product.PriceMember
		inputTutor.PriceNonmember = 0
	}

	err = db.Debug().Model(&inputTutor).Where("id=?", id).Updates(&inputTutor).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/product")
		return
	}

	c.Redirect(http.StatusFound, "/product-user")
}

func IndexProductUser(c *gin.Context, db *gorm.DB) {
	var productusers []model.ProductUser
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	query := db.Debug().Model(&model.ProductUser{}).
		Preload("Product").
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
	var tot float64 = 0
	query.Find(&productusers)
	for i := range productusers {
		fmt.Println("Quantity:", productusers[i].Quantity)
		fmt.Println("PV:", productusers[i].Product.Pv)
		tot += (float64(productusers[i].Quantity) * float64(productusers[i].Product.Pv)) * 5 / 100
	}

	formattedPrice := strconv.FormatFloat(tot, 'f', 2, 64)

	// Pembulatan nilai untuk harga
	formattedPriceFloat, _ := strconv.ParseFloat(formattedPrice, 64)

	// Format dengan pemisah ribuan
	formattedString := formatCurrency(formattedPriceFloat)
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "product_user.html", gin.H{"tot": formattedString, "Error": session.Get("error"), "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}
func formatCurrency(price float64) string {
	formatted := fmt.Sprintf("%.2f", price)

	// Pisahkan ribuan dengan tanda titik
	for i := len(formatted) - 6; i > 0; i -= 3 {
		formatted = formatted[:i] + "." + formatted[i:]
	}
	return formatted
}

func GetDataProductUser(c *gin.Context, db *gorm.DB) {
	month, _ := strconv.Atoi(c.PostForm("month"))
	year, _ := strconv.Atoi(c.PostForm("year"))
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	orderColumn := getColumnProductUser(orderColumnIdx)

	var totalRecords int64
	var productusers []model.ProductUser
	searchQuery, queryParams := buildSearchQueryProductUser(searchValue)
	query := db.Debug().Model(&model.ProductUser{}).
		Preload("User").
		Preload("Product").
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

func getColumnProductUser(idx int) string {
	columnsMapping := map[int]string{
		2: "user_id",
		3: "product_id",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryProductUser(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "product_id LIKE ? OR  user_id LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func DeleteProductUser(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var package1 model.ProductUser
	id := c.Query("id")

	err := db.Where("id=?", id).Delete(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/product-user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/product-user")
}
