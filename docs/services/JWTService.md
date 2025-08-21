# JWTService 說明文件

`JWTService` 負責處理 JSON Web Token (JWT) 的生成和驗證。它用於使用者認證和授權，確保 API 請求的安全性。

## 結構

### `JWTService`

```go
type JWTService struct {
	secretKey      string
	expireDuration time.Duration
}
```

- `secretKey`: 用於簽署和驗證 JWT 的密鑰。
- `expireDuration`: JWT 的有效期限。

### `CustomClaims`

```go
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
```

- `UserID`: 使用者的唯一識別符。
- `Role`: 使用者的角色 (例如 `PLAYER`, `STORE`, `ADMIN`)。
- `jwt.RegisteredClaims`: JWT 標準聲明，包含 `ExpiresAt`, `IssuedAt`, `NotBefore`, `Issuer` 等。

## 建構函式

### `NewJWTService`

```go
func NewJWTService(secretKey string, expireDurationStr string) (*JWTService, error)
```

- **功能**: 建立並回傳一個新的 `JWTService` 實例。
- **參數**:
  - `secretKey` (string): 用於簽署和驗證 JWT 的密鑰。
  - `expireDurationStr` (string): JWT 有效期限的字串表示 (例如 "15m", "24h")。
- **回傳值**:
  - `*JWTService`: 新建立的 `JWTService` 實例。
  - `error`: 如果 `expireDurationStr` 無效，回傳錯誤。

## 方法

### `GenerateToken`

```go
func (s *JWTService) GenerateToken(user *model.User) (string, error)
```

- **功能**: 為給定的使用者生成一個新的 JWT。
- **參數**:
  - `user` (*model.User): 包含使用者 ID 和角色的使用者模型。
- **回傳值**:
  - `string`: 生成的 JWT 字串。
  - `error`: 如果生成 token 失敗，回傳錯誤。
- **內部流程**:
  1. 建立 `CustomClaims` 實例，包含使用者 ID、角色和標準聲明 (如過期時間、發行時間)。
  2. 使用 `jwt.SigningMethodHS256` 簽署 token。

### `ValidateToken`

```go
func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error)
```

- **功能**: 驗證給定的 JWT 字串，並回傳其包含的聲明。
- **參數**:
  - `tokenString` (string): 要驗證的 JWT 字串。
- **回傳值**:
  - `*CustomClaims`: 如果 token 有效，回傳其包含的 `CustomClaims`。
  - `error`: 如果 token 無效 (例如簽名不匹配、過期、格式錯誤)，回傳錯誤。
- **內部流程**:
  1. 使用 `jwt.ParseWithClaims` 解析 token 字串。
  2. 驗證 token 的簽名方法和密鑰。
  3. 檢查 token 是否有效 (例如是否過期)。
