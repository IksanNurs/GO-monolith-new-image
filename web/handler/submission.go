package handler

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"akuntansi/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type institutionResponse struct {
	Data struct {
		Institution struct {
			Name string `json:"name"`
		} `json:"institution"`
	} `json:"data"`
}

// func FetchData(institution_id int32, token string) (string, error) {
// 	url := os.Getenv("AUTH_URL") + "/api/institution/" + strconv.Itoa(int(institution_id))
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	req.Header.Add("Authorization", "Bearer "+token)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("Error in fetching data: %s", resp.Status)
// 	}

// 	var institutionResponse institutionResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&institutionResponse); err != nil {
// 		return "", err
// 	}

// 	if institutionResponse.Data.Institution.Name != "" {
// 		return institutionResponse.Data.Institution.Name, nil
// 	}

// 	return "", nil
// }

type provinceResponse struct {
	Data struct {
		Province struct {
			Name string `json:"name"`
		} `json:"province"`
	} `json:"data"`
}

// func FetchDataProvince(province_id int32, token string) (string, error) {
// 	url := os.Getenv("AUTH_URL") + "/api/province/" + strconv.Itoa(int(province_id))
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println(url)
// 	fmt.Println(token)
// 	req.Header.Add("Authorization", token)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("Error in fetching data: %s", resp.Status)
// 	}

// 	var provinceResponse provinceResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&provinceResponse); err != nil {
// 		return "", err
// 	}

// 	if provinceResponse.Data.Province.Name != "" {
// 		return provinceResponse.Data.Province.Name, nil
// 	}

// 	return "", nil
// }

type districtResponse struct {
	Data struct {
		District struct {
			Name string `json:"name"`
		} `json:"district"`
	} `json:"data"`
}

// func FetchDataDistrict(province_id int32, token string) (string, error) {
// 	url := os.Getenv("AUTH_URL") + "/api/district/" + strconv.Itoa(int(province_id))
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	req.Header.Add("Authorization", token)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("Error in fetching data: %s", resp.Status)
// 	}

// 	var districtResponse districtResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&districtResponse); err != nil {
// 		return "", err
// 	}

// 	if districtResponse.Data.District.Name != "" {
// 		return districtResponse.Data.District.Name, nil
// 	}

// 	return "", nil
// }

func IndexSubmission(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "submission.html", gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func EditSubmission(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	id := c.Param("id")
	var submission model.Submission

	err := db.
		Where("id=?", id).
		First(&submission).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/submission/submission_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Data": submission}); err != nil {
		fmt.Println(err)
	}
}

func ViewSubmission(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	// token := session.Get("token").(string)
	id := c.Param("id")
	var submission model.Submission

	err := db.
		Preload("UserFromID").
		Preload("Institution").
		Preload("Province").
		Preload("District").
		Where("user_id=?", id).
		First(&submission).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// if submission.UserID != 0 {
	// 	e, err := FetchDataUser(int(*&submission.UserID), token)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	submission.Name = &e.Data.User.Name
	// 	submission.Phone = &e.Data.User.Phone
	// 	submission.Email = &e.Data.User.Email
	// }

	// if submission.ProvinceID != nil {
	// 	fmt.Println(*submission.ProvinceID)
	// 	submission.ProvinceName, _ = FetchDataProvince(*submission.ProvinceID, token)
	// }
	// if submission.DistrictID != nil {
	// 	submission.DistrictName, _ = FetchDataDistrict(*submission.DistrictID, token)
	// }
	// if submission.InstitutionID != nil {
	// 	submission.InstitutionName, _ = FetchData(*submission.InstitutionID, token)
	// }
	c.HTML(http.StatusOK, "submission_view.html", gin.H{"NilToEmptyString": NilToZeroValue, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Data": submission})
}

func GetDataSubmission(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	token := session.Get("token").(string)
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameSubmission(orderColumnIdx)

	var totalRecords int64
	var submissions []model.Submission

	searchQuery, queryParams := buildSearchQuerySubmission(searchValue, db, token)
	query := db.Table("submission").
		Select("`submission`.*, _user.name, _user.phone, _user.email").
		Joins("LEFT JOIN _user on (submission.user_id = _user.id)").
		Where(searchQuery, queryParams...).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&submissions)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	for i := range submissions {
		submission := &submissions[i]
		if submission.UserID != 0 {
			// Menjalankan Preload hanya jika PackageID tidak 0
			if result := db.Model(submission).Association("UserFromID"); result.Error == nil {
				result.Find(&submission.UserFromID)
			}
		}
	}
	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            submissions,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNameSubmission(idx int) string {
	columnsMapping := map[int]string{
		2: "_user.name",
		3: "formation_quota",
		4: "basic_score",
		5: "basic_rank",
		6: "advanced_score",
		7: "advanced_rank",
		8: "is_public",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "user_id"
	}
	return colName
}

func buildSearchQuerySubmission(searchValue string, db *gorm.DB, token string) (string, []interface{}) {
	if searchValue != "" {
		if strings.ToLower(searchValue) == "public" {
			searchQuery := "is_public=?"
			params := []interface{}{"1"}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "private" {
			searchQuery := "is_public=?"
			params := []interface{}{"0"}
			return searchQuery, params
		}

		searchQuery := "_user.name LIKE ? OR _user.phone LIKE ? OR _user.email LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}
