/*
* Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
* WSO2 Inc. licenses this file to you under the Apache License,
* Version 2.0 (the "License"); you may not use this file except
* in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied. See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os/exec"
	"bufio"
	"os"
	"strings"
	"errors"
	"github.com/wso2/product-mi-tooling/cmd/utils"
	"path"
	"runtime"
	"log"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

const secretCreateCmdLiteral = "create"
const secretCreateCmdShortDesc = "Encrypt secrets"
const secretCreateCmdLongDesc = "Create secrets based on given arguments"

var file = ""
var algorithm = ""
var inputs = make(map[string]string)
var keystoreInfoFile = utils.GetkeyStoreInfoFileLocation()

var secretCreateCmd = &cobra.Command{
	Use:   secretCreateCmdLiteral,
	Short: secretCreateCmdShortDesc,
	Long:  secretCreateCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ = cmd.Flags().GetString("file")
		algorithm, _ = cmd.Flags().GetString("cipher")
		if validateKeystoreInitialization() {
			handleSecretCmdArguments(args)
		} else {
			fmt.Println("Keystore has not been initialized.")
		}
	},
}

func init() {
	secretCmd.AddCommand(secretCreateCmd)
	secretCreateCmd.SetHelpTemplate(secretCreateCmdLongDesc)
	secretCreateCmd.Flags().StringP("file", "f", "", "from file")
	secretCreateCmd.Flags().StringP("cipher", "c", "", "algorithm")
}

func handleSecretCmdArguments(args []string) {
		// checks for the output type
		if len(args) == 0 {
			// default output will be console
			inputs["secret.output.type"] = "console"
			initSecretInformation()
		} else if len(args) == 1 {
			outputType := args[0]
			inputs["secret.output.type"] = outputType
			initSecretInformation()
		} else {
			fmt.Println("Too many arguments. See the usage guide.")
		}
}

func initSecretInformation() {
	var consoleResult error

	// set encryption algorithm if custom one is used
	if len(strings.TrimSpace(algorithm)) > 0 {
		os.Setenv("secret.encryption.algorithm", algorithm)
	}
	//checks if input mode is file
	if len(strings.TrimSpace(file)) > 0 {
		inputs["secret.input.type"] = "file"
		inputs["secret.input.file"] = utils.NormalizeFilePath(file)
		consoleResult = startConsoleForSecretInfo(false)
	} else {
		inputs["secret.input.type"] = "console"
		consoleResult = startConsoleForSecretInfo(true)
	}
	if consoleResult == nil {
		secretInfoFilePath := path.Join(utils.GetSecurityDirectoryPath(), "secret-info.properties")
		os.Setenv("secret.source", secretInfoFilePath)
		os.Setenv("keystore.source", keystoreInfoFile)
		utils.SetProperties(inputs, secretInfoFilePath)
		execClient()
		os.Remove(secretInfoFilePath)
	} else {
		fmt.Println(consoleResult)
	}
}

func execClient() {
	var stdoutMessage []byte
	// checks if client jar exists
	if len(strings.TrimSpace(getEncryptionClientName())) == 0 {
		errors.New("encryption client could not be found")
	}

	if runtime.GOOS == "windows" {
		windowsCommand := "java -jar " + getEncryptionClientName()
		output := exec.Command("cmd", "/c", windowsCommand)
		stdoutMessage, _ = output.CombinedOutput()
	} else {
		unixCommand := "java -jar " + getEncryptionClientName()
		output := exec.Command("bash", "-c", unixCommand)
		stdoutMessage, _ = output.CombinedOutput()
	}
	fmt.Printf("%s", stdoutMessage)
}

func startConsoleForSecretInfo(isConsoleInput bool) error {
	reader := bufio.NewReader(os.Stdin)

	if isConsoleInput {
		fmt.Printf("Enter plain alias for secret:")
		alias, _ := reader.ReadString('\n')
		inputs["secret.plaintext.alias"] = strings.TrimSpace(alias)

		fmt.Printf("Enter plain text secret:")
		byteSecret, _ := terminal.ReadPassword(int(syscall.Stdin))
		secret := string(byteSecret)
		fmt.Println()

		fmt.Printf("Repeat plain text secret:")
		byteRepeatSecret, _ := terminal.ReadPassword(int(syscall.Stdin))
		repeatSecret := string(byteRepeatSecret)
		fmt.Println()

		if validateSecrets(secret, repeatSecret) {
			inputs["secret.plaintext.secret"] = strings.TrimSpace(secret)
		} else {
			fmt.Println("Entered secret values did not match.")
			startConsoleForSecretInfo(true)
		}
	}

	if utils.IsValidConsoleInput(inputs) {
		return nil
	} else {
		return errors.New("incomplete secret information")
	}
}

func validateKeystoreInitialization() bool {
	if _, err := os.Stat(keystoreInfoFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func getEncryptionClientName() (string) {
	pwd, _ := os.Getwd()
	content, err := os.Open(pwd)
	if err != nil {
		log.Fatal(err)
	}
	files, _ := content.Readdir(-1)
	content.Close()

	var fileName = ""
	for _, file := range files {
		if strings.Contains(file.Name(), "encryption-client") {
			fileName = file.Name()
			break
		}
	}
	return fileName
}

func validateSecrets(secret string, repeatSecret string) bool {
	if secret == repeatSecret {
		return true
	}
	return false
}
