package commands

import (
	"fmt"
	"os"

	"github.com/barelyhuman/commitlog/pkg"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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

	err = os.WriteFile(".commitlog.release", []byte(releaser.String()), os.ModePerm)
	if err != nil {
		return
	}

	openRepo, err := git.PlainOpen(c.String("path"))
	if err != nil {
		return err
	}

	var commitHash plumbing.Hash
	wt, err := openRepo.Worktree()
	if err != nil {
		return err
	}

	if c.Bool("commit") {
		wt.Add(".commitlog.release")
		commitHash, err = wt.Commit("chore: version"+releaser.String(), &git.CommitOptions{})
		if err != nil {
			return err
		}

		_, err = openRepo.CreateTag(releaser.String(), commitHash, &git.CreateTagOptions{})
		if err != nil {
			err = fmt.Errorf("looks like there was error while creating a tag for the version commit, please try again or create a tag manually: %v", err)
			return err
		}
	}

	if c.Bool("push") {
		_, err := wt.Status()
		if err != nil {
			return err
		}

		openRepo.Push(&git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		})
	}

	return err
}
