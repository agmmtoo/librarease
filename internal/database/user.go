package database

import (
	"context"
	"librarease/internal/usecase"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID       `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string          `gorm:"column:name;type:varchar(255)"`
	Email     string          `gorm:"column:email;type:varchar(255)"`
	Phone     string          `gorm:"column:phone;type:varchar(255)"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`

	Staffs        []Staff
	Subscriptions []Subscription
	AuthUser      *AuthUser
}

func (User) TableName() string {
	return "users"
}

func (s *service) ListUsers(ctx context.Context, opt usecase.ListUsersOption) ([]usecase.User, int, error) {
	var (
		users  []User
		uusers []usecase.User
		count  int64
	)

	db := s.db.Model([]User{}).WithContext(ctx)

	if opt.Name != "" {
		db = db.Where("name ILIKE ?", "%"+opt.Name+"%")
	}

	if opt.IDs != nil {
		db = db.Where("id IN ?", opt.IDs)
	}

	var (
		orderIn = "DESC"
		orderBy = "created_at"
	)
	if opt.SortBy != "" {
		orderBy = opt.SortBy
	}
	if opt.SortIn != "" {
		orderIn = opt.SortIn
	}

	err := db.
		Count(&count).
		Offset(opt.Skip).
		Limit(opt.Limit).
		Order(orderBy + " " + orderIn).
		Find(&users).
		Error

	if err != nil {
		return nil, 0, err
	}

	for _, u := range users {
		uusers = append(uusers, u.ConvertToUsecase())
	}

	return uusers, int(count), nil
}

func (s *service) GetUserByID(ctx context.Context, id string, opt usecase.GetUserByIDOption) (usecase.User, error) {
	var u User

	db := s.db.WithContext(ctx).Model(&User{})

	if opt.IncludeStaffs {
		db.Preload("Staffs.Library")
	}
	db.Preload("AuthUser")

	err := db.Where("id = ?", id).First(&u).Error
	if err != nil {
		return usecase.User{}, err
	}

	uu := u.ConvertToUsecase()
	if u.Staffs != nil {
		for _, st := range u.Staffs {
			ust := st.ConvertToUsecase()
			if st.Library != nil {
				l := st.Library.ConvertToUsecase()
				ust.Library = &l
			}
			uu.Staffs = append(uu.Staffs, ust)
		}
	}
	return uu, nil
}

func (s *service) CreateUser(ctx context.Context, user usecase.User) (usecase.User, error) {
	u := User{
		Name:  user.Name,
		Email: user.Email,
	}

	err := s.db.WithContext(ctx).Create(&u).Error
	if err != nil {
		return usecase.User{}, err
	}

	return u.ConvertToUsecase(), nil
}

func (s *service) UpdateUser(ctx context.Context, user usecase.User) (usecase.User, error) {
	u := User{
		ID:   user.ID,
		Name: user.Name,
	}

	err := s.db.WithContext(ctx).Where("id = ?", u.ID).Updates(&u).Error
	if err != nil {
		return usecase.User{}, err
	}

	return usecase.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		return err
	}

	return nil
}

// Convert core model to Usecase
func (u User) ConvertToUsecase() usecase.User {
	var d *time.Time
	if u.DeletedAt != nil {
		d = &u.DeletedAt.Time
	}
	return usecase.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeleteAt:  d,
	}
}
