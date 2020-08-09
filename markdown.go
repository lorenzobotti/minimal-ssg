package main

import (
	"regexp"
	"strings"
)

type ListType int

const (
	ulType ListType = iota
	olType
)

type ListItemInfo struct {
	Type        ListType
	Indentation int
}

//TODO: ordered lists, inline code
func markdownCompile(input []byte) Post {
	//various regexes to interpret the markup
	headingRegex := regexp.MustCompile("^# (.*)")
	subheadingRegex := regexp.MustCompile("^## (.*)")
	boldRegex := regexp.MustCompile(`(\*\*|__)(.*?)(\*\*|__)`)
	italicRegex := regexp.MustCompile(`(\*|_)(.*?)(\*|_)`)
	imgRegex := regexp.MustCompile(`\!\[Image\]\((.*?)\)`)
	linkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	quoteRegex := regexp.MustCompile("^> (.*)")
	lineRegex := regexp.MustCompile(`^ *(-+|\*+) *$`)
	inlineCodeRegex := regexp.MustCompile("`(.*?)`")
	ulListRegex := regexp.MustCompile(`^\s*(\*|\+) +(.*)`)
	olListRegex := regexp.MustCompile(`^ *(\d)\. +(.*)`)

	//this is for the markdown parser, I just made it so that if the first line is similar to:
	//[date: 2020-08-07]
	//it interprets that and saves that information
	dateRegex := regexp.MustCompile(`\[date: (\d+)-(\d+)-(\d+)\]`)

	var endProduct strings.Builder
	lines := strings.Split(string(input), "\n")
	lines = append(lines, "")
	documentTitle := ""
	date := ""

	inMultiLnCode := false
	inList := false
	listNesting := 0
	listStack := make([]ListItemInfo, 0)

	for i, line := range lines {
		if line == "```" {
			inMultiLnCode = !inMultiLnCode
			if inMultiLnCode {
				endProduct.WriteString("<pre>\n")
			} else {
				endProduct.WriteString("</pre>\n")
			}
			continue
		}
		//if we are in a code segment, copy the original text verbatim
		if inMultiLnCode {
			endProduct.WriteString(line)
			endProduct.WriteString("\n")
			continue
		}
		newLine := line

		//lists

		//inList stores whether the previous line was a list element or not
		//since lists of different types can be nested, the way we keep track of where we are
		//is with sort of a "stack". Every type an element has more indentation before it than
		//the previous line, it means it starts a nested list, so we push to the listStack
		//the type of list it is (<ul> or <ol>) and how indented it is (number of
		//preceding whitespace characters). This way, when we can close two nested lists at the same time
		//sort of like:
		//* list 1
		//   * sublist 1
		//   * sublist 2
		//      * subsublist 1
		//* list 2       <-- here we close both sublists at the same time
		isList := ulListRegex.MatchString(line) || olListRegex.MatchString(line)
		if isList {
			lineWhitespace := precedingWhitespace(line)

			elementIsUl := ulListRegex.MatchString(line)

			if !inList || lineWhitespace > listNesting {
				//if it has more whitespace than previous (as in: it is nested inside the previous element)
				if elementIsUl {
					listStack = append(listStack, ListItemInfo{ulType, lineWhitespace})
					newLine = "<ul>\n"
				} else {
					listStack = append(listStack, ListItemInfo{olType, lineWhitespace})
					newLine = "<ol>\n"
				}
				listNesting = precedingWhitespace(line)

			} else if lineWhitespace < listNesting {
				//if it has less whitespace than previous (closes previous nest)
				//the real hard part is when you have to close two or more nested lists at the same time
				//we solve that with listStack which basically stores for every level of nesting the
				//type of list (numbered or not) and how much indentation it has
				//(number of spaces or tabs in the markdown code)
				newLine = ""
				for i := len(listStack) - 1; i >= 0; i-- {
					if listStack[i].Indentation > lineWhitespace {
						if listStack[i].Type == ulType {
							newLine += "</ul>\n"
						} else {
							newLine += "</ol>\n"
						}
						listStack = listStack[:i]
					} else {
						listNesting = lineWhitespace
						break
					}
				}
			} else {
				//just vibing
				newLine = ""
			}

			//actually generate the <li> element
			if elementIsUl {
				newLine += ulListRegex.ReplaceAllString(line, "<li>$2</li>")
			} else {
				newLine += olListRegex.ReplaceAllString(line, "<li>$2</li>")
			}
			inList = true
		}

		newLine = headingRegex.ReplaceAllString(newLine, "<h1>$1</h1>")
		newLine = subheadingRegex.ReplaceAllString(newLine, "<h3>$1</h3>")
		newLine = boldRegex.ReplaceAllString(newLine, "<strong>$2</strong>")
		newLine = italicRegex.ReplaceAllString(newLine, "<em>$2</em>")
		newLine = imgRegex.ReplaceAllString(newLine, `<img src="$1"></img>`)
		newLine = linkRegex.ReplaceAllString(newLine, `<a href="$2">$1</a>`)
		newLine = quoteRegex.ReplaceAllString(newLine, `<blockquote>$1</blockquote>`)
		newLine = inlineCodeRegex.ReplaceAllString(newLine, `<code>$1</code>`)

		if !(headingRegex.MatchString(line) || subheadingRegex.MatchString(line) || quoteRegex.MatchString(line) || lineRegex.MatchString(line) || imgRegex.MatchString(line) || isList) && line != "" {
			newLine = "<p>" + newLine + "</p>"
		}

		if lineRegex.MatchString(line) {
			newLine = "<hr>"
		}

		if headingRegex.MatchString(line) && documentTitle == "" {
			documentTitle = headingRegex.ReplaceAllString(line, "$1")
			newLine = ""
		}

		if i == 0 && dateRegex.MatchString(line) {
			date = dateRegex.ReplaceAllString(line, "$3/$2/$1")
			newLine = ""
		} else if i == 0 {
			date = "unknown date"
		}

		//closes the list from the previous line after having parsed the markdown for this one
		if !isList && inList {
			//if the previous line was a list element, but this one is not, close the list(s)
			for i := 0; i < len(listStack); i++ {
				switch listStack[i].Type {
				case ulType:
					newLine = "</ul>\n" + newLine
				case olType:
					newLine = "</ol>\n" + newLine
				}
				listStack = listStack[i:]
			}
			//reinitialize listStack as an empty slice
			//probably not the way you should do it, but I can't think of anything else right now
			listStack = make([]ListItemInfo, 0)
			inList = false
		}

		endProduct.WriteString(newLine)
		endProduct.WriteString("\n")
	}

	return Post{Content: endProduct.String(), Title: documentTitle, PublicationDate: date}
}

//given a line, it returns the number of spaces and tabs at the start of it
//example:
//precedingWhitespace("    culo") == 4
//precedingWhitespace("  culo") == 2
func precedingWhitespace(input string) int {
	counter := 0
	for _, char := range input {
		//TODO: convert the " " and "\t" to runes (performance increase)
		if string(char) == " " || string(char) == "\t" {
			counter++
		} else {
			break
		}
	}
	return counter
}
