# Seashell CLI

Teams often have struggle to manually setup and configure local development environments and integrate them to their cloud infrastructure. The Seashell CLI provides them with a tool to:
  * Setup software development tools in isolated local environments
  * Configure, integrate and sync local/remote development environments with their cloud infrastructure deployed via Seashell


## Getting started

- Runtime dependencies:
   - Hashicorp Vagrant 2.2.10+
   - Oracle Virtualbox  6.1.14+
- Build dependencies:
   - GNU Make 4.3+
   - Golang 1.15.2+

- To build and install the CLI, simply clone the repository and run the installation command:
```
~$ git clone https://github.com/seashell/cli.git seashell-cli && \
cd seashell-cli && make install
```
- To verify an installation, run:
```
~$ seashell --version
```

## Features
* Isolated and shareable development environments based on Vagrant + Virtualbox
* Fully customizable and extensible virtual machines as base environments
* Compatible with Linux, Windows & Mac OS
* Infrastructure configuration synchronization between the local/remote environments and the cloud

## Development requirements
* Easy installation
* Standard environment based on projects, independent from developer setup 
* Synchronization with the Seashell Cloud Platform for security and convenience
* Configurable and extensible to adapt to existing workflows, e.g. not enforce CI tools, GIT provider or any specific automation tool, other than those required by the Seashell Cloud.


## Roadmap
- [x] Implement basic CLI structure in golang
- [x] Implement dev command structure and subcommands (init, update, list and destroy)
- [x] Implement dev command vagrant interfaces
- [ ] Implement authentication with the Seashell Cloud platform
- [ ] Implement interactions with the Seashell Cloud platform (fetch information, perform actions) 
- [ ] Configure CI/CD with github actions and submit an official release to github
- [ ] Make the seashell CLI compatible with Windows

