package route

import (
	"akuntansi/database"
	"net/http"
	"os"
	"path/filepath"

	"akuntansi/middleware"
	"akuntansi/storage"
	"akuntansi/web/handler"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func GetGinRoute() *gin.Engine {
	database.StartDB()
	db := database.GetDB()
	store := cookie.NewStore([]byte("secret"))
	sess := storage.ConnectAws()
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
	})
	router.Use(sessions.Sessions("mysession", store))

	router.HTMLRender = loadTemplates(os.Getenv("PATH") + "web/templates")
	router.Static("/media", os.Getenv("PATH")+"media")
	router.Static("/css", os.Getenv("PATH")+"web/assets/css")
	router.Static("/js", os.Getenv("PATH")+"web/assets/js")
	router.Static("/webfonts", os.Getenv("PATH")+"/web/assets/webfonts")

	router.GET("/login", func(c *gin.Context) {
		handler.IndexLogin(c)
	})
	router.POST("/oauth/token", func(c *gin.Context) {
		handler.OauthToken(c, db)
	})
	router.GET("/dashboard", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.Index(c, db)
	})
	router.GET("/product-user", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.IndexProductUser(c, db)
	})
	router.GET("/product-user/new", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.NewProductUser(c)
	})
	router.GET("/product-user/new/angsuran", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.NewProductUserAngsuran(c)
	})
	router.GET("/product-user/edit", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.EditProductUser(c, db)
	})

	router.GET("/user", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.IndexUser(c)
	})
	router.GET("/user/new", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.NewUser(c)
	})
	router.GET("/user/edit", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.EditUser(c, db)
	})

	router.GET("/product", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.IndexProduct(c)
	})
	router.GET("/product/new", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.NewProduct(c)
	})
	router.GET("/product/new/stock", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.NewProductProductStock(c)
	})
	router.GET("/product/edit", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.EditProduct(c, db)
	})

	router.GET("/report", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.IndexReport(c)
	})
	router.GET("/report/new", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.NewReport(c)
	})
	router.GET("/report/edit", middleware.AuthMiddleware(db), func(c *gin.Context) {
		handler.EditReport(c, db)
	})

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/login")
	})
	router.GET("/index.php", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/login")
	})
	router.GET("/logout", middleware.AuthMiddleware(db), func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
	})

	api := router.Group("api")
	{
		//
		api.GET("/clear-session-info", func(c *gin.Context) {
			handler.ClearSessionInfoAndError(c)
		})

		api.GET("/users", func(c *gin.Context) {
			handler.ActionGetAllUserCPNS(c, db)
		})
		api.GET("/users/angsuran", func(c *gin.Context) {
			handler.ActionGetAllUserCPNSAngsuran(c, db)
		})
		api.GET("/user/:id", middleware.AuthMiddleware(db), func(c *gin.Context) {
			handler.ActionGetAllUserCPNSByID(c, db)
		})

		//
		api.POST("/data-product-user", func(c *gin.Context) {
			handler.GetDataProductUser(c, db)
		})
		api.GET("/data-product-select", func(c *gin.Context) {
			handler.GetDataSelectProduct(c, db)
		})
		api.POST("/create-product-user", func(c *gin.Context) {
			handler.CreateProductUser(c, db)
		})
		api.POST("/create-product-user-angsuran", func(c *gin.Context) {
			handler.CreateProductUserAngsuran(c, db)
		})
		api.POST("/update-product-user/:id", func(c *gin.Context) {
			handler.UpdateProductUser(c, db)
		})
		api.POST("/delete-product-user", func(c *gin.Context) {
			handler.DeleteProductUser(c, db)
		})

		//
		api.POST("/data-user", func(c *gin.Context) {
			handler.GetDataUser(c, db)
		})
		api.POST("/create-user", func(c *gin.Context) {
			handler.CreateUser(c, db)
		})
		api.POST("/update-user/:id", func(c *gin.Context) {
			handler.UpdateUser(c, db)
		})
		api.POST("/delete-user", func(c *gin.Context) {
			handler.DeleteUser(c, db)
		})

		//
		api.POST("/data-product", func(c *gin.Context) {
			handler.GetDataProduct(c, db)
		})
		api.POST("/create-product", func(c *gin.Context) {
			handler.CreateProduct(c, db)
		})
		api.POST("/create-product-stock", func(c *gin.Context) {
			handler.CreateProductStock(c, db)
		})
		api.POST("/update-product/:id", func(c *gin.Context) {
			handler.UpdateProduct(c, db)
		})
		api.POST("/delete-product", func(c *gin.Context) {
			handler.DeleteProduct(c, db)
		})

		//
		api.POST("/data-report", func(c *gin.Context) {
			handler.GetDataReport(c, db)
		})
		api.POST("/data-report1", func(c *gin.Context) {
			handler.GetDataReport1(c, db)
		})
		api.POST("/create-report", func(c *gin.Context) {
			handler.CreateReport(c, db)
		})
		api.POST("/update-report/:id", func(c *gin.Context) {
			handler.UpdateReport(c, db)
		})
		api.POST("/delete-report", func(c *gin.Context) {
			handler.DeleteReport(c, db)
		})

	}

	return router
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
