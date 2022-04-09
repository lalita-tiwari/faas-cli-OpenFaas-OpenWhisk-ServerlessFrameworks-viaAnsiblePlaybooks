package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var remote_host string

func main() {
	fmt.Println("************************ WELCOME TO FaaS CLI Tool ***********************************")
	fmt.Println("************************ DEVELOPER--> Lalita Tiwari ***********************************")
	fmt.Println("************************ TESTED FOR--> NIX BASED SYSTEMS ***********************************")
	askForK8sChoice()
}

func askForConfirmation(s string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(s + " : [y/n]")

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return "Yes"
		} else if response == "n" || response == "no" {
			return "No"
		}
	}

}
func runCommand(s string) string {
	out, err := exec.Command("bash", "-c", s).Output()
	if err != nil {
		fmt.Printf("%s", err)
		if s == "ansible --version" {
			installAnsible()
		}
	}
	output := string(out[:])
	return output
}
func askForFaaSChoice() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Which FaaS platform would you like to deploy?")
	fmt.Println("1. OpenFaaS")
	fmt.Println("2. OpenWhisk")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response
}
func askForK8sChoice() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Do you have an existing Kubernetes cluster? [y/n]:  ")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" {
		fmt.Println("Please provide Remote IP: ")
		response, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		remote_host = strings.ToLower(strings.TrimSpace(response))
		installAnsible()
		installFrameworks()

	} else {
		fmt.Println("Do you want the tool to create single node K8 managed cluster? [y/n]:  ")
		response, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" {
			fmt.Println("Please provide HOST IP: ")
			response, err = reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			remote_host = strings.ToLower(strings.TrimSpace(response))
			fmt.Println("Installing Ansible on host")
			installAnsible()
			fmt.Println("Installing docker and Kubernetes this will take a while...... ")
			fmt.Println(runCommand("ansible-playbook playbook_install_docker_kubernetes_on_cloud.yml -i " + remote_host + ", -e 'hostname=" + remote_host + "'"))
			fmt.Println("docker and Kubernetes are installed")
			time.Sleep(10)
			fmt.Println("Configuring Kubernetes cluster, it will take a while...... ")
			fmt.Println(runCommand("ansible-playbook playbook_create_K8_cluster.yml -i " + remote_host + ", -e 'master=" + remote_host + "'"))

			installFrameworks()

		}
	}
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	response = strings.ToLower(strings.TrimSpace(response))

	return response
}
func askUserCredentials() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("User: ")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	user := strings.ToLower(strings.TrimSpace(response))
	strings.ToLower(strings.TrimSpace(user))
	fmt.Println("Password: ")
	response, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	//password = strings.ToLower(strings.TrimSpace(response))
}
func installNFS() {
	fmt.Println("Check if ansible is installed ..... ")
	response := runCommand("ansible --version")
	if response != "" {
		fmt.Println("Ansible Found")
	}
	fmt.Println(runCommand("ansible-playbook playbook_install_nfs.yml -i " + remote_host + ","))
}
func installAnsible() {

	fmt.Println("Ansible not found but is needed for installing FaaS frameworks")
	log := askForConfirmation("Please consent to installing ansible")
	if log == "Yes" {
		fmt.Println("Installing Ansible")
		runCommand("pip3 install --user ansible")
		runCommand("pip3 install ansible")
		runCommand("sudo cp `python3 -c \"import sys; print(sys.path)\"| cut -d ',' -f 5 | cut -d ''\\' -f 2| awk -F 'lib' '{print $1}'`bin/* /usr/local/bin/")
		runCommand("sudo chmod -R 777 /usr/local/bin/")
		println("Ansible Installed ... ")
		installNFS()
	}
}

func installFrameworks() {

	val := askForFaaSChoice()
	if val == "1" {
		fmt.Println("Prepare for Installing OpenFaaS ....")

		fmt.Println(runCommand("ansible-playbook playbook_prep_openfass_install.yml -i " + remote_host + ","))
		fmt.Println("Installing OpenFaaS Framework will take some time ....")
		fmt.Println(runCommand("ansible-playbook playbook_install_openfass.yml -i " + remote_host + ","))

		fmt.Println("OpenFaaS serverless Framework is installed and ready to deploy the functions")
		//'PASSWORD=$(kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode) && echo "OpenFaaS admin password: $PASSWORD"'
		fmt.Println("OpenFaaS Framework can be accessed via below link with user name admin and please copy the password from above printed command")
		fmt.Println("http://129.114.25.142:31112")
	} else {
		fmt.Println("Installing OpenWhisk Framework it will take some time ....")
		fmt.Println(runCommand("ansible-playbook playbook_install_openwhisk.yml -i " + remote_host + ","))
		fmt.Println("OpenWhisk Framework is installed")
	}

}
