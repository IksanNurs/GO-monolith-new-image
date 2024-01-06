package handler

import (
	"akuntansi/helper"
	"akuntansi/model"
	"bytes"
	"fmt"
	"html/template"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDataSelectMaterial(c *gin.Context, db *gorm.DB) {
	materialName := c.Query("q")
	var tutor []model.MaterialSelect
	err := db.
		Where("name like ?", "%"+materialName+"%").
		Limit(20).
		Find(&tutor).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan material", http.StatusOK, tutor)
	c.JSON(http.StatusOK, formatter)
}

func IndexMaterial(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)

	c.HTML(http.StatusOK, "material.html", gin.H{"Error": session.Get("error"), "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func NewMaterial(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/material/material_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func NewMaterialContent(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var material model.Material

	err := db.Where("id=?", id).First(&material).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/material/material_content_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": material, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditMaterialContent(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var materialcontent model.MaterialContent
	query := db.Model(&model.MaterialContent{}).
		Where("id=?", id).
		First(&materialcontent)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/material/material_content_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": materialcontent, "NilToEmptyString": NilToZeroValue, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditMaterial(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var material model.Material
	query := db.Model(&model.Material{}).
		Where("id=?", id).
		First(&material)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/material/material_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": material, "NilToEmptyString": NilToZeroValue, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func ViewMaterial(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var material model.Material
	var bundlematerial []model.BundleMaterial
	query := db.Model(&model.Material{}).
		Where("id=?", id).
		First(&material)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	query = db.
		Preload("Bundle").
		Where("material_id=?", id).
		Find(&bundlematerial)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	replaceNewlinesWithBr := func(description string) template.HTML {
		return template.HTML(strings.Replace(description, "\n", "<br/>", -1))
	}

	c.HTML(http.StatusOK, "material_view.html", gin.H{"DataBundle": bundlematerial, "ReplaceNewlinesWithBr": replaceNewlinesWithBr, "AWS_DESTINATION_PREFIX": os.Getenv("AWS_DESTINATION_PREFIX"), "PUBLIC_URL_S3": os.Getenv("PUBLIC_URL_S3"), "Data": material, "NilToEmptyString": NilToZeroValue, "id": id, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")})
}

// func ViewMaterial(c *gin.Context, db *gorm.DB) {
// 	id := c.Param("id")
// 	session := sessions.Default(c)
// 	email := session.Get("email").(string)
// 	var material model.Material
// 	query := db.Model(&model.Material{}).
// 		Preload("Bundle").
// 		Preload("Tutor").
// 		Where("id=?", id).
// 		First(&material)

// 	if query.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
// 		return
// 	}

// 	date := time.Unix(material.StartedAt, 0)
// 	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
// 	if err != nil {
// 		// Handle error jika gagal memuat zona waktu
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
// 		return
// 	}
// 	date = date.In(jakartaLocation)
// 	material.StartedAtTe = date.Format("2006-01-02 15:04")

// 	date1 := time.Unix(material.EndedAt, 0)
// 	date1 = date1.In(jakartaLocation)
// 	material.EndedAtTe = date1.Format("2006-01-02 15:04")

// 	c.HTML(http.StatusOK, "material_view.html", gin.H{ "NilToEmptyString": NilToZeroValue, "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host), "JWTToken": session.Get("token").(string), "Data": material, "TutorID": material.Tutor.ID, "TutorName": material.Tutor.Name, "email": email, "AuthURL":os.Getenv("AUTH_URL"), "Info": session.Get("info"), "Error": session.Get("error")})
// }

func CreateMaterialContent(c *gin.Context, db *gorm.DB) {
	var inputMaterialContent model.InputMaterialContent
	sess := c.MustGet("sess").(*s3.S3)
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int32(userData["user_id"].(float64))

	batch := time.Now().UTC().Unix()
	name := strconv.Itoa(int(batch))
	err := c.ShouldBind(&inputMaterialContent)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/material/view/"+strconv.Itoa(int(inputMaterialContent.MaterialID)))
		return
	}
	inputMaterialContent.CreatedBy = &userID
	err = db.Create(&inputMaterialContent).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/material/view/"+strconv.Itoa(int(inputMaterialContent.MaterialID)))
		return
	}

	file1, header_image, err1 := c.Request.FormFile("file_material")
	if err1 != nil {
		file1 = nil
	}
	if file1 != nil {
		fmt.Println("kjfgd")
		UploadS3Material(sess, file1, header_image, name, db, int(inputMaterialContent.ID))
	}
	c.Redirect(http.StatusFound, "/material/view/"+strconv.Itoa(int(inputMaterialContent.MaterialID)))
}

func UpdateMaterialContent(c *gin.Context, db *gorm.DB) {
	var inputMaterialContent model.InputMaterialContent
	sess := c.MustGet("sess").(*s3.S3)
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int32(userData["user_id"].(float64))
	id := c.Param("id")

	batch := time.Now().UTC().Unix()
	name := strconv.Itoa(int(batch))
	err := c.ShouldBind(&inputMaterialContent)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material/view/"+strconv.Itoa(int(inputMaterialContent.MaterialID)))
		return
	}
	err = db.Model(&inputMaterialContent).Where("id=?", id).Updates(&inputMaterialContent).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material/view/"+strconv.Itoa(int(inputMaterialContent.MaterialID)))
		return
	}
	err = db.Table("material_content").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	}).Error

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material/view/"+strconv.Itoa(int(inputMaterialContent.MaterialID)))
		return
	}
	file1, header_image, err1 := c.Request.FormFile("file_material")
	if err1 != nil {
		file1 = nil
	}
	if file1 != nil {
		fmt.Println("kjfgd")
		id1, _ := strconv.Atoi(id)
		UploadS3Material(sess, file1, header_image, name, db, id1)
	}
	c.Redirect(http.StatusFound, "/material/view/"+strconv.Itoa(int(inputMaterialContent.MaterialID)))
}

