package controllers

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/nfnt/resize"
)

// GetImage handles image retrieval, format conversion, resizing, and download
func GetImage(c *fiber.Ctx) error {
	// Ambil parameter nama file
	imageName := c.Params("filename")

	// Default format adalah PNG
	imageFormat := "png"
	if strings.Contains(imageName, ".") {
		parts := strings.Split(imageName, ".")
		imageName = parts[0]   // Ambil ID tanpa ekstensi
		imageFormat = parts[1] // Ambil format yang diminta
	}

	// Validasi format yang didukung
	allowedFormats := map[string]bool{"png": true, "jpg": true, "webp": true}
	if !allowedFormats[imageFormat] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image format. Supported formats: png, jpg, webp",
		})
	}

	// Path file original (format PNG)
	imagePathOri := filepath.Join("./storage/images/", fmt.Sprintf("%s.png", imageName))

	// Periksa apakah file ada
	if _, err := os.Stat(imagePathOri); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	// Buka file PNG asli
	file, err := os.Open(imagePathOri)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open image",
		})
	}
	defer file.Close()

	// Decode gambar PNG
	img, err := png.Decode(file)
	if err != nil {
		log.Println("Decoding error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode PNG image",
		})
	}

	// Ambil parameter resize
	width, _ := strconv.Atoi(c.Query("width", "0"))
	height, _ := strconv.Atoi(c.Query("height", "0"))

	// Maksimal ukuran 5000x5000
	if width > 5000 {
		width = 5000
	}
	if height > 5000 {
		height = 5000
	}

	// Resize gambar jika parameter diberikan
	if width > 0 || height > 0 {
		img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	}

	// Buat file sementara untuk menyimpan gambar yang sudah dikonversi
	tempFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.%s", imageName, imageFormat))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create temporary file",
		})
	}
	defer os.Remove(tempFile.Name()) // Hapus file sementara setelah selesai

	// Konversi PNG ke format yang diminta
	switch imageFormat {
	case "png":
		err = png.Encode(tempFile, img)
	case "jpg":
		rgba := image.NewRGBA(img.Bounds())
		white := color.White
		draw.Draw(rgba, rgba.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Over)

		c.Set("Content-Type", "image/jpeg")
		return jpeg.Encode(c.Response().BodyWriter(), rgba, &jpeg.Options{Quality: 90})
	case "webp":
		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 90)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create WebP encoder options",
			})
		}

		c.Set("Content-Type", "image/webp")
		return webp.Encode(c.Response().BodyWriter(), img, options)
	}

	// Tangani error encoding
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to encode image",
		})
	}

	// Tutup file sementara setelah selesai menulis
	tempFile.Close()

	// Jika `dl=1`, atur header untuk download
	if c.Query("dl") == "1" {
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", imageName, imageFormat))
	}

	// Set Content-Type sesuai format yang diminta
	c.Set("Content-Type", fmt.Sprintf("image/%s", imageFormat))

	// Kirim file sebagai response
	return c.SendFile(tempFile.Name())
}
