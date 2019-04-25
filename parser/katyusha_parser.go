// Code generated from Katyusha.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // Katyusha

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 19, 79, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7, 4,
	8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 3, 2, 7, 2, 22, 10, 2, 12, 2, 14, 2,
	25, 11, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 34, 10, 3, 3,
	3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 5, 5, 43, 10, 5, 3, 5, 3, 5, 3,
	6, 3, 6, 3, 6, 3, 6, 7, 6, 51, 10, 6, 12, 6, 14, 6, 54, 11, 6, 3, 7, 3,
	7, 3, 7, 3, 7, 7, 7, 60, 10, 7, 12, 7, 14, 7, 63, 11, 7, 3, 8, 3, 8, 3,
	8, 5, 8, 68, 10, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 75, 10, 9, 3, 10,
	3, 10, 3, 10, 2, 2, 11, 2, 4, 6, 8, 10, 12, 14, 16, 18, 2, 3, 4, 2, 11,
	12, 17, 17, 2, 80, 2, 23, 3, 2, 2, 2, 4, 33, 3, 2, 2, 2, 6, 37, 3, 2, 2,
	2, 8, 39, 3, 2, 2, 2, 10, 46, 3, 2, 2, 2, 12, 55, 3, 2, 2, 2, 14, 64, 3,
	2, 2, 2, 16, 69, 3, 2, 2, 2, 18, 76, 3, 2, 2, 2, 20, 22, 5, 4, 3, 2, 21,
	20, 3, 2, 2, 2, 22, 25, 3, 2, 2, 2, 23, 21, 3, 2, 2, 2, 23, 24, 3, 2, 2,
	2, 24, 26, 3, 2, 2, 2, 25, 23, 3, 2, 2, 2, 26, 27, 7, 2, 2, 3, 27, 3, 3,
	2, 2, 2, 28, 34, 5, 6, 4, 2, 29, 34, 5, 8, 5, 2, 30, 34, 5, 10, 6, 2, 31,
	34, 5, 12, 7, 2, 32, 34, 5, 14, 8, 2, 33, 28, 3, 2, 2, 2, 33, 29, 3, 2,
	2, 2, 33, 30, 3, 2, 2, 2, 33, 31, 3, 2, 2, 2, 33, 32, 3, 2, 2, 2, 33, 34,
	3, 2, 2, 2, 34, 35, 3, 2, 2, 2, 35, 36, 7, 3, 2, 2, 36, 5, 3, 2, 2, 2,
	37, 38, 7, 4, 2, 2, 38, 7, 3, 2, 2, 2, 39, 40, 7, 5, 2, 2, 40, 42, 7, 13,
	2, 2, 41, 43, 7, 15, 2, 2, 42, 41, 3, 2, 2, 2, 42, 43, 3, 2, 2, 2, 43,
	44, 3, 2, 2, 2, 44, 45, 5, 18, 10, 2, 45, 9, 3, 2, 2, 2, 46, 47, 7, 6,
	2, 2, 47, 48, 7, 9, 2, 2, 48, 52, 7, 14, 2, 2, 49, 51, 5, 16, 9, 2, 50,
	49, 3, 2, 2, 2, 51, 54, 3, 2, 2, 2, 52, 50, 3, 2, 2, 2, 52, 53, 3, 2, 2,
	2, 53, 11, 3, 2, 2, 2, 54, 52, 3, 2, 2, 2, 55, 56, 7, 7, 2, 2, 56, 57,
	7, 9, 2, 2, 57, 61, 7, 14, 2, 2, 58, 60, 5, 16, 9, 2, 59, 58, 3, 2, 2,
	2, 60, 63, 3, 2, 2, 2, 61, 59, 3, 2, 2, 2, 61, 62, 3, 2, 2, 2, 62, 13,
	3, 2, 2, 2, 63, 61, 3, 2, 2, 2, 64, 65, 7, 8, 2, 2, 65, 67, 7, 9, 2, 2,
	66, 68, 7, 10, 2, 2, 67, 66, 3, 2, 2, 2, 67, 68, 3, 2, 2, 2, 68, 15, 3,
	2, 2, 2, 69, 74, 7, 13, 2, 2, 70, 71, 7, 15, 2, 2, 71, 75, 5, 18, 10, 2,
	72, 73, 7, 16, 2, 2, 73, 75, 7, 18, 2, 2, 74, 70, 3, 2, 2, 2, 74, 72, 3,
	2, 2, 2, 75, 17, 3, 2, 2, 2, 76, 77, 9, 2, 2, 2, 77, 19, 3, 2, 2, 2, 9,
	23, 33, 42, 52, 61, 67, 74,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'\n'", "", "'set'", "'add'", "'remove'", "'list'", "'hosts'", "'--oneline'",
	"", "", "", "", "'='", "'=~'",
}
var symbolicNames = []string{
	"", "", "RUN", "SET", "ADD", "REMOVE", "LIST", "HOSTS", "ONELINE", "DURATION",
	"NUMBER", "IDENTIFIER", "GLOB", "EQUALS", "MATCHES", "STRING", "REGEXP",
	"SKIP_",
}

