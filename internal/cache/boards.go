package cache

import (
	"context"
	"fmt"

	"github.com/bentodev01/grello/internal/data"
	"github.com/go-redis/redis/v8"
)

type BoardCache struct {
	Cache *redis.Client
}

func (c BoardCache) Get(ctx context.Context, id string) data.BoardResult {
	key := fmt.Sprintf("/board/%s", id)
	res := c.Cache.Get(ctx, key)

	value, err := res.Result()
	if err != nil {
		return data.BoardResult{Err: err}
	}

	b := data.Board{}
	b.FromJson(value)

	return data.BoardResult{Board: b}

}

func (c BoardCache) Set(ctx context.Context, b data.Board) error {
	key := fmt.Sprintf("/board/%s", b.ID.Hex())
	value, err := b.ToJson()
	if err != nil {
		return err
	}
	res := c.Cache.Set(ctx, key, value, 0)
	return res.Err()
}
