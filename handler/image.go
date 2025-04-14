package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// 画像ディレクトリのパス
var imgDirPath string

// 指定された名前の画像をDLする。
func DownloadImage(w http.ResponseWriter, r *http.Request) {
	if imgDirPath == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// URLパスから画像の名前を取得
	name := r.PathValue("name")
	if name == "" || strings.Contains(name, ".") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Image name required in last path"))
		return
	}

	// 画像ファイルを開く
	imgPath := path.Join(imgDirPath, name)
	file, err := os.Open(fmt.Sprintf("%s.jpg", imgPath))
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// 画像の最終更新日時を取得。取得できない場合は現在日時を使用
	var modTime time.Time
	imgInfo, err := file.Stat()
	if err == nil {
		modTime = imgInfo.ModTime()
	} else {
		modTime = time.Now()
	}

	// 画像をレスポンスに書き込む
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.jpg", name))
	http.ServeContent(w, r, fmt.Sprintf("%s.jpg", name), modTime, file)
}

func init() {
	imgDirPath = os.Getenv("IMG_DIR_PATH")
	if imgDirPath == "" {
		slog.Error("IMG_DIR_PATH is not set", slog.String("err", "IMG_DIR_PATH is not set"))
	}
}
