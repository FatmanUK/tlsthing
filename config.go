package main

import (
	"os"
	"github.com/go-git/go-git/v5"
)

type Args struct {
	Path string
	Repo string
	Refresh int
}

type ConfigRepo struct {
	Args
	internal_path string
}

func (re ConfigRepo) validate() error {
	// valid repo url?
	// valid path?
	// valid repo+path if applicable?
	// maybe set remote or local mode
	// refresh rate is allowed value, set default if not
	if re.Refresh < 1 {
		re.Refresh = default_refresh
	}
	return nil
}

func (re ConfigRepo) is_cloned() bool {
	return false
}

func (re ConfigRepo) is_ok_repo() bool {
	return false
}

func (re ConfigRepo) clone() error {
	// seems gogit doesn't support pull, just clone afresh each time
	// check for nonexistent directory / mkdir /var/tlsthing with appropriate ownership
	_, err := git.PlainClone(
		re.internal_path,
		false,   // bare=false
		&git.CloneOptions{
			URL: re.Repo,
			Progress: os.Stdout,})
	return err
}

func (re ConfigRepo) pull() error {
	re.erase()
	err := re.clone()
	return err
}

func (re ConfigRepo) archive() error {
	// tar up archive under another name
	return nil
}

func (re ConfigRepo) erase() {
	// erase repo directory ready for clone
	// rm -r /var/tlsthing/repo
}

func (re ConfigRepo) do_nothing() {
}

func configFactory(args Args) (ConfigRepo, error) {
	re := ConfigRepo{args, internal_repo_path}
	err := re.validate()
	ok := false
	//ok = re.is_cloned()
	//if !ok {
	//	re.clone()
	//}
	ok = re.is_ok_repo()
	if ok {
		re.pull()
	} else {
		re.archive()
		re.erase()
		re.clone()
	}

	return re, err
}
