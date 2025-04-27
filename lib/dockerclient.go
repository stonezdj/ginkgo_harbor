package lib

import (
	"fmt"

	"github.com/goharbor/tracker/ginkgo_harbor/envs"
)

// ConcourseCiSuite : Provides some base cases
type DockerClient struct{}

// PushImage : Push image to the registry
func (dc *DockerClient) PushImage(onEnvironment *envs.HarborEnvironment) error {
	fmt.Println("pushing image")
	docker := onEnvironment.DockerClient
	if err := docker.Status(); err != nil {
		return err
	}

	imagePulling := fmt.Sprintf("%s:%s", onEnvironment.GCRProjectName+onEnvironment.ImageName, onEnvironment.ImageTag)
	if err := docker.Pull(imagePulling); err != nil {
		fmt.Printf("failed to pull image  %v, pull image is %v", err, imagePulling)
		return err
	}

	if err := docker.Login(onEnvironment.Account, onEnvironment.Password, onEnvironment.Hostname); err != nil {
		fmt.Printf("failed to login  %v", err)
		return err
	}

	imagePushing := fmt.Sprintf("%s/%s/%s:%s",
		onEnvironment.Hostname,
		onEnvironment.TestingProject,
		onEnvironment.ImageName,
		onEnvironment.ImageTag)

	fmt.Printf("Pushing image %s to %s\n", imagePulling, imagePushing)

	if err := docker.Tag(imagePulling, imagePushing); err != nil {
		return err
	}

	if err := docker.Push(imagePushing); err != nil {
		return err
	}

	return nil
}

// PullImage : Pull image from registry
func (dc *DockerClient) PullImage(onEnvironment *envs.HarborEnvironment) error {
	docker := onEnvironment.DockerClient
	if err := docker.Status(); err != nil {
		return err
	}

	if err := docker.Login(onEnvironment.Account, onEnvironment.Password, onEnvironment.Hostname); err != nil {
		return err
	}

	imagePulling := fmt.Sprintf("%s/%s/%s:%s",
		onEnvironment.Hostname,
		onEnvironment.TestingProject,
		onEnvironment.ImageName,
		onEnvironment.ImageTag)

	if err := docker.Pull(imagePulling); err != nil {
		return err
	}

	return nil
}
