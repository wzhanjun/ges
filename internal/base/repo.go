package base

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
)

var unExpandVarPath = []string{"~", ".", ".."}

// Repo is git repository manager.
type Repo struct {
	url   string
	ref   string
	isTag bool
	home  string
}

// NewRepo new a repository manager.
func NewRepo(url string, ref string, isTag bool) *Repo {
	return &Repo{
		url:   url,
		ref:   ref,
		isTag: isTag,
		home:  GESHomeWithDir("repo/" + repoDir(url)),
	}
}

func (r *Repo) Path() string {
	start := strings.LastIndex(r.url, "/")
	end := strings.LastIndex(r.url, ".git")
	if end == -1 {
		end = len(r.url)
	}
	ref := r.ref
	if ref == "" {
		ref = "main"
	}
	return path.Join(r.home, r.url[start+1:end]+"@"+ref)
}

// Pull fetch the repository from remote url.
func (r *Repo) Pull(ctx context.Context) error {
	if r.isTag {
		cmd := exec.CommandContext(ctx, "git", "fetch", "--tags", "origin")
		cmd.Dir = r.Path()
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git fetch failed: %s, %w", string(out), err)
		}
		cmd = exec.CommandContext(ctx, "git", "checkout", r.ref)
		cmd.Dir = r.Path()
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git checkout failed: %s, %w", string(out), err)
		}
		return nil
	}
	cmd := exec.CommandContext(ctx, "git", "symbolic-ref", "HEAD")
	cmd.Dir = r.Path()
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	cmd = exec.CommandContext(ctx, "git", "pull")
	cmd.Dir = r.Path()
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return err
}

// Clone clones the repository to cache path.
func (r *Repo) Clone(ctx context.Context) error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		return r.Pull(ctx)
	}
	var cmd *exec.Cmd
	if r.ref == "" {
		cmd = exec.CommandContext(ctx, "git", "clone", r.url, r.Path())
	} else {
		cmd = exec.CommandContext(ctx, "git", "clone", "-b", r.ref, r.url, r.Path())
	}

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CopyTo(ctx context.Context, to string, modPath string, ignores []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	mod, err := ModulePath(path.Join(r.Path(), "go.mod"))
	if err != nil {
		return err
	}
	return copyDir(r.Path(), to, []string{mod, modPath}, ignores)
}

func repoDir(url string) string {
	vcsURL, err := ParseVCSUrl(url)
	if err != nil {
		return url
	}
	// check host contains port
	host, _, err := net.SplitHostPort(vcsURL.Host)
	if err != nil {
		host = vcsURL.Host
	}
	for _, p := range unExpandVarPath {
		host = strings.TrimLeft(host, p)
	}
	dir := path.Base(path.Dir(vcsURL.Path))
	url = fmt.Sprintf("%s/%s", host, dir)
	return url
}
