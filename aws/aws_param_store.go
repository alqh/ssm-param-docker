package aws

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pkg/errors"
)

type AWSParamStoreClient interface {
	GetParameters(ctx context.Context, input *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
}

type AWSParamStore struct {
	client AWSParamStoreClient
}

func NewAWSParamStore(client AWSParamStoreClient) *AWSParamStore {
	return &AWSParamStore{client}
}

func (a *AWSParamStore) Parameters(ctx context.Context, paramPath string, paramKeys []string) (map[string]string, error) {
	decrypt := true
	req := ssm.GetParametersInput{
		Names:          make([]string, 0, len(paramKeys)),
		WithDecryption: &decrypt,
	}
	prependPath := paramPath + "/"
	for _, k := range paramKeys {
		req.Names = append(req.Names, prependPath+k)
	}

	res, err := a.client.GetParameters(ctx, &req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch parameters")
	}

	keyToParameterVal := make(map[string]string, len(res.Parameters))
	for _, p := range res.Parameters {
		k := strings.TrimPrefix(*p.Name, prependPath)
		keyToParameterVal[k] = *p.Value
	}

	if len(res.InvalidParameters) > 0 {
		log.Printf("Unable to find parameters %d parameters: %v \n", len(res.InvalidParameters), res.InvalidParameters)
	}

	return keyToParameterVal, nil
}
