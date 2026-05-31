# Databases

We use a few different databases and plan to add more so we do not lose any data.

The databases we have that make up the platform:

- `redis`
- `mariadb`

## Redis

Redis is used for rate limits and cache. 

# mariadb

mariadb is our main database that stores almost all data including:

- `session keys`
- `access tokens`
- `station keys`
- `applications`

## Build the Database

To start the maridb database, you will run this command in your terminal:

```bash
sudo mariadb
CREATE DATABASE cast_onion;
CREATE USER 'cast'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON cast_onion.* TO 'cast'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

This command will start the database as how it is prompted in the codebase.