// Code generated from Herd.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // Herd

import "github.com/antlr/antlr4/runtime/Go/antlr"

// HerdListener is a complete listener for a parse tree produced by HerdParser.
type HerdListener interface {
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

	// EnterRemove is called when entering the remove production.
	EnterRemove(c *RemoveContext)

	// EnterList is called when entering the list production.
	EnterList(c *ListContext)

	// EnterFilter is called when entering the filter production.
	EnterFilter(c *FilterContext)

	// EnterScalar is called when entering the scalar production.
	EnterScalar(c *ScalarContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// EnterArray is called when entering the array production.
	EnterArray(c *ArrayContext)

	// EnterHash is called when entering the hash production.
	EnterHash(c *HashContext)

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

	// ExitRemove is called when exiting the remove production.
	ExitRemove(c *RemoveContext)

	// ExitList is called when exiting the list production.
	ExitList(c *ListContext)

	// ExitFilter is called when exiting the filter production.
	ExitFilter(c *FilterContext)

	// ExitScalar is called when exiting the scalar production.
	ExitScalar(c *ScalarContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)

	// ExitArray is called when exiting the array production.
	ExitArray(c *ArrayContext)

	// ExitHash is called when exiting the hash production.
	ExitHash(c *HashContext)
}