var ruleNames = []string{
	"prog", "line", "run", "set", "add", "remove", "list", "filter", "value",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type KatyushaParser struct {
	*antlr.BaseParser
}

func NewKatyushaParser(input antlr.TokenStream) *KatyushaParser {
	this := new(KatyushaParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Katyusha.g4"

	return this
}

// KatyushaParser tokens.
const (
	KatyushaParserEOF        = antlr.TokenEOF
	KatyushaParserT__0       = 1
	KatyushaParserRUN        = 2
	KatyushaParserSET        = 3
	KatyushaParserADD        = 4
	KatyushaParserREMOVE     = 5
	KatyushaParserLIST       = 6
	KatyushaParserHOSTS      = 7
	KatyushaParserONELINE    = 8
	KatyushaParserDURATION   = 9
	KatyushaParserNUMBER     = 10
	KatyushaParserIDENTIFIER = 11
	KatyushaParserGLOB       = 12
	KatyushaParserEQUALS     = 13
	KatyushaParserMATCHES    = 14
	KatyushaParserSTRING     = 15
	KatyushaParserREGEXP     = 16
	KatyushaParserSKIP_      = 17
)

// KatyushaParser rules.
const (
	KatyushaParserRULE_prog   = 0
	KatyushaParserRULE_line   = 1
	KatyushaParserRULE_run    = 2
	KatyushaParserRULE_set    = 3
	KatyushaParserRULE_add    = 4
	KatyushaParserRULE_remove = 5
	KatyushaParserRULE_list   = 6
	KatyushaParserRULE_filter = 7
	KatyushaParserRULE_value  = 8
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
	p.RuleIndex = KatyushaParserRULE_prog
	return p
}

func (*ProgContext) IsProgContext() {}

func NewProgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgContext {
	var p = new(ProgContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_prog

	return p
}

func (s *ProgContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgContext) EOF() antlr.TerminalNode {
	return s.GetToken(KatyushaParserEOF, 0)
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
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterProg(s)
	}
}

func (s *ProgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitProg(s)
	}
}

func (p *KatyushaParser) Prog() (localctx IProgContext) {
	localctx = NewProgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, KatyushaParserRULE_prog)
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
	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<KatyushaParserT__0)|(1<<KatyushaParserRUN)|(1<<KatyushaParserSET)|(1<<KatyushaParserADD)|(1<<KatyushaParserREMOVE)|(1<<KatyushaParserLIST))) != 0 {
		{
			p.SetState(18)
			p.Line()
		}

		p.SetState(23)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(24)
		p.Match(KatyushaParserEOF)
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
	p.RuleIndex = KatyushaParserRULE_line
	return p
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_line

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
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterLine(s)
	}
}

func (s *LineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitLine(s)
	}
}

func (p *KatyushaParser) Line() (localctx ILineContext) {
	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, KatyushaParserRULE_line)

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
	p.SetState(31)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case KatyushaParserRUN:
		{
			p.SetState(26)
			p.Run()
		}

	case KatyushaParserSET:
		{
			p.SetState(27)
			p.Set()
		}

	case KatyushaParserADD:
		{
			p.SetState(28)
			p.Add()
		}

	case KatyushaParserREMOVE:
		{
			p.SetState(29)
			p.Remove()
		}

	case KatyushaParserLIST:
		{
			p.SetState(30)
			p.List()
		}

	case KatyushaParserT__0:

	default:
	}
	{
		p.SetState(33)
		p.Match(KatyushaParserT__0)
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
	p.RuleIndex = KatyushaParserRULE_run
	return p
}

func (*RunContext) IsRunContext() {}

func NewRunContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RunContext {
	var p = new(RunContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_run

	return p
}

func (s *RunContext) GetParser() antlr.Parser { return s.parser }

func (s *RunContext) RUN() antlr.TerminalNode {
	return s.GetToken(KatyushaParserRUN, 0)
}

func (s *RunContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RunContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RunContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterRun(s)
	}
}

func (s *RunContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitRun(s)
	}
}

func (p *KatyushaParser) Run() (localctx IRunContext) {
	localctx = NewRunContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, KatyushaParserRULE_run)

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
		p.SetState(35)
		p.Match(KatyushaParserRUN)
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
	GetVarvalue() IValueContext

	// SetVarvalue sets the varvalue rule contexts.
	SetVarvalue(IValueContext)

	// IsSetContext differentiates from other interfaces.
	IsSetContext()
}

