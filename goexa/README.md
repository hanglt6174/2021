# 【gin + gorm】sample

```
go run main.go
```

## 注意事項
動かす前に、.envファイルをルートディレクトリに作って以下のように記述してください。`""`はいらないです。

```
mytweet_DBMS=mysql
mytweet_USER=<ユーザー名>
mytweet_PASS=<パスワード>
mytweet_DBNAME=<DB名>
```
Redis for window:
https://redislabs.com/ebook/appendix-a/a-3-installing-on-windows/a-3-2-installing-redis-on-window/