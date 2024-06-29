package routes

import (
	"golang-app/app/controllers"

	"github.com/gin-gonic/gin"

)

func SetupRouter(r *gin.Engine) {

	dashboardController := controllers.NewDashboardController()

	r.GET("/dashboard", dashboardController.Index)

	pendudukController := controllers.NewPendudukController()
	r.GET("/penduduk", pendudukController.Index)

	r.GET("/penduduk/data", pendudukController.DataPenduduk)

	userAdminController := controllers.NewUserAdminController()
	r.GET("/user-admin", userAdminController.Index)

	r.GET("/user-admin/data", userAdminController.GetAdminUsers)
	r.POST("/user-admin/insert", userAdminController.InsertData)
	
    r.DELETE("/user-admin/delete/:id", userAdminController.DeleteData)
    r.GET("/user-admin/getData/:id", userAdminController.GetAdminUser)
    r.PUT("/user-admin/update/:id", userAdminController.UpdateData)
}
