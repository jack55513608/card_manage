# 卡片遊戲寄售平台

## 專案簡介
本專案旨在建立一個以「店家為基礎」的卡片遊戲寄售平台。平台允許玩家將卡片寄放於店家，由店家代為售出，並依照抽成比例進行收益分配。

系統主要提供：
1. 玩家寄售端功能
2. 店家管理與售出端功能
3. 平台管理員後台
4. 金流與抽成模組（初期僅支援現金與店家儲值金）

## 已完成功能列表
- 使用者管理 (註冊、登入、角色權限)
- 店家管理 (新增、修改、抽成比例設定)
- 卡片資料庫管理 (店家可新增/修改卡片標準化資訊)
- 玩家寄售申請與管理
- 交易記錄與清算管理
- JWT 身份驗證與角色權限控制

## 技術選型
- **後端語言**: Go
- **Web 框架**: Gin
- **資料庫**: PostgreSQL
- **資料庫遷移**: Goose
- **依賴管理**: Go Modules
- **環境設定**: Viper

## 環境設定步驟

### 1. 安裝 Go
請確保您的系統已安裝 Go 1.18 或更高版本。您可以從 [Go 官方網站](https://golang.org/doc/install) 下載並安裝。

### 2. 安裝 PostgreSQL
請確保您的系統已安裝 PostgreSQL 資料庫。您可以從 [PostgreSQL 官方網站](https://www.postgresql.org/download/) 下載並安裝。

### 3. 設定資料庫
建立一個新的 PostgreSQL 資料庫，例如 `card_manage_db`。

### 4. 設定專案環境變數
複製 `config.yaml` 並根據您的環境修改資料庫連線資訊、JWT 密鑰等。
```yaml
# config.yaml 範例
server:
  address: "0.0.0.0:8080"
db:
  driver: "postgres"
  source: "postgresql://user:password@localhost:5432/card_manage_db?sslmode=disable"
jwt:
  secret: "your_jwt_secret_key"
  access_token_duration: "15m"
  refresh_token_duration: "24h"
```

### 5. 安裝依賴
在專案根目錄下執行：
```bash
go mod tidy
```

### 6. 運行資料庫遷移
本專案使用 Goose 進行資料庫遷移。請確保您已安裝 Goose：
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

然後，在專案根目錄下執行資料庫遷移：
```bash
goose postgres "postgresql://user:password@localhost:5432/card_manage_db?sslmode=disable" up
```
請將 `user:password@localhost:5432/card_manage_db?sslmode=disable` 替換為您的實際資料庫連線字串。

## 如何啟動服務
在專案根目錄下執行：
```bash
go run cmd/server/main.go
```
服務將會在 `config.yaml` 中設定的地址和端口上啟動 (預設為 `0.0.0.0:8080`)。
