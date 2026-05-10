package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/model"
	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/service"
	pb "github.com/NatalyaAsh/ads-platform/services/ad-search/proto_gen/ad_search/v1"
)

type AdSearchHandler struct {
	pb.UnimplementedAdSearchServiceServer
	adService *service.AdService
}

func NewAdSearchHandler(adService *service.AdService) *AdSearchHandler {
	return &AdSearchHandler{
		adService: adService,
	}
}

// CreateAd
func (h *AdSearchHandler) CreateAd(ctx context.Context, req *pb.CreateAdRequest) (*pb.CreateAdResponse, error) {
	ad, err := h.adService.CreateAd(&model.CreateAdRequest{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		UserID:      uint(req.UserId),
		CategoryID:  uint(req.CategoryId),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateAdResponse{
		Id:          uint32(ad.ID),
		Title:       ad.Title,
		Description: ad.Description,
		Price:       ad.Price,
		UserId:      uint32(ad.UserID),
		CategoryId:  uint32(ad.CategoryID),
		Status:      string(ad.Status),
		Views:       int32(ad.Views),
		CreatedAt:   ad.CreatedAt.String(),
		UpdatedAt:   ad.UpdatedAt.String(),
	}, nil
}

// GetAd
func (h *AdSearchHandler) GetAd(ctx context.Context, req *pb.GetAdRequest) (*pb.GetAdResponse, error) {
	ad, err := h.adService.GetAdByID(uint(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "ad not found")
	}

	categoryName := ""
	if ad.Category.ID != 0 {
		categoryName = ad.Category.Name
	}

	return &pb.GetAdResponse{
		Id:           uint32(ad.ID),
		Title:        ad.Title,
		Description:  ad.Description,
		Price:        ad.Price,
		UserId:       uint32(ad.UserID),
		CategoryId:   uint32(ad.CategoryID),
		CategoryName: categoryName,
		Status:       string(ad.Status),
		Views:        int32(ad.Views),
		CreatedAt:    ad.CreatedAt.String(),
		UpdatedAt:    ad.UpdatedAt.String(),
	}, nil
}

// GetUserAds
func (h *AdSearchHandler) GetUserAds(ctx context.Context, req *pb.GetUserAdsRequest) (*pb.GetUserAdsResponse, error) {
	ads, err := h.adService.GetUserAds(uint(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetUserAdsResponse{}
	for _, ad := range ads {
		response.Ads = append(response.Ads, &pb.Ad{
			Id:          uint32(ad.ID),
			Title:       ad.Title,
			Description: ad.Description,
			Price:       ad.Price,
			UserId:      uint32(ad.UserID),
			CategoryId:  uint32(ad.CategoryID),
			Status:      string(ad.Status),
			Views:       int32(ad.Views),
			CreatedAt:   ad.CreatedAt.String(),
			UpdatedAt:   ad.UpdatedAt.String(),
		})
	}

	return response, nil
}

// UpdateAd
func (h *AdSearchHandler) UpdateAd(ctx context.Context, req *pb.UpdateAdRequest) (*pb.UpdateAdResponse, error) {
	ad, err := h.adService.UpdateAd(&model.UpdateAdRequest{
		ID:          uint(req.Id),
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  uint(req.CategoryId),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	categoryName := ""
	if ad.Category.ID != 0 {
		categoryName = ad.Category.Name
	}

	return &pb.UpdateAdResponse{
		Id:           uint32(ad.ID),
		Title:        ad.Title,
		Description:  ad.Description,
		Price:        ad.Price,
		UserId:       uint32(ad.UserID),
		CategoryId:   uint32(ad.CategoryID),
		CategoryName: categoryName,
		Status:       string(ad.Status),
		Views:        int32(ad.Views),
		UpdatedAt:    ad.UpdatedAt.String(),
	}, nil
}

// DeleteAd
func (h *AdSearchHandler) DeleteAd(ctx context.Context, req *pb.DeleteAdRequest) (*pb.DeleteAdResponse, error) {
	err := h.adService.DeleteAd(uint(req.Id))
	if err != nil {
		return &pb.DeleteAdResponse{Success: false}, nil
	}
	return &pb.DeleteAdResponse{Success: true}, nil
}

// ListAds
func (h *AdSearchHandler) ListAds(ctx context.Context, req *pb.ListAdsRequest) (*pb.ListAdsResponse, error) {
	filters := model.AdListFilters{
		CategoryID: uint(req.CategoryId),
		UserID:     uint(req.UserId),
		Status:     req.Status,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Limit:      int(req.Limit),
		Offset:     int(req.Offset),
	}

	if req.MinPrice != 0 {
		minPrice := req.MinPrice
		filters.MinPrice = &minPrice
	}
	if req.MaxPrice != 0 {
		maxPrice := req.MaxPrice
		filters.MaxPrice = &maxPrice
	}

	ads, total, err := h.adService.ListAds(filters)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.ListAdsResponse{Total: total}
	for _, ad := range ads {
		response.Ads = append(response.Ads, &pb.Ad{
			Id:          uint32(ad.ID),
			Title:       ad.Title,
			Description: ad.Description,
			Price:       ad.Price,
			UserId:      uint32(ad.UserID),
			CategoryId:  uint32(ad.CategoryID),
			Status:      string(ad.Status),
			Views:       int32(ad.Views),
			CreatedAt:   ad.CreatedAt.String(),
			UpdatedAt:   ad.UpdatedAt.String(),
		})
	}

	return response, nil
}

// ListCategories
func (h *AdSearchHandler) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	categories, err := h.adService.ListCategories()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.ListCategoriesResponse{}
	for _, cat := range categories {
		response.Categories = append(response.Categories, &pb.Category{
			Id:   uint32(cat.ID),
			Name: cat.Name,
			Slug: cat.Slug,
		})
	}

	return response, nil
}

// GetCategory
func (h *AdSearchHandler) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
	cat, err := h.adService.GetCategoryByID(uint(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "category not found")
	}

	return &pb.GetCategoryResponse{
		Id:   uint32(cat.ID),
		Name: cat.Name,
		Slug: cat.Slug,
	}, nil
}

func (h *AdSearchHandler) UploadMedia(ctx context.Context, req *pb.UploadMediaRequest) (*pb.UploadMediaResponse, error) {
	// Создаём папку для объявления, если её нет
	uploadDir := filepath.Join("uploads", fmt.Sprintf("ad_%d", req.AdId))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create dir: %v", err)
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(req.FileName)
	uniqueName := fmt.Sprintf("%d_%d%s", req.AdId, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, uniqueName)

	// Сохраняем файл на диск
	if err := os.WriteFile(filePath, req.FileData, 0644); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write file: %v", err)
	}

	// Сохраняем метаданные в MongoDB (позже, пока заглушка)
	// TODO: вызвать mediaRepo.Create

	// Формируем URL для доступа к файлу (можно будет отдавать через статический сервер)
	url := fmt.Sprintf("/uploads/ad_%d/%s", req.AdId, uniqueName)

	return &pb.UploadMediaResponse{
		Id:        uniqueName,
		AdId:      req.AdId,
		Url:       url,
		FileName:  req.FileName,
		FileSize:  int64(len(req.FileData)),
		MimeType:  req.MimeType,
		IsPrimary: req.IsPrimary,
	}, nil
}

func (h *AdSearchHandler) GetMedia(ctx context.Context, req *pb.GetMediaRequest) (*pb.GetMediaResponse, error) {
	// Ищем файлы в папке uploads/ad_<id>
	dir := fmt.Sprintf("uploads/ad_%d", req.AdId)
	files, err := os.ReadDir(dir)
	if err != nil {
		return &pb.GetMediaResponse{Media: []*pb.Media{}}, nil // нет картинок
	}

	var mediaList []*pb.Media
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		mediaList = append(mediaList, &pb.Media{
			Id:        file.Name(),
			AdId:      req.AdId,
			Url:       fmt.Sprintf("/uploads/ad_%d/%s", req.AdId, file.Name()),
			FileName:  file.Name(),
			FileSize:  0, // опционально
			MimeType:  "image/jpeg",
			IsPrimary: len(mediaList) == 0,
		})
	}
	return &pb.GetMediaResponse{Media: mediaList}, nil
}
