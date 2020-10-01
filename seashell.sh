#!/bin/sh

function log {
    echo "$@"
}

function log_info {
    log "[INFO] ==> $@"
}

function log_error {
    log "[ERROR] ==> $@"
}

VAGRANT=$(which vagrant 2> /dev/null)
if [[  ${VAGRANT} == "" ]]; then
    log_error "Vagrant not found, aborting ..."
    exit 1
fi

set -e


function do_init_env {
    log_info "Initiliazing development environment (this might take several minutes on the first run) ..."
    cd ${VAGRANT_DIR} && \
    vagrant up ${VAGRANT_ID} && \
    vagrant ssh ${VAGRANT_ID} -c \
        "cd /workspace && \
        export NOMAD_ADDR=$NOMAD_ADDR && \
        sudo -E su"
}

function do_destroy_env {
    log_info "initiliazing development environment ..."
    cd ${VAGRANT_DIR} && \
    vagrant destroy -f ${VAGRANT_ID}
    rm -rf ${VAGRANT_FILE}
}

function do_setup_vagrantfile {

    mkdir -p "${VAGRANT_DIR}"

    cat > "${VAGRANT_FILE}" <<EOF
Vagrant.configure("2") do |config|
	config.vm.define "seashell-dev"
	config.vm.hostname = "seashell-dev"
	config.vm.box = "generic/debian10"
	config.vm.synced_folder "${PWD}", "/workspace"

	config.vm.provision "shell", inline: <<-SHELL
        sudo su -
        export GOPATH=$HOME/work
        curl -sS -O https://storage.googleapis.com/golang/go1.15.linux-amd64.tar.gz
        tar -xvf go1.15.linux-amd64.tar.gz
        cp -r go /usr/local
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/.profile
        export PATH=$PATH:/usr/local/go/bin
        git clone git://github.com/hashicorp/levant.git /levant
        cd /levant 
        git checkout --detach f58df42192b5cd98fb83ffec01d9229006a3b29e
        make build
        mv bin/levant /bin/levant
        rm -rf /levant
  	SHELL
  	config.vm.provision "shell", inline: <<-SHELL
        sudo su -
        apt-get update && curl -fsSL get.docker.com | sh >/dev/null
        wget -q https://github.com/docker/buildx/releases/download/v0.4.2/buildx-v0.4.2.linux-amd64  --output-document=docker-buildx 
        mkdir -p /root/.docker/cli-plugins
        mv docker-buildx /root/.docker/cli-plugins
        chmod a+x /root/.docker/cli-plugins/docker-buildx
  	SHELL

end

EOF

}

function do_update_env {
    log_info "Updating development environment (this might take several minutes on the first run) ..."
    cd ${VAGRANT_DIR} && \
    vagrant reload ${VAGRANT_ID}
    vagrant provision ${VAGRANT_ID}
}

function env_command {
    #ENV DEFAULTS
    PROJECT_ID="smart-gateway"
    VAGRANT_DIR="${HOME}/.seashell/${PROJECT_ID}"
    VAGRANT_FILE="${VAGRANT_DIR}/Vagrantfile"
    # TODO: fetch nomad addr from project
    NOMAD_ADDR="${NOMAD_ADDR:-http://35.207.108.145:4646}"

    do_setup_vagrantfile
    
    case $1 in
    init| initialize ) 
        # ARGS: nothing
        shift
        do_init_env "$@"
    ;;
    update ) 
        # ARGS: nothing
        shift
        do_update_env "$@"
    ;;
    destroy ) 
        # ARGS: nothing
        shift
        do_destroy_env "$@"
    ;;
    * ) 
        echo "Invalid option ..." 
    ;;
    esac

}


# GENERAL DEFAULTS
VAGRANT_ID="seashell-dev"

case $1 in
    env | environment ) 
        # ARGS: nothing
        shift        
        env_command "$@"
    ;;
    * ) 
        echo "Invalid option ..." 
    ;;
esac
