.PHONY: all run-app build-all

$(eval BRANCH := $(shell git branch --show-current))
$(eval WORK_DIR := /root/bin/hotbox-adm-backend)

all: run-app


.pulled:
	git pull origin $(BRANCH)

build-all:
	go build -o $(WORK_DIR)/app/cmd

ifeq ($(BRANCH), main)
	go build -o $(WORK_DIR)/app/cmd
endif

run-app: .pulled
ifeq ($(BRANCH), qa)
	swag init
	go build -o $(WORK_DIR)/app/cmd
	supervisorctl restart hotbox-adm-backend
endif
ifeq ($(BRANCH), main)
	go build -o $(WORK_DIR)/app/cmd
	supervisorctl restart hotbox-adm-backend
endif

lint:
	gofumpt -l -w .

doc:
	swag init

push:
	git push origin $(BRANCH)

pull:
	git pull origin $(BRANCH)

push-f:
	git push origin $(BRANCH) -f
