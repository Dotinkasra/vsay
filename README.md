
# VOICEVOX Engine CLI tool

VOICEVOX Engineまたはその互換APIサーバをコマンドラインから実行します。

## Installation

```bash
$ git clone https://github.com/Dotinkasra/vsay
$ cd vsay/cmd/vasy
$ go build 
```

## Usage
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

IDか上記のlsオプションで表示された番号を指定できます。
```bash
$ say -id 888753765 -intonation 1.0 -accent 4 -s ./test.wav "こんにちは"
$ say --number 0 --style 5 "こんにちは"
```
## License
MIT