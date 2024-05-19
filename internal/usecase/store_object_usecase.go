package usecase

import (
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"

	"github.com/sirupsen/logrus"
)

type StoreObjectUseCase struct {
	Log *logrus.Logger
}

func NewStoreObjectUseCase(log *logrus.Logger) *StoreObjectUseCase {
	return &StoreObjectUseCase{
		Log: log,
	}
}

// IsImage berfungsi untuk memeriksa apakah file adalah gambar berdasarkan tipe MIME
func (c *StoreObjectUseCase) IsImage(file *multipart.FileHeader) bool {
	switch file.Header.Get("Content-Type") {
	case "image/jpeg", "image/jpg", "image/png":
		return true
	default:
		return false
	}
}

// IsValidImageFormat berfungsi untuk memeriksa apakah file adalah gambar berdasarkan format file
/*
	- buka file;
	- baca 512 byte pertama untuk menentukan tipe gambar;
	- coba dekode gambar menggunakan decoder PNG;
	- coba dekode gambar menggunakan decoder JPEG;
	- jika tidak ada yang berhasil, format gambar tidak didukung.
*/
func (c *StoreObjectUseCase) IsValidImageFormat(file *multipart.FileHeader) bool {
	src, err := file.Open()
	if err != nil {
		c.Log.Error("Error opening file:", err)
		return false
	}
	defer src.Close()

	header := make([]byte, 512)
	if _, err := src.Read(header); err != nil {
		c.Log.Error("Error reading file header:", err)
		return false
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		c.Log.Error("Error seeking to the beginning of the file:", err)
		return false
	}
	if _, err := png.DecodeConfig(src); err == nil {
		return true
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		c.Log.Error("Error seeking to the beginning of the file:", err)
		return false
	}
	if _, err := jpeg.DecodeConfig(src); err == nil {
		return true
	}

	return false
}
