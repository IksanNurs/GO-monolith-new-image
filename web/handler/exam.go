package handler

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"akuntansi/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NilToZeroValue(ptr interface{}) interface{} {
	val := reflect.ValueOf(ptr)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return ""
	}
	return ptr
}

func IndexExam(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var totalRecords int64
	query := db.Model(&model.Exam{}).
		Where("score_summarized IS NULL").
		Count(&totalRecords)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	c.HTML(http.StatusOK, "exam.html", gin.H{"Info": session.Get("info"), "Error": session.Get("error"), "totalRecords": totalRecords, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func EditExam(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	id := c.Param("id")
	var exam model.Exam

	err := db.
		Preload("Package").
		Where("id=?", id).
		First(&exam).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/exam/exam_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Data": exam}); err != nil {
		fmt.Println(err)
	}
}

func ViewExam(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	// token := session.Get("token").(string)
	id := c.Param("id")
	var exam model.Exam
	var examquestion model.ExamQuestion

	err := db.
		Preload("Package").
		Preload("UserFromCreatedBy").
		Where("id=?", id).
		First(&exam).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// if exam.CreatedBy != nil {
	// 	e, err := FetchDataUser(int(*exam.CreatedBy), token)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	exam.Name = &e.Data.User.Name
	// 	exam.Phone = &e.Data.User.Phone
	// 	exam.Email = &e.Data.User.Email
	// }
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle error jika gagal memuat zona waktu
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	date := time.Unix(int64(exam.StartedAt), 0)
	date = date.In(jakartaLocation)
	exam.StartedAtTe = date.Format("2006-01-02 15:04")

	if exam.EndedAt != nil {
		date1 := time.Unix(int64(*exam.EndedAt), 0)
		date1 = date1.In(jakartaLocation)
		exam.EndedAtTe = date1.Format("2006-01-02 15:04")
	}

	db.Where("exam_id=?", id).Order("updated_at desc").First(&examquestion)
	elapsedTime := int32(0)
	if examquestion.UpdatedAt != nil {
		elapsedTime = *examquestion.UpdatedAt - exam.StartedAt
	}
	c.HTML(http.StatusOK, "exam_view.html", gin.H{"elapsedTime": elapsedTime, "NilToEmptyString": NilToZeroValue, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Data": exam, "AWS_DESTINATION_PREFIX": os.Getenv("AWS_DESTINATION_PREFIX"), "PUBLIC_URL_S3": os.Getenv("PUBLIC_URL_S3")})
}

func ActionEditExtraTime(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	e_time := c.PostForm("e_time")
	eTime, _ := strconv.Atoi(e_time)
	result := db.Debug().Table("exam").Where("id=?", id).Updates(map[string]interface{}{
		"extra_time": eTime,
		"updated_at": time.Now().UTC().Unix(),
		"updated_by": userID,
	})
	if result.Error != nil {
		session.Set("error", result.Error.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/exam")
		return
	}

	c.Redirect(http.StatusSeeOther, "/exam")
}

func GetDataExam(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	token := session.Get("token").(string)
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameExam(orderColumnIdx)

	var totalRecords int64
	var exams []model.Exam

	searchQuery, queryParams := buildSearchQueryExam(searchValue, db, token)
	query := db.Table("exam").
		Select("exam.*, _user.name, _user.phone, _user.email, _institution.name, if (package.scoring_method = 1, score_percentaged, score_summarized) as score").
		Joins("left join package on (package_id = package.id)").
		Joins("LEFT JOIN _user on (exam.created_by = _user.id)").
		Joins("LEFT JOIN _institution on (_user.education_institution_id = _institution.id)").
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&exams)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle error jika gagal memuat zona waktu
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range exams {
		exam := &exams[i]
		if exam.PackageID != 0 {
			if result := db.Model(exam).Association("Package"); result.Error == nil {
				result.Find(&exam.Package)
			}
		}

		if exam.CreatedBy != nil {
			if result := db.Model(exam).Association("UserFromCreatedBy"); result.Error == nil {
				result.Find(&exam.UserFromCreatedBy)
			}
		}

		if exam.UserFromCreatedBy.ID != nil {
			if result := db.Model(&exam.UserFromCreatedBy).Association("Institution"); result.Error == nil {
				result.Find(&exam.UserFromCreatedBy.Institution)
			}
		}
		// if exams[i].CreatedBy != nil {
		// 	e, err := FetchDataUser(int(*exams[i].CreatedBy), token)
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	exams[i].Name = &e.Data.User.Name
		// 	exams[i].Email = &e.Data.User.Email
		// 	exams[i].Phone = &e.Data.User.Phone
		// }
		date := time.Unix(int64(exams[i].StartedAt), 0)
		date = date.In(jakartaLocation)
		exams[i].StartedAtTe = date.Format("2006-01-02 15:04")
		if exams[i].EndedAt != nil {
			date1 := time.Unix(int64(*exams[i].EndedAt), 0)
			date1 = date1.In(jakartaLocation)
			exams[i].EndedAtTe = date1.Format("2006-01-02 15:04")
		}
	}
	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            exams,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameExam(idx int) string {
	columnsMapping := map[int]string{
		2: "exam.id",
		3: "_user.name",
		4: "_institution.name",
		5: "package.name",
		6: "extra_time",
		7: "score",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "exam.id"
	}
	return colName
}

func buildSearchQueryExam(searchValue string, db *gorm.DB, token string) (string, []interface{}) {
	if searchValue != "" {
		// var sectionIDs []int
		// var UserID []int
		// err := db.Table("_user").
		// 	Joins("INNER JOIN `order` ON (_user.id=order.user_id)").
		// 	Joins("INNER JOIN `exam` ON (order.id=exam.order_id)").
		// 	Where("name LIKE ? OR phone LIKE ? OR email LIKE ?", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%").
		// 	Limit("20").
		// 	Pluck("_user.id", &UserID).
		// 	Error
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// db.Model(&model.Package{}).Where("name = ? OR name LIKE ?", searchValue, "%"+searchValue+"%").Pluck("id", &sectionIDs)
		//
		searchQuery := "exam.id LIKE ? OR _user.name LIKE ? OR _user.phone LIKE ? OR _user.email LIKE ? OR _institution.name LIKE ? OR package.name LIKE ? OR order_id LIKE ? OR extra_time LIKE ? OR score_percentaged LIKE ? OR score_summarized LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func ActionDeleteExam(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var exam model.Exam
	id := c.Param("id")

	err := db.Where("id=?", id).Delete(&exam).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/exam")
		return
	}

	c.Redirect(http.StatusSeeOther, "/exam")
}

func RecalculateScore(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	var ekam []model.SubExam
	id := c.Param("id")
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))

	err := db.Debug().Where("package_id=?", id).Find(&ekam).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/package/view/"+id)
		return
	}

	if len(ekam) != 0 {
		var count_question int64 = 0
		var count_questionistrue int64 = 0
		score := model.ExamQuestionScore{}
		for _, e := range ekam {
			err = db.Debug().
				Table("exam_question").
				Select("COALESCE(SUM(option.is_true), 0) as sumscore, COALESCE(COUNT(exam_question.id), 0) as count_question").
				Joins("LEFT JOIN option ON exam_question.option_id=option.id").
				Where("exam_id=? AND option_id IS NOT NULL", e.ID).
				First(&score).Error

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			err = db.Debug().
				Model(&model.ExamQuestion{}).
				Select("COALESCE(SUM(option.is_true), 0) as sumscore, COALESCE(COUNT(exam_question.id), 0) as count_question").
				Joins("LEFT JOIN option ON exam_question.option_id=option.id").
				Where("exam_id=?", e.ID).
				Count(&count_question).
				Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			err = db.Debug().
				Model(&model.ExamQuestion{}).
				Select("COALESCE(SUM(option.is_true), 0) as sumscore, COALESCE(COUNT(exam_question.id), 0) as count_question").
				Joins("LEFT JOIN option ON exam_question.option_id=option.id").
				Where("exam_id=?", e.ID).
				Count(&count_questionistrue).
				Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			percentaged := (float64(count_questionistrue) / float64(count_question)) * 100
			percentaged, err = strconv.ParseFloat(fmt.Sprintf("%.2f", percentaged), 64)
			if err != nil {
				fmt.Println(err) // 3.14159265
			}
			result := db.Debug().Table("exam").Where("id=?", e.ID).Updates(map[string]interface{}{
				"score_summarized":  score.SumScore,
				"score_percentaged": percentaged,
				"updated_at":        time.Now().UTC().Unix(),
				"updated_by":        userID,
			})
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
				return
			}

		}
		session.Set("success", "Berhasil recalculated "+strconv.Itoa(len(ekam))+" exam")
		session.Save()
		c.Redirect(http.StatusFound, "/package/view/"+id)
		return
	}
	session.Set("error", "Tidak ada exam di package ini")
	session.Save()
	c.Redirect(http.StatusFound, "/package/view/"+id)
}

func GetDataExamQuestion(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	token := session.Get("token").(string)
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	id := c.Param("id")

	orderColumn := getColumnNameExamQuestion(orderColumnIdx)

	var totalRecords int64
	var exams []model.ExamQuestionData

	searchQuery, queryParams := buildSearchQueryExamQuestion(searchValue, db, token)
	query := db.Table("exam_question").
		Where(searchQuery, queryParams...).
		Where("exam_id=?", id).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&exams)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle error jika gagal memuat zona waktu
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range exams {
		exam := &exams[i]
		if exam.QuestionID != 0 {
			if result := db.Model(exam).Association("Question"); result.Error == nil {
				result.Find(&exam.Question)
			}
		}
		if exam.OptionID != 0 {
			if result := db.Model(exam).Association("Option"); result.Error == nil {
				result.Find(&exam.Option)
			}
		}

		if exams[i].UpdatedAt != nil {
			date := time.Unix(int64(*exams[i].UpdatedAt), 0)
			date = date.In(jakartaLocation)
			exams[i].UpdatedAtTe = date.Format("2006-01-02 15:04")
		}
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            exams,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameExamQuestion(idx int) string {
	columnsMapping := map[int]string{
		0: "sequence",
		1: "question_id",
		2: "option_id",
		3: "updated_at",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "exam_question.id"
	}
	return colName
}

func buildSearchQueryExamQuestion(searchValue string, db *gorm.DB, token string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "sequence LIKE ? OR question_id LIKE ? OR option LIKE ? OR updated_at"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}
