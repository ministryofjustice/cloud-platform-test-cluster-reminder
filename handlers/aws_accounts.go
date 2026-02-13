package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func GetAWSAccountName(ctx context.Context, cfg aws.Config) (string, error) {
	iamClient := iam.NewFromConfig(cfg)
	aliases, err := iamClient.ListAccountAliases(ctx, &iam.ListAccountAliasesInput{})
	if err == nil && len(aliases.AccountAliases) > 0 {
		return formatAccountName(aliases.AccountAliases[0]), nil
	}

	if err != nil {
		return "", err
	}

	return "", fmt.Errorf("aws account alias is not configured")
}

func formatAccountName(alias string) string {
	parts := strings.FieldsFunc(alias, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})

	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	return strings.Join(parts, " ")
}
