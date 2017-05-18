# structor
Uses EL in Go struct tags to compute struct fields.

Status: work in progress, not ready for use.

## In a nutshell

Basic idea is to use simple expression language within Go struct tags to
compute struct fields based on other fields or provided additional context.

It uses reflection and EL interpretation, so, it is relatively slow by design,
and, because of that, should be used for infrequent tasks. For example, during
program initialization. Or in cases when all other alternatives are equally
slow.

Initial implementation uses `"text/template"` as EL, it get fields name, value,
tags, full struct and extra context as a "dot" context and either use special
"set" custom function to save new field's value as object or just stores result
of the template evaluation as string into the field.

Map of the custom functions can be passed to EL engine for use in the
expressions.  "set" is a special custom EL function which stores argument into
the currently processed field. Special case: if annotated field is struct and
result can not be converted to it, then result stored into `.Sub`, and
evaluation executed recursively on the inner struct.

It's simpler to show on examples, see examples in tests and
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
- [ ] bash - should invoke bash and return struct with exit code, stdout,
  stderr
- [ ] encrypt/decrypt - (too many possible variants) ?
- [x] env
- [ ] eval - get expression from other field (`eval .Extra.Expression`) - ?
- [x] fields
- [x] readFile - read file content to []byte or string
- [ ] goquery
- [x] hex/unhex
- [x] match
- [x] math (basic arithmetics)
- [x] replace
- [x] rot13
- [x] set
- [x] split
- [x] standard for package `"text/template"`
- [x] trimSpace
- [ ] xpath
- [x] ... (custom)

## Other ideas

- [ ] errors with stack traces to debug custom functions
- [ ] godoc in separate doc.go file
- [ ] travis CI
- [ ] go expressions ("go/types".Eval()) as EL - ?

## Godoc

[Godoc](http://godoc.org/github.com/nikolay-turpitko/structor)
