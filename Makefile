export GO15VENDOREXPERIMENT=1

exe = github.com/devopsmakers/go-sakila-remora
pkgs = $(shell glide novendor)
cmd = go-sakila-remora

TRAVIS_TAG ?= "0.0.0"
BUILD_DIR=build
COVERAGE_DIR=${BUILD_DIR}/coverage
GODEP=$(GOPATH)/bin/godep

# Build related tasks
.PHONY: all
all: clean deps build

.PHONY: deps
deps:
	go get -u github.com/Masterminds/glide
	glide install

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o $(BUILD_DIR)/$(cmd) $(exe)

.PHONY: run
run: build
	build/go-sakila-remora

.PHONY: clean
clean:
	rm -rfv ./$(BUILD_DIR)

.PHONY: lint
lint: prepare-tests
	golint $(pkgs)

# Test related tasks
.PHONY: test
test: prepare-tests unit-tests coverage-report

.PHONY: prepare-tests
prepare-tests: deps
	mkdir -p ${COVERAGE_DIR}
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/pierrre/gotestcover

.PHONY: unit-tests
unit-tests: prepare-tests
	$(GOPATH)/bin/gotestcover -coverprofile=${COVERAGE_DIR}/unit.cov -short -covermode=atomic $(pkgs)

.PHONY: coverage-report
coverage-report:
	# Writes atomic mode on top of file
	echo 'mode: atomic' > ./${COVERAGE_DIR}/full.cov
	# Collects all coverage files and skips top line with mode
	tail -q -n +2 ./${COVERAGE_DIR}/*.cov >> ./${COVERAGE_DIR}/full.cov
	go tool cover -html=./${COVERAGE_DIR}/full.cov -o ${COVERAGE_DIR}/full.html

# Percona Server related tasks
.PHONY: master
master:
	docker-compose up -d percona_master
	sleep 30
	docker-compose exec -T percona_master 'mysql' -uroot -psecret -hpercona_master -vvv -e'GRANT REPLICATION SLAVE ON *.* TO repl@"%" IDENTIFIED BY "slavepass"\G'
	docker-compose exec -T percona_master 'mysql' -uroot -psecret -hpercona_master -e"SHOW MASTER STATUS\G"

.PHONY: slave
slave:
	docker-compose up -d percona_slave
	sleep 30
	out=`docker-compose exec -T percona_master 'mysql' -uroot -psecret -hpercona_master -e"SHOW MASTER STATUS\G;"` ; \
	file=`grep File <<<"$$out" | cut -f2 -d':' | xargs`  ; \
	position=`grep Position <<<"$$out" | cut -f2 -d':' | xargs` ; \
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"STOP SLAVE;" -vvv ; \
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"change master to master_host=\"percona_master\",master_user=\"repl\",master_password=\"slavepass\",master_log_file=\"$$file\",master_log_pos=$$position;" -vvv
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"START SLAVE;" -vvv

.PHONY: slave-status
slave-status:
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"SHOW SLAVE STATUS\G" -vvv

.PHONY: slave-lag
slave-lag:
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"SHOW SLAVE STATUS\G" -vvv | grep Seconds

.PHONY: start-lag
start-lag:
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"STOP SLAVE;" -vvv
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"CHANGE MASTER TO MASTER_DELAY = 300;" -vvv
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"START SLAVE;" -vvv

.PHONY: stop-lag
stop-lag:
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"STOP SLAVE;" -vvv
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"CHANGE MASTER TO MASTER_DELAY = 0;" -vvv
	docker-compose exec -T percona_slave 'mysql' -uroot -psecret -hpercona_slave -e"START SLAVE;" -vvv

.PHONY: create
create:
	docker-compose exec -T percona_master 'mysql' -uroot -psecret -hpercona_master -e"CREATE DATABASE lag_test;"

.PHONY: drop
drop:
	docker-compose exec -T percona_master 'mysql' -uroot -psecret -hpercona_master -e"DROP DATABASE lag_test;"

.PHONY: stuff
stuff: create drop

.PHONY: down
down:
	docker-compose down

.PHONY: up
up: down master slave
