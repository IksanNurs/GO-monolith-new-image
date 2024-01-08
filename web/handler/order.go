package handler

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"akuntansi/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IndexOrder(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	c.HTML(http.StatusOK, "order.html", gin.H{"Info": session.Get("info"), "Error": session.Get("error"), "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "urllogout": os.Getenv("AUTH_URL") + "/login?client_id=" + fmt.Sprintf("%s"+"%s://%s", "https", c.Request.URL.Scheme, c.Request.Host)})
}

func NewOrder(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/order/order_new.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func EditOrder(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	id := c.Param("id")
	email := session.Get("email").(string)
	var orders model.Order
	var user model.User
	query := db.Model(&model.Order{}).
		Preload("Package").
		Preload("Bundle").
		Where("id=?", id).
		First(&orders)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	err := db.Where("id=?", orders.UserID).First(&user).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	tmpl := template.Must(template.ParseFiles(os.Getenv("PATH_SUB_BASE") + "/order/order_edit.html"))

	if err := tmpl.Execute(c.Writer, gin.H{"EmailData": user.Email, "UserID": orders.UserID, "BundleID": orders.Bundle.ID, "BundleName": orders.Bundle.Name, "PackageID": orders.Package.ID, "PackageName": orders.Package.Name, "Data": orders, "email": email, "AuthURL": os.Getenv("AUTH_ADMIN_URL"), "URL": os.Getenv("AUTH_URL"), "JWTToken": session.Get("token").(string), "Info": session.Get("info"), "Error": session.Get("error")}); err != nil {
		fmt.Println(err)
	}
}

func ActionPrice(c *gin.Context, db *gorm.DB) {
	p := c.Query("package_id")
	b := c.Query("bundle_id")
	var packagePrice model.PackagePrice
	var bundlePrice model.BundlePrice

	if p != "" {
		db.Where("id=?", p).First(&packagePrice)
	}
	if b != "" {
		db.Where("id=?", b).First(&bundlePrice)
	}

	c.JSON(http.StatusOK, gin.H{"packagePrice": packagePrice.Price, "bundlePrice": bundlePrice.Price})
}

func ActionCreateOrder(c *gin.Context, db *gorm.DB) {
	a := c.Query("a")
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	userValue := c.PostForm("user_id")
	packageValue := c.PostForm("package_id")
	bundleValue := c.PostForm("bundle_id")
	priceValue := c.PostForm("price")
	emailList := c.PostForm("email_list")
	data := []string{}
	data1 := []string{}
	var order model.InputOrder

	if a != "" {
		if userValue == "" && emailList == "" {
			session.Set("info", "salah satu field harus diisi, User atau Email List")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
			return
		}
	}

	if packageValue == "" && bundleValue == "" {
		session.Set("info", "salah satu field harus diisi, package atau bundle")
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
		}

		return
	}
	if packageValue != "" && bundleValue != "" {
		session.Set("info", "hanya boldeh isi salah satu field, package atau bundle")
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
		}
		return
	}

	user_id, _ := strconv.Atoi(userValue)
	priceD, _ := strconv.Atoi(priceValue)
	priceF := int32(priceD)
	package_id, _ := strconv.Atoi(packageValue)
	packageid := int32(package_id)
	bundle_id, _ := strconv.Atoi(bundleValue)
	paidat := int32(time.Now().UTC().Unix())
	bundleid := int32(bundle_id)
	order.UserID = int32(user_id)
	order.CreatedBy = int32(userID)
	order.PaidAt = &paidat
	order.Amount = &priceF
	if packageValue != "" && bundleValue == "" {
		firstSuccess := true
		k := model.InputOrder{
			PackageID: &packageid,
			UserID:    int32(user_id),
		}
		order.PackageID = &packageid
		result := db.Where(k).FirstOrCreate(&order)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Operasi Create akan dijalankan karena data tidak ditemukan
			firstSuccess = false
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) && result.Error != nil {
			session.Set("error", result.Error)
			session.Save()
			if a == "" {
				c.Redirect(http.StatusSeeOther, "/order")
			} else {
				c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
			}
			return
		}

		if firstSuccess && result.RowsAffected != 1 {
			if (order.CancelledAt != nil || order.ExpiredAt != nil || order.CheckedOutAt == nil) && order.PaidAt == nil {
				if order.CheckedOutAt == nil && order.ID != 0 {
					err := db.Debug().Table("order").Where("id=? AND checked_out_at IS NULL", order.ID).Updates(map[string]interface{}{
						"cancelled_at": time.Now().UTC().Unix(),
						"updated_at":   time.Now().UTC().Unix(),
					}).Error
					if err != nil {
						session.Set("error", err.Error())
						session.Save()
						if a == "" {
							c.Redirect(http.StatusSeeOther, "/order")
						} else {
							c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
						}
						return
					}
				}
				order.ID = 0
				ero := int32(0)
				order.MidtransPaymentURL = nil
				order.ExpiredAt = &ero
				order.CancelledAt = &ero
				now := int32(time.Now().UTC().Unix())
				order.PaidAt = &now
				err := db.Create(&order).Error
				if err != nil {
					session.Set("error", result.Error)
					session.Save()
					if a == "" {
						c.Redirect(http.StatusSeeOther, "/order")
					} else {
						c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
					}
					return
				}
			} else {
				fmt.Println("jj")
				session.Set("error", "Data order paket gagal tersimpan karena sudah ada sebelumnya")
				session.Save()
				if a == "" {
					c.Redirect(http.StatusSeeOther, "/order")
				} else {
					c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
				}
				return
			}
		}
	} else if packageValue == "" && bundleValue != "" {
		if userValue != "" {
			firstSuccess := true
			k := model.InputOrder{
				BundleID: &bundleid,
				UserID:   int32(user_id),
			}
			order.BundleID = &bundleid
			result := db.Where(k).FirstOrCreate(&order)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Operasi Create akan dijalankan karena data tidak ditemukan
				firstSuccess = false
			}
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) && result.Error != nil {
				session.Set("error", result.Error)
				session.Save()
				if a == "" {
					c.Redirect(http.StatusSeeOther, "/order")
				} else {
					c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
				}
				return
			}

			if firstSuccess && result.RowsAffected != 1 {
				if (order.CancelledAt != nil || order.ExpiredAt != nil || order.CheckedOutAt == nil) && order.PaidAt == nil {
					if order.CheckedOutAt == nil && order.ID != 0 {
						err := db.Debug().Table("order").Where("id=? AND checked_out_at IS NULL", order.ID).Updates(map[string]interface{}{
							"cancelled_at": time.Now().UTC().Unix(),
							"updated_at":   time.Now().UTC().Unix(),
						}).Error
						if err != nil {
							session.Set("error", err.Error())
							session.Save()
							if a == "" {
								c.Redirect(http.StatusSeeOther, "/order")
							} else {
								c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
							}
							return
						}
					}
					order.ID = 0
					ero := int32(0)
					order.MidtransPaymentURL = nil
					order.ExpiredAt = &ero
					order.CancelledAt = &ero
					now := int32(time.Now().UTC().Unix())
					order.PaidAt = &now
					order.CreatedBy = int32(userID)
					err := db.Create(&order).Error
					if err != nil {
						session.Set("error", result.Error)
						session.Save()
						if a == "" {
							c.Redirect(http.StatusSeeOther, "/order")
						} else {
							c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
						}
						return
					}
				} else {
					session.Set("error", "Data order kelas gagal tersimpan karena sudah ada sebelumnya")
					session.Save()
					if a == "" {
						c.Redirect(http.StatusSeeOther, "/order")
					} else {
						c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
					}
					return
				}
			}
		}
		if emailList != "" {
			// emailValue := strings.Split(emailList, "\r\n")
			emailValue := strings.FieldsFunc(emailList, func(r rune) bool {
				return r == '\r' || r == '\n'
			})
			for _, r := range emailValue {
				// trimEmail:=strings.Trim(r, "")
				order.ID = 0
				userStruct := model.UserSelect{}
				err := db.Debug().Where("email = ?", r).First(&userStruct).Error
				if err != nil {
					data = append(data, r)
					continue
				}
				fmt.Println(userStruct.ID)
				firstSuccess := true
				k := model.InputOrder{
					BundleID: &bundleid,
				}
				order.BundleID = &bundleid
				result := db.Debug().Where(k).FirstOrCreate(&order)
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					// Operasi Create akan dijalankan karena data tidak ditemukan
					firstSuccess = false
				}
				if !errors.Is(result.Error, gorm.ErrRecordNotFound) && result.Error != nil {
					session.Set("error", result.Error)
					session.Save()
					if a == "" {
						c.Redirect(http.StatusSeeOther, "/order")
					} else {
						c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
					}
					return
				}

				if firstSuccess && result.RowsAffected != 1 {
					if (order.CancelledAt != nil || order.ExpiredAt != nil || order.CheckedOutAt == nil) && order.PaidAt == nil {
						if order.CheckedOutAt == nil && order.ID != 0 {
							err := db.Debug().Table("order").Where("id=? AND checked_out_at IS NULL", order.ID).Updates(map[string]interface{}{
								"cancelled_at": time.Now().UTC().Unix(),
								"updated_at":   time.Now().UTC().Unix(),
							}).Error
							if err != nil {
								session.Set("error", err.Error())
								session.Save()
								if a == "" {
									c.Redirect(http.StatusSeeOther, "/order")
								} else {
									c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
								}
								return
							}
						}
						order.ID = 0
						ero := int32(0)
						order.MidtransPaymentURL = nil
						order.ExpiredAt = &ero
						order.CancelledAt = &ero
						now := int32(time.Now().UTC().Unix())
						order.PaidAt = &now
						order.CreatedBy = int32(userID)
						err := db.Create(&order).Error
						if err != nil {
							session.Set("error", result.Error)
							session.Save()
							if a == "" {
								c.Redirect(http.StatusSeeOther, "/order")
							} else {
								c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
							}
							return
						}
					} else {
						data1 = append(data1, r)
					}
				}
			}
		}
	}

	if len(data) != 0 {
		formatted := strings.Join(data, ", ")
		session.Set("info", "email "+formatted+" tidak ditemukan")
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
		}
		return
	}
	if len(data1) != 0 {
		formatted := strings.Join(data1, ", ")
		session.Set("info", "data email "+formatted+" sudah ada di kelas")
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
		}
		return
	}
	if a == "" {
		c.Redirect(http.StatusSeeOther, "/order")
	} else {
		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
	}
}

