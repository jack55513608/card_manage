# StoreService 說明文件

`StoreService` 負責處理店家相關的業務邏輯，主要包括建立新的店家資訊。它與 `StoreRepository` 互動以進行資料庫操作。

## 結構

```go
type StoreService struct {
	storeRepo *repository.StoreRepository
}
```

- `storeRepo`: `StoreRepository` 的實例，用於執行店家資料的持久化操作。

## 建構函式

### `NewStoreService`

```go
func NewStoreService(storeRepo *repository.StoreRepository) *StoreService
```

- **功能**: 建立並回傳一個新的 `StoreService` 實例。
- **參數**:
  - `storeRepo`: 必須提供一個 `StoreRepository` 的實例。
- **回傳值**:
  - `*StoreService`: 新建立的 `StoreService` 實例。

## 方法

### `CreateStore`

```go
func (s *StoreService) CreateStore(userID int64, name string, commissionCash, commissionCredit float64) (*model.Store, error)
```

- **功能**: 處理建立新店家的業務邏輯。它會將店家資訊與提供的使用者 ID 關聯起來。
- **參數**:
  - `userID` (int64): 建立店家的使用者 ID。
  - `name` (string): 店家名稱。
  - `commissionCash` (float64): ���金交易的抽成比例。
  - `commissionCredit` (float64): 儲值金交易的抽成比例。
- **回傳值**:
  - `*model.Store`: 如果建立成功，回傳新建立的店家模型。
  - `error`: 如果發生錯誤 (例如資料庫操作失敗)，回傳錯誤資訊。
- **內部流程**:
  1. 建立 `model.Store` 實例，並填入提供的資訊。
  2. 調用 `storeRepo.CreateStore` 將店家資訊儲存到資料庫。
