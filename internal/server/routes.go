package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	e.GET("/websocket", s.websocketHandler)

	var userGroup = e.Group("/api/v1/users")
	userGroup.GET("", s.ListUsers)
	userGroup.POST("", s.CreateUser)
	userGroup.GET("/:id", s.GetUserByID)
	userGroup.PUT("/:id", s.UpdateUser)
	userGroup.DELETE("/:id", s.DeleteUser)

	var libraryGroup = e.Group("/api/v1/libraries")
	libraryGroup.GET("", s.ListLibraries)
	libraryGroup.POST("", s.CreateLibrary)
	libraryGroup.GET("/:id", s.GetLibraryByID)
	libraryGroup.PUT("/:id", s.UpdateLibrary)
	libraryGroup.DELETE("/:id", s.DeleteLibrary)
	libraryGroup.GET("/:id/staffs", s.ListStaffs)
	libraryGroup.POST("/:id/staffs", s.CreateStaff)
	libraryGroup.GET("/:id/staffs/:staff_id", s.GetStaffByID)

	return e
}
