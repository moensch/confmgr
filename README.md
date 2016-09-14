# confmgr

[![Build Status](https://travis-ci.org/moensch/confmgr.svg?branch=master)](https://travis-ci.org/moensch/confmgr)

`confmgr` is a minimalist key-value store REST front-end supporting configuration management needs. It uses Redis as
its primary back end.

Why not just use plain redis, I hear you ask?

confmgr supports dynamic scope-based variable lookups. You may need a different value for a config option
on different containers in different data centers. confmgr can handle that.

## Building

```
go get github.com/moensch/confmgr/cmd/confmgr
```
