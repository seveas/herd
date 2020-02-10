// Code generated from Herd.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // Herd

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 26, 135,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 3, 2, 7, 2, 28, 10, 2, 12, 2, 14, 2, 31, 11, 2, 3, 2, 3, 2, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 40, 10, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5,
	3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 7, 6, 54, 10, 6, 12, 6, 14, 6,
	57, 11, 6, 3, 6, 6, 6, 60, 10, 6, 13, 6, 14, 6, 61, 5, 6, 64, 10, 6, 3,
	7, 3, 7, 3, 7, 3, 7, 7, 7, 70, 10, 7, 12, 7, 14, 7, 73, 11, 7, 3, 7, 6,
	7, 76, 10, 7, 13, 7, 14, 7, 77, 5, 7, 80, 10, 7, 3, 8, 3, 8, 3, 8, 5, 8,
	85, 10, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 92, 10, 9, 3, 10, 3, 10,
	3, 11, 3, 11, 3, 11, 5, 11, 99, 10, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3,
	12, 3, 12, 7, 12, 107, 10, 12, 12, 12, 14, 12, 110, 11, 12, 3, 12, 3, 12,
	5, 12, 114, 10, 12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3,
	13, 3, 13, 3, 13, 7, 13, 126, 10, 13, 12, 13, 14, 13, 129, 11, 13, 3, 13,
	3, 13, 5, 13, 133, 10, 13, 3, 13, 2, 2, 14, 2, 4, 6, 8, 10, 12, 14, 16,
	18, 20, 22, 24, 2, 6, 3, 2, 18, 19, 4, 2, 20, 20, 22, 22, 4, 2, 21, 21,
	23, 23, 4, 2, 16, 18, 24, 24, 2, 142, 2, 29, 3, 2, 2, 2, 4, 39, 3, 2, 2,
	2, 6, 43, 3, 2, 2, 2, 8, 45, 3, 2, 2, 2, 10, 49, 3, 2, 2, 2, 12, 65, 3,
	2, 2, 2, 14, 81, 3, 2, 2, 2, 16, 86, 3, 2, 2, 2, 18, 93, 3, 2, 2, 2, 20,
	98, 3, 2, 2, 2, 22, 113, 3, 2, 2, 2, 24, 132, 3, 2, 2, 2, 26, 28, 5, 4,
	3, 2, 27, 26, 3, 2, 2, 2, 28, 31, 3, 2, 2, 2, 29, 27, 3, 2, 2, 2, 29, 30,
	3, 2, 2, 2, 30, 32, 3, 2, 2, 2, 31, 29, 3, 2, 2, 2, 32, 33, 7, 2, 2, 3,
	33, 3, 3, 2, 2, 2, 34, 40, 5, 6, 4, 2, 35, 40, 5, 8, 5, 2, 36, 40, 5, 10,
	6, 2, 37, 40, 5, 12, 7, 2, 38, 40, 5, 14, 8, 2, 39, 34, 3, 2, 2, 2, 39,
	35, 3, 2, 2, 2, 39, 36, 3, 2, 2, 2, 39, 37, 3, 2, 2, 2, 39, 38, 3, 2, 2,
	2, 39, 40, 3, 2, 2, 2, 40, 41, 3, 2, 2, 2, 41, 42, 7, 3, 2, 2, 42, 5, 3,
	2, 2, 2, 43, 44, 7, 8, 2, 2, 44, 7, 3, 2, 2, 2, 45, 46, 7, 11, 2, 2, 46,
	47, 7, 18, 2, 2, 47, 48, 5, 18, 10, 2, 48, 9, 3, 2, 2, 2, 49, 50, 7, 12,
	2, 2, 50, 63, 7, 15, 2, 2, 51, 55, 9, 2, 2, 2, 52, 54, 5, 16, 9, 2, 53,
	52, 3, 2, 2, 2, 54, 57, 3, 2, 2, 2, 55, 53, 3, 2, 2, 2, 55, 56, 3, 2, 2,
	2, 56, 64, 3, 2, 2, 2, 57, 55, 3, 2, 2, 2, 58, 60, 5, 16, 9, 2, 59, 58,
	3, 2, 2, 2, 60, 61, 3, 2, 2, 2, 61, 59, 3, 2, 2, 2, 61, 62, 3, 2, 2, 2,
	62, 64, 3, 2, 2, 2, 63, 51, 3, 2, 2, 2, 63, 59, 3, 2, 2, 2, 64, 11, 3,
	2, 2, 2, 65, 66, 7, 13, 2, 2, 66, 79, 7, 15, 2, 2, 67, 71, 9, 2, 2, 2,
	68, 70, 5, 16, 9, 2, 69, 68, 3, 2, 2, 2, 70, 73, 3, 2, 2, 2, 71, 69, 3,
	2, 2, 2, 71, 72, 3, 2, 2, 2, 72, 80, 3, 2, 2, 2, 73, 71, 3, 2, 2, 2, 74,
	76, 5, 16, 9, 2, 75, 74, 3, 2, 2, 2, 76, 77, 3, 2, 2, 2, 77, 75, 3, 2,
	2, 2, 77, 78, 3, 2, 2, 2, 78, 80, 3, 2, 2, 2, 79, 67, 3, 2, 2, 2, 79, 75,
	3, 2, 2, 2, 80, 13, 3, 2, 2, 2, 81, 82, 7, 14, 2, 2, 82, 84, 7, 15, 2,
	2, 83, 85, 5, 24, 13, 2, 84, 83, 3, 2, 2, 2, 84, 85, 3, 2, 2, 2, 85, 15,
	3, 2, 2, 2, 86, 91, 7, 18, 2, 2, 87, 88, 9, 3, 2, 2, 88, 92, 5, 18, 10,
	2, 89, 90, 9, 4, 2, 2, 90, 92, 7, 25, 2, 2, 91, 87, 3, 2, 2, 2, 91, 89,
	3, 2, 2, 2, 92, 17, 3, 2, 2, 2, 93, 94, 9, 5, 2, 2, 94, 19, 3, 2, 2, 2,
	95, 99, 5, 18, 10, 2, 96, 99, 5, 22, 12, 2, 97, 99, 5, 24, 13, 2, 98, 95,
	3, 2, 2, 2, 98, 96, 3, 2, 2, 2, 98, 97, 3, 2, 2, 2, 99, 21, 3, 2, 2, 2,
	100, 101, 7, 9, 2, 2, 101, 114, 7, 4, 2, 2, 102, 103, 7, 9, 2, 2, 103,
	108, 5, 20, 11, 2, 104, 105, 7, 5, 2, 2, 105, 107, 5, 20, 11, 2, 106, 104,
	3, 2, 2, 2, 107, 110, 3, 2, 2, 2, 108, 106, 3, 2, 2, 2, 108, 109, 3, 2,
	2, 2, 109, 111, 3, 2, 2, 2, 110, 108, 3, 2, 2, 2, 111, 112, 7, 4, 2, 2,
	112, 114, 3, 2, 2, 2, 113, 100, 3, 2, 2, 2, 113, 102, 3, 2, 2, 2, 114,
	23, 3, 2, 2, 2, 115, 116, 7, 10, 2, 2, 116, 133, 7, 6, 2, 2, 117, 118,
	7, 10, 2, 2, 118, 119, 7, 18, 2, 2, 119, 120, 7, 7, 2, 2, 120, 127, 5,
	20, 11, 2, 121, 122, 7, 5, 2, 2, 122, 123, 7, 18, 2, 2, 123, 124, 7, 7,
	2, 2, 124, 126, 5, 20, 11, 2, 125, 121, 3, 2, 2, 2, 126, 129, 3, 2, 2,
	2, 127, 125, 3, 2, 2, 2, 127, 128, 3, 2, 2, 2, 128, 130, 3, 2, 2, 2, 129,
	127, 3, 2, 2, 2, 130, 131, 7, 6, 2, 2, 131, 133, 3, 2, 2, 2, 132, 115,
	3, 2, 2, 2, 132, 117, 3, 2, 2, 2, 133, 25, 3, 2, 2, 2, 17, 29, 39, 55,
	61, 63, 71, 77, 79, 84, 91, 98, 108, 113, 127, 132,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'\n'", "']'", "','", "'}'", "':'", "", "'['", "'{'", "'set'", "'add'",
	"'remove'", "'list'", "'hosts'", "", "", "", "", "'=='", "'=~'", "'!='",
	"'!~'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "RUN", "SB_OPEN", "CB_OPEN", "SET", "ADD", "REMOVE",
	"LIST", "HOSTS", "DURATION", "NUMBER", "IDENTIFIER", "GLOB", "EQUALS",
	"MATCHES", "NOT_EQUALS", "NOT_MATCHES", "STRING", "REGEXP", "SKIP_",
}

