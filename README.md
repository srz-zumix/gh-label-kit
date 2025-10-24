# gh-label-kit

`gh-label-kit`は、GitHubのラベルを管理したり、プルリクエストに自動でラベルを付けたりするための`gh`拡張コマンドです。

## インストール

以下のコマンドでインストールできます:

```sh
gh extension install srz-zumix/gh-label-kit
```

## 主なコマンド

### `labeler`: プルリクエストへの自動ラベル付け

変更されたファイルやブランチ名をもとに、プルリクエストへ自動でラベルを付けたり剥がしたりします。
設定は`.github/labeler.yml`ファイルで行います。

```sh
gh label-kit labeler <pr-number>
```

### `repo`: リポジトリのラベル管理

リポジトリ間でラベルをコピーしたり、同期したり、一覧表示したりできます。

- `repo copy <dst-repo>`: ラベルをコピー
- `repo sync <dst-repo>`: ラベルを同期
- `repo list`: ラベルを一覧表示

### `issue`: issueのラベル操作

issueに対してラベルを追加、削除、置換、一覧表示できます。

- `issue add <number> <label>`: ラベルを追加
- `issue remove <number> <label>`: ラベルを削除
- `issue set <number> <label>`: すべてのラベルを置き換え
- `issue list <number>`: ラベルを一覧表示
- `issue clear <number>`: すべてのラベルを削除

### `runner`: GitHub Actions ランナーのラベル一覧

リポジトリやOrganizationのGitHub Actionsランナーのラベルを一覧表示します。

```sh
gh label-kit runner list
```

### `milestone`: マイルストーンのラベル一覧

マイルストーンに含まれるIssueやプルリクエストのラベルを一覧表示します。

```sh
gh label-kit milestone list <milestone>
```

---

より詳細な使い方は、各コマンドの`--help`フラグを参照してください。
例: `gh label-kit labeler --help`
