package couchbase

// Сюда пишем реализацию для коуча
import (
	"fmt"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"github.com/skvoch/burlington-backend/tree/master/internal/repository"
	"strconv"
	"time"
)

const (
	areaBucketName = "area_bucket"
	entityBucketName = "entity_bucket"
)

func New(user, password, host string) (*Repositories, error) {
	cluster, err := gocb.Connect(host, gocb.ClusterOptions{
		Username: user,
		Password: password,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to conncet to couchbase: %v", err)
	}

	areaBucket := cluster.Bucket(areaBucketName)
	if err := areaBucket.WaitUntilReady(time.Second*5, nil); err != nil {
		return nil, fmt.Errorf("failed to conncet to area bucket: %w", err)
	}

	entityBucket := cluster.Bucket(entityBucketName)
	if err := entityBucket.WaitUntilReady(time.Second*5, nil); err != nil{
		return nil, fmt.Errorf("failed to conncet to entity bucket: %w", err)
	}

	return &Repositories{
		cluster: cluster,
		area: &AreaRepo{
			bucket: areaBucket,
		},
		entity: &EntityRepo{
			bucket: entityBucket,
		},
	}, nil
}

// Union repos
type Repositories struct {
	area   *AreaRepo
	entity *EntityRepo

	cluster *gocb.Cluster
}
func (r *Repositories)Entities() repository.EntityRepository{
	var res repository.EntityRepository = r.entity
	return res
}
func (r *Repositories)Areas() repository.AreaRepository{
	var res repository.AreaRepository = r.area
	return res
}


//Area
type AreaRepo struct {
	bucket *gocb.Bucket
	lastId	int
}

func (a *AreaRepo)Create(area models.Area) (string, error){
	defer func(){
		a.lastId += 1
	}()
	if area.ID == ""{
		area.ID = strconv.Itoa(a.lastId)
		return strconv.Itoa(a.lastId), a.Set(strconv.Itoa(a.lastId),area)
	}
	return area.ID, a.Set(area.ID, area)

}

func (a *AreaRepo) Get (id string) (models.Area, error){
	collection := a.bucket.DefaultCollection()
	getRes, err := collection.Get(id, &gocb.GetOptions{})
	if err != nil{
		return models.Area{}, fmt.Errorf("failed to get object with id: %v, err: %v", id, err)
	}
	var res models.Area
	if err := getRes.Content(&res); err != nil{
		return models.Area{}, fmt.Errorf("failed to connect objects, err: %v", err)
	}
	return res, nil
}

func (a *AreaRepo)Set(id string, area models.Area) error{
	collection := a.bucket.DefaultCollection()
	_, err := collection.Upsert(id, area, &gocb.UpsertOptions{})
	if err != nil{
		return fmt.Errorf("failed to upsert object %v, with id %v, err : %v", area, id, err)
	}
	return nil
}

//entity
type EntityRepo struct {
	bucket 	*gocb.Bucket
	lastId	int
}

func (e *EntityRepo)Create(entity models.Entity) (string, error){
	defer func(){
		e.lastId += 1
	}()
	return strconv.Itoa(e.lastId), e.Set(strconv.Itoa(e.lastId),entity)
}

func (e *EntityRepo)Get(id string) (models.Entity, error){
	collection := e.bucket.DefaultCollection()
	getRes, err := collection.Get(id, &gocb.GetOptions{})
	if err != nil{
		return models.Entity{}, fmt.Errorf("failed to get object with id: %v, err: %v", id, err)
	}
	var res models.Entity
	if err := getRes.Content(&res); err != nil{
		return models.Entity{}, fmt.Errorf("failed to connect objects, err: %v", err)
	}
	return res, nil
}

func (e *EntityRepo)Set(id string, entity models.Entity) error{
	collection := e.bucket.DefaultCollection()
	if _, err := collection.Upsert(id, entity, &gocb.UpsertOptions{}); err != nil{
		return fmt.Errorf("failed to upsert object %v, with id %v, err : %v", entity, id, err)
	}
	return nil
}

func (e *EntityRepo)GetByAreaName(areaName string) []models.Entity{
	//TODO:
	return nil
}

