# dnsovertlsproxy

A simple [DNS over TLS][1] proxy.

Run this application on your computer or server. It will accept standard DNS queries and forward them to a DNS server supporting [DNS over TLS][1].

## Install

### Binary (Linux; macOS; Windows)

Download and install the binary from the [releases](https://github.com/leighmcculloch/dnsovertlsproxy/releases) page.

### Brew (macOS)

```
brew install 4d63/dnsovertlsproxy/dnsovertlsproxy
```

### From Source

```
go get 4d63.com/dnsovertlsproxy
```

## Usage

### General

```
dnsovertlsproxy -listen :53 -server 1.1.1.1:853
```

### Brew (macOS)

Start it as a service and tell brew to configure it to run on boot:

```
sudo brew services start dnsovertlsproxy
```

[1]: https://en.wikipedia.org/wiki/DNS_over_TLS
