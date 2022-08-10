package libcli

import (
	"fmt"

	nblogger "github.com/banaconda/nb-logger"
	"github.com/eiannone/keyboard"
)

type CommandElem struct {
	Regex      string
	Desc       string
	Func       func(args []string)
	commandMap map[string]*CommandElem
}

type GoCli struct {
	logger       nblogger.Logger
	isRunning    bool
	historyPos   int
	historySlice []string
	commandMap   map[string]*CommandElem
	buf          []byte
	cursorPos    int
}

func NewCommandElem(regex string, desc string, f func(args []string)) *CommandElem {
	return &CommandElem{
		Regex:      regex,
		Desc:       desc,
		Func:       f,
		commandMap: make(map[string]*CommandElem),
	}
}

func NewCommandElemWithoutFunc(regex string, desc string) *CommandElem {
	return &CommandElem{
		Regex:      regex,
		Desc:       desc,
		commandMap: make(map[string]*CommandElem),
	}
}

func (cli *GoCli) Init(logger nblogger.Logger) {
	cli.commandMap = make(map[string]*CommandElem, 0)
	cli.logger = logger
	cli.historySlice = make([]string, 0)
	cli.buf = make([]byte, 0)

	cli.defaultCommand()
}

func (cli *GoCli) AddCommandElem(commandElemSlice ...*CommandElem) {
	commandMap := cli.commandMap
	var commandElem *CommandElem = nil
	for _, elem := range commandElemSlice {
		commandElem = commandMap[elem.Regex]
		if commandElem == nil {
			commandElem = &CommandElem{
				Regex:      elem.Regex,
				Desc:       elem.Desc,
				Func:       elem.Func,
				commandMap: make(map[string]*CommandElem),
			}
			commandMap[elem.Regex] = commandElem
		}
		commandMap = commandElem.commandMap
	}
}

func (cli *GoCli) Run() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	// default command
	cli.isRunning = true

	fmt.Println("Press ESC to quit")
	fmt.Printf("# ")
	for cli.isRunning {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		// normal key input
		if char != 0 && char != '?' {
			cli.inputChar(char)
		}

		if char == '?' {
			cli.printHelp()
			continue
		}

		switch key {
		case keyboard.KeyEsc:
			cli.isRunning = false
		case keyboard.KeySpace:
			cli.inputChar(' ')
		case keyboard.KeyBackspace2:
			cli.backspace()
		case keyboard.KeyTab:
			cli.tabCompletion()
		case keyboard.KeyArrowUp:
			cli.selectHistory(-1)
		case keyboard.KeyArrowDown:
			cli.selectHistory(1)
		case keyboard.KeyArrowLeft:
			cli.moveCursor(-1)
		case keyboard.KeyArrowRight:
			cli.moveCursor(1)
		case keyboard.KeyEnter:
			cli.runCommand()
		}

	}
	fmt.Printf("\n")
}