type SetContext struct {
	*antlr.BaseParserRuleContext
	parser   antlr.Parser
	varname  antlr.Token
	varvalue IValueContext
}

func NewEmptySetContext() *SetContext {
	var p = new(SetContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = KatyushaParserRULE_set
	return p
}

func (*SetContext) IsSetContext() {}

func NewSetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetContext {
	var p = new(SetContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_set

	return p
}

func (s *SetContext) GetParser() antlr.Parser { return s.parser }

func (s *SetContext) GetVarname() antlr.Token { return s.varname }

func (s *SetContext) SetVarname(v antlr.Token) { s.varname = v }

func (s *SetContext) GetVarvalue() IValueContext { return s.varvalue }

func (s *SetContext) SetVarvalue(v IValueContext) { s.varvalue = v }

func (s *SetContext) SET() antlr.TerminalNode {
	return s.GetToken(KatyushaParserSET, 0)
}

func (s *SetContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(KatyushaParserIDENTIFIER, 0)
}

func (s *SetContext) Value() IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *SetContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(KatyushaParserEQUALS, 0)
}

func (s *SetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterSet(s)
	}
}

func (s *SetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitSet(s)
	}
}

func (p *KatyushaParser) Set() (localctx ISetContext) {
	localctx = NewSetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, KatyushaParserRULE_set)
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
		p.SetState(37)
		p.Match(KatyushaParserSET)
	}
	{
		p.SetState(38)

		var _m = p.Match(KatyushaParserIDENTIFIER)

		localctx.(*SetContext).varname = _m
	}
	p.SetState(40)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == KatyushaParserEQUALS {
		{
			p.SetState(39)
			p.Match(KatyushaParserEQUALS)
		}

	}
	{
		p.SetState(42)

		var _x = p.Value()

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
	p.RuleIndex = KatyushaParserRULE_add
	return p
}

func (*AddContext) IsAddContext() {}

func NewAddContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AddContext {
	var p = new(AddContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_add

	return p
}

func (s *AddContext) GetParser() antlr.Parser { return s.parser }

func (s *AddContext) GetGlob() antlr.Token { return s.glob }

func (s *AddContext) SetGlob(v antlr.Token) { s.glob = v }

func (s *AddContext) GetFilters() IFilterContext { return s.filters }

func (s *AddContext) SetFilters(v IFilterContext) { s.filters = v }

func (s *AddContext) ADD() antlr.TerminalNode {
	return s.GetToken(KatyushaParserADD, 0)
}

func (s *AddContext) HOSTS() antlr.TerminalNode {
	return s.GetToken(KatyushaParserHOSTS, 0)
}

func (s *AddContext) GLOB() antlr.TerminalNode {
	return s.GetToken(KatyushaParserGLOB, 0)
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
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterAdd(s)
	}
}

func (s *AddContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitAdd(s)
	}
}

func (p *KatyushaParser) Add() (localctx IAddContext) {
	localctx = NewAddContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, KatyushaParserRULE_add)
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
		p.SetState(44)
		p.Match(KatyushaParserADD)
	}
	{
		p.SetState(45)
		p.Match(KatyushaParserHOSTS)
	}
	{
		p.SetState(46)

		var _m = p.Match(KatyushaParserGLOB)

		localctx.(*AddContext).glob = _m
	}
	p.SetState(50)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == KatyushaParserIDENTIFIER {
		{
			p.SetState(47)

			var _x = p.Filter()

			localctx.(*AddContext).filters = _x
		}

		p.SetState(52)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.RuleIndex = KatyushaParserRULE_remove
	return p
}

func (*RemoveContext) IsRemoveContext() {}

func NewRemoveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RemoveContext {
	var p = new(RemoveContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_remove

	return p
}

func (s *RemoveContext) GetParser() antlr.Parser { return s.parser }

func (s *RemoveContext) GetGlob() antlr.Token { return s.glob }

func (s *RemoveContext) SetGlob(v antlr.Token) { s.glob = v }

func (s *RemoveContext) GetFilters() IFilterContext { return s.filters }

func (s *RemoveContext) SetFilters(v IFilterContext) { s.filters = v }

func (s *RemoveContext) REMOVE() antlr.TerminalNode {
	return s.GetToken(KatyushaParserREMOVE, 0)
}

func (s *RemoveContext) HOSTS() antlr.TerminalNode {
	return s.GetToken(KatyushaParserHOSTS, 0)
}

func (s *RemoveContext) GLOB() antlr.TerminalNode {
	return s.GetToken(KatyushaParserGLOB, 0)
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
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterRemove(s)
	}
}

func (s *RemoveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitRemove(s)
	}
}

