package libcli

import (
	"fmt"
	"go-cli/pkg/libutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// get history length
func (cli *GoCli) getHistoryLen() int {
	return len(cli.historySlice)
}

// get buf length of GoCli
func (cli *GoCli) getBufLen() int {
	return len(cli.buf)
}

// get line string of GoCli with trim
func (cli *GoCli) getLineString() string {
	return strings.TrimSpace(string(cli.buf))
}

// get args by line string
func (cli *GoCli) getArgs() []string {
	args := strings.Split(cli.getLineString(), " ")

	// remove empty string
	args = libutil.RemoveEmptyString(args)

	return args
}

// get is empty of last char of GoCli
func (cli *GoCli) getIsLastCharEmpty() bool {
	return cli.getBufLen() == 0 || cli.buf[cli.getBufLen()-1] == ' '
}

// move cursor by direction
func (cli *GoCli) moveCursor(direction int) {
	if cli.cursorPos == 0 && direction < 0 || cli.cursorPos >= cli.getBufLen() && direction > 0 {
		return
	}
	libutil.MoveCursor(direction)

	// update cursor pos
	cli.cursorPos += direction
}

// append rune to GoCli buf
func (cli *GoCli) inputChar(r rune) {
	cli.logger.Info("cli cursor pos: %d, buf len %d", cli.cursorPos, cli.getBufLen())

	if cli.cursorPos == cli.getBufLen() {
		cli.buf = append(cli.buf, byte(r))
	} else {
		cli.buf = append(cli.buf, 0)
		copy(cli.buf[cli.cursorPos+1:], cli.buf[cli.cursorPos:])
		cli.buf[cli.cursorPos] = byte(r)
	}
	libutil.ClearLine()
	fmt.Printf("# %s", cli.buf)
	cli.cursorPos++

	libutil.MoveCursor(cli.cursorPos - cli.getBufLen())
}

// backspace
func (cli *GoCli) backspace() {
	if cli.getBufLen() == 0 {
		return
	}

	if cli.cursorPos == cli.getBufLen() {
		cli.buf = cli.buf[:cli.getBufLen()-1]
	} else {
		copy(cli.buf[cli.cursorPos-1:], cli.buf[cli.cursorPos:])
		cli.buf = cli.buf[:cli.getBufLen()-1]
	}

	libutil.ClearLine()
	fmt.Printf("# %s", cli.buf)
	cli.cursorPos--

	libutil.MoveCursor(cli.cursorPos - cli.getBufLen())
}

func (cli *GoCli) printHelp() {
	commandElemSlice := cli.getCommandElemList(cli.getArgs()...)

	fmt.Printf("\n")

	helpStringSlice := make([]string, 0)
	helpStringMaxLength := 0
	for _, commandElem := range commandElemSlice {
		helpString := libutil.GetRegexHelpString(commandElem.Regex)
		// compare with the longest one
		if helpStringMaxLength < len(helpString) {
			helpStringMaxLength = len(helpString)
		}
		helpStringSlice = append(helpStringSlice, helpString)
	}

	for i, helpString := range helpStringSlice {
		if commandElemSlice[i].Desc != "" {
			fmt.Printf("%*s   \"%s\"\n", helpStringMaxLength+1, helpString, commandElemSlice[i].Desc)
		} else {
			fmt.Printf("%*s\n", helpStringMaxLength+1, helpString)
		}
	}
	fmt.Printf("# %s", cli.buf)

}

// tab completion
func (cli *GoCli) tabCompletion() {
	args := cli.getArgs()
	lastArg := args[len(args)-1] // get last arg
	commandElemSlice := cli.getCommandElemList(cli.getArgs()...)
	// case 1: no match							eg. # history eee
	// case 2: incomplete sentece.				eg. # history cle
	// case 3: complete sentece in regex. 		eg. # histroy 10
	// case 4: already complete sentence. 		eg. # history or # history clear

	// case 1: no match
	if len(commandElemSlice) == 0 {
		cli.logger.Info("no match\n")
		return
	}

	// case 2: incomplete sentence
	if len(commandElemSlice) == 1 && !cli.getIsLastCharEmpty() {
		cli.logger.Info("one match")

		// clear line and auto complete to cli buf
		libutil.ClearLine()                                     // clear line to rewrite cli buf on stdout
		cli.buf = []byte(strings.Join(args[:len(args)-1], " ")) // join all args to cli buf

		cli.logger.Info("len=%d, args=%v, subargs=%v, arg=%s, regex=%s",
			len(args), args, args[:len(args)-1], lastArg, commandElemSlice[0].Regex)

		if len(args) > 1 {
			cli.buf = append(cli.buf, ' ')
		}

		// case 2: exact string completion
		if strings.HasPrefix(commandElemSlice[0].Regex, lastArg) {
			cli.logger.Info("incomplete sentence")

			cli.buf = append(cli.buf, []byte(commandElemSlice[0].Regex)...) // the Regex is the exact string
		} else { // case 3: regex completion
			cli.logger.Info("complete sentence in regex")

			cli.buf = append(cli.buf, []byte(lastArg)...)
		}
		cli.buf = append(cli.buf, ' ')
		fmt.Printf("# %s", string(cli.buf))
	} else { // case 4: already complete sentence
		cli.logger.Info("multi %d", len(commandElemSlice))

		if len(commandElemSlice) == 1 { // only one match
			if libutil.GetRegexHelpString(commandElemSlice[0].Regex) == commandElemSlice[0].Regex {
				cli.buf = append(cli.buf, []byte(commandElemSlice[0].Regex)...) // the Regex is the exact string
				cli.buf = append(cli.buf, ' ')
			} else {
				fmt.Printf("\n")
				helpString := libutil.GetRegexHelpString(commandElemSlice[0].Regex)
				fmt.Printf(" %s", helpString)
			}
		} else { // multi match then print one line help
			fmt.Printf("\n")

			helpStringSlice := make([]string, 0)
			helpStringMaxLength := 0
			for _, commandElem := range commandElemSlice {
				helpString := libutil.GetRegexHelpString(commandElem.Regex)
				// compare with the longest one
				if helpStringMaxLength < len(helpString) {
					helpStringMaxLength = len(helpString)
				}
				helpStringSlice = append(helpStringSlice, helpString)
			}

			for _, helphelpString := range helpStringSlice {
				fmt.Printf("%*s", helpStringMaxLength+1, helphelpString)
			}

			// auto add space
			if cli.getBufLen() != 0 && cli.buf[cli.getBufLen()-1] != ' ' {
				cli.buf = append(cli.buf, ' ')
			}
		}
		fmt.Printf("\n# %s", cli.buf)
	}

	cli.cursorPos = cli.getBufLen()
}

func (cli *GoCli) runCommand() {
	line := cli.getLineString()
	args := cli.getArgs()

	fmt.Printf("\n")

	// if line is empty, just return
	if len(line) == 0 {
		cli.buf = make([]byte, 0)
		fmt.Printf("# ")
		return
	}

	// update history
	if cli.getHistoryLen() == 0 || cli.historySlice[cli.getHistoryLen()-1] != line {
		cli.historySlice = append(cli.historySlice, line)
	}
	cli.historyPos = cli.getHistoryLen()

	// find command function. if not found, print command not found
	commandElem := cli.getCommandElemByExactMatch(args...)
	if commandElem != nil && commandElem.Func != nil {
		commandElem.Func(args)
	} else {
		fmt.Printf("command \"%s\" does not exist\n", line)
	}
	fmt.Printf("# ")

	// clear cli buf
	cli.buf = make([]byte, 0)
	cli.cursorPos = 0
}

// select history by direction
func (cli *GoCli) selectHistory(direction int) {
	if cli.getHistoryLen() == 0 {
		return
	}

	// clear line and auto complete to cli buf
	libutil.ClearLine() // clear line to rewrite cli buf on stdout

	// update history pos
	cli.historyPos += direction
	if cli.historyPos < 0 {
		cli.historyPos = 0
	} else if cli.historyPos >= cli.getHistoryLen() {
		cli.historyPos = cli.getHistoryLen() - 1
	}

	// update cli buf
	cli.buf = []byte(cli.historySlice[cli.historyPos])
	cli.cursorPos = cli.getBufLen()
	fmt.Printf("# %s", cli.buf)
}

func (cli *GoCli) getCommandElemByExactMatch(args ...string) *CommandElem {
	commandMap := cli.commandMap
	var commandElem *CommandElem = nil
	for _, arg := range args {
		commandElem = nil
		for regex := range commandMap {
			match, _ := regexp.MatchString("\\b"+regex+"\\b", arg)
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

// get matched command elem list
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
				if libutil.GetRegexHelpString(regex) == regex && arg != regex { // exact string
					continue
				}

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

	sort.Slice(commandElemSlice,
		func(i, j int) bool { return commandElemSlice[i].Regex < commandElemSlice[j].Regex })

	return commandElemSlice
}

func (cli *GoCli) defaultCommand() {
	// show history
	cli.AddCommandElem(
		NewCommandElemWithoutFunc("history", ""),
		NewCommandElem("show", "show history", func(args []string) {
			fmt.Printf("\n")
			for index, each := range cli.historySlice {
				fmt.Printf(" %d %s\n", index, each)
			}
		}))

	// show history	last n
	cli.AddCommandElem(
		NewCommandElemWithoutFunc("history", "show history"),
		NewCommandElemWithoutFunc("show", ""),
		NewCommandElemWithoutFunc("last", "show history last n"),
		NewCommandElem(libutil.NumberRegex, "last n", func(args []string) {
			fmt.Printf("\n")
			for index, each := range cli.historySlice {
				if value, _ := strconv.Atoi(args[3]); value <= index {
					break
				}
				fmt.Printf(" %d %s\n", index, each)
			}
		}))

	// clear history
	cli.AddCommandElem(
		NewCommandElemWithoutFunc("history", ""),
		NewCommandElem("clear", "clear history", func(args []string) {
			fmt.Printf("\n")
			cli.historySlice = make([]string, 0)
		}))

	// clear screen
	cli.AddCommandElem(
		NewCommandElem("clear", "clear screen", func(args []string) {
			libutil.ClearScreen()
		}))

	// quit
	cli.AddCommandElem(
		NewCommandElem("quit", "quit", func(args []string) {
			cli.isRunning = false
		}))

	// show help
	cli.AddCommandElem(
		NewCommandElem("help", "show help", func(args []string) {
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
		}))
}
