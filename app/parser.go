package main

import (
	"log"
	"strings"
)

type RedirectType int

const (
	RedirectNone   RedirectType = iota
	RedirectStdout              // > or 1>
	RedirectStderr              // 2>
)

var redirectType RedirectType
var redirectFile string

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
			if r == '>' {
				if b.Len() > 0 {
					parsedArgs = append(parsedArgs, b.String())
					b.Reset()
				}
				redirectType = RedirectStdout
				isRedirect = true
				skipNextSpaces = true

				continue
			}

			if (r == '1' || r == '2') && i+1 < n && runes[i+1] == '>' {
				if b.Len() > 0 {
					parsedArgs = append(parsedArgs, b.String())
					b.Reset()
				}

				if r == '1' {
					redirectType = RedirectStdout
				} else {
					redirectType = RedirectStderr
				}

				isRedirect = true
				skipNextSpaces = true
				i++
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

func printSliceWithIndexAndVal(s []string) {
	for i, v := range s {
		log.Println(i, v)
	}
}
