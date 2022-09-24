// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Herd

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type HerdParser struct {
	*antlr.BaseParser
}

var herdParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func herdParserInit() {
	staticData := &herdParserStaticData
	staticData.literalNames = []string{
		"", "'\\n'", "']'", "','", "'}'", "':'", "", "'['", "'{'", "'set'",
		"'add'", "'remove'", "'list'", "'hosts'", "", "", "", "", "'=='", "'=~'",
		"'!='", "'!~'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "", "RUN", "SB_OPEN", "CB_OPEN", "SET", "ADD", "REMOVE",
		"LIST", "HOSTS", "DURATION", "NUMBER", "IDENTIFIER", "GLOB", "EQUALS",
		"MATCHES", "NOT_EQUALS", "NOT_MATCHES", "STRING", "REGEXP", "SKIP_",
	}
	staticData.ruleNames = []string{
		"prog", "line", "run", "set", "add", "remove", "list", "filter", "scalar",
		"value", "array", "hash",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 24, 134, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 1, 0, 5, 0, 26, 8, 0, 10, 0, 12, 0, 29, 9, 0, 1, 0, 1,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 38, 8, 1, 1, 1, 1, 1, 1, 2, 1, 2,
		1, 3, 1, 3, 1, 3, 3, 3, 47, 8, 3, 1, 4, 1, 4, 1, 4, 1, 4, 5, 4, 53, 8,
		4, 10, 4, 12, 4, 56, 9, 4, 1, 4, 4, 4, 59, 8, 4, 11, 4, 12, 4, 60, 3, 4,
		63, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 5, 5, 69, 8, 5, 10, 5, 12, 5, 72, 9,
		5, 1, 5, 4, 5, 75, 8, 5, 11, 5, 12, 5, 76, 3, 5, 79, 8, 5, 1, 6, 1, 6,
		1, 6, 3, 6, 84, 8, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 91, 8, 7, 1,
		8, 1, 8, 1, 9, 1, 9, 1, 9, 3, 9, 98, 8, 9, 1, 10, 1, 10, 1, 10, 1, 10,
		1, 10, 1, 10, 5, 10, 106, 8, 10, 10, 10, 12, 10, 109, 9, 10, 1, 10, 1,
		10, 3, 10, 113, 8, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11,
		1, 11, 1, 11, 1, 11, 5, 11, 125, 8, 11, 10, 11, 12, 11, 128, 9, 11, 1,
		11, 1, 11, 3, 11, 132, 8, 11, 1, 11, 0, 0, 12, 0, 2, 4, 6, 8, 10, 12, 14,
		16, 18, 20, 22, 0, 4, 1, 0, 16, 17, 2, 0, 18, 18, 20, 20, 2, 0, 19, 19,
		21, 21, 2, 0, 14, 16, 22, 22, 142, 0, 27, 1, 0, 0, 0, 2, 37, 1, 0, 0, 0,
		4, 41, 1, 0, 0, 0, 6, 43, 1, 0, 0, 0, 8, 48, 1, 0, 0, 0, 10, 64, 1, 0,
		0, 0, 12, 80, 1, 0, 0, 0, 14, 85, 1, 0, 0, 0, 16, 92, 1, 0, 0, 0, 18, 97,
		1, 0, 0, 0, 20, 112, 1, 0, 0, 0, 22, 131, 1, 0, 0, 0, 24, 26, 3, 2, 1,
		0, 25, 24, 1, 0, 0, 0, 26, 29, 1, 0, 0, 0, 27, 25, 1, 0, 0, 0, 27, 28,
		1, 0, 0, 0, 28, 30, 1, 0, 0, 0, 29, 27, 1, 0, 0, 0, 30, 31, 5, 0, 0, 1,
		31, 1, 1, 0, 0, 0, 32, 38, 3, 4, 2, 0, 33, 38, 3, 6, 3, 0, 34, 38, 3, 8,
		4, 0, 35, 38, 3, 10, 5, 0, 36, 38, 3, 12, 6, 0, 37, 32, 1, 0, 0, 0, 37,
		33, 1, 0, 0, 0, 37, 34, 1, 0, 0, 0, 37, 35, 1, 0, 0, 0, 37, 36, 1, 0, 0,
		0, 37, 38, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39, 40, 5, 1, 0, 0, 40, 3, 1,
		0, 0, 0, 41, 42, 5, 6, 0, 0, 42, 5, 1, 0, 0, 0, 43, 46, 5, 9, 0, 0, 44,
		45, 5, 16, 0, 0, 45, 47, 3, 16, 8, 0, 46, 44, 1, 0, 0, 0, 46, 47, 1, 0,
		0, 0, 47, 7, 1, 0, 0, 0, 48, 49, 5, 10, 0, 0, 49, 62, 5, 13, 0, 0, 50,
		54, 7, 0, 0, 0, 51, 53, 3, 14, 7, 0, 52, 51, 1, 0, 0, 0, 53, 56, 1, 0,
		0, 0, 54, 52, 1, 0, 0, 0, 54, 55, 1, 0, 0, 0, 55, 63, 1, 0, 0, 0, 56, 54,
		1, 0, 0, 0, 57, 59, 3, 14, 7, 0, 58, 57, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0,
		60, 58, 1, 0, 0, 0, 60, 61, 1, 0, 0, 0, 61, 63, 1, 0, 0, 0, 62, 50, 1,
		0, 0, 0, 62, 58, 1, 0, 0, 0, 63, 9, 1, 0, 0, 0, 64, 65, 5, 11, 0, 0, 65,
		78, 5, 13, 0, 0, 66, 70, 7, 0, 0, 0, 67, 69, 3, 14, 7, 0, 68, 67, 1, 0,
		0, 0, 69, 72, 1, 0, 0, 0, 70, 68, 1, 0, 0, 0, 70, 71, 1, 0, 0, 0, 71, 79,
		1, 0, 0, 0, 72, 70, 1, 0, 0, 0, 73, 75, 3, 14, 7, 0, 74, 73, 1, 0, 0, 0,
		75, 76, 1, 0, 0, 0, 76, 74, 1, 0, 0, 0, 76, 77, 1, 0, 0, 0, 77, 79, 1,
		0, 0, 0, 78, 66, 1, 0, 0, 0, 78, 74, 1, 0, 0, 0, 79, 11, 1, 0, 0, 0, 80,
		81, 5, 12, 0, 0, 81, 83, 5, 13, 0, 0, 82, 84, 3, 22, 11, 0, 83, 82, 1,
		0, 0, 0, 83, 84, 1, 0, 0, 0, 84, 13, 1, 0, 0, 0, 85, 90, 5, 16, 0, 0, 86,
		87, 7, 1, 0, 0, 87, 91, 3, 16, 8, 0, 88, 89, 7, 2, 0, 0, 89, 91, 5, 23,
		0, 0, 90, 86, 1, 0, 0, 0, 90, 88, 1, 0, 0, 0, 91, 15, 1, 0, 0, 0, 92, 93,
		7, 3, 0, 0, 93, 17, 1, 0, 0, 0, 94, 98, 3, 16, 8, 0, 95, 98, 3, 20, 10,
		0, 96, 98, 3, 22, 11, 0, 97, 94, 1, 0, 0, 0, 97, 95, 1, 0, 0, 0, 97, 96,
		1, 0, 0, 0, 98, 19, 1, 0, 0, 0, 99, 100, 5, 7, 0, 0, 100, 113, 5, 2, 0,
		0, 101, 102, 5, 7, 0, 0, 102, 107, 3, 18, 9, 0, 103, 104, 5, 3, 0, 0, 104,
		106, 3, 18, 9, 0, 105, 103, 1, 0, 0, 0, 106, 109, 1, 0, 0, 0, 107, 105,
		1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108, 110, 1, 0, 0, 0, 109, 107, 1, 0,
		0, 0, 110, 111, 5, 2, 0, 0, 111, 113, 1, 0, 0, 0, 112, 99, 1, 0, 0, 0,
		112, 101, 1, 0, 0, 0, 113, 21, 1, 0, 0, 0, 114, 115, 5, 8, 0, 0, 115, 132,
		5, 4, 0, 0, 116, 117, 5, 8, 0, 0, 117, 118, 5, 16, 0, 0, 118, 119, 5, 5,
		0, 0, 119, 126, 3, 18, 9, 0, 120, 121, 5, 3, 0, 0, 121, 122, 5, 16, 0,
		0, 122, 123, 5, 5, 0, 0, 123, 125, 3, 18, 9, 0, 124, 120, 1, 0, 0, 0, 125,
		128, 1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 126, 127, 1, 0, 0, 0, 127, 129,
		1, 0, 0, 0, 128, 126, 1, 0, 0, 0, 129, 130, 5, 4, 0, 0, 130, 132, 1, 0,
		0, 0, 131, 114, 1, 0, 0, 0, 131, 116, 1, 0, 0, 0, 132, 23, 1, 0, 0, 0,
		16, 27, 37, 46, 54, 60, 62, 70, 76, 78, 83, 90, 97, 107, 112, 126, 131,
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

// HerdParserInit initializes any static state used to implement HerdParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewHerdParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func HerdParserInit() {
	staticData := &herdParserStaticData
	staticData.once.Do(herdParserInit)
}

// NewHerdParser produces a new parser instance for the optional input antlr.TokenStream.
func NewHerdParser(input antlr.TokenStream) *HerdParser {
	HerdParserInit()
	this := new(HerdParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &herdParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "java-escape"

	return this
}

// HerdParser tokens.
const (
	HerdParserEOF         = antlr.TokenEOF
	HerdParserT__0        = 1
	HerdParserT__1        = 2
	HerdParserT__2        = 3
	HerdParserT__3        = 4
	HerdParserT__4        = 5
	HerdParserRUN         = 6
	HerdParserSB_OPEN     = 7
	HerdParserCB_OPEN     = 8
	HerdParserSET         = 9
	HerdParserADD         = 10
	HerdParserREMOVE      = 11
	HerdParserLIST        = 12
	HerdParserHOSTS       = 13
	HerdParserDURATION    = 14
	HerdParserNUMBER      = 15
	HerdParserIDENTIFIER  = 16
	HerdParserGLOB        = 17
	HerdParserEQUALS      = 18
	HerdParserMATCHES     = 19
	HerdParserNOT_EQUALS  = 20
	HerdParserNOT_MATCHES = 21
	HerdParserSTRING      = 22
	HerdParserREGEXP      = 23
	HerdParserSKIP_       = 24
)

// HerdParser rules.
const (
	HerdParserRULE_prog   = 0
	HerdParserRULE_line   = 1
	HerdParserRULE_run    = 2
	HerdParserRULE_set    = 3
	HerdParserRULE_add    = 4
	HerdParserRULE_remove = 5
	HerdParserRULE_list   = 6
	HerdParserRULE_filter = 7
	HerdParserRULE_scalar = 8
	HerdParserRULE_value  = 9
	HerdParserRULE_array  = 10
	HerdParserRULE_hash   = 11
)

// IProgContext is an interface to support dynamic dispatch.
type IProgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsProgContext differentiates from other interfaces.
	IsProgContext()
}

type ProgContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgContext() *ProgContext {
	var p = new(ProgContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_prog
	return p
}

func (*ProgContext) IsProgContext() {}

func NewProgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgContext {
	var p = new(ProgContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_prog

	return p
}

func (s *ProgContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgContext) EOF() antlr.TerminalNode {
	return s.GetToken(HerdParserEOF, 0)
}

func (s *ProgContext) AllLine() []ILineContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILineContext); ok {
			len++
		}
	}

	tst := make([]ILineContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILineContext); ok {
			tst[i] = t.(ILineContext)
			i++
		}
	}

	return tst
}

