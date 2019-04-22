// Code generated from Katyusha.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // Katyusha

import "github.com/antlr/antlr4/runtime/Go/antlr"

// KatyushaListener is a complete listener for a parse tree produced by KatyushaParser.
type KatyushaListener interface {
	antlr.ParseTreeListener

	// EnterProg is called when entering the prog production.
	EnterProg(c *ProgContext)

	// EnterLine is called when entering the line production.
	EnterLine(c *LineContext)

	// EnterRun is called when entering the run production.
	EnterRun(c *RunContext)

	// EnterSet is called when entering the set production.
	EnterSet(c *SetContext)

	// EnterAdd is called when entering the add production.
	EnterAdd(c *AddContext)

	// EnterFilter is called when entering the filter production.
	EnterFilter(c *FilterContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// ExitProg is called when exiting the prog production.
	ExitProg(c *ProgContext)

	// ExitLine is called when exiting the line production.
	ExitLine(c *LineContext)

	// ExitRun is called when exiting the run production.
	ExitRun(c *RunContext)

	// ExitSet is called when exiting the set production.
	ExitSet(c *SetContext)

	// ExitAdd is called when exiting the add production.
	ExitAdd(c *AddContext)

	// ExitFilter is called when exiting the filter production.
	ExitFilter(c *FilterContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)
}
