package scategory

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/KDarenskii/catalog-service/internal/app/entity"
	"github.com/KDarenskii/catalog-service/internal/app/repository"
	"github.com/KDarenskii/catalog-service/internal/app/service"
)

type srv struct {
	repoCategory repository.Category
	repoProduct  repository.Product
}

func NewService(repoCategory repository.Category, repoProduct repository.Product) service.Category {
	return &srv{
		repoCategory: repoCategory,
		repoProduct:  repoProduct,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.Category, error) {
	var category entity.Category

	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		existing, err := s.repoCategory.List(ctx, &req.Name)
		if err != nil {
			return err
		}

		if len(existing) > 0 {
			return entity.ErrAlreadyExists
		}

		now := time.Now()
		category = entity.Category{
			GUID:      uuid.Must(uuid.NewV4()),
			Name:      req.Name,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := s.repoCategory.Create(ctx, category); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error) {
	var updatedCategory entity.Category

	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		categories, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{guid})
		if err != nil {
			return err
		}

		if len(categories) == 0 {
			return entity.ErrNotFound
		}

		category := categories[0]

		if req.Name != "" && req.Name != category.Name {
			existingCategories, err := s.repoCategory.List(ctx, &req.Name)
			if err != nil {
				return err
			}

			for _, currentCategory := range existingCategories {
				if currentCategory.GUID != guid {
					return entity.ErrAlreadyExists
				}
			}

			category.Name = req.Name
		}

		category.UpdatedAt = time.Now()

		if err := s.repoCategory.Update(ctx, category); err != nil {
			return err
		}

		updatedCategory = category

		return nil
	})
	if err != nil {
		return entity.Category{}, err
	}

	return updatedCategory, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	return s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		categories, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{guid})
		if err != nil {
			return err
		}

		if len(categories) == 0 {
			return entity.ErrNotFound
		}

		products, err := s.repoProduct.List(ctx, nil, &guid, nil, nil)
		if err != nil {
			return err
		}

		if len(products) > 0 {
			return entity.ErrCategoryHasProducts
		}

		return s.repoCategory.Delete(ctx, guid)
	})
}

func (s *srv) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error) {
	return s.repoCategory.GetByGUIDs(ctx, guids)
}

func (s *srv) List(ctx context.Context) ([]entity.Category, error) {
	return s.repoCategory.List(ctx, nil)
}