func (s *ProgContext) Line(i int) ILineContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILineContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILineContext)
}

func (s *ProgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterProg(s)
	}
}

func (s *ProgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitProg(s)
	}
}

func (p *HerdParser) Prog() (localctx IProgContext) {
	this := p
	_ = this

	localctx = NewProgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, HerdParserRULE_prog)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(27)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&7746) != 0 {
		{
			p.SetState(24)
			p.Line()
		}

		p.SetState(29)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(30)
		p.Match(HerdParserEOF)
	}

	return localctx
}

// ILineContext is an interface to support dynamic dispatch.
type ILineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLineContext differentiates from other interfaces.
	IsLineContext()
}

type LineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLineContext() *LineContext {
	var p = new(LineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_line
	return p
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_line

	return p
}

func (s *LineContext) GetParser() antlr.Parser { return s.parser }

func (s *LineContext) Run() IRunContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRunContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRunContext)
}

func (s *LineContext) Set() ISetContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetContext)
}

func (s *LineContext) Add() IAddContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAddContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAddContext)
}

func (s *LineContext) Remove() IRemoveContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRemoveContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRemoveContext)
}

func (s *LineContext) List() IListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListContext)
}

func (s *LineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterLine(s)
	}
}

func (s *LineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitLine(s)
	}
}

