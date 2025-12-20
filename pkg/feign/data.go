package feign

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const WorkData = "work-data"

var (
	DataServiceInstance = DataService{
		log: slog.Default().With("service", WorkData),
	}
)

func init() {
	FeignClients.register(WorkData)
}

type News struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Time      time.Time `json:"time"`
	Content   string    `json:"content"`
}

type DataService struct {
	log *slog.Logger
}

func (d DataService) GetNew(id int) (map[string]any, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/one/%d", id), nil)
	if err != nil {
		return nil, err
	}
	d.log.Info("get news request", "id", id)
	t, err := do[map[string]any](request, WorkData)
	if err != nil {
		return nil, err
	}
	return t, nil
}
