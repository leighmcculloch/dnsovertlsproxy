# dnsovertlsproxy

A simple [DNS over TLS][1] proxy.

Run this application on your computer or server. It will accept standard DNS queries and forward them to a DNS server supporting [DNS over TLS][1].

## Install

### From Source

```
go get 4d63.com/dnsovertlsproxy
```

## Usage

```bash
dnsovertlsproxy -listen :53 -server 1.1.1.1:853
```

[1]: https://en.wikipedia.org/wiki/DNS_over_TLS
