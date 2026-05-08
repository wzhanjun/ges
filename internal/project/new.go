package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/wzhanjun/ges/internal/base"
)

// Project is a project template.
type Project struct {
	Name string
}

// New create a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string, branch string, tag string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return errors.New("project creation cancelled by user")
		}
		e = os.RemoveAll(to)
		if e != nil {
			return e
		}
	}

	ref := branch
	isTag := false
	if tag != "" {
		ref = tag
		isTag = true
	}

	if isTag {
		fmt.Printf("🚀 Creating service %s, layout repo is %s, tag is %s, please wait a moment.\n\n", p.Name, layout, ref)
	} else {
		fmt.Printf("🚀 Creating service %s, layout repo is %s, branch is %s, please wait a moment.\n\n", p.Name, layout, ref)
	}
	repo := base.NewRepo(layout, ref, isTag)
	if err := repo.CopyTo(ctx, to, p.Name, []string{".git", ".github"}); err != nil {
		return err
	}

	// rename cmd/server to cmd/{p.Name}
	//e := os.Rename(
	//	path.Join(to, "cmd", "server"),
	//	path.Join(to, "cmd", p.Name),
	//)
	//if e != nil {
	//	return e
	//}

	base.Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ make run\n"))
	fmt.Println("🤝 Thanks for using GES")
	fmt.Println("📚 Tutorial: https://github.com/wzhanjun/ges")
	return nil
}
