# bitnel/api

Here we go again... Let's get rly rich. ;)

## Database migrations

Some package named [goose](https://bitbucket.org/liamstask/goose) is used for database migrations.

Make sure to copy `db/dbconf.yml.example` to `db/dbconf.yml` and edit as nessecary.

## Configuration

At the moment, `config.json` is read in the current directory.

## Error handling

A typical request handler will implement `apiHandler`.

`type apiHandler func(http.ResponseWriter, *http.Request) *serverError`

It accepts the parameters of an `http.HandlerFunc` with an added return value of an `*serverError`. Server errors are errors that we encounter that are not caused by the user. In these cases, we return a `*serverError`.

`return &serverError{err, "cannot start db tx"}`

A 500 internal server error response will then be given.

For other errors, such as request validation errors, we use `writeError` to return an appropriate reponse to the client.