package routes

import (
	"golang-app/app/controllers"
	"golang-app/app/controllers/api"
	"golang-app/app/middleware"

	"github.com/gin-gonic/gin"

)

func SetupRouter(r *gin.Engine) {
	
	dashboardController := controllers.NewDashboardController()
	userAdminController := controllers.NewUserAdminController()
	userController := controllers.NewUserController()
	authController := controllers.NewAuthController()

	r.GET("/login", middleware.RedirectIfLoggedIn(), authController.LoginPage)
	r.POST("/login", authController.Login)
	r.GET("/logout", authController.Logout)

	
	r.Use(middleware.AuthMiddleware())
	{
		r.GET("/dashboard", dashboardController.Index)

		// ADMIN ROLE ROUTE
		r.Use(middleware.RoleMiddleware("admin"))
		{
			userAdminRoutes := r.Group("/user-admin")
			{
				userAdminRoutes.GET("/", userAdminController.Index)
		
				userAdminRoutes.GET("/data", userAdminController.GetAdminUsers)
				userAdminRoutes.POST("/insert", userAdminController.InsertData)
		
				userAdminRoutes.DELETE("/delete/:id", userAdminController.DeleteData)
				userAdminRoutes.GET("/getData/:id", userAdminController.GetAdminUser)
				userAdminRoutes.PUT("/update/:id", userAdminController.UpdateData)
			}
			userRoutes := r.Group("/user")
			{
				userRoutes.GET("/", userController.Index)
				userRoutes.GET("/data", userController.GetUsers)
				userRoutes.POST("/insert", userController.InsertData)
				userRoutes.DELETE("/delete/:id", userController.DeleteData)
				userRoutes.GET("/getData/:id", userController.GetUser)
				userRoutes.PUT("/update/:id", userController.UpdateData)
			}
		}

	}

	// Login routes
	

	//  API
	apiRoutes := r.Group("/api")
	{
		authController := api.NewAuthController()

		authRoutes := apiRoutes.Group("/auth")
		{
			authRoutes.POST("/register", authController.Register)
			authRoutes.POST("/login", authController.Login)
			authRoutes.POST("/logout", authController.Logout)
			authRoutes.GET("/validate", authController.ValidateToken)
		}
	}
}
