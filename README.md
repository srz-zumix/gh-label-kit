# gh-label-kit

`gh-label-kit`は、GitHubラベルの管理と自動ラベリングのための`gh`拡張コマンドセットです。

## インストール

以下のコマンドでインストールできます:

```sh
gh extension install srz-zumix/gh-label-kit
```

## コマンドと使用法

各コマンドの詳細なオプションについては、`--help`フラグを参照してください。

---

### labeler: PRへの自動ラベル付け

```sh
gh label-kit labeler <pr-number...>
```

変更されたファイルやブランチ名に基づき、設定ファイル（`.github/labeler.yml`）に従ってプルリクエストに自動でラベルを付け外しします。
このコマンドは[actions/labeler][labeler]と同様に動作します。

設定ファイルについての詳細は[docs/labeler-config.md](docs/labeler-config.md)を参照してください。

---

### repo copy: リポジトリ間でラベルをコピー

```sh
gh label-kit repo copy <dst-repository...>
```

指定したリポジトリから、対象リポジトリへ全てのラベルをコピーします。

---

### repo list: ラベルの一覧表示

```sh
gh label-kit repo list
```

リポジトリの全ラベルを一覧表示します。

---

### repo sync: ラベルの同期

```sh
gh label-kit repo sync <dst-repository...>
```

指定したリポジトリから、対象リポジトリへ全てのラベルを同期します。

---

### runner list: GitHub Actionsランナーのラベルを一覧表示

```sh
gh label-kit runner list
```

リポジトリまたはOrganizationに登録されているGitHub Actionsランナーのラベルを一覧表示します。

---

### issue add: Issue/PRにラベルを追加

```sh
gh label-kit issue add <number> <label>...
```

IssueまたはPRに1つ以上のラベルを追加します。

---

### issue clear: Issue/PRから全てのラベルを削除

```sh
gh label-kit issue clear <number>
```

IssueまたはPRから全てのラベルを削除します。

---

### issue list: Issue/PRのラベルを一覧表示

```sh
gh label-kit issue list <number>
```

IssueまたはPRに付けられたラベルを一覧表示します。

---

### issue remove: Issue/PRからラベルを削除

```sh
gh label-kit issue remove <number> <label>...
```

IssueまたはPRから1つ以上のラベルを削除します。

---

### issue set: Issue/PRのラベルを置換

```sh
gh label-kit issue set <number> <label>...
```

IssueまたはPRの全てのラベルを指定したラベルに置き換えます。

---

### milestone list: マイルストーンのラベルを一覧表示

```sh
gh label-kit milestone list <milestone>
```

指定したマイルストーンに含まれるIssueとPRに付けられたラベルを一覧表示します。

[labeler]: https://github.com/actions/labeler
