package utils

import (
	"encoding/json"

	model "generatePDF/models"

	"github.com/valyala/fasthttp"
)

func RespondError(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	resp := model.ErrorResponse{Error: message}
	jsonBytes, _ := json.Marshal(resp)
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonBytes)
}

func RespondSuccess(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	jsonBytes, _ := json.Marshal(data)
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonBytes)
}
