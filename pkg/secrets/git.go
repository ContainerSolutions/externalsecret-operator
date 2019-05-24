package secrets

import (
	"bytes"
	"fmt"
	"log"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GitSecretsManagerBackend TODO
type GitSecretsManagerBackend struct {
	Backend
	repositoryURL string
	repository    *git.Repository
}

// NewGitSecretsManagerBackend Return an instance of GitSecretsManagerBackend
func NewGitSecretsManagerBackend(repository *git.Repository) *GitSecretsManagerBackend {
	// repository, err := cloneRepository(repositoryURL)
	backend := &GitSecretsManagerBackend{repository: repository}
	backend.Init()
	return backend
}

// Init Initialisation of the backend
func (s *GitSecretsManagerBackend) Init(params ...interface{}) error {
	return nil
}

// CloneRepository Clone the given repository URL to memory storage
func cloneRepository(repositoryURL string) (*git.Repository, error) {
	fmt.Printf("git clone %s", repositoryURL)
	repository, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repositoryURL,
	})
	if err != nil {
		log.Fatalf("Git clone failed with %s.\n", err)
		return nil, err
	}
	return repository, nil
}

// Get Returns value from repository file defined by `key` path
func (s *GitSecretsManagerBackend) Get(key string) (string, error) {
	//TODO: `key` value shouldn't allow ".git*" values just in case (security!)

	// ... retrieves the branch pointed by HEAD
	ref, err := s.repository.Head()
	fmt.Println(ref.Hash())
	if err != nil {
		log.Fatalf("Git clone failed with %s.\n", err)
		return "", err
	}

	// ... retrieving the commit object
	commit, err := s.repository.CommitObject(ref.Hash())
	// fmt.Println(commit)
	if err != nil {
		log.Fatalf("Git clone failed with %s.\n", err)
		return "", err
	}

	// ... retrieve the tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		log.Fatalf("Retrieving commit failed with %s.\n", err)
		return "", err
	}

	// ... find the file in the object tree
	entry, err := tree.FindEntry(key)
	if err != nil {
		log.Fatalf("Git clone failed with %s.\n", err)
		return "", err
	}

	// ... read the file blob entry the object tree
	fileEntry, err := tree.TreeEntryFile(entry)
	fileReader, err := fileEntry.Blob.Reader()

	buf := new(bytes.Buffer)
	buf.ReadFrom(fileReader)
	secretValue := buf.String()
	fmt.Println(secretValue)

	if err != nil {
		log.Fatalf("Reading file failed with %s.\n", err)
		return "", err
	}

	// return string(secretValue), nil
	return secretValue, nil
}
