package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/barelyhuman/commitlog/v2/lib"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// commitsByCategory is a collection of commits by a given key
// the _key_ can be `all-changes` or a dynamic key / pattern defined
// by the user
type commitsByCategory struct {
	key     string
	commits []string
}

type Generator struct {
	repo       *git.Repository
	dirPath    string
	startRef   string
	endRef     string
	addPromo   bool
	categories []string
	output     struct {
		stdio    bool
		file     bool
		filePath string
	}
	classifiedCommits []commitsByCategory
	rawCommits        []*object.Commit
}

type GeneratorConfigMod func(*Generator)

func (g *Generator) openRepo() {
	if g.repo != nil {
		return
	}

	r, err := git.PlainOpen(g.dirPath)
	if err != nil {
		log.Fatal("Error opening Repository: ", err)
	}
	g.repo = r
}

func (g *Generator) readCommitsInTags() (err error) {
	// make sure the repo is open
	g.openRepo()

	var commits []*object.Commit

	var latestTagCommit *object.Commit
	var previousTagCommit *object.Commit

	commitsIter, err := g.repo.Log(&git.LogOptions{})
	if err != nil {
		err = fmt.Errorf("[commitlog] Failed to get commits from given hash %v", err)
		return
	}
	defer commitsIter.Close()

	// loop through the commits to get the
	// most recent tagged commit and the 2nd most recent
	// tagged commit
	for {
		c, err := commitsIter.Next()

		if err == io.EOF {
			break
		}

		isTag := lib.IsHashATag(g.repo, c.Hash)
		if isTag {
			if latestTagCommit == nil {
				latestTagCommit = c
			} else if previousTagCommit == nil {
				previousTagCommit = c
				break
			}
		}
	}

	commitsIter, err = g.repo.Log(&git.LogOptions{From: latestTagCommit.Hash})

	if err != nil {
		err = fmt.Errorf("[commitlog] Failed to get commits from given hash %v", err)
		return
	}
	defer commitsIter.Close()

	for {
		c, err := commitsIter.Next()

		if err == io.EOF {
			break
		}

		// do not include the last matching commit
		if previousTagCommit != nil && c.Hash == previousTagCommit.Hash {
			break
		}

		commits = append(commits, c)
	}

	// this will either have commits between 2 tags or all commits if no tags exist
	g.rawCommits = commits

	return

}

func (g *Generator) readCommitsInRange() (err error) {
	// make sure the repo is open
	g.openRepo()

	var commits []*object.Commit
	var commitsIter object.CommitIter

	startHash := lib.GetCommitFromString(g.repo, g.startRef)
	endHash := lib.GetCommitFromString(g.repo, g.endRef)

	if startHash != nil {
		commitsIter, err = g.repo.Log(&git.LogOptions{From: startHash.Hash})
	} else {
		commitsIter, err = g.repo.Log(&git.LogOptions{})
	}

	if err != nil {
		err = fmt.Errorf("[commitlog] Failed to get commits from given hash %v", err)
		return
	}

	for {
		c, err := commitsIter.Next()

		if err == io.EOF {
			break
		}

		if endHash != nil && c.Hash == endHash.Hash {
			break
		}

		commits = append(commits, c)
	}

	g.rawCommits = commits

	return

}

// ReadCommmits will try to collect commits in the given range
// or default to current tag to the previous tag or current commit to
// recent tag, in the same order of priority
func (g *Generator) ReadCommmits() (err error) {
	// make sure the repo is open
	g.openRepo()

	if len(g.startRef) > 0 || len(g.endRef) > 0 {
		err = g.readCommitsInRange()
		if err != nil {
			return
		}
	} else {
		err = g.readCommitsInTags()
		if err != nil {
			return
		}
	}

	return
}

func (g *Generator) Classify() (err error) {
	// write the classification using the existing commits
	if len(g.categories) == 0 {
		allCommits := []string{}

		for _, commit := range g.rawCommits {
			allCommits = append(allCommits, lib.CommitToLog(commit))
		}

		g.classifiedCommits = []commitsByCategory{
			{
				key:     "All Changes",
				commits: allCommits,
			},
		}

		return
	}

	for _, catg := range g.categories {
		rgx, err := regexp.Compile(catg)
		if err != nil {
			break
		}

		catgCommits := []string{}

		for _, commit := range g.rawCommits {
			if !rgx.Match([]byte(commit.Message)) {
				continue
			}

			catgCommits = append(catgCommits, lib.CommitToLog(commit))
		}

		g.classifiedCommits = append(g.classifiedCommits, commitsByCategory{
			key:     catg,
			commits: catgCommits,
		})
	}

	return

}

func (g *Generator) Generate() (err error) {

	var inMardown strings.Builder

	for _, class := range g.classifiedCommits {

		// title
		inMardown.Write([]byte("##"))
		inMardown.Write([]byte(" "))
		inMardown.Write([]byte(class.key))

		inMardown.Write([]byte("\n"))

		// each commit with 2 returns to separate the long description one's
		for _, commitMsg := range class.commits {
			inMardown.Write([]byte(commitMsg))
			inMardown.Write([]byte("\n\n"))
		}

	}

	if g.output.stdio {
		fmt.Println(inMardown.String())
	}

	if g.output.file {
		err = os.WriteFile(g.output.filePath, []byte(inMardown.String()), os.ModePerm)
		if err != nil {
			return
		}
	}

	return
}

func CreateGenerator(path string, mods ...GeneratorConfigMod) *Generator {
	generator := &Generator{
		dirPath: path,
	}

	for _, mod := range mods {
		mod(generator)
	}

	return generator
}

func WithPromo() GeneratorConfigMod {
	return func(g *Generator) {
		g.addPromo = true
	}
}

func WithOutputToFile(filePath string) GeneratorConfigMod {
	return func(g *Generator) {
		g.output.file = true
		g.output.filePath = filePath
	}
}

func WithOutputToStdio() GeneratorConfigMod {
	return func(g *Generator) {
		g.output.stdio = true
	}
}

func WithStartReference(startRef string) GeneratorConfigMod {
	return func(g *Generator) {
		g.startRef = startRef
	}
}

func WithEndReference(endRef string) GeneratorConfigMod {
	return func(g *Generator) {
		g.endRef = endRef
	}
}

func WithCategories(categories string) GeneratorConfigMod {
	return func(g *Generator) {
		parts := strings.Split(categories, ",")
		g.categories = parts
	}
}
