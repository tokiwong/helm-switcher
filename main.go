package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	lib "github.com/tokiwong/helm-switcher/lib"
	"github.com/tokiwong/helm-switcher/modal"
)

const (
	helmURL        = "https://api.github.com/repos/helm/helm/releases?"
	defaultBin     = "/usr/local/bin/helm"
	installFile    = "helm"
	installVersion = "helm_"
	installPath    = "/.helm.versions/"
)

var version = "0.0.5\n"

var clientID = "xxx"
var clientSecret = "xxx"

func main() {

	var client modal.Client

	client.ClientID = clientID
	client.ClientSecret = clientSecret

	custBinPath := getopt.StringLong("bin", 'b', defaultBin, "Custom binary path. For example: /Users/username/bin/helm")
	helpFlag := getopt.BoolLong("help", 'h', "displays help message")
	versionFlag := getopt.BoolLong("version", 'v', "displays the version of helmswitch")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	if *helpFlag {
		usageMessage()
	} else if *versionFlag {
		fmt.Printf("Version: %v\n", version)
	} else {
		if len(args) == 0 {
			helmList, assets := lib.GetAppList(helmURL, &client)
			recentVersions, _ := lib.GetRecentVersions()     //get recent versions from RECENT file
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

		} else if len(args) == 1 {
			semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)
			if semverRegex.MatchString(args[0]) {
				requestedVersion := args[0]

				//check if version is already downloaded before checking if it exists
				/* get current user */
				usr, errCurr := user.Current()
				if errCurr != nil {
					log.Fatal(errCurr)
				}
				/* set installation location */
				installLocation := usr.HomeDir + installPath

				fileInstalled := lib.CheckFileExist(installLocation + installVersion + requestedVersion)

				if fileInstalled {

					/* remove current symlink if exist*/
					symlinkExist := lib.CheckSymlink(*custBinPath)

					if symlinkExist {
						lib.RemoveSymlink(*custBinPath)
					}
					/* set symlink to desired version */
					lib.CreateSymlink(installLocation+installVersion+requestedVersion, *custBinPath)
					fmt.Printf("Switched helm to version %q \n", requestedVersion)
				} else {
					//check if version exist before downloading it
					fmt.Println(requestedVersion + " not found in install path " + installPath)
					fmt.Println("Checking if the version exists...")

					helmList, assets := lib.GetAppList(helmURL, &client)
					exist := lib.VersionExist(requestedVersion, helmList)

					if exist {
						installLocation := lib.Install(helmURL, requestedVersion, assets, custBinPath)
						lib.AddRecent(requestedVersion, installLocation) //add to recent file for faster lookup
					} else {
						fmt.Println("Not a valid helm version")
					}

				}

			} else {
				usageMessage()
			}

		}

	}

}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the helm version as an argument (ex: helmswitch 2.4.13 ), or choose from a menu")
}
