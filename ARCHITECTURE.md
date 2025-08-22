# 系統架構說明文件

本文件旨在提供卡片遊戲寄售平台後端服務的內部架構說明，解釋專案的設計理念、組織方式以及核心模組的運作。

## 1. 分層架構

本專案採用經典的三層架構模式，旨在實現職責分離、提高程式碼的可維護性和可測試性：

- **API 層 (Handler)**: 負責處理 HTTP 請求和回應。它接收來自客戶端的請求，進行請求參數的驗證和綁定，然後將請求轉發給 Service 層處理。API 層不包含業務邏輯，主要負責協調和數據格式轉換。
  - 對應目錄: `internal/api`

- **Service 層 (Business Logic)**: 包含核心業務邏輯。它接收來自 API 層的請求，協調多個 Repository 進行數據操作，執行業務規則驗證，並處理複雜的業務流程。Service 層是業務邏輯的中心。
  - 對應目錄: `internal/service`

- **Repository 層 (Data Access)**: 負責與資料庫進行交互。它封裝了所有數據存取邏輯，提供 CRUD (Create, Read, Update, Delete) 操作的介面，將業務邏輯與底層資料庫實現解耦。Service 層通過 Repository 介面與資料庫通信，而無需關心具體的資料庫操作細節。
  - 對應目錄: `internal/repository`

這種分層設計使得各層職責清晰，修改某一層的實現不會影響到其他層，例如更換資料庫類型或調整 API 接口，都能在不大幅改動其他層的情況下進行。

## 2. 目錄結構

專案的目錄結構遵循 Go 語言的標準佈局，主要目錄及其用途如下：

- `cmd/`: 包含主應用程式的入口點。每個子目錄代表一個獨立的可執行應用程式。
  - `cmd/server/main.go`: 後端 API 服務的啟動入口。

- `db/`: 包含資料庫相關的檔案。
  - `db/migration/`: 存放資料庫遷移腳本 (SQL 檔案)，用於管理資料庫 schema 的版本控制。

- `internal/`: 存放不希望被外部專案直接引用的私有應用程式和函式庫程式碼。這是專案的核心業務邏輯和內部組件的所在地。
  - `internal/api/`: API 處理器 (Handler) 和中間件的定義。
  - `internal/config/`: 應用程式配置的讀取和管理。
  - `internal/model/`: 資料庫模型 (struct) 的定義，對應資料庫中的表結構。
  - `internal/repository/`: 資料庫操作介面和實現，負責數據持久化。
  - `internal/service/`: 業務邏輯服務的實現。

- `config.yaml`: 應用程式的配置檔案，用於設定資料庫連線、JWT 密鑰、伺服器���址等。
- `Dockerfile`: Docker 容器化配置檔案。
- `go.mod`, `go.sum`: Go Modules 依賴管理檔案。
- `prd.txt`: 產品需求文件。

## 3. 核心模組

### 3.1 認證與授權 (JWT)

- **模組**: `internal/service/jwt_service.go`, `internal/api/auth_middleware.go`, `internal/api/role_middleware.go`
- **運作方式**: 
  - 使用者登入成功後，`jwt_service` 會生成一個 JWT (JSON Web Token)，其中包含使用者的 ID、角色等資訊。
  - `auth_middleware` 負責驗證每個請求中的 JWT。它會從請求頭中提取 token，並使用 `jwt_service` 進行驗證。如果 token 無效或缺失，請求將被拒絕。
  - `role_middleware` 則在 `auth_middleware` 之後執行，根據 JWT 中包含的使用者角色，判斷其是否有權限訪問特定的 API 端點。這實現了基於角色的訪問控制 (RBAC)。

### 3.2 設定管理

- **模組**: `internal/config/config.go`
- **運作方式**: 使用 Viper 函式庫從 `config.yaml` 檔案中讀取應用程式的配置。這使得配置可以集中管理，並且易於修改和部署。

### 3.3 資料庫遷移

- **工具**: Goose
- **運作方式**: `db/migration` 目錄下的 SQL 檔案用於管理資料庫 schema 的版本。通過 Goose 工具，可以執行 `up` (應用新的遷移) 和 `down` (回滾遷移) 等操作，確保資料���結構與程式碼版本保持一致，便於開發和部署。

## 4. 資料庫設計與交易 (Transaction)

在 `internal/repository` 層中，許多資料庫操作都可能涉及到多個步驟，為了確保資料的一致性和完整性，專案在必要時使用了資料庫交易 (Transaction)。

- **使用場景**: 
  - 例如，在**建立寄售申請**時，需要在一次資料庫交易中，同時建立一筆父層的 `consignments` (寄售請求) 紀錄以及多筆子層的 `consignment_items` (寄售品項) 紀錄，以確保資料的完整性。
  - 另一個例子是**建立銷售紀錄**，它需要在同一次交易中，建立一筆 `transactions` 紀錄，並將對應的 `consignment_item` 狀態更新為 `SOLD`。

- **實現方式**: Repository 層會提供開始交易、提交交易和回滾交易的方法。Service 層在執行需要原子性操作的業務邏輯時，會調用 Repository 層的交易相關方法來包裹一系列的資料庫操作，確保數據的完整性。

這種設計模式確保了即使在併發或錯誤發生的情況下，資料庫中的數據也能保持正確和一致。
