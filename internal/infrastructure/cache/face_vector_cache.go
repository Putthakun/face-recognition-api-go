package cache

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const hashKey = "face:vectors"

type FaceVectorCache interface {
	LoadFromDB(embeddings map[int64][]float32) error
	Set(empID int64, vector []float32) error
	Remove(empID int64) error
}

type redisFaceCache struct {
	client *redis.Client
}

func NewFaceVectorCache(client *redis.Client) FaceVectorCache {
	return &redisFaceCache{client: client}
}

func (c *redisFaceCache) LoadFromDB(embeddings map[int64][]float32) error {
	if len(embeddings) == 0 {
		log.Println("info: no face vectors to load into Redis")
		return nil
	}

	ctx := context.Background()
	pipe := c.client.Pipeline()
	for empID, vec := range embeddings {
		pipe.HSet(ctx, hashKey, strconv.FormatInt(empID, 10), floatToBytes(vec))
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis pipeline LoadFromDB: %w", err)
	}
	log.Printf("info: loaded %d face vectors into Redis", len(embeddings))
	return nil
}

func (c *redisFaceCache) Set(empID int64, vector []float32) error {
	ctx := context.Background()
	return c.client.HSet(ctx, hashKey, strconv.FormatInt(empID, 10), floatToBytes(vector)).Err()
}

func (c *redisFaceCache) Remove(empID int64) error {
	ctx := context.Background()
	return c.client.HDel(ctx, hashKey, strconv.FormatInt(empID, 10)).Err()
}

func floatToBytes(v []float32) []byte {
	b := make([]byte, len(v)*4)
	for i, f := range v {
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(f))
	}
	return b
}
