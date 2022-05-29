package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetParameter(ctx context.Context, key string) ([]byte, error) {
	return getParameter(ctx, key, false)
}

func GetEncryptedParameter(ctx context.Context, key string) ([]byte, error) {
	return getParameter(ctx, key, true)
}

func getParameter(ctx context.Context, key string, encrypted bool) ([]byte, error) {
	svc := ssm.New(session.Must(session.NewSession()))
	key = stripSSMKeyPrefix(key)
	out, err := svc.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: &encrypted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parameter from ssm: %v", err)
	}
	return []byte(*out.Parameter.Value), nil
}

const ssmKeyPrefix = "ssm://"

func stripSSMKeyPrefix(key string) string {
	if strings.HasPrefix(key, ssmKeyPrefix) {
		return strings.TrimPrefix(key, ssmKeyPrefix)
	}
	return key
}
