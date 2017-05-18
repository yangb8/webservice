# webservice

Welcome to webservice.

Ensure env variables GOPATH and GOROOT are set properly, and this folder is at the following location:
`${GOPATH}/github.com/yangb8`

## Getting Started with webservice

### Requirements

* [Golang](https://golang.org/dl/) 1.8
* [glide](https://github.com/Masterminds/glide)

### Init Project
To get running with webservice and also install the
dependencies, run the following command:

```
make deps
```

### Build

To build the binary for webservice run the command below. This will generate a binary
in the same directory with the name webservice.

```
make
```

### Run

To run webservice:

```
modify config.yml based on your env
./webservice
```


### Test

To test webservice, run the following command:

```
make test
```

### Cleanup

```
make clean
```

## Docker

after build is done

```make docker-image```


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
