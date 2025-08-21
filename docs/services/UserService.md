# UserService 說明文件

`UserService` 負責處理使用者相關的核心業務邏輯，包括使用者註冊和登入。它與 `UserRepository` 互動以進行資料庫操作，並使用 `bcrypt` 函式庫處理密碼雜湊。

## 結構

```go
type UserService struct {
	userRepo *repository.UserRepository
}
```

- `userRepo`: `UserRepository` 的實例，用於執行使用者資料的持久化操作。

## 建構函式

### `NewUserService`

```go
func NewUserService(userRepo *repository.UserRepository) *UserService
```

- **功能**: 建立並回傳一個新的 `UserService` 實例。
- **參數**:
  - `userRepo`: 必須提供一個 `UserRepository` 的實例。
- **回傳值**:
  - `*UserService`: 新建立的 `UserService` 實例。

## 方法

### `Register`

```go
func (s *UserService) Register(email, password, role string) (*model.User, error)
```

- **功能**: 處理新使用者的註冊業務邏輯。它會檢查電子郵件是否已存在，對密碼進行雜湊處理，然後將新使用者資訊儲存到資料庫中。
- **參數**:
  - `email` (string): 使用者的電子郵件地址。
  - `password` (string): 使用者提供的原始密碼。
  - `role` (string): 使用者的角色 (例如 `PLAYER` 或 `STORE`)。此角色應在 API 層進行初步驗證。
- **回傳值**:
  - `*model.User`: 如果註冊成功，回傳新建立的使用者模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `service.ErrEmailExists`: 電子郵件已存在。
    - 其他內部錯誤 (例如密碼雜湊失敗、資料庫操作失敗)。
- **內部流程**:
  1. 調用 `userRepo.GetUserByEmail` 檢查電子郵件是否已註冊。
  2. 使用 `bcrypt.GenerateFromPassword` 對密碼進行雜湊。
  3. 建立 `model.User` 實例。
  4. 調用 `userRepo.CreateUser` 將使用者資訊儲存到資料庫。

### `Login`

```go
func (s *UserService) Login(email, password string) (*model.User, error)
```

- **功能**: 處理使用者登入業務邏輯。它會根據電子郵件查找使用者，並驗證提供的密碼是否正確。
- **參數**:
  - `email` (string): 使用者的電子郵件地址。
  - `password` (string): 使用者提供的原始密碼。
- **回傳值**:
  - `*model.User`: 如果登入成功，回傳匹配的使用者模型。
  - `error`: 如果發生錯誤，回傳錯誤資訊。可能的錯誤包括：
    - `errors.New("invalid credentials")`: 電子郵件不存在或密碼不正確。
    - 其他內部錯誤 (例如資料庫查詢失敗)。
- **內部流程**:
  1. 調用 `userRepo.GetUserByEmail` 根據電子郵件查找使用者。
  2. 使用 `bcrypt.CompareHashAndPassword` 比較雜湊後的密碼與提供的密碼。
