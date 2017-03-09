# go-sakila-remora
A sidecar process written in Go to monitor MySQL slaves / masters for replication errors or lag.

## Welcome!

Hi! and thanks for stopping by this project in it's early days.

Inspired by GitHub's interesting article:
https://githubengineering.com/context-aware-mysql-pools-via-haproxy/

It's similar to methods I've used previously to manage traffic being sent
to Asterisk VoIP servers based on several key metrics.

## Project goals

Ultimately, I'm pretty new to Go and thought this would be an interesting
challenge and an awesome product if everything turns out well.

### Phase 1:

The initial iteration will handle reading YAML configuration, executing some
basic checks against MySQL and presenting a JSON based endpoint for health
checking against. After this iteration I should have a decent BVP - Barely
Viable Product.

### Phase 2:

Let's not get ahead of ourselves.

## Percona master / slave for testing

I needed a really quick way to be able to:
* Spin up a master and slave with replication configured
* Create lag on the slave
* Stop lag on the slave

To do this I've thrown together a really dirty `docker-compose` file and a Makefile
with the following targets:

| Target         | Description                                               |
| -------------- | :-------------------------------------------------------: |
|`make up`       | Brings up the master, slave and sets up replication       |
|`make down`     | Runs `docker-compose down` cleaning everything up         |
|`make start-lag`| Uses `MASTER_DELAY=300` to start the slave lagging        |
|`make stop-lag` | Stops slave lag by setting `MASTER_DELAY=0`               |
|`make slave-lag`| Shows the current value of `Seconds_Behind_Master`        |
|`make stuff`    | Does some SQL stuff - a `CREATE` and `DROP` of a database |

### MySQL Ports
When running the containers you can access MySQL on host: `127.0.0.1` with:
* User: `root`
* Pass: `secret`
* Master Port: `3306`
* Slave Port: `3307`

> Obviously: Please! Please! Please! Don't ever use these settings or config
> In production (or any other application environment). This is for quick and
> dirty testing and development of `go-sakila-remora`.

### Simulating slave lag

To simulate slave lag for testing is super easy:

1. `make up` to bring up a master / slave and configure replication
2. `make start-lag` to start the slave lagging (up to 300s)
3. `make stuff` to do some SQL stuff and things
4. `make slave-lag` a few times to see the lag increasing
5. `make stop-lag` stops the slave lagging
