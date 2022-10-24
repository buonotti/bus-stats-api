package v1

import (
	"github.com/gin-gonic/gin"
)

// RefreshUserToken godoc
// @Summary Refresh token
// @Schemes serviceV1.RefreshRequest serviceV1.RefreshResponse services.ErrorResponse
// @Description Refreshes a user token identified by the given id
// @ID refresh-token
// @Security ApiKeyAuth
// @Tags authentication
// @Accept application/json
// @Produce application/json
// @param Authorization header string true "Authorization"
// @Success 200 {object} v1.RefreshResponse
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /refresh [post]
func RefreshUserToken(c *gin.Context) {
	// Handled by library, only here for docs. Check ./root.go
}
