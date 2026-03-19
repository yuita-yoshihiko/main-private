# Issue #11: feat: Statsを作成するエンドポイントの作成

- URL: https://github.com/yuita-yoshihiko/main-private/issues/11
- 起票者: @yuita-yoshihiko
- 担当: なし
- 状態: OPEN
- ラベル: なし
- 作成日時: 2026-03-19T13:15:37Z

## タスク概要

Statsを作成するエンドポイントを作成する。

## 実装方針

`POST /api/v1/stats` エンドポイントを追加する。
呼び出し時点の item 総数を取得し、`stats` テーブルにスナップショットとして保存して返す。

Clean Architecture に従い、各層を分離して実装する。

## 関連ファイル

| ファイル | 役割 |
|---|---|
| `db-migrater/migrations/000002_create_stats.up.sql` | stats テーブル作成マイグレーション |
| `main/internal/infrastructure/db/schema/schema.sql` | stats テーブル定義追加 |
| `main/internal/infrastructure/db/queries/stats.sql` | SQL クエリ定義 |
| `main/internal/infrastructure/db/sqlc/stats.sql.go` | sqlc 生成コード |
| `main/internal/domain/stats/stats.go` | Stats ドメインモデル |
| `main/internal/domain/stats/repository.go` | Repository インターフェース |
| `main/internal/usecase/stats/usecase.go` | UseCase インターフェース |
| `main/internal/usecase/stats/interactor.go` | ユースケース実装 |
| `main/internal/infrastructure/repository/stats_repository.go` | DB リポジトリ実装 |
| `main/internal/interface/handler/stats_handler.go` | HTTP ハンドラー |
| `main/pkg/server/server.go` | ルーター登録 |
| `main/cmd/api/main.go` | DI 配線 |

## API

### POST /api/v1/stats

リクエストボディ: なし

レスポンス例:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "total_items": 5,
  "created_at": "2026-03-19T13:15:37Z"
}
```

## 残タスク

- [x] エンドポイントの設計
- [x] 実装
- [ ] テスト（DB 接続後に動作確認）
