package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var remote_host string

func main() {
	fmt.Println("************************ WELCOME TO FaaS CLI Tool ***********************************")
	fmt.Println("************************ DEVELOPER--> Lalita Tiwari ***********************************")
	fmt.Println("************************ TESTED FOR--> NIX BASED SYSTEMS ***********************************")
	askForK8sChoice()
}

func main3() {
	fmt.Println("Please provide remote host ip: ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	//ip := runCommand("ip address|grep \"inet 10.\"| awk '{print $2}'|cut -d '/' -f 1")
	runCommand("ansible-playbook playbook_install_nfs.yml -i " + response + ",")
}

func main2() {

	//fmt.Printf("Hi")
	faasChoice := askForFaaSChoice()
	if faasChoice == "1" {
		fmt.Println("you have chosen to install OpenFaaS")
	} else {
		fmt.Println("you have chosen to install OpenWhisk")
	}

	output := runCommand("kubectl config current-context")
	//fmt.Printf(output)
	res := askForConfirmation("Is " + output + " the right Kubernetes context you want to install FaaS?")
	if res == "Yes" {
		response := runCommand("kubectl get nodes -o wide| grep master | awk '{print $6}'")
		fmt.Println("The chosen FaaS will be deployed on " + response)
		fmt.Println("Please enter user credentials for " + response)
		askUserCredentials()
		installNFS()
	}

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
		val := askForFaaSChoice()
		if val == "1" {
			fmt.Println("Installing OpenFaaS it will take some time ....")
			runCommand("ansible-playbook playbook_install_openfass.yml -i " + remote_host + ",")
		} else {
			fmt.Println("Installing OpenWhisk it will take some time ....")
			runCommand("ansible-playbook playbook_install_openfass.yml -i " + remote_host + ",")
		}
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
			fmt.Println("Installing Kubernetes this will take a while...... ")
			runCommand("ansible-playbook playbook_create_K8_cluster.yml -i " + remote_host + ", -e 'master=" + remote_host + "'")
			fmt.Println("Configuring Kubernetes it will take a while...... ")
			runCommand("ansible-playbook playbook_install_docker_kubernetes_on_cloud.yml -i " + remote_host + ", -e 'hostname=" + remote_host + "'")
			val := askForFaaSChoice()
			if val == "1" {
				fmt.Println("Installing OpenFaaS it will take some time ....")
				runCommand("ansible-playbook playbook_install_openfass.yml -i " + remote_host + ",")
			} else {
				fmt.Println("Installing OpenWhisk it will take some time ....")
				runCommand("ansible-playbook playbook_install_openfass.yml -i " + remote_host + ",")
			}
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
	runCommand("ansible-playbook playbook_install_nfs.yml -i " + remote_host + ",")
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
