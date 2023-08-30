package commands

import (
	"github.com/barelyhuman/commitlog/v3/pkg"
	"github.com/urfave/cli/v2"
)

func Commitlog(c *cli.Context) (err error) {
	path := c.String("path")
	addPromo := c.Bool("promo")
	out := c.String("out")
	stdio := c.Bool("stdio")
	startRef := c.String("start")
	endRef := c.String("end")
	categories := c.String("categories")

	gOptions := []pkg.GeneratorConfigMod{}

	if addPromo {
		gOptions = append(gOptions,
			pkg.WithPromo(),
		)
	}

	// either write to a file or to the stdio with true by default
	if len(out) > 0 {
		gOptions = append(gOptions,
			pkg.WithOutputToFile(out),
		)
	} else if stdio {
		gOptions = append(gOptions,
			pkg.WithOutputToStdio(),
		)
	}

	if len(startRef) > 0 {
		gOptions = append(gOptions,
			pkg.WithStartReference(startRef),
		)
	}

	if len(endRef) > 0 {
		gOptions = append(gOptions,
			pkg.WithEndReference(endRef),
		)
	}

	if len(categories) > 0 {
		gOptions = append(gOptions,
			pkg.WithCategories(categories),
		)
	}

	generator := pkg.CreateGenerator(path,
		gOptions...)

	err = generator.ReadCommmits()

	if err != nil {
		return err
	}

	err = generator.Classify()

	if err != nil {
		return err
	}

	return generator.Generate()
}
