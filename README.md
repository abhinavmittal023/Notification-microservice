# Notifications Microservice

## Make a database and write it's configuration in config.json file as provided in config.sample file

## Make a config.json inside configuration folder using syntax given in config.sample file.

## Dependencies (what to install)

```
go get -u github.com/dgrijalva/jwt-go
```

```
go get -u github.com/gin-gonic/gin
```

```
go get -u github.com/jinzhu/gorm
```

```
go get -u github.com/lib/pq
```

```
go get -u github.com/pkg/errors
```

```
go get -u github.com/pressly/goose
```

## MailCatcher is used as a fake SMTP server:
### Step 1: Install MailCatcher and its dependencies

```
apt install build-essential libsqlite3-dev ruby-dev
```

```
gem install mailcatcher
```

### Step 2: Running MailCatcher

run `make run_mail_catcher` to start the smtp server and go to the directed url on browser

## Building the server

run `make setup` for setting up the server for the first time.

## Running the server

run `make run` for running the server

## Running the unit tests

run `make run_unit_tests` for running all unit tests
