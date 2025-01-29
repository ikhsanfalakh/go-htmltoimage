package utils

import (
	"encoding/json"
	"os"
)

// jsonString mengonversi string ke format JSON aman untuk digunakan dalam JavaScript
func JsonString(input string) string {
	b, _ := json.Marshal(input)
	return string(b)
}

// isTooManyRequestsError memeriksa apakah error adalah batas permintaan
/* func IsTooManyRequestsError(err error) bool {
	const requestLimit = 1000

	// Simulasi error batas permintaan
	if renderCount >= requestLimit {
		return true
	}
	return false
} */

func SaveImageToLocalStorage(filePath string, image []byte) error {
	// Buat direktori jika belum ada
	if err := os.MkdirAll("./storage/images", os.ModePerm); err != nil {
		return err
	}

	// Simpan file gambar
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(image); err != nil {
		return err
	}
	return nil
}
