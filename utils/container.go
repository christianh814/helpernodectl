package utils

import (
	"fmt"
	"os/exec"
)

var ImageName string
var ImageVersion string
//setting this up for future docker support possibly
var containerRuntime string = "podman"

//going to covert this to use the podman module in the future
func PullImage(image string, version string){

	fmt.Println("Pulling: " + image)
	//TODO Need to write the output for the image pull
	cmd, err := exec.Command(containerRuntime, "pull", image + ":" + version).Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println(err)
	}

}

//going to covert this to use the podman module in the future
func StartImage(image string, version string, encodedyaml string, containername string){

	fmt.Println("Running: " + image)
	/* TODO:
		- Need to write the output for the image run
		- Check if the image is already running
	*/
	cmd, err := exec.Command(containerRuntime, "run", "-d", "--env=HELPERPOD_CONFIG_YAML=" + encodedyaml, "--net=host", "--name=helpernode-" + containername, image + ":" + version).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(cmd)
	}

}

//going to covert this to use the podman module in the future
func StopImage(containername string){

	fmt.Println("Stopping: helpernode-" + containername)
	/* TODO:
		- Need to write the output for the image run
		- Check if service is already stopped
	*/
	// First, stop container
	exec.Command(containerRuntime, "stop", "helpernode-" + containername).Output()
	// Then, rm the container so we can reuse the name afterwards
	rmcmd, err := exec.Command(containerRuntime, "rm", "--force", "helpernode-" + containername).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(rmcmd)
	}

}
