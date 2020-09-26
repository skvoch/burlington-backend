package couchbase

import (
	"errors"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/couchbase/gocbcore/v9/memd"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"strconv"
)

const (
	idsKey = "idsKey"
)

type IDs struct {
	LastID int64 `json:"last_id"`
}

func (e *EntityRepo) createIDsDocumentIfNotExist() error {
	_, err := e.bucket.DefaultCollection().Get(idsKey, nil)

	var er *gocb.KeyValueError
	if errors.As(err, &e) && er.StatusCode == memd.StatusKeyNotFound {
		ids := &IDs{
			LastID: 1,
		}

		_, err := e.bucket.DefaultCollection().Upsert(
			idsKey, ids, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EntityRepo) generateID() (int64, error) {
	mops := []gocb.MutateInSpec{
		gocb.IncrementSpec("last_id", 1, &gocb.CounterSpecOptions{}),
	}
	incrementResult, err := e.bucket.DefaultCollection().MutateIn(idsKey, mops, &gocb.MutateInOptions{})
	if err != nil {
		return 0, err
	}

	var value int64
	err = incrementResult.ContentAt(0, &value)

	if err != nil {
		return 0, err
	}

	return value, nil
}

func (e *EntityRepo) Get(id string) (models.Entity, error) {
	collection := e.bucket.DefaultCollection()
	getRes, err := collection.Get(id, &gocb.GetOptions{})
	if err != nil {
		return models.Entity{}, fmt.Errorf("failed to get object with id: %v, err: %v", id, err)
	}
	var res models.Entity
	if err := getRes.Content(&res); err != nil {
		return models.Entity{}, fmt.Errorf("failed to connect objects, err: %v", err)
	}
	return res, nil
}

func (e *EntityRepo) Set(entity models.Entity) error {
	if entity.ID == int64(0) {
		id, err := e.generateID()
		if err != nil {
			return fmt.Errorf("failed to generate new id")
		}
		entity.ID = id
	}

	collection := e.bucket.DefaultCollection()
	if _, err := collection.Upsert(strconv.FormatInt(entity.ID, 10), entity, &gocb.UpsertOptions{}); err != nil {
		return fmt.Errorf("failed to upsert object %v, with id %v, err : %v", entity, entity.ID, err)
	}
	return nil
}

func (e *EntityRepo) GetByAreaName(areaName string) []models.Entity {
	//TODO:
	return nil
}
