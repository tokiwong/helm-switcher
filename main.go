package main

import (
	"fmt"
	"os"
	"log"

	"github.com/pborman/getopt"
	"github.com/manifoldco/promptui"
	lib "github.com/tokiwong/helm-switcher/lib"
	"github.com/tokiwong/helm-switcher/modal"
)

const (
	helmURL    = "https://api.github.com/repos/helm/helm/releases?"
	defaultBin = "/usr/local/bin/helm"
)

var CLIENT_ID = "xxx"
var CLIENT_SECRET = "xxx"

func main() {

	var client modal.Client

	client.ClientID = CLIENT_ID
	client.ClientSecret = CLIENT_SECRET

	custBinPath := getopt.StringLong("bin", 'b', defaultBin, "Custom binary path. For example: /Users/username/bin/helm")

	helmList, assets := lib.GetAppList(helmURL, &client)
	recentVersions, _ := lib.GetRecentVersions() //get recent versions from RECENT file
	helmList = append(recentVersions, helmList...)   //append recent versions to the top of the list
	helmList = lib.RemoveDuplicateVersions(helmList) //remove duplicate version

	/* prompt user to select version of helm */
	prompt := promptui.Select{
		Label: "Select helm version",
		Items: helmList,
	}

	_, helmVersion, errPrompt := prompt.Run()

	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	installLocation := lib.Install(helmURL, helmVersion, assets, custBinPath)
	lib.AddRecent(helmVersion, installLocation) //add to recent file for faster lookup
	os.Exit(0)

	fmt.Println(helmList)


}
