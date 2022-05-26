VERSION 0.6

all-docker:
    BUILD ./services/ping+docker
    BUILD ./services/query+docker

all-tidy:
    BUILD ./services/ping+tidy
    BUILD ./services/query+tidy
    BUILD ./libs/minecraft+tidy
    BUILD ./libs/database+tidy

dev-up:
    LOCALLY
    RUN docker-compose up --force-recreate

dev-down:
    LOCALLY
    RUN docker-compose down