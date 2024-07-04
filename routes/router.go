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

	orderController := controllers.NewOrderController()
	belumbayarController := controllers.NewBelumBayarController()
	lunasController := controllers.NewLunasController()
	luarnegeriController := controllers.NewLuarNegeriController()
	dalamnegeriController := controllers.NewDalamNegeriController()
	selesaiController := controllers.NewSelesaiController()
	produkrekomendasiController := controllers.NewProdukRekomendasiController()

	r.GET("/login", middleware.RedirectIfLoggedIn(), authController.LoginPage)
	r.POST("/login", authController.Login)
	r.GET("/logout", authController.Logout)

	//  API
	apiRoutes := r.Group("/api")
	{
		apiAuthController := api.NewAuthController()

		authRoutes := apiRoutes.Group("/auth")
		{
			authRoutes.POST("/register", apiAuthController.Register)
			authRoutes.POST("/login", apiAuthController.Login)
			authRoutes.POST("/logout", apiAuthController.Logout)
			authRoutes.GET("/validate", apiAuthController.ValidateToken)

			// Order

		}
		orderRoutes := apiRoutes.Group("/order")
		{
			// orderRoutes.Use(middleware.ApiAuthMiddleware())
			orderRoutes.POST("/place-order", orderController.CreateOrder)
			orderRoutes.PUT("/update-order/:id", orderController.UpdateOrder)
			orderRoutes.DELETE("/place-order/:id", orderController.DeleteOrder)
		}
	}

	// protectedRoutes := r.Group("/")
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
			orderRoutes := r.Group("/order")
			{
				orderRoutes.GET("/", orderController.Index)
				orderRoutes.GET("/data", orderController.GetOrders)
				orderRoutes.PUT("/confirm/:id", orderController.ConfirmOrder)
				orderRoutes.PUT("/decline/:id", orderController.DeclineOrder)
			}
			belumbayarRoutes := r.Group("/belum-bayar")
			{
				belumbayarRoutes.GET("/", belumbayarController.Index)
				// baruRoutes.GET("/data", orderController.GetAdminUsers)
			}
			lunasRoutes := r.Group("/lunas")
			{
				lunasRoutes.GET("/", lunasController.Index)
				// baruRoutes.GET("/data", orderController.GetAdminUsers)
			}
			luarnegeriRoutes := r.Group("/proses-luar-negeri")
			{
				luarnegeriRoutes.GET("/", luarnegeriController.Index)
				// baruRoutes.GET("/data", orderController.GetAdminUsers)
			}
			dalamnegeriRoutes := r.Group("/proses-dalam-negeri")
			{
				dalamnegeriRoutes.GET("/", dalamnegeriController.Index)
				// baruRoutes.GET("/data", baruController.GetAdminUsers)
			}
			selesaiRoutes := r.Group("/selesai")
			{
				selesaiRoutes.GET("/", selesaiController.Index)
				// baruRoutes.GET("/data", baruController.GetAdminUsers)
			}
			produkrekomendasiRoutes := r.Group("/produk-rekomendasi")
			{
				produkrekomendasiRoutes.GET("/", produkrekomendasiController.Index)
				// baruRoutes.GET("/data", baruController.GetAdminUsers)
			}

		}

		// r.Use(middleware.RoleMiddleware("admin"))
		// {
		// 	userRoutes := r.Group("/baru")
		// 	{
		// 		userRoutes.GET("/", baruController.Index)
		// 	}
		// }

	}
}
