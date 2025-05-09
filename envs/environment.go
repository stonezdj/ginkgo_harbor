package envs

import (
	"fmt"
	"os"
	"strings"

	"github.com/goharbor/tracker/ginkgo_harbor/client"
)

// Environment keeps the testing env info
type HarborEnvironment struct {
	Protocol       string //env var: HTTP_PROTOCOL
	Hostname       string //env var: TESTING_ENV_HOSTNAME
	Account        string //env var: TESTING_ENV_ACCOUNT
	Password       string //env var: TESTING_ENV_PASSWORD
	Admin          string //env var: TESTING_ENV_ADMIN
	AdminPass      string //env var: TESTING_ENV_ADMIN_PASS
	TestingProject string //env var: TESTING_PROJECT_NAME
	GCRProjectName string //env var: TESTING_IMAGE_NAME
	ImageName      string //env var: TESTING_IMAGE_NAME
	ImageTag       string //env var: TESTING_IMAGE_TAG
	ImageDigest    string
	CAFile         string //env var: CA_FILE_PATH
	CertFile       string //env var: CERT_FILE_PATH
	KeyFile        string //env var: KEY_FILE_PATH
	ProxyURL       string //env var: http_proxy, https_proxy, HTTP_PROXY, HTTPS_PROXY

	//API client
	HTTPClient *client.APIClient

	//Docker client
	DockerClient *client.DockerClient

	//Initialize status
	loaded bool
}

func NewHarborEnvironment(hostname string, adminpass string, proxyURL string) *HarborEnvironment {
	env := &HarborEnvironment{
		Protocol:       "https",
		Hostname:       hostname,
		Admin:          "admin",
		AdminPass:      adminpass,
		ImageName:      "nginx",
		ImageTag:       "latest",
		ImageDigest:    "sha256:2f1cd90e00fe2c991e18272bb35d6a8258eeb27785d121aa4cc1ae4235167cfd",
		GCRProjectName: "registry.goharbor.io/harbor-ci/",
		CAFile:         "ca.crt",
		TestingProject: "smoketest",
		Account:        "user001",
		Password:       "Harbor12345",
		ProxyURL:       proxyURL,
	}
	if err := env.Init(); err != nil {
		fmt.Println("Failed to init Harbor environment")
		panic(fmt.Sprintf("Failed to init Harbor environment: %v", err))
	}
	return env
}

func (env *HarborEnvironment) Init() error {
	cfg := client.APIClientConfig{
		Username: env.Admin,
		Password: env.AdminPass,
		CaFile:   env.CAFile,
		CertFile: env.CertFile,
		KeyFile:  env.KeyFile,
		Proxy:    env.ProxyURL,
	}

	httpClient, err := client.NewAPIClient(cfg)
	if err != nil {
		fmt.Println("Failed to create API client")
		return err
	}
	env.HTTPClient = httpClient
	env.DockerClient = &client.DockerClient{}
	return nil
}

// Load test env info
func (env *HarborEnvironment) Load() error {
	host := os.Getenv("TESTING_ENV_HOSTNAME")
	if isNotEmpty(host) {
		env.Hostname = host
	}

	account := os.Getenv("TESTING_ENV_ACCOUNT")
	if isNotEmpty(account) {
		env.Account = account
	}

	pwd := os.Getenv("TESTING_ENV_PASSWORD")
	if isNotEmpty(pwd) {
		env.Password = pwd
	}

	admin := os.Getenv("TESTING_ENV_ADMIN")
	if isNotEmpty(admin) {
		env.Admin = admin
	}

	adminPwd := os.Getenv("TESTING_ENV_ADMIN_PASS")
	if isNotEmpty(adminPwd) {
		env.AdminPass = adminPwd
	}

	pro := os.Getenv("TESTING_PROJECT_NAME")
	if isNotEmpty(pro) {
		env.TestingProject = pro
	}

	imgName := os.Getenv("TESTING_IMAGE_NAME")
	if isNotEmpty(imgName) {
		env.ImageName = imgName
	}

	imgTag := os.Getenv("TESTING_IMAGE_TAG")
	if isNotEmpty(imgTag) {
		env.ImageTag = imgTag
	}

	protocol := os.Getenv("HTTP_PROTOCOL")
	if isNotEmpty(protocol) {
		env.Protocol = protocol
	}

	caFile := os.Getenv("CA_FILE_PATH")
	if isNotEmpty(caFile) {
		env.CAFile = caFile
	}

	keyFile := os.Getenv("KEY_FILE_PATH")
	if isNotEmpty(keyFile) {
		env.KeyFile = keyFile
	}

	certFile := os.Getenv("CERT_FILE_PATH")
	if isNotEmpty(certFile) {
		env.CertFile = certFile
	}

	proxyEnvVar := "https_proxy"
	if env.Protocol == "http" {
		proxyEnvVar = "http_proxy"
	}
	proxyURL := os.Getenv(proxyEnvVar)
	if !isNotEmpty(proxyURL) {
		proxyURL = os.Getenv(strings.ToUpper(proxyEnvVar))
	}
	if isNotEmpty(proxyURL) {
		env.ProxyURL = proxyURL
	}

	if !env.loaded {
		cfg := client.APIClientConfig{
			Username: env.Admin,
			Password: env.AdminPass,
			CaFile:   env.CAFile,
			CertFile: env.CertFile,
			KeyFile:  env.KeyFile,
			Proxy:    env.ProxyURL,
		}

		httpClient, err := client.NewAPIClient(cfg)
		if err != nil {
			return err
		}
		env.HTTPClient = httpClient
		env.DockerClient = &client.DockerClient{}

		env.loaded = true
	}

	return nil
}

// RootURI : The root URI like https://<hostname>
func (env *HarborEnvironment) RootURI() string {
	return fmt.Sprintf("%s://%s", env.Protocol, env.Hostname)
}

func isNotEmpty(str string) bool {
	return len(strings.TrimSpace(str)) > 0
}
