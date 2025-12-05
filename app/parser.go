package main

import (
	"log"
	"strings"
)

type RedirectType int

const (
	RedirectNone         RedirectType = iota
	RedirectStdout                    // > or 1>
	RedirectStderr                    // 2>
	RedirectAppend                    // >>, 1>>
	RedirectStderrAppend              // 2>>
)

var redirectType RedirectType

func parseArgString(args string) (parsedArgs []string, redirectFileName string, redirectType RedirectType) {
	runes := []rune(args)
	n := len(runes)

	// var parsedArgs []string

	isSingleQuote := false
	isDoubleQuote := false
	isRedirect := false
	skipNextSpaces := false

	var rd strings.Builder
	var b strings.Builder

	for i := 0; i < n; i++ {
		r := runes[i]

		if r == '"' && !isSingleQuote {
			isDoubleQuote = !isDoubleQuote
			continue
		}

		if r == '\'' && !isDoubleQuote {
			isSingleQuote = !isSingleQuote
			continue
		}

		if r == '\\' {

			if i == n-1 {

				if isRedirect {
					rd.WriteRune('\\')
				} else {
					b.WriteRune('\\')
				}
				continue
			}

			next := runes[i+1]

			if isSingleQuote {
				b.WriteRune('\\')
				continue
			}

			if isDoubleQuote {
				// only certain escapes in double quotes
				if next == '"' || next == '\\' || next == '$' || next == '`' {
					b.WriteRune(next)
					i++
					continue
				}
				// backslash is literal otherwise
				b.WriteRune('\\')
				continue
			}

			if isRedirect {
				rd.WriteRune(next)
			} else {
				b.WriteRune(next)
			}
			i++
			continue
		}

		if !isSingleQuote && !isDoubleQuote {

			rt, consumed := detectRedirect(runes, i)
			if rt != RedirectNone {
				if b.Len() > 0 {
					parsedArgs = append(parsedArgs, b.String())
					b.Reset()
				}
				redirectType = rt
				isRedirect = true
				skipNextSpaces = true
				i += consumed - 1
				continue
			}

		}

		// consider everything after redirect as redirect file name
		if isRedirect {
			if skipNextSpaces && r == ' ' {
				continue
			}
			skipNextSpaces = false
			if r == ' ' && !isSingleQuote && !isDoubleQuote {
				isRedirect = false
				continue
			}

			rd.WriteRune(r)
			continue
		}

		if r == ' ' && !isSingleQuote && !isDoubleQuote {
			if b.Len() > 0 {
				parsedArgs = append(parsedArgs, b.String())
				b.Reset()
			}
			continue
		}
		b.WriteRune(r)
	}

	if b.Len() > 0 {
		parsedArgs = append(parsedArgs, b.String())
	}

	redirectFileName = strings.TrimSpace(rd.String())
	return parsedArgs, redirectFileName, redirectType
}

func detectRedirect(runes []rune, i int) (RedirectType, int) {
	n := len(runes)

	// first check case 1: 1>> 2>>
	if i+2 < n {
		op := string(runes[i : i+3])
		if op == "1>>" {
			return RedirectAppend, 3
		}
		if op == "2>>" {
			return RedirectStderrAppend, 3
		}
	}

	// Case 2: Simple append >>
	if i+1 < n && string(runes[i:i+2]) == ">>" {
		return RedirectAppend, 2
	}

	// Case 3: Simple redirect: 1>, 2>
	if i+1 < n && (runes[i] == '1' || runes[i] == '2') && runes[i+1] == '>' {
		op := string(runes[i : i+2])
		if op == "1>" {
			return RedirectStdout, 2
		}
		if op == "2>" {
			return RedirectStderr, 2
		}
	}

	// case 4: stdout truncate >
	if runes[i] == '>' {
		return RedirectStdout, 1
	}
	return RedirectNone, 0
}

func printSliceWithIndexAndVal(s []string) {
	for i, v := range s {
		log.Println(i, v)
	}
}
