/*
Copyright © 2023 Alexander Mancevice <alexander.mancevice@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var CliName string = "aws-env"
var CliVersion string = "0.1.0"
var DefaultConfigPath = "/var/task/.aws"
var ShowVersion bool

var UsageTemplate string = `
Usage:
  {{.Name}} [OPTIONS] [ARGS...]
`
var HelpTemplate string = `
{{.Short}}

Usage:
  {{.Name}} [OPTIONS] [ARGS...]

Options:
  -h, --help     show help
  -v, --version  show version

`

var rootCmd = &cobra.Command{
	Use:   CliName,
	Short: "Export ENV variables from AWS ParameterStore & SecretsManager",
	Long:  "Export ENV variables from AWS ParameterStore & SecretsManager",
	Args:  args,
	Run:   run,
}

type ConfigObject struct {
	Exports []struct {
		Secretsmanager string
		Ssm            string
	}
}

type LogWriter struct {
}

type SecretsManagerGetSecretValueApi interface {
	GetSecretValue(ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type SsmGetParametersByPathApi interface {
	GetParametersByPath(ctx context.Context,
		params *ssm.GetParametersByPathInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
}

func (writer LogWriter) Write(bytes []byte) (int, error) {
	return fmt.Fprint(os.Stderr, "ENV "+string(bytes))
}

func GetParametersByPath(ctx context.Context, api SsmGetParametersByPathApi, input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	return api.GetParametersByPath(ctx, input)
}

func GetSecretValue(ctx context.Context, api SecretsManagerGetSecretValueApi, input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return api.GetSecretValue(ctx, input)
}

func GetAwsConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}
	return cfg
}

func GetConfig() ConfigObject {
	// Get config path from $AWS_ENV_CONFIG or use default
	configPath := os.Getenv("AWS_ENV_CONFIG")
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	// Load & return config if it exists
	config := ConfigObject{}
	_, error := os.Stat(configPath)
	if !errors.Is(error, os.ErrNotExist) {
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("unable to read config: %v", err)
		}
		yaml.Unmarshal([]byte(data), &config)
	}
	return config
}

func ExportParameters(path string) {
	// Set up SSM client
	client := ssm.NewFromConfig(GetAwsConfig())

	// Get params
	log.Printf("Ssm:GetParametersByPath Path: %s WithDecryption: true", path)
	input := &ssm.GetParametersByPathInput{Path: aws.String(path), WithDecryption: aws.Bool(true)}
	result, err := GetParametersByPath(context.TODO(), client, input)
	if err != nil {
		log.Fatalf("unable to get parameters: %v", err)
	}

	// Export params
	for _, param := range result.Parameters {
		parts := strings.Split(*param.Name, "/")
		key := parts[len(parts)-1]
		val := param.Value
		ExportVar(key, *val)
	}
}

func ExportSecret(secretId string) {
	// Set up SecretsManager client
	client := secretsmanager.NewFromConfig(GetAwsConfig())

	// Get SecretsManager secret
	log.Printf("SecretsManager:GetSecretValue SecretId: %s", secretId)
	input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretId)}
	result, err := GetSecretValue(context.TODO(), client, input)
	if err != nil {
		log.Fatalf("unable to get secret: %v", err)
	}

	// Parse SecretsManager secret JSON
	var secretJson map[string]string
	secretString := *result.SecretString
	err = json.Unmarshal([]byte(secretString), &secretJson)
	if err != nil {
		log.Fatalf("unable to parse secret JSON: %v", err)
	}

	// Export secret to ENV
	keys := make([]string, 0, len(secretJson))
	for key := range secretJson {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ExportVar(key, secretJson[key])
	}
}

func ExportVar(key string, val string) {
	if os.Getenv(key) == "" {
		// log.Printf("export %s", key)
		os.Setenv(key, val)
		// } else {
		// 	log.Printf("export %s [already exported]", key)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpTemplate(HelpTemplate)
	rootCmd.SetUsageTemplate(UsageTemplate)
	rootCmd.PersistentFlags().BoolVarP(&ShowVersion, "version", "v", false, "show version")
}

func args(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New(CliName + " requires at least 1 argument")
	} else if len(args) >= 1 && args[0][0:1] != "/" {
		return errors.New(CliName + " first arg must be an absolute path")
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	if ShowVersion {
		os.Stdout.WriteString(CliName + " v" + CliVersion + "\n")
		return
	}

	// Set up logger
	log.SetFlags(0)
	log.SetOutput(new(LogWriter))

	// Export from config
	config := GetConfig()
	for _, resource := range config.Exports {
		if resource.Secretsmanager != "" {
			secretId := resource.Secretsmanager
			ExportSecret(secretId)
		}
		if resource.Ssm != "" {
			path := resource.Ssm
			ExportParameters(path)
		}
	}

	// Export from ENV
	exports := strings.Split(os.Getenv("AWS_ENV_EXPORT"), ",")
	secretsmanagerPrefix := "secretsmanager://"
	ssmPrefix := "ssm://"
	for _, export := range exports {
		if strings.HasPrefix(export, secretsmanagerPrefix) {
			secretId := strings.TrimSuffix(strings.TrimPrefix(export, secretsmanagerPrefix), "/")
			ExportSecret(secretId)
		} else if strings.HasPrefix(export, ssmPrefix) {
			path := "/" + strings.TrimPrefix(export, ssmPrefix)
			ExportParameters(path)
		}
	}

	syscall.Exec(args[0], args, os.Environ())
}
