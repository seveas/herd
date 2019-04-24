// Code generated from Katyusha.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // Katyusha

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseKatyushaListener is a complete listener for a parse tree produced by KatyushaParser.
type BaseKatyushaListener struct{}

var _ KatyushaListener = &BaseKatyushaListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseKatyushaListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseKatyushaListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseKatyushaListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseKatyushaListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProg is called when production prog is entered.
func (s *BaseKatyushaListener) EnterProg(ctx *ProgContext) {}

// ExitProg is called when production prog is exited.
func (s *BaseKatyushaListener) ExitProg(ctx *ProgContext) {}

// EnterLine is called when production line is entered.
func (s *BaseKatyushaListener) EnterLine(ctx *LineContext) {}

// ExitLine is called when production line is exited.
func (s *BaseKatyushaListener) ExitLine(ctx *LineContext) {}

// EnterRun is called when production run is entered.
func (s *BaseKatyushaListener) EnterRun(ctx *RunContext) {}

// ExitRun is called when production run is exited.
func (s *BaseKatyushaListener) ExitRun(ctx *RunContext) {}

// EnterSet is called when production set is entered.
func (s *BaseKatyushaListener) EnterSet(ctx *SetContext) {}

// ExitSet is called when production set is exited.
func (s *BaseKatyushaListener) ExitSet(ctx *SetContext) {}

// EnterAdd is called when production add is entered.
func (s *BaseKatyushaListener) EnterAdd(ctx *AddContext) {}

// ExitAdd is called when production add is exited.
func (s *BaseKatyushaListener) ExitAdd(ctx *AddContext) {}

// EnterRemove is called when production remove is entered.
func (s *BaseKatyushaListener) EnterRemove(ctx *RemoveContext) {}

// ExitRemove is called when production remove is exited.
func (s *BaseKatyushaListener) ExitRemove(ctx *RemoveContext) {}

// EnterList is called when production list is entered.
func (s *BaseKatyushaListener) EnterList(ctx *ListContext) {}

// ExitList is called when production list is exited.
func (s *BaseKatyushaListener) ExitList(ctx *ListContext) {}

// EnterFilter is called when production filter is entered.
func (s *BaseKatyushaListener) EnterFilter(ctx *FilterContext) {}

// ExitFilter is called when production filter is exited.
func (s *BaseKatyushaListener) ExitFilter(ctx *FilterContext) {}

// EnterValue is called when production value is entered.
func (s *BaseKatyushaListener) EnterValue(ctx *ValueContext) {}

// ExitValue is called when production value is exited.
func (s *BaseKatyushaListener) ExitValue(ctx *ValueContext) {}
