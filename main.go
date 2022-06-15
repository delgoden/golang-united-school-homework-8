package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/delgoden/golang-united-school-homework-8/service"
	"io"
	"os"
)

type Arguments map[string]string

var (
	flagOperationMissing     = errors.New("-operation flag has to be specified")
	operationCanNotBeHandled = "Operation %s not allowed!"
	flagFileNameMissing      = errors.New("-fileName flag has to be specified")
	flagItemNotProvided      = errors.New("-item flag has to be specified")
	flagIDNotProvided        = errors.New("-id flag has to be specified")

	operations = [4]string{"list", "add", "findById", "remove"}
)

func parseArgs() Arguments {
	var (
		operation string
		item      string
		id        string
		fileName  string
	)

	args := make(Arguments)

	flag.StringVar(&operation, "operation", "", "")
	flag.StringVar(&item, "item", "", "")
	flag.StringVar(&id, "id", "", "")
	flag.StringVar(&fileName, "fileName", "", "")

	flag.Parse()

	args["operation"] = operation
	args["fileName"] = fileName

	if item != "" {
		args["item"] = item
	}

	if id != "" {
		args["id"] = id
	}

	return args
}

func Perform(args Arguments, writer io.Writer) (err error) {

	if args["operation"] == "" {
		return flagOperationMissing
	}
	if args["fileName"] == "" {
		return flagFileNameMissing
	}

	for i, operation := range operations {
		if args["operation"] == operation {
			break
		}
		if i == len(operations)-1 {
			return errors.New(fmt.Sprintf(operationCanNotBeHandled, args["operation"]))
		}
	}

	switch args["operation"] {
	case "list":
		data, err := service.List(args["fileName"])
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(writer, data)
		if err != nil {
			return err
		}
	case "add":
		if args["item"] == "" {
			return flagItemNotProvided
		}
		err = service.Add(args["item"], args["fileName"], writer)
	case "findById":
		if args["id"] == "" {
			return flagIDNotProvided
		}
		err = service.FindById(args["id"], args["fileName"], writer)
	case "remove":
		if args["id"] == "" {
			return flagIDNotProvided
		}
		err = service.Remove(args["id"], args["fileName"], writer)
	}

	return err
}

func main() {

	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
