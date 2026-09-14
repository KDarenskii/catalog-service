package mcatv1

import (
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/KDarenskii/catalog-service/internal/app/entity"
	catalogv1 "github.com/KDarenskii/catalog-service/internal/pkg/grpc/gen/catalog/v1"
)

func ProductToProto(p entity.Product) *catalogv1.Product {
	desc := ""

	if p.Description != nil {
		desc = *p.Description
	}

	return &catalogv1.Product{
		Guid:         p.GUID.String(),
		Name:         p.Name,
		Description:  desc,
		Price:        p.Price,
		CategoryGuid: p.CategoryGUID.String(),
		CreatedAt:    timestamppb.New(p.CreatedAt),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
	}
}

func GUIDsToStrings(guids []uuid.UUID) []string {
	rawGuids := make([]string, len(guids))

	for index, guid := range guids {
		rawGuids[index] = guid.String()
	}

	return rawGuids
}

func GUIDsFromStrings(rawGuids []string) ([]uuid.UUID, error) {
	guids := make([]uuid.UUID, len(rawGuids))

	for index, rawGuid := range rawGuids {
		guid, err := uuid.FromString(rawGuid)
		if err != nil {
			return []uuid.UUID{}, err
		}

		guids[index] = guid
	}

	return guids, nil
}
