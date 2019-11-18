package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"

	"github.mpi-internal.com/Yapo/goms/pkg/infrastructure"
)

var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/pacts", dir)

type PactConf struct {
	BrokerHost   string `env:"BROKER_HOST" envDefault:"http://3.229.36.112"`
	BrokerPort   string `env:"BROKER_PORT" envDefault:"80"`
	ProviderHost string `env:"PROVIDER_HOST" envDefault:"http://localhost"`
	ProviderPort string `env:"PROVIDER_PORT" envDefault:"8080"`
}

// Example Provider Pact: How to run me!
// 1. Start the daemon with `./pact-go daemon`
// 2. cd <pact-go>/examples
// 3. go test -v -run TestProvider
func TestProvider(t *testing.T) {
	fmt.Printf("Pact directory: %+v", pactDir)
	var conf PactConf
	infrastructure.LoadFromEnv(&conf)
	var pact = &dsl.Pact{
		Consumer: "goms",
		Provider: "profile-ms",
		Port:     6666,
	}
	files, err := IOReadDir(pactDir)
	if err != nil {
		fmt.Printf("Error in reading files. Error %+v", err)

	}
	for _, file := range files {
		// Verify the Provider with local Pact Files
		h := types.VerifyRequest{
			ProviderBaseURL:       conf.ProviderHost + ":" + conf.ProviderPort,
			PactURLs:              []string{pactDir + "/" + file},
			CustomProviderHeaders: []string{"Authorization: basic e5e5e5e5e5e5e5"},
		}
		_, err := pact.VerifyProvider(t, h)
		if err != nil {
			fmt.Printf("Error verifying the provider.Error %+v\n", err)
			return
		}
	}
}
func TestSendBroker(t *testing.T) {
	pactPublisher := &dsl.Publisher{}
	var conf PactConf
	var conf2 infrastructure.Config
	type JSONTemp struct {
		Provider string
	}
	var contractJSON JSONTemp

	infrastructure.LoadFromEnv(&conf)
	infrastructure.LoadFromEnv(&conf)
	prometheus := infrastructure.MakePrometheusExporter(
		conf2.PrometheusConf.Port,
		conf2.PrometheusConf.Enabled,
	)
	logger, err := infrastructure.MakeYapoLogger(&conf2.LoggerConf,
		prometheus.NewEventsCollector(
			"goms_service_events_total",
			"events tracker counter for goms service",
		),
	)
	HTTPHandler := infrastructure.NewHTTPHandler(logger)
	httprequest := HTTPHandler.NewRequest().
		SetMethod("GET").
		SetPath("http://3.229.36.112/pacts/provider/profile-ms/consumer/goms/latest")
	publishedContract, err := HTTPHandler.Send(httprequest)
	resp := fmt.Sprintf("%s", publishedContract)
	err = json.Unmarshal([]byte(resp), &contractJSON)
	fmt.Printf("Request preview %+v\n", contractJSON.Provider)

	if err != nil {
		fmt.Printf("Error with request preview %+v\n", err)
	}

	err = pactPublisher.Publish(types.PublishRequest{
		PactURLs:        []string{"./pacts/goms.json"},
		PactBroker:      conf.BrokerHost + ":" + conf.BrokerPort,
		ConsumerVersion: "2.0.0",
		Tags:            []string{"goms"},
	})
	if err != nil {
		fmt.Printf("Error with the Pact Broker server. Error %+v\n", err)
		return
	}

	if err != nil {
		log.Fatalln(err)
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
