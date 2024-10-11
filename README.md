## ここまでやったことメモ

Go言語でServerとClient書くまで

[Connect for Go - Getting started](https://connectrpc.com/docs/go/getting-started)

SvelteKit周り

[Connect for Web - Getting started](https://connectrpc.com/docs/web/getting-started)

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

#### 一覧

通信できるまでにやったことまとめ

- ./cmd/server/main.goのCORS対応の追加(フロントエンド側から操作できるようにするため)
- buf.gen.yamlにtypescriptファイル用の設定追加
- ルートに必要なパッケージをnpm install --save-devで追加
- SvelteKitのプロジェクト追加
- npx buf generateで対応ファイル作成
- SvelteKitのコンポーネント追加

#### やったこと

この辺を参考にしながら

https://connectrpc.com/docs/web/getting-started

今の構成にsvelte-todoを追加したので
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

## MySQLのDBに実際に書き込むようにする

### 前準備としてWindowsに入ってるMySQLをWSL2側から繋ぐ

Windows側にMySQL5.6が入っているためそちらにつなぐ

Windows側とUbuntu側のWSL用のIPアドレスを調べておく

MySQL側で、あらかじめWSL2側のIPを許可しておく

```mysql
GRANT ALL PRIVILEGES ON *.* TO 'root'@'WSL2_IPADDRESS' IDENTIFIED BY 'password'
```

```
sudo apt install mariadb-client

mysql -h WINDOWS_IPADDRESS -u root -p
```

で接続確認

### sql-migrateを入れてマイグレーションする

[sql-migrate](https://github.com/rubenv/sql-migrate)

```cmd
go install github.com/rubenv/sql-migrate/...@latest
```

dbconfig.yaml

https://github.com/ykobayashi-techno/connect-go-todo/blob/main/dbconfig.yml

DB用の設定。

DB作成自体は手動で行う必要があるみたいなのであらかじめDBを作成する

```
CREATE DATABASE todos_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

.envファイルに設定値をあらかじめ書いておいて sql-migrate up を実行

```
export $(grep -v '^#' .env | xargs) && sql-migrate up
```

```mysql
mysql> use todos_dev;
Database changed
mysql> show tables;
+---------------------+
| Tables_in_todos_dev |
+---------------------+
| gorp_migrations     |
| todos               |
+---------------------+
```

migrationができた

sql-migrate down実行のたびに1STEPずつdownされていく

### sqlboilerでDB構造をモデルにしてGoで読み書きできるようにする

todoアプリ、サーバー側の処理にMySQLの書き込み等を入れる

[sqlboiler](https://github.com/volatiletech/sqlboiler)

インストール

```
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
```

sqlboiler.tomlを作成して

https://github.com/ykobayashi-techno/connect-go-todo/blob/main/sqlboiler.toml

DBの設定を行う

`blacklist=["gorp_migrations"]`で設定を行って、sql-migrateが管理しているテーブルを対象外にする

```
sqlboiler mysql
```

でDBの構造に対応したデータが`model`フォルダ内に出来上がるのでこれでGoのサーバーでMySQLの読み書きができる

### 実装

https://github.com/ykobayashi-techno/connect-go-todo/commit/8cdf7bdf898623593deb38ec2a65939b4f2ea165

MySQLに書き込み、読み込みするように実装を変更