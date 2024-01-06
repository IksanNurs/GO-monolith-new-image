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
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDataSelectBundle(c *gin.Context, db *gorm.DB) {
	packageId := c.Query("q")
	var bundle []model.BundleSelect
	err := db.
		Where("name like ?", "%"+packageId+"%").
		Find(&bundle).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan bundle", http.StatusOK, bundle)
	c.JSON(http.StatusOK, formatter)
}

func NewBundle(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}

}

func NewBundlePackage(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundle model.Bundle

	err := db.Where("id=?", id).First(&bundle).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_package_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": bundle, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditBundlePackage(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundlePackage model.BundlePackage
	query := db.Model(&model.BundlePackage{}).
		Preload("Package").
		Where("id=?", id).
		First(&bundlePackage)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	date := time.Unix(int64(bundlePackage.DateStart), 0)
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle error jika gagal memuat zona waktu
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	date = date.In(jakartaLocation)
	bundlePackage.Starte = date.Format("2006-01-02 15:04")
	date1 := time.Unix(int64(bundlePackage.DateEnd), 0)
	date1 = date1.In(jakartaLocation)
	bundlePackage.Ende = date1.Format("2006-01-02 15:04")

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_package_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": bundlePackage, "PackageID": bundlePackage.Package.ID, "PackageName": bundlePackage.Package.Name, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditBundle(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundle model.Bundle

	err := db.Where("id=?", id).First(&bundle).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": bundle, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditBundleAddMaterial(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundle model.Bundle

	err := db.Where("id=?", id).First(&bundle).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_add_material.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": bundle, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditBundleUpdateMaterial(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundlematerial model.BundleMaterial

	err := db.Preload("Bundle").Preload("Material").Where("id=?", id).First(&bundlematerial).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_edit_material.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"id": id, "MaterialName": bundlematerial.Material.Name, "MaterialID": bundlematerial.MaterialID, "Data": bundlematerial, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditBundleAddUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundle model.Bundle

	err := db.Where("id=?", id).First(&bundle).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_add_user.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": bundle, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditBundleEditUser(c *gin.Context, db *gorm.DB) {
	id := c.Query("id")
	order_id := c.Query("order_id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var bundle model.Bundle
	var order model.Order
	var user model.User
	err := db.Where("id=?", id).First(&bundle).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	err1 := db.Where("id=?", order_id).First(&order).Error
	if err1 != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	err = db.Where("id=?", order.UserID).First(&user).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/bundle/bundle_edit_user.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"UserID": order.UserID, "DataOrder": order, "Data": bundle, "EmailData": user.Email, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func ViewBundle(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	var bundle model.Bundle

	err := db.Where("id=?", id).First(&bundle).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	replaceNewlinesWithBr := func(description string) template.HTML {
		return template.HTML(strings.Replace(description, "\n", "<br/>", -1))
	}

	c.HTML(http.StatusOK, "bundle_view.html", gin.H{"UserIDT": userID, "BundleName": bundle.Name, "ReplaceNewlinesWithBr": replaceNewlinesWithBr, "URL_PRESENCE": os.Getenv("URL_PRESENCE"), "Start": os.Getenv("URL_PRESENCE_START"), "End": os.Getenv("URL_PRESENCE_END"), "Data": bundle, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func IndexBundle(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var totalRecords int64
	query := db.Model(&model.Bundle{}).
		Where("is_active=?", 1).
		Count(&totalRecords)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	c.HTML(http.StatusOK, "bundle.html", gin.H{"Info": session.Get("info"), "Error": session.Get("error"), "totalrecords": totalRecords, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL")})
}

func CreateBundlePackage(c *gin.Context, db *gorm.DB) {
	var inputBundlePackage model.InputBundlePackage
	session := sessions.Default(c)
	startedatt := c.PostForm("starte")
	endedatt := c.PostForm("ende")
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	err := c.ShouldBind(&inputBundlePackage)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}
	fmt.Println(startedatt)
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}
	dateParsed, err := time.ParseInLocation("2006-01-02T15:04", startedatt, jakartaLocation)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp := dateParsed.Unix()
	inputBundlePackage.DateStart = int32(timestamp)

	dateParsed1, err1 := time.ParseInLocation("2006-01-02T15:04", endedatt, jakartaLocation)
	if err1 != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp1 := dateParsed1.Unix()
	inputBundlePackage.DateEnd = int32(timestamp1)
	inputBundlePackage.CreatedBy = int32(userID)
	err = db.Debug().Create(&inputBundlePackage).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}

	c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
}

func UpdateBundlePackage(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	startedatt := c.PostForm("starte")
	endedatt := c.PostForm("ende")
	var inputBundlePackage model.InputBundlePackage
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	session := sessions.Default(c)
	err := c.ShouldBind(&inputBundlePackage)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}
	fmt.Println(startedatt)
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}
	dateParsed, err := time.ParseInLocation("2006-01-02T15:04", startedatt, jakartaLocation)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp := dateParsed.Unix()
	inputBundlePackage.DateStart = int32(timestamp)

	dateParsed1, err1 := time.ParseInLocation("2006-01-02T15:04", endedatt, jakartaLocation)
	if err1 != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}

	// Mengonversi time.Time ke timestamp UNIX (int64)
	timestamp1 := dateParsed1.Unix()
	inputBundlePackage.DateEnd = int32(timestamp1)
	err = db.Debug().Model(&inputBundlePackage).Where("id=?", id).Updates(&inputBundlePackage).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
		return
	}
	db.Debug().Table("bundle_package").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	})

	c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundlePackage.BundleID)))
}

