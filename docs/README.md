# Vintage Story Mod Manager (`vsmod`)

A command-line tool designed to manage [mods](https://mods.vintagestory.at/) for [Vintage Story](https://www.vintagestory.at/), a wilderness survival sandbox game, using a simple config file.

## Features

### Download mods

![download command](./download.gif)

### List mods

![list command](./list.gif)

---

## Installation

Download pre-built binaries for Windows, macOS and Linux on the [Releases](https://github.com/stevewm/homelab/releases) page.

## Usage

Run `vsmod --help` to see the available commands and their options.

### Config file

Example config file:

```yaml
# mods.yaml
game_version: 1.20.10 # the game version to check mod compatibility against
mods_dir: /app/VintageStory/Mods # can be relative or absolute (default: $PWD/Mods)
mods:
  - id: primitivesurvival # the path on the mod url on moddb
    version: 3.7.6 # the mod version (without the v)
  - id: ExpandedFoods
    version: 1.7.4
  - id: ACulinaryArtillery
    version: 1.2.5
  - id: CarryOn
    version: 1.8.0-rc.3
  - id: BetterRuins
    version: 0.4.9
  - id: th3dungeon
    version: 0.4.2
```

## Container Image

vsmod is also built as a container image [here](https://github.com/stevewm/vsmod/pkgs/container/vsmod).

## Planned Features

- Mod version updating

## Support

Please raise an issue if you need help or have a feature request. Even better, raise a pull request if you're able.
