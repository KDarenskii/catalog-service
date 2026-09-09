package sproduct

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/KDarenskii/catalog-service/internal/app/entity"
	"github.com/KDarenskii/catalog-service/internal/app/repository/mocks"
	"github.com/KDarenskii/catalog-service/internal/pkg/testutil"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

////////////////////////////////////////////////////////////////////////////////
///// Create SUITE ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type createProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

func (s *createProductSuite) TestCreate() {
	type args struct {
		req entity.RequestProductCreate
	}

	type want struct {
		err error
	}

	someError := errors.New("Some error")

	categoryGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{{GUID: categoryGUID}}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Existing Product",
					Price:        500,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{{Name: "Existing Product"}}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, nil).
					Once()
			},
		},
		{
			name: "product List method returns error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: someError},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, someError).
					Once()
			},
		},
		{
			name: "category GetByGUIDs method returns error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: someError},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, someError).
					Once()
			},
		},
		{
			name: "product Create method returns error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: someError},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{{GUID: categoryGUID}}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).
					Return(someError).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Create(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(tc.args.req.Name, result.Name)
				s.Equal(tc.args.req.Description, result.Description)
				s.Equal(tc.args.req.Price, result.Price)
				s.Equal(tc.args.req.CategoryGUID, result.CategoryGUID)
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////
///// Update SUITE ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type updateProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *updateProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}

func (s *updateProductSuite) TestUpdate() {
	type args struct {
		guid uuid.UUID
		req  entity.RequestProductUpdate
	}

	type want struct {
		product entity.Product
		err     error
	}

	someError := errors.New("some error")

	productGuid := uuid.Must(uuid.NewV4())
	categoryGUID := uuid.Must(uuid.NewV4())
	now := time.Now()
	createdAt := now.AddDate(0, 0, -5)

	fullUpdatedProduct := entity.Product{
		GUID:         productGuid,
		Name:         "New Product Name",
		Description:  testutil.PtrString("A new product description"),
		Price:        5000,
		CategoryGUID: categoryGUID,
		CreatedAt:    createdAt,
	}

	partialUpdatedProduct := entity.Product{
		GUID:         productGuid,
		Name:         "New Product Name",
		Description:  testutil.PtrString("Old product description"),
		Price:        1000,
		CategoryGUID: categoryGUID,
		CreatedAt:    createdAt,
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "full update",
			args: args{
				req: entity.RequestProductUpdate{
					Name:         "New Product Name",
					Description:  testutil.PtrString("A new product description"),
					Price:        5000,
					CategoryGUID: categoryGUID,
				},
				guid: productGuid,
			},
			want: want{
				err:     nil,
				product: fullUpdatedProduct,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{
						{
							GUID:         args.guid,
							Name:         "Old product name",
							Description:  nil,
							Price:        1000,
							CategoryGUID: uuid.Must(uuid.NewV4()),
							UpdatedAt:    now.AddDate(0, 0, -1),
							CreatedAt:    createdAt,
						},
					}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{{GUID: args.req.CategoryGUID}}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						want := fullUpdatedProduct
						want.UpdatedAt = p.UpdatedAt
						return reflect.DeepEqual(p, want)
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "partial update",
			args: args{
				req: entity.RequestProductUpdate{
					Name:         "New Product Name",
					Description:  nil,
					Price:        0,
					CategoryGUID: uuid.Nil,
				},
				guid: productGuid,
			},
			want: want{
				err:     nil,
				product: partialUpdatedProduct,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				oldProduct := partialUpdatedProduct
				oldProduct.Name = "Old Product Name"

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						productWithUpdatedAt := partialUpdatedProduct
						productWithUpdatedAt.UpdatedAt = p.UpdatedAt
						return reflect.DeepEqual(p, productWithUpdatedAt)
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "not found",
			args: args{
				req: entity.RequestProductUpdate{
					Name: "New Product Name",
				},
				guid: productGuid,
			},
			want: want{
				err:     entity.ErrNotFound,
				product: entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
		{
			name: "duplicate name",
			args: args{
				req: entity.RequestProductUpdate{
					Name: "New Product Name",
				},
				guid: productGuid,
			},
			want: want{
				err:     entity.ErrAlreadyExists,
				product: entity.Product{GUID: productGuid},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid, Name: "Old Product Name"}}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{{Name: "New Product Name"}}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductUpdate{
					CategoryGUID: categoryGUID,
				},
				guid: productGuid,
			},
			want: want{
				err:     entity.ErrNotFound,
				product: entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid, Name: "Old Product Name"}}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, nil).
					Once()
			},
		},
		{
			name: "product GetByGUIDs method returns error",
			args: args{
				req:  entity.RequestProductUpdate{},
				guid: productGuid,
			},
			want: want{
				err:     someError,
				product: entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, someError).
					Once()
			},
		},
		{
			name: "product List method returns error",
			args: args{
				req: entity.RequestProductUpdate{
					Name: "New product name",
				},
				guid: productGuid,
			},
			want: want{
				err:     someError,
				product: entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, someError).
					Once()
			},
		},
		{
			name: "category GetByGUIDs method returns error",
			args: args{
				req: entity.RequestProductUpdate{
					CategoryGUID: categoryGUID,
				},
				guid: productGuid,
			},
			want: want{
				err:     someError,
				product: entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, someError).
					Once()
			},
		},
		{
			name: "product Update method returns error",
			args: args{
				req: entity.RequestProductUpdate{
					Price: 5000,
				},
				guid: productGuid,
			},
			want: want{
				err:     someError,
				product: entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return true
					})).
					Return(someError).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Update(s.ctx, tc.args.guid, tc.args.req)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(tc.args.guid, result.GUID)
				s.Equal(tc.want.product.Name, result.Name)
				s.Equal(tc.want.product.Description, result.Description)
				s.Equal(tc.want.product.Price, result.Price)
				s.Equal(tc.want.product.CategoryGUID, result.CategoryGUID)
				s.Equal(tc.want.product.CreatedAt, result.CreatedAt)
				s.Greater(result.UpdatedAt, result.CreatedAt)
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////
///// GetByGUIDs SUITE ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type getByGUIDsProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *getByGUIDsProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestGetByGUIDsProductSuite(t *testing.T) {
	suite.Run(t, new(getByGUIDsProductSuite))
}

