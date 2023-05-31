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
	_ "embed"
	"os"
	"strconv"

	grpcserver "github.com/dvaumoron/puzzlegrpcserver"
	redisclient "github.com/dvaumoron/puzzleredisclient"
	"github.com/dvaumoron/puzzlesaltserver/saltserver"
	pb "github.com/dvaumoron/puzzlesaltservice"
)

//go:embed version.txt
var version string

func main() {
	// should start with this, to benefit from the call to godotenv
	s := grpcserver.Make(saltserver.SaltKey, version)

	saltLen, err := strconv.Atoi(os.Getenv("SALT_LENGTH"))
	if err != nil {
		s.Logger.Fatal("Failed to parse SALT_LENGTH")
	}

	rdb := redisclient.Create(s.Logger)

	pb.RegisterSaltServer(s, saltserver.New(rdb, saltLen, s.Logger))

	s.Start()
}
