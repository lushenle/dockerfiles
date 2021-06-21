## container ops

- build

```bash
$ GO111MODULE=on go build -o cops
```

- get all contaiers

```bash
# localhost
$ cops
# remote docker host
$ cops -H <DockerHost> -P <REST API PORT>
```

- stop container

```bash
# stop a container
$ cops  -name <ContainerName> -ops stop
# or
$ cops -H <DockerHost> -P <REST API PORT> -name <ContainerName> -ops stop

# stop all container
$ cops -s all
# or
$ cops -H <DockerHost> -P <REST API PORT> -s all
```

- remove contaier

```bash
# remove a container
$ cops  -name <ContainerName> -ops remove
# or
$ cops -H <DockerHost> -P <REST API PORT> -name <ContainerName> -ops remove

# stop all container
$ cops -r all
# or
$ cops -H <DockerHost> -P <REST API PORT> -r all
```