func CreateBundle(c *gin.Context, db *gorm.DB) {
	var inputBundle model.InputBundle
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	err := c.ShouldBind(&inputBundle)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundle.ID)))
		return
	}
	inputBundle.CreatedBy = int32(userID)
	err = db.Debug().Create(&inputBundle).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundle.ID)))
		return
	}

	c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundle.ID)))
}

func UpdateBundle(c *gin.Context, db *gorm.DB) {
	var inputBundle model.InputBundle
	id := c.Param("id")
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	err := c.ShouldBind(&inputBundle)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+id)
		return
	}

	err = db.Debug().Model(&inputBundle).Where("id=?", id).Updates(&inputBundle).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+id)
		return
	}

	db.Debug().Table("bundle").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	})
	c.Redirect(http.StatusFound, "/bundle/view/"+id)
}

func CreateBundleMaterial(c *gin.Context, db *gorm.DB) {
	var inputBundleMaterial model.InputBundleMaterial
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int32(userData["user_id"].(float64))
	err := c.ShouldBind(&inputBundleMaterial)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundleMaterial.BundleID)))
		return
	}
	inputBundleMaterial.CreatedBy = &userID
	err = db.Debug().Create(&inputBundleMaterial).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundleMaterial.BundleID)))
		return
	}

	c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundleMaterial.BundleID)))
}

