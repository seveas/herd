// Code generated from Herd.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // Herd

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseHerdListener is a complete listener for a parse tree produced by HerdParser.
type BaseHerdListener struct{}

var _ HerdListener = &BaseHerdListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseHerdListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseHerdListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseHerdListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseHerdListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProg is called when production prog is entered.
func (s *BaseHerdListener) EnterProg(ctx *ProgContext) {}

// ExitProg is called when production prog is exited.
func (s *BaseHerdListener) ExitProg(ctx *ProgContext) {}

// EnterLine is called when production line is entered.
func (s *BaseHerdListener) EnterLine(ctx *LineContext) {}

// ExitLine is called when production line is exited.
func (s *BaseHerdListener) ExitLine(ctx *LineContext) {}

// EnterRun is called when production run is entered.
func (s *BaseHerdListener) EnterRun(ctx *RunContext) {}

// ExitRun is called when production run is exited.
func (s *BaseHerdListener) ExitRun(ctx *RunContext) {}

// EnterSet is called when production set is entered.
func (s *BaseHerdListener) EnterSet(ctx *SetContext) {}

// ExitSet is called when production set is exited.
func (s *BaseHerdListener) ExitSet(ctx *SetContext) {}

// EnterAdd is called when production add is entered.
func (s *BaseHerdListener) EnterAdd(ctx *AddContext) {}

// ExitAdd is called when production add is exited.
func (s *BaseHerdListener) ExitAdd(ctx *AddContext) {}

// EnterRemove is called when production remove is entered.
func (s *BaseHerdListener) EnterRemove(ctx *RemoveContext) {}

// ExitRemove is called when production remove is exited.
func (s *BaseHerdListener) ExitRemove(ctx *RemoveContext) {}

// EnterList is called when production list is entered.
func (s *BaseHerdListener) EnterList(ctx *ListContext) {}

// ExitList is called when production list is exited.
func (s *BaseHerdListener) ExitList(ctx *ListContext) {}

// EnterFilter is called when production filter is entered.
func (s *BaseHerdListener) EnterFilter(ctx *FilterContext) {}

// ExitFilter is called when production filter is exited.
func (s *BaseHerdListener) ExitFilter(ctx *FilterContext) {}

// EnterValue is called when production value is entered.
func (s *BaseHerdListener) EnterValue(ctx *ValueContext) {}

// ExitValue is called when production value is exited.
func (s *BaseHerdListener) ExitValue(ctx *ValueContext) {}
