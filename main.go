package main

import (
    "os"
    "io"
    "bytes"
    "text/template"
    "strings"
    "fmt"
    
    flags "github.com/jessevdk/go-flags"
    log "github.com/Sirupsen/logrus"
)

var version string = "undef"

type Options struct {
    Debug  bool   `          long:"debug"  description:"enable debug"`
    Fail   bool   `          long:"fail"   description:"fail if var not found" default:"true"`
    Input  string `short:"i" long:"input"  description:"input file"            default:"-"`
    Output string `short:"o" long:"output" description:"output file"           default:"-"`
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
    
    // return the value of an environment variable. if an extra argument is
    // provided, it will be used as the default if there's no value for 'key'.
    getenv := func(key string, args ...string) (val string, err error) {
        val = os.Getenv(key)
        
        if val == "" {
            if len(args) > 0 {
                val = args[0]
            } else if opts.Fail {
                err = fmt.Errorf("no value for %s and no default provided", key)
            }
        }
        
        return
    }

    // map of functions to be provided to the template
    funcs := template.FuncMap{
        "env": getenv,
        "split": strings.Split,
    }
    
    // parse the template, with provided functions
    tmpl, err := template.New(opts.Input).Funcs(funcs).Parse(string(inBuf.Bytes()))
    checkError("unable to parse template", err)
    
    // execute the template
    err = tmpl.Execute(ofp, nil)
    checkError("error executing template", err)
}
