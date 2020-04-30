package main

import (
	"os"
	"os/exec"
	//	"github.com/go-git/go-git/v5"
	//	"github.com/go-git/go-git/v5/plumbing/object"
)

// CheckGit returns nil if directory contains git working git tree, error otherwise
func CheckGit(directory string) error {
	current, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(current)
	err = os.Chdir(directory)
	if err != nil {
		return err
	}
	_, err = exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return err
	}
	return nil
}

// PushRemote performs git steps on supplied file in directory: stage, commit, push
func PushRemote(directory string, fileToCommit string, comments []string) (string, error) {
	// Opens an already existing repository.
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}
	defer os.Chdir(current)
	err = os.Chdir(directory)
	if err != nil {
		return "Cannot set directory", err
	}
	_, err = exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return "Status error:", err
	}

	// Adds the new file to the staging area.
	_, err = exec.Command("git", "add", fileToCommit).Output()
	if err != nil {
		return "Staging error", err
	}

	// We can verify the current status of the worktree using the method Status.
	//Info("git status --porcelain")
	_, err = exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return "Status staging error:", err
	}
	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	comment := "-m\""
	for i, c := range comments {
		comment += c
		if i < len(comments)-1 {
			comment += ","
		}
	}
	comment += "\""
	_, err = exec.Command("git", "commit", comment).Output()
	if err != nil {
		return "Commit error:", err
	}
	_, err = exec.Command("git", "push").Output()
	if err != nil {
		return "Push error:", err
	}
	return "Success", nil
}
