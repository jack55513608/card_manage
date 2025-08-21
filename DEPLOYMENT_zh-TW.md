# 部署指南

本文件為 `card-manage` 專案提供了一份完整的設定與部署指南，涵蓋了本地開發環境以及在 Google Cloud Platform (GCP) 上的部署流程。

整個部署流程已透過 `Makefile`、`Docker` 和 `golang-migrate` 等工具實現自動化。

## 1. 前置需求

在開始之前，請確保您的本機已安裝以下工具：

- **Go**: 1.22.5 或更高版本。
- **Docker & Docker Compose**: 用於運行本地資料庫及應用程式容器化。
- **Google Cloud SDK (`gcloud`)**: 用於與 GCP 平台互動。
  - 請遵循 [官方文件](https://cloud.google.com/sdk/docs/install) 進行安裝。
  - 安裝後，請加裝 `cloud_sql_proxy` 元件：
    ```bash
    gcloud components install cloud_sql_proxy
    ```
- **`golang-migrate` CLI**: 用於管理資料庫結構的遷移 (Migration)。
  - 請使用以下指令安裝：
    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

## 2. 環境設定

### 2.1. 本地環境設定

應用程式會從 `config.yaml` 檔案讀取設定。對於本地開發，預設的設定檔已和 `docker-compose.yml` 完美對應，無需額外修改。

### 2.2. GCP 環境設定

在部署�� GCP 之前，您必須先設定 `Makefile` 中的變數。

1.  **打開專案根目錄的 `Makefile`**。
2.  **找到 `GCP Environment` 區塊**，並將所有預留位置的值替換為您真實的 GCP 專案資訊：
    - `GCP_PROJECT_ID`: 您的 Google Cloud 專案 ID。
    - `GCP_REGION`: 您要部署資源的區域 (例如 `asia-east1`)。
    - `GCP_SERVICE_ACCOUNT`: Cloud Run 服務要使用的服務帳號。
    - `GCP_CLOUD_SQL_INSTANCE`: 您的 Cloud SQL 實例的完整連線名稱。
    - `GCP_DB_URL`: 請更新此 URL 中的資料庫使用者和名稱，使其與您的 Cloud SQL 設定相符。

## 3. 本地開發流程

透過以下指令，您可以在本機快速啟動、管理和測試服務。

### 步驟一：啟動資料庫
此指令會使用 Docker Compose 啟動一個 PostgreSQL 容器。

```bash
make db-up
```

### 步驟二：執行資料庫遷移
此指令會在本地資料庫中建立所有必要的資料表。

```bash
make migrate-up
```

### 步驟三：運行應用程式
此指令會編譯並運行 Go 服務。服務啟動後，可透過 `http://localhost:8080` 存取。

```bash
make run
```

### (可選) 關閉資料庫
當您完成開發後，可以使用此指令關閉資料庫容器。

```bash
make db-down
```

## 4. GCP 部署流程

此流程詳細說明如何將服務部署到 Google Cloud Run，並連接到 Cloud SQL 資料庫。

### 步驟一：GCP 平台驗證
首先，登入您的 Google 帳號，並設定 Docker 以便能推送到 Google Artifact Registry。此步驟只需執行一次。

```bash
gcloud auth login
gcloud config set project your-gcp-project-id
make gcp-auth
```

### 步驟二：建置並推送 Docker 映像檔
此指令會建置應用程式的 Docker 映像檔，並將其推送到您專案的 Artifact Registry。

```bash
make gcp-build-push
```

### 步驟三：部署至 Cloud Run
此指令會將容器映像檔部署到 Cloud Run，並將其與指定的 Cloud SQL 實例連接。

```bash
make g-deploy
```

### 步驟四：在 Cloud SQL 上執行資料庫遷移
若要更新 Cloud SQL 上的資料庫結構，您必須先透過 Cloud SQL Auth Proxy 建立連線。

1.  **啟動 Proxy**：請**開啟一個新的、獨立的終端機視窗**，並執行以下指令。請保持此視窗開啟。
    ```bash
    make cloud-sql-proxy
    ```

2.  **設定資料庫密碼**：在您**原本的終端機視窗**中，將您的 Cloud SQL 資料庫密碼設定為環境變數。
    ```bash
    export GCP_DB_PASSWORD='your-cloud-sql-db-password'
    ```

3.  **執行遷移**：現在，執行遷移指令。
    ```bash
    make gcp-migrate-up
    ```

完成以上步��後，您的服務將成功部署在 Cloud Run 上，並連接到已完成遷移的 Cloud SQL 資料庫。

## 5. 資料庫遷移管理

本專案使用 `golang-migrate` 工具來管理資料庫結構的變更。

- **建立新的遷移檔案**：
  ```bash
  migrate create -ext sql -dir db/migration -seq a_descriptive_name
  ```
  此指令會在 `db/migration` 目錄下建立新的 `up` 和 `down` SQL 檔案。

- **還原上一次的遷移**：
  - 本地資料庫: `make migrate-down`
  - GCP 資料庫: `make gcp-migrate-down` (需要 Proxy 正在運行且已設定密碼)
