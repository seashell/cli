SHELL = sh

BOX_NAME := seashell.box

default: release

release: clean build package

clean:
	@vagrant destroy -f
	@rm -rf $(BOX_NAME)

build:
	@vagrant up

package:
	@vagrant package --output $(BOX_NAME)


