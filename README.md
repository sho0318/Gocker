# Gocker - Minimal Container Runtime in Go

## 実装されている機能

### 1. リソースの隔離 (Namespaces)
`syscall.SysProcAttr` を利用し、以下のリソースをホストから隔離します。
* **UTS**: 独自のホスト名（`gocker-container`）を設定。
* **PID**: コンテナ内のプロセスIDを1から開始させ、ホストのプロセスから隠蔽。
* **NS (Mount)**: 独自のマウントテーブルを持ち、ホストのマウントポイントに影響を与えません。

### 2. リソースの制限 (Cgroups v2)
`/sys/fs/cgroup` を直接操作し、以下の制限を課しています。
* **PID制限**: 最大20プロセスまで（Fork爆弾などのリソース枯渇対策）。
* **メモリ制限**: 最大10MBまで。
* **スワップ制限**: スワップを禁止（`0`）。

### 3. ファイルシステムの制御
* **chroot**: 指定したディレクトリ（`./rootfs`）をルートディレクトリとして認識。
* **procfs**: コンテナ内で `ps` コマンド等が正しく動作するよう、専用の `/proc` をマウント。

## 事前準備

1.  **OS**: Linux (Cgroup v2 が有効なカーネル)。
2.  **rootfs**: コンテナのルートとなるファイルシステム。
    * プロジェクト直下に `rootfs` ディレクトリを作成し、最小限のLinux環境（Alpine Linuxのrootfsなど）を配置してください。

```bash
mkdir rootfs
cd rootfs
curl -L -o alpine.tar.gz https://dl-cdn.alpinelinux.org/alpine/v3.18/releases/x86_64/alpine-minirootfs-3.18.4-x86_64.tar.gz
tar -xzf alpine.tar.gz
rm alpine.tar.gz
```

## 使い方

1. ビルド

プログラムをコンパイルします

```bash
go build -o gocker main.go
```

2. 実行

Namespaceの作成やCgroupの操作にはroot権限が必要です

```bash
sudo ./gocker run /bin/sh 
```