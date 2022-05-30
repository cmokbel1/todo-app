package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetParameter(ctx context.Context, key string) ([]byte, error) {
	return getParameter(ctx, key, false)
}

func GetEncryptedParameter(ctx context.Context, key string) ([]byte, error) {
	return getParameter(ctx, key, true)
}

func getParameter(bgCtx context.Context, key string, encrypted bool) ([]byte, error) {
	ctx, cancel := context.WithTimeout(bgCtx, time.Second*5)
	defer cancel()
	svc := ssm.New(session.Must(session.NewSession()))
	key = stripPrefix(key)
	out, err := svc.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: &encrypted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parameter from ssm: %v", err)
	}
	return []byte(*out.Parameter.Value), nil
}

const ParamStorePrefix = "awsparamstore://"

func stripPrefix(key string) string {
	if strings.HasPrefix(key, ParamStorePrefix) {
		return strings.TrimPrefix(key, ParamStorePrefix)
	}
	return key
}
