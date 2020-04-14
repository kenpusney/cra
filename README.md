# CRA

Concurrent Request Agent

## Usage

```
go get github.com/kenpusney/cra

cra https://the-api-you-need-to-visit.com/api/v1
```

Then you can use CRA DSL to send multiple request in single HTTP transaction.

## DSL

Example:

```json
{
  "mode": "cascaded",
  "requests": [
    {
      "endpoint": "/user",
      "method": "post",
      "type": "json",
      "body": {
        "username": "userxxxxxx",
        "password": "pwdxxxxxxx"
      }
    },
    {
      "endpoint": "/user/{id}",
      "method": "get"
    }
  ]
}
```

In the example above, we've defined two cascading requests 
(execution order guaranteed, and subsequent request depends on
its predecessor). CRA will run these 2 requests one by one,
create a user and then retrieves its profile.

## Why golang

 - **static linked distribution**: it is important this kind of tools runs
   on same platform without any dependency. golang provides the mechanism.
 - **concurrency-native support**: goroutine are simpler than thread based
   concurrency.
 - **better ecosystem for web applications**

## TODO

- [ ] More options
- [ ] Implement cascading
- [ ] Add doc
- [ ] Bypassing headers

## Q &amp; A

