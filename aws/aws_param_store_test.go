package aws_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/require"

	"github.com/alqh/ssm-param-docker/aws"
)

type mockGetParametersAPI func(ctx context.Context, input *ssm.GetParametersInput) (*ssm.GetParametersOutput, error)

func (m mockGetParametersAPI) GetParameters(ctx context.Context, input *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
	return m(ctx, input)
}

func TestAWSParamStore_Parameters(t *testing.T) {
	t.Run("Returns list of values for the parameters in a given path", func(t *testing.T) {
		client := mockGetParametersAPI(func(ctx context.Context, input *ssm.GetParametersInput) (*ssm.GetParametersOutput, error) {
			t.Helper()

			expectedParameters := []string{
				"/domain_name/service_name/service_id",
				"/domain_name/service_name/service_secret",
			}
			require.Equal(t, expectedParameters, input.Names)
			require.True(t, *input.WithDecryption)

			resultValue := map[string]string{
				"/domain_name/service_name/service_id":     "my-service-id",
				"/domain_name/service_name/service_secret": "my-service-secret",
			}

			ar := ssm.GetParametersOutput{
				Parameters: make([]types.Parameter, 0, len(resultValue)),
			}

			for k, v := range resultValue {
				kVal := k
				vVal := v
				ar.Parameters = append(ar.Parameters, types.Parameter{
					Name:  &kVal,
					Type:  types.ParameterTypeSecureString,
					Value: &vVal,
				})
			}

			return &ar, nil
		})

		a := aws.NewAWSParamStore(client)

		res, err := a.Parameters(context.Background(), "/domain_name/service_name", []string{"service_id", "service_secret"})
		require.NoError(t, err)

		require.Equal(t, "my-service-id", res["service_id"])
		require.Equal(t, "my-service-secret", res["service_secret"])
	})
}
