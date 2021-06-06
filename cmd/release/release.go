package release

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	clog "github.com/barelyhuman/commitlog/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var releaseCmd *flag.FlagSet
var major *bool
var minor *bool
var patch *bool
var beta *string
var tag *string

var semverPrompt = &survey.Select{
	Message: "Choose a semver version:",
	Options: []string{"major", "minor", "patch", "none"},
	Default: "none",
}

var betaPrompt = &survey.Confirm{
	Message: "Is it a beta release?",
}

var betaSuffixPrompt = &survey.Input{
	Message: "Enter the exiting beta suffix, will be also used for the any beta suffix?",
	Default: "beta",
}

// TagVersion - struct holding the broken down tag
type TagVersion struct {
	major string
	minor string
	patch string
	beta  string
}

// Install - add flags and other options
func Install() {
	releaseCmd = flag.NewFlagSet("release", flag.ExitOnError)
	major = releaseCmd.Bool("major", false, "If release is a major one, will increment the x.0.0 ")
	minor = releaseCmd.Bool("minor", false, "If release is a minor one, will increment the 0.x.0 ")
	patch = releaseCmd.Bool("patch", false, "If release is a patch, will increment the 0.0.x ")
	beta = releaseCmd.String("beta", "beta", "If the release is a beta, to add/increment tag with `-beta.x` or mentioned string")
	tag = releaseCmd.String("tag", "", "The Tag to be taken as base")
}

// Run - execute the command
func Run(args []string) {

	var tagToUse = *tag

	isBeta := needsQuestionnaire(args)
	err := releaseCmd.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if tagToUse == "" {
		tagToUse = getTagString()
	}

	createRelease(tagToUse, *major, *minor, *patch, *beta, isBeta)
}

// needsQuestionnaire - Check semver and beta if no args were supplied
func needsQuestionnaire(args []string) bool {
	var semver string
	var isBeta bool

	if len(args) < 1 {
		err := survey.AskOne(semverPrompt, &semver)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		err = survey.AskOne(betaPrompt, &isBeta)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		err = survey.AskOne(betaSuffixPrompt, beta)
		if err != nil {
			fmt.Println(err.Error())
			return false
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

	return isBeta
}

func createRelease(tagString string, increaseMajor bool, increaseMinor bool, increasePatch bool, betaSuffix string, isBeta bool) {
	version, hasV := breakTag(tagString)
	releaseTagString := ""
	isIncreasedSemver := false

	majorAsInt, err := strconv.ParseInt(version.major, 10, 32)
	if err != nil {
		log.Fatal("Error converting to number on version.major", version)
	}
	minorAsInt, err := strconv.ParseInt(version.minor, 10, 32)
	if err != nil {
		log.Fatal("Error converting to number on version.minor", version)
	}
	patchAsInt, err := strconv.ParseInt(version.patch, 10, 32)
	if err != nil {
		log.Fatal("Error converting to number on version.patch", version)
	}
	betaAsInt, err := strconv.ParseInt(version.beta, 10, 32)
	if err != nil {
		log.Fatal("Error converting to number on version.beta", version)
	}

	if hasV {
		releaseTagString += "v"
	}

	if increaseMajor {
		majorAsInt++
		minorAsInt = 0
		patchAsInt = 0
		isIncreasedSemver = true
	}

	if increaseMinor {
		minorAsInt++
		patchAsInt = 0
		isIncreasedSemver = true
	}

	if increasePatch {
		patchAsInt++
		isIncreasedSemver = true
	}

	releaseTagString += fmt.Sprintf("%d.%d.%d", majorAsInt, minorAsInt, patchAsInt)

	if isBeta {
		betaAsInt++
		if isIncreasedSemver {
			betaAsInt = 0
		}
		releaseTagString += fmt.Sprintf("-%s.%d", betaSuffix, betaAsInt)
	}

	fmt.Println(releaseTagString)

	isConfirmed := confirmRelease(releaseTagString)

	if !isConfirmed {
		return
	}

	repo := clog.OpenRepository(".")

	setTag(repo, releaseTagString)
}

func tagExists(tag string, r *git.Repository) bool {
	tagFoundErr := "tag was found"
	tags, err := r.TagObjects()
	if err != nil {
		log.Printf("get tags error: %s", err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		log.Printf("iterate tags error: %s", err)
		return false
	}
	return res
}

func setTag(r *git.Repository, tag string) (bool, error) {
	if tagExists(tag, r) {
		log.Printf("tag %s already exists", tag)
		return false, nil
	}
	log.Printf("Set tag %s", tag)
	h, err := r.Head()
	if err != nil {
		log.Printf("get HEAD error: %s", err)
		return false, err
	}

	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: tag,
	})

	if err != nil {
		log.Printf("create tag error: %s", err)
		return false, err
	}

	return true, nil
}

func confirmRelease(tag string) bool {
	var confirm bool

	confirmReleasePrompt := &survey.Confirm{
		Message: "Do you want me to create the following tag: " + tag + " ?",
	}

	err := survey.AskOne(confirmReleasePrompt, &confirm)

	if err != nil {
		log.Fatalln(err)
	}
	return confirm
}

func breakTag(tagString string) (*TagVersion, bool) {
	hasV := false
	version := &TagVersion{}
	tagSplits := strings.Split(tagString, ".")

	version.major = tagSplits[0]
	version.minor = tagSplits[1]
	version.patch = tagSplits[2]
	version.beta = tagSplits[3]

	// Check if the major version has the letter `v` in the tag
	if len(version.major) > 1 && strings.Contains(version.major, "v") {
		version.major = version.major[len("v"):]
		hasV = true
	}

	if len(version.patch) > 1 && strings.Contains(version.patch, "-"+*beta) {
		version.patch = strings.Replace(version.patch, "-"+*beta, "", -1)
	}

	return version, hasV
}

func getTagString() string {
	currentRepository := clog.OpenRepository(".")
	var tagRef *plumbing.Reference
	if *tag == "" {
		tagRef, _, _ = clog.GetLatestTagFromRepository(currentRepository)
	}
	onlyTag := tagRef.Name().Short()
	return onlyTag
}
