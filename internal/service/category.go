package service

import (
	"context"
	"io"

	"github.com/celsopires1999/grpc_v2/internal/database"
	"github.com/celsopires1999/grpc_v2/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	dbCategory := database.Category{
		Name:        in.GetName(),
		Description: in.GetDescription(),
	}

	category, err := c.CategoryDB.Create(dbCategory.Name, dbCategory.Description)
	if err != nil {
		return nil, err
	}

	out := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return out, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {

	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}
	out := &pb.CategoryList{}
	for _, category := range categories {
		out.Categories = append(out.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
	return out, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Find(in.GetId())
	if err != nil {
		return nil, err
	}
	out := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
	return out, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}
		category, err := c.CategoryDB.Create(in.GetName(), in.GetDescription())
		if err != nil {
			return err
		}
		out := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}

		categories.Categories = append(categories.Categories, out)
	}
}

func (c *CategoryService) CreateCategoryStreamBidirecional(stream pb.CategoryService_CreateCategoryStreamBidirecionalServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		category, err := c.CategoryDB.Create(in.GetName(), in.GetDescription())
		if err != nil {
			return err
		}
		out := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
		stream.Send(out)
	}
}
