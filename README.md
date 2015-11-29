# `with-docker-machine`

Run commands in an environment defined by `docker-machine`.

``` console
$ with-docker-machine -h
usage: with-docker-machine COMMAND [... ARGS]

Run COMMAND in an environment defined by docker-machine.

  -machine string
    	docker machine name (default "default")
  -version
    	print version and exit

Arguments:
  COMMAND  the command to run (typically 'docker')
  ARGS     optional arguments to COMMAND
```

You can invoke single commands or create a shell configured to interact with a Docker machine:

``` console
$ with-docker-machine docker images
$ exec with-docker-machine /bin/bash
```

This simplifies interaction with the [Docker Toolbox][], which does not support single-command invocation and depends on your shell:

``` console
$ eval $(docker-machine env default)
$ docker images
```

## Install

Grab the [latest release][].

[docker toolbox]: https://www.docker.com/docker-toolbox
[latest release]: https://github.com/whilp/with-docker-machine/releases/latest
