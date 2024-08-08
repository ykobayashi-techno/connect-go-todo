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

### SvelteKitから

この辺を参考にしながら
https://connectrpc.com/docs/node/implementing-services

今の構成にtodo-svelte-appを追加したので
buf.gen.yamlを編集

buf generate用にnpmのライブラリをグローバルに入れておく

```
npm install -g @bufbuild/protoc-gen-es @bufbuild/protoc-gen-connect-es
```

これでbuf generateできる
