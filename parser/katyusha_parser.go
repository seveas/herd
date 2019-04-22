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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 14, 56, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7, 4,
	8, 9, 8, 3, 2, 7, 2, 18, 10, 2, 12, 2, 14, 2, 21, 11, 2, 3, 2, 3, 2, 3,
	3, 3, 3, 3, 3, 5, 3, 28, 10, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3,
	5, 5, 5, 37, 10, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 7, 6, 45, 10, 6,
	12, 6, 14, 6, 48, 11, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 2, 2,
	9, 2, 4, 6, 8, 10, 12, 14, 2, 3, 4, 2, 8, 9, 13, 13, 2, 54, 2, 19, 3, 2,
	2, 2, 4, 27, 3, 2, 2, 2, 6, 31, 3, 2, 2, 2, 8, 33, 3, 2, 2, 2, 10, 40,
	3, 2, 2, 2, 12, 49, 3, 2, 2, 2, 14, 53, 3, 2, 2, 2, 16, 18, 5, 4, 3, 2,
	17, 16, 3, 2, 2, 2, 18, 21, 3, 2, 2, 2, 19, 17, 3, 2, 2, 2, 19, 20, 3,
	2, 2, 2, 20, 22, 3, 2, 2, 2, 21, 19, 3, 2, 2, 2, 22, 23, 7, 2, 2, 3, 23,
	3, 3, 2, 2, 2, 24, 28, 5, 6, 4, 2, 25, 28, 5, 8, 5, 2, 26, 28, 5, 10, 6,
	2, 27, 24, 3, 2, 2, 2, 27, 25, 3, 2, 2, 2, 27, 26, 3, 2, 2, 2, 27, 28,
	3, 2, 2, 2, 28, 29, 3, 2, 2, 2, 29, 30, 7, 3, 2, 2, 30, 5, 3, 2, 2, 2,
	31, 32, 7, 4, 2, 2, 32, 7, 3, 2, 2, 2, 33, 34, 7, 5, 2, 2, 34, 36, 7, 10,
	2, 2, 35, 37, 7, 12, 2, 2, 36, 35, 3, 2, 2, 2, 36, 37, 3, 2, 2, 2, 37,
	38, 3, 2, 2, 2, 38, 39, 5, 14, 8, 2, 39, 9, 3, 2, 2, 2, 40, 41, 7, 6, 2,
	2, 41, 42, 7, 7, 2, 2, 42, 46, 7, 11, 2, 2, 43, 45, 5, 12, 7, 2, 44, 43,
	3, 2, 2, 2, 45, 48, 3, 2, 2, 2, 46, 44, 3, 2, 2, 2, 46, 47, 3, 2, 2, 2,
	47, 11, 3, 2, 2, 2, 48, 46, 3, 2, 2, 2, 49, 50, 7, 10, 2, 2, 50, 51, 7,
	12, 2, 2, 51, 52, 5, 14, 8, 2, 52, 13, 3, 2, 2, 2, 53, 54, 9, 2, 2, 2,
	54, 15, 3, 2, 2, 2, 6, 19, 27, 36, 46,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'\n'", "", "'set'", "'add'", "'hosts'", "", "", "", "", "'='",
}
var symbolicNames = []string{
	"", "", "RUN", "SET", "ADD", "HOSTS", "DURATION", "NUMBER", "IDENTIFIER",
	"GLOB", "EQUALS", "STRING", "SKIP_",
}

var ruleNames = []string{
	"prog", "line", "run", "set", "add", "filter", "value",
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
	KatyushaParserHOSTS      = 5
	KatyushaParserDURATION   = 6
	KatyushaParserNUMBER     = 7
	KatyushaParserIDENTIFIER = 8
	KatyushaParserGLOB       = 9
	KatyushaParserEQUALS     = 10
	KatyushaParserSTRING     = 11
	KatyushaParserSKIP_      = 12
)

// KatyushaParser rules.
const (
	KatyushaParserRULE_prog   = 0
	KatyushaParserRULE_line   = 1
	KatyushaParserRULE_run    = 2
	KatyushaParserRULE_set    = 3
	KatyushaParserRULE_add    = 4
	KatyushaParserRULE_filter = 5
	KatyushaParserRULE_value  = 6
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
	p.SetState(17)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<KatyushaParserT__0)|(1<<KatyushaParserRUN)|(1<<KatyushaParserSET)|(1<<KatyushaParserADD))) != 0 {
		{
			p.SetState(14)
			p.Line()
		}

		p.SetState(19)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(20)
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
	p.SetState(25)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case KatyushaParserRUN:
		{
			p.SetState(22)
			p.Run()
		}

	case KatyushaParserSET:
		{
			p.SetState(23)
			p.Set()
		}

	case KatyushaParserADD:
		{
			p.SetState(24)
			p.Add()
		}

	case KatyushaParserT__0:

	default:
	}
	{
		p.SetState(27)
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
		p.SetState(29)
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
		p.SetState(31)
		p.Match(KatyushaParserSET)
	}
	{
		p.SetState(32)

		var _m = p.Match(KatyushaParserIDENTIFIER)

		localctx.(*SetContext).varname = _m
	}
	p.SetState(34)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == KatyushaParserEQUALS {
		{
			p.SetState(33)
			p.Match(KatyushaParserEQUALS)
		}

	}
	{
		p.SetState(36)

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
		p.SetState(38)
		p.Match(KatyushaParserADD)
	}
	{
		p.SetState(39)
		p.Match(KatyushaParserHOSTS)
	}
	{
		p.SetState(40)

		var _m = p.Match(KatyushaParserGLOB)

		localctx.(*AddContext).glob = _m
	}
	p.SetState(44)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == KatyushaParserIDENTIFIER {
		{
			p.SetState(41)

			var _x = p.Filter()

			localctx.(*AddContext).filters = _x
		}

		p.SetState(46)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.EnterRule(localctx, 10, KatyushaParserRULE_filter)

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
		p.Match(KatyushaParserIDENTIFIER)
	}
	{
		p.SetState(48)
		p.Match(KatyushaParserEQUALS)
	}
	{
		p.SetState(49)
		p.Value()
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
	p.EnterRule(localctx, 12, KatyushaParserRULE_value)
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
		p.SetState(51)
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
