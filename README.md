# Gocker

## 実装されている機能

### 1. リソースの隔離 (Namespaces)
`syscall.SysProcAttr` を利用し、以下のリソースをホストから隔離します。
* **UTS**: 独自のホスト名（`gocker-container`）を設定。
* **PID**: コンテナ内のプロセスIDを1から開始させ、ホストのプロセスから隠蔽。
* **NS (Mount)**: 独自のマウントテーブルを持ち、ホストのマウントポイントに影響を与えません。
* **NET (Network)**: 独自のネットワーク名前空間を作成し、ホストのネットワークから隔離。

### 2. リソースの制限 (Cgroups v2)
`/sys/fs/cgroup` を直接操作し、以下の制限を課しています。
* **PID制限**: 最大20プロセスまで（Fork爆弾などのリソース枯渇対策）。
* **メモリ制限**: 最大10MBまで。
* **スワップ制限**: スワップを禁止（`0`）。

### 3. ファイルシステムの制御
* **OverlayFS**: 読み取り専用のベースイメージ（`./rootfs`）の上に書き込み可能なレイヤーを重ねて使用。
  * コンテナごとに独立した書き込み可能なファイルシステムを提供
  * ベースイメージは変更されず、複数のコンテナで安全に共有可能
  * 書き込みデータは `/tmp/gocker/<container-id>/` 配下に保存
* **chroot**: 指定したディレクトリ（`./rootfs`）をルートディレクトリとして認識。
* **procfs**: コンテナ内で `ps` コマンド等が正しく動作するよう、専用の `/proc` をマウント。

### 4. ネットワークの隔離と接続
* **ブリッジネットワーク**: ホスト上に `br0` という仮想ブリッジを作成し、コンテナ間の通信基盤を構築。
* **veth pair**: 仮想イーサネットペアを作成し、一方をコンテナ内に、もう一方をホストのブリッジに接続。
* **NAT設定**: `iptables` を使用してコンテナから外部ネットワークへのアクセスを実現。
* **IP設定**: コンテナに `172.18.0.1/24` のIPアドレスを割り当て、ブリッジをゲートウェイとして構成。

## 事前準備

1.  **OS**: Linux (Cgroup v2 が有効なカーネル)。
2.  **必要なツール**: 
    * `ip` コマンド（iproute2パッケージ）
    * `iptables` コマンド
    * `nsenter` コマンド（util-linuxパッケージ）
3.  **rootfs**: コンテナのルートとなるファイルシステム。
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
go build -o gocker cmd/gocker/main.go
```

2. 実行

Namespaceの作成やCgroupの操作にはroot権限が必要です

```bash
sudo ./gocker run /bin/sh 
```

3. ネットワーク接続の確認

コンテナ内で以下のコマンドを実行してネットワーク接続を確認できます

```bash
# IPアドレスの確認
ip addr show

# ゲートウェイへのpingテスト
ping -c 3 172.18.0.1

# 外部ネットワークへの接続確認
ping -c 3 8.8.8.8
ping -c 3 google.com
```

## プロジェクト構成

```
gocker/
├── cmd/
│   └── gocker/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── config.go                # 設定管理
│   ├── helpers.go               # ユーティリティ関数
│   ├── container/
│   │   └── runtime.go           # コンテナランタイム
│   ├── filesystem/
│   │   └── overlay.go           # OverlayFS管理
│   ├── host/
│   │   ├── cgroup.go            # Cgroup管理
│   │   └── launcher.go          # コンテナ起動処理
│   └── network/
│       ├── bridge.go            # ブリッジネットワーク管理
│       ├── networkManager.go    # ネットワーク全体管理
│       └── veth.go              # vethペア管理
└── rootfs/                      # コンテナのルートファイルシステム
```
