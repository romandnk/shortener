package main

import (
	"github.com/romandnk/shortener/internal/app"
	"go.uber.org/fx"
)

//	@title			URL shortener project
//	@version		1.0
//	@description	Swagger API for Golang Project URL Shortener.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API [Roman] Support
//	@license.name	romandnk
//	@license.url	https://github.com/romandnk/shortener

// @BasePath	/api/v1/
func main() {
	fx.New(app.NewApp()).Run()
}
