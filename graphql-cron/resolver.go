package graphql_cron

import (
	"context"
	"log"
	"os"
	"root/models"
	"time"

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
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CronUpdateAttends(ctx context.Context) (string, error) {
	if envLoadErr != nil {
		log.Fatal(envLoadErr)
	}

	db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
	defer db.Close()

	db.Exec("DELETE FROM attends WHERE date_start_time = (SELECT min(date_start_time) FROM attends);")

	db.Find(&r.staffs)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc).Format(time.RFC3339)

	for _, staff := range r.staffs {
		id := ksuid.New()
		db.Where("staff_id = ?", staff.ID).Order("date_start_time desc").Limit(1).Find(&r.attends)
		r.attend = r.attends[0]
		t, _ := time.Parse(time.RFC3339, r.attend.DateStartTime)

		newAttend := &models.Attend{
			ID:            id.String(),
			StaffID:       staff.ID,
			InTime:        "9:00",
			OutTime:       "17:00",
			IsAttend:      false,
			DateStartTime: time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, loc).Format(time.RFC3339),
			CreatedAt:     &now,
			UpdatedAt:     now,
		}

		db.Create(&newAttend)
	}

	return "Date update successful", err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Test(ctx context.Context) (string, error) {
	panic("not implemented")
}
