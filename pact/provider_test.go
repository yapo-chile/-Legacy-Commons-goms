package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"

	"github.mpi-internal.com/Yapo/suggested-ads/pkg/infrastructure"
)

var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/pacts", dir)

type PactConf struct {
	BrokerHost string `env:"BROKER_HOST" envDefault:"http://3.229.36.112"`
	BrokerPort string `env:"BROKER_PORT" envDefault:"80"`
}

// Example Provider Pact: How to run me!
// 1. Start the daemon with `./pact-go daemon`
// 2. cd <pact-go>/examples
// 3. go test -v -run TestProvider
func TestProvider(t *testing.T) {

	var conf PactConf
	infrastructure.LoadFromEnv(&conf)
	var pact = &dsl.Pact{
		Consumer: "goms",
		Provider: "profile-ms",
	}
	files, err := IOReadDir(pactDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {

		// Verify the Provider with local Pact Files
		h := types.VerifyRequest{
			ProviderBaseURL:       "http://10.15.1.78:7987",
			PactURLs:              []string{file},
			CustomProviderHeaders: []string{"Authorization: basic e5e5e5e5e5e5e5"},
		}
		_, err := pact.VerifyProvider(t, h)
		if err != nil {
			fmt.Printf("Error verifying the provider.")
		}
		pactPublisher := &dsl.Publisher{}
		err = pactPublisher.Publish(types.PublishRequest{
			PactURLs:        []string{"./mocks"},
			PactBroker:      conf.BrokerHost + ":" + conf.BrokerPort,
			ConsumerVersion: "1.0.0",
			Tags:            []string{"goms"},
		})
		if err != nil {
			fmt.Printf("Error with the Pact Broker server. Error %+v", err)
		}

	}

	if err != nil {
		fmt.Printf("Error listing files %q: %v\n", dir, err)
		return
	}

}
func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}
