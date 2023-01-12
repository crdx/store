# store

**store** is the client-side component of a RESTful stringly-typed key/value store for Go (server-side component not included).

It supports several primitive operations that can be used with any server that implements the API defined in the [API](#api) section.

## Installation

```sh
go get github.com/crdx/store
```

## Example

Error handling has been omitted for the sake of brevity.

```go
import (
    "fmt"

    "github.com/crdx/store"
)

func main() {
    store := store.New("https://.../api/store", "5d41402abc4b2a76b9719d911017c592")

    status, _ := store.Set("foo", "bar")
    status, _ := store.Append("foo", "baz")
    status, _ := store.Delete("foo")

    value, _  := store.Get("foo")
    value, _  := store.GetOrDefault("foo", "(not set)")

    list, _ := store.List()
    for _, k := range list {
        fmt.Println(k)
    }
}
```

## Methods

### `New(baseUrl, apiToken)`

Instantiates a new `store.Store` using `baseUrl` as the base URL and `apiToken` as the API token.

Returns `*store.Store`.

### `Set(key, value)`

Sets the value of a key.

Returns a human-readable string stating which action the server carried out.

### `Append(key, value)`

Appends a line to the current value of a key.

Returns a human-readable string stating which action the server carried out.

### `Get(key)`

Returns the value of a key.

### `GetOrDefault(key, defaultValue)`

Returns the value of a key.

If the value is empty then `defaultValue` is returned.

### `Delete(key)`

Deletes a key.

Returns a human-readable string stating which action the server carried out.

### `List()`

Returns a slice containing each of the keys currently stored.

## API

This section describes the API implementation.

Authentication is done with the standard Bearer token pattern: a header of the form `Authorization: Bearer $TOKEN` is sent along with every request.

### `GET /`

Returns a list of all keys stored.

```sh
curl $BASE_URL --json --oauth2-bearer $API_TOKEN
```

The JSON payload for a successful response must be of the following form.

```json
{
    "success": true,
    "items": [
        { "k": "foo" },
        { "k": "bar" },
    ]
}
```

If a server-side error occurs then `success` must be set to `false` and a `message` key must be provided with the server's human-readable assessment of the error.

### `GET /:k`

Returns the value of a key.

```sh
curl $BASE_URL/foo --json --oauth2-bearer $API_TOKEN
```

The JSON payload for a successful response must be of the following form.

```json
{
    "success": true,
    "value": "bar"
}
```

If a server-side error occurs then `success` must be set to `false` and a `message` key must be provided with the server's human-readable assessment of the error.

### `POST /:k`

Sets the value of a key.

```sh
curl $BASE_URL/foo -d '{ "value": "bar" }' --json --oauth2-bearer $API_TOKEN
```

The JSON payload for a successful response must be of the following form.

```json
{
    "success": true,
    "message": "added"
}
```

If a server-side error occurs then `success` must be set to `false` and the `message` key must be set to the server's human-readable assessment of the error.

### `DELETE /:k`

Deletes a key.

```sh
curl -X DELETE $BASE_URL/foo --json --oauth2-bearer $API_TOKEN
```

The JSON payload for a successful response must be of the following form.

```json
{
    "success": true,
    "message": "deleted"
}
```

If a server-side error occurs then `success` must be set to `false` and the `message` key must be set to the server's human-readable assessment of the error.

## Contributions

Open an [issue](https://github.com/crdx/store/issues) or send a [pull request](https://github.com/crdx/store/pulls).

## Licence

[GPLv3](LICENCE).
