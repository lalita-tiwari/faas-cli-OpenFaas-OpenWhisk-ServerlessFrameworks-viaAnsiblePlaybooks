
# A command line tool to implement the Serverless Frameworks in private clouds

This directory contains code for a CLI (Command Line Interface) tool which is created to automate the implementation
of serverless FaaS frameworks OpenFaaS and openWhisk on top of the kubernetes cluster node.

This tool is tested on macOS and compatible with intel and Apple Silicon.

### Pre-requisite to run this CLI

1) create an instance in your private cloud and allow the below ports
2) For SSH (PORT 22)
  - Egress	IPv4	Any	Any	0.0.0.0/0	-	-	Delete Rule
  - Egress	IPv6	Any	Any	::/0	-	-	Delete Rule
  - Ingress	IPv4	TCP	22 (SSH)	0.0.0.0/0,
3) For kubernetes,docker and frameworks 
 -  Egress	IPv4	Any	Any	0.0.0.0/0	-	-	Delete Rule
 - Egress	IPv6	Any	Any	::/0	-	-	Delete Rule
 -  Ingress	IPv4	TCP	443 (HTTPS)	0.0.0.0/0	-	-	Delete Rule
 -  Ingress	IPv4	TCP	2379 - 2380	0.0.0.0/0	-	Etcd server client API (kube-apiserver, etcd)	Delete Rule
 -  Ingress	IPv4	TCP	5000	0.0.0.0/0	-	Docker private registry	Delete Rule
 -  Ingress	IPv4	TCP	6443	0.0.0.0/0	-	Kubernetes API server (all)	Delete Rule
 -  Ingress	IPv4	TCP	8090	0.0.0.0/0	-	platform agent, master and worker nodes	Delete Rule
 -  Ingress	IPv4	TCP	8091	0.0.0.0/0	-	Platform API Server	Delete Rule
 -  Ingress	IPv4	TCP	10250	0.0.0.0/0	-	Kubelet API	Delete Rule
 -  Ingress	IPv4	TCP	10251	0.0.0.0/0	-	Kubernetes kube-scheduler	Delete Rule
 -  Ingress	IPv4	TCP	10252	0.0.0.0/0	-	Kubernetes kube-controller-manager	Delete Rule
 -  Ingress	IPv4	TCP	10255	0.0.0.0/0	-	Kubernetes kubelet	Delete Rule
 -  Ingress	IPv4	TCP	10257	0.0.0.0/0	-	Kube-controller-manager	Delete Rule
 -  Ingress	IPv4	TCP	10259	0.0.0.0/0	-	Kube-scheduler (Self)	Delete Rule
 -  Ingress	IPv4	TCP	30000 - 32767	0.0.0.0/0	-	Exposed Node Ports	Delete Rule
 -  Ingress	IPv4	UDP	8285	0.0.0.0/0	-	flannel vxlan	Delete Rule
 -  Ingress	IPv4	UDP	8472	0.0.0.0/0	-	Flannel overlay network

### CLI is written in GO language using GoLand

You can install the CLI with a curl utility script, brew or by downloading the binary from the releases page.
In go.mod file define the name of your CLI tool
--> To build,
go build faas-cli(tool name defined in go.mod file)

### To run the cli tool
Once downloaded you'll get the faas-cli command line zip folder.
unzip the folder.

To allow permission run 
--> chmod +x CLI_Project

--> To Run
./CLI_Project

### This CLI is using below ansible playbooks to install and setup the packages on cloud Virtual Machines
1) An ansible playbook for Docker and Kubernetes cluster.
2) An ansible playbook to install ansible 
3) An ansible playbook for NFS Mounting (Network File Sharing) --> install_nfs.yml
4) An ansible playbook to install openFaaS framework


### Accessing OpenFaaS Frameworks

OpenFaaS Framework can be accessed via below link with username admin and password will be provided in the end of running this tool.
http://129.114.25.142:31112

### Creating and executing a function via OpenFaaS

1) cc@team7-vm1-km:~/helm$ faas-cli new --lang python3 hello-test
   Folder: hello-test created.
   It will create a new python3 function using faas-cli with the name of hello-test where the --lang parameter specify the
   language as python3

2) cc@team7-vm1-km:~/helm$ kubectl get svc -n openfaas

To check the gateway cluster IP and port

3) cc@team7-vm1-km:~/helm/hello-test$ ls
   check the files created under the hello-test folder

4) update the gateway ip: port and the repository address of your docker hub in hello-test.yml file
   write your function logic in handler.py

5) faas-cli build -f ./hello-test.yml
   It will build the dockler image of your function.

6) faas-cli deploy -f ./hello-test.yml
   Deploy the image,

7) faas-cli publish -f ./hello-test.yml
   publish the image

8) post and execute the function by passing the arguments to the function
   curl -X POST http://10.110.204.22:8080/function/hello-test --data-binary '{"docs":[{"x":"y","a":"b"}]}'

9) OpenFaaS GUI to invoke the function
   http://129.114.25.12:31112/ui/


### OpenWhisk Installation Steps:

The name of the main function file you create in OpenWhisk should be __main__.py

main method in the function --args takes dict(dictionary) and return dict only.

create requirements.txt, put your package module to import in it

### RUN the below command to create virtual environment, this cmd gives you virtualenv folder containing packaging of environment.

docker run --rm -v "$PWD:/tmp" openwhisk/python3 action bash-c "cd tmp && virtualenv virtualenv
&& source virtualenv/bin/activate && pip install -r requirements.txt"


### Next,put python file and virtual-env folder together and zip it

### create function command
wsk action create owf --kind python:3 owf.zip -i

### create trigger command (To invoke a function you would need to create a trigger)

wsk trigger create MyKafkaTrigger -f /whisk.system/messaging/kafkaFeed -p
brokers "[\"129.114.25.12:9092\"]" -p topic openwhisk -p isJSONData true -i

### create RULE to associate the trigger with function

wsk rule create owt /guest/MyKafkaTrigger owf -i


### steps for the JFrog Artifactory

jfrog rt c rt-server -url https://tiwarilalita.jfrog.io/artifactory 
--user <USERNAME> --password <PASSWORD> --interactive=false


     JFrog Artifactory URL(https://tiwarilalita.jfrog.io/artifactory/)
     JFrog Distribution URL(https://tiwarilalita.jfrog.io/artifactory/distribution/)
     JFrog Xray URL(https://tiwarilalita.jfrog.io/artifactory/xray/)
     JFrog Mission Control URL(https://tiwarilalita.jfrog.io/artifactory/mc/)
     JFrog Pipelines URL(https://tiwarilalita.jfrog.io/artifactory/pipelines/)
