package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	model "generatePDF/models"
	"generatePDF/utils"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/valyala/fasthttp"
)

var receiptTemplate *template.Template
var cssContent string

func init() {
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"formatMoney": func(amount float64) string {
			return fmt.Sprintf("%.2f", amount)
		},
		"multiply": func(qty int, price float64) float64 {
			return float64(qty) * price
		},
		"formatPhone": func(phone string) string {
			if len(phone) == 10 {
				return phone[:3] + "-" + phone[3:6] + "-" + phone[6:]
			}
			return phone
		},
	}

	receiptTemplate = template.Must(
		template.New("receipt.html").Funcs(funcMap).ParseFiles("templates/receipt.html"),
	)

	cssBytes, err := os.ReadFile("templates/receipt.css")
	if err != nil {
		panic(fmt.Sprintf("failed to load CSS: %v", err))
	}
	cssContent = string(cssBytes)
}

type templateData struct {
	model.GeneratePDFRequest
	CSS template.CSS
}

func GeneratePDFHandler(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if len(body) == 0 {
		utils.RespondError(ctx, fasthttp.StatusBadRequest, "request body is required")
		return
	}

	var req model.GeneratePDFRequest
	if err := json.Unmarshal(body, &req); err != nil {
		utils.RespondError(ctx, fasthttp.StatusBadRequest, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}

	data := templateData{
		GeneratePDFRequest: req,
		CSS:                template.CSS(cssContent),
	}

	var htmlBuf bytes.Buffer
	if err := receiptTemplate.Execute(&htmlBuf, data); err != nil {
		utils.RespondError(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf("template render failed: %s", err.Error()))
		return
	}

	pdfBytes, err := htmlToPDF(htmlBuf.String())
	if err != nil {
		utils.RespondError(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf("PDF generation failed: %s", err.Error()))
		return
	}

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("receipt_%s.pdf", req.ReceiptNo))
	if err := os.WriteFile(tmpFile, pdfBytes, 0644); err != nil {
		utils.RespondError(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf("failed to save temp PDF: %s", err.Error()))
		return
	}

	go exec.Command("open", "-a", "Google Chrome", tmpFile).Run()

	utils.RespondSuccess(ctx, fasthttp.StatusOK, model.SuccessResponse{
		Message: "PDF opened in Chrome",
		File:    tmpFile,
	})
}

func htmlToPDF(htmlContent string) ([]byte, error) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)...,
	)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	taskCtx, cancel = context.WithTimeout(taskCtx, 30*time.Second)
	defer cancel()

	var pdfBuf []byte
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	); err != nil {
		return nil, err
	}

	return pdfBuf, nil
}