func (p *HerdParser) Line() (localctx ILineContext) {
	this := p
	_ = this

	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, HerdParserRULE_line)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(37)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case HerdParserRUN:
		{
			p.SetState(32)
			p.Run()
		}

	case HerdParserSET:
		{
			p.SetState(33)
			p.Set()
		}

	case HerdParserADD:
		{
			p.SetState(34)
			p.Add()
		}

	case HerdParserREMOVE:
		{
			p.SetState(35)
			p.Remove()
		}

	case HerdParserLIST:
		{
			p.SetState(36)
			p.List()
		}

	case HerdParserT__0:

	default:
	}
	{
		p.SetState(39)
		p.Match(HerdParserT__0)
	}

	return localctx
}

// IRunContext is an interface to support dynamic dispatch.
type IRunContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRunContext differentiates from other interfaces.
	IsRunContext()
}

type RunContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRunContext() *RunContext {
	var p = new(RunContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_run
	return p
}

func (*RunContext) IsRunContext() {}

func NewRunContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RunContext {
	var p = new(RunContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_run

	return p
}

func (s *RunContext) GetParser() antlr.Parser { return s.parser }

func (s *RunContext) RUN() antlr.TerminalNode {
	return s.GetToken(HerdParserRUN, 0)
}

func (s *RunContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RunContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RunContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterRun(s)
	}
}

func (s *RunContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitRun(s)
	}
}

