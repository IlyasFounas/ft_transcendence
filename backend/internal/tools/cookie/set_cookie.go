package cookie

import (
	"time"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func Set(c *gin.Context, path string, 
	name string, token string, ttl time.Duration, strict bool) {
	sameSite := http.SameSiteLaxMode
	if strict {
		sameSite = http.SameSiteStrictMode
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     path,
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "prod",
		SameSite: sameSite,
	})
}