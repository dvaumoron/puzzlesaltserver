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
	"log"

	pb "github.com/dvaumoron/puzzlesaltservice"
	"github.com/go-redis/redis/v8"
)

const redisCallMsg = "Failed during Redis call :"
const generateMsg = "Failed to generate :"

var errInternal = errors.New("internal service error")

// server is used to implement puzzlesaltservice.SaltServer
type server struct {
	pb.UnimplementedSaltServer
	rdb *redis.Client
	len int
}

func New(rdb *redis.Client, saltLen int) pb.SaltServer {
	return server{rdb: rdb, len: saltLen}
}

func (s server) LoadOrGenerate(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	login := request.Login
	salt, err := s.rdb.Get(ctx, login).Result()
	if err == nil {
		return &pb.Response{Salt: salt}, nil
	}
	if err != redis.Nil {
		log.Println(redisCallMsg, err)
		return nil, errInternal
	}

	saltBuffer := make([]byte, s.len)
	_, err = rand.Read(saltBuffer)
	if err != nil {
		log.Println(generateMsg, err)
		return nil, errInternal
	}
	salt = string(saltBuffer)
	if err = s.rdb.Set(ctx, login, salt, 0).Err(); err != nil {
		log.Println(redisCallMsg, err)
		return nil, errInternal
	}
	return &pb.Response{Salt: salt}, nil
}
