// Utilities

package xpipe

import (
)

//// first
//
//  Sends through the first datum encountered.  Filters all other datums out.
//
type FirstProcess struct {
    FoundDatum      bool
}

// Configures the process using the arguments from pipeline definition
func (p *FirstProcess) Config(args []ConfigArg) error {
    p.FoundDatum = false
    return nil
}

func (p *FirstProcess) Open(ctx *ProcessContext, sink ProcessSink) error {
    p.FoundDatum = false
    return SendOpen(sink, ctx)
}

func (p *FirstProcess) Close(ctx *ProcessContext, sink ProcessSink) error {
    return SendClose(sink, ctx)
}

// Applies the process with the specific datum.
func (p *FirstProcess) Apply(ctx *ProcessContext, in Datum, sink ProcessSink) error {
    if !p.FoundDatum {
        p.FoundDatum = true
        return SendToSink(sink, ctx, in)
    } else {
        return nil
    }
}
