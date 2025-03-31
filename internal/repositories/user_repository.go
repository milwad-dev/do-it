package repositories

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// GetUserIdFromContext => Get userId from context
func GetUserIdFromContext(r *http.Request) interface{} {
	return r.Context().Value("userID").(jwt.MapClaims)["user_id"]
}
