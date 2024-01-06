package handler

import (
	"akuntansi/model"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserResponse struct {
	Data struct {
		User []struct {
			Text string `json:"text"`
		} `json:"user"`
	} `json:"data"`
}

// func fetchData(user_id int32, token string) (string, error) {
// 	url := os.Getenv("AUTH_URL") + "/api/users?q=" + strconv.Itoa(int(user_id))
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

// 	var userResponse UserResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
// 		return "", err
// 	}

// 	if len(userResponse.Data.User) > 0 {
// 		return userResponse.Data.User[0].Text, nil
// 	}

// 	return "", nil
// }

func GetDataPresence(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	meetingID, _ := strconv.Atoi(c.Param("id"))
	orderColumn := getColumnNamePresence(orderColumnIdx)

	var totalRecords int64
	var presences []model.Presence
	var err error
	searchQuery, queryParams := buildSearchQueryPresence(searchValue)
	query := db.Debug().Model(&model.Presence{}).
		Preload("Meeting").
		Preload("UserFromID").
		Where(searchQuery, queryParams...).
		Where("meeting_id=?", meetingID).
		Count(&totalRecords).
		Limit(pageSize).Offset(page).
		Order(orderColumn + " " + orderDir).
		Find(&presences)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	for i := range presences {
		jakartaLocation, err1 := time.LoadLocation("Asia/Jakarta")
		if err1 != nil {
			// Handle error jika gagal memuat zona waktu
			c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
			return
		}
		if presences[i].Start != 0 {
			date := time.Unix(int64(presences[i].Start), 0)
			date = date.In(jakartaLocation)
			presences[i].Starte = date.Format("2006-01-02 15:04")
		}
		if presences[i].End != 0 {
			date1 := time.Unix(int64(presences[i].End), 0)
			date1 = date1.In(jakartaLocation)
			presences[i].Ende = date1.Format("2006-01-02 15:04")
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
		"data":            presences,
		"pages":           numPages,
	}

	c.JSON(http.StatusOK, response)
}

func getColumnNamePresence(idx int) string {
	columnsMapping := map[int]string{
		2: "meeting_id",
		3: "user_id",
		4: "start",
		5: "end",
		6: "rate",
		7: "comment",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryPresence(searchValue string) (string, []interface{}) {
	if searchValue != "" {
		searchQuery := "meeting_id LIKE ? OR  user_id LIKE ? OR start LIKE ? OR end LIKE ? OR rate LIKE ? OR comment LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}