func UpdateBundleMaterial(c *gin.Context, db *gorm.DB) {
	var inputBundleMaterial model.InputBundleMaterial
	session := sessions.Default(c)
	id := c.Param("id")
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int32(userData["user_id"].(float64))
	err := c.ShouldBind(&inputBundleMaterial)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundleMaterial.BundleID)))
		return
	}

	err = db.Debug().Model(&inputBundleMaterial).Where("id=?", id).Updates(&inputBundleMaterial).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundleMaterial.BundleID)))
		return
	}

	db.Debug().Table("bundle").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	})
	c.Redirect(http.StatusFound, "/bundle/view/"+strconv.Itoa(int(inputBundleMaterial.BundleID)))
}
func GetDataBundle(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameBundle(orderColumnIdx)

	var totalRecords int64
	var bundles []model.Bundle

	searchQuery, queryParams := buildSearchQueryBundle(searchValue)
	query := db.Model(&model.Bundle{}).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&bundles)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            bundles,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameBundle(idx int) string {
	columnsMapping := map[int]string{
		2: "id",
		3: "name",
		4: "price",
		5: "type",
		6: "is_public",
		7: "is_active",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryBundle(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		if strings.ToLower(searchValue) == "active" {
			searchQuery := "is_active=1"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "inactive" {
			searchQuery := "is_active=0"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "public" {
			searchQuery := "is_public=1"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "private" {
			searchQuery := "is_public=0"
			var params []interface{}
			return searchQuery, params
		}
		searchQuery := "name LIKE ? OR description LIKE ? OR price LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func GetDataBundlePackage(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	bundleID, _ := strconv.Atoi(c.Param("id"))

	orderColumn := getColumnNameBundlePackage(orderColumnIdx)

	var totalRecords int64
	var bundlespackage []model.BundlePackage

	searchQuery, queryParams := buildSearchQueryBundlePackage(searchValue)
	query := db.Model(&model.BundlePackage{}).
		Preload("Package").
		Preload("Bundle").
		Where(searchQuery, queryParams...).
		Where("bundle_id=?", bundleID).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&bundlespackage)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	for i := range bundlespackage {
		countExam := int64(0)
		jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			// Handle error jika gagal memuat zona waktu
			c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
			return
		}
		date := time.Unix(int64(bundlespackage[i].DateStart), 0)
		date = date.In(jakartaLocation)
		bundlespackage[i].Starte = date.Format("2006-01-02 15:04")
		date1 := time.Unix(int64(bundlespackage[i].DateEnd), 0)
		date1 = date1.In(jakartaLocation)
		bundlespackage[i].Ende = date1.Format("2006-01-02 15:04")
		err = db.Debug().
			Table("exam").
			Joins("LEFT JOIN `order` ON exam.order_id=order.id").
			Where("exam.package_id=? AND order.bundle_id=?", bundlespackage[i].PackageID, bundlespackage[i].BundleID).
			Count(&countExam).
			Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		bundlespackage[i].Count = int(countExam)
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            bundlespackage,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameBundlePackage(idx int) string {
	columnsMapping := map[int]string{
		2: "id",
		3: "package_id",
		4: "date_start",
		5: "date_end",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryBundlePackage(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "id LIKE ? OR bundle_id LIKE ? OR package_id LIKE ? OR date_start LIKE ? OR date_end LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func GetDataBundleMeeting(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	bundleID, _ := strconv.Atoi(c.Param("id"))
	orderColumn := getColumnNameBundleMeeting(orderColumnIdx)

	var totalRecords int64
	var meetings []model.Meeting

	searchQuery, queryParams := buildSearchQueryBundleMeeting(searchValue)
	query := db.Model(&model.Meeting{}).
		Preload("Bundle").
		Preload("Tutor").
		Where(searchQuery, queryParams...).
		Where("bundle_id=?", bundleID).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&meetings)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	for i := range meetings {
		countMeeting := int64(0)
		countRecording := int64(0)
		// Mengonversi timestamp UNIX ke zona waktu "Asia/Jakarta" pada StartedAt
		date := time.Unix(meetings[i].StartedAt, 0)
		jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			// Handle error jika gagal memuat zona waktu
			c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
			return
		}
		date = date.In(jakartaLocation)
		meetings[i].StartedAtTe = date.Format("2006-01-02 15:04")

		// Mengonversi timestamp UNIX ke zona waktu "Asia/Jakarta" pada EndedAt
		date1 := time.Unix(meetings[i].EndedAt, 0)
		date1 = date1.In(jakartaLocation)
		meetings[i].EndedAtTe = date1.Format("2006-01-02 15:04")
		err = db.Debug().
			Table("presence").
			Where("start IS NOT NULL AND meeting_id=?", meetings[i].ID).
			Count(&countMeeting).
			Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		meetings[i].Count = int(countMeeting)
		err = db.Debug().
			Table("recording").
			Where("meeting_id=?", meetings[i].ID).
			Count(&countRecording).
			Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if countRecording > 0 {
			meetings[i].IsRecording = 1
		}
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            meetings,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameBundleMeeting(idx int) string {
	columnsMapping := map[int]string{
		2: "sequence",
		3: "name",
		4: "started_at",
		5: "ended_at",
		6: "tutor_id",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryBundleMeeting(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "bundle_id LIKE ? OR name LIKE ? OR started_at LIKE ? OR ended_at LIKE ? OR tutor_id LIKE ? OR url LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func GetDataOrderBundle(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	ID := c.Param("id")

	orderColumn := getColumnNameOrder1(orderColumnIdx)

	var totalRecords int64
	var orders []model.Order
	searchQuery, queryParams := buildSearchQueryOrder1(searchValue, db, ID)
	query := db.Model(&model.Order{}).
		Preload("Bundle").
		Preload("UserFromID").
		Preload("UserFromID.Institution").
		Preload("UserFromCreatedBy").
		Select("`order`.*, _user.name, _user.phone, _user.email, _institution.name, if (paid_at is not null, paid_at, if (expired_at is not null, expired_at, if (cancelled_at is not null, cancelled_at, if (checked_out_at is not null, checked_out_at, order.created_at)))) as order_time").
		Joins("LEFT JOIN _user on (order.user_id = _user.id)").
		Joins("LEFT JOIN _institution on (_user.education_institution_id = _institution.id)").
		Where("bundle_id=? AND ((paid_at IS NOT NULL AND checked_out_at IS NOT NULL) OR checked_out_at IS NULL)", ID).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&orders)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle error jika gagal memuat zona waktu
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	for i := range orders {
		// e, err := FetchDataUser(int(orders[i].UserID), token)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// orders[i].Name = &e.Data.User.Name
		// orders[i].Email = &e.Data.User.Email
		// orders[i].Phone = &e.Data.User.Phone
		// orders[i].InstitutionName = &e.Data.User.InstitutionName
		if orders[i].OrderTime != nil {
			date := time.Unix(int64(*orders[i].OrderTime), 0)
			date = date.In(jakartaLocation)
			orders[i].PaidString = date.Format("2006-01-02 15:04")
		}
		// if orders[i].CreatedBy != nil {
		// 	e1, err1 := FetchDataUser(int(*orders[i].CreatedBy), token)
		// 	if err1 != nil {
		// 		fmt.Println(err.Error())
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	orders[i].NameCreatedBy = &e1.Data.User.Name
		// }
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            orders,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func GetDataBundleMaterial(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	bundleID, _ := strconv.Atoi(c.Param("id"))
	orderColumn := getColumnNameBundleMaterial(orderColumnIdx)

	var totalRecords int64
	var bundlematerials []model.BundleMaterial

	searchQuery, queryParams := buildSearchQueryBundleMaterial(searchValue)
	query := db.Model(&model.BundleMaterial{}).
		Preload("Material").
		Where(searchQuery, queryParams...).
		Where("bundle_id=?", bundleID).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&bundlematerials)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            bundlematerials,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameBundleMaterial(idx int) string {
	columnsMapping := map[int]string{
		2: "id",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryBundleMaterial(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "id LIKE ?"
		params := []interface{}{"%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func ActionDeleteBundlePackage(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var bundle_package model.BundlePackage
	var bundle_package1 []model.BundlePackage
	var ekam model.Exam
	id := c.Query("id")
	bd_id := c.Query("bd_id")
	pkg_id := c.Query("pkg_id")

	db.Debug().Where("bundle_id=? AND package_id=?", bd_id, pkg_id).Find(&bundle_package1)
	if len(bundle_package1) > 1 {
		err := db.Where("id=?", id).Delete(&bundle_package).Error
		if err != nil {
			session.Set("error", err.Error())
			session.Save()
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
			return
		}
	} else {
		err := db.Debug().Joins("LEFT JOIN `order` on exam.order_id = order.id").Where("order.bundle_id = ? and exam.package_id = ?", bd_id, pkg_id).First(&ekam).Error
		if err != nil {
			err = db.Where("id=?", id).Delete(&bundle_package).Error
			if err != nil {
				session.Set("error", err.Error())
				session.Save()
				c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
				return
			}
		}
	}

	c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
}

func ActionDeleteBundleOrder(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var order model.Order
	id := c.Query("id")
	bd_id := c.Query("bd_id")

	err := db.Where("id=?", id).Delete(&order).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
		return

	}

	c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
}

func ActionDeleteBundlePackageExam(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var ekam model.ExamDelete
	id := c.Query("id")
	bd_id := c.Query("bd_id")
	err := db.Debug().Where("package_id=?", id).Delete(&ekam).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
		return
	}

	c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
}

func ActionDeleteBundleMaterial(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var ekam model.BundleMaterial
	id := c.Query("id")
	bd_id := c.Query("bd_id")
	err := db.Debug().Where("id=?", id).Delete(&ekam).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
		return
	}

	c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
}
