# コーディング規約

## 全般

- エラーは握りつぶさない。必ず呼び出し元に返すか、適切なHTTPステータスで応答する
- コメントは「なぜ」が自明でないときのみ書く。コードの説明コメントは書かない
- 不要な抽象化・早すぎる汎用化はしない

## テスト

- handler / service / repository のいずれかを追加・変更した場合、対応するテストを必ず追加または修正する
- テストは `make test`（`go test ./...`）で全件パスすることを確認してから作業完了とする

### 層ごとのテスト方針

| 層 | 方針 | 理由 |
| -- | ---- | ---- |
| repository | `test-db`（PostgreSQL）で実際にDB操作を検証する | 本番と同じエンジンで方言の差異を防ぐ。SQLite は使わない |
| service | `UserRepository` 等を手書きモックに置き換えてテストする | DBに依存せず、ビジネスロジックに集中できる |
| handler | `XxxService` を手書きモックに置き換え、`httptest` でHTTPレスポンスを検証する | サービス層に依存せず、HTTPの入出力に集中できる |

### repository テストの規則

- `TestMain` で `test-db` に接続し `AutoMigrate` を実行する（パッケージ全体で1回）
- 各テスト関数の冒頭で対象テーブルを `TRUNCATE ... RESTART IDENTITY CASCADE` してデータを分離する
- `test-db` コンテナは `docker compose up -d test-db` で起動する。`backend` コンテナは起動時に自動で `depends_on` する

### モックの書き方

外部ライブラリは使わず、インターフェースを実装した構造体を手書きする。

```go
type mockUserRepo struct {
    findByIDFn func(id uint) (*model.User, error)
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
    return m.findByIDFn(id)
}
```

## Goスタイル

- 命名は Go 標準に従う（`camelCase`、略語は全大文字: `ID`, `URL`, `HTTP`）
- インターフェースは `service/` と `repository/` に定義し、`main.go` で実装を注入する
- レスポンスのJSONエンコードは `writeJSON()` ヘルパーを使う

## モデル

- GORMモデルは `model/` に置く
- `gorm.Model` の埋め込みは使わず、必要なフィールドを明示的に定義する
- JSON出力したくないフィールドは `json:"-"` タグを付ける（例: `Password`）

## マイグレーション

- モデル変更時は必ずマイグレーションファイルを生成してから適用する
- `atlas.sum` と `migrations/` は常にセットでコミットする

## ハンドラー

### リクエスト

- リクエストボディは必ず専用の構造体（`xxxRequest`）で受け取る
- 入力バリデーションはリクエスト構造体の `validate() error` メソッドに実装する
- ハンドラー内でバリデーションロジックを直接書かない

```go
type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (r loginRequest) validate() error {
    if r.Email == "" || r.Password == "" {
        return errors.New("email and password are required")
    }
    return nil
}
```

### レスポンス

- モデルをそのまま返す場合は専用のレスポンス構造体は不要
- モデル以外の構造でレスポンスを返す場合は専用の構造体（`xxxResponse`）を定義する
- `map[string]string` 等をレスポンスに直接使わない

```go
// モデル以外の構造 → 構造体を定義する
type loginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// モデルをそのまま返す場合 → 構造体不要
writeJSON(w, http.StatusOK, user)
```
