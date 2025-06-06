# Vintage Story Mod Manager (`vsmod`)

A command-line tool designed to manage [mods](https://mods.vintagestory.at/) for [Vintage Story](https://www.vintagestory.at/), a wilderness survival sandbox game, using a simple config file.

> [!WARNING]
> Early version. Expect bugs!

## Features

### Download mods

![download command](./download.gif)

#### Mod and Game version compatibility

Each time a new game version is released mod authors must manually mark their mods as compatible with it. Many don't bother to do this, breaking reliable compatibility checks.

To work around this, `game_version` supports [semver ranges](https://github.com/Masterminds/semver#checking-version-constraints). For example, setting it to `~1.20.0` will mean as long as a mod is marked as compatible for `1.20.0`, `1.20.1`, `1.20.2`, etc., it will be downloaded.

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