func (p *KatyushaParser) Remove() (localctx IRemoveContext) {
	localctx = NewRemoveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, KatyushaParserRULE_remove)
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
		p.SetState(53)
		p.Match(KatyushaParserREMOVE)
	}
	{
		p.SetState(54)
		p.Match(KatyushaParserHOSTS)
	}
	{
		p.SetState(55)

		var _m = p.Match(KatyushaParserGLOB)

		localctx.(*RemoveContext).glob = _m
	}
	p.SetState(59)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == KatyushaParserIDENTIFIER {
		{
			p.SetState(56)

			var _x = p.Filter()

			localctx.(*RemoveContext).filters = _x
		}

		p.SetState(61)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IListContext is an interface to support dynamic dispatch.
type IListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOneline returns the oneline token.
	GetOneline() antlr.Token

	// SetOneline sets the oneline token.
	SetOneline(antlr.Token)

	// IsListContext differentiates from other interfaces.
	IsListContext()
}

type ListContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	oneline antlr.Token
}

func NewEmptyListContext() *ListContext {
	var p = new(ListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = KatyushaParserRULE_list
	return p
}

func (*ListContext) IsListContext() {}

func NewListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListContext {
	var p = new(ListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_list

	return p
}

func (s *ListContext) GetParser() antlr.Parser { return s.parser }

func (s *ListContext) GetOneline() antlr.Token { return s.oneline }

func (s *ListContext) SetOneline(v antlr.Token) { s.oneline = v }

func (s *ListContext) LIST() antlr.TerminalNode {
	return s.GetToken(KatyushaParserLIST, 0)
}

func (s *ListContext) HOSTS() antlr.TerminalNode {
	return s.GetToken(KatyushaParserHOSTS, 0)
}

func (s *ListContext) ONELINE() antlr.TerminalNode {
	return s.GetToken(KatyushaParserONELINE, 0)
}

func (s *ListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterList(s)
	}
}

func (s *ListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitList(s)
	}
}

func (p *KatyushaParser) List() (localctx IListContext) {
	localctx = NewListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, KatyushaParserRULE_list)
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
		p.SetState(62)
		p.Match(KatyushaParserLIST)
	}
	{
		p.SetState(63)
		p.Match(KatyushaParserHOSTS)
	}
	p.SetState(65)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == KatyushaParserONELINE {
		{
			p.SetState(64)

			var _m = p.Match(KatyushaParserONELINE)

			localctx.(*ListContext).oneline = _m
		}

	}

	return localctx
}

// IFilterContext is an interface to support dynamic dispatch.
type IFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFilterContext differentiates from other interfaces.
	IsFilterContext()
}

type FilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterContext() *FilterContext {
	var p = new(FilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = KatyushaParserRULE_filter
	return p
}

func (*FilterContext) IsFilterContext() {}

func NewFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterContext {
	var p = new(FilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_filter

	return p
}

func (s *FilterContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(KatyushaParserIDENTIFIER, 0)
}

func (s *FilterContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(KatyushaParserEQUALS, 0)
}

func (s *FilterContext) Value() IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *FilterContext) MATCHES() antlr.TerminalNode {
	return s.GetToken(KatyushaParserMATCHES, 0)
}

func (s *FilterContext) REGEXP() antlr.TerminalNode {
	return s.GetToken(KatyushaParserREGEXP, 0)
}

func (s *FilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterFilter(s)
	}
}

func (s *FilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitFilter(s)
	}
}

func (p *KatyushaParser) Filter() (localctx IFilterContext) {
	localctx = NewFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, KatyushaParserRULE_filter)

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
		p.SetState(67)
		p.Match(KatyushaParserIDENTIFIER)
	}
	p.SetState(72)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case KatyushaParserEQUALS:
		{
			p.SetState(68)
			p.Match(KatyushaParserEQUALS)
		}
		{
			p.SetState(69)
			p.Value()
		}

	case KatyushaParserMATCHES:
		{
			p.SetState(70)
			p.Match(KatyushaParserMATCHES)
		}
		{
			p.SetState(71)
			p.Match(KatyushaParserREGEXP)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.RuleIndex = KatyushaParserRULE_value
	return p
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = KatyushaParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(KatyushaParserNUMBER, 0)
}

func (s *ValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(KatyushaParserSTRING, 0)
}

func (s *ValueContext) DURATION() antlr.TerminalNode {
	return s.GetToken(KatyushaParserDURATION, 0)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(KatyushaListener); ok {
		listenerT.ExitValue(s)
	}
}

func (p *KatyushaParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, KatyushaParserRULE_value)
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
		p.SetState(74)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<KatyushaParserDURATION)|(1<<KatyushaParserNUMBER)|(1<<KatyushaParserSTRING))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}
