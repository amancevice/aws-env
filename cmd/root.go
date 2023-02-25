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
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/cobra"
)

var SecretId string = os.Getenv("AWS_SECRET")
var Version string = "0.1.0"
var ShowVersion bool
var rootCmd = &cobra.Command{
	Use:   "aws-secretsmanager-env",
	Short: "Export a SecretsManager JSON secret to a Lambda runtime ENV",
	Long:  "Export a SecretsManager JSON secret to a Lambda runtime ENV",
	Args:  args,
	Run:   run,
}

type SecretsManagerGetSecretValueAPI interface {
	GetSecretValue(ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

func GetSecretValue(ctx context.Context, api SecretsManagerGetSecretValueAPI, input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return api.GetSecretValue(ctx, input)
}

func ExportSecret() {
	// Set up SecretsManager client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}
	client := secretsmanager.NewFromConfig(cfg)

	// Get SecretsManager secret
	input := &secretsmanager.GetSecretValueInput{SecretId: aws.String(SecretId)}
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
	for key, val := range secretJson {
		os.Setenv(key, val)
	}
}

func args(cmd *cobra.Command, args []string) error {
	return nil
}

func run(cmd *cobra.Command, args []string) {
	if ShowVersion {
		os.Stdout.WriteString("aws-secretsmanager-env v" + Version + "\n")
	} else {
		fmt.Printf("EXPORT SecretId: %s\n", SecretId)
		ExportSecret()
		syscall.Exec(args[0], args, os.Environ())
	}
}

func init() {
	rootCmd.SetHelpTemplate(`{{.Short}}

To use this executable in Lambda you must set the ENV variables:
  - AWS_SECRET
  - AWS_LAMBDA_EXEC_WRAPPER

Usage:
  {{.Name}} [OPTIONS] ARGS...

Flags:
  -h, --help      help for {{.Name}}
  -v, --version   show version
`)
	rootCmd.SetUsageTemplate(`
Usage:
  {{.Name}} [OPTIONS] ARGS...
`)
	rootCmd.PersistentFlags().BoolVarP(&ShowVersion, "version", "v", false, "show version")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
