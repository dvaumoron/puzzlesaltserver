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
package main

import (
	"log"
	"net"
	"os"
	"strconv"

	redisclient "github.com/dvaumoron/puzzleredisclient"
	"github.com/dvaumoron/puzzlesaltserver/saltserver"
	pb "github.com/dvaumoron/puzzlesaltservice"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if godotenv.Overload() == nil {
		log.Println("Loaded .env file")
	}

	saltLen, err := strconv.Atoi(os.Getenv("SALT_LENGTH"))
	if err != nil {
		log.Fatal("Failed to parse SALT_LENGTH")
	}

	lis, err := net.Listen("tcp", ":"+os.Getenv("SERVICE_PORT"))
	if err != nil {
		log.Fatal("Failed to listen :", err)
	}

	rdb := redisclient.Create()

	s := grpc.NewServer()
	pb.RegisterSaltServer(s, saltserver.New(rdb, saltLen))
	log.Println("Listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve :", err)
	}
}
