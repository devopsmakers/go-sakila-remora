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
