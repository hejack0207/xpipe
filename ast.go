// Pipeline definition file parser

package xpipe

// --------------------------------------------------------------------------
// Ast


// Top level script
type AstScript struct {
    Items       []AstItem
}

// Creates a new runtime from the script
func (a *AstScript) Configure(rt *Runtime) error {
    for _, i := range a.Items {
        if err := i.Config(rt) ; err != nil {
            return err
        }
    }

    return nil
}

// ---------------------------------------------------------------------------
//

type AstItem    interface {
    // Configure part of the runtime
    Config(rt *Runtime) error
}

type AstNamespaceMapping struct {
    Prefix      string
    Url         string
}

func (a *AstNamespaceMapping) Config(rt *Runtime) error {
    rt.AddNamespaceMapping(a.Prefix, a.Url)
    return nil
}

// ---------------------------------------------------------------------------
//

// A pipeline definition
type AstPipeline struct {
    Processes   *AstProcess
}

func (a *AstPipeline) Config(rt *Runtime) error {
    pl := NewPipeline()
    err := a.Processes.AddToPipeline(rt, pl)
    if err != nil {
        return err
    }

    rt.AddPipeline(pl)
    return nil
}

// ---------------------------------------------------------------------------
//

// A process invocation
type AstProcess struct {
    Name        string
    Args        []AstProcessArg
    Next        *AstProcess
}

func (a *AstProcess) AddToPipeline(rt *Runtime, pl *Pipeline) error {
    pargs := make([]ConfigArg, len(a.Args))
    for i, arg := range a.Args {
        pargs[i] = arg.ToConfigArg()
    }

    p, err := rt.Registry.NewProcess(a.Name, pargs)
    if err != nil {
        return err
    }

    pl.Append(p)
    if a.Next != nil {
        return a.Next.AddToPipeline(rt, pl)
    } else {
        return nil
    }
}

// -------------------------------------------------------------------------
//

// Ast Process arguments
type AstProcessArg interface {
    // Converts the argument into a ConfigArg
    ToConfigArg()   ConfigArg
}

// A literal argument
type AstLiteralProcessArg struct {
    L               ConfigArg
}

func (a *AstLiteralProcessArg) ToConfigArg() ConfigArg {
    return a.L
}
