package utils

import (
	"fmt"

	"github.com/gin-contrib/cors"
)

func CorsConfigStringify(cc *cors.Config) string {
	if cc == nil {
		return ""
	}
	return fmt.Sprintf(`
Cors Config:
	AllowAllOrigins: %v
	AllowOrigins:    %v
	AllowMethods:    %v
	AllowHeaders:    %v`, cc.AllowAllOrigins, cc.AllowOrigins, cc.AllowMethods, cc.AllowHeaders)

}
