package file

import (
	"bytes"
	"context"
	"fmt"
	"gustav-tgbot/internal/model"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	directoryStorage = "data/storage/images"
)

type UserRepository interface {
	GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error)
}

type FileRepository interface {
	CreateFile(ctx context.Context, file *model.File) (int64, error)
	GetFilesByUserID(ctx context.Context, userID int64) ([]model.File, error)
}

type File struct {
	FileRepository FileRepository
	UserRepository UserRepository

	log *slog.Logger
}

func (f *File) GetFilesByUserID(ctx context.Context, userID int64) (map[int64][]byte, error) {

	files, err := f.FileRepository.GetFilesByUserID(ctx, userID)
	if err != nil {
		f.log.Error("failed to get files ", "user_id", userID, "error", err)
		return nil, err
	}
	m := make(map[int64][]byte)
	for _, file := range files {

		photo, err := os.ReadFile(*file.ShortUrlUpgrade)
		if err != nil {
			f.log.Error("failed to read image", "user_id", userID, "file url", *file.ShortUrlUpgrade, "error", err)
			continue
		}

		m[file.CreatedAt] = photo
		f.log.Info("file info", "user_id", userID, "file id", file.ID, "file url", file.ShortUrlUpgrade)
	}
	return m, nil
}

func (f *File) StoreFile(ctx context.Context, filePath, fileURL string, userID int64) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fileURL)
	if err != nil {
		f.log.Error("failed to download file", "user_id", userID, "fileURL", fileURL, "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download file: http status %d", resp.StatusCode)
	}

	// Формируем имя файла.
	ext := filepath.Ext(filePath)
	filename := fmt.Sprintf(
		"%d_%d_%s",
		time.Now().UnixNano(),
		userID,
		filepath.Base(filePath),
	)

	file := model.NewFile(filename, ext, userID)

	// Путь сохранения оригинального изображения.
	path := filepath.Join(directoryStorage, filename)

	// Создаём файл.
	dst, err := os.Create(path)
	if err != nil {
		f.log.Error("failed to create file", "user_id", userID, "path", path, "error", err)
		return err
	}

	_, copyErr := io.Copy(dst, resp.Body)
	closeErr := dst.Close()

	if copyErr != nil {
		f.log.Error("failed to save original file", "user_id", userID, "path", path, "error", copyErr)
		return copyErr
	}

	if closeErr != nil {
		return closeErr
	}

	f.log.Info("original file stored", "user_id", userID, "path", path)

	// Сохраняем путь в модели.
	file.SetShortUrlOriginal(&path)

	// Отправляем изображение на обработку.
	if err := f.processImage(ctx, path, filename); err != nil {
		f.log.Error("failed to process image", "user_id", userID, "path", path, "error", err)
		return err
	}

	resPath := filepath.Join(directoryStorage+"/results", filename)
	file.SetShortUrlUpgrade(&resPath)

	// Сохраняем информацию об оригинале в БД.
	id, err := f.FileRepository.CreateFile(ctx, file)
	if err != nil {
		f.log.Error("failed to create file in database", "user_id", userID, "error", err)
		return err
	}

	f.log.Info("file stored in database", "user_id", userID, "file_id", id)

	return nil
}

func (f *File) processImage(ctx context.Context, imagePath, filename string) error {
	const inpaintURL = "http://localhost:8081/inpaint"

	// Открываем оригинальное изображение.
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer imageFile.Close()

	// Формируем multipart body.
	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return fmt.Errorf("create multipart file: %w", err)
	}

	if _, err := io.Copy(part, imageFile); err != nil {
		return fmt.Errorf("copy image to multipart: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	// Создаём HTTP-запрос.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		inpaintURL,
		&requestBody,
	)
	if err != nil {
		return fmt.Errorf("create inpaint request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send image to inpaint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"inpaint server returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	// Создаём директорию для результатов.
	resultsDir := filepath.Join(directoryStorage, "results")

	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("create results directory: %w", err)
	}

	// Сохраняем результат с исходным именем.
	resultPath := filepath.Join(resultsDir, filename)

	resultFile, err := os.Create(resultPath)
	if err != nil {
		return fmt.Errorf("create result file: %w", err)
	}

	_, copyErr := io.Copy(resultFile, resp.Body)
	closeErr := resultFile.Close()

	if copyErr != nil {
		return fmt.Errorf("save processed image: %w", copyErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close result file: %w", closeErr)
	}

	f.log.Info(
		"image successfully processed",
		"original", imagePath,
		"result", resultPath,
	)

	return nil
}

func NewFileService(fileRepository FileRepository, userRepository UserRepository, log *slog.Logger) *File {
	return &File{
		FileRepository: fileRepository,
		UserRepository: userRepository,
		log:            log,
	}
}
