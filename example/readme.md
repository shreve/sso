SSO Example
===========

To run the example, you need to run the SSO server as well as a client server.

The simplest way to do this is with the following commands:

    $ go run example/main.go
    $ cd example/web && python -m http.server

To see how it works, let's check out how it works with different domains.

1. Set the environment variable `CLIENT_DOMAINS=localhost:8000` before running
   the SSO server.
2. Open the client and see that you are able to register and log in with the
   form.
3. Change the port of the client site by specifying it at the end of the python
   command (e.g. `python -m http.server 8001`).
4. See that authentication doesn't work.
5. Restart the auth server specifying the new `CLIENT_DOMAINS` value.
6. See that you are already logged in on the new domain.
