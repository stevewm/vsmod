# Vintage Story Mod Manager (`vsmod`)

A command-line tool designed to manage [mods](https://mods.vintagestory.at/) for [Vintage Story](https://www.vintagestory.at/), a wilderness survival sandbox game, using a simple config file.

## Features

### Download mods

![download command](./download.gif)

### List mods

![list command](./list.gif)

### Hooks

Pre and post-run hooks can be used to run arbitrary commands before or after a vsmod command runs. For example, you can use a pre-run hook to clear the mods directory before downloading new mods. Hooks also support templating values from the config file. See the documentation for Go's [text/template](https://pkg.go.dev/text/template) package for how to use templating.

## Installation

Download pre-built binaries for Windows, macOS and Linux on the [Releases](https://github.com/stevewm/homelab/releases) page.

## Usage

Run `vsmod --help` to see the available commands and their options.

An example configuration file can be found in [examples/](./examples/mods.yaml).

## Container Image

vsmod is also built as a container image [here](https://github.com/stevewm/vsmod/pkgs/container/vsmod).

## Planned Features

- Mod version updating

## Support

Please raise an issue if you need help or have a feature request. Even better, raise a pull request if you're able.
