# voodfy-transcocder

Voodfy Transcoder is an open source alternative to the existent cloud encoding services. It is a Queue system that encode videos by tasks.

Join us on our [public Discord channel](https://discord.gg/UjNNkf) for news, discussions, and status updates.

*Warning* This project is still **pre-release** and is not ready for production usage.

## Table of Contents

-   [Prerequisites](#prerequisites)
-   [Design](#design)
-   [Installation](#installation)
-   [Contributing](#contributing)
-   [Changelog](#changelog)
-   [License](#license)

## Prerequisites

To build from source, you need to have Go 1.14 or newer installed.

## Design

Voodfy Transcoder is a queue system used on core of Voodfy to ingest videos using Livepeeer and sendding to IPFS/Filecoin using Powergate

Here's a high-level overview of the main components, and how Powergate interacts with IPFS and a Filecoin client:
![Voodfy Transcoder Design](https://s3.us-west-2.amazonaws.com/secure.notion-static.com/b0fcb3fc-8898-4c69-b38b-f21cd6c3e2f4/voodfy-transcoder.svg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAT73L2G45O3KS52Y5%2F20200712%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20200712T173057Z&X-Amz-Expires=86400&X-Amz-Signature=52f87a17aed16cd34b7c145abe8a0fffd5989360d17d8ae9212966b57b950dcd&X-Amz-SignedHeaders=host&response-content-disposition=filename%20%3D%22voodfy-transcoder.svg%22)

Note in the diagram that the Lotus, Filecoin, IPFS, and Powrgate client node _doesn't need_ to be in the same host where Voodfycli is running. They can, but isn't necessary.

To build and install the CLI, run:
```bash
$ make build-cli
```
The binary will be placed automatically in `$GOPATH/bin` which in general is in `$PATH`, so you can immediately run `voodfycli` in your terminal.

You can run `voodfycli` with the `--help` flag to see the available commands:

```
$ voodfycli --help
A client for transcode video and storage on Filecoin


NAME:
   voodfycli - voodfycli it is the command line interface to add task on voodfy transcoder

USAGE:
   voodfycli [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Leandro Barbosa <contact@voodfy.com>

COMMANDS:
   add, a   add a video to transcode
   ping, p  ping the queue
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

```

## Example

Adding a video inside a folder called /tmp/5f07cca0e64d1e3d92544425/55c6659d3df67bd98c0c2a53f645305d/mountains.mp4

```

$ REDIS_BROKER="localhost:6379" REDIS_RESULT="localhost:6379" voodfycli add local teste mountains.mp4 5f07cca0e64d1e3d92544425 55c6659d3df67bd98c0c2a53f645305d

```

## Installation

Voodfy transcoder installation involves running external dependencies, and wiring them correctly with Filecoin/Lotus IPFS and Powergate.

Please copy the `conf/app.ini` example hosted in this repository

### External dependencies
Powergate needs external dependencies in order to provide full functionality, in particular a synced Filecoin client and a IPFS node and a Redis running.

#### Filecoin client
Currently, we support the Lotus Filecoin client.

Fully syncing a Lotus node can take time, so be sure to check you're fully synced doing `./lotus sync status`.

### IPFS node
A running IPFS node is needed.

Its [Dockerhub repository](https://hub.docker.com/r/ipfs/go-ipfs) if you want to run a contanerized version. Currently we're supporting v0.5.1. The API endpoint should be accessible to Powergate (port 5001, by default).

### Queue system
To build the Powergate server, run:
```bash
$ make up
```

We'll soon provide better information about the integration with Lotus/Filecoin, IPFS and Powergate configurations, stay tuned! 📻

## Contributing

This project is a work in progress. As such, there's a few things you can do right now to help out:

-   **Ask questions**! We'll try to help. Be sure to drop a note (on the above issue) if there is anything you'd like to work on and we'll update the issue to let others know. Also [get in touch](https://discord.gg/UjNNkf) on Discord.
-   **Open issues**, [file issues](https://github.com/Voodfy/voodfy-transcoder/issues), submit pull requests!
-   **Perform code reviews**. More eyes will help a) speed the project along b) ensure quality and c) reduce possible future bugs.
-   **Take a look at the code**. Contributions here that would be most helpful are **top-level comments** about how it should look based on your understanding. Again, the more eyes the better.
-   **Add tests**. There can never be enough tests.

Before you get started, be sure to read our [contributors guide](./CONTRIBUTING.md) and our [contributor covenant code of conduct](./CODE_OF_CONDUCT.md).

## Changelog

[Changelog is published to Releases.](https://github.com/Voodfy/voodfy-transcoder/releases)

## License

[MIT](LICENSE)
