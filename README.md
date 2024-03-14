# Snippetbox from Let' Go

Developing snippetbox example application following along "Let's Go" book from Alex Edwards

## Run the application

Since we are going to run the database in a container, we need to create a network for the application and the database to communicate.

```bash
dk network create snippetbox
```

Run the database in a container with the following command:

```bash
dk run --rm -d \
    --name snippetbox-mysql \
    -e MYSQL_ROOT_PASSWORD=secret \
    -e MYSQL_DATABASE=snippetbox \
    -p 3306:3306 \
    --user 1000:1000 \
    --network snippetbox \
    -v $PWD/db-init.sql:/docker-entrypoint-initdb.d/db-init.sql \
    -v $PWD/data:/var/lib/mysql \
    mysql:8.3
```

Then we can run the application with the following command:

```bash
dk run -it --rm --network snippetbox -w $PWD -v $PWD:$PWD -p 4000:4000 docker.io/cosmtrek/air air
```