func (p *HerdParser) Run() (localctx IRunContext) {
	this := p
	_ = this

	localctx = NewRunContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, HerdParserRULE_run)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(41)
		p.Match(HerdParserRUN)
	}

	return localctx
}

// ISetContext is an interface to support dynamic dispatch.
type ISetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVarname returns the varname token.
	GetVarname() antlr.Token

	// SetVarname sets the varname token.
	SetVarname(antlr.Token)

	// GetVarvalue returns the varvalue rule contexts.
	GetVarvalue() IScalarContext

	// SetVarvalue sets the varvalue rule contexts.
	SetVarvalue(IScalarContext)

	// IsSetContext differentiates from other interfaces.
	IsSetContext()
}

type SetContext struct {
	*antlr.BaseParserRuleContext
	parser   antlr.Parser
	varname  antlr.Token
	varvalue IScalarContext
}

func NewEmptySetContext() *SetContext {
	var p = new(SetContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_set
	return p
}

func (*SetContext) IsSetContext() {}

func NewSetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetContext {
	var p = new(SetContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_set

	return p
}

func (s *SetContext) GetParser() antlr.Parser { return s.parser }

func (s *SetContext) GetVarname() antlr.Token { return s.varname }

func (s *SetContext) SetVarname(v antlr.Token) { s.varname = v }

func (s *SetContext) GetVarvalue() IScalarContext { return s.varvalue }

func (s *SetContext) SetVarvalue(v IScalarContext) { s.varvalue = v }

func (s *SetContext) SET() antlr.TerminalNode {
	return s.GetToken(HerdParserSET, 0)
}

func (s *SetContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(HerdParserIDENTIFIER, 0)
}

func (s *SetContext) Scalar() IScalarContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IScalarContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IScalarContext)
}

func (s *SetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterSet(s)
	}
}

func (s *SetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitSet(s)
	}
}

func (p *HerdParser) Set() (localctx ISetContext) {
	this := p
	_ = this

	localctx = NewSetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, HerdParserRULE_set)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(43)
		p.Match(HerdParserSET)
	}
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == HerdParserIDENTIFIER {
		{
			p.SetState(44)

			var _m = p.Match(HerdParserIDENTIFIER)

			localctx.(*SetContext).varname = _m
		}
		{
			p.SetState(45)

			var _x = p.Scalar()

			localctx.(*SetContext).varvalue = _x
		}

	}

	return localctx
}

