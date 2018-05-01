# jhop [![codecov](https://codecov.io/gh/cooldrip/jhop/branch/master/graph/badge.svg)](https://codecov.io/gh/cooldrip/jhop) [![Build Status](https://travis-ci.org/cooldrip/jhop.svg?branch=master)](https://travis-ci.org/cooldrip/jhop)

🏎Create fake REST API in one sec.

## Install

```
~ $ go get github.com/cooldrip/jhop/cmd/jhop
```

## CLI Example

Create a file `recipes.json`:

```json
{
  "recipes": [
    { "id": 1, "prep_time": "1h", "difficulty": "hard" },
    { "id": 2, "prep_time": "15m", "difficulty": "easy" }
  ]
}
```

Passing the JSON file to `jhop`:

```
~ $ jhop recipes.json
```

Now you can go to `localhost:6000/recipes` and get the collection:

```json
{
  "recipes": [
    { "id": 1, "prep_time": "1h", "difficulty": "hard" },
    { "id": 2, "prep_time": "15m", "difficulty": "easy" }
  ]
}
```

or you can just get a single recipe `localhost:6000/recipes/1`:

```json
{ "id": 1, "prep_time": "1h", "difficulty": "hard" }
```

## API example

```go
  f, err := os.Open("recipes.json")
  if err != nil {
    log.Fatalf("file opening failed: %s", err)
  }
  h, err := jhop.NewHandler(f)
  if err != nil {
    log.Fatalf("handler initialization failed: %s", err)
  }

  s := httptest.NewServer(h)
  defer s.Close()

  resp, err := http.Get(fmt.Sprintf("%s/recipes/1", s.URL))
  if err != nil {
    log.Fatalf("request to test server failed: %s", err)
  }

  io.Copy(os.Stdout, resp.Body) // {"difficulty":"hard","id":1,"prep_time":"1h"}
```

## CLI usage

```
NAME:
   jhop - Create fake REST API in one sec.

USAGE:
   jhop [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value   Set port (default: "6000")
   --host value   Set host (default: "localhost")
   --help, -h     show help
   --version, -v  print the version
```

## TODO

* [x] Single resource
* [ ] Custom routes
* [ ] Middleware
