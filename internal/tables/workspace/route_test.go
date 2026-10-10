package workspace

import (
	"testing"

	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
)

func TestWorkspaceRoutes(t *testing.T) {
	previousMode := gin.Mode()
	t.Cleanup(func() { gin.SetMode(previousMode) })
	gin.SetMode(gin.TestMode)
	router := gin.New()
	WorkspaceRoutes(router.Group("/api"), nil, security.NewJWTService("test-secret"))

	found := map[string]map[string]bool{}
	for _, route := range router.Routes() {
		if found[route.Path] == nil {
			found[route.Path] = map[string]bool{}
		}
		found[route.Path][route.Method] = true
	}

	expected := map[string][]string{
		"/api/workspace":              {"GET", "POST"},
		"/api/workspace/:workspaceID": {"GET", "PATCH", "DELETE"},
	}
	for path, methods := range expected {
		for _, method := range methods {
			if !found[path][method] {
				t.Errorf("expected %s %s route to be registered", method, path)
			}
		}
	}
}
