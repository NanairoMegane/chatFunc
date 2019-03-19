package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

/*
HTMLテンプレートをサーブするためのハンドラ
*/
type templateHandler struct {
	once     sync.Once          //HTMLテンプレートを１度だけコンパイルするための指定
	filename string             //テンプレートとしてHTMLファイル名を指定
	tmpl     *template.Template //テンプレート
}

/* templateHandlerをhttp.Handleに適合させるため、ServeHttpを実装する */
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// テンプレートディレクトリを指定する
	path, err := filepath.Abs("./templates/")
	if err != nil {
		panic(err)
	}

	// 指定された名称のテンプレートファイルを一度だけコンパイルする
	t.once.Do(
		func() {
			t.tmpl = template.Must(template.ParseFiles(path + t.filename))
		})

	t.tmpl.Execute(w, nil)
}
