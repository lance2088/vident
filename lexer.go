package main

import (
	"container/list"
	"fmt"
	"unicode/utf8"
	"strings"
)

/**
 * This represents each different type of token
 * that our Lexer can produce
 */
const (
	END_OF_FILE = iota
	IDENTIFIER
	OPERATOR
	NUMBER
	STRING
	CHARACTER
	SEPARATOR
	UNKNOWN
)

type Lexer struct {
	input         string
	input_length  int
	pos           int
	line_number   int
	char_number   int
	start_pos     int
	current_char  rune
	running       bool
	buffer        string
	skipped_chars int
	char_width    int
	token_stream  list.List
}

type Token struct {
	token_type  int
	content     string
	line_number int
	char_number int
}

/**
 * Typical way in which a Lexer works,
 * flushes the buffer, but returns the old
 * value in the buffer.
 */
func (self *Lexer) flushBuffer() string {
	result := self.buffer
	self.buffer = ""
	return result
}

/**
 * Create a new token with the given token type and
 * content.
 */
func (self *Lexer) createToken(token_type int, content string) {
	token := &Token{}
	token.token_type = token_type
	token.content = content
	token.line_number = self.line_number
	token.char_number = self.char_number
	fmt.Printf("adding token type %d and content %s to stream\n", token.token_type, token.content)
	self.token_stream.PushBack(token)
}

/**
 * Peek ahead in our input stream,
 * returns a character.
 */
func (self *Lexer) peek(ahead int) rune {
	result, _ := utf8.DecodeRuneInString(self.input[self.pos + ahead:])
	return result
}

func (self *Lexer) createLexer(input string) {
	self.input = input
	self.pos = 0
	self.input_length = len(self.input)
	self.line_number = 1
	self.char_number = 1
	self.start_pos = 0
	self.skipped_chars = 0
	self.current_char, self.char_width = utf8.DecodeRuneInString(self.input[self.pos:])
	self.running = true
}

/**
 * @return if the given character is a number
 */
func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

/**
 * @return if the given character is "junk" character,
 * i.e. anything below the ASCII code of 32.
 */
func isLayout(c rune) bool {
	return c <= 32
}

/**
 * @return if the given character is an
 * uppercase OR lowercase letter.
 */
func isLetter(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

/**
 * @return if the given character is a 
 * letter or a digit
 */
func isLetterOrDigit(c rune) bool {
	return isDigit(c) || isLetter(c)
}

/**
 * @return if the given character is
 * an operator, in this case + - * / =
 */
func isOperator(c rune) bool {
	return strings.Contains("+-*/=", string(c))
}

/**
 * @return if the given character is a separator,
 * i.e. , {} ()
 */
func isSeparator(c rune) bool {
	return strings.Contains(",{}()", string(c))
}

/**
 * This eats up all the "junk" characters, feels weird
 * writing that because I'm British. Eitherway, it will
 * eat the comments, and will eat all of the useless characters
 * such as spaces, tabs, newlines, etc.
 */
func (self *Lexer) skipLayoutAndComments() {
	// eat the junk stuffs
	for isLayout(self.current_char) {
		self.consumeCharacter()
	}

	// coment opener
	if self.current_char == '#' {
		self.consumeCharacter()

		// keep eating till newline
		for self.current_char != '\n' {
			self.consumeCharacter()
		}

		// eat more layout and junk chars
		for isLayout(self.current_char) {
			self.consumeCharacter()
		}
	}
}

/**
 * Consumes the character in the input stream
 */
func (self *Lexer) consumeCharacter() {
	if (self.pos + self.skipped_chars) > self.input_length {
		self.running = false
	}

	if !isLayout(self.current_char) {
		self.buffer = self.buffer + string(self.input[self.pos])
	}

	self.pos++
	self.current_char, self.char_width = utf8.DecodeRuneInString(self.input[self.pos:])
	self.char_number = self.char_number + 1
}

func (self *Lexer) recognizeNumberToken() {
	// consume the first char, either a
	// dot or a decimal
	self.consumeCharacter()

	// .52
	if self.current_char == '.' {					// this is for decimals, i.e .15
		self.consumeCharacter()
		for isDigit(self.current_char) {
			self.consumeCharacter()
		}
	} else {										// 5.12, 6.12, 5.1233123, etc.
		for isDigit(self.current_char) {
			if self.peek(1) == '.' {
				self.consumeCharacter()
				for isDigit(self.current_char) {
					self.consumeCharacter()
				}
			}
			self.consumeCharacter()
		}
	}

	// push our token back, the content is of the
	// buffer
	self.createToken(NUMBER, self.flushBuffer())
}

func (self *Lexer) recognizeIdentifierToken() {
	// consume our first character to get
	// the lexer going
	self.consumeCharacter()
	
	// eat up all of the letters and digits
	for isLetterOrDigit(self.current_char) {
		self.consumeCharacter()
	}
	
	for self.current_char == '_' && isLetterOrDigit(self.peek(1)) {
		self.consumeCharacter()
		for isLetterOrDigit(self.current_char) {
			self.consumeCharacter()
		}
	}

	self.createToken(IDENTIFIER, self.flushBuffer())
}

func (self *Lexer) recognizeSeparatorToken() {
	self.consumeCharacter()
	self.createToken(SEPARATOR, self.flushBuffer())
}

func (self *Lexer) recognizeStringToken() {
	self.consumeCharacter() // eat "

	for self.current_char != '"' {
		self.consumeCharacter()
	}

	self.consumeCharacter() // eat "

	self.createToken(STRING, self.flushBuffer())
}

func (self *Lexer) recognizeCharacterToken() {
	self.consumeCharacter()

	if isLetterOrDigit(self.current_char) {
		self.consumeCharacter()
	}

	self.consumeCharacter()

	self.createToken(CHARACTER, self.flushBuffer())
}

func (self *Lexer) recognizeOperatorToken() {
	self.consumeCharacter()
	self.createToken(OPERATOR, self.flushBuffer())
}

func (self *Lexer) getNextToken() {
	self.start_pos = 0
	for isLayout(self.current_char) {
		self.consumeCharacter()
		self.skipped_chars++
	}
	self.start_pos = self.pos

	if isDigit(self.current_char) || self.current_char == '.' {
		self.recognizeNumberToken()
	} else if isLetterOrDigit(self.current_char) || self.current_char == '_' {
		self.recognizeIdentifierToken()
	} else if self.current_char == '"' {
		self.recognizeStringToken()
	} else if self.current_char == '\'' {
		self.recognizeCharacterToken()
	} else if isOperator(self.current_char) {
		self.recognizeOperatorToken()
	} else if isSeparator(self.current_char) {
		self.recognizeSeparatorToken()
	} else {
		fmt.Printf("unknown token type %d, aka %c\n", self.current_char, self.current_char)
		self.running = false
		return
	}
}

func (self *Lexer) startLexing() {
	for self.running {
		self.getNextToken()
	}

	self.createToken(END_OF_FILE, "<EOF>")
}
