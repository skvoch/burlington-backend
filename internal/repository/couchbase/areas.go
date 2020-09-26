package couchbase

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
)

func (a *AreaRepo) GetAreaList() ([]string, error) {
	return nil, nil
}

func (a *AreaRepo) Get(id string) (models.Area, error) {
	collection := a.bucket.DefaultCollection()
	getRes, err := collection.Get(id, &gocb.GetOptions{})
	if err != nil {
		return models.Area{}, fmt.Errorf("failed to get object with id: %v, err: %v", id, err)
	}
	var res models.Area
	if err := getRes.Content(&res); err != nil {
		return models.Area{}, fmt.Errorf("failed to connect objects, err: %v", err)
	}
	return res, nil
}

func (a *AreaRepo) Set(id string, area models.Area) error {
	collection := a.bucket.DefaultCollection()
	_, err := collection.Upsert(id, area, &gocb.UpsertOptions{})
	if err != nil {
		return fmt.Errorf("failed to upsert object %v, with id %v, err : %v", area, id, err)
	}
	return nil
}
