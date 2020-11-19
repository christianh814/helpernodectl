package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var ImageName string
var ImageVersion string
var containerRuntime string = "podman"
type Config struct {
	Services []string `yaml:"disableservice"`
}

//going to covert this to use the podman module in the future
func PullImage(image string, version string) {

	fmt.Println("Pulling: " + image)
	cmd, err := exec.Command(containerRuntime, "pull", image+":"+version).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

//going to covert this to use the podman module in the future
func StartImage(image string, version string, encodedyaml string, containername string) {

	fmt.Println("Running: " + image)
	cmd, err := exec.Command(containerRuntime, "run", "--rm", "-d", "--env=HELPERPOD_CONFIG_YAML="+encodedyaml, "--net=host", "--name=helpernode-"+containername, image+":"+version).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

}

//going to covert this to use the podman module in the future
func StopImage(containername string) {

	fmt.Println("Stopping: helpernode-" + containername)
	// First, stop container
	exec.Command(containerRuntime, "stop", "helpernode-"+containername).Output()
	// Then, rm the container so we can reuse the name afterwards
	exec.Command(containerRuntime, "rm", "--force", "helpernode-"+containername).Output()
}

//check if an image is running. Return true if it is
func IsImageRunning(containername string) bool {

	// output of all of all running containers
	cmd, err := exec.Command("podman", "ps", "--format", "{{.Names}}").Output()

	// check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on "\n" (space)
	s := strings.Split(strings.TrimSuffix(string(cmd), "\n"), "\n")
	_, found := Find(s, containername)
	return found
}

// checking if service is running
func IsServiceRunning(servicename string) bool {
	// check if the service is active
	activestate, err := exec.Command("systemctl", "show", "-p", "ActiveState", servicename).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", activestate, err)
		os.Exit(53)
	}
	// return the status
	as := strings.TrimSuffix(strings.Split(string(activestate), "=")[1], "\n")
	return as == "active"
}

// checking if service is running
func IsServiceEnabled(servicename string) bool {
	// check if the service is active
	enabledstate, err := exec.Command("systemctl", "is-enabled", servicename).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", enabledstate, err)
		os.Exit(53)
	}
	// return the status
	es := strings.TrimSuffix(string(enabledstate), "\n")
	return es == "enabled"
}

// stopping service
func StopService(servicename string) {

	// stop the service only if it's running
	if IsServiceRunning(servicename) {
		fmt.Println("Stopping service: " + servicename)
		//Stop the service with systemd
		cmd, err := exec.Command("systemctl", "stop", servicename).Output()
		// Check to see if the stop was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// stopping service
func StartService(servicename string) {

	// start the service only if it isn't running
	if !IsServiceRunning(servicename) {
		fmt.Println("Starting service: " + servicename)
		//Start the service with systemd
		cmd, err := exec.Command("systemctl", "start", servicename).Output()
		// Check to see if the start was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// disable service
func DisableService(servicename string) {

	// Disable only if it needs to be
	if IsServiceEnabled(servicename) {
		fmt.Println("Disabling service: " + servicename)
		//Stop the service with systemd
		cmd, err := exec.Command("systemctl", "disable", servicename).Output()
		// Check to see if the stop was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// enable service
func EnableService(servicename string) {

	// Enable only if it needs to be
	if !IsServiceEnabled(servicename) {
		fmt.Println("Enabling service: " + servicename)
		//enable the service with systemd
		cmd, err := exec.Command("systemctl", "enable", servicename).Output()
		// Check to see if the enable was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// get current firewalld rules and return as a slice of string
func GetCurrentFirewallRules() []string {

	// get list of ports currently configured
	cmd, err := exec.Command("firewall-cmd", "--list-ports").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on " " (space)
	s := strings.Split(strings.TrimSuffix(string(cmd), "\n"), " ")

	// get the list of services currenly configured
	scmd, err := exec.Command("firewall-cmd", "--list-services").Output()

	// check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", scmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on " " (space)
	svc := strings.Split(strings.TrimSuffix(string(scmd), "\n"), " ")

	// create a new array based on this new svc array. We will be converting service names to port output
	// simiar to what we got with: firewall-cmd --list--ports
	var ns = []string{}

	// range over the service, find out it's port and append it to the array we just created
	for _, v := range svc {
		lc, err := exec.Command("firewall-cmd", "--service", v, "--get-ports", "--permanent").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", lc, err)
			os.Exit(253)
		}
		nv := strings.TrimSuffix(string(lc), "\n")
		if strings.Contains(nv, " ") {
			ls := strings.Split(nv, " ")
			for _, l := range ls {
				ns = append(ns, l)
			}
		} else {
			ns = append(ns, nv)
		}
	}

	// append this new array of string into the original
	for _, v := range ns {
		s = append(s, v)
	}

	// Let's return this slice of string
	return s
}

func OpenPort(port string) {

	// Open Ports using the port number
	cmd, err := exec.Command("firewall-cmd", "--add-port", port, "--permanent", "-q").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running add-port command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// Reload the firewall to get the most up to date table
	rcmd, err := exec.Command("firewall-cmd", "--reload").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running reload command %s: %s\n", rcmd, err)
		os.Exit(253)
	}
}

func DoISkip(service string) bool {
	/* This fuction takes in the config
	   and checks to see if the service
	   passed in is in the YAML file.
	*/

	var config Config

	//read in the config file
	source, err := ioutil.ReadFile(viper.ConfigFileUsed())

	// check for errors
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
		os.Exit(253)
	}

	// get the contents of the YAML file
	err = yaml.Unmarshal(source, &config)

	// check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling  yaml: %s\n", err)
		os.Exit(253)
	}

	// return if the service passed in is in the YAML
	_, found := Find(config.Services, service)
	return found
}
//
