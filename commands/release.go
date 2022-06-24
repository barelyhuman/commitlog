package commands

import (
	"fmt"
	"os"

	"github.com/barelyhuman/commitlog/pkg"
	"github.com/urfave/cli/v2"
)

func Release(c *cli.Context) (err error) {

	_, err = os.ReadFile(".commitlog.release")

	if os.IsNotExist(err) {
		err = fmt.Errorf("couldn't find the release file, please run the `--init` flag first")
		return
	}

	fileData, err := os.ReadFile(".commitlog.release")

	if err != nil {
		err = fmt.Errorf("error reading the version file: %v", err)
		return err
	}

	versionString := string(fileData)

	releaserOpts := []pkg.ReleaserMod{}

	if c.Bool("major") {
		releaserOpts = append(releaserOpts, pkg.WithMajorIncrement(), pkg.WithMinorReset(), pkg.WithPatchReset())
	}

	if c.Bool("minor") {
		releaserOpts = append(releaserOpts, pkg.WithMinorIncrement(), pkg.WithPatchReset())
	}

	if c.Bool("patch") {
		releaserOpts = append(releaserOpts, pkg.WithPatchIncrement())
	}

	releaser, err := pkg.CreateNewReleaser(versionString, releaserOpts...)

	if err != nil {
		return err
	}

	// TODO:
	// add commit
	// add tagging
	// add push
	// methods to the `releaser`

	err = os.WriteFile(".commitlog.release", []byte(releaser.String()), os.ModePerm)

	return err
}
