# tlsthing

Keeps yer internal TLS certs in order.

## director

Talks to the database which keeps the certificate registrations. Instructs clients to go get certificate updates.

### env vars

TLSTHING_PORT='2443'  Listening port for clients.

DATABASE_PORT='5432'  Connection port for database. Postgres only at the moment.

DATABASE_HOST='localhost'  Remote host for database.

DATABASE_NAME='tlsthing'  Name for database, because maybe you already have a database with this name.

DATABASE_TLSMODE='disable'  TLS mode, or "SSL" mode as Postgres calls it.

VAULT_ADDR='https://localhost'  Vault address.

VAULT_TOKEN=''  A Vault token with access to the database creds role.

POSTGRES_USERNAME='postgres'  Static Postgres username.

POSTGRES_PASSWORD='temppw'  Static Postgres password.

## client

Talks to the director and fetches certificate updates from a Vault CA.

### env vars

TLSTHING_PORT='2443'  Connection port for director.

TLSTHING_HOST='localhost'  Remote host for director.

VAULT_ADDR='https://localhost'  Vault address.

VAULT_TOKEN=''  A Vault token with access to the database creds role.
