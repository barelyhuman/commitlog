package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/barelyhuman/commitlog/lib"
	"github.com/barelyhuman/commitlog/pkg"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/urfave/cli/v2"
)

func Release(c *cli.Context) (err error) {

	fileDir := c.String("path")
	filePath := path.Join(fileDir, ".commitlog.release")

	if c.Bool("init") {
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			err = nil
			os.WriteFile(filePath, []byte("v0.0.0"), os.ModePerm)
			fmt.Println("[commitlog] Initialized commitlog release")
		} else {
			err = fmt.Errorf(".commitlog.release already exists, cannot override")
		}
		return
	}

	_, err = os.Stat(filePath)

	if os.IsNotExist(err) {
		err = fmt.Errorf("couldn't find the release file, please run the `--init` flag first")
		return
	}

	fileData, err := os.ReadFile(filePath)

	if err != nil {
		err = fmt.Errorf("error reading the version file: %v", err)
		return err
	}

	versionString := string(fileData)

	releaserOpts := []pkg.ReleaserMod{
		// add in the pre-tag,
		// will be used only if the pre flag
		// is true
		pkg.WithPreTag(c.String("pre-tag")),
	}

	if c.Bool("major") {
		releaserOpts = append(releaserOpts, pkg.WithMajorIncrement(), pkg.WithMinorReset(), pkg.WithPatchReset())
	}

	if c.Bool("minor") {
		releaserOpts = append(releaserOpts, pkg.WithMinorIncrement(), pkg.WithPatchReset())
	}

	if c.Bool("patch") {
		releaserOpts = append(releaserOpts, pkg.WithPatchIncrement())
	}

	if c.Bool("pre") {
		releaserOpts = append(releaserOpts, pkg.WithPrereleaseIncrement())
	}

	releaser, err := pkg.CreateNewReleaser(versionString, releaserOpts...)

	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(releaser.String()), os.ModePerm)
	if err != nil {
		return
	}

	var gitRepo *git.Repository
	var toTagHash plumbing.Hash
	var repoWt *git.Worktree

	if c.Bool("commit") || c.Bool("push") {
		gitRepo, err = git.PlainOpen(c.String("path"))
		if err != nil {
			return err
		}

		repoWt, err = gitRepo.Worktree()

		if err != nil {
			return err
		}
	}

	if c.Bool("commit") {
		msg := "chore: " + releaser.String()
		repoWt.Add(filePath)
		toTagHash, err = repoWt.Commit(msg, &git.CommitOptions{})
		if err != nil {
			return err
		}

		_, err = gitRepo.CreateTag(releaser.String(), toTagHash, &git.CreateTagOptions{
			Message: msg,
		})
		if err != nil {
			err = fmt.Errorf("looks like there was error while creating a tag for the version commit, please try again or create a tag manually: %v", err)
			return err
		}
	}

	if c.Bool("push") {
		_, err = repoWt.Status()
		if err != nil {
			return err
		}

		cmd := exec.Command("git", "push")
		cmd.Dir = fileDir
		err = lib.Command(cmd)

		if err != nil {
			return err
		}

		cmd = exec.Command("git", "push", "--tags")
		cmd.Dir = fileDir

		if err = lib.Command(cmd); err != nil {
			return err
		}
	}

	return err
}
