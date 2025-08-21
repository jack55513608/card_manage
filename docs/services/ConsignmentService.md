# ConsignmentService 說明文件

`ConsignmentService` 負責處理卡片寄售相關的業務邏輯，包括寄售申請的建立、查詢和狀態更新。它與 `ConsignmentRepository`、`CardRepository` 和 `StoreRepository` 互動以進行資料庫操作，並執行相關的驗證。

## 結構

```go
type ConsignmentService struct {
	consignmentRepo *repository.ConsignmentRepository
	cardRepo        *repository.CardRepository
	storeRepo       *repository.StoreRepository
}
```

- `consignmentRepo`: `ConsignmentRepository` 的實例，用於執行寄售資料的持久化操作。
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
- **參數**:
  - `consignmentRepo`: 必須提供一個 `ConsignmentRepository` 的實例。
  - `cardRepo`: 必須提供一個 `CardRepository` 的實例。
  - `storeRepo`: 必須提供一個 `StoreRepository` 的實例。
- **回傳值**:
  - `*ConsignmentService`: 新建立的 `ConsignmentService` 實例。

## 方法

### `CreateConsignment`

```go
func (s *ConsignmentService) CreateConsignment(playerID, storeID, cardID int64, quantity int) (*model.Consignment, error)
```

- **功能**: 允許玩家建立新的卡片寄售申請。它會驗證所選卡片是否屬於指定的店家。
- **參數**:
  - `playerID` (int64): 提交寄售申請的玩家 ID。
  - `storeID` (int64): 寄售目標店家的 ID。
  - `cardID` (int64): 寄售卡片的 ID。
  - `quantity` (int): 寄售的卡片數量。
- **回傳值**:
  - `*model.Consignment`: 如果建立成功，回傳新建立的寄售模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrInvalidCardForStore`: 所選卡片不屬於指定的店家。
    - 其他內部錯誤 (例如資料庫操作失敗)。
- **內部流程**:
  1. 調用 `cardRepo.GetCardByID` 驗證卡片是否存在且屬於指定的店家。
  2. 建立 `model.Consignment` 實例，狀態預設為 `PENDING`。
  3. 調用 `consignmentRepo.CreateConsignment` 將寄售資訊儲存到資料庫。

### `ListConsignmentsForUser`

```go
func (s *ConsignmentService) ListConsignmentsForUser(userID int64, userRole string) ([]model.Consignment, error)
```

- **功能**: 根據使用者的角色列出相關的寄售申請。玩家可以查看自己的寄售，店家可以查看其店家的所有寄售。
- **參數**:
  - `userID` (int64): 請求列出寄售的使用者 ID。
  - `userRole` (string): 使用者的角色 (`PLAYER` 或 `STORE`)。
- **回傳值**:
  - `[]model.Consignment`: 相關的寄售列表。如果沒有找到，回傳空切片。
  - `error`: 如果發生錯誤，回傳錯誤資訊。
- **內部流程**:
  - 如果 `userRole` 是 `PLAYER`，調用 `consignmentRepo.ListConsignmentsByPlayer`。
  - 如果 `userRole` 是 `STORE`，先調用 `storeRepo.GetStoreByUserID` 查找店家，然後調用 `consignmentRepo.ListConsignmentsByStore`。

### `UpdateConsignmentStatus`

```go
func (s *ConsignmentService) UpdateConsignmentStatus(storeUserID, consignmentID int64, newStatus model.ConsignmentStatus) (*model.Consignment, error)
```

- **功能**: 允許店家更新指定寄售申請的狀態。它會驗證操作者是否為該寄售所屬店家的擁有者。
- **參數**:
  - `storeUserID` (int64): 執行更新操作的店家使用者 ID。
  - `consignmentID` (int64): 要更新的寄售申請 ID。
  - `newStatus` (model.ConsignmentStatus): 新的寄售狀態 (例如 `LISTED`, `SOLD`, `CLEARED`)。
- **回傳值**:
  - `*model.Consignment`: 如果更新成功，回傳更新後的寄售模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrConsignmentNotFound`: 寄售申請不存在。
    - `service.ErrForbidden`: 使用者無權限更新此寄售申請。
    - 其他內部錯誤。
- **內部流程**:
  1. 調用 `consignmentRepo.GetConsignmentByID` 查找寄售申請。
  2. 調用 `verifyStoreOwnership` 驗證 `storeUserID` 是否擁有該寄售所屬的店家。
  3. 調用 `consignmentRepo.UpdateConsignmentStatus` 更新資料庫中的寄售狀態。

### 輔助方法

### `verifyStoreOwnership`

```go
func (s *ConsignmentService) verifyStoreOwnership(userID, storeID int64) error
```

- **功能**: 內部輔助方法，用於檢查給定的使用者是否擁有指定的店家。此方法與 `CardService` 中的同名方法功能相同，未來可考慮重構為共享服務。
- **參數**:
  - `userID` (int64): 使用者 ID。
  - `storeID` (int64): 店家 ID。
- **回傳值**:
  - `error`: 如果使用者不擁有該店家，回傳 `service.ErrForbidden`；否則回傳 `nil`。
