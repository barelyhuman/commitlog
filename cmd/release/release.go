package release

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"

	clog "github.com/barelyhuman/commitlog/log"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"

	survey "github.com/AlecAivazis/survey/v2"
)

var (
	releaseCmd *flag.FlagSet
	major      *bool
	minor      *bool
	patch      *bool
	beta       *bool
	betaSuffix *string
)

const configFileName = ".commitlog.release"

var semverPrompt = &survey.Select{
	Message: "Choose a semver version (choose none for prerelease/beta increments):",
	Options: []string{"major", "minor", "patch", "none"},
	Default: "none",
}

var betaPrompt = &survey.Confirm{
	Message: "Is it a beta release?",
}

var betaSuffixPrompt = &survey.Input{
	Message: "Enter the exiting beta suffix (if any) (eg: beta or dev or canary) :",
	Default: "",
}

var confirmCreation = &survey.Confirm{
	Message: "Do you want me to create a commit for the new version?:",
}

type Config struct {
	version *TagVersion
}

// TagVersion - struct holding the broken down tag
type TagVersion struct {
	major       int64
	minor       int64
	patch       int64
	beta        bool
	betaSuffix  string
	betaVersion int64
}

// Install - add flags and other options
func Install() {
	releaseCmd = flag.NewFlagSet("release", flag.ExitOnError)
	major = releaseCmd.Bool("major", false, "If release is a *major* one, will increment the x.0.0 ")
	minor = releaseCmd.Bool("minor", false, "If release is a *minor* one, will increment the 0.x.0 ")
	patch = releaseCmd.Bool("patch", false, "If release is a *patch*, will increment the 0.0.x ")
	beta = releaseCmd.Bool("beta", false, "If the release is a beta/prerelease")
	betaSuffix = releaseCmd.String("beta-suffix", "", "If the release is a beta, to add/increment tag with `-beta.x` or mentioned string")
}

func bail(err error) {
	if err != nil {
		fail(err.Error())
		panic(1)
	}
}

func bullet(msg string) {
	color.New(color.FgWhite).Add(color.Bold).Println(msg)
}

func fail(msg string) {
	color.New(color.FgRed).Add(color.Bold).Println(msg)
}

func success(msg string) {
	color.New(color.FgGreen).Add(color.Bold).Println(msg)
}

func dim(msg string) {
	color.White(msg)
}

// Run - execute the command
func Run(args []string) {

	err := releaseCmd.Parse(args)

	bail(err)

	isInitialised, err := checkReleaseInit()
	bail(err)

	if !isInitialised {
		bullet("Initializing Commitlog Release")
		err := initialiseRelease()
		bail(err)
		success("Created, .commitlog.release")
		dim("Please, modify the file to match the current latest version or leave it as is if starting with a new project")
	}

	config, err := readConfigFile()
	bail(err)

	askQuestions(args)

	err = createRelease(config, *major, *minor, *patch, *beta, *betaSuffix)

	bail(err)
}

// checkReleaseInit - check if release was already initialised in the particular folder, if not, it'll return false
func checkReleaseInit() (bool, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return false, err
	}

	hasConfigFile := false

	for _, file := range files {
		if file.Name() == configFileName {
			hasConfigFile = true
		}
	}

	return hasConfigFile, nil
}

// initialiseRelease - create the .commitlog.release file in the current folder
// TODO: need to add a question before doing so
func initialiseRelease() error {
	err := os.WriteFile(configFileName, []byte("0.0.0"), os.ModePerm)
	return err
}

// readConfigFile - read the file and parse it as a config, limited to version details for now
func readConfigFile() (Config, error) {
	config := Config{}
	dataInBytes, err := os.ReadFile(configFileName)
	if err != nil {
		return config, err
	}

	asString := string(dataInBytes)

	success("Current Version:")
	bullet(asString)

	version, _ := breakTag(asString)

	config.version = version

	return config, nil
}