var ruleNames = []string{
	"prog", "line", "run", "set", "add", "remove", "list", "filter", "scalar",
	"value", "array", "hash",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type HerdParser struct {
	*antlr.BaseParser
}

func NewHerdParser(input antlr.TokenStream) *HerdParser {
	this := new(HerdParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Herd.g4"

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILineContext)(nil)).Elem())
	var tst = make([]ILineContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILineContext)
		}
	}

	return tst
}

func (s *ProgContext) Line(i int) ILineContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILineContext)(nil)).Elem(), i)

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

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<HerdParserT__0)|(1<<HerdParserRUN)|(1<<HerdParserSET)|(1<<HerdParserADD)|(1<<HerdParserREMOVE)|(1<<HerdParserLIST))) != 0 {
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRunContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRunContext)
}

func (s *LineContext) Set() ISetContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISetContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISetContext)
}

func (s *LineContext) Add() IAddContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAddContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAddContext)
}

func (s *LineContext) Remove() IRemoveContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRemoveContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRemoveContext)
}

func (s *LineContext) List() IListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IListContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IScalarContext)(nil)).Elem(), 0)

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
	localctx = NewSetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, HerdParserRULE_set)

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFilterContext)(nil)).Elem())
	var tst = make([]IFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFilterContext)
		}
	}

	return tst
}

