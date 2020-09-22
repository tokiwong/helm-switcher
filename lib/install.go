package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"

	"github.com/tokiwong/helm-switcher/modal"
)

const (
	helmURL        = "https://get.helm.sh/"
	installFile    = "helm"
	installVersion = "helm_"
	binLocation    = "/usr/local/bin/helm"
	installPath    = "/.helm.versions/"
	recentFile     = "RECENT"
)

var (
	installLocation  = "/tmp"
	installedBinPath = "/tmp"
)

func init() {
	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	/* set installation location */
	installLocation = usr.HomeDir + installPath

	/* set default binary path for helm */
	installedBinPath = binLocation

	/* find helm binary location if helm is already installed*/
	cmd := NewCommand("helm")
	next := cmd.Find()

	/* overrride installation default binary path if helm is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		installedBinPath = path
	}

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(installedBinPath)

	if symlinkExist {
		RemoveSymlink(installedBinPath)
	}
	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)
}

//Install : Install the provided version in the argument
func Install(url string, appversion string, assets []modal.Repo, userBinPath *string) string {

	/* If user provided bin path use user one instead of default */
	if userBinPath != nil {
		installedBinPath = *userBinPath
	}

	pathDir := Path(installedBinPath)     //get path directory from binary path
	binDirExist := CheckDirExist(pathDir) //check bin path exist

	if !binDirExist {
		fmt.Printf("Binary path does not exist: %s\n", pathDir)
		fmt.Printf("Please create binary path: %s for helm installation\n", pathDir)
		os.Exit(1)
	}

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(installedBinPath)

	if symlinkExist {
		RemoveSymlink(installedBinPath)
	}

	/* if selected version already exist, */
	/* proceed to download it from the helm release page */
	//url := helmURL + "v" + helmversion + "/" + "helm" + "_" + goos + "_" + goarch

	goarch := runtime.GOARCH
	goos := runtime.GOOS
	urlDownload := ""
	chkDownload := ""

	for _, v := range assets {

		if v.TagName == "v"+appversion {
			if len(v.Assets) > 0 {
				for _, b := range v.Assets {

					matchedOS, _ := regexp.MatchString(goos, b.BrowserDownloadURL)
					matchedARCH, _ := regexp.MatchString(goarch, b.BrowserDownloadURL)
					if matchedOS && matchedARCH {
						// urlDownload = b.BrowserDownloadURL
						urlDownload = "https://get.helm.sh/helm-" + v.TagName + "-" + goos + "-" + goarch + ".tar.gz"
						chkDownload = urlDownload + ".sha256"
						break
					}
				}
			}
			break
		}
	}

	fileInstalled, _ := DownloadFromURL(installLocation, urlDownload)
	tarRead, readErr := os.Open(fileInstalled)
	if readErr != nil {
		fmt.Println("Expected a location, found " + fileInstalled)
	}

	chkInstalled, _ := DownloadFromURL(installLocation, chkDownload)
	verifySha := VerifyChecksum(fileInstalled, chkInstalled)
	if verifySha != true {
		log.Fatal("didn't pass the verify step")
	}

	/* untar the downloaded file*/
	Untar(installLocation, tarRead)
	binDir := installLocation + "/" + goos + "-" + goarch + "/helm"

	/* rename file to helm version name - helm_x.x.x */
	RenameFile(binDir, installLocation+installVersion+appversion)

	err := os.Chmod(installLocation+installVersion+appversion, 0755)
	if err != nil {
		log.Println(err)
	}

	/* set symlink to desired version */
	CreateSymlink(installLocation+installVersion+appversion, installedBinPath)
	fmt.Printf("Switched helm to version %q \n", appversion)
	return installLocation
}

// AddRecent : add to recent file
func AddRecent(requestedVersion string, installLocation string) {

	semverRegex := regexp.MustCompile(`\d+(\.\d+){2}\z`)

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				RemoveFiles(installLocation + recentFile)
				CreateRecentFile(requestedVersion)
				return
			}
		}

		versionExist := VersionExist(requestedVersion, lines)

		if !versionExist {
			if len(lines) >= 3 {
				_, lines = lines[len(lines)-1], lines[:len(lines)-1]

				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			} else {
				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			}
		}

	} else {
		CreateRecentFile(requestedVersion)
	}
}

// GetRecentVersions : get recent version from file
func GetRecentVersions() ([]string, error) {

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return nil, errRead
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				RemoveFiles(installLocation + recentFile)
				return nil, errRead
			}
		}
		return lines, nil
	}
	return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string) {
	WriteLines([]string{requestedVersion}, installLocation+recentFile)
}

// ValidVersionFormat : returns valid version format
/* For example: 0.1.2 = valid
// For example: 0.1.2-beta1 = valid
// For example: 0.1.2-alpha = valid
// For example: a.1.2 = invalid
// For example: 0.1. 2 = invalid
*/
func ValidVersionFormat(version string) bool {

	// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
	// Follow https://semver.org/spec/v1.0.0-beta.html
	// Check regular expression at https://rubular.com/r/ju3PxbaSBALpJB
	semverRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+)(-[a-zA-z]+\d*)?$`)

	return semverRegex.MatchString(version)
}
