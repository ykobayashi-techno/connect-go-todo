## メモ

[go connect - getting started](https://connectrpc.com/docs/go/getting-started/)

### apiを生やす

なんか増やしたければ`todo/v1/todo.proto`を編集してデータや処理を追加する

そのあと

```cmd
buf lint

buf generate
```

### サーバー実行

サーバー起動

```cmd
cmd/server/main.go
```

アクセス確認

```
grpcurl -protoset <(buf build -o -) -plaintext -d '{}' localhost:8080 todo.v1.TodoService.GetAllTasks
```

とかやると結果が返ってくる
