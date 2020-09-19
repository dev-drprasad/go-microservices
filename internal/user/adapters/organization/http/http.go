package http

import (
	"context"
	"encoding/json"
	"gomicroservices/internal/organization/model"
	"log"
	"net/http"
	"time"
)

type HTTPRepo struct{}

func (repo HTTPRepo) GetBranch(ctx context.Context, id uint64) *model.Branch {
	url := "http://localhost:9092/api/v1/branches/1"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	branch := model.Branch{}
	dec := json.NewDecoder(res.Body)
	dec.Decode(&branch)

	return &branch
}