func CreateMaterial(c *gin.Context, db *gorm.DB) {
	var inputMaterial model.InputMaterial
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int32(userData["user_id"].(float64))
	err := c.ShouldBind(&inputMaterial)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material")
		return
	}
	inputMaterial.CreatedBy = &userID
	err = db.Create(&inputMaterial).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material")
		return
	}

	c.Redirect(http.StatusFound, "/material")
}

func UpdateMaterial(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var inputMaterial model.InputMaterial
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	err := c.ShouldBind(&inputMaterial)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material")
		return
	}

	err = db.Model(&inputMaterial).Where("id=?", id).Updates(&inputMaterial).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material")
		return
	}
	err = db.Table("material").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	}).Error

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material")
		return
	}
	c.Redirect(http.StatusFound, "/material")
}

func ActionDeleteMaterial(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var material model.Material
	id := c.Param("id")

	err := db.Where("id=?", id).Delete(&material).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/material")
		return
	}

	c.Redirect(http.StatusSeeOther, "/material")
}

func GetDataMaterial(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	token := session.Get("token").(string)
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameMaterial(orderColumnIdx)

	var totalRecords int64
	var materials []model.Material

	searchQuery, queryParams := buildSearchQueryMaterial(searchValue, db, token)
	query := db.Table("material").
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&materials)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            materials,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameMaterial(idx int) string {
	columnsMapping := map[int]string{
		2: "id",
		3: "name",
		4: "description",
		5: "is_open",
		6: "is_active",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryMaterial(searchValue string, db *gorm.DB, token string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "id LIKE ? OR name LIKE ? OR description LIKE ? OR is_open LIKE ? OR is_active LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func GetDataMaterialContent(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	token := session.Get("token").(string)
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	id := c.Param("id")
	orderColumn := getColumnNameMaterialContent(orderColumnIdx)

	var totalRecords int64
	var materialcontents []model.MaterialContent

	searchQuery, queryParams := buildSearchQueryMaterialContent(searchValue, db, token)
	query := db.Table("material_content").
		Where("material_id=?", id).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&materialcontents)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            materialcontents,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameMaterialContent(idx int) string {
	columnsMapping := map[int]string{
		2: "id",
		3: "name",
		4: "type",
		5: "`group`",
		6: "file",
		7: "sequence",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryMaterialContent(searchValue string, db *gorm.DB, token string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "name LIKE ? OR description LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func UploadS3Material(sess *s3.S3, file multipart.File, header_image *multipart.FileHeader, name string, db *gorm.DB, id int) {
	size := header_image.Size
	var path string
	_ = header_image.Filename
	buffer := make([]byte, size)
	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	db.Table("material_content").Where("id=?", id).Updates(map[string]interface{}{
		"file": name + "_" + header_image.Filename,
	})
	path = os.Getenv("AWS_DESTINATION_PREFIX") + "/material_content/" + name + "_" + header_image.Filename
	params := &s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("AWS_BUCKET")),
		Key:           aws.String(path),
		ACL:           aws.String("public-read"),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	_, err := sess.PutObject(params)
	if err != nil {
		// handle error
		fmt.Println(err)

	}

}

func ActionDeleteMaterialContent(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var materialcontent model.MaterialContent
	id := c.Query("id")
	m_id := c.Query("m_id")

	err := db.Where("id=?", id).Delete(&materialcontent).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/material/view/"+m_id)
		return
	}

	c.Redirect(http.StatusSeeOther, "/material/view/"+m_id)
}