// IAddContext is an interface to support dynamic dispatch.
type IAddContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetGlob returns the glob token.
	GetGlob() antlr.Token

	// SetGlob sets the glob token.
	SetGlob(antlr.Token)

	// GetFilters returns the filters rule contexts.
	GetFilters() IFilterContext

	// SetFilters sets the filters rule contexts.
	SetFilters(IFilterContext)

	// IsAddContext differentiates from other interfaces.
	IsAddContext()
}

type AddContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	glob    antlr.Token
	filters IFilterContext
}

func NewEmptyAddContext() *AddContext {
	var p = new(AddContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_add
	return p
}

func (*AddContext) IsAddContext() {}

func NewAddContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AddContext {
	var p = new(AddContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_add

	return p
}

func (s *AddContext) GetParser() antlr.Parser { return s.parser }

func (s *AddContext) GetGlob() antlr.Token { return s.glob }

func (s *AddContext) SetGlob(v antlr.Token) { s.glob = v }

func (s *AddContext) GetFilters() IFilterContext { return s.filters }

func (s *AddContext) SetFilters(v IFilterContext) { s.filters = v }

func (s *AddContext) ADD() antlr.TerminalNode {
	return s.GetToken(HerdParserADD, 0)
}

func (s *AddContext) HOSTS() antlr.TerminalNode {
	return s.GetToken(HerdParserHOSTS, 0)
}

func (s *AddContext) GLOB() antlr.TerminalNode {
	return s.GetToken(HerdParserGLOB, 0)
}

func (s *AddContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(HerdParserIDENTIFIER, 0)
}

func (s *AddContext) AllFilter() []IFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilterContext); ok {
			len++
		}
	}

	tst := make([]IFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilterContext); ok {
			tst[i] = t.(IFilterContext)
			i++
		}
	}

	return tst
}

func (s *AddContext) Filter(i int) IFilterContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterContext)
}

func (s *AddContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AddContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AddContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterAdd(s)
	}
}

func (s *AddContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitAdd(s)
	}
}

func (p *HerdParser) Add() (localctx IAddContext) {
	this := p
	_ = this

	localctx = NewAddContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, HerdParserRULE_add)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(48)
		p.Match(HerdParserADD)
	}
	{
		p.SetState(49)
		p.Match(HerdParserHOSTS)
	}
	p.SetState(62)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(50)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*AddContext).glob = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == HerdParserIDENTIFIER || _la == HerdParserGLOB) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*AddContext).glob = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(54)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserIDENTIFIER {
			{
				p.SetState(51)

				var _x = p.Filter()

				localctx.(*AddContext).filters = _x
			}

			p.SetState(56)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == HerdParserIDENTIFIER {
			{
				p.SetState(57)

				var _x = p.Filter()

				localctx.(*AddContext).filters = _x
			}

			p.SetState(60)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}

	return localctx
}

// IRemoveContext is an interface to support dynamic dispatch.
type IRemoveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetGlob returns the glob token.
	GetGlob() antlr.Token

	// SetGlob sets the glob token.
	SetGlob(antlr.Token)

	// GetFilters returns the filters rule contexts.
	GetFilters() IFilterContext

	// SetFilters sets the filters rule contexts.
	SetFilters(IFilterContext)

	// IsRemoveContext differentiates from other interfaces.
	IsRemoveContext()
}

type RemoveContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	glob    antlr.Token
	filters IFilterContext
}

func NewEmptyRemoveContext() *RemoveContext {
	var p = new(RemoveContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_remove
	return p
}

func (*RemoveContext) IsRemoveContext() {}

func NewRemoveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RemoveContext {
	var p = new(RemoveContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_remove

	return p
}

func (s *RemoveContext) GetParser() antlr.Parser { return s.parser }

func (s *RemoveContext) GetGlob() antlr.Token { return s.glob }

func (s *RemoveContext) SetGlob(v antlr.Token) { s.glob = v }

func (s *RemoveContext) GetFilters() IFilterContext { return s.filters }

func (s *RemoveContext) SetFilters(v IFilterContext) { s.filters = v }

func (s *RemoveContext) REMOVE() antlr.TerminalNode {
	return s.GetToken(HerdParserREMOVE, 0)
}

func (s *RemoveContext) HOSTS() antlr.TerminalNode {
	return s.GetToken(HerdParserHOSTS, 0)
}

func (s *RemoveContext) GLOB() antlr.TerminalNode {
	return s.GetToken(HerdParserGLOB, 0)
}

func (s *RemoveContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(HerdParserIDENTIFIER, 0)
}

func (s *RemoveContext) AllFilter() []IFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilterContext); ok {
			len++
		}
	}

	tst := make([]IFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilterContext); ok {
			tst[i] = t.(IFilterContext)
			i++
		}
	}

	return tst
}

