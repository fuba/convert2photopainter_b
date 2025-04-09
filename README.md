# PhotoPainter (B) 画像変換ツール

PhotoPainter (B) 用の画像を変換するシンプルなコマンドラインツールです。一般的な画像形式から、PhotoPainter (B) で表示可能な6色 (黒, 白, 赤, 緑, 青, 黄) の24-bit BMP形式に変換します。

## 機能

- 一般的な画像形式（JPG, PNG, GIF, BMP等）をサポート
- Floyd-Steinbergディザリングアルゴリズムによる色変換
- 800×480の解像度に自動調整
- 縦長画像の自動回転（90度回転して適切な向きに調整）
- 単一ファイルまたはディレクトリ内の複数ファイルの一括処理

## 使用方法

```
PhotoPainter (B) 画像変換ツール
使用方法: photoconvert [オプション] 入力ファイル/ディレクトリ

オプション:
  -o <dir>       出力先ディレクトリ（指定しない場合は入力と同じディレクトリ）
  -batch         バッチモード（ディレクトリ内のすべての画像を処理）
  -r <res>       解像度（800x480 または 480x800、デフォルト: 800x480）
  -v             詳細なログ出力
  -h             このヘルプメッセージを表示

例:
  単一ファイルの変換:
    photoconvert input.jpg
  出力先ディレクトリ指定:
    photoconvert -o /path/to/output input.jpg
  バッチ処理:
    photoconvert -batch /path/to/images/
  すべてのオプション:
    photoconvert -batch -o /path/to/output -r 800x480 -v /path/to/images/
```

## インストール

リポジトリをクローンし、ビルドします：

```bash
git clone https://github.com/example/photopainter.git
cd photopainter
go build -o photoconvert ./cmd/photoconvert/
```

## プロジェクト構造

```
/convert2photopainter_b/
  ├── cmd/
  │   └── photoconvert/
  │       └── main.go            # コマンドラインインターフェース
  ├── internal/
  │   ├── convert/
  │   │   └── convert.go         # 画像変換処理
  │   ├── dither/
  │   │   └── dither.go          # Floyd-Steinbergディザリング
  │   └── resize/
  │       └── resize.go          # 画像リサイズ処理
  ├── go.mod                     # 依存関係管理
  └── README.md                  # プロジェクト説明
```

## 処理の流れ

1. 画像を読み込み
2. 指定した解像度（デフォルト: 800×480）にリサイズ
3. Floyd-Steinbergディザリングアルゴリズムを適用して6色に減色
4. 24-bit BMP形式で保存

## 注意事項

- 出力画像の色数は PhotoPainter (B) の仕様に合わせて6色（黒, 白, 赤, 緑, 青, 黄）に制限されています
- 色の変換にはFloyd-Steinbergディザリングを使用して、より自然な見た目を実現しています
- 画像は自動的にリサイズされますが、元のアスペクト比と異なる場合はクロップされる場合があります
- MACでTFカードを作成する場合、隠しファイルが生成される場合があるため、変換後にTFカードのルートディレクトリからfileList.txtとindex.txtを削除するとともに、picフォルダ内の隠しファイルも削除することをおすすめします
