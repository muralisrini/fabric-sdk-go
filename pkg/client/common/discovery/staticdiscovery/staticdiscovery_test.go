/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package staticdiscovery

import (
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	fabImpl "github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp/test/mockmsp"
	"github.com/stretchr/testify/assert"
)

func TestStaticDiscovery(t *testing.T) {

	configBackend, err := config.FromFile("../../../../../test/fixtures/config/config_test.yaml")()
	if err != nil {
		t.Fatalf(err.Error())
	}

	config1, err := fabImpl.ConfigFromBackend(configBackend...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	discoveryProvider, err := New(config1)
	if err != nil {
		t.Fatalf("Failed to  setup discovery provider: %s", err)
	}
	discoveryProvider.Initialize(mocks.NewMockContext(mockmsp.NewMockSigningIdentity("user1", "Org1MSP")))

	discoveryService, err := discoveryProvider.CreateDiscoveryService("mychannel")
	if err != nil {
		t.Fatalf("Failed to setup discovery service: %s", err)
	}

	peers, err := discoveryService.GetPeers()
	if err != nil {
		t.Fatalf("Failed to get peers from discovery service: %s", err)
	}

	// One peer is configured for "mychannel"
	expectedNumOfPeeers := 1
	if len(peers) != expectedNumOfPeeers {
		t.Fatalf("Expecting %d, got %d peers", expectedNumOfPeeers, len(peers))
	}

}

func TestStaticDiscoveryWhenChannelIsEmpty(t *testing.T) {
	configBackend, err := config.FromFile("../../../../../test/fixtures/config/config_test.yaml")()
	if err != nil {
		t.Fatalf(err.Error())
	}

	config1, err := fabImpl.ConfigFromBackend(configBackend...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	discoveryProvider, _ := New(config1)
	discoveryProvider.Initialize(mocks.NewMockContext(mockmsp.NewMockSigningIdentity("user1", "Org1MSP")))

	_, err = discoveryProvider.CreateDiscoveryService("")
	assert.Error(t, err, "expecting error when channel ID is empty")
}

func TestStaticLocalDiscovery(t *testing.T) {
	configBackend, err := config.FromFile("../../../../../test/fixtures/config/config_test.yaml")()
	assert.NoError(t, err)

	config1, err := fabImpl.ConfigFromBackend(configBackend...)
	assert.NoError(t, err)

	discoveryProvider, err := New(config1)
	assert.NoError(t, err)

	clientCtx := mocks.NewMockContext(mockmsp.NewMockSigningIdentity("user1", "Org1MSP"))
	discoveryProvider.Initialize(clientCtx)

	discoveryService, err := discoveryProvider.CreateLocalDiscoveryService()
	assert.NoError(t, err)

	localCtx := mocks.NewMockLocalContext(clientCtx, discoveryProvider)

	err = discoveryService.(*localDiscoveryService).Initialize(localCtx)
	assert.NoError(t, err)

	peers, err := discoveryService.GetPeers()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(peers))
}
