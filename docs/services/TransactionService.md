# TransactionService 說明文件

`TransactionService` 負責處理交易相關的業務邏輯，特別是建立新的交易紀錄。此服務的關鍵在於它會使用**資料庫交易 (Database Transaction)** 來確保操作的原子性，即在建立交易紀錄的同時更新寄售狀態，保證這些操作要麼全部成功，要麼全部失敗，以維護資料的一致性。

## 結構

```go
type TransactionService struct {
	repo            *repository.TransactionRepository
	consignmentRepo *repository.ConsignmentRepository
	storeRepo       *repository.StoreRepository
	db              *sql.DB
}
```

- `repo`: `TransactionRepository` 的實例，用於執行交易資料的持久化操作。
- `consignmentRepo`: `ConsignmentRepository` 的實例，用於更新寄售狀態。
- `storeRepo`: `StoreRepository` 的實例，用於獲取店家資訊和抽成比例。
- `db`: `*sql.DB` 的實例，用於管理資料庫交易。

## 建構函式

### `NewTransactionService`

```go
func NewTransactionService(
	repo *repository.TransactionRepository,
	consignmentRepo *repository.ConsignmentRepository,
	storeRepo *repository.StoreRepository,
	db *sql.DB,
) *TransactionService
```

- **功能**: 建立並回��一個新的 `TransactionService` 實例。
- **參數**:
  - `repo`: 必須提供一個 `TransactionRepository` 的實例。
  - `consignmentRepo`: 必須提供一個 `ConsignmentRepository` 的實例。
  - `storeRepo`: 必須提供一個 `StoreRepository` 的實例。
  - `db`: 必須提供一個 `*sql.DB` 的實例，用於開啟資料庫交易。
- **回傳值**:
  - `*TransactionService`: 新建立的 `TransactionService` 實例。

## 方法

### `CreateTransaction`

```go
func (s *TransactionService) CreateTransaction(storeUserID, consignmentID int64, price float64, paymentMethod model.PaymentMethod) (*model.Transaction, error)
```

- **功能**: 建立一筆新的交易紀錄，並將相關的寄售狀態更新為 `SOLD`。此操作在單一資料庫交易中執行，確保原子性。
- **參數**:
  - `storeUserID` (int64): 執行交易的店家使用者 ID。
  - `consignmentID` (int64): 相關的寄售申請 ID。
  - `price` (float64): 實際售出價格。
  - `paymentMethod` (model.PaymentMethod): 支付方式 (`CASH` 或 `CREDIT`)。
- **回傳值**:
  - `*model.Transaction`: 如果交易成功，回傳新建立的交易模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrConsignmentNotFound`: 寄售申請不存在。
    - `service.ErrConsignmentAlreadySold`: 寄售申請已售出或已清算。
    - `service.ErrForbidden`: 店家無權限操作此寄售。
    - 其他內部錯誤 (例如資料庫交易失敗)。
- **內部流程 (資料庫交易)**:
  1. **驗證**: 獲取寄售資訊，驗證其存在性、狀態，並確認 `storeUserID` 擁有該寄售所屬的店家。
  2. **計算抽成比例**: 根據 `paymentMethod` 和店家的設定，確定適用的抽成比例。
  3. **建立交易模型**: 準備 `model.Transaction` 實例。
  4. **開啟資料庫交易**: 調用 `s.db.Begin()` 開始一個新的資料庫交易。
  5. **defer Rollback**: 使用 `defer tx.Rollback()` 確保在函式結束時，如果交易未被明確提交，則會自動回滾。這是一個重要的錯誤處理機制。
  6. **建立交易紀錄**: 調用 `s.repo.CreateTransactionInTx` (需要傳入交易物件) 將交易紀錄儲存到資料庫。
  7. **更新寄售狀態**: 調用 `s.consignmentRepo.UpdateConsignmentStatusInTx` (需要傳入交易物件) 將寄售狀態更新為 `SOLD`。
  8. **提交交易**: 如果上述所有操作都成功，調用 `tx.Commit()` 提交整個交易。如果任何一步失敗，`defer tx.Rollback()` 將會回滾所有操作。
