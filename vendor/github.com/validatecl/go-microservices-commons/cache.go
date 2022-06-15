package commons

import (
	"context"
	"time"

	"github.com/bluele/gcache"
	"github.com/go-kit/kit/endpoint"
)

// KeyResolver funcion que obtiene una key (string) a partir de un interface
type KeyResolver func(interface{}) string

// MakeCacheEndpointMiddleware crea un middleware de cache
func MakeCacheEndpointMiddleware(keyResolve KeyResolver, cacheSize int, ttl time.Duration) endpoint.Middleware {
	cache := gcache.New(cacheSize).LRU().Build()

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, in interface{}) (res interface{}, err error) {
			key := keyResolve(in)

			res, err = cache.Get(key)

			if res != nil {
				return res, nil
			}

			defer func() {
				if err == nil {
					cache.SetWithExpire(key, res, ttl)
				}
			}()

			return next(ctx, in)
		}
	}
}
