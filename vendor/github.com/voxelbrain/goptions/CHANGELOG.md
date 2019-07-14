# Changelog
## 2.5.6
### Bug fixes

* Unexported fields are now ignored

### Minor changes

* Examples for Verbs and Remainder in documentation

## 2.5.4
### Bugfixes

* Fix typo in documentation

## 2.5.3
### Bugfixes

* Remove placeholders from LICENSE
* Add CONTROBUTORS

## 2.5.2
### Bugfixes

* Bring `examples/readme_example.go` and `README.md` up to date
* Rewrite formatter

## 2.5.1
### Bugfixes

* Make arrays of `goptions.Marshaler` work

## 2.5.0
### New features

* Add support for `int32` and `int64`
* Add support for `float32` and `float64`

### Bugfixes

* Fix a bug where the name of a unknown type would not be properly
  printed
* Fix checks whether to use `os.Stdin` or `os.Stdout` when "-" is given for a
  `*os.File`
* Fix an test example where the output to `os.Stderr` is apparently
  not evaluated anymore.

## 2.4.1
### Bugfixes

* Code was not compilable due to temporary [maintainer](http://github.com/surma) idiocy
  (Thanks [akrennmair](http://github.com/akrennmair))

## 2.4.0
### New features

* Gave `goptions.FlagSet` a `ParseAndFail()` method

## 2.3.0
### New features

* Add support for `time.Duration`

## 2.2.0
### New features

* Add support for `*net.TCPAddr`
* Add support for `*net/url.URL`

### Bugfixes

* Fix behaviour of `[]bool` fields

## 2.1.0
### New features

* `goptions.Verbs` is of type `string` and will have selected verb name as value
  after parsing.

## 2.0.0
### Breaking changes

* Disallow multiple flag names for one member
* Remove `accumulate` option in favor of generic array support

### New features

* Add convenience function `ParseAndFail` to make common usage of the library
  a one-liner (see `readme_example.go`)
* Add a `Marshaler` interface to enable thrid-party types
* Add support for slices (and thereby for mutiple flag definitions)

### Minor changes

* Refactoring to get more flexibility
* Make a flag's default value accessible in the template context
