# API仕様

## 共通仕様

- Content-Type: `application/json`
- 認証が必要なエンドポイントは `Authorization: Bearer <access_token>` ヘッダーが必要
- エラーレスポンスはすべてプレーンテキスト（`text/plain`）

### HTTPステータスコード

| コード | 意味 |
| ------ | ---- |
| 200 | 成功 |
| 201 | 作成成功 |
| 204 | 成功（レスポンスボディなし） |
| 400 | リクエスト不正（バリデーションエラー等） |
| 401 | 認証エラー |
| 404 | リソースが見つからない |
| 500 | サーバー内部エラー |

---

## 認証不要エンドポイント

### GET /health

ヘルスチェック。

**レスポンス 200**

```json
{ "status": "ok" }
```

---

### POST /auth/login

メールアドレスとパスワードでログインし、JWT トークンを発行する。

**リクエスト**

```json
{
  "email": "user@example.com",
  "password": "plaintext_password"
}
```

| フィールド | 型 | 必須 | 説明 |
| --------- | -- | ---- | ---- |
| email | string | ✓ | メールアドレス |
| password | string | ✓ | パスワード（平文） |

**レスポンス 200**

```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

| フィールド | 説明 |
| --------- | ---- |
| access_token | 有効期限15分。認証が必要なAPIに使用する |
| refresh_token | 有効期限7日。アクセストークン再発行に使用する |

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | email または password が空 |
| 401 | 認証情報が正しくない |

---

### POST /auth/refresh

リフレッシュトークンを使ってアクセストークンを再発行する。

**リクエスト**

```json
{
  "refresh_token": "eyJ..."
}
```

**レスポンス 200**

```json
{
  "access_token": "eyJ..."
}
```

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | refresh_token が空 |
| 401 | トークンが無効または期限切れ、あるいはアクセストークンを渡した場合 |

---

## 認証必要エンドポイント

以降すべてのエンドポイントに `Authorization: Bearer <access_token>` が必要。

---

### GET /auth/me

認証中のユーザー情報を返す。

**レスポンス 200**（User オブジェクト）

```json
{
  "id": 1,
  "email": "user@example.com",
  "email_verified_at": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

> `password` フィールドは JSON に含まれない。

---

## 会社 (companies)

### GET /companies

会社一覧を返す。

**レスポンス 200**

```json
[
  {
    "id": 1,
    "name": "ACME Corp",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

---

### POST /companies

会社を作成する。

**リクエスト**

```json
{ "name": "ACME Corp" }
```

| フィールド | 型 | 必須 | 説明 |
| --------- | -- | ---- | ---- |
| name | string | ✓ | 会社名 |

**レスポンス 201**（作成した Company オブジェクト）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | name が空 |

---

### GET /companies/{id}

指定した会社を返す。

**レスポンス 200**（Company オブジェクト）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | id が数値でない |
| 404 | 会社が存在しない |

---

### PUT /companies/{id}

会社情報を更新する。

**リクエスト**

```json
{ "name": "New Name" }
```

**レスポンス 200**（更新後の Company オブジェクト）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | id が数値でない / name が空 |
| 404 | 会社が存在しない |

---

### DELETE /companies/{id}

会社を削除する。

**レスポンス 204**（ボディなし）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | id が数値でない |

---

## チーム (teams)

### GET /companies/{id}/teams

指定した会社のチーム一覧を返す。会社が存在しない場合も空配列を返す。

**レスポンス 200**

```json
[
  {
    "id": 10,
    "company_id": 1,
    "name": "開発チーム",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | id が数値でない |

---

## Todo (todos)

### Todo オブジェクト

```json
{
  "id": 100,
  "team_id": 10,
  "parent_id": null,
  "assignee_id": null,
  "title": "タスク名",
  "description": "詳細説明",
  "status": "not_started",
  "due_date": null,
  "story_points": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

| フィールド | 型 | 説明 |
| --------- | -- | ---- |
| id | number | 主キー |
| team_id | number | 所属チームID |
| parent_id | number \| null | 親Todo の id（階層構造） |
| assignee_id | number \| null | 担当者の `company_members.id`（チームメンバーのみ設定可） |
| title | string | タイトル |
| description | string | 説明 |
| status | string | `not_started` / `in_progress` / `completed` |
| due_date | string \| null | 期限日（RFC3339形式） |
| story_points | number \| null | ストーリーポイント（工数見積もり） |

---

### GET /companies/{companyID}/teams/{teamID}/todos

指定チームのTodo一覧を返す。`{teamID}` が `{companyID}` の会社に属していない場合は404。

**レスポンス 200**（Todo 配列）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | companyID / teamID が数値でない |
| 404 | チームが存在しない、またはそのチームが指定会社に属していない |

---

### POST /companies/{companyID}/teams/{teamID}/todos

指定チームに Todo を作成する。

**リクエスト**

```json
{
  "title": "タスク名",
  "description": "詳細説明",
  "parent_id": null,
  "assignee_id": null,
  "status": "not_started",
  "due_date": "2024-06-30T00:00:00Z",
  "story_points": 3
}
```

| フィールド | 型 | 必須 | 説明 |
| --------- | -- | ---- | ---- |
| title | string | ✓ | タイトル |
| description | string | - | 説明（省略時は空文字） |
| parent_id | number | - | 親Todo の id |
| assignee_id | number | - | 担当者の `company_members.id`（チームメンバーのみ） |
| status | string | - | 省略時は `not_started` |
| due_date | string | - | RFC3339形式 |
| story_points | number | - | 正の整数 |

**レスポンス 201**（作成した Todo オブジェクト）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | title が空 / assignee がそのチームのメンバーでない |
| 404 | チームが存在しない、またはそのチームが指定会社に属していない |

---

### GET /todos/{id}

Todo を1件取得する。

**レスポンス 200**（Todo オブジェクト）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | id が数値でない |
| 404 | Todo が存在しない |

---

### PUT /todos/{id}

Todo を更新する。リクエストボディはPOSTと同じ形式。

**レスポンス 200**（更新後の Todo オブジェクト）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | id が数値でない / title が空 / assignee がチームメンバーでない |
| 404 | Todo が存在しない |

---

### DELETE /todos/{id}

Todo を削除する。

**レスポンス 204**（ボディなし）

**エラー**

| ステータス | 条件 |
| ---------- | ---- |
| 400 | id が数値でない |
