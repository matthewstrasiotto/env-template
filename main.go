package main

import (
    "os"
    "io"
    "bytes"
    "text/template"
    "strings"
    
    flags "github.com/jessevdk/go-flags"
    log "github.com/Sirupsen/logrus"
)

var version string = "undef"

type Options struct {
    Debug  bool   `          long:"debug"  description:"enable debug"`
    Input  string `short:"i" long:"input"  description:"input file"  default:"-"`
    Output string `short:"o" long:"output" description:"output file" default:"-"`
}

func main() {
    log.SetOutput(os.Stderr)

    var opts Options
    
    _, err := flags.Parse(&opts)
    if err != nil {
        os.Exit(1)
    }
    
    if opts.Debug {
        log.SetLevel(log.DebugLevel)
    }
    
    log.Debugf("input file: %v", opts.Input)
    log.Debugf("output file: %v", opts.Output)
    
    // open the input stream
    var ifp *os.File
    if opts.Input == "-" {
        ifp = os.Stdin
    } else {
        ifp, err = os.Open(opts.Input)
        checkError("unable to open input", err)
        defer ifp.Close()
    }
    
    // open the output stream
    var ofp *os.File
    if opts.Output == "-" {
        ofp = os.Stdout
    } else {
        ofp, err = os.Create(opts.Output)
        checkError("unable to open output", err)
        defer ofp.Close()
    }
    
    // read in all input data
    inBuf := bytes.NewBuffer(nil)
    io.Copy(inBuf, ifp)
    
    // map of functions to be provided to the template
    funcs := template.FuncMap{
        "env": os.Getenv,
        "split": strings.Split,
    }
    
    // parse the template, with provided functions
    tmpl, err := template.New("env").Funcs(funcs).Parse(string(inBuf.Bytes()))
    checkError("unable to parse template", err)
    
    // execute the template
    err = tmpl.Execute(ofp, nil)
    checkError("error executing template", err)
}
