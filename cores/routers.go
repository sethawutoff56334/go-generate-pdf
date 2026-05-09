package cores

import (
	service "generatePDF/services"

	"github.com/fasthttp/router"
)

func NewRouter() *router.Router {
	r := router.New()

	r.POST("/generate-pdf", service.GeneratePDFHandler)
	return r
}
