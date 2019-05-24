package secrets

import (
	"testing"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func TestGitSecretsManagerBackendGet(t *testing.T) {
	secretKey := "README.md"
	secretValue := "secretValue"
	expectedValue := secretValue
	Convey("Given a Git repository", t, func() {

		// Create a new repository in memory
		storage := memory.NewStorage()
		filesystem := memfs.New()

		file, err := filesystem.Create(secretKey)
		_, err = file.Write([]byte(secretValue))
		So(err, ShouldBeNil)

		repository, err := git.Init(storage, filesystem)
		So(err, ShouldBeNil)

		worktree, err := repository.Worktree()
		So(err, ShouldBeNil)

		// Adds the new file to the staging area.
		_, err = worktree.Add(secretKey)
		So(err, ShouldBeNil)

		// Commit
		_, err = worktree.Commit("example go-git commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "John Doe",
				Email: "john@doe.org",
				When:  time.Now(),
			},
		})
		So(err, ShouldBeNil)

		Convey("Given an initialized GitSecretsManagerBackend", func() {

			backend := NewGitSecretsManagerBackend(repository)

			Convey("When retrieving a secret", func() {

				actualValue, err := backend.Get(secretKey)
				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					So(actualValue, ShouldEqual, expectedValue)
				})
			})
		})
	})
}
