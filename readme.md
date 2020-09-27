SSO
===

Automatic multi-site single sign on a la Stack Exchange

This library is a proof of concept not yet suitable for production.

Check out the example directory to see how to use it.

## Configuration

Configuration is currently done by environment variables.

| Variable | Default | Description |
|----------|---------|-------------|
| CLIENT_DOMAINS | | Comma-separated list of domains which are allowed to use this authentication sever |
| SECURE_ONLY | true | Only include cookies over secure channels (https) |
| BCRYPT_COST | 10 | Bcrypt cost cycles. This primarily effects how long password hashing takes. Less than 10 becomes too fast (< 50 ms). More than 15 becomes very slow (> 1 s). |
| AUTH_DOMAIN | localhost | The domain of this server. This appears on cookies and tokens. |
| DATABASE_URL | ./auth.db | The path of the SQLite database |


## How it works

First, the client app website will run the client script, which creates an
iframe for the SSO website. If the SSO website recognizes the referrer as a
client domain, it will initialize communication with parent via `postMessage()`.

If the browser has a session via a cookie with the auth domain, regardless of
which client it was started with, the SSO server will provide a token to the
sso website, which will be passed to the app website.

The app server should then verify the token matches a pre-shared signing key in
order to trust the token.

The app website can issue requests to the auth server via passed messages in 
order to register an account, login, and logout.



        Client                                       Server
       ========                                     ========

    ┌─────────────┐                             ┌─────────────┐
    │ 4.          │                             │ 3.          │
    │   App       │    Pre-shared signing key   │   SSO Auth  │
    │   Server    │<────────────────────────────┤   Server    │
    │             │                             │             │
    └─────────────┘                             └───┬─────────┘
          ^                                         |     ^
          │                                     Set │     |
          │                                  Cookie │     |
          │ JWT                                     │     | Forward
          │ Auth                              Issue │     | Messages
          │                                     JWT │     |
          │                                         │     |
          │                                         v     |
    ┌─────┴───────┐                             ┌─────────┴───┐
    │ 1.          │       ┌─────────────┐       │ 2.          │
    │   App       │       │             │       │   SSO       │
    │   Website   │<─────>│   iframe    │<─────>│   Website   │
    │             │       └─────────────┘       │             │
    └─────────────┘      Passing messages       └─────────────┘
                         via postMessage()


