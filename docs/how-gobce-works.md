# gobce の仕組みをやさしく説明

このドキュメントは、`gobce` が内部的に何をやっているのかをできるだけ平易に説明するためのものです。

## まず一言でいうと

`gobce` は次の 2 つを組み合わせて、分岐カバレッジ (C1) を「推定」します。

- `go test -coverprofile` の結果
- ソースコードの構文解析 (AST)

ここで大事なのは、`gobce` の C1 は **推定値 (estimated)** だということです。

## そもそも coverprofile とは？

`go test -coverprofile coverage.out` を実行すると、
「どのコード範囲が何回実行されたか」が記録されたテキストファイルができます。

ただし、これは主に statement/block ベースの情報です。
「if の true 側と false 側が両方通ったか」を直接そのまま持っているわけではありません。

## なぜ AST 解析が必要？

coverprofile だけだと「分岐そのもの」が見えにくいので、`gobce` はコードを読みます。

具体的には、Go の標準機能でコードを AST として解析して、
以下のような分岐候補を探します。

- `if` の true 側 / false 側
- `switch` / `type switch` の各 case
- `for` / `range` の「ループ本体に入ったか・入らなかったか」

つまり、

1. コード上で「分岐の場所」を見つける
2. coverprofile で「その範囲が実行されたか」を当てる

という 2 段構えです。

## 実際の処理フロー

`gobce analyze --coverprofile coverage.out --format json` で、概ね次を行います。

1. `coverage.out` を読む
2. 各行をパースして、ファイル・開始行・終了行・実行回数を取り出す
3. 必要なら import path 形式のファイル名をローカル実ファイルパスに解決する
4. 各 Go ファイルを AST 解析して分岐候補を集める
5. 分岐候補の行範囲と coverage の行範囲が重なっているかを見て、covered/uncovered を判定する
6. 結果を集計して JSON 出力する
   - `statementCoverage`
   - `estimatedBranchCoverage`
   - `uncoveredBranches`

## 図解: 1つの if をどう判定するか

ここでは、最小例で `gobce` の見え方を対応させます。

対象コード:

```go
func score(v int) int {
	if v > 10 {
		return 1
	} else {
		return 2
	}
}
```

テストが `v=20` しか通していない場合のイメージ:

```text
if 条件: v > 10
├─ true 側: 実行された
└─ false 側: 実行されていない
```

coverprofile 側では、ざっくり次の情報が入ります。

```text
sample.go:4.13,5.3 ... count=1   (true 側の範囲)
sample.go:6.8,8.3  ... count=0   (false 側の範囲)
```

`gobce` は AST で

- `if_true_path` の行範囲
- `if_false_path` の行範囲

を作り、coverprofile の範囲と重ねて判定します。

結果イメージ:

```json
{
  "uncoveredBranches": [
    {
      "file": "sample.go",
      "line": 6,
      "kind": "if_false_path"
    }
  ]
}
```

つまり「`if` がある」だけではなく、
**`if` の両側のコードがテスト実行時に実際に実行されたか**を `gobce` は見ようとしています。

## 出力項目の見方

- `statementCoverage`  
  `go test` のカバレッジ情報から算出した statement ベースの割合。

- `estimatedBranchCoverage`  
  `gobce` が分岐候補に対して covered/uncovered を推定して計算した割合。

- `uncoveredBranches`  
  テスト実行時にまだ実行されていないと推定された分岐の一覧。
  どのファイル・どの行・どの kind かが入ります。

## 「推定」になる理由

Go の coverprofile は branch coverage 専用のフォーマットではないため、
分岐の完全な実行情報を直接は持っていません。

そのため `gobce` は、
「コード構造 + 実行された行範囲」から妥当な推定をしている、という位置づけです。

## よくある疑問

### Q. statementCoverage が高いのに estimatedBranchCoverage が低いのはなぜ？

A. 直線的な処理のコードはテスト実行時に実行されていても、条件分岐の片側だけしかテスト実行時に実行されていない可能性があります。
`gobce` はその差分を見える化するためのツールです。

### Q. まず何を見ればいい？

A. `uncoveredBranches` から見れば OK です。  
分岐の未到達ポイントが具体的に出るので、テスト追加の優先順位を付けやすくなります。

## 開発者向けの読み進め順

実装を読むなら、次の順がおすすめです。

1. `analyze.go` (公開 API の入口)
2. `internal/analyzer/analyze.go` (処理の本流)
3. `internal/analyzer/coverprofile.go` (coverage の入力)
4. `internal/analyzer/branches.go` (分岐候補の抽出)
5. `internal/analyzer/coverage.go` (集計)

この順で追うと、全体像を掴みやすいです。
