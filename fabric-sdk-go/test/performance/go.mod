// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/hyperledger/fabric-sdk-go/test/performance

replace github.com/hyperledger/fabric-sdk-go => ../../

require (
	github.com/golang/protobuf v1.4.2
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23
	github.com/hyperledger/fabric-sdk-go v0.0.0-00010101000000-000000000000
	github.com/hyperledger/fabric-sdk-go/test/integration v0.0.0-20200909154308-842c4b3ea51e // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20201026091529-146b70c837a4
	google.golang.org/grpc v1.31.0
)

go 1.14
