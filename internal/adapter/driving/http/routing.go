package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/NekruzRakhimov/auth_service/api/docs"
)

func (s *Server) endpoints() {
	s.router.GET("/ping", s.Ping)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	internalG := s.router.Group("/internal")
	internalG.GET("/emails", s.GetAllUsersEmails)

	authG := s.router.Group("/auth")
	authG.POST("/sign-up", s.SignUp)
	authG.POST("/sign-in", s.SignIn)
	authG.GET("/refresh", s.RefreshTokenPair)
}

// Ping
// @Summary Health-check
// @Description Проверка сервиса
// @Tags Ping
// @Produce json
// @Success 200 {object} CommonResponse
// @Router /ping [get]
func (s *Server) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}
