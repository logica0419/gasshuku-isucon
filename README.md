# gasshuku-isucon

> 「ISUCON」は、LINE株式会社の商標または登録商標です。  
> <https://isucon.net/>

traP 2023春合宿記念 オリジナルISUCON

## 🙏 Waiting For Contribution 🙏

### Frontend

- フロントエンドを作って下さる方を募集しています

### 別言語実装

- `webapp/{任意の言語}`に、Go以外の実装を作って下さる方を募集しています

### ベンチマーカー

- バランス調整がまだまともでない部分が多数あるので、IssueでのアドバイスやPRお待ちしております

### Provisioning

- Provisioningのためのスクリプト類は絶賛準備中です
- もし協力して下さる方がいたら大変ありがたいです

## 遊び方 (08/26時点)

1. 手元でもどこかのサーバーでも良いので、GoとMySQLがインストールされたLinux環境を用意する
2. MySQLにパスワード「isucon」で「isucon」という名前のユーザーを作り、「isulibrary」という名前でデータベースを用意する
3. リポジトリを任意の場所にクローンする
4. <https://github.com/logica0419/gasshuku-isucon/releases/tag/initial-data> から
   1. `1_data.sql`をダウンロードし`webapp/sql`内に置く
   2. `init_data.json`をダウンロードし`bench/repository`内に置く
5. `make run-go`を実行するなどの方法で、`webapp/go/main.go`を恒常的に立ち上げた状態にする
6. `make run-bench`を実行するなどの方法で、`bench/main.go`を実行してベンチマークを開始する
