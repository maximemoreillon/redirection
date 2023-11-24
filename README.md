# Redirection service

A simple service that redirects requests of specific routes to defined target URLs.
The redirection rules are set in the config.yml with the following format:

```yml
# Login
- path: /login
  target: http://10.115.1.100:30097/v2/auth/login
- path: /v2/login
  target: http://10.115.1.100:30097/v2/auth/login
- path: /v3/login
  target: http://10.115.1.100:30097/v3/auth/login

# Whoami
- path: /whoami
  target: http://10.115.1.100:30097/v1/users/self
- path: /v2/whoami
  target: http://10.115.1.100:30097/v2/users/self
- path: /v3/whoami
  target: http://10.115.1.100:30097/v3/users/self
```
