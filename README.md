# sailor
![tests](https://github.com/bjornaer/sailor/actions/workflows/push.yaml/badge.svg)[![Go Report Card](https://goreportcard.com/badge/github.com/bjornaer/sailor)](https://goreportcard.com/report/github.com/bjornaer/sailor)

*Ahoy matey!* Welcome to sailor, your very own Port Domain Service!

This is a service which runs a very rudimentary Gin Gonic server, whose sole purpose is to process a data file with the world's ports information and store it nicely in a DB.

---
### Run 

To run this with a *Redis* instance as a DB first setup a redis locally (however you prefer) set the ENV variable `DB_ADDR=<your-redis-address>` and `DB_BACKEND="redis"`


To run this using a map in memory simply follow the below instruction:



To run this _bad boy_ simply use the [Makefile](./Makefile)'s command

```sh
make [darwin|windows|linux]
```

and then invoque the binary!

You could do `go run .cmd/PortDomainService/` *savvy*?

---
### API


This service runs a simple REST api. The endpoints in place are to POST a Port, GET a Port, and execute a processing of a Port Data File.

The data file can be passed in as an ENV var under the name `PORTS_FILE`.

- `POST /port`
- `GET /port/<port-id>`
- `GET /process`
- `GET /`

#### EXAMPLES
- `POST /port`
    ```sh
    curl -d '{"name": "Apia","city": "Apia","country": "Samoa","alias": [],"regions": [],"coordinates": [-171.7513551,-13.8506958],"province": "Upolu","timezone": "Pacific/Apia","unlocs": ["WSAPW"],"code": "77777"}' -H "Content-Type: application/json" -X POST http://localhost:8081/port
    ```
- `GET /port/<port-id>`
    ```sh
    curl http://localhost:8081/port/WSAPW
    ```
- `GET /process`
    ```sh
    curl http://localhost:8081/process
    ```
- `GET /`
    ```sh
    curl http://localhost:8081/
    ```
---
### Docker


You can run the app with the docker compose file provided in this codebase. Please make sure you have Docker setup.

1. pull a redis docker image: `docker pull bitnami/redis:latest`
2. run `docker-compose up`
3. ???
4. profit

---
### Tests

To execute tests please refer to the [Makefile](./Makefile) once again and invoke

```sh
make test
```

You can also make a check for _suspicious constructs_ by calling:

```sh
make vet
```

this generates a report file with any findings.

---
### Development

During the development please be sure to run and add tests, as well as formating the code.

If you don't have your favorite code editor already set up for this, you can format the code by running:

```sh
make fmt
```

To clean up any generated file you can run

```sh
make clean
```
---
###  TODO/Author Notes

I had the original goal of having etcd as the backend DB, but I ran into issues using the client, and wasted some time on a rabbit hole chasing that -- I would like to make it work

- Add more testing, I covered the bare minimum -> added a table test for the router, tested the port file ingestion code
- better error handling, I basically return 500s everytime for any error, messaging could be better and clear
- logging, I really did not spend time adding logging around the code, besides the default logger inside Gin not much useful stuff.
---
