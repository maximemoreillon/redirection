# Redirection service

A simple service that redirects requests of specific routes to defined target URLs.
The redirection rules can be set using environment variables or using the a config.yml file in the config fiolder.
Configuring via environment variables handles the `/` route while using `config.yml` allows more granular control over the path.

## Environment variables

```
TARGET=https://yahoo.co.jp
WARNING="true"
```

## config.yml

```yml
- path: /google/
  target: https://google.com
- path: /moreillon
  target: https://maximemoreillon.com
```
