# SD Auto Generation Tool - 使い方ガイド

Stable Diffusion WebUIのAPIを利用して、キャラクターとシチュエーションの組み合わせで自動的に画像を生成するツールです。

## セットアップ

1. **Go 1.23.0以上** をインストール
2. **Stable Diffusion WebUI** をインストールし、`--api` オプション付きで起動
3. リポジトリをクローン: `git clone https://github.com/YuuuDays/stable_diffusion_auto_create_go.git`
4. 依存関係インストール: `go mod tidy`
5. ビルド: `go build -o sd-auto-new.exe`

## 使い方

1. `./sd-auto-new.exe` を実行
2. SD WebUI API接続を確認
3. メニューから「1. 生成モード」を選択
4. キャラクターを選択（番号入力、-1でランダム）
5. シチュエーションカテゴリを選択
6. 各シチュエーションの繰り返し回数を設定
7. カテゴリ全体の繰り返し回数を設定
8. 生成条件を確認して実行
9. Ctrl+Cで中断可能

## データファイル

### ネガティブプロンプト
`config/sd_config.txt` に設定。デフォルトでscore_6などの品質向上プロンプトが入っています。

### シチュエーション
`situations/` フォルダにカテゴリごとのフォルダを作成し、その中にシチュエーションファイルを置く。

#### ファイル命名規則
- ファイル名: `{数字}_{シチュエーション名}.txt` （例: `01_classroom.txt`）
- 内容: JSON配列 `[{"en": "プロンプト", "ja": "日本語名"}]`

### キャラクター
`src/character.txt` にJSON配列で定義。

## 注意点

- `.env` はAPI URLなどの環境設定
- `config/sd_config.txt` はSDパラメータ（Git追跡対象）
- 生成画像は `output/YYYY-MM-DD/` に保存
- エラー時はSD WebUIが起動しているか確認