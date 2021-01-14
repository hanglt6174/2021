# 【gin + gorm】sample

```
go run main.go
```

## 注意事項
動かす前に、

.envファイルをルートディレクトリに作って以下のように記述してください。`""`はいらないです。

Dockerを使用しない場合は以下の環境変数を定義してください。
```
MYSQL_DATABASE=<ユーザー名>
MYSQL_USER=<パスワード>
MYSQL_PASSWORD=<DB名>
```
Redis for window:
https://redislabs.com/ebook/appendix-a/a-3-installing-on-windows/a-3-2-installing-redis-on-window/