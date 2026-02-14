# Redirection service

A simple service written in Go used for redirecting HTTP requests.

It can be configured using either environment variables or a configuration file.

## Using environment variables

To redirect users directly with a 307 code:

```
docker run \
  -e REDIRECTION_TARGET_URL=http://new.example.com \
  moreillon/redirection
```

To show users a warning that the content has been moved:

```
docker run \
  -e REDIRECTION_TARGET_URL=http://new.example.com \
  -e REDIRECTION_WARNING=true \
  moreillon/redirection

```

## Using configuration file

When the REDIRECTION_TARGET_URL environment variable is not set, configuration can be done using a file names `config.yml` such as the following:

```
# Will redirect to https://example.com/abc
- path: /abc
  target: https://example.com
  warn: true
```

In this case, the command becomes:

```
docker run \
  -v "$(pwd)"/config.yml:/app/config/config.yml \
  moreillon/redirection
```

## development
