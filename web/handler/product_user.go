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
	"time"

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

	session := sessions.Default(c)
	userID := session.Get("id").(int32)
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/product_user/productuser_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"userID":userID,"AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
		fmt.Println(err)
	}

}

func CreateProductUser(c *gin.Context, db *gorm.DB) {
	var inputTutor model.InputProductUser
	var product model.Product
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
	if *inputTutor.Quantity != productuser.Quantity {
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
	session := sessions.Default(c)
	userID := session.Get("id").(int32)
	var productusers []model.ProductUser
	month, _ := strconv.Atoi(c.Query("month"))
	day, _ := strconv.Atoi(c.Query("day"))
	year, _ := strconv.Atoi(c.Query("year"))
	query := db.Debug().Model(&model.ProductUser{}).
		Preload("Product").
		Find(&productusers)

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

	var tot float64 = 0
	var totmember float64 = 0
	var totnonmember float64 = 0
	var totpaid float64 = 0
	var totunpaid float64 = 0
	query.Find(&productusers)
	for i := range productusers {
		fmt.Println("Quantity:", productusers[i].Quantity)
		fmt.Println("PV:", productusers[i].Product.Pv)
		tot += (float64(productusers[i].Quantity) * float64(productusers[i].Product.Pv)) * 5 / 100
		if productusers[i].CategoriPrice == 1 {
			totmember += (float64(productusers[i].Quantity) * float64(productusers[i].Product.PriceMember))
		}
		if productusers[i].CategoriPrice == 2 {
			di := float64(0)
			d := float64(productusers[i].Product.PriceNonmember) * float64(productusers[i].Quantity)
			if productusers[i].Diskon > 0 {
				dis := ((float64(productusers[i].Diskon) / 100) * float64(productusers[i].Product.PriceNonmember))
				di = d - dis
				totnonmember += di
				continue
			}
			totnonmember += d
		}
		if productusers[i].Paid != 0 {
			totpaid += float64(productusers[i].Paid)
		}
		if productusers[i].Unpaid != 0 {
			totunpaid += float64(productusers[i].Unpaid)
		}
	}

	formattedPrice := strconv.FormatFloat(tot, 'f', 2, 64)
	formattedPrice1 := strconv.FormatFloat(totmember, 'f', 2, 64)
	formattedPrice2 := strconv.FormatFloat(totnonmember, 'f', 2, 64)
	formattedPrice3 := strconv.FormatFloat(totpaid, 'f', 2, 64)
	formattedPrice4 := strconv.FormatFloat(totunpaid, 'f', 2, 64)

	// Pembulatan nilai untuk harga
	formattedPriceFloat, _ := strconv.ParseFloat(formattedPrice, 64)
	formattedPriceFloat1, _ := strconv.ParseFloat(formattedPrice1, 64)
	formattedPriceFloat2, _ := strconv.ParseFloat(formattedPrice2, 64)
	formattedPriceFloat3, _ := strconv.ParseFloat(formattedPrice3, 64)
	formattedPriceFloat4, _ := strconv.ParseFloat(formattedPrice4, 64)

	// Format dengan pemisah ribuan
	formattedString := formatCurrency(formattedPriceFloat)
	formattedString1 := formatCurrency(formattedPriceFloat1)
	formattedString2 := formatCurrency(formattedPriceFloat2)
	formattedString3 := formatCurrency(formattedPriceFloat3)
	formattedString4 := formatCurrency(formattedPriceFloat4)
	c.HTML(http.StatusOK, "product_user.html", gin.H{"userID":userID,"totpaid": formattedString3, "totunpaid": formattedString4, "totnonmember": formattedString2, "totmember": formattedString1, "tot": formattedString, "Error": session.Get("error"), "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
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
	day, _ := strconv.Atoi(c.PostForm("day"))
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

	query.Count(&totalRecords).
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
			// Handle error jika gagal memuat zona waktu
			c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
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
