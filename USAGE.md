# SD Auto Generation Tool - 使い方ガイド

## 概要

このツールは、Stable Diffusion WebUIのAPIを利用して、キャラクターとシチュエーションの組み合わせで自動的に画像を生成するGoプログラムです。コマンドラインインターフェースで操作し、生成条件を設定してバッチ生成が可能です。

## 必要条件

- **Go 1.23.0以上**
- **Stable Diffusion WebUI**
  - APIモードで起動可能
  - デフォルトURL: `http://127.0.0.1:7860`
- **Windows/Linux/macOS** (Goが対応するOS)

## インストール

1. **リポジトリのクローンまたはダウンロード**
   ```
   # このフォルダをそのまま使用
   cd C:\Users\(省略)\stible_diffsion_new_go
   ```

2. **依存関係のインストール**
   ```bash
   go mod tidy
   ```

3. **ビルド**
   ```bash
   go build -o sd-auto-new.exe
   ```

## セットアップ

### 1. Stable Diffusion WebUIの準備

1. Stable Diffusion WebUIをインストール
2. `--api` オプション付きで起動
   ```bash
   # WebUIのフォルダで
   webui-user.bat --api
   # または
   python launch.py --api
   ```
3. WebUIが `http://127.0.0.1:7860` で起動していることを確認

### 2. データファイルの準備

#### キャラクター定義 (`src/character.txt`)
JSON配列形式でキャラクターを定義：
```json
[
  { "en": "fischl (genshin_impact)", "ja": "フィッシュル" },
  { "en": "<lora:curearcanashadow_v1.0_IL:1>Ruruka Moria,blonde hair,red eyes,purple eyes,ahoge,multicolored hair,black dress,black footwear,black thighhighs,hair ornament,artist:rella,age 13,", "ja": "キュアアルカナ・シャドウ" }
]
```

#### シチュエーション定義 (`situation/` フォルダ)
各カテゴリフォルダ内にシチュエーションファイルを作成：
```
## 重要！situationのtxtファイルは先頭に{'数字'_'シチュエーション名'.txt}になるようにお願いします
situation/
├── school/(個々のフォルダ名はなんでもよい)
│   ├── 01_classroom.txt
│   └── 02_playground.txt
└── home/
    ├── 01_livingroom.txt
    └── 02_bedroom.txt
```

各シチュエーションファイルの内容例 (`01_classroom.txt`)：
```json
[{"en": "school uniform, classroom, daytime", "ja": "制服・教室・昼間"}]
```

### 3. 設定ファイル (`.env`)
生成パラメータを調整：
```env
NEGATIVE_PROMPT=score_6, score_5, score_4, source_anime, source_cartoon, watermark, text, signature, blurry, lowres, bad anatomy, bad hands, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, artist name
STEPS=20
CFG_SCALE=7
WIDTH=768
HEIGHT=1024
SAMPLER_NAME=DPM++ 2M Karras
SEED=-1
```

## 使い方

### プログラムの起動

```bash
# ビルド済みの場合
./sd-auto-new.exe

# 直接実行の場合
go run .
```

### 操作フロー

1. **SD接続確認**
   - プログラム起動時に自動でSD WebUI APIに接続確認
   - 接続できない場合はエラーメッセージが表示され終了

2. **モード選択**
   ```
   ━━━━━━━━━━━━━━━━━━━━━━━━━━
   🎨 SD Auto Generation Tool
   ━━━━━━━━━━━━━━━━━━━━━━━━━━
   1. 生成モード
   2. 設定モード
   0. 終了
   ━━━━━━━━━━━━━━━━━━━━━━━━━━
   選択 >>
   ```

3. **生成モードの場合**
   - **キャラクター選択**: 番号入力で選択（-1でランダム）
   - **シチュエーションカテゴリ選択**: フォルダ一覧から選択
   - **繰り返し回数設定**: 各シチュエーションの回数 + カテゴリ全体の回数
   - **条件確認**: 生成枚数などを確認
   - **生成実行**: 進捗表示しながら画像生成

