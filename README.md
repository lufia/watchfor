watchfor
========

INSTALLATION
------------

	go get github.com/lufia/watchfor

SYNOPSIS
--------

	watchfor [-s .ext] [-t .ext] [-c command] [dir]

OPTIONS
-------

|option|説明|デフォルト|
|------|----|----------|
|-a addr|HTTPを待ち受けるアドレス|:8080|
|-s .ext|監視するソースファイル拡張子|(未設定)|
|-t .ext|監視するターゲットファイル拡張子|(未設定)|
|-c command|ソースファイルに変更があったとき実行するコマンド|(未設定)|

DESCRIPTION
-----------

最新のファイルをブラウザ上に表示します。
監視しているファイルに変更があればすぐに、ブラウザ上の表示も更新します。

一般的な使い方としては、ソースファイルを手元で編集します。
その間、ブラウザでターゲットファイルを表示しておきます。
通常であれば http://localhost:8080/<ターゲットファイル名> にアクセスします。
ソースファイルに変更があれば、-cオプションで指定したコマンドを実行させます。
ブラウザで表示しているターゲットファイルは、
コマンドが終わったときに自動リロードされますので、
エディタから手を離さずに最新のプレビューを確認する、という形になるでしょう。

cオプションのコマンドからは、特別に以下の環境変数を参照可能です。

|変数名|意味|
|------|----|
|source|変更のあったソースファイル名|
|target|更新するべきターゲットファイル名|

EXAMPLE
-------

カレントディレクトリの.pumlファイルを監視して、
変更があればPlantUMLで画像生成。
ブラウザにはtest.pngの最新版を表示する。

	watchfor -s .puml -t .png -c 'plantuml $source' .
	http://localhost:8080/test.png

portディレクトリを監視、通常のテキストファイルをブラウザに表示する。

	watchfor -s .md -t .md port
	http://localhost:8080/port/README.md

サブディレクトリは監視対象に含みません。

BUGS
----

* ソースとターゲットファイルは拡張子を除いて同じ名前である必要があります。
	* ターゲット生成に時間がかかる場合、コマンドの終わりを待つためです。
* fsnotifyパッケージがまだ正式版ではないため動作しない環境があります。
* 古いブラウザ等では動作しない可能性が非常に高いです。
