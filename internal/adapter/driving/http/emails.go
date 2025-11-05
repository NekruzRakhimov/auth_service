package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmailsResponse struct {
	Emails []string `json:"emails"`
}

// GetAllUsersEmails
// @Summary Список email'ов
// @Description Список email'ов
// @Tags Internal
// @Produce json
// @Success 200 {object} EmailsResponse
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Failure 500 {object} CommonError
// @Router /internal/emails [get]
func (s *Server) GetAllUsersEmails(c *gin.Context) {
	emails, err := s.uc.EmailsGetter.GetAll(c)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, EmailsResponse{
		Emails: emails,
	})
}
