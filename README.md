# TodoApp - Backend

Go製のTodoアプリバックエンドです。

- Go 1.26
- PostgreSQL 18.4

コンテナで実行するためローカルへのインストールは不要です。

## 必要条件

- Docker / Docker Compose

## セットアップ（ローカル）

1. リポジトリをクローンする

```bash
git clone git@github.com:tamu-b-TodoApp/backend.git
```

2. 環境変数を設定する

```bash
cp .env.example .env
```

`JWT_SECRET` はこの時点では空欄のままで問題ありません（手順4で自動生成します）。

3. コンテナを立ち上げる

```bash
docker compose up -d
```

4. JWT_SECRET を生成する

```bash
docker compose exec backend make gen-jwt-secret
```

`.env` に `JWT_SECRET` が書き込まれます。

5. マイグレーション・シードを実行する

```bash
docker compose exec backend make migrate-apply
docker compose exec backend make seed
```

6. 動作確認

サーバーを起動します。

```bash
docker compose exec backend make dev
```

ヘルスチェックエンドポイントにリクエストして `200 OK` が返れば成功です。

```bash
curl http://localhost:8080/health
```

## テスト

```bash
go test ./...
```

## 開発コマンド

| コマンド                        | 説明                                |
| ------------------------------- | ----------------------------------- |
| `make dev`                      | 開発サーバー起動                    |
| `make migrate-apply`            | マイグレーション実行                |
| `make migrate-diff name=<name>` | マイグレーションファイル生成        |
| `make migrate-reset`            | DB をリセットして再マイグレーション |
| `make seed`                     | シードデータ投入                    |
| `make gen-jwt-secret`           | JWT_SECRET 生成・更新               |

## ライセンス

[MIT](LICENSE)