4. **設定モードの場合**
   - 現在のSDパラメータを表示
   - 変更したい項目を選択して新しい値を入力
   - `.env`ファイルに保存

### サンプル実行例

```
━━━━━━━━━━━━━━━━━━━━━━━━━━
🎨 SD Auto Generation Tool
━━━━━━━━━━━━━━━━━━━━━━━━━━
1. 生成モード
2. 設定モード
0. 終了
━━━━━━━━━━━━━━━━━━━━━━━━━━
選択 >> 1

👤 キャラクター選択

📋 キャラクター一覧:
━━━━━━━━━━━━━━━━━━━━━━━━━━
  0. フィッシュル
━━━━━━━━━━━━━━━━━━━━━━━━━━

キャラクター番号を選択 (-1でランダム) >> 0
✅ 選択: フィッシュル

📁 シチュエーションカテゴリ選択

📋 カテゴリ一覧:
━━━━━━━━━━━━━━━━━━━━━━━━━━
  0. category1 (2個のシチュエーション)
━━━━━━━━━━━━━━━━━━━━━━━━━━

カテゴリ番号を選択 >> 0
✅ 選択: category1

各シチュエーションの繰り返し回数を設定:
  制服・教室・昼間 の回数 >> 2
  公園・ベンチ・夕方 の回数 >> 1

このカテゴリ全体を何回繰り返しますか？ >> 1

📋 生成条件確認
━━━━━━━━━━━━━━━━━━━━━━━━━━
キャラクター: フィッシュル
カテゴリ: category1
カテゴリ繰り返し: 1回
シチュエーション詳細:
  制服・教室・昼間: 2回 × 1回 = 2枚
  公園・ベンチ・夕方: 1回 × 1回 = 1枚
合計生成枚数: 3枚
━━━━━━━━━━━━━━━━━━━━━━━━━━

この条件で生成しますか？ (y/n) >> y

📁 出力先: output/2026-04-09
🎨 合計 3 枚の画像を生成します

🔄 カテゴリ 1/1 回目
📸 生成中 1/3 - キャラ:フィッシュル, シチュ:制服・教室・昼間
📸 生成中 2/3 - キャラ:フィッシュル, シチュ:制服・教室・昼間
📸 生成中 3/3 - キャラ:フィッシュル, シチュ:公園・ベンチ・夕方

✅ すべての生成が完了しました！
```

## 出力ファイル

生成された画像は `output/YYYY-MM-DD/` フォルダに保存されます。

- ファイル名形式: `{キャラ名(ja)}_{シチュ名}_{連番}.png`
- 例: `フィッシュル_制服・教室・昼間_001.png`

## 中断とキャンセル

- **Ctrl+C** で生成を中断可能
- 中断時は現在の生成をキャンセルし、SD APIにも停止リクエストを送信

## トラブルシューティング

### SD WebUIに接続できない
- SD WebUIが `--api` オプション付きで起動しているか確認
- URLが `http://127.0.0.1:7860` であるか確認
- ファイアウォールやポート競合がないかチェック

### 画像が生成されない
- SD WebUIのモデルが読み込まれているか確認
- `.env` のパラメータが適切か確認
- プロンプトの内容を確認

### ビルドエラー
- Goのバージョンが1.23.0以上か確認
- `go mod tidy` を実行して依存関係を更新

### データファイルのエラー
- JSON形式が正しいか確認（JSONLintなどで検証）
- ファイルパスが正しいか確認
- 文字コードがUTF-8か確認

## 拡張・カスタマイズ

### 新しいキャラクターの追加
`src/character.txt` にJSONオブジェクトを追加：
```json
{ "en": "new_character_prompt", "ja": "新しいキャラクター" }
```

### 新しいシチュエーションの追加
1. `situation/` 内に新しいフォルダ作成
2. フォルダ内に `01_situation.txt` などのファイル作成
3. JSON形式でプロンプト定義

### パラメータの調整
`.env` ファイルを編集するか、設定モードで変更

## サポート

バグ報告や機能リクエストは、プログラムのコメントやrequirements.mdを参照してください。