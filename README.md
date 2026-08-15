# Craft Course Material Recommender

This Go 1.22.12 backend recommends fabric, toolkits, and semi-finished materials from a learner's enrolled courses, purchase history, preferred difficulty, material-type preferences, and ratings. It uses an in-memory catalog and deterministic fixture data, and reports recommendation reasons plus a confusion matrix, precision, recall, and accuracy.

## Run

Use the built-in fixture:

```sh
CGO_ENABLED=0 go run ./cmd/craftmaterials
```

Read a local CSV:

```sh
CGO_ENABLED=0 go run ./cmd/craftmaterials -input fixtures/learners.csv -limit 3
```

Start the backend upload endpoint:

```sh
CGO_ENABLED=0 go run ./cmd/craftmaterials -listen :8080
curl -F 'csv=@fixtures/learners.csv' 'http://localhost:8080/recommend?limit=3'
```

`POST /recommend` accepts either a `text/csv` request body or a multipart upload in the `csv` field. `GET /health` returns service status.

## CSV schema

The required columns are:

```text
learner_id,enrolled_courses,purchased_materials,preferred_difficulty,preferred_kinds,ratings
```

List fields use semicolons. Ratings use `kind=value` pairs. Supported difficulties are `beginner`, `intermediate`, and `advanced`; supported kinds are `fabric`, `toolkit`, and `semi_finished`.

## Verify

```sh
CGO_ENABLED=0 go test -count=1 ./...
```

The lifecycle acceptance case exercises archived business documents as part of the end-to-end regression suite.

