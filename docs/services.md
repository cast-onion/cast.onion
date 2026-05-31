# Services

We have a few different services built into the backend and client:

- `api`
- `cdn`
- `graphql`
- `website`
- `websocket`
- `keygen` - for randomized station keys, access tokens, ws session keys
- `sql` - database
- `email-sender` - for confirmation
- `cloudflare` - dns

## API

The API server is used to run the v1 API and following components like:

- `cdn`
- `graphql`
- `websocket`

The following 3 components are built into the API and run together.

## Website

The website is built using Typescript and Svelte for simplicity.

- `typscript` - connection to websockets and backend
- `svelte` - simple & sustainability web developement

## Email Sender

This function sends emails to the station owner / applicant whether their station application got approved or denied, or go revoked or suspended.

## Keygen

The keygen component is used to generate all:

- `session keys` - for the websocket and other API components
- `access tokens` - for certain permissions
- `station keys` - specified key for a station

## SQL

SQL is used to run the whole database, more info at docs/database.md