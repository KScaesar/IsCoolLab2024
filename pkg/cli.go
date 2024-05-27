package pkg

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func CliSetUsage(command *cobra.Command, group, text string) {
	// help
	// https://github.com/spf13/cobra/blob/v1.8.0/command.go#L593-L596
	// https://github.com/spf13/cobra/blob/v1.8.0/command.go#L459-L468

	// usage
	// https://github.com/spf13/cobra/blob/v1.8.0/command.go#L552-L581
	// https://github.com/spf13/cobra/blob/v1.8.0/command.go#L433-L441

	usage := fmt.Sprintf("%v\n", text)
	command.SetHelpCommandGroupID(group)
	command.UsageString()
	command.SetUsageTemplate(usage)
	return
}

func CliSetActivePrompt(command *cobra.Command, text string) {
	// https://github.com/spf13/cobra/blob/main/site/content/active_help.md

	command.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string
		nArgs := len(args)
		nFlag := cmd.Flags().NFlag()
		// fmt.Printf("[Debug] {nArgs=%v nFlag=%v}\n", nArgs, nFlag)

		paramPrompts := CliParse(text)[1:]
		if nArgs+nFlag >= len(paramPrompts) {
			comps = cobra.AppendActiveHelp(comps, "command has all needed params")
		} else {
			comps = cobra.AppendActiveHelp(comps, paramPrompts[nArgs+nFlag])
		}

		return comps, cobra.ShellCompDirectiveNoFileComp
	}
}

func CliParse(text string) []string {
	const (
		stateFree = 1 << iota
		stateBrackets
		stateQuotation
		stateFlag
		stateFlagHasValue
	)

	var buffer strings.Builder
	var result []string

	prevState := stateFree
	currState := stateFree
	flagHasValue := func() bool {
		return currState&stateFlagHasValue != 0
	}

	var prevChar, nextChar rune
	for i := 0; i < len(text); i++ {
		char := rune(text[i])

		if i+1 < len(text) {
			nextChar = rune(text[i+1])
		}

		// FSM:
		// when char event and current=(state1 && stateX && stateY),
		// then transition from state1 to state2
		switch char {
		case ' ':
			switch currState {
			case stateQuotation:
				buffer.WriteRune(char)
			case stateFlag:
				if !flagHasValue() && prevState == stateBrackets {
					currState ^= stateFlagHasValue
					buffer.WriteRune(char)
				} else {
					currState ^= stateFlagHasValue
					currState = stateFree
					if buffer.Len() > 0 {
						result = append(result, buffer.String())
						buffer.Reset()
					}
				}
			default:
				currState = stateFree
				if buffer.Len() > 0 {
					result = append(result, buffer.String())
					buffer.Reset()
				}
			}

		case '-':
			if (prevChar == ' ' || prevChar == '[') && nextChar == '-' {
				if currState == stateBrackets {
					prevState = currState
				}
				currState = stateFlag
			}
			buffer.WriteRune(char)

		// for prompt
		case '[':
			if currState == stateFree {
				currState = stateBrackets
			}
			buffer.WriteRune(char)
		case ']':
			if currState != stateFlag {
				currState = stateFree
			}
			buffer.WriteRune(char)
		case '?':
			currState = stateFree
			buffer.WriteRune(char)

		// for command
		case '"':
			if currState == stateQuotation {
				currState = prevState
			} else {
				prevState = currState
				currState = stateQuotation
			}
		case '=':
			if currState == stateFlag && !flagHasValue() {
				currState ^= stateFlagHasValue
			}
			buffer.WriteRune(char)

		default:
			buffer.WriteRune(char)
		}

		prevChar = char
		nextChar = 0
	}

	if buffer.Len() > 0 {
		result = append(result, buffer.String())
		buffer.Reset()
	}

	return result
}
