
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
   ls, l    Show speakers
   say, s   Say something
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value            Host address (default: "http://localhost")
   --port value, -p value  Port number (default: 10101)
   --help, -h              show help
```

```bash
$ vsay ls
0: Anneli
  0: 888753760: ノーマル
  1: 888753761: 通常
  2: 888753762: テンション高め
  3: 888753763: 落ち着き
  4: 888753764: 上機嫌
  5: 888753765: 怒り・悲しみ
1: white
  0: 706073888: ノーマル
```

例えばVOICEVOXのポートを指定する場合は以下のようになります。
```bash
$ vsay -p 50021 ls
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
```

```bash
$ vsay say -h
NAME:
   vsay say - Say something

USAGE:
   vsay say [command options]

OPTIONS:
   --id ID                  Style ID. This takes priority over the speaker number option. (default: 0)
   --number ls              The speaker number as displayed by the ls command. (default: 0)
   --style ls               The style number as displayed by the ls command. (default: 0)
   --accent index           Specify the accent by its index in the string. (default: -1)
   --speed 0.5 to 2.0       Set the speaking speed. Valid range: 0.5 to 2.0. (default: 1)
   --intonation 0.0 to 2.0  Set the intonation, affecting the style's strength. Valid range: 0.0 to 2.0. (default: 1)
   --tempo 0.0 to 2.0       Set the tempo. Valid range: 0.0 to 2.0. (default: 1)
   --pitch -0.15 to 0.15    Set the pitch. Valid range: -0.15 to 0.15. (default: 0)
   --save PATH, -s PATH     Specify the PATH to save the audio file.
   --quiet, -q              Don't play audio. (default: false)
   --help, -h               show help
```

IDか上記の`vsay ls`オプションで表示された番号を指定できます。 

```bash
$ vsay say -id 888753765 -intonation 1.0 -accent 4 -s ./test.wav "こんにちは"
$ vsay say -n 0 -style 5 "こんにちは"
```

パイプを使って標準入力に話したいテキストを渡すこともできます。 

```bash
$ echo "こんにちは" | vsay say -id 888753765 -intonation 1.0 -accent 4 
$ cat ./text.txt | vsay say -id 888753765 -intonation 1.0 -accent 4 
```

# 未実装の機能
 - ユーザー辞書機能
 - 設定機能
 - プリセット機能
 - クエリ編集機能

## License
MIT