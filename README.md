# stream games

## About

Twitch-only(for now) chat game based on viewer's messages ![alt text](https://github.com/toxicOctopus/sg/blob/master/kappa.png?raw=true).

## Required configs

- `config/env/values.json`

example config @ `test_fixtures/values.json`


- `config/centrifugo.json`

https://centrifugal.github.io/centrifugo/server/configuration/

simplest usable config(with ur creds):

```
{
  "v3_use_offset": true,
  "token_hmac_secret_key": "my_secret",
  "api_key": "my_api_key",
  "admin_password": "password",
  "admin_secret": "secret",
  "admin": true,
  "publish": true,
  "presence": true,
  "namespaces": [
    {
      "name": "public",
      "publish": true,
      "history_size": 10,
      "history_lifetime": 300,
      "history_recover": true
    }
  ]
}
```