func (s *RemoveContext) Filter(i int) IFilterContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterContext)
}

func (s *RemoveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RemoveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RemoveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterRemove(s)
	}
}

func (s *RemoveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitRemove(s)
	}
}

func (p *HerdParser) Remove() (localctx IRemoveContext) {
	this := p
	_ = this

	localctx = NewRemoveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, HerdParserRULE_remove)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(64)
		p.Match(HerdParserREMOVE)
	}
	{
		p.SetState(65)
		p.Match(HerdParserHOSTS)
	}
	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(66)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*RemoveContext).glob = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == HerdParserIDENTIFIER || _la == HerdParserGLOB) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*RemoveContext).glob = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(70)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserIDENTIFIER {
			{
				p.SetState(67)

				var _x = p.Filter()

				localctx.(*RemoveContext).filters = _x
			}

			p.SetState(72)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.SetState(74)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == HerdParserIDENTIFIER {
			{
				p.SetState(73)

				var _x = p.Filter()

				localctx.(*RemoveContext).filters = _x
			}

			p.SetState(76)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}

	return localctx
}

// IListContext is an interface to support dynamic dispatch.
type IListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOpts returns the opts rule contexts.
	GetOpts() IHashContext

	// SetOpts sets the opts rule contexts.
	SetOpts(IHashContext)

	// IsListContext differentiates from other interfaces.
	IsListContext()
}

type ListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	opts   IHashContext
}

func NewEmptyListContext() *ListContext {
	var p = new(ListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_list
	return p
}

func (*ListContext) IsListContext() {}

func NewListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListContext {
	var p = new(ListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_list

	return p
}

func (s *ListContext) GetParser() antlr.Parser { return s.parser }

func (s *ListContext) GetOpts() IHashContext { return s.opts }

func (s *ListContext) SetOpts(v IHashContext) { s.opts = v }

func (s *ListContext) LIST() antlr.TerminalNode {
	return s.GetToken(HerdParserLIST, 0)
}

func (s *ListContext) HOSTS() antlr.TerminalNode {
	return s.GetToken(HerdParserHOSTS, 0)
}

func (s *ListContext) Hash() IHashContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHashContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHashContext)
}

func (s *ListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterList(s)
	}
}

func (s *ListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitList(s)
	}
}

func (p *HerdParser) List() (localctx IListContext) {
	this := p
	_ = this

	localctx = NewListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, HerdParserRULE_list)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(80)
		p.Match(HerdParserLIST)
	}
	{
		p.SetState(81)
		p.Match(HerdParserHOSTS)
	}
	p.SetState(83)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == HerdParserCB_OPEN {
		{
			p.SetState(82)

			var _x = p.Hash()

			localctx.(*ListContext).opts = _x
		}

	}

	return localctx
}

// IFilterContext is an interface to support dynamic dispatch.
type IFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetKey returns the key token.
	GetKey() antlr.Token

	// GetComp returns the comp token.
	GetComp() antlr.Token

	// GetRx returns the rx token.
	GetRx() antlr.Token

	// SetKey sets the key token.
	SetKey(antlr.Token)

	// SetComp sets the comp token.
	SetComp(antlr.Token)

	// SetRx sets the rx token.
	SetRx(antlr.Token)

	// GetVal returns the val rule contexts.
	GetVal() IScalarContext

	// SetVal sets the val rule contexts.
	SetVal(IScalarContext)

	// IsFilterContext differentiates from other interfaces.
	IsFilterContext()
}

type FilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	key    antlr.Token
	comp   antlr.Token
	val    IScalarContext
	rx     antlr.Token
}

