# Resume Assets (`resume-assets`)

A small Go HTTP service for serving authorized static assets. It looks up asset metadata in MongoDB, checks the access policy, then streams the object from an S3-compatible storage service using `github.com/nawaphonOHM/whatever/v3`.

## HTTP routes

| Method | Path | Description |
|---|---|---|
| `GET` | `/assets/{filename}` | Streams one asset whose MongoDB `access` value is exactly `ALLOW`. `filename` is one decoded path segment. |
| `GET` | `/health` | Framework health probe. |
| `GET` | `/ready` | Framework readiness probe. |

The only application route is the asset `GET`. Nested asset paths, other HTTP methods, SPA routes, index fallback, and asset modification APIs are not provided. Empty, dot, traversal, and slash-containing filename values are rejected.

## MongoDB metadata

The server reads the fixed MongoDB collection named `assets`. Documents use this schema:

```json
{
  "filename": "app.js",
  "location": "s3://assets-bucket/path/app.js",
  "access": "ALLOW"
}
```

- `filename` is the exact lookup key and must be unique. Create a MongoDB unique index on `filename` before using the service; the server does not modify the database or create indexes automatically.
- `location` identifies an object as described below.
- `access` must be exactly `ALLOW` for storage access. `DENY`, a missing field, lowercase values, or any other value (including a non-string value) is rejected with HTTP `403`.

Missing metadata and missing storage objects return `404`. Invalid metadata and backend failures return a generic `500` response; internal error details are logged, not returned to clients.

## Supported object locations

All locations are fetched through the configured S3-compatible storage client:

| `location` | Bucket | Object key |
|---|---|---|
| `s3://bucket-name/path/app.js` | `bucket-name` | `path/app.js` |
| `gs://bucket-name/path/app.js` | `bucket-name` | `path/app.js` |
| `s3:///path/app.js` | `OHM9996_STORAGE_BUCKET` | `path/app.js` |
| `gs:///path/app.js` | `OHM9996_STORAGE_BUCKET` | `path/app.js` |
| `path/app.js` or `/path/app.js` | `OHM9996_STORAGE_BUCKET` | `path/app.js` |

Locations with unsupported schemes or without an object key are invalid. Leading slashes on plain keys are removed; nested object key names within `location` are supported even though nested HTTP routes are not.

Successful responses stream the object and forward available `Content-Type`, `Content-Length`, `ETag`, and HTTP-formatted `Last-Modified` metadata.

## Configuration

An optional `.env` file (or `configs/.env`) is loaded at startup. Existing process environment variables take precedence. The shared library reads the `OHM9996_*` MongoDB, storage, and server settings; commonly used values include:

| Variable | Purpose |
|---|---|
| `OHM9996_SERVER_HOST` | HTTP bind address. |
| `OHM9996_SERVER_PORT` | HTTP port; the shared framework default is `8080`. |
| `OHM9996_SERVER_API_GROUP` | Must be `raw` so the asset route is exactly `/assets/{filename}`; defaults to `raw`. |
| `OHM9996_MONGODB_HOST` | MongoDB server address. |
| `OHM9996_MONGODB_PORT` | MongoDB server port; the shared library default is `27017`. |
| `OHM9996_MONGODB_DATABASE` | Database containing the `assets` collection. |
| `OHM9996_MONGODB_USERNAME` / `OHM9996_MONGODB_PASSWORD` | Optional MongoDB credentials. |
| `OHM9996_STORAGE_BUCKET` | Default bucket used for plain keys and bucket-less S3/GS URIs. |
| `OHM9996_STORAGE_REGION` | Storage region; the shared library default is `us-east-1`. |
| `OHM9996_STORAGE_ENDPOINT` | Optional custom S3-compatible endpoint (for example, MinIO). |
| `OHM9996_STORAGE_ACCESS_KEY_ID` / `OHM9996_STORAGE_SECRET_ACCESS_KEY` | Optional static storage credentials; no secrets are embedded in the application. |
| `OHM9996_STORAGE_USE_PATH_STYLE` | Enable path-style addressing for compatible storage services when needed. |

The shared clients also support their standard `OHM9996_*` timeout, pool, and authentication settings. Consult the `whatever/v3` configuration documentation for the complete list.

## Build, run, and test

```bash
go mod verify
go test -race ./...
go vet ./...
go build ./...
```

Existing Make targets are preserved:

```bash
make run
make test
make vet
make build
make lint
make tidy
```

Tests use injected fakes and the framework router; no live MongoDB or S3-compatible service is required.