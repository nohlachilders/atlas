# ATLAS
`atlas` is currently a very testable template for an HTTP server written
in Go, with a basic CRUD interface for a user account. Plans are to
expand it into a demo CRUD application that primarily deals with the storing
and display of GeoJSON with a very light frontend using [leaflet.js](https://leafletjs.com)
and [HTMX](https://htmx.org).

### Features and Technical Details:
- Main server process is executable via code for testability.
- Support for middleware layers for flexible endpoint outcomes, an example authorization
layer with JWTs and refresh tokens has been implemented.
- Configured via environment variables, or via test code.
- Integration testing pipeline uses a PostgreSQL database.
- Database migration is managed via [`goose`](https://github.com/pressly/goose).
- Database queries are written in raw SQL, and converted into type-safe Go via
the excellent [`sqlc`](https://sqlc.dev/). Database schema is automatically detected
via `sqlc`'s built in `goose` integration.


### Running
#### Configuration
The following environment variables are used to configure `atlas`:
- `ATLAS_PORT`: Port to use. Will default to `":8080"` if not supplied.
- `ATLAS_DB_URL`: Connection string to a PostgreSQL database. The testing environment
(`go test ./...` or `./scripts/test.sh`) will default to `"postgresql://localhost:5432/atlas?sslmode=disable"`
- `ATLAS_PLATFORM`: A string used to configure certain development features of `atlas`.
If set to `"dev"`, will enable the `/reset` endpoint to reset the database, and will
enable endpoints to return more verbose errors that would not be advisable in production.
- `ATLAS_SECRET`: A string used to sign and validate JWTs for authorization. 
Reccomended to set to something cryptographically secure.

#### Endpoints
- `GET /healthz`: Check server availability.
- `POST /reset`: If `ATLAS_PLATFORM` is set to `dev`, this will reset the database.
- `POST /users`: Create a user with a JSON body `{"username":"USERNAME","password":"PASSWORD"}`.
- `GET /users`: Authenticated endpoint, supply JWT in an Authorization header with
format `"Bearer TOKEN"`. Returns user data for the logged in user.
- `PUT /users`: Authenticated endpoint. Update user email/password with the same
format as above.
- `DELETE /users`: Authenticated endpoint. Delete logged in user.
- `POST /login`: Supply credentials in the above format to recieve a JSON containing
user info, including a refresh token.
- `POST /refresh`: Supply a request with a refresh token in the above
Authorization/Bearer token header format to receive a JSON containing a JWT.
- `POST /revoke`: Supply a request with a refresh token in the above
Authorization/Bearer token header format to revoke the given refresh token.



