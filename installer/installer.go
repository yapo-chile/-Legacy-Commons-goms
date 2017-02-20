package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wingedpig/loom"
	"os"
	"os/exec"
	"path/filepath"
)

// Structure to describe Installer's configuration like:
//
// Name: name for the service
//
// FilesPath: the path to where the RPM files are located
//
// PoyaRepoPath: the path in poya's repo server where the repo is placed
type ServiceConfig struct {
	Name         string
	FilesPath    string
	PoyaRepoPath string
	ExecPath     string
}

// Structure to describe the server's related to a poya configuration
//
// Service: the server where the service is installed
//
// Repo: Repo and Logs server
type PoyaConfig struct {
	Payments loom.Config `json:"payments"`
	Repo     loom.Config `json:"repo"`
}

// Structure to describe main installer's configuration json file
//
// Service: Installer configuration structure
//
// Staging: Loom structure to Staging's Repo server
//
// Poya1-4: Poya enviroment structure
type Config struct {
	Service ServiceConfig
	Staging loom.Config
	Poya1   PoyaConfig
	Poya2   PoyaConfig
	Poya3   PoyaConfig
	Poya4   PoyaConfig
}

// function to parse and load data from json configuration file.
// Returns a Config structure with the values.
func loadConfig(configPath string) Config {
	file, _ := os.Open(configPath)
	decoder := json.NewDecoder(file)
	configStruct := Config{}
	err := decoder.Decode(&configStruct)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configStruct
}

// Global variable to represent the configuration
var config Config

// Copy a tar zipped file with all the RPM files placed in "Config.Service.FilesPath" path into
// Staging server
func copyPackages() {
	var tgzFileName string
	RPMfiles := make([]string, 0)

	// open folder
	d, err := os.Open(config.Service.FilesPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// for each file in build/RPMS/x86_64
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".rpm" {
			RPMfiles = append(RPMfiles, file.Name())
			if file.Name()[:8] == config.Service.Name {
				tgzFileName = file.Name()
			}
		}
	}

	// create tgz
	tgzFileName = tgzFileName[:len(tgzFileName)-11] + ".tgz"
	args := []string{"czf", tgzFileName, "-C", config.Service.FilesPath}
	args = append(args, RPMfiles...)
	_, err = exec.Command("tar", args...).Output()
	if err != nil {
		fmt.Println("failed tgz creation", err)
		os.Exit(1)
	}

	// copy file to ch9stg
	config.Staging.Put(tgzFileName, "/tmp/"+tgzFileName)
	config.Staging.Run("rm -f /tmp/goms.tgz")
	config.Staging.Run("ln -s /tmp/" + tgzFileName + " /tmp/goms.tgz")
}

// Get a Poya Loom Structure from a PoyaName string
func getPoya(poyaName string) PoyaConfig {
	var poya PoyaConfig
	switch poyaName {
	case "poya1":
		poya = config.Poya1
	case "poya2":
		poya = config.Poya2
	case "poya3":
		poya = config.Poya3
	case "poya4":
		poya = config.Poya4
	}
	return poya
}

// Copy RPM files placed in "Config.Service.FilesPath" to a specific Poya Repo server
// (described by poyaName), update the repo server, then update "*goms*" packages
// in Payments of Poya environment.   Finally restart creditos service
func deployToPoya(poyaName string) {
	poya := getPoya(poyaName)
	// copy files to poya repo
	poya.Repo.Run("mkdir -p /tmp/.goms")
	poya.Repo.Put(config.Service.FilesPath+"/*.rpm", "/tmp/.goms")
	poya.Repo.Sudo("mv /tmp/.goms/*.rpm " + config.Service.PoyaRepoPath + "/x86_64")
	poya.Repo.Run("rm -Rf /tmp/.goms")
	poya.Repo.Sudo("createrepo " + config.Service.PoyaRepoPath)
	// update packages in server
	poya.Payments.Sudo("yum clean all")
	poya.Payments.Sudo("yum --disablerepo=\"*\" --enablerepo=\"yapo\" update -y \"*goms*\"")
	poya.Payments.Sudo(config.Service.ExecPath + " restart")
}

// Run "rpm -qa *goms*" command in Payment server
func listPoyaRpm(poyaName string) {
	poya := getPoya(poyaName)
	poya.Payments.Run("rpm -qa *goms*")
}

// The main function start parsing command line flags, the options recognized are this:
//
// -config=/path/to/file.json: path to json config file. Default value: ./installer-config.json
//
// -copy=[true|false]: copy RPM packages in a tar zipped file to Staging Repo server.
// Default value: false
//
// -list-rpm=[true|false]: list the RPM creditos packages installed in Payments server of a
// specific Poya environment. Default value: false
//
// -deploy=[true|false]: update the creditos RPM packages in Poya's Repo server, then update
// creditos packages in Payments servers.
//
// -poya=[poya1|poya2|poya3|poya4]: the Poya environment where the packages will be installed,
// this is a required parameter if "-deploy" or "-list-rpm" flags are true. Default value: ""
func main() {
	var configPath = flag.String("config", "installer-config.json", "Path to JSON config file")
	var copy = flag.String("copy", "false", "Copy packages to ch9stg")
	var listRpm = flag.String("list-rpm", "false", "List goms packages in poya")
	var deploy = flag.String("deploy", "false", "Deploy packages to poya")
	var poya = flag.String("poya", "", "Poya where packages should be deployed")
	flag.Parse()

	config = loadConfig(*configPath)

	if *copy == "true" {
		fmt.Println("==> copy packages")
		copyPackages()
	}

	if *deploy == "true" {
		if *poya != "" {
			fmt.Printf("==> deploy packages to %s\n", *poya)
			deployToPoya(*poya)
		} else {
			fmt.Println("Missing poya name")
		}
	}

	if *listRpm == "true" {
		if *poya != "" {
			fmt.Printf("==> list packages in %s\n", *poya)
			listPoyaRpm(*poya)
		} else {
			fmt.Println("Missing poya name")
		}
	}
}
