package testlib

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"gotest.tools/assert"
)

const TLS_SECRET_PASSWORD_YAML_TEMPLATE = `---
apiVersion: v1
kind: Secret
metadata:
  name: %s
  namespace: %s
apiVersion: v1
data:
  %s: %s
  password: %s
`

const TLS_SECRET_NO_PASSWORD_YAML_TEMPLATE = `---
apiVersion: v1
kind: Secret
metadata:
  name: %s
  namespace: %s
apiVersion: v1
data:
  %s: %s
`

func ReadAll(path string) ([]byte, error) {
	file, ferr := os.Open(path)
	if ferr != nil {
		return nil, ferr
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	content, rerr := ioutil.ReadAll(reader)
	if rerr != nil {
		return nil, rerr
	}
	return content, nil
}

func readAsBase64(path string) (string, error) {
	content, err := ReadAll(path)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

func createSecretDecl(path string, namespace string, name string, key string) (string, error) {
	base64, err := readAsBase64(path)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf(TLS_SECRET_NO_PASSWORD_YAML_TEMPLATE,
		name, namespace, key, base64)
	return text, nil
}

func createSecretPassDecl(path string, namespace string, name string, key string, password string) (string, error) {
	base64, err := readAsBase64(path)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf(TLS_SECRET_PASSWORD_YAML_TEMPLATE,
		name, namespace, key, base64, password)
	return text, nil
}

func verifySecretFields(t *testing.T, namespaceName string, secretName string, fields ...string) {
	secret := GetSecret(t, namespaceName, secretName)
	for _, field := range fields {
		_, ok := secret.Data[field]
		assert.Check(t, ok)
	}
}

func CreateSecret(t *testing.T, namespaceName string, certName string, secretName string, keyDir string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	if keyDir == "" {
		keyDir = filepath.Join("..", "..", "keys")
	}
	keyFile := filepath.Join(keyDir, certName)

	secretString, err := createSecretDecl(keyFile, namespaceName, secretName, certName)
	assert.NilError(t, err)

	k8s.KubectlApplyFromString(t, kubectlOptions, secretString)
	AddTeardown(TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromString(t, kubectlOptions, secretString) })

	fields := []string{certName}
	verifySecretFields(t, namespaceName, secretName, fields...)
}

func CreateSecretWithPassword(t *testing.T, namespaceName string, certName string, secretName string, password string, keyDir string) {
	kubectlOptions := k8s.NewKubectlOptions("", "")
	kubectlOptions.Namespace = namespaceName

	if keyDir == "" {
		keyDir = filepath.Join("..", "..", "keys")
	}
	keyFile := filepath.Join(keyDir, certName)

	secretString, err := createSecretPassDecl(keyFile, namespaceName, secretName, certName, password)
	assert.NilError(t, err)

	k8s.KubectlApplyFromString(t, kubectlOptions, secretString)
	AddTeardown(TEARDOWN_SECRETS, func() { k8s.KubectlDeleteFromString(t, kubectlOptions, secretString) })

	fields := []string{certName, "password"}
	verifySecretFields(t, namespaceName, secretName, fields...)
}
