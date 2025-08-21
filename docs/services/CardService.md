# CardService 說明文件

`CardService` 負責處理卡片資料庫相關的業務邏輯，包括卡片的建立、查詢、更新和刪除。它與 `CardRepository` 和 `StoreRepository` 互動以進行資料庫操作，並執行權限驗證。

## 結構

```go
type CardService struct {
	cardRepo  *repository.CardRepository
	storeRepo *repository.StoreRepository
}
```

- `cardRepo`: `CardRepository` 的實例，用於執行卡片資料的持久化操作。
- `storeRepo`: `StoreRepository` 的實例，用於驗證使用者與店家的關聯。

## 建構函式

### `NewCardService`

```go
func NewCardService(cardRepo *repository.CardRepository, storeRepo *repository.StoreRepository) *CardService
```

- **功能**: 建立並回傳一個新的 `CardService` 實例。
- **參數**:
  - `cardRepo`: 必須提供一個 `CardRepository` 的實例。
  - `storeRepo`: 必須提供一個 `StoreRepository` 的實例。
- **回傳值**:
  - `*CardService`: 新建立的 `CardService` 實例。

## 方法

### `CreateCard`

```go
func (s *CardService) CreateCard(userID int64, name, series, rarity, cardNumber string) (*model.Card, error)
```

- **功能**: 為指定使用者所屬的店家建立一張新卡片。
- **參數**:
  - `userID` (int64): 建立卡片的使用者 ID (必須是店家角色)。
  - `name` (string): 卡片名稱。
  - `series` (string): 卡片系列。
  - `rarity` (string): 卡片稀有度。
  - `cardNumber` (string): 卡片編號。
- **回傳值**:
  - `*model.Card`: 如果建立成功，回傳新建立的卡片模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrStoreNotFound`: 找不到與使用者關聯的店家。
    - 其他內部錯誤 (例如資料庫操作失敗)。
- **內部流程**:
  1. 調用 `storeRepo.GetStoreByUserID` 查找使用者所屬的店家。
  2. 建立 `model.Card` 實例。
  3. 調用 `cardRepo.CreateCard` 將卡片資訊儲存到資料庫。

### `GetCard`

```go
func (s *CardService) GetCard(userID, cardID int64) (*model.Card, error)
```

- **功能**: 根據卡片 ID 取得單一卡片資訊，並驗證使用者是否有權限查看。
- **參數**:
  - `userID` (int64): 請求查看卡片的使用者 ID。
  - `cardID` (int64): 要查詢的卡片 ID。
- **回傳值**:
  - `*model.Card`: 如果找到卡片且使用者有權限，回傳卡片模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrCardNotFound`: 卡片不存在。
    - `service.ErrForbidden`: 使用者無��限查看此卡片。
    - 其他內部錯誤。
- **內部流程**:
  1. 調用 `cardRepo.GetCardByID` 查找卡片。
  2. 調用 `verifyStoreOwnership` 驗證使用者是否擁有該卡片所屬的店家。

### `ListCardsByCurrentUser`

```go
func (s *CardService) ListCardsByCurrentUser(userID int64) ([]model.Card, error)
```

- **功能**: 列出當前使用者所屬店家的所有卡片資訊。
- **參數**:
  - `userID` (int64): 請求列出卡片的使用者 ID。
- **回傳值**:
  - `[]model.Card`: 該店家下的所有卡片列表。如果沒有店家，回傳空切片。
  - `error`: 如果發生錯誤，回傳錯誤資訊。
- **內部流程**:
  1. 調用 `storeRepo.GetStoreByUserID` 查找使用者所屬的店家。
  2. 調用 `cardRepo.ListCardsByStore` 列出該店家的所有卡片。

### `UpdateCard`

```go
func (s *CardService) UpdateCard(userID, cardID int64, name, series, rarity, cardNumber string) (*model.Card, error)
```

- **功能**: 更新指定卡片的資訊，並驗證使用者是否有權限更新。
- **參數**:
  - `userID` (int64): 請求更新卡片的使用者 ID。
  - `cardID` (int64): 要更新的卡片 ID。
  - `name` (string): 新的卡片名稱。
  - `series` (string): 新的卡片系列。
  - `rarity` (string): 新的卡片稀有度。
  - `cardNumber` (string): 新的��片編號。
- **回傳值**:
  - `*model.Card`: 如果更新成功，回傳更新後的卡片模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrCardNotFound`: 卡片不存在。
    - `service.ErrForbidden`: 使用者無權限更新此卡片。
    - 其他內部錯誤。
- **內部流程**:
  1. 調用 `cardRepo.GetCardByID` 查找卡片。
  2. 調用 `verifyStoreOwnership` 驗證使用者是否擁有該卡片所屬的店家。
  3. 更新卡片模型的欄位。
  4. 調用 `cardRepo.UpdateCard` 更新資料庫中的卡片資訊。

### `DeleteCard`

```go
func (s *CardService) DeleteCard(userID, cardID int64) error
```

- **功能**: 刪除指定卡片，並驗證使用者是否有權限刪除。
- **參數**:
  - `userID` (int64): 請求刪除卡片的使用者 ID。
  - `cardID` (int64): 要刪除的卡片 ID。
- **回傳值**:
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrCardNotFound`: 卡片不存在。
    - `service.ErrForbidden`: 使用者無權限刪除此卡片。
    - 其他內部錯誤。
- **內部流程**:
  1. 調用 `cardRepo.GetCardByID` 查找卡片。
  2. 調用 `verifyStoreOwnership` 驗證使用者是否擁有該卡片所屬的店家。
  3. 調用 `cardRepo.DeleteCard` 從資料庫中刪���卡片。

### 輔助方法

### `verifyStoreOwnership`

```go
func (s *CardService) verifyStoreOwnership(userID, storeID int64) error
```

- **功能**: 內部輔助方法，用於檢查給定的使用者是否擁有指定的店家。
- **參數**:
  - `userID` (int64): 使用者 ID。
  - `storeID` (int64): 店家 ID。
- **回傳值**:
  - `error`: 如果使用者不擁有該店家，回傳 `service.ErrForbidden`；否則回傳 `nil`。
