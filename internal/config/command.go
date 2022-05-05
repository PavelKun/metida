package config

import (
	"flag"
	"fmt"
	"os"
)

type CommandLineI interface {
	IfDebagOn() bool
	IfPprofOn() bool
	IfLogRequestsOn() bool
}

// CommandLine contains all params from os.Args.
type CommandLine struct {
	debag       bool // debag = true in mode dev
	pprof       bool //
	logRequests bool
}

// NewCommandLine отвечает за настройки запуска данного проекта
// Тут должны задаваться уровни логирования(включения дебаг лога для гина и запа)
func NewCommandLine() (*CommandLine, error) {
	var debag, pprof bool
	var logRequests bool

	if len(os.Args[1:]) < 1 {
		CommandHelp()
		return &CommandLine{}, fmt.Errorf("Не задан режим запуска")
	}

	modeProd := flag.NewFlagSet("prod", flag.ExitOnError)
	modeProd.BoolVar(&pprof, "pprof", false, "profiling mode")
	modeProd.BoolVar(&logRequests, "logRequests", false, "logging all requests")

	modeDev := flag.NewFlagSet("dev", flag.ExitOnError)
	modeDev.BoolVar(&pprof, "pprof", false, "profiling mode")

	switch os.Args[1] {
	case "prod":
		debag = false
		if err := modeProd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Ошибка аргументов")
		}
	case "dev":
		debag = true
		if err := modeProd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Ошибка аргументов") // тут логера пока нету
		}

	case "help":
		CommandHelp()
		return &CommandLine{}, fmt.Errorf("Для дальнейшей работы выберите режим запуска")
	}

	return &CommandLine{debag: debag, pprof: pprof, logRequests: logRequests}, nil
}

func (o *CommandLine) IfDebagOn() bool {
	return o.debag
}

func (o *CommandLine) IfPprofOn() bool {
	return o.pprof
}

func (o *CommandLine) IfLogRequestsOn() bool {
	return o.logRequests
}

func CommandHelp() {
	fmt.Println(`
Перечень основных команд:
        dev: [режим разработки, по умолчанию включает расширенный логгер, и профилировщик]
           flags: -pprof [включение профилирофщика] (true or false)

        prod:
           flags: -pprof [включение профилирофщика] (true or false)
                  -logRequests [включениe логирования всех запросов] (true or false)	

`)
}
