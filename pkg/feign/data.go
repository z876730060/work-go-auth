package feign

import (
	"fmt"
	"net/http"
	"time"
)

const WorkData = "work-data"

func init() {
	fClients.Register(WorkData)
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

func GetNew(id int) (*News, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/one/%d", id), nil)
	if err != nil {
		return nil, err
	}

	t, err := do[*News](request, WorkData)
	if err != nil {
		return nil, err
	}
	return t, nil
}
