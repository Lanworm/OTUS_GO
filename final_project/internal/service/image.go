package service

import (
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type ImgParams struct {
	Width  uint `validate:"required,gt=0,lte=9999"`
	Height uint `validate:"required,gt=0,lte=9999"`
	URL    string
}

func ResizeImg(imgParams *ImgParams) (img image.Image, err error) {
	// Директория для сохранения изображений
	saveDir := "./images/"

	response, err := http.Get(imgParams.URL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	img, _, err = image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	newImg := resize.Resize(imgParams.Width, imgParams.Height, img, resize.Lanczos3)
	// Генерируем уникальный идентификатор для файла
	imageID := uuid.New()

	// Полный путь к файлу для сохранения в указанной директории
	filePath := filepath.Join(saveDir, imageID.String()+"_resized.jpg")

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Создаем новый файл с уникальным идентификатором для сохранения измененного изображения
	outputFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {

		}
	}(outputFile)

	// Сохраняем измененное изображение в формате JPEG
	EncodeErr := jpeg.Encode(outputFile, newImg, nil)
	if EncodeErr != nil {
		return nil, EncodeErr
	}

	return newImg, nil
}
