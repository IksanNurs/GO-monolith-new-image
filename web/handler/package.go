package handler

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"akuntansi/helper"
	"akuntansi/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IndexPackage(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var totalRecords int64
	query := db.Model(&model.Package{}).
		Where("is_independent=1").
		Count(&totalRecords)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	c.HTML(http.StatusOK, "package.html", gin.H{"totalIndependent": totalRecords, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")})
}

func NewPackage(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/package/package_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditPackage(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var package1 model.Package

	err := db.Preload("Version").Where("id=?", id).First(&package1).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/package/package_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"VersionName": package1.Version.Name, "VersionID": package1.VersionID, "Data": package1, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditPackageSection(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var package1 model.Package

	err := db.Preload("Version").Preload("Version.Section").Where("id=?", id).First(&package1).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/package/package_edit_section.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": package1, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func ViewPackage(c *gin.Context, db *gorm.DB) {
	packageId := c.Param("id")
	id, _ := strconv.Atoi(packageId)
	var oldPackage model.Package
	var packageQuestionCount []model.PackageQuestionCount
	err := db.Preload("Version").
		Where("id=?", id).
		Find(&oldPackage).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	email := session.Get("email").(string)

	replaceNewlinesWithBr := func(description string) template.HTML {
		return template.HTML(strings.Replace(description, "\n", "<br/>", -1))
	}

	err = db.Debug().Preload("Section").Select("package_question.*, count(id) as count").Where("package_id=?", id).Group("section_id").Find(&packageQuestionCount).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "package_view.html", gin.H{"PackageName": oldPackage.Name, "packageQuestionCount": packageQuestionCount, "ReplaceNewlinesWithBr": replaceNewlinesWithBr, "Error": session.Get("error"), "Success": session.Get("success"), "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "AWS_DESTINATION_PREFIX": os.Getenv("AWS_DESTINATION_PREFIX"), "PUBLIC_URL_S3": os.Getenv("PUBLIC_URL_S3"), "oldPackage": oldPackage})
}

func CreatePackage(c *gin.Context, db *gorm.DB) {
	var inputPackage model.InputPackage
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	err := c.ShouldBind(&inputPackage)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/package/view/"+strconv.Itoa(int(inputPackage.ID)))
		return
	}
	inputPackage.CreatedBy = int32(userID)
	err = db.Debug().Create(&inputPackage).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/package/view/"+strconv.Itoa(int(inputPackage.ID)))
		return
	}

	c.Redirect(http.StatusFound, "/package/view/"+strconv.Itoa(int(inputPackage.ID)))
}

func UpdatePackage(c *gin.Context, db *gorm.DB) {
	var inputPackage model.InputPackage
	id := c.Param("id")
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	session := sessions.Default(c)
	err := c.ShouldBind(&inputPackage)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/package/view/"+id)
		return
	}

	err = db.Debug().Model(&inputPackage).Where("id=?", id).Updates(&inputPackage).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/package/view/"+id)
		return
	}
	db.Debug().Table("package").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	})

	c.Redirect(http.StatusFound, "/package/view/"+id)
}

func PackageQuestion(c *gin.Context, db *gorm.DB) {
	packageId := c.Param("id")
	id, _ := strconv.Atoi(packageId)
	var section []model.Section
	var version model.PackageVersion
	err := db.
		Where("id=?", id).
		First(&version).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "package_question_new.html", gin.H{"Error": err.Error()})
		return
	}

	err = db.Debug().
		Where("version_id=?", version.VersionID).
		Find(&section).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "package_question_new.html", gin.H{"Error": err.Error()})
		return
	}

	session := sessions.Default(c)
	email := session.Get("email").(string)
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/package/package_question_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "package_id": id, "section": section, "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func ActionPackageQuestion(c *gin.Context, db *gorm.DB) {
	packageId := c.Param("id")
	id, _ := strconv.Atoi(packageId)
	sectionid, _ := strconv.Atoi(c.PostForm("section_id"))
	var pquestion []model.PackageQuestion1
	var question model.Question1
	session := sessions.Default(c)
	data := []string{}
	// email := session.Get("email").(string)

	questionsid := strings.Split(c.PostForm("id_question"), " ")
	for _, r := range questionsid {
		if strings.Contains(r, "-") {
			questionsiddash := strings.Split(r, "-")
			start, _ := strconv.Atoi(questionsiddash[0])
			end, _ := strconv.Atoi(questionsiddash[1])
			for i := start; i <= end; i++ {
				err := db.
					Where("package_id=?", id).
					Find(&pquestion).Error
				if err != nil {
					session.Set("error", err.Error())
					session.Save()
					c.Redirect(http.StatusSeeOther, "/package/view/"+packageId)
					return
				}

				db.Where("id=?", i).First(&question)
				if question.IsActive != 1 {
					data = append(data, strconv.Itoa(i))
				} else {
					if question.MergedTo != nil {
						i = int(*question.MergedTo)
					}
					sequence := int32(len(pquestion) + 1)
					packagequestion := model.PackageQuestion1{
						PackageID:  int32(id),
						Sequence:   &sequence,
						QuestionID: int32(i),
						SectionID:  int32(sectionid),
					}
					err = db.Create(&packagequestion).Error
					if err != nil {
						data = append(data, strconv.Itoa(i))
					}
				}
			}
		} else {
			err := db.
				Where("package_id=?", id).
				Find(&pquestion).Error
			if err != nil {
				session.Set("error", err.Error())
				session.Save()
				c.Redirect(http.StatusSeeOther, "/package/view/"+packageId)
				return
			}

			sequence := int32(len(pquestion) + 1)
			idQuestion, _ := strconv.Atoi(r)
			db.Where("id=?", idQuestion).First(&question)
			if question.IsActive != 1 {
				data = append(data, strconv.Itoa(idQuestion))
			} else {
				if question.MergedTo != nil {
					idQuestion = int(*question.MergedTo)
				}
				packagequestion := model.PackageQuestion1{
					PackageID:  int32(id),
					Sequence:   &sequence,
					QuestionID: int32(idQuestion),
					SectionID:  int32(sectionid),
				}
				err = db.Create(&packagequestion).Error
				if err != nil {
					data = append(data, strconv.Itoa(idQuestion))
				}
			}
		}

	}

	if len(data) != 0 {
		formatted := strings.Join(data, ", ")
		session.Set("info", "id question "+formatted+" tidak ada atau tidak aktif")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package/view/"+packageId)
		return
	}

	c.Redirect(http.StatusSeeOther, "/package/view/"+packageId)
}

