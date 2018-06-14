# Simple file writer

A small Go application, mainly for demo purposes, that takes an input string from parameters, the standard input, or environment variables, and writes it into a file. After this, it waits for a `SIGINT`, `SIGTERM` or `SIGHUP` signal to exit.

```
Simple file writer
------------------

Usage:

  ./write [TARGET] [DATA]

      Writes DATA to the TARGET file.
      By default, it creates any intermediate directories
      necessary, and waits until interrupted.
      DATA can be '-' to read it from the standard input.

Options as environment variables:

  TARGET      The target file
  DATA        The data to write
  TRIM        Trim spaces at the start and end of DATA
  ECHO        Echo the DATA written
  NO_CREATE   Do not create new directories
  NO_READY    Do not create a ready state file
  NO_WAIT     Exit once the file is written
```

Run it as:

```shell
$ docker run --rm -it -v target:/var/static \
     -e DATA='Some content'                 \
     -e TARGET='/var/static/output.txt'     \
     -e NO_WAIT=1 -e NO_READY=1 -e ECHO=1   \
     rycus86/write
```

Or alternatively:

```shell
$ docker run --rm -it -v target:/var/static \
     -e NO_WAIT=1 -e NO_READY=1 -e ECHO=1   \
     rycus86/write /var/static/output.txt 'Some content'
```

Or alternatively:

```shell
$ echo 'Some content' |                       \
    docker run --rm -i -v target:/var/static  \
       -e NO_WAIT=1 -e NO_READY=1 -e ECHO=1   \
       rycus86/write /var/static/output.txt -
# note the missing `-t`, plus the `-` at the end
```

The Docker container comes with health-checking enabled, that will report healthy once the file is successfully written. You can use this in Compose or Swarm stacks to write the file once, then just hang around not doing anything.
