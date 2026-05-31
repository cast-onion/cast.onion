# Repository Breakdown

Below will explain how the whole client and service is designed and what each component is meant for.

## Repository File Tree

```
├── app/ - client
├── cdn-files/ - cdn files that can be accessed through the cdn on a browser
├── cmd/
│   └── server/ - the "runner" of the server
├── docs/ - full documentation on the repository
├── email/ - rust email sender
├── internal/
│   ├── api/
│   │   ├── graphql/ - graphql host
│   │   ├── rest/ - v1 rest api host
│   │   ├── web/ - web api host
│   │   └── websocket/ - websocket host
│   ├── cache/ - redis cache
│   ├── cdn/ - hosts the cdn service
│   ├── config/ - configuration
│   ├── db/
│   │   ├── migrations/ - updates sql
│   │   └── queries/ - selects which query will be ran
│   ├── model/ - models for stations and applications
│   ├── room/ - guest invites
│   ├── service/ - handles all requests
│   └── stream/ - service that handles the broadcast
├── lang/
│    └── en.json - official English translation
├── pkg/
│    ├── cdn/ - lead to cdn site
│    └── keygen/ - generate keys and tokens
├── scripts/ - test scripts
├── terms/ - includes all of the terms we use per user
├── web/
│   └── src/ - svelte
│       ├── components/ - svelte components
│       └── lib/ - nav bars
│       └── routes/ - routes throughout the site
│       └── types/ - index.ts
├── .dockerignore - hides Cargo.lock and target files/folders from Docker
├── .env - setting certain variables
├── .gitignore - hides Cargo.lock and target files/folders from Git commit
├── .prettierrc - Prettier formatting
├── config.json - setting the login and password for admin dashboard
├── Dockerfile - file that Docker will run as
├── go.mod
├── go.sum
├── migrations.go
└── README.md
```

## Services

The official cast.onion repository runs all of the services within 3 different ports:

- `api`
- `websocket`
- `cdn`
- `graphql`
- `website`
- `keygen`
- `email-sender`

## Languages

This repository is a full stack project including:

- `Go` - Backend
- `Rust` - Client
- `Typscript` - Website
- `Svelte` - Website

## Documentation

You can find full documentation at /docs.