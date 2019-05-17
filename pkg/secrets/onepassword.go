package secrets

import (
	"encoding/json"
	"fmt"
)

type OnePasswordBackend struct {
	Backend
	OnePasswordClient
}

func NewOnePasswordBackend(vault string) *OnePasswordBackend {
	backend := &OnePasswordBackend{}
	backend.Init()
	return backend
}

func (b *OnePasswordBackend) Init(params ...interface{}) error {
	fmt.Println("Initializing 1password backend.")

	return nil
}

// Parse the 'fields' array in the result of the command below. Return the 'v' property of the field object of which the 'n' property matches parameter key.
//
// $ op get item "key"
//
func (b *OnePasswordBackend) Get(key string) (string, error) {
	opItemString := b.OnePasswordClient.Get(key)

	var opItem OpItem

	json.Unmarshal([]byte(opItemString), &opItem)

	var value = opItem.Details.Sections[0].Fields[0].V

	return value, nil
}

type OnePasswordClient struct {
}

func (c *OnePasswordClient) Get(key string) string {
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
