***************************
qsplit
***************************
Quoted splitting for Golang
===========================

Qsplit splits a byte-slice into its constituent non-whitespace chunks,
keeping quoted chunks together, in the manner of the shell (mostly).

::
   
    `foo bar baz`      -> {`foo`, `bar`, `baz`}
    `   foo \tbar baz` -> {`foo`, `bar`, `baz`}
    `'foo bar' baz`    -> {`foo bar`, `baz`}
    `a b ‹c d "e f"›`  -> {`a`, `b`, `c d "e f"`}
    `a b'cd e'f`       -> {`a`, `b'cd`, `e'f`}

See the `package doc <http://godoc.org/github.com/firepear/qsplit>`_ for more
information.

* Current version: 2.2.3 (2019-02-25)

* Install: :code:`go get github.com/firepear/qsplit`
