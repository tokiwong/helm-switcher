package lib_test

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"testing"

	lib "github.com/tokiwong/helm-switcher/lib"
)

// TestDownloadFromURL_FileNameMatch : Check expected filename exist when downloaded
func TestDownloadFromURL_FileNameMatch(t *testing.T) {

	helmURL := "https://github.com/helm/helm/releases/download/"
	installVersion := "helm_"
	installPath := "/.helm.versions_test/"
	goarch := runtime.GOARCH
	goos := runtime.GOOS

	// get current user
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	fmt.Printf("Current user: %v \n", usr.HomeDir)
	installLocation := usr.HomeDir + installPath

	// create /.helm.versions_test/ directory to store code
	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		log.Printf("Creating directory for helm: %v", installLocation)
		err = os.MkdirAll(installLocation, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for helm: %v", installLocation)
			panic(err)
		}
	}

	/* test download lowest helm version */
	lowestVersion := "0.13.9"

	url := helmURL + "v" + lowestVersion + "/" + installVersion + goos + "_" + goarch
	expectedFile := usr.HomeDir + installPath + installVersion + goos + "_" + goarch
	installedFile, _ := lib.DownloadFromURL(installLocation, url)

	if installedFile == expectedFile {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Download file mismatches expected file")
	}

	/* test download latest helm version */
	latestVersion := "0.14.11"

	url = helmURL + "v" + latestVersion + "/" + installVersion + goos + "_" + goarch
	expectedFile = usr.HomeDir + installPath + installVersion + goos + "_" + goarch
	installedFile, _ = lib.DownloadFromURL(installLocation, url)

	if installedFile == expectedFile {
		t.Logf("Expected file name %v", expectedFile)
		t.Logf("Downloaded file name %v", installedFile)
		t.Log("Download file name matches expected file")
	} else {
		t.Logf("Expected file name %v", expectedFile)
		t.Logf("Downloaded file name %v", installedFile)
		t.Error("Downoad file name mismatches expected file")
	}

	cleanUp(installLocation)
}

// TestDownloadFromURL_FileExist : Check expected file exist when downloaded
func TestDownloadFromURL_FileExist(t *testing.T) {

	helmURL := "https://github.com/helm/helm/releases/download/"
	installVersion := "helm_"
	installPath := "/.helm.versions_test/"
	goarch := runtime.GOARCH
	goos := runtime.GOOS

	// get current user
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	fmt.Printf("Current user: %v \n", usr.HomeDir)
	installLocation := usr.HomeDir + installPath

	// create /.helm.versions_test/ directory to store code
	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		log.Printf("Creating directory for helm: %v", installLocation)
		err = os.MkdirAll(installLocation, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for helm: %v", installLocation)
			panic(err)
		}
	}

	/* test download lowest helm version */
	lowestVersion := "0.13.9"

	url := helmURL + "v" + lowestVersion + "/" + installVersion + goos + "_" + goarch
	expectedFile := usr.HomeDir + installPath + installVersion + goos + "_" + goarch
	installedFile, _ := lib.DownloadFromURL(installLocation, url)

	if checkFileExist(expectedFile) {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Downoad file mismatches expected file")
	}

	/* test download latest helm version */
	latestVersion := "0.14.11"

	url = helmURL + "v" + latestVersion + "/" + installVersion + goos + "_" + goarch
	expectedFile = usr.HomeDir + installPath + installVersion + goos + "_" + goarch
	installedFile, _ = lib.DownloadFromURL(installLocation, url)

	if checkFileExist(expectedFile) {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Downoad file mismatches expected file")
	}

	cleanUp(installLocation)
}

func TestDownloadFromURL_Valid(t *testing.T) {

	helmURL := "https://github.com/helm/helm/releases/download/"

	url, err := url.ParseRequestURI(helmURL)
	if err != nil {
		t.Errorf("Valid URL provided:  %v", err)
		t.Errorf("Invalid URL %v", err)
	} else {
		t.Logf("Valid URL from %v", url)
	}
}
