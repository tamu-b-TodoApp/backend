# コーディング規約

## 全般

- エラーは握りつぶさない。必ず呼び出し元に返すか、適切なHTTPステータスで応答する
- コメントは「なぜ」が自明でないときのみ書く。コードの説明コメントは書かない
- 不要な抽象化・早すぎる汎用化はしない

## テスト

- handler / service / repository のいずれかを追加・変更した場合、対応するテストを必ず追加または修正する
- テストは `make test`（`go test ./...`）で全件パスすることを確認してから作業完了とする
- repository テストは `test-db`（PostgreSQL）を使う。

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
