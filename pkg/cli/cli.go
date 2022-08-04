package cli

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	nblogger "github.com/banaconda/nb-logger"
	"github.com/eiannone/keyboard"
)

func clearLine() {
	fmt.Printf("\033[2K\r")
}

func clearScreen() {
	fmt.Printf("\033[2J")
}

type CommandElem struct {
	Regex      string
	Desc       string
	Func       func(args []string)
	commandMap map[string]*CommandElem
}

type GoCli struct {
	logger       nblogger.Logger
	isRunning    bool
	historyLen   int
	historyPos   int
	historySlice []string
	commandMap   map[string]*CommandElem
}

func (cli *GoCli) Init(logger nblogger.Logger) {
	cli.commandMap = make(map[string]*CommandElem, 0)
	cli.logger = logger
	cli.historySlice = make([]string, 0)
}

func (cli *GoCli) AddCommand(f func(args []string), desc string, regexs ...string) {
	commandMap := cli.commandMap
	var commandElem *CommandElem = nil
	for _, regex := range regexs {
		commandElem = commandMap[regex]
		if commandElem == nil {
			commandElem = &CommandElem{
				Regex:      regex,
				commandMap: make(map[string]*CommandElem),
			}
			commandMap[regex] = commandElem
		}
		commandMap = commandElem.commandMap
	}
	commandElem.Func = f
	commandElem.Desc = desc
	cli.logger.Info("add command %s %s\n", commandElem.Regex, commandElem.Desc)
}

func (cli *GoCli) getCommandElem(args ...string) *CommandElem {
	commandMap := cli.commandMap
	var commandElem *CommandElem = nil
	for _, arg := range args {
		commandElem = nil
		for regex := range commandMap {
			match, _ := regexp.MatchString(regex, arg)
			if match {
				commandElem = commandMap[regex]
				commandMap = commandElem.commandMap
				break
			}
		}

		if commandElem == nil {
			return nil
		}
	}

	return commandElem
}

