# Hosts translator

This repository contains hosts file to different DNS systems translator.

## Why

While it may sound crazy - why do anyone need to translate hosts file? - it has own specific use case.

I have some internal DNS servers powered by PowerDNS and actually do not want to tinker with hosts file for adding ability to use systems and servers at work. So instead of modifying hosts file on every system I have I just launching hosts translator and it seamlessly (more or less) updating DNS server data with provided hosts file.

## Caveats

### PowerDNS storage provider

1. It assumes that server ID is 'localhost'. It is 'kind-of-default', especially if native database replication is used for zones distribution.
2. It assumes that every domain has only one address.
3. It assumes IPv4 usage. Some additional love is required for proper AAAA support.

## Installation

TBW

## Configuration

TBW