func NewEmptyFilterContext() *FilterContext {
	var p = new(FilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_filter
	return p
}

func (*FilterContext) IsFilterContext() {}

func NewFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterContext {
	var p = new(FilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_filter

	return p
}

func (s *FilterContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterContext) GetKey() antlr.Token { return s.key }

func (s *FilterContext) GetComp() antlr.Token { return s.comp }

func (s *FilterContext) GetRx() antlr.Token { return s.rx }

func (s *FilterContext) SetKey(v antlr.Token) { s.key = v }

func (s *FilterContext) SetComp(v antlr.Token) { s.comp = v }

func (s *FilterContext) SetRx(v antlr.Token) { s.rx = v }

func (s *FilterContext) GetVal() IScalarContext { return s.val }

func (s *FilterContext) SetVal(v IScalarContext) { s.val = v }

func (s *FilterContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(HerdParserIDENTIFIER, 0)
}

func (s *FilterContext) Scalar() IScalarContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IScalarContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IScalarContext)
}

func (s *FilterContext) REGEXP() antlr.TerminalNode {
	return s.GetToken(HerdParserREGEXP, 0)
}

func (s *FilterContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(HerdParserEQUALS, 0)
}

func (s *FilterContext) NOT_EQUALS() antlr.TerminalNode {
	return s.GetToken(HerdParserNOT_EQUALS, 0)
}

func (s *FilterContext) MATCHES() antlr.TerminalNode {
	return s.GetToken(HerdParserMATCHES, 0)
}

func (s *FilterContext) NOT_MATCHES() antlr.TerminalNode {
	return s.GetToken(HerdParserNOT_MATCHES, 0)
}

func (s *FilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterFilter(s)
	}
}

func (s *FilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitFilter(s)
	}
}

func (p *HerdParser) Filter() (localctx IFilterContext) {
	this := p
	_ = this

	localctx = NewFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, HerdParserRULE_filter)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)

		var _m = p.Match(HerdParserIDENTIFIER)

		localctx.(*FilterContext).key = _m
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case HerdParserEQUALS, HerdParserNOT_EQUALS:
		{
			p.SetState(86)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*FilterContext).comp = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == HerdParserEQUALS || _la == HerdParserNOT_EQUALS) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*FilterContext).comp = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(87)

			var _x = p.Scalar()

			localctx.(*FilterContext).val = _x
		}

	case HerdParserMATCHES, HerdParserNOT_MATCHES:
		{
			p.SetState(88)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*FilterContext).comp = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == HerdParserMATCHES || _la == HerdParserNOT_MATCHES) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*FilterContext).comp = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(89)

			var _m = p.Match(HerdParserREGEXP)

			localctx.(*FilterContext).rx = _m
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IScalarContext is an interface to support dynamic dispatch.
type IScalarContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsScalarContext differentiates from other interfaces.
	IsScalarContext()
}

type ScalarContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyScalarContext() *ScalarContext {
	var p = new(ScalarContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_scalar
	return p
}

func (*ScalarContext) IsScalarContext() {}

func NewScalarContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScalarContext {
	var p = new(ScalarContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_scalar

	return p
}

func (s *ScalarContext) GetParser() antlr.Parser { return s.parser }

func (s *ScalarContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(HerdParserNUMBER, 0)
}

func (s *ScalarContext) STRING() antlr.TerminalNode {
	return s.GetToken(HerdParserSTRING, 0)
}

func (s *ScalarContext) DURATION() antlr.TerminalNode {
	return s.GetToken(HerdParserDURATION, 0)
}

func (s *ScalarContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(HerdParserIDENTIFIER, 0)
}

func (s *ScalarContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScalarContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ScalarContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterScalar(s)
	}
}

func (s *ScalarContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitScalar(s)
	}
}

func (p *HerdParser) Scalar() (localctx IScalarContext) {
	this := p
	_ = this

	localctx = NewScalarContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, HerdParserRULE_scalar)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(92)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4308992) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_value
	return p
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) Scalar() IScalarContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IScalarContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IScalarContext)
}

func (s *ValueContext) Array() IArrayContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrayContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrayContext)
}

func (s *ValueContext) Hash() IHashContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHashContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHashContext)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitValue(s)
	}
}

