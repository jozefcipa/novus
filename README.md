<p align="center">
  <img src="./assets/banner.png">
</p>

Briefly introduce your CLI tool and its purpose. Explain how it helps developers deploy websites on localhost more efficiently.

## Overview
- show an animation gif (defining routes in config file, then switch to console, and run novus serve)
- what it does, how it helps you

- combines `mkcert`, `nginx` and `DNSMasq` to provide a simple way to work with your local web applications using regular HTTPS URLs instead of several `localhost` addresses with different ports.

## Install
Installing Novus is as simple as running 
```bash
$ brew tap jozefcipa/novus
$ brew install novus
```

## How to use
- `novus init`
- then define your routes in the config
```yaml
appName: my-app
routes:
  - domain: my-frontend.test
    upstream: http://localhost:3000
  - domain: my-api.test
    upstream: http://localhost:5000
```

- run `novus serve`

It will ask for your password as it performs some `sudo` calls (for managing DNS resolvers)


## Commands
Explain how to use your CLI tool, including command syntax and available options. Provide examples of common use cases to help users get started quickly.

- novus init
- novus serve
- novus status
- novus stop
- novus pause [app]
- novus resume [app]
- novus remove [app]

## Notes
Do not use top level domains (TLD) defined by [IANA](https://www.iana.org/domains/root/db)
This will result in redirecting all URLs using the given TLD to localhost
e.g. my.local.website.com -> all .com websites will stop working

Instead, prefer `.test` or anything else that works for you

`.local` doesn't work (on MacOS) - [Apple article](https://support.apple.com/en-us/101471)
recommend using `.test`

`.dev` - do not use either, this is now a valid TLD domain

## License
Novus is released under the MIT license ([LICENSE](./LICENSE)).
