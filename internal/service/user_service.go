package service

import (
	"context"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
)

type userRepository interface {
	Create(ctx context.Context, input domain.RegisterUserInput) (domain.User, error)
	List(ctx context.Context, filter domain.ListUsersFilter) ([]domain.User, error)
	GetByEmail(ctx context.Context, email string, includeDeleted bool) (domain.User, error)
	GetByID(ctx context.Context, id int64, includeDeleted bool) (domain.User, error)
	Update(ctx context.Context, id int64, input domain.UpdateUserInput) (domain.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

////// methods
////// methods
////// methods

func (s *UserService) Create(ctx context.Context, input domain.RegisterUserInput) (domain.User, error) {
	if err := input.NormalizeAndValidate(); err != nil {
		return domain.User{}, err
	}
	return s.repo.Create(ctx, input)
}

func (s *UserService) List(ctx context.Context, filter domain.ListUsersFilter) ([]domain.User, error) {
	filter.Normalize()
	return s.repo.List(ctx, filter)
}

func (s *UserService) GetByEmail(ctx context.Context, email string, includeDeleted bool) (domain.User, error) {
	// if id <= 0 {
	// 	return domain.User{}, domain.ErrInvalidUserID
	// }
	return s.repo.GetByEmail(ctx, email, includeDeleted)
}

func (s *UserService) GetByID(ctx context.Context, id int64, includeDeleted bool) (domain.User, error) {
	if id <= 0 {
		return domain.User{}, domain.ErrInvalidUserID
	}
	return s.repo.GetByID(ctx, id, includeDeleted)
}

func (s *UserService) Update(ctx context.Context, id int64, input domain.UpdateUserInput) (domain.User, error) {
	if id <= 0 {
		return domain.User{}, domain.ErrInvalidUserID
	}
	if err := input.NormalizeAndValidate(); err != nil {
		return domain.User{}, err
	}
	return s.repo.Update(ctx, id, input)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidUserID
	}
	return s.repo.Delete(ctx, id)
}
