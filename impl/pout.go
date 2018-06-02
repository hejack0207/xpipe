// Pipelines that print out things

package xpipe

import (
    "fmt"
)

//// print
//
//  Prints the results of each datum to standard out.
//
//  ```
//  xpath "/xpath/result" | print
//  ```
type PrintProcess struct {
}

// Configures the process using the arguments from pipeline definition
func (p *PrintProcess) Config(args []ConfigArg) error {
    return nil
}

func (p *PrintProcess) Open(ctx *ProcessContext, sink ProcessSink) error {
    return SendOpen(sink, ctx)
}

func (p *PrintProcess) Close(ctx *ProcessContext, sink ProcessSink) error {
    return SendClose(sink, ctx)
}

// Applies the process with the specific datum.
func (p *PrintProcess) Apply(ctx *ProcessContext, in Datum, sink ProcessSink) error {
    switch in.(type) {
    case DocDatum:
        fmt.Print(in.String())
    default:
        fmt.Println(in.String())
    }
    return SendToSink(sink, ctx, in)
}

//// printfile
//
//  Prints the filename if at least one datum is encountered.
//
type PrintFileProcess struct {
    hasDatum        bool
}

// Configures the process using the arguments from pipeline definition
func (p *PrintFileProcess) Config(args []ConfigArg) error {
    return nil
}

func (p *PrintFileProcess) Open(ctx *ProcessContext, sink ProcessSink) error {
    p.hasDatum = false
    return SendOpen(sink, ctx)
}

func (p *PrintFileProcess) Close(ctx *ProcessContext, sink ProcessSink) error {
    if p.hasDatum {
        fmt.Println(ctx.Filename)
    }
    return SendClose(sink, ctx)
}

func (p *PrintFileProcess) Apply(ctx *ProcessContext, in Datum, sink ProcessSink) error {
    p.hasDatum = true
    return SendToSink(sink, ctx, in)
}

//// printemptyfile
//
//  Prints the filename if no datums are encountered.
//
type PrintEmptyFileProcess struct {
    hasDatum        bool
}

// Configures the process using the arguments from pipeline definition
func (p *PrintEmptyFileProcess) Config(args []ConfigArg) error {
    return nil
}

func (p *PrintEmptyFileProcess) Open(ctx *ProcessContext, sink ProcessSink) error {
    p.hasDatum = false
    return SendOpen(sink, ctx)
}

func (p *PrintEmptyFileProcess) Close(ctx *ProcessContext, sink ProcessSink) error {
    if !p.hasDatum {
        fmt.Println(ctx.Filename)
    }
    return SendClose(sink, ctx)
}

func (p *PrintEmptyFileProcess) Apply(ctx *ProcessContext, in Datum, sink ProcessSink) error {
    p.hasDatum = true
    return SendToSink(sink, ctx, in)
}
