# Bash nats

Bash nats is a layer built in top of go, to allow bash commands manage nats communication.

## Build status

* master:  [![CircleCI Master](https://circleci.com/gh/ErnestIO/bash-nats/tree/master.svg?style=svg&circle-token=627e89c447fe342aff9815ca146b081a37c075ad)](https://circleci.com/gh/r3labs/bash-nats/tree/master)
* develop: [![CircleCI Develop](https://circleci.com/gh/ErnestIO/bash-nats/tree/develop.svg?style=svg&circle-token=627e89c447fe342aff9815ca146b081a37c075ad)](https://circleci.com/gh/r3labs/bash-nats/tree/develop)

## Installation

```
go get git.r3labs.io:libraries/bash-nats
```

## Useing it

Generally bash-nats can be used like:
```
$ bash-nats subject manager arguments
```
Where:
- subject : the nats subject which will be listening at
- manager : the bash command that will be called in order to manage the message
- arguments : extra arguments the manager may need

So, a real example can look like:
```
bash-nats create-instance jruby adapter.rb
```

## Running Tests

```
make test
```

## Custom NATS server

If you have your nats running on a non default host you can define this host on the environment variable NATS_URI


## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/). 

## Copyright and License

Code and documentation copyright since 2015 r3labs.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).
