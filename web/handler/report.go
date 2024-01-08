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

func NewReport(c *gin.Context) {

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/report/report_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL")}); err != nil {
		fmt.Println(err)
	}

}

func EditReport(c *gin.Context, db *gorm.DB) {
	id := c.Query("id")
	session := sessions.Default(c)
	var package1 model.Report

	err := db.Where("id=?", id).First(&package1).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/report/report_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Data": package1, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func CreateReport(c *gin.Context, db *gorm.DB) {
	var inputTutor model.Report
	session := sessions.Default(c)
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/report")
		return
	}
	err = db.Debug().Create(&inputTutor).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/report")
		return
	}

	c.Redirect(http.StatusFound, "/report")
}

func UpdateReport(c *gin.Context, db *gorm.DB) {
	var inputTutor model.Report
	session := sessions.Default(c)
	id := c.Param("id")
	err := c.ShouldBind(&inputTutor)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/report")
		return
	}
	err = db.Debug().Model(&inputTutor).Where("id=?", id).Updates(&inputTutor).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/report")
		return
	}

	c.Redirect(http.StatusFound, "/report")
}

func IndexReport(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "report.html", gin.H{"Error": session.Get("error"), "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func GetDataReport(c *gin.Context, db *gorm.DB) {
	month, _ := strconv.Atoi(c.PostForm("month"))
	day, _ := strconv.Atoi(c.PostForm("day"))
	year, _ := strconv.Atoi(c.PostForm("year"))
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	orderColumn := getColumnReport(orderColumnIdx)

	var totalRecords int64
	var reportusers []model.Report
	searchQuery, queryParams := buildSearchQueryReport(searchValue)
	query := db.Debug().Model(&model.Report{}).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&reportusers)
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
		Find(&reportusers)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	for i := range reportusers {
		jakartaLocation, err1 := time.LoadLocation("Asia/Jakarta")
		if err1 != nil {
			// Handle error jika gagal memuat zona waktu
			c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
			return
		}
		if reportusers[i].CreatedAt != 0 {
			date := time.Unix(int64(reportusers[i].CreatedAt), 0)
			date = date.In(jakartaLocation)
			reportusers[i].CreatedAt_t = date.Format("2006-01-02 15:04")
		}
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            reportusers,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func GetDataReport1(c *gin.Context, db *gorm.DB) {
	month, _ := strconv.Atoi(c.PostForm("month"))
	day, _ := strconv.Atoi(c.PostForm("day"))
	year, _ := strconv.Atoi(c.PostForm("year"))
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	orderColumn := getColumnReport(orderColumnIdx)

	var totalRecords int64
	var reportusers []model.Report
	var reportusersbaru []model.Report
	searchQuery, queryParams := buildSearchQueryReport(searchValue)
	query := db.Debug().Model(&model.Report{}).
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&reportusers)
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
		Find(&reportusers)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	var totalk1 int
	var totalk2 int
	var totalk3 int
	var totalk4 int
	db.Table("report").Where("categori_id=?", 1).Select("SUM(price)").Scan(&totalk1)
	db.Table("report").Where("categori_id=?", 2).Select("SUM(price)").Scan(&totalk2)
	db.Table("report").Where("categori_id=?", 3).Select("SUM(price)").Scan(&totalk3)
	db.Table("report").Where("categori_id=?", 4).Select("SUM(price)").Scan(&totalk4)

	labausaha := totalk1 - totalk2
	labadiluarusaha := totalk3 - totalk4
	lababersih := labausaha - labadiluarusaha

	i1 := false
	i2 := false
	i3 := false
	i4 := false
	i5 := false
	i6 := false
	k1 := []model.Report{}
	db.Where("categori_id=?", 1).Order("id desc").Find(&k1)
	k2 := []model.Report{}
	db.Where("categori_id=?", 2).Order("id desc").Find(&k2)
	k3 := []model.Report{}
	db.Where("categori_id=?", 3).Order("id desc").Find(&k3)
	k4 := []model.Report{}
	db.Where("categori_id=?", 4).Order("id desc").Find(&k4)
	k3gab := len(k1) + len(k2) + len(k3)
	k4gab := len(k1) + len(k2) + len(k3) + len(k4)
	for i := range reportusers {
		if !i1 {
			k := model.Report{}
			err := db.Where("id=? AND categori_id=?", reportusers[i].ID, 1).First(&k).Error
			if err == nil {
				repo := model.Report{
					Name:       "Pendapatan Usaha",
					CategoriID: 6,
				}
				reportusersbaru = append(reportusersbaru, repo)
				i1 = true
			}
		}
		if i+1 == len(k1) && len(k2) == 0 && len(k1) != 0 {
			repo1 := model.Report{
				Name:       "Laba Usaha",
				CategoriID: 5,
				Price:      int32(labausaha),
			}
			reportusersbaru = append(reportusersbaru, reportusers[i])
			reportusersbaru = append(reportusersbaru, repo1)
			if len(k4) == 0 && len(k3) == 0 {
				repo1 := model.Report{
					Name:       "Laba Bersih",
					CategoriID: 5,
					Price:      int32(lababersih),
				}
				reportusersbaru = append(reportusersbaru, repo1)
			}
			continue
		}
		if !i2 {
			k := model.Report{}
			err := db.Where("id=? AND categori_id=?", reportusers[i].ID, 2).First(&k).Error
			if err == nil {
				repo := model.Report{
					Name:       "Beban Usaha",
					CategoriID: 6,
				}
				reportusersbaru = append(reportusersbaru, repo)
				i2 = true
			}
		}
		if i+1 == len(k1)+len(k2) && !i5 && i+1 != len(k1) {
			k := model.Report{}
			err := db.Where("id=? AND categori_id=?", reportusers[i].ID, 2).First(&k).Error
			if err == nil {
				repo := model.Report{
					Name:       "Jumlah Beban Usaha",
					CategoriID: 5,
					Price:      int32(-totalk2),
				}
				repo1 := model.Report{
					Name:       "Laba Usaha",
					CategoriID: 5,
					Price:      int32(labausaha),
				}
				reportusersbaru = append(reportusersbaru, reportusers[i])
				reportusersbaru = append(reportusersbaru, repo)
				reportusersbaru = append(reportusersbaru, repo1)
				if len(k4) == 0 && len(k3) == 0 {
					repo1 := model.Report{
						Name:       "Laba Bersih",
						CategoriID: 5,
						Price:      int32(lababersih),
					}
					reportusersbaru = append(reportusersbaru, repo1)
				}
				i5 = true
				continue
			}
		}
		if !i3 {
			k := model.Report{}
			err := db.Where("id=? AND categori_id=?", reportusers[i].ID, 3).First(&k).Error
			if err == nil {
				// repo := model.Report{
				// 	Name: "Jumlah beban Usaha",
				// }
				// reportusersbaru = append(reportusersbaru, repo)
				repo1 := model.Report{
					Name:       "Pendapatan di Luar Usaha",
					CategoriID: 6,
				}
				reportusersbaru = append(reportusersbaru, repo1)
				// repo2 := model.Report{
				// 	Name: "",
				// }
				// reportusersbaru = append(reportusersbaru, repo2)
				// repo3 := model.Report{
				// 	Name: "Pendapatan di Luar Usaha",
				// }
				// reportusersbaru = append(reportusersbaru, repo3)
				i3 = true
			}
		}
		if i+1 == k3gab && len(k4) == 0 && len(k3) != 0 {
			repo1 := model.Report{
				Name:       "Laba di Luar Usaha",
				CategoriID: 5,
				Price:      int32(-labadiluarusaha),
			}
			repo2 := model.Report{
				Name:       "Laba Bersih",
				CategoriID: 5,
				Price:      int32(lababersih),
			}
			reportusersbaru = append(reportusersbaru, reportusers[i])
			reportusersbaru = append(reportusersbaru, repo1)
			reportusersbaru = append(reportusersbaru, repo2)
			continue
		}
		if !i4 {
			k := model.Report{}
			err := db.Where("id=? AND categori_id=?", reportusers[i].ID, 4).First(&k).Error
			if err == nil {
				repo := model.Report{
					Name:       "Beban di Luar Usaha",
					CategoriID: 6,
				}
				reportusersbaru = append(reportusersbaru, repo)
				i4 = true
			}
		}
		if i+1 == k4gab && !i6 && i+1 != k3gab {
			k := model.Report{}
			err := db.Where("id=? AND categori_id=?", reportusers[i].ID, 4).First(&k).Error
			if err == nil {

				repo1 := model.Report{
					Name:       "Laba di Luar Usaha",
					CategoriID: 5,
					Price:      int32(-labadiluarusaha),
				}
				repo2 := model.Report{
					Name:       "Laba Bersih",
					CategoriID: 5,
					Price:      int32(lababersih),
				}
				reportusersbaru = append(reportusersbaru, reportusers[i])
				reportusersbaru = append(reportusersbaru, repo1)
				reportusersbaru = append(reportusersbaru, repo2)
				i6 = true
				continue
			}
		}
		reportusersbaru = append(reportusersbaru, reportusers[i])
	}

	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            reportusersbaru,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnReport(idx int) string {
	columnsMapping := map[int]string{
		2: "name",
		3: "categori_id",
		4: "price",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "categori_id"
	}
	return colName
}

func buildSearchQueryReport(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "name LIKE ? OR categori_id LIKE ? OR price LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func DeleteReport(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var package1 model.Report
	id := c.Query("id")

	err := db.Where("id=?", id).Delete(&package1).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/report")
		return
	}

	c.Redirect(http.StatusSeeOther, "/report")
}
