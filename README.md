# sep-golang
Smart Energy Profile using Go

## Setup

### Host
For testing you will need to ensure your server address is recognized on your system. Simply add it to your localhost for internal testing in the **/etc/hosts** file for linux systems.

```shell
127.0.0.1   egot.internal.com
```

### SSL
I found a super handy tool for setting yourself as a CA and generating tls certificates for clients and servers.

* https://smallstep.com/docs/step-cli/reference/

After you have step and step-ca installed simply issue a ca, server, and client sertificate using the very handy tutorial and you are up and ready for tls.

* https://smallstep.com/hello-mtls/doc/combined/go/go

If you follow the basic example settings for the CA setup you need to modify the default certificate duration to be greater than 24 hours. Use the following code, but verify that the provisioner identity is correct for your installation. 

* https://smallstep.com/docs/step-ca/provisioners/#remote-provisioner-management

```shell
step ca provisioner update you@smallstep.com \
   --x509-min-dur=24h \
   --x509-max-dur=8760h \
   --x509-default-dur=8760h
```

### Golang

#### Import packages
Go allows you to import github directly into your code.

- sep: IEEE std 2030.5 Smart Energy Profile models

Example:

 ```shell
import "github.com/Tylores/egot/sep"
```

#### Install programs

Directly install programs into your go bin directory which can then be executed directly

- client: interface for DER
- crawler: tester to check server for known services
- core: server microservice for main DER registration
- flowreservation: server microservice for FlowReservationRequest/Responses
- operator: Utility server for Grid Service Providers to participate in grid services

Example:

```shell
go install github.com/Tylores/egot/cmd/crawler@latest
crawler
```

#### Run programs

during development it is easier to direct call your programs to test, check the tools folder for utilies to run scenarios


```shell
go mod tidy
go run ./cmd/crawler
```
