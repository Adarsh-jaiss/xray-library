# XRay


## Install 

```bash
go get github.com/thesaas-company/xray@latest
```

## Docs

- [Go Docs](https://pkg.go.dev/github.com/thesaas-company/xray)
- [Example](./example)

## Getting started 

- Run a MySQL Server
```bash
docker run -d --name mysql-employees \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=college \
  genschsa/mysql-employees
```
- Set Password in env variable
```bash
export DB_PASSWORD=college
```
- Run mysql example 
```
go mod tidy
go run example/mysql/main.go
```

## Maintainer
- @Adarsh-jaiss [Adarsh Jaiss]
