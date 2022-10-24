package v1

import "github.com/gin-gonic/gin"

// LogoutUser godoc
// @Summary Logs a user out
// @Schemes services.ErrorResponse
// @Description Logs the current user out
// @ID logout-user
// @Tags authentication
// @Accept application/json
// @Produces application/json
// @Param data body v1.LoginRequest true "content"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /logout [post]
func LogoutUser(c *gin.Context) {
	// Handled by library, only here for docs. Check ./root.go
}
