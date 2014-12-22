# simple env var-driven template utility

I use Docker a lot and passing configuration through environment variables is
awesome, but not all applications support that kind of configuration.  I've
previously used python, sed, and `eval`'ing files in the context of a shell
script.  I got tired of reinventing the wheel (and having to figure out how I
did it last time), so I wrote this.

## usage

    # env-template < /path/to/template > /path/to/generated-file
    # env-template -i /path/to/template > /path/to/generated-file
    # env-template -i /path/to/template -o /path/to/generated-file

## supported functions

* `env` -- retrieves an environment variable
* `split` -- splits a string on a substring; [`strings#Split`](http://golang.org/pkg/strings/#Split)

## unsupported functions (future?)

* `envs` -- return a map of all environment variable names; needs some regex love for filtering

## why not use reuse something else?

Several other utilities perform similar functionality.  [`confd`][confd] and
[`consul-template`][consul-template] are both able to render templates from
environment variables.  `consul-template` requires a connection to a consul
agent.  `confd` supports etcd, consul and environment variable backends, but
requires a configuration file.  I wanted something that has less overhead and a
simpler command-line interface.

[consul-template]: https://github.com/hashicorp/consul-template/
[confd]: https://github.com/kelseyhightower/confd/
