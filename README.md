<p align="center">
  <img src="./assets/banner.png">
</p>

Novus is a CLI tool that enables local sites to use SSL and actual domain names, instead of `localhost`.

[gif] (defining routes in config file, then switch to console, and run novus serve)

## Overview

Novus streamlines managing of numerous `localhost` services by providing a simple way to define regular domain names instead. It comes with built-in HTTPS support so all domains are secure by default.

In the background it’s just good old **Nginx** acting as a proxy and **DNSMasq** for defining custom domain resolvers. No more `/etc/hosts` manipulation. SSL certificates are automatically managed and renewed for you by **mkcert**.

All you have to do is **map your [localhost](http://localhost) URLs to the DNS domains**. The rest is up to Novus and you can enjoy a seamless production-like experience on your machine 💯.

## Installing

Installing Novus is very simple and can be done in two steps.

```bash
$ brew tap jozefcipa/novus
$ brew install novus
```

You can verify Novus has been install by running

```bash
$ novus -v
```

## Usage

To start using Novus, run `novus init`.

It creates a `novus.yml` configuration file that you can open in your editor and define your domains mapping.

**Example configuration:**

```yaml
appName: my-app
routes:
  - domain: my-frontend.test
    upstream: http://localhost:3000
  - domain: my-api.test
    upstream: http://localhost:4000
```

Once you’re done, just call `novus serve` and you can start using nice HTTPs domains locally.

**Note:** It will ask for your password as it performs some `sudo` calls (for managing DNS resolvers).

## Commands

Here is the list of all available commands.<br/>
You can run them by calling `novus [command]`

| Command | Description |
| ------- | ----------- |
| init | Initializes the Novus proxy. Install the necessary binaries and creates a configuration file (novus.yml) |
| serve | Reads the configuration file, updates DNS, creates SSL certificates and registers routes. |
| status | Shows Novus status and all registered apps. |
| stop | Stops Novus routing. |
| pause&nbsp;[app] | Pauses routing of a specific app. (Needed if there are multiple apps defined with the conflicting domains) |
| resume&nbsp;[app] | Starts routing the paused app again. |
| remove&nbsp;[app] | Removes app configuration from Novus and stops routing. |

## Notes

💡 **Prefer** `.test` or any other postfix that is not a TLD domain

❌ **Do not use** top level domains (TLD) defined by [IANA](https://www.iana.org/domains/root/db) <br/>
👉 This will result in redirecting all URLs using the given TLD to localhost
    e.g. `my.local.website.com` -> all `*.com` websites will stop working (**!**)

❌ MacOS doesn't work well with `.local` TLDs. See [Apple article](https://support.apple.com/en-us/101471) for more.

❌  **Do not use** `.dev` domain either, this is now a valid TLD domain

## **License**

Novus is released under the MIT license. See [LICENSE](./LICENSE).