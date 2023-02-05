<a name="readme-top"></a>

<br />
<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![Build Status][build-shield]][build-url]

<h3 align="center">Itun2socks</h3>
  <p align="center">
    <br />
    <a href="https://github.com/igoogolx/itun2socks/wiki"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://igoogolx.github.io/lux-dashboard/">View Demo</a>
    .
    <a href="https://github.com/igoogolx/itun2socks/issues">Report Bug</a>
    ·
    <a href="https://github.com/igoogolx/itun2socks/issues">Request Feature</a>
  </p>
</div>


[contributors-shield]: https://img.shields.io/github/contributors/igoogolx/itun2socks.svg
[contributors-url]: https://github.com/igoogolx/itun2socks/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/igoogolx/itun2socks.svg
[forks-url]: https://github.com/igoogolx/itun2socks/network/members
[stars-shield]: https://img.shields.io/github/stars/igoogolx/itun2socks.svg
[stars-url]: https://github.com/igoogolx/itun2socks/stargazers
[issues-shield]: https://img.shields.io/github/issues/igoogolx/itun2socks.svg
[issues-url]: https://github.com/igoogolx/itun2socks/issues
[license-shield]: https://img.shields.io/github/license/igoogolx/itun2socks.svg
[license-url]: https://github.com/igoogolx/itun2socks/blob/main/LICENSE
[build-shield]: https://github.com/igoogolx/itun2socks/actions/workflows/build.yml/badge.svg
[build-url]: https://github.com/igoogolx/itun2socks/actions/workflows/build.yml


### Build with lwip stack
`go build -v -ldflags="-s -w" -trimpath -tags="with_lwip"`

### Supported net stack
Build with tag `gvisor`,`lwip`,`system`

### Enable debug mode
Build with tag `debug`