func GetDataPackage(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNamePackage(orderColumnIdx)

	var totalRecords int64
	var questions []model.Package

	searchQuery, queryParams := buildSearchQueryPackage(searchValue)
	query := db.Model(&model.Package{}).
		Preload("Version").
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&questions)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            questions,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNamePackage(idx int) string {
	columnsMapping := map[int]string{
		2: "id",
		3: "name",
		4: "duration",
		5: "price",
		6: "shuffle_type",
		7: "is_independent",
		8: "is_active",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryPackage(searchValue string) (string, []interface{}) {
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
		if strings.ToLower(searchValue) == "lengkap" {
			searchQuery := "is_complete=1"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "mini to" {
			searchQuery := "is_complete=0"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "berurutan" {
			searchQuery := "shuffle_type=1"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "acak keseluruhan" {
			searchQuery := "shuffle_type=2"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "acak per section" {
			searchQuery := "shuffle_type=3"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "acak per topik" {
			searchQuery := "shuffle_type=4"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "dependent" {
			searchQuery := "is_independent=0"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "independent" {
			searchQuery := "is_independent=1"
			var params []interface{}
			return searchQuery, params
		}
		searchQuery := "id LIKE ? OR name LIKE ? OR description LIKE ? OR duration LIKE ? OR price LIKE ? OR version_id LIKE ? OR scoring_method LIKE ? OR minimum_score LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func GetDataPackageQuestion(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	packageID, _ := strconv.Atoi(c.Param("id"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNamePackageQuestion(orderColumnIdx)

	var totalRecords int64
	var questions []model.PackageQuestion

	searchQuery, queryParams := buildSearchQueryPackageQuestion(searchValue, db)
	query := db.Model(&model.PackageQuestion{}).
		Preload("Question").
		Preload("Question.Tutor").
		Preload("Section").
		Preload("Question.Option").
		Joins("LEFT JOIN question q ON package_question.question_id=q.id").
		Joins("LEFT JOIN section s ON package_question.section_id=s.id").
		Where(searchQuery, queryParams...).
		Where("package_id=?", packageID).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&questions)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	ekam := model.Exam{}
	var del int32
	err := db.Debug().Where("package_id=?", packageID).First(&ekam).Error
	if err != nil {
		del = 1
	} else {
		del = 0
	}
	for i := range questions {
		questions[i].IsDelete = del

	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            questions,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNamePackageQuestion(idx int) string {
	columnsMapping := map[int]string{
		0: "package_question.sequence",
		2: "q.question",
		3: "s.name",
		4: "q.explanation",
		5: "q.topic",
		6: "q.tutor_name",
		7: "package_question.question_id",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "q.question"
	}
	return colName
}

func buildSearchQueryPackageQuestion(searchValue string, db *gorm.DB) (string, []interface{}) {
	if searchValue != "" {
		// Ubah logika pencarian untuk mencari nama Section
		var sectionIDs []int
		db.Model(&model.Section{}).Where("name = ? OR name LIKE ?", searchValue, "%"+searchValue+"%").Pluck("id", &sectionIDs)
		var tutors []int
		db.Model(&model.Tutor{}).Where("name = ? OR name LIKE ?", searchValue, "%"+searchValue+"%").Pluck("id", &tutors)
		searchQuery := "q.tutor_id IN (?) OR section_id IN (?) OR question_id LIKE ? OR package_question.sequence LIKE ? OR q.batch LIKE ? OR q.question LIKE ? OR q.explanation LIKE ? OR q.topic LIKE ? OR q.subtopic LIKE ? OR q.tutor_name LIKE ?"
		params := []interface{}{tutors, sectionIDs, "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func GetDataSelectPackage(c *gin.Context, db *gorm.DB) {
	packageId := c.Query("q")
	var package1 []model.PackageSelect
	err := db.
		Where("name like ? AND is_active=?", "%"+packageId+"%", 1).
		Find(&package1).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	formatter := helper.APIResponse("berhasil menampilkan data package", http.StatusOK, package1)
	c.JSON(http.StatusOK, formatter)
}

// func GetDataSelectPackageBundle(c *gin.Context, db *gorm.DB) {
// 	packageId := c.Query("q")
// 	var package1 []model.PackageSelect
// 	err := db.
// 		Where("name like ? AND 	is_independent=?", "%"+packageId+"%", 0).
// 		Find(&package1).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
// 		return
// 	}
// 	formatter := helper.APIResponse("berhasil menampilkan data package", http.StatusOK, package1)
// 	c.JSON(http.StatusOK, formatter)
// }

func ActionPackageSection(c *gin.Context, db *gorm.DB) {
	packageId := c.Param("id")
	session := sessions.Default(c)
	data := []string{}
	// email := session.Get("email").(string)
	var package1 model.Package

	err := db.Preload("Version").Preload("Version.Section").Where("id=?", packageId).First(&package1).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	for i := range package1.Version.Section {
		if c.PostForm("sequence["+strconv.Itoa(i)+"]") != "" {
			indsequence := strings.Split(c.PostForm("sequence["+strconv.Itoa(i)+"]"), " ")
			for _, r := range indsequence {
				if strings.Contains(r, "-") {
					questionsiddash := strings.Split(r, "-")
					start, _ := strconv.Atoi(questionsiddash[0])
					end, _ := strconv.Atoi(questionsiddash[1])
					for j := start; j <= end; j++ {
						err := db.Debug().Table("package_question").Where("package_id=? AND sequence=?", packageId, j).Updates(map[string]interface{}{
							"section_id": package1.Version.Section[i].ID,
						}).Error
						if err != nil {
							data = append(data, strconv.Itoa(j))
						}
					}
				} else {
					numsequwnce, _ := strconv.Atoi(r)
					err := db.Debug().Table("package_question").Where("package_id=? AND sequence=?", packageId, numsequwnce).Updates(map[string]interface{}{
						"section_id": package1.Version.Section[i].ID,
					}).Error
					if err != nil {
						data = append(data, strconv.Itoa(numsequwnce))
					}

				}

			}
		}

	}

	if len(data) != 0 {
		formatted := strings.Join(data, ", ")
		session.Set("info", "sequence "+formatted+" tidak ada")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package/view/"+packageId)
		return
	}

	c.Redirect(http.StatusSeeOther, "/package/view/"+packageId)
}

func ActionDeletePackage(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var package1 model.Package
	id := c.Param("id")

	err := db.Where("id=?", id).Delete(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package")
		return
	}

	c.Redirect(http.StatusSeeOther, "/package")
}

func ActionDuplicatePackage(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	// email := session.Get("email").(string)
	var package1 model.InputPackage
	var packageQuestion []model.InputPackageQuestion
	id := c.Param("id")

	err := db.Where("id=?", id).First(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package/view/"+id)
		return
	}
	package1.Name = package1.Name + " - Duplicated"
	ero := int32(0)
	package1.IsIndependent = &ero
	package1.ID = 0
	err = db.Create(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package/view/"+id)
		return
	}
	err = db.Where("package_id=?", id).Find(&packageQuestion).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package/view/"+id)
		return
	}

	for _, pq := range packageQuestion {
		pq.PackageID = package1.ID
		pq.CreatedBy = int32(userID)
		err = db.Create(&pq).Error
		if err != nil {
			session.Set("error", err.Error())
			session.Save()
			c.Redirect(http.StatusSeeOther, "/package/view/"+id)
			return
		}
	}

	c.Redirect(http.StatusSeeOther, "/package/view/"+strconv.Itoa(int(package1.ID)))
}

func ActionDeletePackageQuestion(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var package1 model.PackageQuestion
	var package2 model.PackageQuestion
	id := c.Param("id")
	err := db.Where("id=?", id).First(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package/view/"+strconv.Itoa(int(package1.PackageID)))
		return
	}
	err = db.Where("id=?", id).Delete(&package2).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/package/view/"+strconv.Itoa(int(package1.PackageID)))
		return
	}

	c.Redirect(http.StatusSeeOther, "/package/view/"+strconv.Itoa(int(package1.PackageID)))
}