func (cli *GoCli) getCommandElemList(args ...string) []*CommandElem {
	commandElemSlice := make([]*CommandElem, 0)
	commandMap := cli.commandMap
	var commandElem *CommandElem = nil
	for _, arg := range args {
		cli.logger.Info("arg=%s", arg)
		commandElem = nil
		for regex := range commandMap {
			cli.logger.Info("regex=%s", regex)
			match, _ := regexp.MatchString(regex, arg)
			if match {
				cli.logger.Info("match")
				commandElem = commandMap[regex]
				commandMap = commandElem.commandMap
				break
			}
		}

		if commandElem == nil {
			cli.logger.Info("commandElem nil")
			for key := range commandMap {
				if !strings.HasPrefix(key, arg) {
					continue
				}
				commandElemSlice = append(commandElemSlice, commandMap[key])
			}

			break
		}
	}

	if commandElem != nil {
		cli.logger.Info("commandElem not nil")
		if len(commandElem.commandMap) == 0 {
			commandElemSlice = append(commandElemSlice, commandElem)
		} else {
			for key := range commandElem.commandMap {
				commandElemSlice = append(commandElemSlice, commandElem.commandMap[key])
			}
		}
	}

	return commandElemSlice
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
	cli.AddCommand(
		func(args []string) {
			fmt.Printf("\n")
			for index, each := range cli.historySlice {
				fmt.Printf(" %d %s\n", index, each)
			}
		},
		"show command history",
		"history")

	cli.AddCommand(
		func(args []string) {
			fmt.Printf("\n")
			for index, each := range cli.historySlice {
				if value, _ := strconv.Atoi(args[1]); value <= index {
					break
				}
				fmt.Printf(" %d %s\n", index, each)
			}
		},
		"show command history in range",
		"history", "[1-9][0-9]*")

	cli.AddCommand(
		func(args []string) {
			fmt.Printf("\n")
			cli.historySlice = make([]string, 0)
		},
		"clear command history",
		"history", "clear")

	cli.AddCommand(
		func(args []string) {
			clearScreen()
		},
		"clear screen",
		"clear")

	cli.AddCommand(
		func(args []string) {
			cli.isRunning = false
		},
		"quit cli",
		"quit")

	cli.AddCommand(
		func(args []string) {
			commandElemSlice := make([]*CommandElem, 0)
			maxLength := 0
			for key := range cli.commandMap {
				commandElem := cli.commandMap[key]
				commandElemSlice = append(commandElemSlice, commandElem)
				if len(commandElem.Regex) > maxLength {
					maxLength = len(commandElem.Regex)
				}
			}
			sort.Slice(commandElemSlice, func(i, j int) bool { return commandElemSlice[i].Regex < commandElemSlice[j].Regex })

			for _, commandElem := range commandElemSlice {
				fmt.Printf("%-*s %s\n", maxLength+1, commandElem.Regex, commandElem.Desc)
			}
		},
		"show help",
		"help")

	fmt.Println("Press ESC to quit")
	fmt.Printf("# ")

	buf := make([]byte, 0)
	for cli.isRunning {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if char != 0 && char != '?' {
			fmt.Printf("%c", char)
			buf = append(buf, byte(char))
		}

		maxLength := 0
		bufLen := len(buf)
		line := strings.Trim(string(buf), " ")
		args := strings.Split(line, " ")
		next := true
		if bufLen > 0 {
			next = buf[bufLen-1] == ' '
		}

		commandElemSlice := cli.getCommandElemList(args...)
		sort.Slice(commandElemSlice, func(i, j int) bool { return commandElemSlice[i].Regex < commandElemSlice[j].Regex })
		for _, commandElem := range commandElemSlice {
			if maxLength < len(commandElem.Regex) {
				maxLength = len(commandElem.Regex)
			}
		}

		if char == '?' {
			fmt.Printf("\n")
			for _, commandElem := range commandElemSlice {
				fmt.Printf("%-*s %s\n", maxLength+1, commandElem.Regex, commandElem.Desc)
			}
			fmt.Printf("# %s", buf)
		}

		switch key {
		case keyboard.KeyEsc:
			cli.isRunning = false
		case keyboard.KeySpace:
			fmt.Printf(" ")
			buf = append(buf, ' ')
		case keyboard.KeyBackspace2:
			if len(buf) > 0 {
				fmt.Printf("\b \b")
				buf = buf[:len(buf)-1]
			}
		case keyboard.KeyTab:
			// case 1: no match							eg. # history eee
			// case 2: incomplete sentece.				eg. # history cle
			// case 3: complete sentece in regex. 		eg. # histroy 10
			// case 4: complete sentence. 				eg. # history or # history clear

			// case 1: no match
			if len(commandElemSlice) == 0 {
				cli.logger.Info("no match\n")
				break
			}

			// case 2: incomplete sentence
			if len(commandElemSlice) == 1 && !next {
				cli.logger.Info("one match")
				clearLine()
				buf = []byte(strings.Join(args[:len(args)-1], " "))
				arg := args[len(args)-1]
				cli.logger.Info("len=%d, args=%v, subargs=%v, arg=%s", len(args), args, args[:len(args)-1], arg)

				// case 2: incomplete sentence
				if strings.HasPrefix(commandElemSlice[0].Regex, arg) {
					cli.logger.Info("incomplete sentence")
					if len(args) > 1 {
						buf = append(buf, ' ')
					}
					buf = append(buf, []byte(commandElemSlice[0].Regex)...)
				} else { // case 3: complete sentece in regex
					cli.logger.Info("complete sentence in regex")
					buf = append(buf, ' ')
					buf = append(buf, []byte(arg)...)
				}
				buf = append(buf, ' ')
				fmt.Printf("# %s", string(buf))
			} else { // case 4: complete sentence
				cli.logger.Info("multi %d", len(commandElemSlice))
				sort.Slice(commandElemSlice, func(i, j int) bool { return commandElemSlice[i].Regex < commandElemSlice[j].Regex })

				fmt.Printf("\n")
				for _, commandElem := range commandElemSlice {
					fmt.Printf("%*s", maxLength+1, commandElem.Regex)
				}
				fmt.Printf("\n# %s", buf)
			}
		case keyboard.KeyArrowUp:
			if cli.historyPos > 0 {
				cli.historyPos--
				clearLine()
				buf = []byte(cli.historySlice[cli.historyPos])
				fmt.Printf("# %s", string(buf))
			}
		case keyboard.KeyArrowDown:
			if cli.historyPos < cli.historyLen-1 {
				cli.historyPos++
				clearLine()
				buf = []byte(cli.historySlice[cli.historyPos])
				fmt.Printf("# %s", string(buf))
			}

		case keyboard.KeyEnter:
			fmt.Printf("\n")
			cli.historyLen = len(cli.historySlice)
			cli.historyPos = cli.historyLen
			if len(line) == 0 {
				buf = make([]byte, 0)
				fmt.Printf("# ")
				break
			}

			if cli.historyLen == 0 || cli.historySlice[cli.historyLen-1] != line {
				cli.historySlice = append(cli.historySlice, line)
				cli.historyLen = len(cli.historySlice)
				cli.historyPos = cli.historyLen
			}

			commandElem := cli.getCommandElem(args...)
			if commandElem != nil && commandElem.Func != nil {
				commandElem.Func(args)
			} else {
				fmt.Printf("command \"%s\" not exist\n", line)
			}
			fmt.Printf("# ")

			buf = make([]byte, 0)
		}
	}
	fmt.Printf("\n")
}
