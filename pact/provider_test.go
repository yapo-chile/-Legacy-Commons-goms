package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"

	"github.mpi-internal.com/Yapo/goms/pkg/infrastructure"
)

const (
	errorNotFound string = "404"
)

type PactConf struct {
	BrokerHost        string `env:"PACT_BROKER_HOST" envDefault:"http://pact-broker.dev.yapo.cl"`
	BrokerPort        string `env:"PACT_BROKER_PORT" envDefault:"80"`
	ProviderHost      string `env:"PACT_PROVIDER_HOST" envDefault:"http://localhost"`
	ProviderPort      string `env:"PACT_PROVIDER_PORT" envDefault:"8080"`
	PactPath          string `env:"PACTS_PATH" envDefault:"./pacts"`
	PactProvidersPath string `env:"PACT_PROVIDERS_PATH" envDefault:"./mocks"`
}

// Detail has the consumer version details of a pact test
type Detail struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	Href  string `json:"href"`
}

// ConsumerVersion represents the data of the consumer version of a pact test
type ConsumerVersion struct {
	Details Detail `json:"pb:consumer-version"`
}

// JSONTemp stores the variables from the json of the pact test
// that we are going to use
type JSONTemp struct {
	Links        ConsumerVersion `json:"_links"`
	Interactions []interface{}   `json:"interactions"`
}

// A temporary logger is created to be used by the HTTPhandler
type loggerMock struct {
	*testing.T
}

// It has all the logger functions that a normal logger has
func (m *loggerMock) Debug(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Info(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Warn(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Error(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Crit(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Success(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}

// Example Provider Pact: How to run me!
// 1. Start the daemon with `./pact-go daemon`
// 2. cd <pact-go>/examples
// 3. go test -v -run TestProvider

func TestProvider(t *testing.T) {
	var conf PactConf
	infrastructure.LoadFromEnv(&conf)
	fmt.Printf("Pact directory: %+v ", conf.PactPath)

	var pact = &dsl.Pact{
		Consumer: "goms",
	}
	files, err := ioutil.ReadDir(conf.PactPath)
	if err != nil {
		fmt.Printf("Error while reading microservice contract. Error %+v", err)
	}
	for _, file := range files {
		// Verify the Provider with local Pact Files
		h := types.VerifyRequest{
			ProviderBaseURL:       conf.ProviderHost + ":" + conf.ProviderPort,
			PactURLs:              []string{conf.PactPath + "/" + file.Name()},
			CustomProviderHeaders: []string{"Authorization: basic e5e5e5e5e5e5e5"},
		}
		_, err := pact.VerifyProvider(t, h)
		if err != nil {
			fmt.Printf("Error verifying the provider.Error %+v\n", err)
			return
		}
	}
}

func TestSendBroker(*testing.T) {
	pactPublisher := &dsl.Publisher{}
	var conf PactConf
	newVer := 1.0

	infrastructure.LoadFromEnv(&conf)
	fmt.Printf("Conf: %+v\n", conf)

	files, err := ioutil.ReadDir(conf.PactProvidersPath)
	if err != nil {
		fmt.Printf("Error while reading mock files. Error %+v", err)
	}

	// Publishing providers pacts
	for _, file := range files {
		f := conf.PactProvidersPath + "/" + file.Name()
		fmt.Printf("Provider file: %+v\n", f)

		currentPact, currentVer, err := getContractInfo(conf.BrokerHost +
			"/pacts/provider/" + fileWithoutExtension(file.Name()) + "/consumer/goms/latest")
		if err != nil && !strings.Contains(err.Error(), errorNotFound) {
			fmt.Printf("Error getting the contract from the broker: +%v\n", err)
			return
		}

		newPact, err := getProviderPactFile(f)
		if err != nil {
			fmt.Printf("Error generating the contract from the local file: +%v\n", err)
			return
		}

		if currentPact != nil && !reflect.DeepEqual(currentPact, newPact) {
			fmt.Printf("Newer version available! old version: %+v\n", currentVer)
			newVer = currentVer + 0.1
		}

		if currentPact == nil || (newVer > currentVer) {
			fmt.Printf("Publishing pact, version: %.1f\n", newVer)
			err := pactPublisher.Publish(types.PublishRequest{
				PactURLs:        []string{f},
				PactBroker:      conf.BrokerHost + ":" + conf.BrokerPort,
				ConsumerVersion: fmt.Sprintf("%.1f", newVer),
				Tags:            []string{"goms"},
			})
			if err != nil {
				fmt.Printf("Error with the Pact Broker server. Error %+v\n", err)
				return
			}
		}
	}
}

func fileWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

// Method that gets the last version of the contract between consumer (goms) and it's provider
// Returns the contract and the last version of the contract or error
func getContractInfo(url string) (interface{}, float64, error) {
	var confContracts infrastructure.Config
	var result JSONTemp
	infrastructure.LoadFromEnv(&confContracts)
	logger := &loggerMock{}

	HTTPHandler := infrastructure.NewHTTPHandler(logger)
	httprequest := HTTPHandler.NewRequest().
		SetMethod("GET").
		SetPath(url)

	publishedContract, err := HTTPHandler.Send(httprequest)
	if err != nil {
		return nil, -1, err
	}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", publishedContract)), &result)
	if (err != nil) || (len(result.Interactions) < 1) {
		return nil, -1, err
	}

	pactResponse := result.Interactions[0].(map[string]interface{})
	delete(pactResponse, "_id")
	versionFloat, err := strconv.ParseFloat(result.Links.Details.Name, 64)
	if err != nil {
		return nil, -1, err
	}
	fmt.Printf("Current version %+v\n", versionFloat)
	return pactResponse, versionFloat, nil
}

func getProviderPactFile(file string) (interface{}, error) {
	var result JSONTemp
	providerPact, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	resp := fmt.Sprintf("%s", providerPact)
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return nil, err
	}

	pactResponse := result.Interactions[0].(map[string]interface{})
	return pactResponse, err
}
