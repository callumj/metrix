# Metrix

Metrix is a basic metric logging utility. 

If you're looking for something reliable, look at [statsd](https://github.com/etsy/statsd/), this project exists as a way for me to better understand Go and applications that handle large requests.

Right now it handles the anonymous stats for [Extended](https://itunes.apple.com/us/app/extended/id836630098).

It uses Redis hashes segmented by "DDMMYYYY"

## Configuration

`config.yml` should look like this (everything is optional). Defaults are included

```YAML
sentry: "HTTP address to Sentry DSN"
listen: ":8080"
redis:
  server: ":6379"
  password: ""
```

As noted above Metrix supports [Sentry](https://getsentry.com) for exception logging.