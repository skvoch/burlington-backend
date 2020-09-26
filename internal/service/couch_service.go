package service

import (
	"encoding/json"
	"fmt"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"github.com/skvoch/burlington-backend/tree/master/internal/repository/couchbase"
	"strings"
)

var rep *couchbase.Repositories

const (
	name = "Administrator"
	pswd = "password"
	host = "192.168.0.58"
)

func RepInit(repositories *couchbase.Repositories) (*couchbase.Repositories, error){
	if repositories != nil{
		return repositories, nil
	}
	rep, err := couchbase.New(name, pswd, host)
	if err != nil{
		return nil, fmt.Errorf("Can't create repository, err :%w", err)
	}
	return rep, nil
}
//area
func GetArea(id string)(models.Area, error){
	rep, err := RepInit(rep)
	if err != nil{
		return models.Area{}, fmt.Errorf("invalid repository, err: %w", err)
	}
	area := rep.Areas()
	return area.Get(id)
}

func SetArea(jsonObj string) (string, error){
	rep, err := RepInit(rep)
	area := rep.Areas()
	var obj models.Area
	dec := json.NewDecoder(strings.NewReader(jsonObj))
	if err := dec.Decode(&obj); err != nil {
		return "", fmt.Errorf("failed to create json decoder, err:%w", err)
	}
	id, err := area.Create(obj)
	if err != nil{
		return "", fmt.Errorf("failed to upsert object, err: %w", err)
	}
	return fmt.Sprintf("object with id %v was added successfully", id), nil
}

func SetAreaModel(jsonObj models.Area) (string, error){
	rep, err := RepInit(rep)
	if err != nil{
		return "", fmt.Errorf("Repository problem, err : %w", err)
	}
	area := rep.Areas()
	id, err := area.Create(jsonObj)
	if err != nil{
		return "", fmt.Errorf("failed to upsert object, err: %w", err)
	}
	return fmt.Sprintf("object with id %v was added successfully", id), nil
}

//entity

func GetEntity(id string)(models.Entity, error){
	rep, err := RepInit(rep)
	if err != nil{
		return models.Entity{}, fmt.Errorf("Repository problem, err : %w", err)
	}
	entity := rep.Entities()
	return entity.Get(id)
}

func SetEntity(model models.Entity) (string, error) {
	rep, err := RepInit(rep)
	if err != nil{
		return "", fmt.Errorf("Repository problem, err : %w", err)
	}
	entity := rep.Entities()
	id, err := entity.Create(model)
	if err != nil{
		return "", fmt.Errorf("failed to upsert object, err: %w", err)
	}
	return fmt.Sprintf("object with id %v was added successfully", id), nil
}