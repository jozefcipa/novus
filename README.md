# novus

<p align="center">
  <img src="./assets/gopher.png" width="200">
</p>

Briefly introduce your CLI tool and its purpose. Explain how it helps developers deploy websites on localhost more efficiently.

## Overview
- what it does, how it helps
- what it uses under the hood (nginx, dnsmasq, mkcert)

## Installing
- describe steps to install the binary
- through brew

## Commands
Explain how to use your CLI tool, including command syntax and available options. Provide examples of common use cases to help users get started quickly.

- describe all commands

- novus serve
- novus serve --create-config
- novus status
- novus stop
- novus trust

## Configuration
Describe any configuration options or settings that users can customize to fit their development environment. Provide examples and explanations for each option.

- describe `novus.yml` config file

```yaml
routes:
  - domain: my-frontend.test
    upstream: http://localhost:3000
  - domain: my-api.test
    upstream: http://localhost:5050
  - domain: my-search-api.test
    upstream: http://localhost:8000
```


## License
Novus is released under the MIT license. See [LICENSE](./LICENSE)
