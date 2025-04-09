package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/example/convert2photopainter/internal/convert"
)

func main() {
	// コマンドライン引数の解析
	outputDir := flag.String("o", "", "出力先ディレクトリ（バッチモードでは必須）")
	batchMode := flag.Bool("batch", false, "バッチモード（ディレクトリ内のすべての画像を処理）")
	maxDepth := flag.Int("depth", 3, "バッチモード時の最大サブディレクトリ探索深度")
	resolution := flag.String("r", "800x480", "解像度（800x480 または 480x800）")
	autoRotate := flag.Bool("rotate", true, "縦長画像を自動的に90度回転して適切な方向に調整する")
	verbose := flag.Bool("v", false, "詳細なログ出力")
	help := flag.Bool("h", false, "ヘルプの表示")

	flag.Parse()

	// ヘルプの表示
	if *help || flag.NArg() == 0 {
		printHelp()
		return
	}

	// 入力ファイル/ディレクトリのパス
	inputPath := flag.Arg(0)

	// バッチモードで出力先が指定されていない場合はエラー
	if *batchMode && *outputDir == "" {
		fmt.Println("エラー: バッチモードでは -o オプションで出力先ディレクトリの指定が必須です")
		os.Exit(1)
	}

	// 出力ディレクトリの設定
	outDir := *outputDir
	if outDir == "" {
		// 出力先が指定されていない場合は入力と同じディレクトリを使用（バッチモード以外）
		outDir = filepath.Dir(inputPath)
	}

	// 解像度の確認
	width, height := 800, 480
	if *resolution == "480x800" {
		width, height = 480, 800
	} else if *resolution != "800x480" {
		fmt.Println("警告: サポートされていない解像度です。800x480 を使用します。")
	}

	// 変換オプション
	options := convert.Options{
		Width:      width,
		Height:     height,
		AutoRotate: *autoRotate,
		Verbose:    *verbose,
	}

	// バッチモードまたは単一ファイルモード
	if *batchMode {
		// バッチ処理の前にディレクトリ連番をリセット
		processedDirs = make(map[string]string)
		nextDirNum = 1
		
		// ディレクトリ内のすべての画像を処理
		processDirectory(inputPath, outDir, options, *maxDepth, 0)
	} else {
		// 単一ファイルの処理
		err := processSingleFile(inputPath, outDir, options, "")
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			os.Exit(1)
		}
	}
}

// グローバル変数で処理済みのディレクトリを管理
var processedDirs = make(map[string]string)
var nextDirNum = 1

// ディレクトリパスを処理して連番を作成
func processDirPath(inputDir, basePath string) string {
	// 既に処理済みのディレクトリかチェック
	if num, exists := processedDirs[inputDir]; exists {
		return num
	}
	
	// 新しいディレクトリ番号を割り当て
	dirNum := fmt.Sprintf("%04d", nextDirNum)
	processedDirs[inputDir] = dirNum
	nextDirNum++
	
	return dirNum
}

