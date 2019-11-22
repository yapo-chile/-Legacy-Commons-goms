package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
type Detail struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	Href  string `json:"href"`
}
type ConsumerVersion struct {
	Details Detail `json:"pb:consumer-version"`
}
type JSONTemp struct {
	Links        ConsumerVersion `json:"_links"`
	Interactions []interface{}   `json:"interactions"`
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
	newVer := 0.0
	//var contractJSON JSONTemp

	infrastructure.LoadFromEnv(&conf)

	oldPactResponse, currentVer, err := getContractInfo(conf.BrokerHost +
		"/pacts/provider/profile-ms/consumer/goms/latest")

	if err != nil {
		fmt.Printf("Error getting old contract %+v\n", err)
	}
	newPactResponse, err := getJSONPactFile(pactDir)
	if err != nil {
		fmt.Printf("Error getting pact response to send %+v\n", err)
	}
	if oldPactResponse != nil {
		newVer = currentVer + 0.1
	}
	fmt.Printf("Old pact %+v\n", oldPactResponse)
	fmt.Printf("New pact %+v\n", newPactResponse)

	err = pactPublisher.Publish(types.PublishRequest{
		PactURLs:        []string{"./pacts/goms.json"},
		PactBroker:      conf.BrokerHost + ":" + conf.BrokerPort,
		ConsumerVersion: fmt.Sprintf("%.1f", newVer),
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
func getContractInfo(url string) (interface{}, float64, error) {

	var conf2 infrastructure.Config

	infrastructure.LoadFromEnv(&conf2)
	prometheus := infrastructure.MakePrometheusExporter(
		conf2.PrometheusConf.Port,
		conf2.PrometheusConf.Enabled,
	)
	logger, _ := infrastructure.MakeYapoLogger(&conf2.LoggerConf,
		prometheus.NewEventsCollector(
			"goms_service_events_total",
			"events tracker counter for goms service",
		),
	)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	HTTPHandler := infrastructure.NewHTTPHandler(logger)
	httprequest := HTTPHandler.NewRequest().
		SetMethod("GET").
		SetPath(url)
	publishedContract, err := HTTPHandler.Send(httprequest)
	if err != nil {

		return nil, -1, err
	}
	resp := fmt.Sprintf("%s", publishedContract)

	var result JSONTemp
	json.Unmarshal([]byte(resp), &result)

	pactRsponse := result.Interactions[1].(map[string]interface{})["response"]

	versionFloat, err := strconv.ParseFloat(result.Links.Details.Name, 64)
	if err != nil {

		return nil, -1, err
	}
	return pactRsponse, versionFloat, nil

}
func getJSONPactFile(pactDir string) (interface{}, error) {
	var result JSONTemp
	file, err := IOReadDir(pactDir)
	if err != nil {
		fmt.Printf("Error in reading files. Error %+v", err)
	}
	pactFileToSend, err := ioutil.ReadFile(pactDir + "/" + file[0])
	if err != nil {
		return nil, err
	}
	resp := fmt.Sprintf("%s", pactFileToSend)
	json.Unmarshal([]byte(resp), &result)
	pactResponse := result.Interactions[1].(map[string]interface{})["response"]
	return pactResponse, err
}
