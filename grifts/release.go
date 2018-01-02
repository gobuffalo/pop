package grifts

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	. "github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = Desc("release", "Generates a CHANGELOG and creates a new GitHub release based on what is in the version.go file.")
var _ = Add("release", func(c *Context) error {
	Run("shoulders", c)
	v, err := findVersion()
	if err != nil {
		return err
	}

	err = installBin()
	if err != nil {
		return err
	}

	err = localTest()
	if err != nil {
		return err
	}

	err = tagRelease(v)
	if err != nil {
		return err
	}

	if err := commitAndPush(v); err != nil {
		return errors.WithStack(err)
	}
	return runReleaser(v)
})

func installBin() error {
	cmd := exec.Command("go", "install", "-v", "./soda")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func localTest() error {
	cmd := exec.Command("./test.sh")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func tagRelease(v string) error {
	cmd := exec.Command("git", "tag", v)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "push", "origin", "--tags")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func commitAndPush(v string) error {
	cmd := exec.Command("git", "push", "origin", "master")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func runReleaser(v string) error {
	cmd := exec.Command("goreleaser", "--rm-dist")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func findVersion() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	vfile, err := ioutil.ReadFile(filepath.Join(pwd, "./soda/cmd/version.go"))
	if err != nil {
		return "", err
	}

	//var Version = "0.4.0"
	re := regexp.MustCompile(`const Version = "(.+)"`)
	matches := re.FindStringSubmatch(string(vfile))
	if len(matches) < 2 {
		return "", errors.New("failed to find the version!")
	}
	v := matches[1]
	return v, nil
}