func (s *AddContext) Filter(i int) IFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFilterContext)(nil)).Elem(), i)

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
		p.SetState(47)
		p.Match(HerdParserADD)
	}
	{
		p.SetState(48)
		p.Match(HerdParserHOSTS)
	}
	p.SetState(61)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(49)

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
		p.SetState(53)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserIDENTIFIER {
			{
				p.SetState(50)

				var _x = p.Filter()

				localctx.(*AddContext).filters = _x
			}

			p.SetState(55)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.SetState(57)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == HerdParserIDENTIFIER {
			{
				p.SetState(56)

				var _x = p.Filter()

				localctx.(*AddContext).filters = _x
			}

			p.SetState(59)
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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFilterContext)(nil)).Elem())
	var tst = make([]IFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFilterContext)
		}
	}

	return tst
}

func (s *RemoveContext) Filter(i int) IFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFilterContext)(nil)).Elem(), i)

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
		p.SetState(63)
		p.Match(HerdParserREMOVE)
	}
	{
		p.SetState(64)
		p.Match(HerdParserHOSTS)
	}
	p.SetState(77)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(65)

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
		p.SetState(69)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserIDENTIFIER {
			{
				p.SetState(66)

				var _x = p.Filter()

				localctx.(*RemoveContext).filters = _x
			}

			p.SetState(71)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.SetState(73)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == HerdParserIDENTIFIER {
			{
				p.SetState(72)

				var _x = p.Filter()

				localctx.(*RemoveContext).filters = _x
			}

			p.SetState(75)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHashContext)(nil)).Elem(), 0)

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
		p.SetState(79)
		p.Match(HerdParserLIST)
	}
	{
		p.SetState(80)
		p.Match(HerdParserHOSTS)
	}
	p.SetState(82)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == HerdParserCB_OPEN {
		{
			p.SetState(81)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IScalarContext)(nil)).Elem(), 0)

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
		p.SetState(84)

		var _m = p.Match(HerdParserIDENTIFIER)

		localctx.(*FilterContext).key = _m
	}
	p.SetState(89)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case HerdParserEQUALS, HerdParserNOT_EQUALS:
		{
			p.SetState(85)

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
			p.SetState(86)

			var _x = p.Scalar()

			localctx.(*FilterContext).val = _x
		}

	case HerdParserMATCHES, HerdParserNOT_MATCHES:
		{
			p.SetState(87)

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
			p.SetState(88)

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
		p.SetState(91)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<HerdParserDURATION)|(1<<HerdParserNUMBER)|(1<<HerdParserIDENTIFIER)|(1<<HerdParserSTRING))) != 0) {
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IScalarContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IScalarContext)
}

func (s *ValueContext) Array() IArrayContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayContext)
}

func (s *ValueContext) Hash() IHashContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHashContext)(nil)).Elem(), 0)

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

	p.SetState(96)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case HerdParserDURATION, HerdParserNUMBER, HerdParserIDENTIFIER, HerdParserSTRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(93)
			p.Scalar()
		}

	case HerdParserSB_OPEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(94)
			p.Array()
		}

	case HerdParserCB_OPEN:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(95)
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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IValueContext)(nil)).Elem())
	var tst = make([]IValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IValueContext)
		}
	}

	return tst
}

func (s *ArrayContext) Value(i int) IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), i)

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
	p.SetState(111)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(98)
			p.Match(HerdParserSB_OPEN)
		}
		{
			p.SetState(99)
			p.Match(HerdParserT__1)
		}

	case 2:
		{
			p.SetState(100)
			p.Match(HerdParserSB_OPEN)
		}
		{
			p.SetState(101)
			p.Value()
		}
		p.SetState(106)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserT__2 {
			{
				p.SetState(102)
				p.Match(HerdParserT__2)
			}
			{
				p.SetState(103)
				p.Value()
			}

			p.SetState(108)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(109)
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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IValueContext)(nil)).Elem())
	var tst = make([]IValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IValueContext)
		}
	}

	return tst
}

func (s *HashContext) Value(i int) IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), i)

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
	p.SetState(130)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(113)
			p.Match(HerdParserCB_OPEN)
		}
		{
			p.SetState(114)
			p.Match(HerdParserT__3)
		}

	case 2:
		{
			p.SetState(115)
			p.Match(HerdParserCB_OPEN)
		}
		{
			p.SetState(116)
			p.Match(HerdParserIDENTIFIER)
		}
		{
			p.SetState(117)
			p.Match(HerdParserT__4)
		}
		{
			p.SetState(118)
			p.Value()
		}
		p.SetState(125)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == HerdParserT__2 {
			{
				p.SetState(119)
				p.Match(HerdParserT__2)
			}
			{
				p.SetState(120)
				p.Match(HerdParserIDENTIFIER)
			}
			{
				p.SetState(121)
				p.Match(HerdParserT__4)
			}
			{
				p.SetState(122)
				p.Value()
			}

			p.SetState(127)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(128)
			p.Match(HerdParserT__3)
		}

	}

	return localctx
}
