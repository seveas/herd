// Code generated from Katyusha.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 13, 112,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 7,
	3, 37, 10, 3, 12, 3, 14, 3, 40, 11, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3,
	5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 6, 7, 57, 10,
	7, 13, 7, 14, 7, 58, 3, 8, 3, 8, 6, 8, 63, 10, 8, 13, 8, 14, 8, 64, 3,
	9, 6, 9, 68, 10, 9, 13, 9, 14, 9, 69, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11,
	3, 11, 7, 11, 78, 10, 11, 12, 11, 14, 11, 81, 11, 11, 3, 11, 3, 11, 3,
	11, 3, 11, 3, 11, 7, 11, 88, 10, 11, 12, 11, 14, 11, 91, 11, 11, 3, 11,
	5, 11, 94, 10, 11, 3, 12, 3, 12, 6, 12, 98, 10, 12, 13, 12, 14, 12, 99,
	3, 13, 6, 13, 103, 10, 13, 13, 13, 14, 13, 104, 3, 14, 3, 14, 5, 14, 109,
	10, 14, 3, 14, 3, 14, 2, 2, 15, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15,
	9, 17, 10, 19, 11, 21, 12, 23, 2, 25, 2, 27, 13, 3, 2, 9, 3, 2, 12, 12,
	3, 2, 50, 59, 6, 2, 48, 48, 67, 92, 97, 97, 99, 124, 7, 2, 48, 48, 50,
	59, 67, 92, 97, 97, 99, 124, 7, 2, 44, 44, 47, 48, 50, 59, 67, 92, 99,
	124, 6, 2, 12, 12, 14, 15, 41, 41, 94, 94, 6, 2, 12, 12, 14, 15, 36, 36,
	94, 94, 2, 121, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2,
	9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2,
	2, 17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 27, 3, 2, 2,
	2, 3, 29, 3, 2, 2, 2, 5, 31, 3, 2, 2, 2, 7, 41, 3, 2, 2, 2, 9, 45, 3, 2,
	2, 2, 11, 49, 3, 2, 2, 2, 13, 56, 3, 2, 2, 2, 15, 60, 3, 2, 2, 2, 17, 67,
	3, 2, 2, 2, 19, 71, 3, 2, 2, 2, 21, 93, 3, 2, 2, 2, 23, 95, 3, 2, 2, 2,
	25, 102, 3, 2, 2, 2, 27, 108, 3, 2, 2, 2, 29, 30, 7, 12, 2, 2, 30, 4, 3,
	2, 2, 2, 31, 32, 7, 116, 2, 2, 32, 33, 7, 119, 2, 2, 33, 34, 7, 112, 2,
	2, 34, 38, 3, 2, 2, 2, 35, 37, 10, 2, 2, 2, 36, 35, 3, 2, 2, 2, 37, 40,
	3, 2, 2, 2, 38, 36, 3, 2, 2, 2, 38, 39, 3, 2, 2, 2, 39, 6, 3, 2, 2, 2,
	40, 38, 3, 2, 2, 2, 41, 42, 7, 117, 2, 2, 42, 43, 7, 103, 2, 2, 43, 44,
	7, 118, 2, 2, 44, 8, 3, 2, 2, 2, 45, 46, 7, 99, 2, 2, 46, 47, 7, 102, 2,
	2, 47, 48, 7, 102, 2, 2, 48, 10, 3, 2, 2, 2, 49, 50, 7, 106, 2, 2, 50,
	51, 7, 113, 2, 2, 51, 52, 7, 117, 2, 2, 52, 53, 7, 118, 2, 2, 53, 54, 7,
	117, 2, 2, 54, 12, 3, 2, 2, 2, 55, 57, 9, 3, 2, 2, 56, 55, 3, 2, 2, 2,
	57, 58, 3, 2, 2, 2, 58, 56, 3, 2, 2, 2, 58, 59, 3, 2, 2, 2, 59, 14, 3,
	2, 2, 2, 60, 62, 9, 4, 2, 2, 61, 63, 9, 5, 2, 2, 62, 61, 3, 2, 2, 2, 63,
	64, 3, 2, 2, 2, 64, 62, 3, 2, 2, 2, 64, 65, 3, 2, 2, 2, 65, 16, 3, 2, 2,
	2, 66, 68, 9, 6, 2, 2, 67, 66, 3, 2, 2, 2, 68, 69, 3, 2, 2, 2, 69, 67,
	3, 2, 2, 2, 69, 70, 3, 2, 2, 2, 70, 18, 3, 2, 2, 2, 71, 72, 7, 63, 2, 2,
	72, 20, 3, 2, 2, 2, 73, 79, 7, 41, 2, 2, 74, 75, 7, 94, 2, 2, 75, 78, 11,
	2, 2, 2, 76, 78, 10, 7, 2, 2, 77, 74, 3, 2, 2, 2, 77, 76, 3, 2, 2, 2, 78,
	81, 3, 2, 2, 2, 79, 77, 3, 2, 2, 2, 79, 80, 3, 2, 2, 2, 80, 82, 3, 2, 2,
	2, 81, 79, 3, 2, 2, 2, 82, 94, 7, 41, 2, 2, 83, 89, 7, 36, 2, 2, 84, 85,
	7, 94, 2, 2, 85, 88, 11, 2, 2, 2, 86, 88, 10, 8, 2, 2, 87, 84, 3, 2, 2,
	2, 87, 86, 3, 2, 2, 2, 88, 91, 3, 2, 2, 2, 89, 87, 3, 2, 2, 2, 89, 90,
	3, 2, 2, 2, 90, 92, 3, 2, 2, 2, 91, 89, 3, 2, 2, 2, 92, 94, 7, 36, 2, 2,
	93, 73, 3, 2, 2, 2, 93, 83, 3, 2, 2, 2, 94, 22, 3, 2, 2, 2, 95, 97, 7,
	37, 2, 2, 96, 98, 10, 2, 2, 2, 97, 96, 3, 2, 2, 2, 98, 99, 3, 2, 2, 2,
	99, 97, 3, 2, 2, 2, 99, 100, 3, 2, 2, 2, 100, 24, 3, 2, 2, 2, 101, 103,
	7, 34, 2, 2, 102, 101, 3, 2, 2, 2, 103, 104, 3, 2, 2, 2, 104, 102, 3, 2,
	2, 2, 104, 105, 3, 2, 2, 2, 105, 26, 3, 2, 2, 2, 106, 109, 5, 25, 13, 2,
	107, 109, 5, 23, 12, 2, 108, 106, 3, 2, 2, 2, 108, 107, 3, 2, 2, 2, 109,
	110, 3, 2, 2, 2, 110, 111, 8, 14, 2, 2, 111, 28, 3, 2, 2, 2, 15, 2, 38,
	58, 64, 69, 77, 79, 87, 89, 93, 99, 104, 108, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'\n'", "", "'set'", "'add'", "'hosts'", "", "", "", "'='",
}

var lexerSymbolicNames = []string{
	"", "", "RUN", "SET", "ADD", "HOSTS", "NUMBER", "IDENTIFIER", "GLOB", "EQUALS",
	"STRING", "SKIP_",
}

var lexerRuleNames = []string{
	"T__0", "RUN", "SET", "ADD", "HOSTS", "NUMBER", "IDENTIFIER", "GLOB", "EQUALS",
	"STRING", "COMMENT", "SPACES", "SKIP_",
}

type KatyushaLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewKatyushaLexer(input antlr.CharStream) *KatyushaLexer {

	l := new(KatyushaLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Katyusha.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// KatyushaLexer tokens.
const (
	KatyushaLexerT__0       = 1
	KatyushaLexerRUN        = 2
	KatyushaLexerSET        = 3
	KatyushaLexerADD        = 4
	KatyushaLexerHOSTS      = 5
	KatyushaLexerNUMBER     = 6
	KatyushaLexerIDENTIFIER = 7
	KatyushaLexerGLOB       = 8
	KatyushaLexerEQUALS     = 9
	KatyushaLexerSTRING     = 10
	KatyushaLexerSKIP_      = 11
)
