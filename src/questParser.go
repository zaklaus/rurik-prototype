package main

/*
	Quest language parser

	WARN: The parser is not Unicode-aware! Use ASCII characters only!
*/

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/zaklaus/rurik/src/system"
)

// Quest language keywords
const (
	kwTitle      = "title"
	kwBriefing   = "briefing"
	kwResources  = "qrc"
	kwMessage    = "message"
	kwVideo      = "video"
	kwSound      = "sound"
	kwStage      = "stage"
	kwStages     = "qst"
	kwTask       = "task"
	kwSet        = "set"
	kwAbove      = "above"
	kwBelow      = "below"
	kwEquals     = "equals"
	kwNotEquals  = "!equals"
	kwComment    = "$-"
	kwScope      = ":"
	kwLeftBrace  = "("
	kwRightBrace = ")"
)

const (
	tkIdentifier = iota
	tkInteger
	tkSeparator
	tkEndOfFile
)

type questToken struct {
	kind  int
	text  string
	value int

	wordPos int
}

type questTask struct {
	name     string
	commands []questCmd
	pc       int
	isDone   bool
}

type questCmd struct {
	name string
	args []string
}

type questResource struct {
	kind    int
	content string
}

const (
	qrMessage = iota
	qrSound
	qrVideo
	qrStage
)

var questResourceKinds = map[string]int{
	kwMessage: qrMessage,
	kwSound:   qrSound,
	kwVideo:   qrVideo,
	kwStage:   qrStage,
}

type questParser struct {
	data            []byte
	textPos         int
	lastWordPos     int
	allowWhitespace bool
}

func (p *questParser) at(idx int) rune {
	return rune(p.data[idx])
}

func (p *questParser) skipWhitespace() {
	for p.textPos < len(p.data) && (isWhitespace(p.at(p.textPos))) {
		p.textPos++
	}
}

func (p *questParser) skipSeparators() {
	for t := p.peekToken(); t.kind == tkSeparator; t = p.peekToken() {
		p.parseToken()
	}
}

func (p *questParser) peekChar() rune {
	if p.textPos >= len(p.data)-1 {
		return 0
	}

	return p.at(p.textPos)
}

func (p *questParser) nextChar() rune {
	r := p.peekChar()
	p.textPos++

	return r
}

func (p *questParser) parseToken() questToken {
	p.skipWhitespace()

	if p.textPos >= len(p.data)-1 {
		return p.tokenEndOfFile()
	}

	var buf string
	p.lastWordPos = p.textPos

	al := p.allowWhitespace
	brc := 0

	if string(p.peekChar()) == kwLeftBrace {
		p.allowWhitespace = true
	}

	for r := p.peekChar(); r != 0 && (!isWhitespace(r) || p.allowWhitespace) && r != '\n' && string(r) != kwScope; r = p.peekChar() {
		buf += string(p.nextChar())

		if string(r) == kwLeftBrace {
			brc++
		} else if string(r) == kwRightBrace {
			brc--

			if brc == 0 {
				break
			}
		}
	}

	p.allowWhitespace = al

	if len(buf) == 0 && string(p.peekChar()) == kwScope {
		p.nextChar()
		return p.tokenIdentifier(kwScope)
	} else if len(buf) == 0 && p.peekChar() == '\n' {
		sep := 1
		p.nextChar()
		p.skipWhitespace()

		for p.peekChar() == '\n' {
			sep++
			p.nextChar()
			p.skipWhitespace()
		}

		return p.tokenSeparator(sep)
	}

	if val, err := strconv.Atoi(buf); err == nil {
		return p.tokenInteger(val)
	}

	return p.tokenIdentifier(buf)
}

func (p *questParser) nextIdentifier() string {
	p.skipSeparators()
	ident := p.parseToken()

	if ident.kind != tkIdentifier {
		log.Fatalf("Token at '%d' invalid! Expected Identifier.\n", ident.wordPos)
		return ""
	}

	return ident.text
}

func (p *questParser) nextWord() string {
	p.skipSeparators()
	t := p.parseToken()

	if t.kind != tkIdentifier && t.kind != tkInteger {
		log.Fatalf("Word at '%d' invalid! Expected Word.\n", t.wordPos)
		return ""
	}

	return t.text
}

func (p *questParser) nextString() string {
	p.skipSeparators()
	p.allowWhitespace = true
	var buf string

	for tk := p.peekToken(); tk.kind != tkEndOfFile && tk.kind != tkSeparator; tk = p.peekToken() {
		buf += p.nextWord()
	}

	p.allowWhitespace = false
	return buf
}

func (p *questParser) nextTextBlock() string {
	p.skipSeparators()
	p.allowWhitespace = true
	var buf string

	for sep := p.peekToken(); sep.kind != tkEndOfFile &&
		sep.kind != tkSeparator ||
		(sep.kind == tkSeparator && sep.value < 2); sep = p.peekToken() {

		if sep.kind == tkSeparator {
			buf += "\n"
			p.parseToken()
		} else {
			buf += p.nextWord()
		}
	}

	p.allowWhitespace = false
	return buf
}

