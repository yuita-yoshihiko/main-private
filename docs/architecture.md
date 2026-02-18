# Go Backend API — アーキテクチャ設計ドキュメント

## 概要

`main/` ディレクトリに Clean Architecture に基づいた Go バックエンド API のスケルトンを構築する。

| 項目 | 採用技術 |
|---|---|
| 言語 | Go 1.26+ |
| HTTPルーター | [chi v5](https://github.com/go-chi/chi) |
| DB操作 | [sqlc](https://sqlc.dev/) (コード自動生成) |
| データベース | PostgreSQL 18 |
| 認証 | なし |
| 環境構築 | Docker / docker-compose |

---

## ディレクトリ構成

```
main-private/
├── docs/
│   └── architecture.md          ← このファイル
│
├── db-migrater/                  ← マイグレーション専用モジュール（独立したGoモジュール）
│   ├── go.mod
│   ├── main.go                   ← up/down サブコマンドCLI
│   └── migrations/
│       ├── 000001_create_items.up.sql
│       └── 000001_create_items.down.sql
│
└── main/                         ← APIアプリケーション本体
    ├── go.mod
    ├── sqlc.yaml
    ├── Makefile
    ├── Dockerfile
    ├── docker-compose.yml
    ├── .env.example
    ├── .gitignore
    │
    ├── cmd/
    │   └── api/
    │       └── main.go           ← エントリポイント（依存のワイヤリング・サーバー起動）
    │
    ├── internal/
    │   ├── domain/               ← Layer 1: ドメイン層（外部依存ゼロ）
    │   │   └── item/
    │   │       ├── item.go       ← Itemエンティティ・バリデーション・ドメインエラー
    │   │       └── repository.go ← Repositoryインターフェース（ポート）
    │   │
    │   ├── usecase/              ← Layer 2: ユースケース層
    │   │   └── item/
    │   │       ├── usecase.go    ← UseCaseインターフェース
    │   │       └── interactor.go ← ユースケース実装
    │   │
    │   ├── interface/            ← Layer 3: インターフェースアダプター層
    │   │   └── handler/
    │   │       ├── item_handler.go   ← Item CRUD HTTPハンドラー
    │   │       ├── health_handler.go ← GET /health
    │   │       └── response.go       ← JSONレスポンスヘルパー・エラーマッピング
    │   │
    │   └── infrastructure/       ← Layer 4: インフラ層（DB・外部サービス）
    │       ├── db/
    │       │   ├── schema/
    │       │   │   └── schema.sql    ← DDL（sqlc generate用）
    │       │   ├── queries/
    │       │   │   └── item.sql      ← sqlcアノテーション付きSQLクエリ
    │       │   └── sqlc/             ← sqlc generate で自動生成（手動編集禁止）
    │       │       ├── db.go
    │       │       ├── models.go
    │       │       ├── querier.go
    │       │       └── item.sql.go
    │       └── repository/
    │           └── item_repository.go ← Repositoryインターフェース実装（sqlc利用）
    │
    └── pkg/                      ← 層に依存しないユーティリティ
        ├── config/
        │   └── config.go         ← 環境変数 → Config struct
        ├── logger/
        │   └── logger.go         ← slog設定（dev: text、prod: JSON）
        └── server/
            └── server.go         ← chiルーター構築・ミドルウェア登録
```

---

## Clean Architecture 依存関係

```
┌─────────────────────────────────────────────┐
│              Domain Layer                    │
│   entity (Item) + Repository interface       │
│   外部依存ゼロ                                │
└───────────────────┬─────────────────────────┘
                    │ ← 依存方向（内側へ）
┌───────────────────▼─────────────────────────┐
│             UseCase Layer                    │
│   UseCase interface + Interactor             │
│   domain/ のみに依存                          │
└───────────────────┬─────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
┌───────▼───────┐     ┌─────────▼─────────────┐
│ Interface     │     │ Infrastructure Layer   │
│ Handler Layer │     │ Repository (sqlc)      │
│ HTTP ↔ UseCase│     │ domain interfaceを実装 │
└───────────────┘     └───────────────────────┘
```

**依存ルール**: 外側の層は内側の層に依存する。逆方向の依存は禁止。

---

## データフロー（POST /api/v1/items の例）

```
1. HTTPクライアント → POST /api/v1/items
2. chi Router → ItemHandler.Create
3. ItemHandler: JSONデコード → CreateItemRequest DTO
4. ItemHandler: UseCase.CreateItem(ctx, name, description) 呼び出し
5. Interactor: domain.Item 構築 → Validate() でバリデーション
6. Interactor: Repository.Create(ctx, item) 呼び出し
7. ItemRepository: domain.Item → sqlc.CreateItemParams 変換
8. ItemRepository: q.CreateItem(ctx, params) 実行（sqlc生成コード）
9. PostgreSQL: INSERT → 作成レコード返却
10. 上位層へ戻る: sqlcモデル → domainモデル → JSONレスポンスDTO
11. HTTPクライアント ← 201 Created + Item JSON
```

各層の境界でデータ変換（マッピング）が行われる。これにより DB スキーマ変更がハンドラー層に影響しない。

---

## 主要インターフェース

### `internal/domain/item/repository.go`

```go
type Repository interface {
    Create(ctx context.Context, item Item) (Item, error)
    FindByID(ctx context.Context, id uuid.UUID) (Item, error)
    FindAll(ctx context.Context) ([]Item, error)
    Update(ctx context.Context, item Item) (Item, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### `internal/usecase/item/usecase.go`

```go
type UseCase interface {
    CreateItem(ctx context.Context, name, description string) (domain.Item, error)
    GetItem(ctx context.Context, id uuid.UUID) (domain.Item, error)
    ListItems(ctx context.Context) ([]domain.Item, error)
    UpdateItem(ctx context.Context, id uuid.UUID, name, description string) (domain.Item, error)
    DeleteItem(ctx context.Context, id uuid.UUID) error
}
```

---

## API エンドポイント

| Method | Path | 説明 | Status |
|---|---|---|---|
| GET | `/health` | ヘルスチェック | 200 |
| POST | `/api/v1/items` | Item作成 | 201 |
| GET | `/api/v1/items` | Item一覧取得 | 200 |
| GET | `/api/v1/items/{id}` | Item単体取得 | 200 |
| PUT | `/api/v1/items/{id}` | Item更新 | 200 |
| DELETE | `/api/v1/items/{id}` | Item削除 | 204 |

### エラーレスポンスのHTTPステータスマッピング

| エラー | HTTPステータス |
|---|---|
| `domain.ErrNotFound` | 404 Not Found |
| `domain.ErrNameRequired` | 400 Bad Request |
| 不正なUUID（パスパラメータ） | 400 Bad Request |
| `json.SyntaxError` | 400 Bad Request |
| その他 | 500 Internal Server Error |

---

## 依存パッケージ

### `main/go.mod`

| パッケージ | 用途 |
|---|---|
| `github.com/go-chi/chi/v5` | HTTPルーター |
| `github.com/jackc/pgx/v5` | PostgreSQLドライバ（pgxpool） |
| `github.com/google/uuid` | UUID生成・パース |

標準ライブラリ使用: `encoding/json`, `log/slog`, `context`, `os`

### `db-migrater/go.mod`

| パッケージ | 用途 |
|---|---|
| `github.com/golang-migrate/migrate/v4` | マイグレーション実行 |
| `github.com/jackc/pgx/v5` | PostgreSQL接続 |

---

## sqlc 設定（`main/sqlc.yaml`）

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "internal/infrastructure/db/schema/schema.sql"
    queries: "internal/infrastructure/db/queries/item.sql"
    gen:
      go:
        package: "sqlcdb"
        out: "internal/infrastructure/db/sqlc"
        sql_package: "pgx/v5"
        emit_interface: true      # Querierインターフェース生成（テスト容易性）
        emit_json_tags: true
        json_tags_case_style: "camel"
        emit_empty_slices: true   # 空リストを null でなく [] で返す
        overrides:
          - db_type: "uuid"
            go_type: { import: "github.com/google/uuid", type: "UUID" }
          - db_type: "pg_catalog.timestamptz"
            go_type: { import: "time", type: "Time" }
```

---

## Docker 構成

```
docker-compose.yml
├── app           ← Goアプリ（マルチステージビルド、最終イメージ ~15MB）
│   └── depends_on: postgres (condition: service_healthy)
└── postgres      ← PostgreSQL 17 Alpine
    └── healthcheck: pg_isready
```

- `app` コンテナは postgres の healthcheck が通るまで起動を待機
- `postgres_data` ボリュームにより `docker-compose down` してもデータは保持
- データも削除する場合は `docker-compose down -v`

---

## Makefile ターゲット

| コマンド | 説明 |
|---|---|
| `make run` | アプリ起動（`go run ./cmd/api/...`） |
| `make build` | バイナリビルド（`bin/server`） |
| `make test` | テスト実行（race detector付き） |
| `make sqlc-gen` | `sqlc generate` でDB操作コードを再生成 |
| `make migrate-up` | マイグレーション適用 |
| `make migrate-down` | マイグレーション1段階ロールバック |
| `make migrate-create` | 新しいマイグレーションファイルペアを作成 |
| `make docker-up` | Docker起動（`--build`付き） |
| `make docker-down` | Docker停止（データ保持） |
| `make docker-down-v` | Docker停止（ボリューム削除） |
| `make docker-logs` | appコンテナのログをフォロー |

---

## 実装順序

1. ブランチ作成: `feat-go-clean-architecture-skeleton`
2. `db-migrater/` セットアップ（go.mod, migration SQL, main.go）
3. `main/internal/infrastructure/db/schema/schema.sql` 作成
4. `main/internal/infrastructure/db/queries/item.sql` 作成
5. `main/sqlc.yaml` 作成 → `sqlc generate` 実行
6. `main/go.mod` + `go get` で依存解決
7. `internal/domain/item/` 作成（entity + repository interface）
8. `internal/usecase/item/` 作成（usecase interface + interactor）
9. `pkg/config/`, `pkg/logger/` 作成
10. `internal/infrastructure/repository/item_repository.go` 作成
11. `internal/interface/handler/` 作成
12. `pkg/server/server.go` 作成
13. `cmd/api/main.go` 作成（全依存をワイヤリング）
14. `Dockerfile`, `docker-compose.yml`, `.env.example`, `.gitignore`, `Makefile` 作成

---

## 動作確認手順

```bash
# Docker起動・マイグレーション適用
make docker-up
make migrate-up

# ヘルスチェック
curl http://localhost:8080/health
# → {"status":"ok"}

# Item作成
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name":"test item","description":"hello"}'
# → 201 + Item JSON

# 一覧取得
curl http://localhost:8080/api/v1/items
# → 200 + Item配列 JSON
```

---

## 注意事項

- `sqlc/` 以下は自動生成ファイルのため**手動編集禁止**（`sqlc generate` で上書きされる）
- `schema.sql`（sqlc用）と `db-migrater/migrations/` のDDLは**常に同期を保つ**
- `.env` は `.gitignore` に追加必須（`.env.example` のみコミット）
- pgx/v5 の UUID型は `sqlc.yaml` の override で `uuid.UUID` にマップ
- パスパラメータのUUIDは必ず `uuid.Parse()` で検証し、不正なら400を返す
- `db-migrater` は独立したGoモジュールとして管理（メインアプリに migrate 依存を持ち込まない）
