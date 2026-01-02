package usecases

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fiap-161/tc-golunch-core-service/internal/product/entity"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/entity/enum"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/gateway"
	apperror "github.com/fiap-161/tc-golunch-core-service/internal/shared/errors"
	"github.com/google/uuid"
)

type UseCases struct {
	productGateway gateway.Gateway
}

func Build(productGateway gateway.Gateway) *UseCases {
	return &UseCases{productGateway: productGateway}
}

func (u *UseCases) CreateProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	isValidCategory := enum.IsValidCategory(string(product.Category))
	if !isValidCategory {
		return entity.Product{}, &apperror.ValidationError{Msg: "Invalid category"}
	}

	if err := product.Validate(); err != nil {
		return entity.Product{}, err
	}

	saved, createErr := u.productGateway.Create(ctx, product)
	if createErr != nil {
		return entity.Product{}, &apperror.InternalError{Msg: createErr.Error()}
	}

	return saved, nil
}

func (u *UseCases) ListCategories(ctx context.Context) []enum.Category {
	return enum.GetAllCategories()
}

func (u *UseCases) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	publicURL := os.Getenv("PUBLIC_URL")

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	if !allowedTypes[contentType] {
		return "", &apperror.ValidationError{Msg: "Only JPEG and PNG images are allowed"}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("error resetting file: %w", err)
	}

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return "", fmt.Errorf("error creating upload dir: %w", err)
		}
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
	fullPath := filepath.Join(uploadDir, filename)

	if err := saveFile(fileHeader, fullPath); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	return fmt.Sprintf("%s/uploads/%s", publicURL, filename), nil
}

func saveFile(fileHeader *multipart.FileHeader, dest string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (u *UseCases) GetAllByCategory(ctx context.Context, category string) ([]entity.Product, error) {
	isValidCategory := enum.IsValidCategory(category)
	invalidCategory := !isValidCategory && category != ""

	if invalidCategory {
		return []entity.Product{}, &apperror.ValidationError{Msg: "Invalid category"}
	}

	result, err := u.productGateway.GetAllByCategory(ctx, category)
	if err != nil {
		return []entity.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	return result, nil
}

func (u *UseCases) Update(ctx context.Context, productId string, product entity.Product) (entity.Product, error) {

	_, findErr := u.FindByID(ctx, productId)
	if findErr != nil {
		return entity.Product{}, findErr
	}

	updated, updateErr := u.productGateway.Update(ctx, productId, product)
	if updateErr != nil {
		return entity.Product{}, updateErr
	}

	return updated, nil
}

func (u *UseCases) FindByID(ctx context.Context, productId string) (entity.Product, error) {
	if _, err := uuid.Parse(productId); err != nil {
		return entity.Product{}, &apperror.ValidationError{Msg: "Invalid UUID format for product ID"}
	}

	found, err := u.productGateway.FindByID(ctx, productId)
	if err != nil {
		return entity.Product{}, err
	}

	return found, nil
}

func (u *UseCases) Delete(ctx context.Context, productId string) error {
	if _, err := uuid.Parse(productId); err != nil {
		return &apperror.ValidationError{Msg: "invalid UUID format for product ID"}
	}
	_, findErr := u.FindByID(ctx, productId)
	if findErr != nil {
		return findErr
	}

	deleteErr := u.productGateway.Delete(ctx, productId)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (u *UseCases) FindByIDs(ctx context.Context, productIdList []string) ([]entity.Product, error) {
	if len(productIdList) == 0 {
		return nil, &apperror.ValidationError{Msg: "productIdList cannot be empty"}
	}

	// UUIDs validation
	for _, id := range productIdList {
		if _, err := uuid.Parse(id); err != nil {
			return nil, &apperror.ValidationError{Msg: fmt.Sprintf("Invalid UUID format: %s", id)}
		}
	}

	products, err := u.productGateway.FindByIDs(ctx, productIdList)
	if err != nil {
		return nil, err
	}

	return products, nil

}
