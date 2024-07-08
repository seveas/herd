// Code generated from Herd.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type HerdLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var HerdLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func herdlexerLexerInit() {
	staticData := &HerdLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'\\n'", "']'", "','", "'}'", "':'", "", "'['", "'{'", "'set'",
		"'add'", "'remove'", "'list'", "'hosts'", "", "", "", "", "'=='", "'=~'",
		"'!='", "'!~'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "RUN", "SB_OPEN", "CB_OPEN", "SET", "ADD", "REMOVE",
		"LIST", "HOSTS", "DURATION", "NUMBER", "IDENTIFIER", "GLOB", "EQUALS",
		"MATCHES", "NOT_EQUALS", "NOT_MATCHES", "STRING", "REGEXP", "SKIP_",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "RUN", "SB_OPEN", "CB_OPEN",
		"SET", "ADD", "REMOVE", "LIST", "HOSTS", "DURATION", "NUMBER", "IDENTIFIER",
		"GLOB", "EQUALS", "MATCHES", "NOT_EQUALS", "NOT_MATCHES", "STRING",
		"REGEXP", "COMMENT", "SPACES", "SKIP_",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 24, 200, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5,
		1, 5, 1, 5, 1, 5, 5, 5, 69, 8, 5, 10, 5, 12, 5, 72, 9, 5, 1, 6, 1, 6, 1,
		7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10,
		1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1,
		12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 13, 3, 13, 105, 8, 13, 1, 13,
		4, 13, 108, 8, 13, 11, 13, 12, 13, 109, 1, 13, 1, 13, 4, 13, 114, 8, 13,
		11, 13, 12, 13, 115, 3, 13, 118, 8, 13, 1, 13, 4, 13, 121, 8, 13, 11, 13,
		12, 13, 122, 1, 14, 1, 14, 3, 14, 127, 8, 14, 1, 14, 4, 14, 130, 8, 14,
		11, 14, 12, 14, 131, 1, 15, 1, 15, 5, 15, 136, 8, 15, 10, 15, 12, 15, 139,
		9, 15, 1, 15, 1, 15, 3, 15, 143, 8, 15, 1, 16, 4, 16, 146, 8, 16, 11, 16,
		12, 16, 147, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1,
		19, 1, 20, 1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 5, 21, 166, 8, 21,
		10, 21, 12, 21, 169, 9, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 22, 5,
		22, 177, 8, 22, 10, 22, 12, 22, 180, 9, 22, 1, 22, 1, 22, 1, 23, 1, 23,
		4, 23, 186, 8, 23, 11, 23, 12, 23, 187, 1, 24, 4, 24, 191, 8, 24, 11, 24,
		12, 24, 192, 1, 25, 1, 25, 3, 25, 197, 8, 25, 1, 25, 1, 25, 0, 0, 26, 1,
		1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11,
		23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39, 20,
		41, 21, 43, 22, 45, 23, 47, 0, 49, 0, 51, 24, 1, 0, 10, 1, 0, 10, 10, 1,
		0, 48, 57, 3, 0, 104, 104, 109, 109, 115, 115, 3, 0, 65, 90, 95, 95, 97,
		122, 5, 0, 45, 46, 48, 58, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90,
		95, 95, 97, 122, 2, 0, 65, 90, 97, 122, 6, 0, 42, 42, 45, 46, 48, 57, 63,
		63, 65, 90, 97, 122, 4, 0, 10, 10, 12, 13, 34, 34, 92, 92, 4, 0, 10, 10,
		12, 13, 47, 47, 92, 92, 215, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5,
		1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13,
		1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0,
		21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0,
		0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0,
		0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1, 0,
		0, 0, 0, 45, 1, 0, 0, 0, 0, 51, 1, 0, 0, 0, 1, 53, 1, 0, 0, 0, 3, 55, 1,
		0, 0, 0, 5, 57, 1, 0, 0, 0, 7, 59, 1, 0, 0, 0, 9, 61, 1, 0, 0, 0, 11, 63,
		1, 0, 0, 0, 13, 73, 1, 0, 0, 0, 15, 75, 1, 0, 0, 0, 17, 77, 1, 0, 0, 0,
		19, 81, 1, 0, 0, 0, 21, 85, 1, 0, 0, 0, 23, 92, 1, 0, 0, 0, 25, 97, 1,
		0, 0, 0, 27, 120, 1, 0, 0, 0, 29, 126, 1, 0, 0, 0, 31, 142, 1, 0, 0, 0,
		33, 145, 1, 0, 0, 0, 35, 149, 1, 0, 0, 0, 37, 152, 1, 0, 0, 0, 39, 155,
		1, 0, 0, 0, 41, 158, 1, 0, 0, 0, 43, 161, 1, 0, 0, 0, 45, 172, 1, 0, 0,
		0, 47, 183, 1, 0, 0, 0, 49, 190, 1, 0, 0, 0, 51, 196, 1, 0, 0, 0, 53, 54,
		5, 10, 0, 0, 54, 2, 1, 0, 0, 0, 55, 56, 5, 93, 0, 0, 56, 4, 1, 0, 0, 0,
		57, 58, 5, 44, 0, 0, 58, 6, 1, 0, 0, 0, 59, 60, 5, 125, 0, 0, 60, 8, 1,
		0, 0, 0, 61, 62, 5, 58, 0, 0, 62, 10, 1, 0, 0, 0, 63, 64, 5, 114, 0, 0,
		64, 65, 5, 117, 0, 0, 65, 66, 5, 110, 0, 0, 66, 70, 1, 0, 0, 0, 67, 69,
		8, 0, 0, 0, 68, 67, 1, 0, 0, 0, 69, 72, 1, 0, 0, 0, 70, 68, 1, 0, 0, 0,
		70, 71, 1, 0, 0, 0, 71, 12, 1, 0, 0, 0, 72, 70, 1, 0, 0, 0, 73, 74, 5,
		91, 0, 0, 74, 14, 1, 0, 0, 0, 75, 76, 5, 123, 0, 0, 76, 16, 1, 0, 0, 0,
		77, 78, 5, 115, 0, 0, 78, 79, 5, 101, 0, 0, 79, 80, 5, 116, 0, 0, 80, 18,
		1, 0, 0, 0, 81, 82, 5, 97, 0, 0, 82, 83, 5, 100, 0, 0, 83, 84, 5, 100,
		0, 0, 84, 20, 1, 0, 0, 0, 85, 86, 5, 114, 0, 0, 86, 87, 5, 101, 0, 0, 87,
		88, 5, 109, 0, 0, 88, 89, 5, 111, 0, 0, 89, 90, 5, 118, 0, 0, 90, 91, 5,
		101, 0, 0, 91, 22, 1, 0, 0, 0, 92, 93, 5, 108, 0, 0, 93, 94, 5, 105, 0,
		0, 94, 95, 5, 115, 0, 0, 95, 96, 5, 116, 0, 0, 96, 24, 1, 0, 0, 0, 97,
		98, 5, 104, 0, 0, 98, 99, 5, 111, 0, 0, 99, 100, 5, 115, 0, 0, 100, 101,
		5, 116, 0, 0, 101, 102, 5, 115, 0, 0, 102, 26, 1, 0, 0, 0, 103, 105, 5,
		45, 0, 0, 104, 103, 1, 0, 0, 0, 104, 105, 1, 0, 0, 0, 105, 107, 1, 0, 0,
		0, 106, 108, 7, 1, 0, 0, 107, 106, 1, 0, 0, 0, 108, 109, 1, 0, 0, 0, 109,
		107, 1, 0, 0, 0, 109, 110, 1, 0, 0, 0, 110, 117, 1, 0, 0, 0, 111, 113,
		5, 46, 0, 0, 112, 114, 7, 1, 0, 0, 113, 112, 1, 0, 0, 0, 114, 115, 1, 0,
		0, 0, 115, 113, 1, 0, 0, 0, 115, 116, 1, 0, 0, 0, 116, 118, 1, 0, 0, 0,
		117, 111, 1, 0, 0, 0, 117, 118, 1, 0, 0, 0, 118, 119, 1, 0, 0, 0, 119,
		121, 7, 2, 0, 0, 120, 104, 1, 0, 0, 0, 121, 122, 1, 0, 0, 0, 122, 120,
		1, 0, 0, 0, 122, 123, 1, 0, 0, 0, 123, 28, 1, 0, 0, 0, 124, 125, 5, 48,
		0, 0, 125, 127, 5, 120, 0, 0, 126, 124, 1, 0, 0, 0, 126, 127, 1, 0, 0,
		0, 127, 129, 1, 0, 0, 0, 128, 130, 7, 1, 0, 0, 129, 128, 1, 0, 0, 0, 130,
		131, 1, 0, 0, 0, 131, 129, 1, 0, 0, 0, 131, 132, 1, 0, 0, 0, 132, 30, 1,
		0, 0, 0, 133, 137, 7, 3, 0, 0, 134, 136, 7, 4, 0, 0, 135, 134, 1, 0, 0,
		0, 136, 139, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 137, 138, 1, 0, 0, 0, 138,
		140, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 140, 143, 7, 5, 0, 0, 141, 143,
		7, 6, 0, 0, 142, 133, 1, 0, 0, 0, 142, 141, 1, 0, 0, 0, 143, 32, 1, 0,
		0, 0, 144, 146, 7, 7, 0, 0, 145, 144, 1, 0, 0, 0, 146, 147, 1, 0, 0, 0,
		147, 145, 1, 0, 0, 0, 147, 148, 1, 0, 0, 0, 148, 34, 1, 0, 0, 0, 149, 150,
		5, 61, 0, 0, 150, 151, 5, 61, 0, 0, 151, 36, 1, 0, 0, 0, 152, 153, 5, 61,
		0, 0, 153, 154, 5, 126, 0, 0, 154, 38, 1, 0, 0, 0, 155, 156, 5, 33, 0,
		0, 156, 157, 5, 61, 0, 0, 157, 40, 1, 0, 0, 0, 158, 159, 5, 33, 0, 0, 159,
		160, 5, 126, 0, 0, 160, 42, 1, 0, 0, 0, 161, 167, 5, 34, 0, 0, 162, 163,
		5, 92, 0, 0, 163, 166, 9, 0, 0, 0, 164, 166, 8, 8, 0, 0, 165, 162, 1, 0,
		0, 0, 165, 164, 1, 0, 0, 0, 166, 169, 1, 0, 0, 0, 167, 165, 1, 0, 0, 0,
		167, 168, 1, 0, 0, 0, 168, 170, 1, 0, 0, 0, 169, 167, 1, 0, 0, 0, 170,
		171, 5, 34, 0, 0, 171, 44, 1, 0, 0, 0, 172, 178, 5, 47, 0, 0, 173, 174,
		5, 92, 0, 0, 174, 177, 9, 0, 0, 0, 175, 177, 8, 9, 0, 0, 176, 173, 1, 0,
		0, 0, 176, 175, 1, 0, 0, 0, 177, 180, 1, 0, 0, 0, 178, 176, 1, 0, 0, 0,
		178, 179, 1, 0, 0, 0, 179, 181, 1, 0, 0, 0, 180, 178, 1, 0, 0, 0, 181,
		182, 5, 47, 0, 0, 182, 46, 1, 0, 0, 0, 183, 185, 5, 35, 0, 0, 184, 186,
		8, 0, 0, 0, 185, 184, 1, 0, 0, 0, 186, 187, 1, 0, 0, 0, 187, 185, 1, 0,
		0, 0, 187, 188, 1, 0, 0, 0, 188, 48, 1, 0, 0, 0, 189, 191, 5, 32, 0, 0,
		190, 189, 1, 0, 0, 0, 191, 192, 1, 0, 0, 0, 192, 190, 1, 0, 0, 0, 192,
		193, 1, 0, 0, 0, 193, 50, 1, 0, 0, 0, 194, 197, 3, 49, 24, 0, 195, 197,
		3, 47, 23, 0, 196, 194, 1, 0, 0, 0, 196, 195, 1, 0, 0, 0, 197, 198, 1,
		0, 0, 0, 198, 199, 6, 25, 0, 0, 199, 52, 1, 0, 0, 0, 19, 0, 70, 104, 109,
		115, 117, 122, 126, 131, 137, 142, 147, 165, 167, 176, 178, 187, 192, 196,
		1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// HerdLexerInit initializes any static state used to implement HerdLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewHerdLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func HerdLexerInit() {
	staticData := &HerdLexerLexerStaticData
	staticData.once.Do(herdlexerLexerInit)
}

// NewHerdLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewHerdLexer(input antlr.CharStream) *HerdLexer {
	HerdLexerInit()
	l := new(HerdLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &HerdLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "Herd.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// HerdLexer tokens.
const (
	HerdLexerT__0        = 1
	HerdLexerT__1        = 2
	HerdLexerT__2        = 3
	HerdLexerT__3        = 4
	HerdLexerT__4        = 5
	HerdLexerRUN         = 6
	HerdLexerSB_OPEN     = 7
	HerdLexerCB_OPEN     = 8
	HerdLexerSET         = 9
	HerdLexerADD         = 10
	HerdLexerREMOVE      = 11
	HerdLexerLIST        = 12
	HerdLexerHOSTS       = 13
	HerdLexerDURATION    = 14
	HerdLexerNUMBER      = 15
	HerdLexerIDENTIFIER  = 16
	HerdLexerGLOB        = 17
	HerdLexerEQUALS      = 18
	HerdLexerMATCHES     = 19
	HerdLexerNOT_EQUALS  = 20
	HerdLexerNOT_MATCHES = 21
	HerdLexerSTRING      = 22
	HerdLexerREGEXP      = 23
	HerdLexerSKIP_       = 24
)