func (p *HerdParser) Value() (localctx IValueContext) {
	this := p
	_ = this

	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, HerdParserRULE_value)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(97)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case HerdParserDURATION, HerdParserNUMBER, HerdParserIDENTIFIER, HerdParserSTRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(94)
			p.Scalar()
		}

	case HerdParserSB_OPEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(95)
			p.Array()
		}

	case HerdParserCB_OPEN:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(96)
			p.Hash()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IArrayContext is an interface to support dynamic dispatch.
type IArrayContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArrayContext differentiates from other interfaces.
	IsArrayContext()
}

type ArrayContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrayContext() *ArrayContext {
	var p = new(ArrayContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_array
	return p
}

func (*ArrayContext) IsArrayContext() {}

func NewArrayContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrayContext {
	var p = new(ArrayContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_array

	return p
}

func (s *ArrayContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrayContext) SB_OPEN() antlr.TerminalNode {
	return s.GetToken(HerdParserSB_OPEN, 0)
}

func (s *ArrayContext) AllValue() []IValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueContext); ok {
			len++
		}
	}

	tst := make([]IValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueContext); ok {
			tst[i] = t.(IValueContext)
			i++
		}
	}

	return tst
}

func (s *ArrayContext) Value(i int) IValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *ArrayContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrayContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrayContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterArray(s)
	}
}

func (s *ArrayContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitArray(s)
	}
}

func (p *HerdParser) Array() (localctx IArrayContext) {
	this := p
	_ = this

	localctx = NewArrayContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, HerdParserRULE_array)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(99)
			p.Match(HerdParserSB_OPEN)
		}
		{
			p.SetState(100)
			p.Match(HerdParserT__1)
		}

	case 2:
		{
			p.SetState(101)
			p.Match(HerdParserSB_OPEN)
		}
		{
			p.SetState(102)
			p.Value()
		}
		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserT__2 {
			{
				p.SetState(103)
				p.Match(HerdParserT__2)
			}
			{
				p.SetState(104)
				p.Value()
			}

			p.SetState(109)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(110)
			p.Match(HerdParserT__1)
		}

	}

	return localctx
}

// IHashContext is an interface to support dynamic dispatch.
type IHashContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHashContext differentiates from other interfaces.
	IsHashContext()
}

type HashContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHashContext() *HashContext {
	var p = new(HashContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = HerdParserRULE_hash
	return p
}

func (*HashContext) IsHashContext() {}

func NewHashContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HashContext {
	var p = new(HashContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = HerdParserRULE_hash

	return p
}

func (s *HashContext) GetParser() antlr.Parser { return s.parser }

func (s *HashContext) CB_OPEN() antlr.TerminalNode {
	return s.GetToken(HerdParserCB_OPEN, 0)
}

func (s *HashContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(HerdParserIDENTIFIER)
}

func (s *HashContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(HerdParserIDENTIFIER, i)
}

func (s *HashContext) AllValue() []IValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueContext); ok {
			len++
		}
	}

	tst := make([]IValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueContext); ok {
			tst[i] = t.(IValueContext)
			i++
		}
	}

	return tst
}

func (s *HashContext) Value(i int) IValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *HashContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HashContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HashContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.EnterHash(s)
	}
}

func (s *HashContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(HerdListener); ok {
		listenerT.ExitHash(s)
	}
}

func (p *HerdParser) Hash() (localctx IHashContext) {
	this := p
	_ = this

	localctx = NewHashContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, HerdParserRULE_hash)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(114)
			p.Match(HerdParserCB_OPEN)
		}
		{
			p.SetState(115)
			p.Match(HerdParserT__3)
		}

	case 2:
		{
			p.SetState(116)
			p.Match(HerdParserCB_OPEN)
		}
		{
			p.SetState(117)
			p.Match(HerdParserIDENTIFIER)
		}
		{
			p.SetState(118)
			p.Match(HerdParserT__4)
		}
		{
			p.SetState(119)
			p.Value()
		}
		p.SetState(126)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserT__2 {
			{
				p.SetState(120)
				p.Match(HerdParserT__2)
			}
			{
				p.SetState(121)
				p.Match(HerdParserIDENTIFIER)
			}
			{
				p.SetState(122)
				p.Match(HerdParserT__4)
			}
			{
				p.SetState(123)
				p.Value()
			}

			p.SetState(128)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(129)
			p.Match(HerdParserT__3)
		}

	}

	return localctx
}