// ディレクトリ内のすべての画像を処理
func processDirectory(inputDir, outputDir string, options convert.Options, maxDepth, currentDepth int) {
	// 最大深度チェック
	if currentDepth > maxDepth {
		return
	}

	// ディレクトリの存在確認
	fileInfo, err := os.Stat(inputDir)
	if err != nil {
		fmt.Printf("エラー: ディレクトリが見つかりません: %v\n", err)
		os.Exit(1)
	}
	if !fileInfo.IsDir() {
		fmt.Println("エラー: 指定されたパスはディレクトリではありません")
		os.Exit(1)
	}

	// 出力ディレクトリが存在しない場合は作成
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("エラー: 出力ディレクトリの作成に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// サポートされている画像形式の拡張子
	supportedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
	}

	// サブディレクトリリストを収集
	var subDirs []string
	dirEntries, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("エラー: ディレクトリの読み取りに失敗しました: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			entryPath := filepath.Join(inputDir, entry.Name())
			subDirs = append(subDirs, entryPath)
		}
	}

	// サブディレクトリを名前でソート
	if len(subDirs) > 0 {
		// サブディレクトリを名前順にソート
		sort.Strings(subDirs)

		// サブディレクトリを再帰的に処理
		for _, subDir := range subDirs {
			// サブディレクトリを処理
			processDirectory(subDir, outputDir, options, maxDepth, currentDepth+1)
		}
	}

	// 現在のディレクトリのファイルリストを収集
	var imageFiles []string
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("エラー: ディレクトリの読み取りに失敗しました: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			entryName := entry.Name()
			ext := strings.ToLower(filepath.Ext(entryName))
			if supportedExt[ext] {
				entryPath := filepath.Join(inputDir, entryName)
				imageFiles = append(imageFiles, entryPath)
			}
		}
	}

	// 画像ファイルを名前順にソート
	sort.Strings(imageFiles)

	// 親ディレクトリの連番を取得（メインの処理ディレクトリを基準とする）
	// main関数内で設定されたパスを基準点として使用
	dirPrefix := processDirPath(inputDir, flag.Arg(0))

	// ファイルを連番で処理
	processedCount := 0
	errorCount := 0

	for i, imagePath := range imageFiles {
		// 連番ファイル名を生成 (0001, 0002, ...)
		fileNum := fmt.Sprintf("%s_%04d", dirPrefix, i+1)
		
		// ファイルを処理 - すべて同じ出力ディレクトリに保存
		err := processSingleFile(imagePath, outputDir, options, fileNum)
		if err != nil {
			fmt.Printf("エラー: %s の処理に失敗しました: %v\n", filepath.Base(imagePath), err)
			errorCount++
		} else {
			processedCount++
		}
	}

	if options.Verbose && (processedCount > 0 || errorCount > 0) {
		fmt.Printf("%s: %d 個のファイルを処理しました（エラー: %d）\n", inputDir, processedCount, errorCount)
	}
}

// 単一ファイルを処理
func processSingleFile(inputPath, outputDir string, options convert.Options, fileNum string) error {
	if options.Verbose {
		fmt.Printf("処理中: %s\n", inputPath)
	}

	// ファイルが存在するか確認
	_, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("ファイルが見つかりません: %v", err)
	}

	// 出力ディレクトリが存在しない場合は作成
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	// fileNumを使用して連番のファイル名を生成
	outputName := fileNum

	// 出力ファイルパスを生成（24bit BMPファイル）
	outputPath := filepath.Join(outputDir, outputName+".bmp")

	// 変換処理を実行
	err = convert.ConvertImage(inputPath, outputPath, options)
	if err != nil {
		return fmt.Errorf("画像変換に失敗しました: %v", err)
	}

	if options.Verbose {
		fmt.Printf("保存しました: %s\n", outputPath)
	}

	return nil
}

// ヘルプメッセージの表示
func printHelp() {
	fmt.Println("PhotoPainter (B) 画像変換ツール")
	fmt.Println("使用方法: photoconvert [オプション] 入力ファイル/ディレクトリ")
	fmt.Println("\nオプション:")
	fmt.Println("  -o <dir>       出力先ディレクトリ（バッチモードでは必須）")
	fmt.Println("  -batch         バッチモード（ディレクトリ内のすべての画像を処理）")
	fmt.Println("  -depth <n>     バッチモード時の最大サブディレクトリ探索深度（デフォルト: 3）")
	fmt.Println("  -r <res>       解像度（800x480 または 480x800、デフォルト: 800x480）")
	fmt.Println("  -rotate=false  縦長画像の自動回転を無効化（デフォルト: 有効）")
	fmt.Println("  -v             詳細なログ出力")
	fmt.Println("  -h             このヘルプメッセージを表示")
	fmt.Println("\n例:")
	fmt.Println("  単一ファイルの変換:")
	fmt.Println("    photoconvert input.jpg")
	fmt.Println("  出力先ディレクトリ指定:")
	fmt.Println("    photoconvert -o /path/to/output input.jpg")
	fmt.Println("  バッチ処理:")
	fmt.Println("    photoconvert -batch -o /path/to/output /path/to/images/")
	fmt.Println("  サブディレクトリも含めたバッチ処理:")
	fmt.Println("    photoconvert -batch -depth 3 -o /path/to/output /path/to/images/")
	fmt.Println("  すべてのオプション:")
	fmt.Println("    photoconvert -batch -o /path/to/output -depth 2 -r 800x480 -v /path/to/images/")
}