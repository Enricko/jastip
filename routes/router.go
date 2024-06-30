package routes

import (
	"golang-app/app/controllers"
	"golang-app/app/controllers/api"

	"github.com/gin-gonic/gin"

)

func SetupRouter(r *gin.Engine) {

	dashboardController := controllers.NewDashboardController()

	r.GET("/dashboard", dashboardController.Index)

	userAdminController := controllers.NewUserAdminController()



	userAdminRoutes := r.Group("/user-admin")
	{
		userAdminRoutes.GET("/", userAdminController.Index)
	
		userAdminRoutes.GET("/data", userAdminController.GetAdminUsers)
		userAdminRoutes.POST("/insert", userAdminController.InsertData)
		
		userAdminRoutes.DELETE("/delete/:id", userAdminController.DeleteData)
		userAdminRoutes.GET("/getData/:id", userAdminController.GetAdminUser)
		userAdminRoutes.PUT("/update/:id", userAdminController.UpdateData)
	}

	// New UserController routes
	userController := controllers.NewUserController()
	r.GET("/user", userController.Index)
	r.GET("/user/data", userController.GetUsers)
	r.POST("/user/insert", userController.InsertData)
	r.DELETE("/user/delete/:id", userController.DeleteData)
	r.GET("/user/getData/:id", userController.GetUser)
	r.PUT("/user/update/:id", userController.UpdateData)

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
