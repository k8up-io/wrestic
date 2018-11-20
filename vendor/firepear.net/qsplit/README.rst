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

See the `package doc <http://godoc.org/firepear.net/qsplit>`_ for more
information.

* Current version: 2.2.2 (2016-03-30)

* Install: :code:`go get firepear.net/qsplit`

* `Coverage report <http://firepear.net/qsplit/coverage.html>`_

* `Github <http://github.com/firepear/qsplit/>`_

