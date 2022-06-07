# sailor
![tests](https://github.com/bjornaer/sailor/actions/workflows/push.yaml/badge.svg)[![Go Report Card](https://goreportcard.com/badge/github.com/bjornaer/sailor)](https://goreportcard.com/report/github.com/bjornaer/sailor)

*Ahoy matey!* Welcome to sailor, your very own Port Domain Service!

This is a service which runs a very rudimentary Gin Gonic server, whose sole purpose is to process a data file with the world's ports information and store it nicely in a DB.

### Run 

To run this with a "redis" instance as a DB first setup a redis locally (however you prefer) set the ENV variable `REDIS_ADDR` and then modify the code in `main.go` to initialize DB with "redis"

To run this using a map in memory simply follow the below instruction:



To run this _bad boy_ simply use the [Makefile](./Makefile)'s command

```sh
make <YOUR_OS>
```

and then call on the binary!

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

### API

This service runs a simple REST api. The endpoints in place are to POST a Port, GET a Port, and execute a processing of a Port Data File.

The data file can be passed in as an ENV var under the name `PORTS_FILE`.
