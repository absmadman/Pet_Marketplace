serv = server
db = database_api
cache = redis-db

HELP_FUNC = \
	%help; while(<>){push@{$$help{$$2//'options'}},[$$1,$$3] \
	if/^([\w-_]+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/}; \
    print"$$_:\n", map"  $$_->[0]".(" "x(20-length($$_->[0])))."$$_->[1]\n",\
    @{$$help{$$_}},"\n" for keys %help; \

all: ##@APP application in docker container
	docker-compose-api

docker-compose-api: ##@APP runs application in docker container
	docker build --no-cache -t $(serv) .
	docker-compose up

clean-data: ##@DB clean a database saved data
	rm -rf pkg/repository/db/pgdata
	rm -rf pkg/repository/redis/data
	rm -rf pkg/repository/redis/redis.conf

docker-stop-api: ##@SERVER stops containers
	docker stop $(db)
	docker stop $(serv)
	docker stop $(cache)

docker-clean-api: docker-stop-api ##@SERVER delete server, database and cache containers
	docker rm $(db)
	docker rm $(serv)
	docker rm $(cache)

server-logs: ##@SERVER show logs from server container
	docker logs $(serv)

database-logs:  ##@DB show logs from database container
	docker logs $(db)

cache-logs: ##@CACHE show logs from cache container
	docker logs $(cache)

all-logs: database-logs server-logs cache-logs ##@APP show logs from server and db containers together

help: ##@APP Show help info
	@echo -e "Usage: make [target] ...\n"
	@perl -e '$(HELP_FUNC)' $(MAKEFILE_LIST)