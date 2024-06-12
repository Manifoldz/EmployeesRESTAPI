package handler

import (
	//_ "github.com/Manifoldz/EmployeesRESTAPI/docs"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/service"
	"github.com/gin-gonic/gin"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services}
}

// метод инициализации всех эндпоинтов
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		lists := api.Group("/employees")
		{
			lists.POST("/", h.createEmployee)
			lists.GET("/", h.getAllEmployees)
			lists.PUT("/:id", h.updateEmployeeById)
			lists.DELETE("/:id", h.deleteEmployeeById)

		}
	}
	return router
}