func (s *getByGUIDsProductSuite) TestGetByGUIDs() {
	type args struct {
		guids []uuid.UUID
	}

	type want struct {
		products []entity.Product
		err      error
	}

	productGUID := uuid.Must(uuid.NewV4())
	productsGuids := []uuid.UUID{productGUID}

	expectedProduct := entity.Product{
		GUID:         productGUID,
		Name:         "Test Product",
		Description:  testutil.PtrString("A test product"),
		Price:        1000,
		CategoryGUID: uuid.Must(uuid.NewV4()),
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "single product found",
			args: args{
				guids: productsGuids,
			},
			want: want{
				products: []entity.Product{expectedProduct},
				err:      nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{expectedProduct}, nil).
					Once()
			},
		},
		{
			name: "not found returns empty slice",
			args: args{guids: productsGuids},
			want: want{err: nil, products: []entity.Product{}},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.GetByGUIDs(s.ctx, tc.args.guids)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tc.want.products, result)
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////
///// GetByGUIDs SUITE ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type deleteProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *deleteProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func (s *deleteProductSuite) TestDelete() {
	type args struct {
		guid uuid.UUID
	}

	type want struct {
		err error
	}

	productGuid := uuid.Must(uuid.NewV4())

	deleteErr := errors.New("delete error")

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{guid: productGuid},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(ctx2 context.Context) error) error {
						return fn(ctx)
					}).Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).Return(nil).Once()
			},
		},
		{
			name: "not found",
			args: args{guid: productGuid},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(ctx2 context.Context) error) error {
						return fn(ctx)
					}).Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
		{
			name: "delete error",
			args: args{guid: productGuid},
			want: want{err: deleteErr},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(ctx2 context.Context) error) error {
						return fn(ctx)
					}).Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: productGuid}}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).
					Return(deleteErr).
					Once()
			},
		},
		{
			name: "product GetByGUIDs method error",
			args: args{guid: productGuid},
			want: want{err: deleteErr},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(ctx2 context.Context) error) error {
						return fn(ctx)
					}).Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, deleteErr).
					Once()

			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			err := s.srv.Delete(s.ctx, tc.args.guid)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////
///// List SUITE ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type listProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *listProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestListProductSuite(t *testing.T) {
	suite.Run(t, new(listProductSuite))
}

func (s *listProductSuite) TestList() {
	type args struct {
		req entity.RequestProductList
	}

	type want struct {
		products []entity.Product
		err      error
	}

	expectedProduct1 := entity.Product{
		GUID:         uuid.Must(uuid.NewV4()),
		Name:         "Test Product 1",
		Description:  testutil.PtrString("A test product"),
		Price:        1000,
		CategoryGUID: uuid.Must(uuid.NewV4()),
	}

	expectedProduct2 := entity.Product{
		GUID:         uuid.Must(uuid.NewV4()),
		Name:         "Test Product 2",
		Description:  testutil.PtrString("A test product"),
		Price:        500,
		CategoryGUID: uuid.Must(uuid.NewV4()),
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductList{
					CategoryGUID: nil,
					MinPrice:     nil,
					MaxPrice:     nil,
				},
			},
			want: want{
				products: []entity.Product{expectedProduct1, expectedProduct2},
				err:      nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx, (*string)(nil), args.req.CategoryGUID, args.req.MinPrice, args.req.MaxPrice).
					Return([]entity.Product{expectedProduct1, expectedProduct2}, nil).
					Once()
			},
		},
		{
			name: "empty result",
			args: args{
				req: entity.RequestProductList{
					CategoryGUID: nil,
					MinPrice:     nil,
					MaxPrice:     nil,
				},
			},
			want: want{err: nil, products: []entity.Product{}},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx, (*string)(nil), args.req.CategoryGUID, args.req.MinPrice, args.req.MaxPrice).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.List(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tc.want.products, result)
			}
		})
	}
}
