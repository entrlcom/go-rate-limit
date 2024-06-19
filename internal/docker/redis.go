package docker

import (
	"context"

	"github.com/redis/rueidis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type Redis struct {
	client    rueidis.Client
	container *redis.RedisContainer
}

func NewRedis(ctx context.Context, opts ...RedisOption) (Redis, error) {
	ops := make([]testcontainers.ContainerCustomizer, 0, len(opts)+1)

	ops = append(ops, WithRedisVersion(""))

	for _, opt := range opts {
		ops = append(ops, opt)
	}

	container, err := redis.RunContainer(ctx, ops...)
	if err != nil {
		return Redis{}, err
	}

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		return Redis{}, err
	}

	opt, err := rueidis.ParseURL(uri)
	if err != nil {
		return Redis{}, err
	}

	client, err := rueidis.NewClient(opt)
	if err != nil {
		return Redis{}, err
	}

	redisDocker := Redis{
		client:    client,
		container: container,
	}

	return redisDocker, nil
}

func (x Redis) GetClient() rueidis.Client {
	return x.client
}

func (x Redis) Shutdown(ctx context.Context) error {
	return x.container.Terminate(ctx)
}

type RedisOption testcontainers.ContainerCustomizer

func WithRedisVersion(version string) RedisOption {
	if len(version) == 0 {
		version = "latest"
	}

	return testcontainers.WithImage("redis/redis-stack-server:" + version)
}
