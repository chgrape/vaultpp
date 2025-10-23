package vault

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

type VaultService struct {
	Client *api.Client
}

func AuthVault(client *api.Client) error {
	if os.Getenv("ENV") == "dev" {
		client.SetToken(os.Getenv("VAULT_TOKEN"))
	} else if os.Getenv("ENV") == "prod" {
		secret := &approle.SecretID{FromEnv: "SECRET_ID"}
		auth, err := approle.NewAppRoleAuth(os.Getenv("ROLE_ID"), secret)
		if err != nil {
			log.Fatalf("Unable to initialize AppRole: %v", err)
			return err
		}

		info, err := client.Auth().Login(context.Background(), auth)
		if err != nil {
			log.Fatalf("Unable to login: %v", err)
			return err
		}
		if info == nil {
			log.Fatalf("No auth info")
			return err
		}
	}

	return nil
}

func (s *VaultService) FetchSecret(path string, key string) (string, error) {
	secret, err := s.Client.Logical().Read(fmt.Sprintf("secret/data/%s", path))
	if err != nil {
		return "", err
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no data found at %s", path)
	}

	data := secret.Data["data"].(map[string]any)

	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("value doesn't exist for this key %s", key)
	}

	return value, nil
}
