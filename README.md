## ここまでやったことメモ

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

### SvelteKitアプリ準備

この辺を参考にしながら
https://connectrpc.com/docs/node/implementing-services

今の構成にtodo-svelte-appを追加したので
buf.gen.yamlを編集

プロジェクトのルートで

```
npm install --save-dev @bufbuild/buf @connectrpc/protoc-gen-connect-es@"^1.0.0" @bufbuild/protoc-gen-es@"^1.0.0"
```

後に

```
npm create svelte@latest svelte-todo
```

```
npx buf generate --path ./todo/v1/todo.proto
```

行うと、svelte-todo/src/genにファイルが出力される

svelte-todo内で

```
npm install @connectrpc/connect@"^1.0.0" @connectrpc/connect-web@"^1.0.0" @bufbuild/protobuf@"^1.0.0"
```

### サーバー立ち上げ

./cmd/server/

```
go run main.go
```

### アプリ立ち上げ

SvelteKit側で各種処理記述

./svelte-todo/
```
npm run dev
```
