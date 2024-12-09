
# VOICEVOX Engine CLI tool

VOICEVOX Engineまたはその互換APIサーバへのリクエストをコマンドラインから実行します。  
音声の再生、保存が可能です。

- [voicevox_engine](https://github.com/VOICEVOX/voicevox_engine)
- [COEIROINK(v1)](https://coeiroink.com/)
- [AivisSpeech-Engine](https://github.com/Aivis-Project/AivisSpeech-Engine)

## インストール

```bash
$ git clone https://github.com/Dotinkasra/vsay
$ cd vsay/cmd/vasy
$ go build 
```

## 使用方法
デフォルトのポート番号はAivisSpeechのものになります（10101）。  
あらかじめVOICEVOX、COEIROINK、AivisSpeechのいずれかのエンジンを起動しておく必要があります。

```bash
$ vsay -h
NAME:
   vsay - Synthesized voice is played from the terminal.

USAGE:
   vsay [global options] command [command options]

COMMANDS:
   say      Say something
   dict     Show dictionary
   install  Show version
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### 音声を合成する

```bash
$ vsay say -h 
NAME:
   vsay say - Say something

USAGE:
   vsay say command [command options]

COMMANDS:
   ls, l    Show speakers
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --host value             Host address (default: "http://localhost")
   --port value, -p value   Port number (default: 10101)
   --id ID                  Style ID. This takes priority over the speaker number option. (default: 0)
   --number ls, -n ls       The speaker number as displayed by the ls command. (default: 0)
   --style ls               The style number as displayed by the ls command. (default: 0)
   --accent index           Specify the accent by its index in the string. (default: -1)
   --speed 0.5 to 2.0       Set the speaking speed. Valid range: 0.5 to 2.0. (default: 1)
   --intonation 0.0 to 2.0  Set the intonation, affecting the style's strength. Valid range: 0.0 to 2.0. (default: 1)
   --tempo 0.0 to 2.0       Set the tempo. Valid range: 0.0 to 2.0. (default: 1)
   --pitch -0.15 to 0.15    Set the pitch. Valid range: -0.15 to 0.15. (default: 0)
   --save PATH, -s PATH     Specify the PATH to save the audio file.
   --quiet, -q              Don't play audio. (default: false)
   --b64, -b                Outputs audio as base64 encoding to Stdout. (default: false)
   --help, -h               show help
```
---
#### 話者の表示 
ここで表示される番号は`vsay say`コマンドのIDとして使用できます。
例えばVOICEVOXのポートを指定する場合は以下のようになります。

```bash
$ vsay say ls
0: Anneli
        0: 888753760: ノーマル
        1: 888753761: 通常
        2: 888753762: テンション高め
        3: 888753763: 落ち着き
        4: 888753764: 上機嫌
        5: 888753765: 怒り・悲しみ
1: peach
        0: 933744512: ノーマル
```

以下は`VOICEVOX`を起動し、ポート番号に`50021`を指定した場合です。
```bash
$ vsay say ls -p 50021
0: 四国めたん
        0: 2: ノーマル
        1: 0: あまあま
        2: 6: ツンツン
        3: 4: セクシー
        4: 36: ささやき
        5: 37: ヒソヒソ
1: ずんだもん
        0: 3: ノーマル
        1: 1: あまあま
        2: 7: ツンツン
        3: 5: セクシー
        4: 22: ささやき
        5: 38: ヒソヒソ
        6: 75: ヘロヘロ
        7: 76: なみだめ
2: 春日部つむぎ
        0: 8: ノーマル
3: 雨晴はう
        0: 10: ノーマル
4: 波音リツ
        0: 9: ノーマル
        1: 65: クイーン
5: 玄野武宏
        0: 11: ノーマル
        1: 39: 喜び
        2: 40: ツンギレ
        3: 41: 悲しみ
```

#### 音声の作成
`id`オプションには固有のスタイルIDを指定できます。  
`number`オプションと`style`オプションは、`ls`サブコマンドの要素番号です。  
<br>
**`id`オプションは常に`number`オプションと`style`オプションより優先されます。**
```bash
$ vsay say -id 888753765 -intonation 1.0 -accent 4 -s ./test.wav "こんにちは"
$ vsay say -n 0 -style 5 "こんにちは"
```

パイプを使って標準入力に話したいテキストを渡すこともできます。 

```bash
$ echo "こんにちは" | vsay say -id 888753765 -intonation 1.0 -accent 4 
$ cat ./text.txt | vsay say -id 888753765 -intonation 1.0 -accent 4 
```

合成した音声を`Base64`として標準出力することも可能です。
```bash
$ echo "こんにちは" | vsay say -id 888753765 -intonation 1.0 -accent 4 -b -q >> b64out.txt
```

### ユーザー辞書機能
```bash
$ vsay dict -h
NAME:
   vsay dict - Show dictionary

USAGE:
   vsay dict command [command options]

COMMANDS:
   add, a     Add word
   delete, r  Remove word
   ls, l      Show dictionary
   help, h    Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

#### 現在の辞書一覧を表示する
```bash
$ vsay dict ls
dc94a187-9881-43c9-a9c1-cebbf774a96d:
        ID: 1348
        単語: 担々麺
        読み: タンタンメン
        アクセント: 3
```

#### 辞書へ追加する
```bash
$ vsay dict add -h
NAME:
   vsay dict add - Add word

USAGE:
   vsay dict add [command options]

OPTIONS:
   --host value                           Host address (default: "http://localhost")
   --port value, -p value                 Port number (default: 10101)
   --surface word, -w word                The surface form of the word
   --pronunciation katakana, -y katakana  Pronunciation of words (katakana)
   --accent value, -a value               Accented type (refers to where the sound goes down) (default: 0)
   --type value, -t value                 One of the following: PROPER_NOUN, COMMON_NOUN, VERB, ADJECTIVE,SUFFIX
   --priority 0 to 10                     Word priority (integer from 0 to 10) (default: 0)
   --help, -h                             show help
```
以下は`星野瑠美衣`を登録する例です。
```bash
$ vsay dict add -w "星野瑠美衣" -y "ホシノルビー" -a 4 -t proper_noun 
Success
"c81c2a68-xxxx-xxxx-xxxx-xxxxxxxxx"
```

#### 辞書から削除する
```bash
$ vsay dict delete -h
NAME:
   vsay dict delete - Remove word

USAGE:
   vsay dict delete [command options]

OPTIONS:
   --host value            Host address (default: "http://localhost")
   --port value, -p value  Port number (default: 10101)
   --help, -h              show help
```
先ほど登録した`星野瑠美衣`を削除する場合。
```bash
$ vsay dict delete c81c2a68-xxxx-xxxx-xxxx-xxxxxxxxx
Success
```

#### （試験的）モデルデータをインストールする
**AivisSpeechにしか対応していません**
```bash
$ vsay install -h
NAME:
   vsay install - Show version

USAGE:
   vsay install [command options]

OPTIONS:
   --host value            Host address (default: "http://localhost")
   --port value, -p value  Port number (default: 10101)
   --path path, -i path    If installing from a local file, specify the file path.
   --link URL, -l URL      Specify if you wanna install from a URL
   --help, -h              show help
```
```bash
$ vsay install -i /Users/name/Downloads/model.aivmx
```

# 未実装の機能
 - アクセントをフレーズ毎に指定する
 - 設定機能
 - プリセット機能
 - クエリ編集機能

# License
MIT