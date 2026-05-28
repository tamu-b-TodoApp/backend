# AGENTS.md — TodoApp Backend 作業ガイド

このファイルはAIエージェント（Claude等）がこのリポジトリで作業する際の指示書です。
**作業開始前に必ずこのファイルを読んでください。**

---

## プロジェクト概要

Go製のTodoAppバックエンドAPI。JWT認証（アクセストークン＋リフレッシュトークン）を備えたRESTful API。

| 項目             | 内容                                            |
| ---------------- | ----------------------------------------------- |
| 言語             | Go 1.26                                         |
| ORM              | GORM                                            |
| DB               | PostgreSQL                                      |
| マイグレーション | Atlas (atlas.hcl)                               |
| 認証             | JWT (golang-jwt/jwt/v5) + bcrypt                |
| HTTPサーバー     | 標準ライブラリ `net/http`（フレームワークなし） |

---

## ディレクトリ構成

```
backend/
├── main.go                  # エントリーポイント。DI（依存注入）はここで行う
├── model/                   # GORMモデル定義
│   ├── todo.go
│   ├── user.go
│   ├── team.go
│   └── refresh_token.go
├── internal/
│   ├── handler/             # HTTPハンドラー（リクエスト/レスポンス処理）
│   ├── service/             # ビジネスロジック
│   ├── repository/          # DB操作
│   ├── middleware/          # 認証ミドルウェア等
│   └── router/              # ルーティング登録
├── migrations/              # Atlasが管理するSQLマイグレーションファイル
├── cmd/
│   ├── atlas/main.go        # Atlas用スキーマ出力コマンド
│   ├── seed/main.go         # シードデータ投入
│   └── gen-jwt-secret/main.go # JWT秘密鍵生成
├── atlas.hcl                # Atlas設定（local / prod 環境）
├── docker-compose.yml       # ローカル開発環境
└── Makefile                 # よく使うコマンド
```

---

## アーキテクチャ

3層構成。各層はインターフェースを介して依存する。

```
Handler → Service → Repository → DB
```

- **Handler**: リクエストのデコード・バリデーション、レスポンスのエンコード
- **Service**: ビジネスロジック。インターフェースで定義
- **Repository**: GORM経由のDB操作。インターフェースで定義
- **DI**: `main.go` で全依存を組み立てて注入する

新しい機能を追加するときは必ずこの3層に分けて実装する。

---

## よく使うコマンド（Makefile）

> **重要**: `go` および `make` コマンドは必ず `backend` コンテナ内で実行すること。
>
> ```bash
> docker compose exec backend <コマンド>
> ```

| コマンド                                  | 内容                                |
| ----------------------------------------- | ----------------------------------- |
| `make dev`                                | 開発サーバー起動 (`go run main.go`) |
| `make build`                              | バイナリビルド → `./server`         |
| `make run`                                | ビルド済みバイナリ実行              |
| `make seed`                               | シードデータ投入                    |
| `make gen-jwt-secret`                     | JWT秘密鍵生成                       |
| `make migrate-apply ENV=local`            | マイグレーション適用（ローカル）    |
| `make migrate-diff name=<name> ENV=local` | マイグレーションファイル生成        |
| `make migrate-reset ENV=local`            | DBリセット＋マイグレーション再適用  |

---

## マイグレーション手順

Atlasを使用。マイグレーションファイルは **手書き禁止**。必ず以下の手順で生成する。

### スキーマ変更の流れ

1. `model/` 配下のGORMモデルを変更する
2. Docker Compose が起動中であることを確認（`atlas-dev-db` コンテナが必要）
3. マイグレーションファイルを生成（コンテナ内で実行）：
   ```bash
   docker compose exec backend make migrate-diff name=<変更内容の名前> ENV=local
   ```
4. `migrations/` に生成されたSQLを確認する
5. マイグレーションを適用（コンテナ内で実行）：
   ```bash
   docker compose exec backend make migrate-apply ENV=local
   ```

> `atlas.sum` は Atlas が自動管理するファイル。手動で編集しない。

---

## コーディング規約

### 全般

- エラーは握りつぶさない。必ず呼び出し元に返すか、適切なHTTPステータスで応答する
- コメントは「なぜ」が自明でないときのみ書く。コードの説明コメントは書かない
- 不要な抽象化・早すぎる汎用化はしない

### Goスタイル

- 命名は Go 標準に従う（`camelCase`、略語は全大文字: `ID`, `URL`, `HTTP`）
- インターフェースは `service/` と `repository/` に定義し、`main.go` で実装を注入する
- レスポンスのJSONエンコードは `writeJSON()` ヘルパー（`internal/handler/todo.go`）を使う

### モデル

- GORMモデルは `model/` に置く
- `gorm.Model` の埋め込みは使わず、必要なフィールドを明示的に定義する（現行の実装に合わせる）
- JSON出力したくないフィールドは `json:"-"` タグを付ける（例: `Password`）

### マイグレーション

- モデル変更時は必ずマイグレーションファイルを生成してから適用する
- `atlas.sum` と `migrations/` は常にセットでコミットする

---

## 注意事項

- `JWT_SECRET` は必須環境変数。未設定の場合サーバーが起動しない
- Atlas のマイグレーションは `atlas-dev-db` コンテナ（Docker内）が必要。ローカル単体では `migrate-diff` が動かない
- `make migrate-reset` はDBを完全にリセットする破壊的操作。本番環境では使用禁止
