package couchbase

// Сюда пишем реализацию для коуча
import (
	"fmt"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/skvoch/burlington-backend/tree/master/internal/repository"
	"time"
)

const (
	areaBucketName   = "area_bucket"
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
	if err := entityBucket.WaitUntilReady(time.Second*5, nil); err != nil {
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

func (r *Repositories) Entities() repository.EntityRepository {
	return r.entity
}
func (r *Repositories) Areas() repository.AreaRepository {
	return r.area
}

//Area
type AreaRepo struct {
	bucket *gocb.Bucket
	lastId int
}

//entity
type EntityRepo struct {
	bucket *gocb.Bucket
	lastId int
}