func ActionUpdateOrder(c *gin.Context, db *gorm.DB) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	session := sessions.Default(c)
	a := c.Query("a")
	id := c.Query("id")
	// email := session.Get("email").(string)
	userValue := c.PostForm("user_id")
	packageValue := c.PostForm("package_id")
	bundleValue := c.PostForm("bundle_id")
	priceValue := c.PostForm("price")
	var order model.InputOrder

	if packageValue == "" && bundleValue == "" {
		session.Set("info", "salah satu field harus diisi, package atau bundle")
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
		}
		return
	}
	if packageValue != "" && bundleValue != "" {
		session.Set("info", "hanya boldeh isi salah satu field, package atau bundle")
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
		}
		return
	}

	user_id, _ := strconv.Atoi(userValue)
	priceD, _ := strconv.Atoi(priceValue)
	priceF := int32(priceD)
	package_id, _ := strconv.Atoi(packageValue)
	packageid := int32(package_id)
	bundle_id, _ := strconv.Atoi(bundleValue)
	paidat := int32(time.Now().UTC().Unix())
	bundleid := int32(bundle_id)
	order.UserID = int32(user_id)
	order.UpdatedBy = int32(userID)
	order.PaidAt = &paidat
	order.Amount = &priceF
	if packageValue != "" && bundleValue == "" {
		order.PackageID = &packageid
		order.BundleID = nil
	} else if packageValue == "" && bundleValue != "" {
		order.BundleID = &bundleid
		order.PackageID = nil
	}

	criteria := map[string]interface{}{
		"Amount":    order.Amount,
		"PackageID": order.PackageID,
		"BundleID":  order.BundleID,
		"UpdatedBy": order.UpdatedBy,
		"UserID":    order.UserID,
		"PaidAt":    order.PaidAt,
	}
	if order.UserID == 0 {
		criteria = map[string]interface{}{
			"Amount":    order.Amount,
			"PackageID": order.PackageID,
			"BundleID":  order.BundleID,
			"UpdatedBy": order.UpdatedBy,
			"PaidAt":    order.PaidAt,
		}
	}

	// err := db.Debug().Model(&order).Where("id=?", id).Updates(&order).Error
	err := db.Debug().Model(&order).Where("id=?", id).Updates(criteria).Error

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
		}
		return
	}

	if a == "" {
		c.Redirect(http.StatusSeeOther, "/order")
	} else {
		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bundleValue)
	}
}

