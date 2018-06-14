package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

func main() {
	write()
	waitOrExit()
}

func write() {
	var target, data string

	args := os.Args[1:]

	// parse command line args
	switch len(args) {
	case 2:
		data = args[1]
		fallthrough
	case 1:
		target = args[0]
	}

	// print help if needed
	switch target {
	case "check":
		if contents, err := ioutil.ReadFile("/.ready"); err != nil {
			println("Failed to read the ready state file")
			os.Exit(1)
		} else if string(contents) == "OK" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}

	case "-h", "--help", "help":
		help()
		os.Exit(0)

	}

	// use the TARGET env variable if not given as parameter
	if fromEnv, ok := os.LookupEnv("TARGET"); ok && target == "" {
		target = fromEnv
	}
	if target == "" {
		println("No target given")
		os.Exit(1)
	}

	// use the DATA env variable if not given as parameter
	if fromEnv, ok := os.LookupEnv("DATA"); ok && data == "" {
		data = fromEnv
	}
	if data == "-" {
		if in, err := ioutil.ReadAll(os.Stdin); err != nil {
			println("Failed to read standard input")
			os.Exit(1)
		} else {
			data = string(in)
		}
	}
	if data == "" {
		println("No data given")
		os.Exit(1)
	}

	// create intermediate directories if needed
	if os.Getenv("NO_CREATE") == "" {
		_, err := os.Stat(path.Dir(target))
		if err != nil {
			if p, ok := err.(*os.PathError); !ok {
				println("Unexpected error when checking", p.Path, ":", err.Error())
				os.Exit(1)
			} else if err := os.MkdirAll(p.Path, os.ModePerm); err != nil {
				println("Failed to create directory at", p.Path, ":", err.Error())
				os.Exit(1)
			}
		}
	}

	// trim the spaces if needed
	if os.Getenv("TRIM") != "" {
		data = strings.TrimSpace(data)
	}

	// create/overwrite the target file
	f, err := os.Create(target)
	if err != nil {
		println("Failed to create the target file:", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	// write DATA into the file
	if _, err := f.WriteString(data); err != nil {
		println("Failed to write the target file:", err.Error())
		os.Exit(1)
	}

	if os.Getenv("NO_READY") == "" {
		// create ready state file
		r, err := os.Create("/.ready")
		if err != nil {
			println("Failed to create ready state file")
			os.Exit(1)
		}
		defer r.Close()

		// write ready state
		if _, err := r.WriteString("OK"); err != nil {
			println("Failed to write ready state")
			os.Exit(1)
		}
	}

	// echo back the DATA if needed
	if os.Getenv("ECHO") != "" {
		println("Data:")
		println(data)
	}

	println("Target file written to", target)
}

func waitOrExit() {
	// wait for an exit signal if needed
	if os.Getenv("NO_WAIT") == "" {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-c
	}
}

func help() {
	println(`Simple file writer
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
`)
}
