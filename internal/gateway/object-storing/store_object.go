package object_storing

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// object store, tempat penyimpanan objek
// object storing, proses menyimpan data dalam bentuk objek yang dapat diakses melalui protokol tertentu seperti HTTP
// store object, Mengacu pada tindakan menyimpan objek ke suatu tempat
type StoreObject struct {
	Client *minio.Client
	Log    *logrus.Logger
	Config *viper.Viper
}

func NewUserObject(client *minio.Client, config *viper.Viper, log *logrus.Logger) *StoreObject {
	return &StoreObject{
		Client: client,
		Log:    log,
		Config: config,
	}
}

func (s *StoreObject) StoreFromFileHeader(ctx context.Context, uploadFile *multipart.FileHeader, objectName string) error {
	src, err := uploadFile.Open()
	if err != nil {
		s.Log.Error("Error opening file:", err)
		return err
	}
	defer src.Close()

	_, err = s.Client.PutObject(ctx, s.Config.GetString("minio.bucket"), objectName, src, uploadFile.Size, minio.PutObjectOptions{
		ContentType: uploadFile.Header.Get("Content-Type"),
	})
	if err != nil {
		s.Log.Error("Error storing object:", err)
		return err
	}

	return nil
}

// get url object
func (s *StoreObject) GetURLObject(objectName string) string {
	return s.Config.GetString("minio.objectUrl") + "/" + s.Config.GetString("minio.bucket") + "/" + objectName
}