func ActionDeleteOrder(c *gin.Context, db *gorm.DB) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := int(userData["user_id"].(float64))
	session := sessions.Default(c)
	// email := session.Get("email").(string)
	var order model.Order
	id := c.Param("id")

	err := db.Where("id=? AND created_by=?", id, userID).Delete(&order).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/order")
		return
	}

	c.Redirect(http.StatusSeeOther, "/order")
}

func GetDataOrder(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	token := session.Get("token").(string)
	page, _ := strconv.Atoi(c.PostForm("start"))
	pageSize, _ := strconv.Atoi(c.PostForm("length"))
	searchValue := c.PostForm("search[value]")
	orderColumnIdx, _ := strconv.Atoi(c.PostForm("order[0][column]"))
	orderDir := c.PostForm("order[0][dir]")

	orderColumn := getColumnNameOrder(orderColumnIdx)

	var totalRecords int64
	var orders []model.Order
	searchQuery, queryParams := buildSearchQueryOrder(searchValue, db, token)
	query := db.Model(&model.Order{}).
		Preload("Package").
		Preload("Bundle").
		Preload("UserFromID").
		Preload("UserFromID.Institution").
		Preload("UserFromCreatedBy").
		Select("`order`.*, bundle.name, package.name, _user.name, _user.phone, _user.email, _institution.name, concat(coalesce(bundle.name, ''), coalesce(package.name, '')) as bp, if (paid_at is not null, paid_at, if (expired_at is not null, expired_at, if (cancelled_at is not null, cancelled_at, if (checked_out_at is not null, checked_out_at, order.created_at)))) as order_time").
		Joins("LEFT JOIN bundle on (bundle_id = bundle.id)").
		Joins("LEFT JOIN package on (package_id = package.id)").
		Joins("LEFT JOIN _user on (order.user_id = _user.id)").
		Joins("LEFT JOIN _institution on (_user.education_institution_id = _institution.id)").
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

func getColumnNameOrder(idx int) string {
	columnsMapping := map[int]string{
		2: "order.id",
		3: "_user.name",
		4: "_institution.name",
		5: "bp",
		6: "amount",
		7: "order_time",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "order.id"
	}
	return colName
}

func getColumnNameOrder1(idx int) string {
	columnsMapping := map[int]string{
		2: "id",
		3: "_user.name",
		4: "_institution.name",
		5: "amount",
		6: "order_time",
	}

	colName, ok := columnsMapping[idx]
	if !ok {
		colName = "id"
	}
	return colName
}

func buildSearchQueryOrder(searchValue string, db *gorm.DB, token string) (string, []interface{}) {
	if searchValue != "" {
		if strings.ToLower(searchValue) == "success" {
			searchQuery := "paid_at IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "cancel" {
			searchQuery := "cancelled_at IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "cancel" {
			searchQuery := "expired_at IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}

		if strings.ToLower(searchValue) == "pending" {
			searchQuery := "checked_out_at IS NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "admin" {
			searchQuery := "midtrans_payment_url IS NULL AND created_by!=user_id"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "midtrans" {
			searchQuery := "midtrans_payment_url IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "auto" {
			searchQuery := "midtrans_payment_url IS NULL AND created_by=user_id"
			var params []interface{}
			return searchQuery, params
		}
		// e, err := FetchDataUserOption(searchValue, token)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// fmt.Println("e.Data.User")
		// fmt.Println(e.Data.User)
		// var UserID []int
		// err := db.Table("`order`").
		// 	Joins("INNER JOIN `_user` ON (user_id=_user.id)").
		// 	Where("name LIKE ? OR phone LIKE ? OR email LIKE ?", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%").
		// 	Limit("20").
		// 	Pluck("_user.id", &UserID).
		// 	Error
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// db.Model(&model.Package{}).Where("name LIKE ?", "%"+searchValue+"%").Pluck("id", &packages)
		// db.Model(&model.Bundle{}).Where("name LIKE ?", "%"+searchValue+"%").Pluck("id", &bundles)
		searchQuery := "order.id LIKE ? OR _user.name LIKE ? OR _user.phone LIKE ? OR _user.email LIKE ? OR _institution.name LIKE ? OR package.name LIKE ? OR bundle.name LIKE ? OR amount LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func buildSearchQueryOrder1(searchValue string, db *gorm.DB, bundle_id string) (string, []interface{}) {
	if searchValue != "" {
		if strings.ToLower(searchValue) == "success" {
			searchQuery := "paid_at IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "cancel" {
			searchQuery := "cancelled_at IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "cancel" {
			searchQuery := "expired_at IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}

		if strings.ToLower(searchValue) == "pending" {
			searchQuery := "checked_out_at IS NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "admin" {
			searchQuery := "midtrans_payment_url IS NULL AND created_by!=user_id"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "midtrans" {
			searchQuery := "midtrans_payment_url IS NOT NULL"
			var params []interface{}
			return searchQuery, params
		}
		if strings.ToLower(searchValue) == "auto" {
			searchQuery := "midtrans_payment_url IS NULL AND created_by=user_id"
			var params []interface{}
			return searchQuery, params
		}
		// e, err := FetchDataUserOption(searchValue, token)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// fmt.Println("e.Data.User")
		// fmt.Println(e.Data.User)
		searchQuery := "order.id LIKE ? OR _user.name LIKE ? OR _user.phone LIKE ? OR _user.email LIKE ? OR _institution.name LIKE ? OR amount LIKE ?"
		params := []interface{}{"%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%", "%" + searchValue + "%"}
		return searchQuery, params
	}
	return "", nil
}

func ClearSessionInfoAndError(c *gin.Context) {
	// Hapus sesi di sini
	session := sessions.Default(c)
	session.Delete("info")
	session.Delete("error")
	session.Delete("success")
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Session cleared"})
}

type UserResponseOption1 struct {
	Data struct {
		User []int `json:"user"`
	} `json:"data"`
}

type UserResponseOption struct {
	Data struct {
		User struct {
			Id              int    `json:"id"`
			Email           string `json:"email"`
			Name            string `json:"name"`
			Phone           string `json:"phone"`
			InstitutionName string `json:"institution_name"`
		} `json:"user"`
	} `json:"data"`
}

// func FetchDataUserOption(user_id string, token string) (UserResponseOption1, error) {
// 	var userResponse UserResponseOption1
// 	url := os.Getenv("AUTH_URL") + "/api/user-option/" + user_id
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return userResponse, err
// 	}

// 	req.Header.Add("Authorization", "Bearer "+token)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return userResponse, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return userResponse, fmt.Errorf("Error in fetching data: %s", resp.Status)
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
// 		return userResponse, err
// 	}

// 	return userResponse, nil
// }

// func FetchDataUser(user_id int, token string) (UserResponseOption, error) {
// 	var userResponse UserResponseOption
// 	url := os.Getenv("AUTH_URL") + "/api/user/" + strconv.Itoa(user_id)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return userResponse, err
// 	}

// 	fmt.Println(token)
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return userResponse, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
// 		return userResponse, fmt.Errorf("Error in fetching data: %s", resp.Status)
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
// 		return userResponse, err
// 	}

// 	return userResponse, nil
// }

func CancelOrder(c *gin.Context, db *gorm.DB) {
	a := c.Query("a")
	id := c.Query("id")
	bd_id := c.Query("bd_id")
	session := sessions.Default(c)
	// err := fetchCancel(id)
	// if err != nil {
	// 	session.Set("error", err.Error())
	// 	session.Save()
	// 	if a == "" {
	// 		c.Redirect(http.StatusSeeOther, "/order")
	// 	} else {
	// 		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
	// 	}
	// 	return
	// }
	err := db.Debug().Table("order").Where("id=? AND checked_out_at IS NULL", id).Updates(map[string]interface{}{
		"cancelled_at": time.Now().UTC().Unix(),
		"updated_at":   time.Now().UTC().Unix(),
	}).Error
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		if a == "" {
			c.Redirect(http.StatusSeeOther, "/order")
		} else {
			c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
		}
		return
	}

	if a == "" {
		c.Redirect(http.StatusSeeOther, "/order")
	} else {
		c.Redirect(http.StatusSeeOther, "/bundle/view/"+bd_id)
	}
}

// func fetchCancel(id string) error {
// 	// URL endpoint
// 	var url string

// 	if os.Getenv("MIDTRANS_ENV") != "midtrans.Production" {
// 		url = "https://api.sandbox.midtrans.com/v2/" + id + "/cancel"
// 	} else {
// 		url = "https://api.midtrans.com/v2/" + id + "/cancel"
// 	}

// 	// Membuat permintaan GET
// 	req, err := http.NewRequest("POST", url, nil)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(url)
// 	// Membuat Basic Auth string
// 	authString := os.Getenv("SERVER_KEY") + ":" + ""
// 	authEncoded := base64.StdEncoding.EncodeToString([]byte(authString))

// 	// Menambahkan header Authorization
// 	req.Header.Add("Authorization", "Basic "+authEncoded)

// 	// Melakukan permintaan HTTP
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Memeriksa status respons
// 	if resp.StatusCode != http.StatusOK {
// 		// Baca badan respons untuk melihat pesan kesalahan
// 		body, readErr := ioutil.ReadAll(resp.Body)
// 		if readErr != nil {
// 			return fmt.Errorf("Error in fetching data: %s (Cannot read response body)", resp.Status)
// 		}
// 		return fmt.Errorf("Error in fetching data: %s, Response Body: %s", resp.Status, string(body))
// 	}
// 	return nil
// }
