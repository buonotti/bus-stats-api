package v1

import "github.com/gin-gonic/gin"

// LoginUser godoc
// @Summary Logs a user in
// @Schemes services.ErrorResponse serviceV1.LoginRequest
// @Description Logs a user in using the provided credentials
// @ID login-user
// @Tags authentication
// @Accept application/json
// @Produces application/json
// @Param data body v1.LoginRequest true "content"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /login [post]
func LoginUser(c *gin.Context) {
	// Handled by library, only here for docs. Check ./root.go
}
