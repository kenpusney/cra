# CRA

Concurrent Request Agent

## Usage

```
> go get github.com/kenpusney/cra
> cra https://the-api-you-need-to-visit.com/api/v1
Started at port: 9511 Proxying: https://the-api-you-need-to-visit.com/api/v1
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
      },
      "cascading": {
        "id": "$.id"
      }
    },
    {
      "endpoint": "/user/{{id}}",
      "method": "get"
    }
  ]
}
```

In the example above, we've defined two cascading requests 
(execution order guaranteed, and subsequent request depends on
its predecessor). CRA will run these 2 requests one by one,
create a user and then retrieves its profile.

## Modes

- Sequential (`seq`): Executing requests in given order.
- Concurrent (`con`): Executing every request in a goroutine, 
  response may out of order.
- Cascaded (`cascaded`): Executing requests in given order, 
  and subsequent request can be parameterized using previous response.
- **[TODO]** Batch (`batch`): Executing requests in batching using data provided
  in request or from previous response.
  
### Cascaded mode

In cascaded mode, you can save response's data to a context using jsonpath, and
retrieve it in subsequent requests using Mustache template.

For example:
```javascript
{
  "mode": "cascaded",
  "requests": [
    {
      "type": "json",
      "endpoint": "/test.json",
      // cascading: save to context
      "cascading": {
        "value": "$.value"
      }
    },
    {
      "type": "json",
      // retrieve the value from context
      "endpoint": "/{{value}}.json"
    }
  ]
}
```

### Batch mode

In batch mode, you need to specify a list of data as a seed to run the batch
request, either in request spec or using response.

For example, if you need to delete all expired resource:
```javascript
{
  "mode": "batch",
  "requests": [
    {
      "type": "json",
      "endpoint": "/resources?status=expired",
      "seed": {
        // must be an array
        "id": "$.resources[:].id"
      }
    },
    {
      "type": "json",
      "method": "delete",
      "endpoint": "/resource/{{id}}",
      "batch": ["id"]
    }
  ]
}
```

## CRA request schema

see `core/types.go`

## Why golang

 - **static linked distribution**: it is important this kind of tools runs
   on same platform without any dependency. golang provides the mechanism.
 - **concurrency-native support**: goroutine are simpler than thread based
   concurrency.
 - **better ecosystem for web applications**

## TODO

- [ ] More options
- [X] Implement cascading
- [ ] Add doc
- [ ] Bypassing headers
- [ ] ID generating strategy

## Q &amp; A

