# Moderari

![build](https://github.com/fjah/Moderari/workflows/Go/badge.svg) ![discord](https://img.shields.io/discord/714026885390008401) ![go version](https://img.shields.io/github/go-mod/go-version/fjah/Moderari) ![commit activity](https://img.shields.io/github/commit-activity/m/fjah/Moderari)

Moderari is a Discord bot focused that brings

- really easy self-hosting;
- speed;
- scalability; and
- powerfulness

together. Its core is *really* lightweight, using about 25MB RAM despite running both a bot and a HTTP server.

## Self-hosting

Moderari is licensed under v2 of the Apache License, so self-hosting is permitted. You'll just need two things:

- The built Moderari executable; and
- a Redis server.

These can be obtained from [the releases page](https://github.com/fjah/Moderari/releases) and [the Redis website](https://redis.io) respectively.

You can also deploy Moderari with [Docker](https://docker.io); see the `build/package` directory. Simply clone this repository (get `git` if you haven't already), install [Docker](https://docker.io) and [Redis](https://redis.io), then run `./build/package/update.sh`. This is the recommended way of deploying Moderari in *production*.
