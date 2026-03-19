---
name: gh-issue-search
description: |
  This skill should be used when the user asks to "Issue一覧を見せて", "Issueを検索して",
  "Issue確認して", "Issueの詳細を見せて", "バグ一覧", "アサインされたIssue",
  "Issue.mdを作って", "Issueをまとめて", "/gh-issue-search",
  or needs to search, view, or summarize GitHub Issues via gh CLI.
---

# GitHub Issue 検索・取得・変換スキル

gh CLI を使って GitHub Issue の検索・詳細取得・Issue.md 変換を行います。

## Mode Selection

ユーザーの指示から以下の3モードを自動判定して実行する。

| モード | トリガー例 | 主なコマンド |
|--------|-----------|-------------|
| 検索・一覧 | 「Issue一覧を見せて」「bugラベルのIssue」 | `gh issue list` |
| 詳細取得 | 「Issue #3 の詳細」「コメントも見せて」 | `gh issue view` |
| Issue.md 変換 | 「Issue.md を作って」「進捗をまとめて」 | `gh issue view` → テンプレート出力 |

## Token Efficiency Rules

**必ず守ること:**

- `gh issue list` は `--limit 20` をデフォルトとする（全件取得しない）
- `--comments` は詳細取得・Issue.md 変換モード時のみ使用
- JSON 出力 + jq で必要なフィールドのみ取得する

---

## Mode 1: 検索・一覧

Issue をフィルタ条件付きで一覧取得する。

### Step 1: フィルタ条件の組み立て

ユーザーの指示からオプションを自動判定する。

```bash
# 基本一覧
gh issue list --limit 20

# ラベルフィルタ
gh issue list --label bug --state open

# アサインフィルタ
gh issue list --assignee @me

# ステートフィルタ
gh issue list --state closed

# キーワード検索
gh issue list --search "robots.txt"

# 複合条件
gh issue list --assignee @me --label bug --state open

# JSON 出力（構造化データ）
gh issue list --json number,title,state,labels,assignees,createdAt --limit 20
```

### Step 2: 結果を一覧テンプレートで表示

```
## Issue 一覧

| # | タイトル | ラベル | アサイン | 状態 |
|---|---------|--------|---------|------|
| <number> | <title> | <labels> | <assignees> | <state> |
```

## Mode 2: 詳細取得

特定の Issue の詳細情報を取得する。

### Step 1: Issue の詳細取得

```bash
# 基本詳細
gh issue view <number>

# コメント付き
gh issue view <number> --comments

# JSON 出力
gh issue view <number> --json title,body,state,labels,assignees,author,comments,createdAt
```

### Step 2: 内容を整理して報告

## Mode 3: Issue.md 変換

Issue の内容を取得し、構造化されたドキュメントに変換する。

### Step 1: Issue データの取得

```bash
# 本文 + メタデータ
gh issue view <number> --json title,body,state,labels,assignees,author,createdAt

# コメント付き（進捗まとめの場合）
gh issue view <number> --json title,body,state,labels,assignees,author,comments,createdAt
```

コメント取得時は投稿者情報を必ず含める:

```bash
gh issue view <number> --json comments --jq '.comments[] | {author: .author.login, createdAt, body}'
```

### Step 2: テンプレートに変換

ユーザーが独自テンプレートを指定した場合はそれに従う。未指定時は以下を使用。

#### バグ・課題の場合（デフォルト）

```markdown
# Issue #<number>: <title>

- URL: <url>
- 起票者: @<author>
- 担当: @<assignees>
- 状態: <state>
- ラベル: <labels>
- 作成日時: <createdAt>

## 問題
（Issue のタイトルと本文から要約）

## 仮説
（考えられる原因を2-3個）

## 関連ファイル
（調査すべきファイルパス）

## 再現手順
（Issue の再現手順をそのまま転記）
```

#### 進捗まとめの場合（コメントあり）

```markdown
# Issue #<number>: <title>

- URL: <url>
- 起票者: @<author>
- 担当: @<assignees>
- 状態: <state>
- ラベル: <labels>
- 作成日時: <createdAt>
- 更新日時: <updatedAt>

## タスク概要
（Issue の本文から）

## 進捗状況
- **@<author>** (YYYY-MM-DD): コメント内容の要約
- **@<author>** (YYYY-MM-DD): コメント内容の要約

## 残タスク
（未完了の項目）
```

## JSON Output Patterns

よく使う構造化データ取得パターン:

```bash
# ラベル名のみ抽出
gh issue list --json number,title,labels --jq '.[] | {number, title, labels: [.labels[].name]}'

# アサイン名のみ抽出
gh issue list --json number,title,assignees --jq '.[] | {number, title, assignees: [.assignees[].login]}'

# コメント数付き一覧
gh issue list --json number,title,comments --jq '.[] | {number, title, commentCount: (.comments | length)}'
```

## Important Notes

- **gh CLI 未認証の場合**: `gh auth status` で認証状態を確認し、未認証ならユーザーに通知
- **リポジトリ外で実行した場合**: エラーをユーザーに通知
- **大量の Issue**: `--limit` で絞り込み、必要に応じて `--search` で検索条件を追加
- **コメントが多い Issue**: 主要な進捗のみ要約し、全文転記は避ける
- **出力先**: Issue.md 変換時、ユーザーが指定しない場合は `docs/bugs/issue-<number>.md` に保存
