package AttendsInitialDataProvider

import (
	"root/models"
	"time"

	"github.com/segmentio/ksuid"
)

type AttendsInitialDataProvider struct {
	attends []*models.Attend
	StaffId string
}

func Summon(staffId string) *AttendsInitialDataProvider {
	return &AttendsInitialDataProvider{
		StaffId: staffId,
	}
}

func (self *AttendsInitialDataProvider) ServeData() []*models.Attend {

	t := time.Now()
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc).Format(time.RFC3339)

	for i := 0; i < 14; i++ {
		id := ksuid.New()
		self.attends = append(self.attends, &models.Attend{
			ID:            id.String(),
			StaffID:       self.StaffId,
			InTime:        "9:00",
			OutTime:       "17:00",
			IsAttend:      false,
			DateStartTime: time.Date(t.Year(), t.Month(), t.Day()+i, 0, 0, 0, 0, loc).Format(time.RFC3339),
			CreatedAt:     &now,
			UpdatedAt:     now,
		})
	}

	return self.attends
}
