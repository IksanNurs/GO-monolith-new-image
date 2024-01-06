package handler

import (
	"akuntansi/module/institution"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type institutionHandler struct {
	institutionService institution.Service
}

func NewIntitutionHandler(institutionService institution.Service) *institutionHandler {
	return &institutionHandler{institutionService}
}

func (h *institutionHandler) Index(c *gin.Context) {

	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "institution.html", gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL")})
}

func (h *institutionHandler) NewInstitution(c *gin.Context) {
	// provinces, err := h.institutionService.FetchAllProvince()
	// if err != nil {
	// 	c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
	// 	return
	// }

	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "institution_new.html", gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL")})
}

func (h *institutionHandler) CreateInstitution(c *gin.Context) {
	var input institution.FormCreateInstitution

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "institution_new.html", input)
		return
	}

	_, err = h.institutionService.Create(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err})
		return
	}
	c.Redirect(http.StatusFound, "/admin/institution")
}

func (h *institutionHandler) EditInstitution(c *gin.Context) {
	userId := c.Param("id")
	id, _ := strconv.Atoi(userId)

	oldInstitution, err := h.institutionService.GetInstitutionByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "institution_edit.html", gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "old": oldInstitution})
}

func (h *institutionHandler) UpdateInstitution(c *gin.Context) {
	institutionId := c.Param("id")
	id, _ := strconv.Atoi(institutionId)

	var input institution.FormCreateInstitution
	err := c.ShouldBind(&input)
	if err != nil {
		fmt.Print(err)
	}

	_, err = h.institutionService.Update(input, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/admin/institution")
}

func (h *institutionHandler) DeleteInstitution(c *gin.Context) {
	institutionID := c.Param("id")
	id, _ := strconv.Atoi(institutionID)

	_, err := h.institutionService.Delete(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/institution")
}

// func (h *contentHandler) DeleteContent(c *gin.Context) {
// 	contentId := c.Param("id")
// 	id, _ := strconv.Atoi(contentId)

// 	_, err := h.contentService.DeleteContent(id)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	c.Redirect(http.StatusFound, "/admin/contents")
// }

func (h *institutionHandler) GetData(c *gin.Context, db *sql.DB) {
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")
	fmt.Println(orderDir)

	// Mapping kolom pada DataTable ke kolom pada query
	columnsMapping := map[int]string{
		3: "name",
		4: "code",
		5: "type",
	}

	orderColumn, ok := columnsMapping[orderColumnIdx]
	if !ok {
		orderColumn = "id" // Default kolom pengurutan jika tidak ada yang cocok
	}

	// Konstruksi kueri SQL untuk mencari data dengan server-side paging dan pengurutan
	query := "SELECT id, name, code, type, province_name FROM institution WHERE 1=1"
	if searchValue != "" {
		query += " AND (id LIKE '%" + searchValue + "%' OR name LIKE '%" + searchValue + "%' OR type LIKE '%" + searchValue + "%' OR province_name LIKE '%" + searchValue + "%' OR code LIKE '%" + searchValue + "%')"
	}

	countQuery := strings.Replace(query, "SELECT id, name, code, type, province_name", "SELECT COUNT(id)", 1)

	// Eksekusi kueri untuk menghitung jumlah total data tanpa paging
	var totalRecords int
	err := db.QueryRow(countQuery).Scan(&totalRecords)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing query"})
		return
	}

	// Eksekusi kueri untuk mendapatkan data dengan paging dan pengurutan
	query += " ORDER BY " + orderColumn + " " + orderDir + " LIMIT ? OFFSET ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing query"})
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(pageSize, page)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing query"})
		return
	}
	defer rows.Close()

	data := []institution.InstitutionDatatable{}
	for rows.Next() {
		var d institution.InstitutionDatatable
		err := rows.Scan(&d.ID, &d.Name, &d.Code, &d.Type, &d.Province_name)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
			return
		}
		data = append(data, d)
	}

	// Hitung jumlah halaman berdasarkan total data dan ukuran halaman
	numPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	// Format data dalam bentuk JSON untuk dikirim ke DataTable
	response := map[string]interface{}{
		"draw":            c.PostForm("draw"),
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords, // Total data yang cocok dengan kriteria pencarian (tanpa paging)
		"data":            data,
		"pages":           numPages, // Jumlah halaman yang ada
	}

	c.JSON(http.StatusOK, response)
}
