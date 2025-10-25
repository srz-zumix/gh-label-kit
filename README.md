# gh-label-kit

`gh-label-kit`は、GitHubラベルの管理と自動ラベル付けのための一連の`gh`拡張コマンドです。

## インストール

このツールをインストールするには、次のコマンドを使用します。

```sh
gh extension install srz-zumix/gh-label-kit
```

## コマンドと使用法

---

### labeler: PRの自動ラベル付け

```sh
gh label-kit labeler <pr-number...> [--repo <owner/repo>] [--config <path>] [--sync] [--dryrun] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>] [--name-only] [--ref <string>]
```

変更されたファイル、ブランチ名、およびYAML設定ファイル（デフォルト：.github/labeler.yml）に基づいて、GitHubプルリクエストにラベルを自動的に追加または削除します。
glob/regexパターンとラベル削除のためのsyncLabelsオプションをサポートしています。このコマンドは[actions/labeler][labeler]と同じように動作します。

- --color: diff出力で色を使用します（auto|never|always、デフォルト：auto）
- --config: ラベラー設定YAMLファイルへのパス（デフォルト：.github/labeler.yml）
- --dryrun/-n: ドライラン：実際にはラベルを設定しません
- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --name-only: チーム名のみを出力します
- --ref: リポジトリから設定をロードするためのGit参照（ブランチ、タグ、またはコミットSHA）
- --repo/-R: 「owner/repo」形式のターゲットリポジトリ
- --sync: どの条件にも一致しないラベルを削除します
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

`labeler`コマンドは、YAML設定ファイルを使用してラベリングルールを定義します。設定形式は[actions/labeler][labeler]と互換性があり、`color`および`codeowners`機能の追加サポートがあります。

詳細な設定ドキュメントについては、[docs/labeler-config.md](docs/labeler-config.md)を参照してください。

---

### repo copy: リポジトリ間でラベルをコピー

```sh
gh label-kit repo copy <dst-repository...> [--repo <owner/repo>] [--force]
```

ソースリポジトリから宛先リポジトリにすべてのラベルをコピーします。宛先にラベルが既に存在する場合、--forceが指定されていない限りスキップされます。

- --force/-f: 宛先リポジトリの既存のラベルを上書きします
- --repo/-R: 「owner/repo」形式のリポジトリ（ソースリポジトリ）

---

### repo list: ラベルの一覧表示

```sh
gh label-kit repo list [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

指定されたリポジトリのすべてのラベルを一覧表示します。

- --color: diff出力で色を使用します（always|never|auto、デフォルト：auto）
- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --repo/-R: 「owner/repo」形式のリポジトリ
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

---

### repo sync: ラベルの差分を同期

```sh
gh label-kit repo sync <dst-repository...> [--repo <owner/repo>] [--force]
```

ソースリポジトリから宛先リポジトリにすべてのラベルを同期します。宛先にラベルが既に存在する場合、--forceが指定されていれば更新されます。

- --force/-f: 宛先リポジトリの既存のラベルを上書きします
- --repo/-R: 「owner/repo」形式のリポジトリ（ソースリポジトリ）

---

### runner list: GitHub Actionsランナーのラベルを一覧表示

```sh
gh label-kit runner list [--repo <owner/repo>] [--owner <organization>] [--format <json>] [--jq <expression>] [--template <string>]
```

指定されたリポジトリのすべてのGitHub Actionsランナーのラベルを一覧表示します。

- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --owner: 組織名を指定します
- --repo/-R: 「owner/repo」形式のリポジトリ
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

---

### issue add: issueにラベルを追加

```sh
gh label-kit issue add <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

リポジトリのissueに1つ以上のラベルを追加します。

- --color: diff出力で色を使用します（always|never|auto、デフォルト：auto）
- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --repo/-R: 「owner/repo」形式のリポジトリ
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

---

### issue clear: issueからすべてのラベルを削除

```sh
gh label-kit issue clear <number> [--repo <owner/repo>]
```

リポジトリのissueからすべてのラベルを削除します。

- --repo/-R: 「owner/repo」形式のリポジトリ

---

### issue list: issueのラベルを一覧表示

```sh
gh label-kit issue list <number> [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

リポジトリのissueに添付されているすべてのラベルを一覧表示します。

- --color: diff出力で色を使用します（always|never|auto、デフォルト：auto）
- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --repo/-R: 「owner/repo」形式のリポジトリ
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

---

### issue remove: issueからラベルを削除

```sh
gh label-kit issue remove <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

リポジトリのissueから1つ以上のラベルを削除します。

- --color: diff出力で色を使用します（always|never|auto、デフォルト：auto）
- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --repo/-R: 「owner/repo」形式のリポジトリ
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

---

### issue set: issueのラベルを設定（すべて置き換え）

```sh
gh label-kit issue set <number> <label>... [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

リポジトリのissueのすべてのラベルを設定（置き換え）します。

- --color: diff出力で色を使用します（always|never|auto、デフォルト：auto）
- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --repo/-R: 「owner/repo」形式のリポジトリ
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

---

### milestone list: マイルストーンのラベルを一覧表示

```sh
gh label-kit milestone list <milestone> [--repo <owner/repo>] [--color <auto|always|never>] [--format <json>] [--jq <expression>] [--template <string>]
```

指定されたマイルストーン内のissueとPRに添付されているすべてのラベルを一覧表示します。

- --color: diff出力で色を使用します（always|never|auto、デフォルト：auto）
- --format: 出力形式（json）
- --jq: jq式を使用してJSON出力をフィルタリングします
- --repo/-R: 「owner/repo」形式のリポジトリ
- --template/-t: Goテンプレートを使用してJSON出力をフォーマットします

[labeler]: https://github.com/actions/labeler
