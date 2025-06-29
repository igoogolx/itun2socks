<a name="readme-top"></a>

<br />
<div align="center">

[![Report][report-shield]][report-url]
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL License][license-shield]][license-url]
[![Build Status][build-shield]][build-url]
[![Winget Version][winget-shield]][winget-url]

<h3 align="center">Itun2socks</h3>
The engine that powers the <a href="https://github.com/igoogolx/lux"><strong>lux</strong></a>.

  <p align="center">
    <br />
    <a href="https://github.com/igoogolx/itun2socks/issues">Report Bug</a>
    ·
    <a href="https://github.com/igoogolx/itun2socks/issues">Request Feature</a>
  </p>
</div>


[report-shield]: https://goreportcard.com/badge/github.com/igoogolx/itun2socks
[report-url]: https://goreportcard.com/report/github.com/igoogolx/itun2socks
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
[winget-shield]: https://img.shields.io/winget/v/igoogolx.itun2socks
[winget-url]: https://github.com/microsoft/winget-cli


### Install by winget(Windows only)

`winget install igoogolx.itun2socks`

### Build with gvisor stack

The release is built with gvisor stack by default, you can build it with the following command:

`go build -v -ldflags="-s -w" -trimpath -tags="with_gvisor"`

### Supported net stack
Build with tag `with_gvisor`,`with_lwip`,`with_system`

### Enable debug mode
Build with tag `debug`

### Debug with pprof
`go tool pprof -http :8080 http://localhost:9000/debug/pprof/heap`
