# Processes

### first

```
first
```

Sends through the first datum encountered.  Filters all other datums out.

### print

```
print
```

Prints the results of each datum to standard out.

```
xpath "/xpath/result" | print
```

### printemptyfile

```
printemptyfile
```

Prints the filename if no datums are encountered.

### printfile

```
printfile
```

Prints the filename if at least one datum is encountered.

### settext

```
settext <newvalue>
```

Changes the contents of nodes and document datums to *newvalue*.  All other datums
are sent to the sink unchanged.

### thisdoc

```
thisdoc
```

Selects the document of nodes and sends them to the sink.  Any documents encountered
will pass through unmodified.  All other datums are filtered out.

### xpath

```
xpath <expr>
```

Selects nodes based on an XPath expression and sends the results to the sink.  If the XPath
expression produces scalars (e.g. string, numbers, booleans), these are sent as single datums
to the sink as well.

XPaths are executed over incoming node and document datums.  All other incomming datums are
filtered out.

