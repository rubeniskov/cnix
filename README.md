<img src="./docs/img/logo.svg">

A simple and extensible tool for makers to control [3D printers](https://en.wikipedia.org/wiki/3D_printing) and [CNC routers](https://en.wikipedia.org/wiki/CNC_router).


> Pronunced /tsːnɪks/ came from acronyms CNC + UNIX. The logo is composed by triangles using a 3D cube  as guide, this is motivated due the volumen of machining workarea. The X also describes the cartesian system of coordinates which is responsible of the movement through axis.

## Motivation 

For some reason [3D printers](https://en.wikipedia.org/wiki/3D_printing) and [CNC routers](https://en.wikipedia.org/wiki/CNC_router) have been separated paths in software development, providing very competent tools in both areas, but wasting such usefull functionality between them, increasing the learning curve for those who own both or want to getting started. 

There are good news, both machines are not so different, the have many shared functionalities. Even the `firmware` is the same or almost the same with minor changes, so why not create a "simple-extensible-multi-machine-control-funcy-app"? as someone says `one app to rule them all`.

## Detailed design

<img src="./docs/img/simple-diagram.svg">

## Building

To simplify the module compilation there is a `Makefile` with a variety of tools
ie, for build in local 

```shell
make build
```

this will compile  

```shell
cnix --help
```

if you are a docker user that don't have [golang](https://golang.org/) installed you can use the following command.

```shell
make dist
```

This will create all targets on a `dist` directory. Just find the binary that fits on your machine. 

> The following command will help you a bit.

```shell
ls dist/* |grep --color $(uname -a|awk -v a=$(uname -m) '{printf("%s-%s\n",tolower($1),a=="x86_64"?"amd64":a)}')
```

## Contribution

Checkout this [document](./CONTRIBUTION.md) with all the details and steps to become a contributor!

## Packaging
TODO
At the moment the packaging methods allowed are:
- [debian](./docs/debian)
- [macos](./docs/debian)

## Docker

TODO