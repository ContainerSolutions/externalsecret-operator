package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	ctrl "sigs.k8s.io/controller-runtime"
)

const validObjChars = "0123456789abcdefghijklmnopqrstuvwxyz"

var (
	log = ctrl.Log.WithName("asm")
)

// RandomBytes generate random bytes
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RandomInt returns a random int64
func RandomInt() (int64, error) {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(validObjChars))))
	if err != nil {
		return 0, err
	}

	return randomInt.Int64(), nil
}

// RandomStringObjectSafe returns a random string that is safe to use as an k8s object Name
//  https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
func RandomStringObjectSafe(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", err
	}

	for i := range b {
		randomInt, err := RandomInt()
		if err != nil {
			return "", err
		}
		b[i] = validObjChars[randomInt]
	}
	return string(b), nil

}

// AWSCredentials represents expected credentials
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

/* GetAWSSession returns an aws.session.Session based on the parameters or environment variables
* If parameters are not present or incomplete (secret key, access key AND region)
* then let default config loading order to go on:
* https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
 */
func GetAWSSession(parameters map[string]interface{}, creds []byte, defaultRegion string) (*session.Session, error) {
	awsCreds := &AWSCredentials{}
	if err := json.Unmarshal(creds, awsCreds); err != nil {
		log.Error(err, "Unmarshalling failed")
		return nil, err
	}

	region, ok := parameters["region"].(string)
	if !ok {
		log.Error(nil, "AWS region parameter missing")
		return nil, fmt.Errorf("AWS region parameter missing")
	}

	if region == "" {
		region = defaultRegion
	}

	return session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			awsCreds.AccessKeyID,
			awsCreds.SecretAccessKey,
			awsCreds.SessionToken),
	})
}
