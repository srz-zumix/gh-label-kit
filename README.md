# gh-label-kit

`gh-label-kit`は、GitHubのラベル管理と自動ラベリングのための`gh`拡張コマンドセットです。

## インストール

以下のコマンドでツールをインストールできます:

```sh
gh extension install srz-zumix/gh-label-kit
```

## コマンドと使用方法

---

### `labeler`: PRの自動ラベリング

```sh
gh label-kit labeler <pr-number...> [--repo <owner/repo>] [--config <path>] [--sync] [--dryrun] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>] [--name-only] [--ref <string>]
```

変更されたファイル、ブランチ名、およびYAML設定ファイル（デフォルト: `.github/labeler.yml`）に基づいて、GitHubプルリクエストにラベルを自動的に追加または削除します。
このコマンドは[actions/labeler][labeler]と同じように動作します。

**オプション:**
- `--config`: ラベラー設定YAMLファイルへのパス（デフォルト: `.github/labeler.yml`）
- `--dryrun/-n`: ドライラン：実際にはラベルを設定しません
- `--sync`: どの条件にも一致しないラベルを削除します
- その他、出力フォーマットやリポジトリ指定などのオプションがあります。

詳細な設定については、[docs/labeler-config.md](docs/labeler-config.md)を参照してください。

---

### リポジトリのラベル管理 (`repo`)

#### `repo copy`: リポジトリ間でラベルをコピー

```sh
gh label-kit repo copy <dst-repository...> [--repo <owner/repo>] [--force]
```
ソースリポジトリから宛先リポジトリにすべてのラベルをコピーします。

#### `repo list`: ラベルの一覧表示

```sh
gh label-kit repo list [--repo <owner/repo>] [--format <json>]
```
指定されたリポジトリのすべてのラベルを一覧表示します。

#### `repo sync`: ラベルの同期

```sh
gh label-kit repo sync <dst-repository...> [--repo <owner/repo>] [--force]
```
ソースリポジトリから宛先リポジトリにすべてのラベルを同期します。

---

### issueのラベル管理 (`issue`)

#### `issue add`: issueにラベルを追加

```sh
gh label-kit issue add <number> <label>... [--repo <owner/repo>]
```

#### `issue clear`: issueからすべてのラベルを削除

```sh
gh label-kit issue clear <number> [--repo <owner/repo>]
```

#### `issue list`: issueのラベルを一覧表示

```sh
gh label-kit issue list <number> [--repo <owner/repo>]
```

#### `issue remove`: issueからラベルを削除

```sh
gh label-kit issue remove <number> <label>... [--repo <owner/repo>]
```

#### `issue set`: issueのラベルをすべて置き換え

```sh
gh label-kit issue set <number> <label>... [--repo <owner/repo>]
```

---

### その他のコマンド

#### `runner list`: GitHub Actionsランナーのラベルを一覧表示

```sh
gh label-kit runner list [--repo <owner/repo>] [--owner <organization>]
```

#### `milestone list`: マイルストーンのラベルを一覧表示

```sh
gh label-kit milestone list <milestone> [--repo <owner/repo>]
```
指定されたマイルストーン内のissueとPRに付けられたすべてのラベルを一覧表示します。

[labeler]: https://github.com/actions/labeler