// createRelease - create a release based on the read config and given parameters
func createRelease(config Config, incMajor bool, incMinor bool, incPatch bool, incBeta bool, betaSuffixString string) error {
	updatedVersion := TagVersion{
		major:       config.version.major,
		minor:       config.version.minor,
		patch:       config.version.patch,
		beta:        incBeta,
		betaSuffix:  betaSuffixString,
		betaVersion: config.version.betaVersion,
	}

	resetBeta := false

	if incMajor {
		updatedVersion.major++
		updatedVersion.minor = 0
		updatedVersion.patch = 0
		resetBeta = true
	}

	if incMinor {
		updatedVersion.minor++
		updatedVersion.patch = 0
		resetBeta = true
	}

	if incPatch {
		updatedVersion.patch++
		resetBeta = true
	}

	updatedVersionString := fmt.Sprintf("%v.%v.%v", updatedVersion.major, updatedVersion.minor, updatedVersion.patch)

	if updatedVersion.beta {
		if resetBeta {
			updatedVersion.betaVersion = 0
		} else {
			updatedVersion.betaVersion++
		}

		var sb strings.Builder

		if len(updatedVersion.betaSuffix) > 1 {
			sb.WriteString(updatedVersionString + fmt.Sprintf("-%v.", updatedVersion.betaSuffix))
		} else {
			sb.WriteString(updatedVersionString + "-")
		}
		sb.WriteString(fmt.Sprintf("%v", updatedVersion.betaVersion))
		updatedVersionString = sb.String()
	}

	success("New Version")
	bullet(updatedVersionString)

	var confirmed bool

	survey.AskOne(confirmCreation, &confirmed)

	if !confirmed {
		fail("✖ Cancelled")
		return nil
	}

	err := writeUpdatedVersion(updatedVersionString)
	if err != nil {
		return err
	}

	success("✔ Updated Version")

	return nil
}

func writeUpdatedVersion(versionString string) error {
	err := os.WriteFile(configFileName, []byte(versionString), fs.ModePerm)
	if err != nil {
		return err
	}

	err = createCommit(versionString)
	return err
}

func createCommit(versionString string) error {
	repo := clog.OpenRepository(".")

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = wt.Add(configFileName)
	if err != nil {
		return err
	}

	commit, err := wt.Commit(versionString, &git.CommitOptions{})
	if err != nil {
		return err
	}

	_, err = repo.CreateTag(versionString, commit, &git.CreateTagOptions{
		Message: versionString,
	})

	return err
}

// askQuestions - Check semver and beta if no args were supplied
func askQuestions(args []string) error {
	var semver string

	if len(args) < 1 {
		err := survey.AskOne(semverPrompt, &semver)

		if err != nil {
			return err
		}

		err = survey.AskOne(betaPrompt, beta)

		if err != nil {
			return err
		}

		if *beta {
			err = survey.AskOne(betaSuffixPrompt, betaSuffix)
			if err != nil {
				return err
			}
		}

		switch semver {
		case "major":
			{
				*major = true
				break
			}
		case "minor":
			{
				*minor = true
				break
			}
		case "patch":
			{
				*patch = true
				break
			}
		}
	}

	return nil
}

// breakTag - break the given semver version string into proper version values, does support breaking semver pre-release strings
func breakTag(tagString string) (*TagVersion, bool) {
	hasV := false
	version := &TagVersion{}
	tagSplits := strings.Split(tagString, ".")

	majorStringSplit := strings.Split(tagSplits[0], "")

	if len(majorStringSplit) > 1 {
		hasV = true
		major, err := strconv.ParseInt(majorStringSplit[1], 10, 32)
		bail(err)
		version.major = major
	} else {
		major, err := strconv.ParseInt(majorStringSplit[0], 10, 32)
		bail(err)
		version.major = major
	}

	minor, err := strconv.ParseInt(tagSplits[1], 10, 32)
	bail(err)
	version.minor = minor

	if len(tagSplits) > 3 {
		version.beta = true
		betaV, err := strconv.ParseInt(tagSplits[3], 10, 32)
		bail(err)
		version.betaVersion = betaV
	} else {
		version.betaVersion = -1
	}

	patchStringSplit := strings.Split(tagSplits[2], "-")

	patch, err := strconv.ParseInt(patchStringSplit[0], 10, 32)
	bail(err)
	version.patch = patch

	if len(patchStringSplit) > 1 {
		version.betaSuffix = patchStringSplit[1]
	}

	return version, hasV
}
