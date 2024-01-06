package handler

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"akuntansi/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm"
)

func IndexQuestion(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "question.html", gin.H{"Info": session.Get("info"), "Error": session.Get("error"), "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "AWS_DESTINATION_PREFIX": os.Getenv("AWS_DESTINATION_PREFIX"), "PUBLIC_URL_S3": os.Getenv("PUBLIC_URL_S3")})
}

func NewQuestion(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	data := []int{}
	loop, _ := strconv.Atoi(os.Getenv("OPTION_LOOP"))
	for i := 0; i < loop; i++ {
		data = append(data, i)
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/question/question_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"Loop": data, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Query("id")
	pv_id := c.Query("pv_id")
	session := sessions.Default(c)
	email := session.Get("email").(string)
	var question model.Question1

	err := db.Preload("Option").Preload("Tutor").Where("id=?", id).First(&question).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	var convertedOptions []map[string]interface{}

	for _, opt := range question.Option {
		optionData := map[string]interface{}{
			"ID":          opt.ID,
			"Name":        opt.Name,
			"Description": opt.Description,
			"IsTrue":      opt.IsTrue,
			"FileImage":   opt.FileImage,
			// Tambahkan bidang lain jika diperlukan
		}
		convertedOptions = append(convertedOptions, optionData)
	}

	jsonOptions, err := json.Marshal(convertedOptions)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/question/question_edit.html"))

	// Mendaftarkan fungsi escapeJS setelah template.Must
	// tmpl.Funcs(template.FuncMap{
	// 	"escapeJS": escapeJS,
	// })
	if err := tmpl.Execute(c.Writer, gin.H{"Options": string(jsonOptions), "NilToEmptyString": NilToZeroValue, "TutorName": question.Tutor.Name, "TutorID": question.Tutor.ID, "Pv_id": pv_id, "Data": question, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func Import(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/question/import.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL")}); err != nil {
		fmt.Println(err)
	}
}

func CreateQuestion(c *gin.Context, db *gorm.DB) {
	var inputQuestion model.InputQuestion
	sess := c.MustGet("sess").(*s3.S3)
	var option model.InputOption
	session := sessions.Default(c)
	// seen := make(map[int]bool)
	count_question := 4
	err := c.ShouldBind(&inputQuestion)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/question")
		return
	}
	batch := time.Now().UTC().Unix()
	name := strconv.Itoa(int(batch))
	file1, header_image, err1 := c.Request.FormFile("file_image")
	if err1 != nil {
		file1 = nil
	}
	file3, header_image3, err3 := c.Request.FormFile("file_image_ekplanation")
	if err3 != nil {
		file3 = nil
	}

	td := db.Begin()
	if td.Error != nil {
		log.Fatal(td.Error)
	}
	inputQuestion.Batch = int32(batch)
	err = td.Debug().Create(&inputQuestion).Error
	if err != nil {
		td.Rollback()
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/question")
		return
	}

	if file1 != nil {
		if strings.Contains(header_image.Filename, "(") {
			header_image.Filename = strings.ReplaceAll(header_image.Filename, "(", "")
			header_image.Filename = strings.ReplaceAll(header_image.Filename, ")", "")
		}
		UploadS3OptionQuestion(sess, file1, header_image, name, "question")
	}
	if file3 != nil {
		if strings.Contains(header_image3.Filename, "(") {
			header_image3.Filename = strings.ReplaceAll(header_image3.Filename, "(", "")
			header_image3.Filename = strings.ReplaceAll(header_image3.Filename, ")", "")
		}
		UploadS3OptionQuestion(sess, file3, header_image3, name, "question")
	}

	for i := 0; i < count_question; i++ {
		if c.PostForm("Option["+strconv.Itoa(i)+"].name") != "" {
			k, _ := strconv.Atoi(c.PostForm("Option[" + strconv.Itoa(i) + "].is_true"))
			//fmt.Println(c.PostForm("Option[" + strconv.Itoa(i) + "].description"))
			// fmt.Println(c.PostForm("Option[" + strconv.Itoa(i) + "].is_true"))
			// if seen[k] {
			// 	session.Set("error", "is_true tidak boleh duplicate")
			// 	session.Save()
			// 	c.Redirect(http.StatusFound, "/question")

			// 	return
			// }
			// seen[k] = true
			if k > count_question {
				session.Set("error", "is_true tidak boleh lebih dari "+strconv.Itoa(count_question))
				session.Save()
				c.Redirect(http.StatusFound, "/question")
			}
			option.QuestionID = inputQuestion.ID
			option.Name = c.PostForm("Option[" + strconv.Itoa(i) + "].name")
			description := c.PostForm("Option[" + strconv.Itoa(i) + "].description")
			option.Description = &description
			isTrue, _ := strconv.Atoi(c.PostForm("Option[" + strconv.Itoa(i) + "].is_true"))
			iSTrue := int32(isTrue)
			option.IsTrue = &iSTrue
			err := td.Debug().Create(&option).Error
			if err != nil {
				td.Rollback()
				fmt.Println(err.Error())
				session.Set("error", err.Error())
				session.Save()
				c.Redirect(http.StatusFound, "/question")
				return
			}
			file2, header_image2, err1 := c.Request.FormFile("Option[" + strconv.Itoa(i) + "].file")
			if err1 != nil {
				file2 = nil
			}
			if file2 != nil {
				if strings.Contains(header_image2.Filename, "(") {
					header_image2.Filename = strings.ReplaceAll(header_image2.Filename, "(", "")
					header_image2.Filename = strings.ReplaceAll(header_image2.Filename, ")", "")
				}
				UploadS3OptionQuestion(sess, file2, header_image2, name, "option")
			}

		}

	}

	if err := td.Commit().Error; err != nil {
		td.Rollback() // Mengembalikan transaksi jika terjadi kesalahan
		log.Fatal(err)
	}
	c.Redirect(http.StatusFound, "/question")
}

func UpdateQuestion(c *gin.Context, db *gorm.DB) {
	var inputQuestion model.InputQuestion
	var count_question int = 5
	sess := c.MustGet("sess").(*s3.S3)
	// seen := make(map[int]bool)
	id := c.Query("id")
	pv_id := c.Query("pv_id")
	session := sessions.Default(c)
	err := c.ShouldBind(&inputQuestion)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		if pv_id == "" {
			c.Redirect(http.StatusFound, "/question")
		} else {
			c.Redirect(http.StatusSeeOther, "/package/view/"+pv_id)
		}
		return
	}
	name := strconv.Itoa(int(inputQuestion.Batch))
	file1, header_image, err1 := c.Request.FormFile("file_image")
	if err1 != nil {
		file1 = nil
	}
	file3, header_image3, err3 := c.Request.FormFile("file_image_ekplanation")
	if err3 != nil {
		file3 = nil
	}

	err = db.Debug().Model(&inputQuestion).Where("id=?", id).Updates(&inputQuestion).Error
	if err != nil {
		fmt.Println(err.Error())
		session.Set("error", err.Error())
		session.Save()
		if pv_id == "" {
			c.Redirect(http.StatusFound, "/question")
		} else {
			c.Redirect(http.StatusSeeOther, "/package/view/"+pv_id)
		}
		return
	}

	if file1 != nil {
		if strings.Contains(header_image.Filename, "(") {
			header_image.Filename = strings.ReplaceAll(header_image.Filename, "(", "")
			header_image.Filename = strings.ReplaceAll(header_image.Filename, ")", "")
		}
		UploadS3OptionQuestionUpdate(sess, file1, header_image, id, name, "question")
	}
	if file3 != nil {
		if strings.Contains(header_image3.Filename, "(") {
			header_image3.Filename = strings.ReplaceAll(header_image3.Filename, "(", "")
			header_image3.Filename = strings.ReplaceAll(header_image3.Filename, ")", "")
		}
		UploadS3OptionQuestionUpdate(sess, file3, header_image3, id, name, "ekplanation")
	}

	// err = db.Debug().
	// 	Model(&model.Option{}).
	// 	Where("question_id = ?", id).
	// 	Find(&count_question).
	// 	Error
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	for i := 0; i < count_question; i++ {
		modelOption := model.Option{}
		criteria := ""
		if i == 0 {
			criteria = "A"
		}
		if i == 1 {
			criteria = "B"
		}
		if i == 2 {
			criteria = "C"
		}
		if i == 3 {
			criteria = "D"
		}
		if i == 4 {
			criteria = "E"
		}
		err := db.Where("name=? AND question_id=?", criteria, id).First(&modelOption).Error
		if err != nil {
			if c.PostForm("Option["+strconv.Itoa(i)+"].name") != "" {
				description := c.PostForm("Option[" + strconv.Itoa(i) + "].description")
				questionID, _ := strconv.Atoi(id)
				istrue, _ := strconv.Atoi(c.PostForm("Option[" + strconv.Itoa(i) + "].is_true"))
				istrueInt32 := int32(istrue)
				option := model.InputOption{
					Name:        criteria,
					QuestionID:  int32(questionID),
					Description: &description,
					IsTrue:      &istrueInt32,
				}
				err := db.Debug().Create(&option).Error
				if err != nil {
					fmt.Println(err)
				}
				file2, header_image2, err1 := c.Request.FormFile("Option[" + strconv.Itoa(i) + "].file")
				if err1 != nil {
					file2 = nil
				}
				if file2 != nil {
					if strings.Contains(header_image2.Filename, "(") {
						header_image2.Filename = strings.ReplaceAll(header_image2.Filename, "(", "")
						header_image2.Filename = strings.ReplaceAll(header_image2.Filename, ")", "")
					}
					UploadS3OptionQuestion(sess, file2, header_image2, name, "option")
				}
			}
		}
		if c.PostForm("Option["+strconv.Itoa(i)+"].name") != "" {
			k, _ := strconv.Atoi(c.PostForm("Option[" + strconv.Itoa(i) + "].is_true"))
			if k > count_question {
				session.Set("error", "is_true tidak boleh lebih dari "+strconv.Itoa(count_question))
				session.Save()
				if pv_id == "" {
					c.Redirect(http.StatusSeeOther, "/question")
				} else {
					c.Redirect(http.StatusSeeOther, "/package/view/"+pv_id)
				}
				return
			}
			db.Debug().Table("option").Where("name=? AND question_id=?", c.PostForm("Option["+strconv.Itoa(i)+"].name"), id).Updates(map[string]interface{}{
				"description": c.PostForm("Option[" + strconv.Itoa(i) + "].description"),
				"is_true":     c.PostForm("Option[" + strconv.Itoa(i) + "].is_true"),
			})
			file2, header_image2, err1 := c.Request.FormFile("Option[" + strconv.Itoa(i) + "].file")
			if err1 != nil {
				file2 = nil
			}
			if file2 != nil {
				if strings.Contains(header_image2.Filename, "(") {
					header_image2.Filename = strings.ReplaceAll(header_image2.Filename, "(", "")
					header_image2.Filename = strings.ReplaceAll(header_image2.Filename, ")", "")
				}
				UploadS3OptionQuestionUpdate(sess, file2, header_image2, c.PostForm("Option["+strconv.Itoa(i)+"].id"), name, "option")
			}
		} else {
			err := db.Debug().Where("name=? AND question_id=?", criteria, id).Delete(&model.Option{}).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	db.Debug().Table("question").Where("id=?", id).Updates(map[string]interface{}{
		"updated_at": time.Now().UTC().Unix(),
	})

	if pv_id == "" {
		c.Redirect(http.StatusFound, "/question")
	} else {
		c.Redirect(http.StatusSeeOther, "/package/view/"+pv_id)
	}
}

func ActionImport(c *gin.Context, db *gorm.DB) {
	var question model.InputQuestion
	var option model.Option
	session := sessions.Default(c)
	var tutorname string
	var tutor model.Tutor
	batch := time.Now().UTC().Unix()
	name := strconv.Itoa(int(batch))
	// file2, err2 := c.FormFile("f_image_o")
	// if err2 != nil {
	// 	file2 = nil
	// }
	// if file2 != nil {
	// 	uploadedFilePath := os.Getenv("PATH_FOLDER") + "/uploded1.zip"
	// 	if err := c.SaveUploadedFile(file2, uploadedFilePath); err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
	// 		return
	// 	}

	// 	extractPath := os.Getenv("PATH_FOLDER")

	// 	err := unzip(uploadedFilePath, extractPath, name+"_option")
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
	// 		return
	// 	}

	// 	// Upload folder to S3
	// 	err = uploadToS3(extractPath, os.Getenv("AWS_BUCKET"), os.Getenv("AWS_DESTINATION_PREFIX"), name+"_option", uploadedFilePath, c)
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error uploading to S3: %s", err.Error()))
	// 		return
	// 	}
	// }

	file1, err1 := c.FormFile("f_image_q")
	if err1 != nil {
		file1 = nil
	}
	if file1 != nil {
		uploadedFilePath := os.Getenv("PATH_FOLDER") + "/uploded.zip"
		if err := c.SaveUploadedFile(file1, uploadedFilePath); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		extractPath := os.Getenv("PATH_FOLDER")

		err := unzip(uploadedFilePath, extractPath, name)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		// Upload folder to S3
		err = uploadToS3(extractPath, os.Getenv("AWS_BUCKET"), os.Getenv("AWS_DESTINATION_PREFIX"), name, uploadedFilePath, c)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error uploading to S3: %s", err.Error()))
			return
		}
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
	file, err := c.FormFile("f_exel")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening uploaded file: " + err.Error()})
		return
	}
	defer uploadedFile.Close()

	// Read the Excel file from the uploaded file
	xlFile, err := xlsx.OpenReaderAt(uploadedFile, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading Excel file: " + err.Error()})
		return
	}

	// Assuming the first sheet is the one with data
	if len(xlFile.Sheets) > 0 {
		// Iterate through rows (skipping the header)
		for i, row := range xlFile.Sheets[0].Rows[1:] {
			tutor = model.Tutor{}
			// name := row.Cells[0].String()
			// status, _ := row.Cells[1].Int()
			if i == 0 && len(row.Cells) < 17 {
				session.Set("error", "format file excel tidak sesuai: jumlah column tidak sesuai template")
				session.Save()
				c.Redirect(http.StatusSeeOther, "/question")
				return
			}
			isEmptyRow := true
			rowLength := len(row.Cells)
			if rowLength >= 17 {
				rowLength = 17
			}
			for cellIndex := 0; cellIndex < rowLength; cellIndex++ {
				if row.Cells[cellIndex].String() != "" {
					isEmptyRow = false
				}
			}
			if isEmptyRow {
				continue
			}

			if i != 0 && len(row.Cells) < 17 {
				continue
			}
			if row.Cells[1].String() == "" && row.Cells[2].String() == "" {
				session.Set("error", "format file excel tidak sesuai: salah satu column harus berisi pada baris "+strconv.Itoa(i+1))
				session.Save()
				c.Redirect(http.StatusSeeOther, "/question")
				return
			}
			if row.Cells[1].String() == "" && (row.Cells[7].String() != "" || row.Cells[8].String() != "" || row.Cells[9].String() != "" || row.Cells[10].String() != "" || row.Cells[11].String() != "") {
				session.Set("error", "format file excel tidak sesuai: soal tidak ada pada baris "+strconv.Itoa(i+1))
				session.Save()
				c.Redirect(http.StatusSeeOther, "/question")
				return
			}
			if row.Cells[1].String() == "" || row.Cells[2].String() == "" {
				continue
			}
			if (row.Cells[7].String() != "" && row.Cells[12].String() == "") || (row.Cells[8].String() != "" && row.Cells[13].String() == "") || (row.Cells[9].String() != "" && row.Cells[14].String() == "") || (row.Cells[10].String() != "" && row.Cells[15].String() == "") || (row.Cells[11].String() != "" && row.Cells[16].String() == "") {
				session.Set("error", "format file excel tidak sesuai: terdapat opsi jawaban yang belum ada bobotnya pada soal baris "+strconv.Itoa(i+1))
				session.Save()
				c.Redirect(http.StatusSeeOther, "/question")
				return
			}
			questionText := row.Cells[1].String() + "\n\n" + row.Cells[2].String()
			if row.Cells[1].String() == "-" {
				questionText = row.Cells[2].String()
			}
			subtopic := row.Cells[4].String()
			// trimSection := strings.Trim(row.Cells[3].String(), " \n\t\r")
			// err := db.Debug().Where("name = ?", trimSection).First(&section).Error
			// if err != nil {
			// 	tx.Rollback()
			// 	c.HTML(http.StatusOK, "import.html", gin.H{"error": "penulisan " + row.Cells[3].String() + " kurang tepat, periksa lagi column section di file excelnya"})
			// 	return
			// }
			trimTutor := strings.Trim(row.Cells[5].String(), " \n\t\r")
			err = db.Where("name = ?", trimTutor).First(&tutor).Error
			if err != nil {
				tutorname = row.Cells[6].String()
			}

			// questionText := strings.Replace(row.Cells[1].String(), "\n", "<br/>", -1)
			// questionText = html.UnescapeString(questionText)
			ekplanation := row.Cells[3].String()
			if ekplanation == "-" {
				ekplanation = ""
			}
			// ekplanation := strings.Replace(row.Cells[2].String(), "\n", "<br/>", -1)
			// ekplanation = html.UnescapeString(ekplanation)
			if tutor.ID != 0 {
				question = model.InputQuestion{
					Question:    questionText,
					Explanation: &ekplanation,
					// SectionID:   section.ID,
					Batch:    int32(batch),
					Topic:    row.Cells[4].String(),
					Subtopic: &subtopic,
					TutorID:  &tutor.ID,
				}
			} else {
				question = model.InputQuestion{
					Question:    questionText,
					Explanation: &ekplanation,
					// SectionID:   section.ID,
					Batch:     int32(batch),
					Topic:     row.Cells[4].String(),
					Subtopic:  &subtopic,
					TutorName: &tutorname,
				}
			}
			err = tx.Create(&question).Error
			if err != nil {
				tx.Rollback()
				session.Set("error", err.Error())
				session.Save()
				c.Redirect(http.StatusSeeOther, "/question")
				return
			}
			if cell := row.Cells[8]; cell != nil && cell.String() != "" {
				IsTrue, _ := row.Cells[12].Int()
				description := row.Cells[7].String()
				option = model.Option{
					QuestionID:  question.ID,
					Name:        "A",
					IsTrue:      int32(IsTrue),
					Description: &description,
				}
				err = tx.Create(&option).Error
				if err != nil {
					tx.Rollback()
					session.Set("error", err.Error())
					session.Save()
					c.Redirect(http.StatusSeeOther, "/question")
					return
				}
			}
			if cell := row.Cells[8]; cell != nil && cell.String() != "" {
				IsTrue, _ := row.Cells[13].Int()
				description := row.Cells[8].String()
				option = model.Option{
					QuestionID:  question.ID,
					Name:        "B",
					IsTrue:      int32(IsTrue),
					Description: &description,
				}
				err = tx.Create(&option).Error
				if err != nil {
					tx.Rollback()
					session.Set("error", err.Error())
					session.Save()
					c.Redirect(http.StatusSeeOther, "/question")
					return
				}
			}
			if cell := row.Cells[9]; cell != nil && cell.String() != "" {
				description := row.Cells[9].String()
				IsTrue, _ := row.Cells[14].Int()
				option = model.Option{
					QuestionID:  question.ID,
					Name:        "C",
					IsTrue:      int32(IsTrue),
					Description: &description,
				}
				err = tx.Create(&option).Error
				if err != nil {
					tx.Rollback()
					session.Set("error", err.Error())
					session.Save()
					c.Redirect(http.StatusSeeOther, "/question")
					return
				}
			}
			if cell := row.Cells[10]; cell != nil && cell.String() != "" {
				IsTrue, _ := row.Cells[15].Int()
				description := row.Cells[10].String()
				option = model.Option{
					QuestionID:  question.ID,
					Name:        "D",
					IsTrue:      int32(IsTrue),
					Description: &description,
				}
				err = tx.Create(&option).Error
				if err != nil {
					tx.Rollback()
					session.Set("error", err.Error())
					session.Save()
					c.Redirect(http.StatusSeeOther, "/question")
					return
				}
			}
			if cell := row.Cells[11]; cell != nil && cell.String() != "" {
				description := row.Cells[11].String()
				IsTrue, _ := row.Cells[16].Int()
				option = model.Option{
					QuestionID:  question.ID,
					Name:        "E",
					IsTrue:      int32(IsTrue),
					Description: &description,
				}
				err = tx.Create(&option).Error
				if err != nil {
					tx.Rollback()
					session.Set("error", err.Error())
					session.Save()
					c.Redirect(http.StatusSeeOther, "/question")
					return
				}
			}

		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // Mengembalikan transaksi jika terjadi kesalahan
		log.Fatal(err)
	}

	// fmt.Println("Imported Users:", questions[0].Option)

	c.Redirect(http.StatusSeeOther, "/question")
}

func GetDataQuestion(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameQuestion(orderColumnIdx)

	var totalRecords int64
	var questions []model.Question

	searchQuery, queryParams := buildSearchQueryQuestion(searchValue, db)

	query := db.Model(&model.Question{}).
		Preload("Option").
		Preload("Tutor").
		Preload("PackageQuestion").
		Preload("PackageQuestion.Package").
		Where(searchQuery, queryParams...). // Menggunakan ... untuk melempar nilai parameter secara individual
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

func getColumnNameQuestion(idx int) string {
	columnsMapping := map[int]string{
		2:  "id",
		3:  "question",
		4:  "explanation",
		5:  "topic",
		6:  "tutor_name",
		7:  "is_active",
		8:  "id",
		9:  "competency_area",
		10: "domain",
		11: "science",
		12: "nursing_process",
		13: "effort",
		14: "need",
		15: "body_system",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryQuestion(searchValue string, db *gorm.DB) (string, []interface{}) {
	if searchValue != "" {
		var tutors []int
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
		db.Model(&model.Tutor{}).Where("name = ? OR name LIKE ?", searchValue, "%"+searchValue+"%").Pluck("id", &tutors)
		searchQuery := "tutor_id IN (?) OR id LIKE ? OR batch LIKE ? OR question LIKE ? OR explanation LIKE ? OR topic LIKE ? OR subtopic LIKE ? OR tutor_name LIKE ? OR competency_area LIKE ? OR domain LIKE ? OR science LIKE ? OR nursing_process LIKE ? OR effort LIKE ? OR need LIKE ? OR body_system LIKE ?"
		params := []interface{}{tutors, "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func hasPrefixWithGHJK(filePath string) bool {
	// Membagi path menjadi direktori dan nama file
	_, file := filepath.Split(filePath)

	// Mengambil nama file tanpa ekstensi
	fileName := strings.TrimSuffix(file, filepath.Ext(file))

	// Mengambil karakter pertama dari nama file
	firstChar := strings.ToLower(fileName[:1])

	// Mengecek apakah karakter pertama mengandung G, H, I, J, atau K
	return firstChar == "g" || firstChar == "h" || firstChar == "i" || firstChar == "j" || firstChar == "k"
}

func uploadToS3(sourcePath, bucketName, destinationPrefix string, folderpath string, zippath string, c *gin.Context) error {
	svc := c.MustGet("sess").(*s3.S3)

	err := filepath.Walk(sourcePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(sourcePath, filePath)
		if strings.HasPrefix(relPath, folderpath+string(filepath.Separator)) {
			key := filepath.Join(destinationPrefix, relPath)
			key = filepath.ToSlash(key)
			// if hasPrefixWithGHJK(key) {
			// 	replacement := "_option"
			// 	key = strings.Replace(key, "/"+folderpath+"/", "/"+folderpath+""+replacement+"/", 1)
			// }
			fmt.Println(key)
			// Lewati direktori
			if info.IsDir() {
				return nil
			}

			// Buka file lokal
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			// Get the file size
			fileInfo, _ := file.Stat()
			fileSize := fileInfo.Size()

			// Set content type based on file extension
			contentType := getContentType(filePath)

			// Upload the file to S3
			params := &s3.PutObjectInput{
				Bucket:        aws.String(bucketName),
				Key:           aws.String(key),
				ACL:           aws.String("public-read"),
				Body:          file,
				ContentType:   aws.String(contentType),
				ContentLength: aws.Int64(fileSize),
			}
			_, err = svc.PutObject(params)
			if err != nil {
				// handle error
				fmt.Print(err)

			}
			// _, err = sv.PutObject(&s3.PutObjectInput{
			// 	Bucket:        aws.String(bucketName),
			// 	Key:           aws.String(key),
			// 	ACL:           aws.String("public-read"),
			// 	Body:          file,
			// 	ContentType:   aws.String(contentType),
			// 	ContentLength: aws.Int64(fileSize),
			// })
			// if err != nil {
			// 	return err
			// }
		}

		return nil
	})
	if err != nil {
		return err
	}
	if err := removeFolderAndContents(folderpath); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := os.Remove(zippath); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func getContentType(filePath string) string {
	// You can implement a logic to determine content type based on file extension
	// For example: application/pdf, image/jpeg, text/plain, etc.
	// Return a default content type if the logic is not implemented.
	return "application/octet-stream"
}

func unzip(zipFilePath, extractPath string, name string) error {
	var extractedFilePath string
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Construct the full path for the extracted file, removing original folder structure
		if !strings.Contains(file.Name, ".png") {
			extractedFilePath = filepath.Join(extractPath, name, strings.Replace(filepath.Base(file.Name), strings.Split(file.Name, ".")[len(strings.Split(file.Name, "."))-1], "png", -1))
		} else {
			extractedFilePath = filepath.Join(extractPath, name, filepath.Base(file.Name))
		}

		// If it's a directory, skip it
		if file.FileInfo().IsDir() {
			continue
		}

		// Create the parent directory
		if err := os.MkdirAll(filepath.Dir(extractedFilePath), os.ModePerm); err != nil {
			return err
		}

		// Open the source file from the zip archive
		sourceFile, err := file.Open()
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// Create the target file
		targetFile, err := os.Create(extractedFilePath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		// Copy the content from source to target
		_, err = io.Copy(targetFile, sourceFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fungsi untuk menghapus folder beserta isinya
func removeFolderAndContents(folderPath string) error {
	return os.RemoveAll(folderPath)
}

func UploadS3OptionQuestion(sess *s3.S3, file multipart.File, header_image *multipart.FileHeader, name string, status string) {
	size := header_image.Size
	var path string
	_ = header_image.Filename
	buffer := make([]byte, size)
	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	if status == "question" {
		path = os.Getenv("AWS_DESTINATION_PREFIX") + "/" + name + "/" + header_image.Filename
		if !strings.Contains(header_image.Filename, ".png") {
			path = os.Getenv("AWS_DESTINATION_PREFIX") + "/" + name + "/" + strings.Replace(header_image.Filename, strings.Split(header_image.Filename, ".")[len(strings.Split(header_image.Filename, "."))-1], "png", -1)
		}
	}
	if status == "option" {
		path = os.Getenv("AWS_DESTINATION_PREFIX") + "/" + name + "/" + header_image.Filename
		if !strings.Contains(header_image.Filename, ".png") {
			path = os.Getenv("AWS_DESTINATION_PREFIX") + "/" + name + "/" + strings.Replace(header_image.Filename, strings.Split(header_image.Filename, ".")[len(strings.Split(header_image.Filename, "."))-1], "png", -1)
		}
	}
	fmt.Println(path)
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

func UploadS3OptionQuestionUpdate(sess *s3.S3, file multipart.File, header_image *multipart.FileHeader, namefile string, name string, status string) {
	size := header_image.Size
	var path string
	_ = header_image.Filename
	buffer := make([]byte, size)
	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	if status == "question" {
		path = os.Getenv("AWS_DESTINATION_PREFIX") + "/" + name + "/" + namefile + "_q.png"
	}
	if status == "option" {
		path = os.Getenv("AWS_DESTINATION_PREFIX") + "/" + name + "/" + namefile + "_o.png"
	}
	if status == "ekplanation" {
		path = os.Getenv("AWS_DESTINATION_PREFIX") + "/" + name + "/" + namefile + "_e.png"
	}
	fmt.Println(path)
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
