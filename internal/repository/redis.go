package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"machine-telemetry-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisRepo struct {
	client 			*redis.Client
	ttl    			time.Duration
	logger 			*zap.SugaredLogger
	
}

func NewRedisRepo (addr, password string, logger *zap.SugaredLogger) *RedisRepo {
	rdb := redis.NewClient(&redis.Options{
		Addr:     	addr,
		Password: 	password,
		DB:       	0,
	})
	
	return &RedisRepo{
		client: 	rdb,
		ttl:    	15 * time.Minute,
		logger: 	logger,
	}
}

// func (r *RedisRepo) GetLockOld(ctx context.Context, eventId string) (bool, error) {
// 	val, err := r.client.Get(ctx, eventId).Result()
// 	if err != nil {
// 		// logger будет обрабатывать
// 		return false, err
// 	}
// 	return val == "1", nil
// }

func (r *RedisRepo) AcquireEventLock(ctx context.Context, eventId string) (bool, error) {
	key := fmt.Sprintf("telemetry:event:%s", eventId)
	fmt.Println("redis lock:", key)
	
	ok, err := r.client.SetNX(ctx, key, "1", 24*time.Hour).Result()
	if err != nil {
		r.logger.Errorf("redis setnx error %v",err)
		return false, err
		
	}


	return ok, nil
}
func (r *RedisRepo) ReleaseEventLock(ctx context.Context, eventId string) error {
	key := fmt.Sprintf("telemetry:event:%s", eventId)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisRepo) IncrementOverloadStreak(ctx context.Context, machineId string) (int64, error) {
	key := fmt.Sprintf("machine:streak:%s", machineId)
	cnt, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	_ = r.client.Expire(ctx, key, 30*time.Minute).Err()
	return cnt, nil
}
func (r *RedisRepo) ResetOverloadStreak(ctx context.Context, machineId string) error {
	key := fmt.Sprintf("machine:streak:%s", machineId)
	return r.client.Del(ctx, key).Err()
}


func (r *RedisRepo) GetMachineState(ctx context.Context, machineId string) (*domain.MachineStateDTO, error) {
	key := fmt.Sprintf("machine:state:%s", machineId)
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return  nil, nil
		}
		return nil, err
	}

	state := domain.MachineStateDTO{}
	if err := json.Unmarshal([]byte(val), &state); err != nil {
		return nil, err
	}
	
	return &state, nil
}
func (r *RedisRepo) SetMachineState(ctx context.Context, state domain.MachineStateDTO) error {
	key := fmt.Sprintf("machine:state:%s", state.MachineId)
	
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *RedisRepo) DeleteMachineState(ctx context.Context, machineId string) error {
	key := fmt.Sprintf("machine:state:%s", machineId)
	return r.client.Del(ctx, key).Err()
}
