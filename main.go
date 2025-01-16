package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/v1/image", func(c *fiber.Ctx) error {
		// Mendapatkan konten HTML dan CSS dari permintaan
		type Request struct {
			HTML string `json:"html" form:"html"`
			CSS  string `json:"css" form:"css"`
			URL  string `json:"url" form:"url"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		// Validasi parameter
		if req.URL == "" && req.HTML == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Either 'url' or 'html' parameter is required",
			})
		}

		var img []byte
		var err error

		if req.URL != "" {
			// Jika URL diberikan, ambil screenshot dari URL
			img, err = generateImageFromURL(req.URL, req.CSS)
		} else {
			// Jika HTML diberikan, gabungkan HTML dan CSS lalu hasilkan gambar
			fullHTML := fmt.Sprintf(`
                <html>
                <head>
                    <style>%s</style>
                </head>
                <body>%s</body>
                </html>
            `, req.CSS, req.HTML)
			img, err = generateImageFromHTML(fullHTML)
		}

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate image",
			})
		}

		// Mengembalikan gambar sebagai respons
		c.Set("Content-Type", "image/png")
		c.Set("Content-Length", fmt.Sprintf("%d", len(img)))
		return c.Send(img)
	})

	log.Fatal(app.Listen(":3000"))
}

func generateImageFromHTML(htmlContent string) ([]byte, error) {
	// Membuat konteks ChromeDP
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Menyiapkan variabel untuk menyimpan data gambar
	var imageBuf []byte

	// Menjalankan tugas ChromeDP
	err := chromedp.Run(ctx,
		// Membuka halaman kosong
		chromedp.Navigate("about:blank"),
		// Menyuntikkan konten HTML dengan menjalankan JavaScript
		chromedp.Evaluate(`document.write(`+jsonString(htmlContent)+`); document.close();`, nil),
		// Menunggu hingga halaman selesai dimuat
		chromedp.WaitReady("body"),
		// Mengambil screenshot halaman
		chromedp.CaptureScreenshot(&imageBuf),
	)
	if err != nil {
		return nil, err
	}

	return imageBuf, nil
}

// generateImageFromURL menghasilkan gambar dari URL menggunakan ChromeDP
func generateImageFromURL(url, css string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var imageBuf []byte
	err := chromedp.Run(ctx,
		// Membuka URL
		chromedp.Navigate(url),
		// Menyuntikkan CSS jika ada
		chromedp.ActionFunc(func(ctx context.Context) error {
			if css != "" {
				// Membuat JavaScript untuk menyuntikkan CSS
				js := fmt.Sprintf(`
					let style = document.createElement("style");
					style.innerHTML = %s;
					document.head.appendChild(style);
				`, jsonString(css))
				expErr := chromedp.Evaluate(js, nil).Do(ctx)
				return expErr
			}
			return nil
		}),
		// Menunggu hingga halaman selesai dimuat
		chromedp.WaitReady("body"),
		// Mengambil screenshot halaman
		chromedp.CaptureScreenshot(&imageBuf),
	)
	return imageBuf, err
}

// jsonString mengonversi string ke format JSON aman untuk digunakan dalam JavaScript
func jsonString(input string) string {
	b, _ := json.Marshal(input)
	return string(b)
}
