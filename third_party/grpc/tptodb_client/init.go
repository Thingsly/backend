/*
 *
 * Copyright 2015 gRPC authors.
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

// Package main implements a client for Greeter service.
package tptodb

import (
	"fmt"
	"log"

	pb "github.com/Thingsly/backend/third_party/grpc/tptodb_client/grpc_tptodb"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var TptodbClient pb.ThingsPanelClient

// GrpcTptodbInit initializes the gRPC connection to the Tptodb server
func GrpcTptodbInit() {
	fmt.Println("Initializing gRPC connection to Tptodb...")
	var conn *grpc.ClientConn
	var err error
	grpcHost := viper.GetString("grpc.tptodb_server") // Get the gRPC server address from configuration

	// Try to connect to the gRPC server
	conn, err = grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Failed to load Redis configuration:", err)
		return
	}

	// If connection is successful, initialize the gRPC client
	TptodbClient = pb.NewThingsPanelClient(conn)
}
