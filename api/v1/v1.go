package apiv1

import (
	authenticate "github.com/MyriadFlow/gateway/api/v1/authenticate"
	avatar "github.com/MyriadFlow/gateway/api/v1/avatar"
	webxr "github.com/MyriadFlow/gateway/api/v1/webxr"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes Use the given Routes
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		authenticate.ApplyRoutes(v1)
		webxr.ApplyRoutes(v1)
		avatar.ApplyRoutes(v1)
	}
}
