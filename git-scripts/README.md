# Git Scripts
A collection of scripts that are focused around Git.

#### tagging.go
_@contributors -- Mikey L, Raymond Jones._

Basic script that when passed a version code and a directory will get the current version and increment it. Currently only works with [Semantic Versioning](https://semver.org/).

##### Usage
```
./tagging.go -version=<Version> -dir=<ProjectDirectory>
```