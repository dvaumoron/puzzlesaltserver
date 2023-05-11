/*
 *
 * Copyright 2023 puzzlesaltserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package saltserver

import (
	"context"
	"crypto/rand"
	"errors"

	pb "github.com/dvaumoron/puzzlesaltservice"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const SaltKey = "puzzleSalt"

const redisCallMsg = "Failed during Redis call"
const generateMsg = "Failed to generate"

var errInternal = errors.New("internal service error")

// server is used to implement puzzlesaltservice.SaltServer
type server struct {
	pb.UnimplementedSaltServer
	rdb    *redis.Client
	len    int
	logger *otelzap.Logger
}

func New(rdb *redis.Client, saltLen int, logger *otelzap.Logger) pb.SaltServer {
	return server{rdb: rdb, len: saltLen, logger: logger}
}

func (s server) LoadOrGenerate(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	logger := s.logger.Ctx(ctx)
	login := request.Login
	salt, err := s.rdb.Get(ctx, login).Result()
	if err == nil {
		return &pb.Response{Salt: []byte(salt)}, nil
	}
	if err != redis.Nil {
		logger.Error(redisCallMsg, zap.Error(err))
		return nil, errInternal
	}

	saltBuffer := make([]byte, s.len)
	_, err = rand.Read(saltBuffer)
	if err != nil {
		logger.Error(generateMsg, zap.Error(err))
		return nil, errInternal
	}
	salt = string(saltBuffer)
	if err = s.rdb.Set(ctx, login, salt, 0).Err(); err != nil {
		logger.Error(redisCallMsg, zap.Error(err))
		return nil, errInternal
	}
	return &pb.Response{Salt: saltBuffer}, nil
}
