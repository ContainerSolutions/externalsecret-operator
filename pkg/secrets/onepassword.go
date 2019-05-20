package secrets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type OnePasswordBackend struct {
	Backend
	OnePasswordClient
}

func NewOnePasswordBackend(vault string, client OnePasswordClient) *OnePasswordBackend {
	backend := &OnePasswordBackend{}
	backend.OnePasswordClient = client
	backend.Init()
	return backend
}

func (b *OnePasswordBackend) Init(params ...interface{}) error {
	return nil
}

// Call the 1password client and parse the 'fields' array in the output. Return the 'v' property of the field object of which the 'n' property matches parameter key.
func (b *OnePasswordBackend) Get(key string) (string, error) {
	opItemString := b.OnePasswordClient.Get(key)

	var opItem OpItem

	json.Unmarshal([]byte(opItemString), &opItem)

	var value = opItem.Details.Sections[0].Fields[0].V

	return value, nil
}

type OnePasswordClient interface {
	Get(key string) string
}

type OnePasswordCliClient struct {
}

// Invoke $ op get item 'key'
func (c OnePasswordCliClient) Get(key string) string {
	cmd := exec.Command("/usr/local/bin/op", "get", "item", key)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(string(stderr.Bytes()))
		fmt.Println(string(stdout.Bytes()))
		log.Fatalf("/usr/local/bin/op get item '"+key+"' failed with %s\n", err)
	}
	return string(stdout.Bytes())
}

type StubOnePasswordClient struct {
}

// Return a static JSON output for $ op get item 'testkey'
func (c StubOnePasswordClient) Get(key string) string {
	return `{
		"uuid": "xyz",
		"templateUuid": "003",
		"trashed": "N",
		"createdAt": "2019-05-17T12:40:36Z",
		"updatedAt": "2019-05-17T12:40:58Z",
		"changerUuid": "uvw",
		"itemVersion": 1,
		"vaultUuid": "abc",
		"details": {
		  "fields": [],
		  "notesPlain": "",
		  "sections": [
			{
			  "fields": [
				{
				  "k": "string",
				  "n": "efg",
				  "t": "",
				  "v": "testvalue"
				}
			  ],
			  "name": "Section_hij",
			  "title": ""
			}
		  ]
		},
		"overview": {
		  "URLs": [],
		  "ainfo": "",
		  "pbe": 0,
		  "pgrng": false,
		  "ps": 0,
		  "tags": [],
		  "title": "testkey",
		  "url": ""
		}
	  }
	`
}

type OpItem struct {
	Details Details
}

type Details struct {
	Sections []Section
}

type Section struct {
	Fields []Field
}
type Field struct {
	N string
	V string
}
