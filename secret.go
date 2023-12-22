import (
    "context"
    "log"

    "cloud.google.com/go/secretmanager"
)

func getLatestEnabledSecret(ctx context.Context, client *secretmanager.Client, secretName string) (string, error) {
    // Create the request to access the secret.
    accessRequest := &secretmanager.AccessSecretVersionRequest{}

    // Retrieve the latest version.
    result, err := client.AccessSecretVersion(ctx, accessRequest)
    if err != nil {
        return "", err
    }

    // Check if the retrieved version is enabled.
    if result.Payload.State != secretmanager.SecretVersion_ENABLED {
        // If the latest version is disabled, find the most recent enabled version.
        versions, err := client.ListSecretVersions(ctx, &secretmanager.ListSecretVersionsRequest{
            Parent: result.Secret,
        })
        if err != nil {
            return "", err
        }

        for _, version := range versions {
            if version.State == secretmanager.SecretVersion_ENABLED {
                // Retrieve the most recent enabled version.
                result, err = client.AccessSecretVersion(ctx, &secretmanager.AccessSecretVersionRequest{
                    Name: version.Name,
                })
                if err != nil {
                    return "", err
                }
                break
            }
        }
    }

    // The secret payload is available in result.Payload.Data.
    secretData := string(result.Payload.Data)
    return secretData, nil
}

func main() {
    ctx := context.Background()
    client, err := secretmanager.NewClient(ctx)
    if err != nil {
        log.Fatalf("Failed to create secret manager client: %v", err)
    }
    defer client.Close()

    // Replace "your-secret-name" with the actual name of your secret.
    secretName := "projects/your-project-id/secrets/your-secret-name"

    secretData, err := getLatestEnabledSecret(ctx, client, secretName)
    if err != nil {
        log.Fatalf("Failed to access secret: %v", err)
    }

    log.Printf("Secret data: %s", secretData)
}

