# cast.onion

Welcome to the official repository for cast.onion, the private, application-based internet radio network.

# Developing

```bash
sudo mariadb
CREATE DATABASE cast_onion;
CREATE USER 'cast'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON cast_onion.* TO 'cast'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```
to start

```bash
go mod tidy

go run ./cmd/server
```

## Issues

If you find any issues within the client or service, please make a issue in [Issues](https://github.com/ItsBr0dyy/cast.onion/issues/new).

Include info like:

- The issue (SQL injection, etc.)
- Description of the issue
- Where the issue takes place (Client or service(API, GraphQL, CDN, etc.))
- A image example of the issue
