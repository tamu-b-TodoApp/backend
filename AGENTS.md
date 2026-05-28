# AGENTS.md — TodoApp Backend 作業ガイド

このファイルはAIエージェント（Claude等）がこのリポジトリで作業する際の指示書です。
**作業開始前に必ず以下を読んでください。**

1. このファイル（AGENTS.md）— 作業ルール・ガイドライン
2. [docs/coding-conventions.md](docs/coding-conventions.md) — コーディング規約
3. [docs/api.md](docs/api.md) — API仕様（エンドポイント追加・変更時に更新すること）

指示で不明点がある場合、推測せず、必ず実装者に確認して下さい。

また、仕様として残すべきものはdocsフォルダに残して下さい。

---

## プロジェクト概要

Go製のTodoAppバックエンドAPI。会社・チーム・従業員（company_members）を管理し、チームごとに階層構造のTodoを扱うRESTful API。

| 項目             | 内容                                            |
| ---------------- | ----------------------------------------------- |
| 言語             | Go 1.26                                         |
| ORM              | GORM                                            |
| DB               | PostgreSQL                                      |
| マイグレーション | Atlas (atlas.hcl)                               |
| 認証             | JWT (golang-jwt/jwt/v5) + bcrypt                |
| HTTPサーバー     | 標準ライブラリ `net/http`（フレームワークなし） |

---

## ドメインモデル概要

- **会社 (companies)** を中心に、チームと従業員が属する
- **ユーザー (users)** は複数の会社に所属できる。会社内での所属は **company_members** で管理する
- **チーム (teams)** は会社に属する。チームメンバーは **team_members** で管理する（company_members との多対多）
- **Todo (todos)** はチームに属し、`parent_id` による自己参照で階層構造を持つ。担当者は team_members の中から設定する

詳細は [docs/er-diagram.md](docs/er-diagram.md) を参照。

---

## ディレクトリ構成

```
backend/
├── main.go                    # エントリーポイント。DI（依存注入）はここで行う
├── model/                     # GORMモデル定義
├── internal/
│   ├── handler/               # HTTPハンドラー（リクエスト/レスポンス処理）
│   ├── service/               # ビジネスロジック
│   ├── repository/            # DB操作
│   ├── middleware/            # 認証ミドルウェア等
│   └── router/                # ルーティング登録
├── migrations/                # Atlasが管理するSQLマイグレーションファイル
├── cmd/
│   ├── atlas/main.go          # Atlas用スキーマ出力コマンド
│   ├── seed/main.go           # シードデータ投入
│   └── gen-jwt-secret/main.go # JWT秘密鍵生成
├── docs/                      # ドキュメント
│   └── er-diagram.md          # ER図
├── atlas.hcl                  # Atlas設定（local / prod 環境）
├── docker-compose.yml         # ローカル開発環境
└── Makefile                   # よく使うコマンド
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
| `make test`                               | 全テスト実行                        |
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

## APIエンドポイント一覧

※ 実装に合わせて随時更新する。

### 認証不要

| メソッド | パス            | 説明                   |
| -------- | --------------- | ---------------------- |
| `GET`    | `/health`       | ヘルスチェック         |
| `POST`   | `/auth/login`   | ログイン（JWT発行）    |
| `POST`   | `/auth/refresh` | アクセストークン再発行 |

### 認証必要（`Authorization: Bearer <access_token>`）

| メソッド | パス                                              | 説明                   |
| -------- | ------------------------------------------------- | ---------------------- |
| `GET`    | `/auth/me`                                        | 認証ユーザー情報取得   |
| `GET`    | `/companies`                                      | 会社一覧取得           |
| `POST`   | `/companies`                                      | 会社作成               |
| `GET`    | `/companies/{id}`                                 | 会社取得               |
| `PUT`    | `/companies/{id}`                                 | 会社更新               |
| `DELETE` | `/companies/{id}`                                 | 会社削除               |
| `GET`    | `/companies/{id}/teams`                           | チーム一覧取得         |
| `GET`    | `/companies/{companyID}/teams/{teamID}/todos`     | Todo一覧取得           |
| `POST`   | `/companies/{companyID}/teams/{teamID}/todos`     | Todo作成               |
| `GET`    | `/todos/{id}`                                     | Todo取得               |
| `PUT`    | `/todos/{id}`                                     | Todo更新               |
| `DELETE` | `/todos/{id}`                                     | Todo削除               |

### 認証の仕組み

- アクセストークン: JWT（`type: "access"`）、有効期限15分
- リフレッシュトークン: JWT（`type: "refresh"`）、有効期限7日、DBに保存しない（ステートレス）
- ログアウト: サーバー側APIなし。フロントエンドでトークンを削除する
- ミドルウェア: `internal/middleware/auth.go`

---

## 環境変数

`.env` ファイルに設定する。`Makefile` が自動で読み込む。

| 変数名 | 必須 | 説明 |
| ------ | ---- | ---- |
| `DB_HOST` | ✓ | PostgreSQL ホスト名（Docker内では `db`） |
| `DB_USER` | ✓ | PostgreSQL ユーザー名 |
| `DB_PASSWORD` | ✓ | PostgreSQL パスワード |
| `DB_NAME` | ✓ | PostgreSQL データベース名 |
| `JWT_SECRET` | ✓ | JWT署名鍵。`make gen-jwt-secret` で生成する |
| `TEST_DB_HOST` | - | テスト用DB ホスト名。`docker-compose.yml` が `backend` コンテナへ自動注入するため `.env` への記載不要 |
| `TEST_DB_NAME` | - | テスト用DB 名。同上 |

---

## テスト

```bash
docker compose exec backend make test
```

### コンテナ構成

| コンテナ | 用途 |
| -------- | ---- |
| `db` | 開発用 PostgreSQL |
| `test-db` | テスト専用 PostgreSQL（DB名: `test`） |
| `atlas-dev-db` | マイグレーションファイル生成用 PostgreSQL |

- repository テストは `test-db` に接続する。`TestMain` で `AutoMigrate` + 各テスト前に `TRUNCATE` を行い分離する
- service / handler テストは手書きモックを使う。外部依存なしで実行できる
- `make test`（`go test ./...`）で全件パスすることを確認してから作業完了とする

---

## ドキュメント

| ファイル | 内容 |
| -------- | ---- |
| [docs/api.md](docs/api.md) | 全エンドポイントのリクエスト・レスポンス仕様 |
| [docs/er-diagram.md](docs/er-diagram.md) | ER図 |
| [docs/coding-conventions.md](docs/coding-conventions.md) | コーディング規約 |

---

## 注意事項

- `JWT_SECRET` は必須環境変数。未設定の場合サーバーが起動しない
- Atlas のマイグレーションは `atlas-dev-db` コンテナ（Docker内）が必要。ローカル単体では `migrate-diff` が動かない
- `make migrate-reset` はDBを完全にリセットする破壊的操作。本番環境では使用禁止
- Todoのアサインは**チームメンバーのみ**設定可能。サービス層で検証すること
