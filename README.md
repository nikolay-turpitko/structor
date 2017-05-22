[![Build Status](https://travis-ci.org/nikolay-turpitko/structor.svg?branch=master)](https://travis-ci.org/nikolay-turpitko/structor)

[![Travis CI](https://img.shields.io/travis/nikolay-turpitko/structor/master.svg?style=flat-square)](https://travis-ci.org/nikolay-turpitko/structor)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/nikolay-turpitko/structor/blob/master/LICENSE.md)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/nikolay-turpitko/structor)
[![Coverage Status](http://codecov.io/github/nikolay-turpitko/structor/coverage.svg?branch=master)](http://codecov.io/github/nikolay-turpitko/structor?branch=master)

# structor
Uses EL in Go struct tags to compute struct fields.

## In a nutshell

Basic idea is to use simple expression language (EL) within Go struct tags to
compute struct fields based on other fields or provided additional context.

It can be viewed/used as an advanced struct builder or constructor. Hence the
name. Another possible use case is a "struct walker" or "struct field visitor".

Implementation uses reflection and EL interpretation, so, it is relatively slow
by design, and, because of that, should be used for infrequent tasks. For
example, during program initialization. Or in cases when all other alternatives
are equally slow.

Initial implementation uses `"text/template"` as EL, it get fields name, value,
tags, full struct and extra context as a "dot" context and either use special
"set" custom function to save new field's value as object or just stores result
of the template evaluation as string into the field.

Map of the custom functions can be passed to EL engine for use in the
expressions.  "set" is a special custom EL function which stores argument into
the currently processed field. Special case: if annotated field is struct and
result can not be converted to it, then result stored into `.Sub`, and
evaluation executed recursively on the inner struct.

Simple namespaces and "import" mechanism provided for custom functions - whole
map of custom functions can be combined from several maps with (optional)
custom prefixes, prepended to every function name within map. Maps for provided
functions placed into separate Go packages and can be imported as needed. When
used with [Glide](https://github.com/Masterminds/glide), they can be imported
as subpackages, allowing to optimize dependencies.

It's simpler to illustrate on examples, see examples in tests and
[Godoc](http://godoc.org/github.com/nikolay-turpitko/structor).

Possible use cases:

- take a struct, filled up from configuration file using
  `github.com/spf13/viper` and compute additional fields:

  * replace referencies in settings to their values (for example base URL in
    othre URLs); also: templates in the settings themeselves, not in the tags
    of structure;
  * load short text files into string fields (with file names from config);
  * decode passwords;
  * extract data from environment variables, listed in config;
  * execute bash scripts to compute fields (iconv, openssl, etc);
  * parse complex types from string representation;

- use engines like regexp, xpath or goquery to extract pieces of data from
  text, xml, html etc. formats into fields (scraping data from html pages, text
  files or emails); for long and complex expressions it can be convenient to
  use multiline tags and process whole tag value as single expression;

- registering a custom EL "interpreter", which executes arbitrary custom logic
  in its Execute() method, it is possible to use structor as a "struct walker"
  or "field visitor", which will traverse structure fields and invoke custom
  function for every field, containing custom tag; it can be convenient with
  ability to use multiline tags;

## Ideas of functions, available in expressions

- [x] atoi
- [x] base64/unbase64
- [x] exec - invoke external process (shell, for instance)
- [x] encrypt/decrypt
- [x] env
- [ ] eval - get expression from other field (`eval "interpreter-name"
  .Extra.Expression`) - ? Use case: expressions in config file (arithmetics,
  string concatenation, etc).
- [x] fields
- [x] readFile - read file content to []byte or string
- [x] goquery
- [x] hex/unhex
- [x] match
- [x] math (basic arithmetics)
- [x] replace
- [x] rot13
- [x] set
- [x] split
- [x] standard for package `"text/template"`
- [x] trimSpace
- [x] xpath
- [x] ... (custom)

## Other ideas

- [ ] travis CI
- [ ] go expressions ("go/types".Eval()) as EL - ? Could be simpler than
  default EL interpreter.

## Godoc

[Godoc](http://godoc.org/github.com/nikolay-turpitko/structor)
