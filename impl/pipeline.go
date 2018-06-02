// Base types for a pipeline.

package xpipe

import (
)

// A process factory
type ProcessFactory func() Process

type ProcessSink interface {
    // Accepts data from a process
    Accept(ctx *ProcessContext, d Datum) error

    // Sends an open signal
    Open(ctx *ProcessContext) error

    // Sends a close signal
    Close(ctx *ProcessContext) error
}



// A single process.
type Process interface {

    // Configures the process using the arguments from pipeline definition
    Config(args []ConfigArg) error

    // Called when the process is opened
    Open(ctx *ProcessContext, sink ProcessSink) error

    // Called when the process is closed
    Close(ctx *ProcessContext, sink ProcessSink) error

    // Applies the process with the specific datum.
    Apply(ctx *ProcessContext, in Datum, sink ProcessSink) error
}


// Utility method for safely sending something to a sink
func SendToSink(sink ProcessSink, ctx *ProcessContext, d Datum) error {
    if sink != nil {
        return sink.Accept(ctx, d)
    } else {
        return nil
    }
}

func SendOpen(sink ProcessSink, ctx *ProcessContext) error {
    if sink != nil {
        return sink.Open(ctx)
    } else {
        return nil
    }
}

func SendClose(sink ProcessSink, ctx *ProcessContext) error {
    if sink != nil {
        return sink.Close(ctx)
    } else {
        return nil
    }
}

// ----------------------------------------------------------------------

// A pipline chain
type PipelineChain struct {
    Process         Process
    Next            ProcessSink
}

// Implementation of the ProcessSink interface.  This will forward the data to
// the next process in the chain if one is defined.
func (pc *PipelineChain) Accept(ctx *ProcessContext, out Datum) error {
    return pc.Process.Apply(ctx, out, pc.Next)
}

func (pc *PipelineChain) Open(ctx *ProcessContext) error {
    return pc.Process.Open(ctx, pc.Next)
}

func (pc *PipelineChain) Close(ctx *ProcessContext) error {
    return pc.Process.Close(ctx, pc.Next)
}

// ----------------------------------------------------------------------

// A pipeline
type Pipeline struct {
    Start           *PipelineChain
    End             *PipelineChain
}

// Creates a new pipeline
func NewPipeline() *Pipeline {
    return &Pipeline{nil, nil}
}

// Appends a process to the end of the pipeline.
func (p *Pipeline) Append(proc Process) {
    pc := &PipelineChain{proc, nil}

    if (p.Start == nil) && (p.End == nil) {
        // Pipeline is empty
        p.Start = pc
        p.End = pc
    } else if (p.Start != nil) && (p.End != nil) {
        p.End.Next = pc
        p.End = pc
    } else {
        panic("Invariant violated: either start or end is not nil")
    }
}

// Appends the elements of another pipeline to this pipeline.  This does
// not modify the original pipeline.
//
// TODO: Redo this.  This is stupid!!
func (p *Pipeline) AppendPipeline(other *Pipeline) {
    var isChain bool

    if other == nil {
        return
    }

    pc := other.Start
    for pc != nil {
        p.Append(pc.Process)
        if pc, isChain = pc.Next.(*PipelineChain) ; !isChain {
            break
        }
    }
}

// Prepends a process to the start of the pipeline.
func (p *Pipeline) Prepend(proc Process) {
    pc := &PipelineChain{proc, nil}

    if (p.Start == nil) && (p.End == nil) {
        // Pipeline is empty
        p.Start = pc
        p.End = pc
    } else if (p.Start != nil) && (p.End != nil) {
        pc.Next = p.Start
        p.Start = pc
    } else {
        panic("Invariant violated: either start or end is not nil")
    }
}

// Opens the pipeline
func (p *Pipeline) Open(ctx *ProcessContext) error {
    if p.Start != nil {
        return p.Start.Open(ctx)
    } else {
        return nil
    }
}

// Closes the pipeline
func (p *Pipeline) Close(ctx *ProcessContext) error {
    if p.Start != nil {
        return p.Start.Close(ctx)
    } else {
        return nil
    }
}

// Sends a datum to the pipeline
func (p *Pipeline) Accept(ctx *ProcessContext, d Datum) error {
    if p.Start != nil {
        return p.Start.Accept(ctx, d)
    } else {
        return nil
    }
}

// Starts the pipeline
func (p *Pipeline) WithDatum(ctx *ProcessContext, d Datum) error {
    if p.Start != nil {
        var err error

        err = p.Open(ctx)
        if err != nil {
            return err
        }

        err = p.Accept(ctx, d)
        if err != nil {
            return err
        }

        err = p.Close(ctx)
        if err != nil {
            return err
        }
    }

    return nil
}
