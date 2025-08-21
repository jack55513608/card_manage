# SettlementService 說明文件

`SettlementService` 負責處理玩家收益清算相關的業務邏輯。它允許玩家申請清算其在特定店家的收益，並確保清算過程中的資料一致性，特別是將相關寄售狀態更新為已清算。此服務也使用**資料庫交易**來保證操作的原子性。

## 結構

```go
type SettlementService struct {
	repo            *repository.SettlementRepository
	consignmentRepo *repository.ConsignmentRepository
	storeRepo       *repository.StoreRepository
	db              *sql.DB
}
```

- `repo`: `SettlementRepository` 的實例，用於執行清算資料的持久化操作。
- `consignmentRepo`: `ConsignmentRepository` 的實例，用於更新寄售狀態。
- `storeRepo`: `StoreRepository` 的實例 (雖然在此服務中目前未直接使用，但作為依賴注入)。
- `db`: `*sql.DB` 的實例，用於管理資料庫交易。

## 建構函式

### `NewSettlementService`

```go
func NewSettlementService(
	repo *repository.SettlementRepository,
	consignmentRepo *repository.ConsignmentRepository,
	storeRepo *repository.StoreRepository,
	db *sql.DB,
) *SettlementService
```

- **功能**: 建立並回傳一個新的 `SettlementService` 實例。
- **參數**:
  - `repo`: 必須提供一個 `SettlementRepository` 的實例。
  - `consignmentRepo`: 必須提供一個 `ConsignmentRepository` 的實例。
  - `storeRepo`: 必須提供一個 `StoreRepository` 的實例。
  - `db`: 必須提供一個 `*sql.DB` 的實例，用於開啟資料庫交易。
- **回傳值**:
  - `*SettlementService`: 新建立的 `SettlementService` 實例。

## 方法

### `CreateSettlement`

```go
func (s *SettlementService) CreateSettlement(playerID, storeID int64) (*model.Settlement, error)
```

- **功能**: 允許玩家為其在指定店家的已售出交易申請清算。它會計算玩家的總收益，建立清算紀錄，並將所有相關的寄售狀態更新為 `CLEARED`。此操作在單一資料庫交易中執行。
- **參數**:
  - `playerID` (int64): 申請清算的玩家 ID。
  - `storeID` (int64): 申請清算的店家 ID。
- **回傳值**:
  - `*model.Settlement`: 如果清算申請成功，回傳新建立的清算模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrNoUnsettledTransactions`: 沒有可供清算的交易。
    - 其他內部錯誤 (例如資料庫交易失敗)。
- **內部流程 (資料庫交易)**:
  1. **獲取未清算交易**: 調用 `s.repo.GetUnsettledTransactions` 獲取指定玩家在指定店家的所有未清算 (已售出) 交易。
  2. **計算總收益**: 遍歷所有未清算交易，根據交易價格和抽成比例計算玩家的實際收益總額。
  3. **建立清算模型**: 準備 `model.Settlement` 實例，狀態預設為 `REQUESTED`。
  4. **開啟資料庫交易**: 調用 `s.db.Begin()` 開始一個新的資料庫交易。
  5. **defer Rollback**: 使用 `defer tx.Rollback()` 確保在函式結束時，如果交易未被明確提交，則會自動回滾。
  6. **建立清算紀錄**: 調用 `s.repo.CreateSettlement` (需要傳入交易物件) 將清算紀錄儲存到資料庫。
  7. **更新寄售狀態**: 遍歷所有與清算相關的寄售 ID，調用 `s.consignmentRepo.UpdateConsignmentStatusInTx` (需要傳入交易物件) 將其狀態更新為 `CLEARED`。
  8. **提交交易**: 如果上述所有操作都成功，調用 `tx.Commit()` 提交整個交易。

### `CompleteSettlement`

```go
func (s *SettlementService) CompleteSettlement(storeUserID, settlementID int64) (*model.Settlement, error)
```

- **功能**: 允許店家將清算申請標記為已完成。此方法**尚未實作**。
- **參數**:
  - `storeUserID` (int64): 執行完成操作的店家使用者 ID。
  - `settlementID` (int64): 要完成的清算申請 ID。
- **回傳值**:
  - `*model.Settlement`: (預期) 如果完成成功，回傳更新後的清算模型。
  - `error`: (預期) 如果發生錯誤，回傳錯誤資訊。
