package gqlgen

import (
	"context"
	"fmt"
	"log"
	"os"
	"root/middlewares/auth"
	"root/models"
	"root/servant/AttendsInitialDataProvider"
	"time"

	"github.com/hako/branca"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var envLoadErr = godotenv.Load()

type Resolver struct {
	staffs  []*models.Staff
	staff   *models.Staff
	attends []*models.Attend
	attend  *models.Attend
	auth    *models.Auth
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// que
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateStaff(ctx context.Context, name string, age int) (*models.Staff, error) {
	if envLoadErr != nil {
		log.Fatal(envLoadErr)
	}
	
	if user := auth.UserContextExtracter(ctx); user == nil || !user.IsAdmin {
		return r.staff, fmt.Errorf("Access denied")
	}

	id := ksuid.New()
	aidp := AttendsInitialDataProvider.Summon(id.String())
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc).Format(time.RFC3339)

	r.staff = &models.Staff{
		ID:        id.String(),
		Name:      name,
		Age:       age,
		CreatedAt: &now,
		UpdatedAt: now,
	}

	r.attends = aidp.ServeData()

	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
	defer db.Close()

	db.Create(&r.staff)

	for _, attend := range r.attends {
		db.Create(&attend)
	}

	return r.staff, err
}
func (r *mutationResolver) DeleteStaff(ctx context.Context, id string) (*models.Staff, error) {
	if envLoadErr != nil {
		log.Fatal(envLoadErr)
	}

	if user := auth.UserContextExtracter(ctx); user == nil || !user.IsAdmin {
		return r.staff, fmt.Errorf("Access denied")
	}

	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
	defer db.Close()

	db.Where("id = ?", id).First(&r.staffs)
	r.staff = r.staffs[0]
	db.Where("staff_id = ?", id).Delete(&r.attends)
	db.Where("id = ?", id).Delete(&r.staffs)

	return r.staff, err
}
func (r *mutationResolver) UpdateStaffProfile(ctx context.Context, id string, name string, age int) (*models.Staff, error) {
	if envLoadErr != nil {
		log.Fatal(envLoadErr)
	}

	if user := auth.UserContextExtracter(ctx); user == nil || !user.IsAdmin {
		return r.staff, fmt.Errorf("Access denied")
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc).Format(time.RFC3339)

	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
	defer db.Close()

	db.Where("id = ?", id).First(&r.staffs)
	r.staff = r.staffs[0]
	db.Model(&r.staff).Where("id = ?", id).Update("name", name)
	db.Model(&r.staff).Where("id = ?", id).Update("age", age)
	db.Model(&r.staff).Where("id = ?", id).Update("updated_at", now)

	return r.staff, err
}
func (r *mutationResolver) UpdateAttend(ctx context.Context, id string, is_attend bool, in_time string, out_time string) (*models.Attend, error) {
	if envLoadErr != nil {
		log.Fatal(envLoadErr)
	}

	if user := auth.UserContextExtracter(ctx); user == nil || !user.IsAdmin {
		return r.attend, fmt.Errorf("Access denied")
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc).Format(time.RFC3339)

	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
	defer db.Close()

	db.Where("id = ?", id).First(&r.attends)
	r.attend = r.attends[0]
	db.Model(&r.attend).Where("id = ?", id).Update("is_attend", is_attend)
	db.Model(&r.attend).Where("id = ?", id).Update("in_time", in_time)
	db.Model(&r.attend).Where("id = ?", id).Update("out_time", out_time)
	db.Model(&r.attend).Where("id = ?", id).Update("updated_at", now)

	return r.attend, err
}

// func (r *mutationResolver) CreateUser(ctx context.Context, name string, isAdmin bool, email string, password string) (*models.User, error) {
// 	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
// 	defer db.Close()

// 	id := ksuid.New()

// 	user := &models.User{
// 		ID: id.String(),
// 		Name: name,
// 		IsAdmin: isAdmin,
// 		Email: &email,
// 		Password: &password,
// 	}
// 	db.Create(&user)

// 	return user, err
// }

type queryResolver struct{ *Resolver }

func (r *queryResolver) Staffs(ctx context.Context) ([]*models.Staff, error) {
	if envLoadErr != nil {
		log.Fatal(envLoadErr)
	}

	if user := auth.UserContextExtracter(ctx); user == nil || !user.IsAdmin {
		return r.staffs, fmt.Errorf("Access denied")
	}

	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
	defer db.Close()

	db.Order("created_at asc").Find(&r.staffs)

	for _, staff := range r.staffs {
		db.Where("staff_id = ?", staff.ID).Order("date_start_time asc").Limit(14).Find(&staff.Attends)
	}

	return r.staffs, err
}
func (r *queryResolver) GetCertificate(ctx context.Context, email string, password string) (*models.Auth, error) {
	if envLoadErr != nil {
		log.Fatal(envLoadErr)
	}
	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
	defer db.Close()

	var users []*models.User

	db.Where("email = ? AND password = ?", email, password).First(&users)
	if len(users) == 0 {
		return &models.Auth{}, fmt.Errorf("Authentication Error")
	}

	user := users[0]
	secretKey := os.Getenv("SECRET_KEY")
	b := branca.NewBranca(secretKey)

	userEmail := *user.Email
	userPassword := *user.Password

	token, err := b.EncodeToString(userEmail + "/" + userPassword)
	if err != nil {
		return &models.Auth{}, fmt.Errorf("Token generation failed")
	}

	return &models.Auth{
		UserID:   user.ID,
		UserName: user.Name,
		Token:    token,
	}, err
}
