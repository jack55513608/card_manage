# ConsignmentService 說明文件

`ConsignmentService` 負責處理卡片寄售相關的業務邏輯。它實現了以「寄售請求 (Request)」為單位，對「寄售品項 (Item)」進行獨立狀態管理的複雜流程。

## 結構

```go
type ConsignmentService struct {
	consignmentRepo *repository.ConsignmentRepository
	cardRepo        *repository.CardRepository
	storeRepo       *repository.StoreRepository
}
```

- `consignmentRepo`: `ConsignmentRepository` 的實例，用於執行寄售請求和品項的持久化操作。
- `cardRepo`: `CardRepository` 的實例，用於驗證卡片資訊。
- `storeRepo`: `StoreRepository` 的實例，用於驗證使用者與店家的關聯。

## 建構函式

### `NewConsignmentService`

```go
func NewConsignmentService(
	consignmentRepo *repository.ConsignmentRepository,
	cardRepo *repository.CardRepository,
	storeRepo *repository.StoreRepository,
) *ConsignmentService
```

- **功能**: 建立並回傳一個新的 `ConsignmentService` 實例。

## 方法

### `CreateConsignment`

```go
func (s *ConsignmentService) CreateConsignment(playerID, storeID int64, cardIDs []int64) (*model.Consignment, error)
```

- **功能**: 允許玩家建立一個新的寄售請求，其中可包含多個寄售品項。
- **參數**:
  - `playerID` (int64): 提交寄售申請的玩家 ID。
  - `storeID` (int64): 寄售目標店家的 ID。
  - `cardIDs` ([]int64): 一個包含多個卡片 ID 的切片，代表玩家希望寄售的所有卡片。
- **回傳值**:
  - `*model.Consignment`: 如果建立成功，回傳新建立的寄售請求模型，其中會包含所有子品項的資訊。
  - `error`: 如果發生錯誤，回傳錯誤資訊。
- **內部流程**:
  1. 建立一個 `model.Consignment` 實例，狀態預設為 `PROCESSING`。
  2. 根據傳入的 `cardIDs` 列表，為每張卡片建立一個對應的 `model.ConsignmentItem` 實例，其初始狀態為 `PENDING`。
  3. 調用 `consignmentRepo.CreateConsignment`，在一次資料庫交易中，將寄售請求和所有寄售品項儲存到資料庫。

### `UpdateConsignmentItemStatus`

```go
func (s *ConsignmentService) UpdateConsignmentItemStatus(storeUserID, itemID int64, newStatus model.ConsignmentItemStatus, reason string) (*model.ConsignmentItem, error)
```

- **功能**: 允許店家核可或拒絕一個指定的寄售品項。它會驗證操作者是否為該品項所屬店家的擁有者。
- **參數**:
  - `storeUserID` (int64): 執行更新操作的店家使用者 ID。
  - `itemID` (int64): 要更新的寄售品項 ID。
  - `newStatus` (model.ConsignmentItemStatus): 新的品項狀態，只能是 `APPROVED` 或 `REJECTED`。
  - `reason` (string): 當狀態更新為 `REJECTED` 時，可以提供拒絕原因。
- **回傳值**:
  - `*model.ConsignmentItem`: 如果更新成功，回傳更新後的寄售品項模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrConsignmentItemNotFound`: 寄售品項不存在。
    - `service.ErrForbidden`: 使用者無權限更新此品項。
    - `service.ErrCannotUpdateStatus`: 品項的當前狀態不允許更新 (例如，不是 `PENDING` 狀態)。
- **內部流程**:
  1. 調用 `consignmentRepo.GetConsignmentItemByID` 查找寄售品項。
  2. 調用 `consignmentRepo.GetConsignmentByID` 獲取父層的寄售請求，以取得 `storeID`。
  3. 調用 `verifyStoreOwnership` 驗證 `storeUserID` 是否擁有該店家。
  4. 驗證狀態轉換是否合法 (只能從 `PENDING` 更新為 `APPROVED` 或 `REJECTED`)。
  5. 調用 `consignmentRepo.UpdateConsignmentItemStatus` 更新資料庫中的品項狀態。

### 輔助方法

### `verifyStoreOwnership`

```go
func (s *ConsignmentService) verifyStoreOwnership(userID, storeID int64) error
```

- **功能**: 內部輔助方法，用於檢查給定的使用者是否擁有指定的店家。