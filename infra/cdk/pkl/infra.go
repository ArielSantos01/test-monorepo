package pkl

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var awsEnvs map[string]string

type awsProfile struct {
	name   string
	config map[string]string
}

func Pkl(service string, isLocal bool) (map[string]any, error) {
	mainFile := service + "/config/app/app.pkl"

	println("Searching pkl file:", mainFile)

	if !Exists(mainFile) {
		return map[string]any{}, nil
	}

	println("Building pkl:", service)

	var envVars []string
	envs, err := getAWSEnvs(isLocal) // {}
	if err != nil {
		return map[string]any{}, err
	}

	for key, value := range envs {
		envVars = append(envVars, key+"="+value)
	}

	cmd := "pkl eval -f json " + mainFile
	command := exec.Command("bash", "-c", cmd)
	command.Env = append(os.Environ(), envVars...)

	var output []byte
	output, err = command.CombinedOutput()
	if err != nil {
		println("Error ejecutando el comando")
		println("Error: ", err)
		return map[string]any{}, err
	}

	var result map[string]any
	err = json.Unmarshal(output, &result)
	if err != nil {
		println("Error unmarshalling")
		println("Error: ", err)
		return map[string]any{}, err
	}

	return result, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getAWSEnvs(isLocal bool) (map[string]string, error) {
	if !isLocal { //false -> true
		return map[string]string{}, nil // si no es local se toman las variables de entorno de la instancia
	}

	if len(awsEnvs) > 0 {
		return awsEnvs, nil
	}

	profile, err := GetAWSProfileCredentials("draftea-dev") //profile.config
	if err != nil {
		return nil, err
	}

	accountID, err := getAWSAccountID() // solo el id 1223456789
	if err != nil {
		return nil, err
	}

	awsEnvs = map[string]string{
		"AWS_ACCESS_KEY_ID":     profile["aws_access_key_id"],
		"AWS_SECRET_ACCESS_KEY": profile["aws_secret_access_key"],
		"AWS_DEFAULT_REGION":    "us-east-2",
		"AWS_ACCOUNT":           accountID,
		"STAGE":                 "dev",
	}

	return awsEnvs, nil
}

func getAWSAccountID() (string, error) {
	cmd := "aws sts get-caller-identity --query Account --output text --profile draftea-dev"
	command := exec.Command("bash", "-c", cmd)
	output, err := command.CombinedOutput()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetAWSProfileCredentials(profile string) (map[string]string, error) {
	profiles, err := readAWSCredentials(os.ExpandEnv("$HOME/.aws/credentials"))
	if err != nil {
		return nil, err
	}

	for _, p := range profiles {
		if p.name == profile {
			return p.config, nil
		}
	}

	return nil, errors.New("profile not found")
}

func readAWSCredentials(filepath string) ([]awsProfile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var profiles []awsProfile
	var currentProfile *awsProfile

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if currentProfile != nil {
				profiles = append(profiles, *currentProfile)
			}
			profileName := strings.Trim(line, "[]")
			currentProfile = &awsProfile{
				name:   profileName,
				config: make(map[string]string),
			}
		} else if currentProfile != nil {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				currentProfile.config[key] = value
			}
		}
	}

	if currentProfile != nil {
		profiles = append(profiles, *currentProfile)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func getReplace(pklConfig map[string]any) map[string]string {
	result := make(map[string]string)
	for k, v := range pklConfig {
		result["{{pkl:"+k+"}}"] = fmt.Sprintf("%v", v)
		//"{{pkl:TURBO_TICKETS_API_GATEWAY_ID}}": "12313"
	}
	return result
}

func ReadConfig[T any](path string, pklConfig map[string]any) (T, error) {
	var resource T

	replaced := make(map[string]string)
	if pklConfig != nil {
		replaced = getReplace(pklConfig)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return resource, err
	}

	configPath := filepath.Join(workingDir, "..", "cmd", path, "/config.json")

	var configData []byte
	configData, err = os.ReadFile(configPath)
	if err != nil {
		return resource, err
	}

	strConfigData := string(configData)

	for k, v := range replaced {
		strConfigData = strings.ReplaceAll(strConfigData, k, v)
	}

	configData = []byte(strConfigData)

	if err = json.Unmarshal(configData, &resource); err != nil {
		return resource, err
	}
	return resource, nil
}
