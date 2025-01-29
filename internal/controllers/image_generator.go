package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	config "go-htmlcsstoimage/configs"
	model "go-htmlcsstoimage/internal/models"
	util "go-htmlcsstoimage/pkg/utils"
)

func GenerateImage(c *fiber.Ctx) error {
	// Mendapatkan konten HTML dan CSS dari permintaan

	var req model.GenerateImageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":      "Bad Request",
			"statusCode": 400,
			"message":    "Invalid request payload",
		})
	}

	// Validasi parameter
	if req.URL == "" && req.HTML == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":      "Bad Request",
			"statusCode": 400,
			"message":    "HTML or URL is required",
		})
	}

	var img []byte
	var err error

	if req.URL != "" {
		img, err = generateImageFromURL(req)
	} else {
		img, err = generateImageFromHTML(req)
	}

	if err != nil {
		// if util.IsTooManyRequestsError(nil) {
		// 	return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
		// 		"error":      "Plan limit exceeded",
		// 		"statusCode": 429,
		// 		"message":    fmt.Sprintf("You've used %d of your 1000 renders. Upgrade your plan to continue.", renderCount),
		// 	})
		// }
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"statusCode": 500,
			"message":    err.Error(),
		})
	}

	// Simpan gambar ke storage
	imageID := uuid.New().String()
	imageFileName := fmt.Sprintf("%s.png", imageID)
	imagePath := fmt.Sprintf("./storage/images/%s", imageFileName)

	if err := util.SaveImageToLocalStorage(imagePath, img); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"statusCode": 500,
			"message":    "Failed to save image to storage",
		})
	}

	// Simpan informasi gambar ke database
	getUserID, _ := strconv.Atoi(c.Locals("UserID").(string))
	imageRecord := model.Image{
		ImageName: imageID,
		UserID:    getUserID,
		CreatedAt: time.Now(),
		Type:      "png",
	}

	if err := model.SaveImageToDatabase(imageRecord); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"statusCode": 500,
			"message":    "Failed to save image information to database",
		})
	}

	// Kirimkan respons ke klien
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Image generated successfully",
		"url":     config.AppEnv.AppURL + "/storage/images/" + imageRecord.ImageName + "." + imageRecord.Type,
	})
}

func generateImageFromHTML(req model.GenerateImageRequest) ([]byte, error) {
	// Membuat konteks ChromeDP
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Menyiapkan variabel untuk menyimpan data gambar
	var imageBuf []byte
	fullHTML := fmt.Sprintf(`
		<html>
		<head>
			<link href="https://fonts.googleapis.com/css2?family=%s&display=swap" rel="stylesheet">
			<style>%s</style>
		</head>
		<body>%s</body>
		</html>
	`, req.GoogleFonts, req.CSS, req.HTML)

	// Menjalankan tugas ChromeDP
	err := chromedp.Run(ctx,
		// Membuka halaman kosong
		chromedp.Navigate("about:blank"),
		// Menyuntikkan konten HTML
		chromedp.Evaluate(`document.write(`+util.JsonString(fullHTML)+`); document.close();`, nil),
		// Menunggu hingga halaman selesai dimuat
		chromedp.WaitReady("body"),
		// Menambahkan delay jika diminta
		chromedp.Sleep(time.Duration(req.MsDelay)*time.Millisecond),
		// Menangkap tangkapan layar dari elemen tertentu atau seluruh halaman
		func() chromedp.Action {
			if req.Selector != "" {
				return chromedp.Screenshot(req.Selector, &imageBuf, chromedp.NodeVisible)
			}
			return chromedp.CaptureScreenshot(&imageBuf)
		}(),
	)
	return imageBuf, err
}

// generateImageFromURL menghasilkan gambar dari URL menggunakan ChromeDP
func generateImageFromURL(req model.GenerateImageRequest) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var imageBuf []byte
	err := chromedp.Run(ctx,
		// Mengatur viewport jika diminta
		chromedp.EmulateViewport(req.ViewportWidth, req.ViewportHeight, chromedp.EmulateScale(req.DeviceScale)),
		// Membuka URL
		chromedp.Navigate(req.URL),
		// Menyuntikkan CSS jika ada
		chromedp.ActionFunc(func(ctx context.Context) error {
			if req.CSS != "" {
				// Membuat JavaScript untuk menyuntikkan CSS
				js := fmt.Sprintf(`
					let style = document.createElement("style");
					style.innerHTML = %s;
					document.head.appendChild(style);
				`, util.JsonString(req.CSS))
				expErr := chromedp.Evaluate(js, nil).Do(ctx)
				return expErr
			}
			return nil
		}),
		// Menunggu hingga halaman selesai dimuat
		chromedp.WaitReady("body"),
		// Menambahkan delay jika diminta
		chromedp.Sleep(time.Duration(req.MsDelay)*time.Millisecond),
		// Menangkap tangkapan layar
		func() chromedp.Action {
			if req.Selector != "" {
				return chromedp.Screenshot(req.Selector, &imageBuf, chromedp.NodeVisible)
			}
			if req.FullScreen {
				return chromedp.FullScreenshot(&imageBuf, 100)
			}
			return chromedp.CaptureScreenshot(&imageBuf)
		}(),
	)
	return imageBuf, err
}

func GetImageOld(c *fiber.Ctx) error {
	// Ambil nama file dari URL
	imageName := c.Params("filename") // Mengambil sisa path, termasuk subdirektori jika ada

	// Tentukan path ke direktori penyimpanan gambar
	storagePath := "./storage/images"

	// Gabungkan direktori dengan nama file
	imagePath := filepath.Join(storagePath, imageName)

	// Periksa apakah file gambar ada
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	// Kirim file gambar sebagai respons
	return c.SendFile(imagePath)
}
