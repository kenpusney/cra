# CRA

CRA 是用来组合后端API请求的工具。

## 用法

```
> go get github.com/kenpusney/cra
> cra https://the-api-you-need-to-visit.com/api/v1
Started at port: 9511 Proxying: https://the-api-you-need-to-visit.com/api/v1
```

直接把请求按照特定格式（参考下文）发给 `http://localhost:9511` 就可以一次发送和接收多个请求了。

## 解决的问题

- 使用RESTful架构设计的后端接口偏向特定功能资源，无法做到灵活应对前端请求
- 针对大量数据和请求的场景，前后端通信频繁造成网络负担过重


## 请求示例

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

这里我们定义了一套层叠的请求，后者依赖前者返回的数据。CRA会依次执行这两个请求，中间遇到任意错误都会中断同时将请求返回。

## 请求模式

- 顺序执行 (`seq`): 按照请求给定的顺序依次执行请求
- 并发执行 (`con`): 使用goroutine批量运行请求，返回结果可能是乱序的
- 层叠执行 (`cascaded`): 按照顺序执行给定的请求，并且可以保存返回结果以便后续请求参数化调用
- **[TODO]** 批量执行 (`batch`): 使用请求中提供的数据或者请求获取到的数据进行批量操作
  
### 层叠执行模式（cascaded）

层叠模式中，你可以通过jsonpath把一部分返回的数据保存在上下文中，通过Mustache模版进行接下来请求的参数化处理。

例如：
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

### 批量执行模式（batch）

批量模式中，你需要指定一组数据作为批量执行的种子。

比如下面这个例子中，删除所有过期的资源：
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