func (p *questParser) nextNumber() int {
	p.skipSeparators()
	tk := p.parseToken()

	if tk.kind != tkInteger {
		log.Fatalf("Number at '%d' invalid! Expected Number.\n", tk.wordPos)
		return -1
	}

	return tk.value
}

func (p *questParser) expect(ident string) bool {
	p.skipSeparators()
	ok := true

	tk := p.parseToken()

	if tk.kind != tkIdentifier || strings.ToLower(tk.text) != ident {
		log.Fatalf("Unexpected token '%v'! Expected: '%s'.\n", tk, ident)
		ok = false
	}

	return ok
}

func (p *questParser) peekToken() questToken {
	op := *p

	tk := p.parseToken()

	*p = op
	return tk
}

func (p *questParser) checkResourceKind(kind string) bool {
	_, ok := questResourceKinds[strings.ToLower(kind)]
	return ok
}

func (p *questParser) parseResources() (res map[int]questResource) {
	res = map[int]questResource{}

	p.skipSeparators()

	for resKind := p.peekToken(); resKind.kind != tkEndOfFile && p.checkResourceKind(resKind.text); resKind = p.peekToken() {
		p.parseToken()
		p.expect(kwScope)
		resourceID := p.nextNumber()
		kind, _ := questResourceKinds[strings.ToLower(resKind.text)]
		content := p.nextTextBlock()

		res[resourceID] = questResource{
			kind:    kind,
			content: content,
		}

		p.skipSeparators()
	}

	return
}

func (p *questParser) parseTasks() (res []questTask) {
	res = []questTask{}

	// Handle headless main task (entry point)
	res = append(res, questTask{
		name:     "<entry-point>",
		commands: p.parseTask(),
	})

	p.skipSeparators()

	for t := p.peekToken(); t.kind == tkIdentifier; t = p.peekToken() {
		kw := strings.ToLower(p.nextIdentifier())

		if kw != kwTask {
			log.Fatalf("Invalid task found at '%d'!\n", t.wordPos)
			return
		}

		taskName := p.nextIdentifier()
		p.expect(kwScope)

		res = append(res, questTask{
			name:     taskName,
			commands: p.parseTask(),
		})

		log.Printf("Task '%s' has been added!", taskName)

		p.skipSeparators()
	}

	return
}

func (p *questParser) parseTask() (res []questCmd) {
	res = []questCmd{}
	p.skipSeparators()

	for t := p.peekToken(); t.kind == tkIdentifier; t = p.peekToken() {
		// end of the line
		if t.text == kwTask {
			break
		}

		cmd := strings.ToLower(p.nextIdentifier())

		args := []string{}

		for pt := p.peekToken(); pt.kind != tkEndOfFile && pt.kind != tkSeparator; pt = p.peekToken() {
			args = append(args, p.nextWord())
		}

		res = append(res, questCmd{
			name: cmd,
			args: args,
		})

		p.skipSeparators()
	}

	return
}

// questDef describes the quest definition file and the opcodes
type questDef struct {
	title     string
	briefing  string
	resources map[int]questResource
	tasks     []questTask
}

func parseQuest(questName string) *questDef {
	questAsset := system.FindAsset(fmt.Sprintf("quests/%s.qst", strings.ToLower(questName)))

	if questAsset == nil {
		log.Fatalf("Quest '%s' could not be found!\n", questName)
		return nil
	}

	parser := questParser{
		data: questAsset.Data,
	}

	def := &questDef{}

	for t := parser.peekToken(); t.kind != tkEndOfFile; t = parser.peekToken() {
		parser.skipSeparators()
		ident := parser.nextIdentifier()
		parser.expect(kwScope)

		switch strings.ToLower(ident) {
		case kwTitle:
			def.title = parser.nextString()
		case kwBriefing:
			def.briefing = parser.nextTextBlock()
		case kwResources:
			def.resources = parser.parseResources()
		case kwStages:
			def.tasks = parser.parseTasks()
			return def
		default:
			log.Fatalf("Undefined token at '%d'! It says: '%s'.\n", t.wordPos, ident)
			return def
		}
	}

	return def
}

func isAlpha(c rune) bool {
	return unicode.IsLetter(c)
}

func isNumber(c rune) bool {
	return unicode.IsNumber(c)
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isNumber(c)
}

func isWhitespace(c rune) bool {
	return unicode.IsSpace(c) && c != '\n'
}

func (p *questParser) tokenEndOfFile() questToken {
	return questToken{
		kind:    tkEndOfFile,
		wordPos: len(p.data),
	}
}

func (p *questParser) tokenInteger(v int) questToken {
	return questToken{
		kind:    tkInteger,
		text:    strconv.Itoa(v),
		value:   v,
		wordPos: p.lastWordPos,
	}
}

func (p *questParser) tokenIdentifier(s string) questToken {
	return questToken{
		kind:    tkIdentifier,
		text:    s,
		wordPos: p.lastWordPos,
	}
}

func (p *questParser) tokenSeparator(sep int) questToken {
	return questToken{
		kind:    tkSeparator,
		value:   sep,
		wordPos: p.lastWordPos,
	}
}
