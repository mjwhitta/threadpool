# gomk

## What is this?

A "simple" `Makefile` solution for Go projects. It attempts to
automate a lot of tasks for you. There are targets such as
`reportcard` which will "grade" your code and provide warnings where
necessary. The default target is `build` and it will compile
`./cmd/*`. However, it can be redefined in your top-level `Makefile`
if needed. See below for an example.

## How to install

Open a terminal, navigate to your Go repo, and run the following:

```
$ git submodule add https://gitlab.com/mjwhitta/gomk.git
```

## Usage

After adding gomk as a submodule, you can create a `Makefile` with the
following contents:

```
# Include gomk if it's been checked-out: git submodule update --init
-include gomk/main.mk

# Override the default build recipe
build: reportcard dir
    @go build -o "$(OUT)" ./cmd/justonething

# Override the default debug recipe with an empty recipe (@true is
# required)
debug: reportcard dir
    @true

# Add new recipes specific to your project
superclean: clean
	@echo "Clean up extra stuff"
```

## Links

- [Source](https://gitlab.com/mjwhitta/gomk)
