# Go C1 カバレッジ計測ツール gobce

Status: 初期コンセプトノート
Date: 2026-04-25

## 位置づけ

`gobce` は Go 向けの C1 branch coverage estimator であり、独立した OSS library / CLI として提供する。

`gobce` は Golang Branch Coverage Estimator の略である。Go 標準の coverage 情報をもとに、C1 branch coverage を軽量に推定することを目的にする。

重要な方針:

```text
gobce は独立した別リポジトリとして開発する。
任意の CI / quality gate から利用できる汎用ツールとして設計する。
```

## 目的

Go 標準の coverage は `go test -coverprofile` による statement/block coverage を提供する。しかし、Production Readiness の観点では「分岐がどれだけテストされているか」というより強いシグナルが欲しい。

`gobce` は Go の coverprofile と AST / 軽量 CFG 解析を組み合わせ、C1 branch coverage を推定する。

```text
go test -coverprofile + Go AST/CFG analysis
  -> estimated branch coverage
  -> uncovered branch findings
```

## Standalone CLI

```bash
go test ./... -coverprofile coverage.out

go run ./cmd/gobce analyze \
  --coverprofile coverage.out \
  --format json \
  --output gobce-result.json
```

出力例:

```json
{
  "language": "go",
  "statementCoverage": 82.1,
  "estimatedBranchCoverage": 68.4,
  "uncoveredBranches": [
    {
      "file": "internal/order/validator.go",
      "line": 42,
      "kind": "if_false_path"
    }
  ]
}
```

## 初期解析スコープ

MVP で扱う branch candidate:

```text
if / else
switch / case / default
type switch
select
for / range の loop body entered vs not entered
&& / || による short-circuit boolean expression
```

将来的に扱う branch candidate:

```text
panic / recover paths
defer-related paths
error-return conventions
table-test branch attribution
generated-code filtering improvements
```

## アルゴリズム概要

```text
1. Go coverprofile を parse する。
2. go/parser と go/ast で Go package / file を parse する。
3. 対応構文から branch candidate span を作る。
4. coverage block と branch candidate span を対応付ける。
5. 各 branch side が実行されたかを推定する。
6. estimated C1 percentage を計算する。
7. file, line, kind, recommendation 付きの uncovered branch finding を出す。
```

## 重要な制約

最初のバージョンでは、結果を exact C1 ではなく estimated C1 と呼ぶ。Go の coverprofile は branch coverage 専用フォーマットではないため、シグナルの性質を正直に表現する。

```text
statement coverage = measured
estimated C1       = inferred from coverprofile + source analysis
```
