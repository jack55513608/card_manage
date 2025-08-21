# API 參考文件

本文件提供卡片遊戲寄售平台後端 API 的詳細參考資訊，旨在協助前端開發者或其他服務整合。

## 認證與授權

所有需要認證的 API 都必須在請求頭中包含 `Authorization: Bearer <token>`。`<token>` 是使用者登入後獲取的 JWT。

### 角色權限

系統定義了以下角色：
- `PLAYER`: 玩家，可以寄售卡片、查詢銷售紀錄、申請清算。
- `STORE`: 店家，可以管理卡片資料庫、設定售價、處理銷售、管理收益、處理清算。
- `ADMIN`: 平台管理員，擁有最高權限，可以管理店家、監控平台數據。

每個 API 端點都會明確標示所需的權限。

---

## 1. 使用者相關 API

### 1.1 註冊新使用者

- **HTTP 方法**: `POST`
- **路徑**: `/api/users/register`
- **功能描述**: 允許新使用者註冊為玩家或店家。
- **需要的權限**: 公開 (無需認證)

- **請求 (Request)**:
  - **Header**: `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "email": "user@example.com",
      "password": "your_secure_password",
      "role": "PLAYER" // 或 "STORE"
    }
    ```
  - **參數說明**:
    - `email` (string, required): 使用者電子郵件，必須是有效格式。
    - `password` (string, required): 密碼，至少 8 個字元。
    - `role` (string, required): 使用者角色，只能是 `PLAYER` 或 `STORE`。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "message": "user created successfully",
      "user_id": 1
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request` (例如：請求格式錯誤、email 無效、密碼不符合要求、role 無效)
    ```json
    {
      "error": "invalid email or password"
    }
    ```
  - **狀態碼**: `409 Conflict` (例如：email 已存在)
    ```json
    {
      "error": "email already exists"
    }
    ```
  - **狀態碼**: `500 Internal Server Error` (例如：伺服器內部錯誤)

### 1.2 使用者登入

- **HTTP 方法**: `POST`
- **路徑**: `/api/users/login`
- **功能描述**: 使用者登入並獲取 JWT 認證 token。
- **需要的權限**: 公開 (無需認證)

- **請求 (Request)**:
  - **Header**: `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "email": "user@example.com",
      "password": "your_secure_password"
    }
    ```
  - **參數說明**:
    - `email` (string, required): 使用者電子郵件。
    - `password` (string, required): 密碼。

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request` (例如：請求格式錯誤、email 無效)
  - **狀態碼**: `401 Unauthorized` (例如：email 或密碼不正確)
    ```json
    {
      "error": "invalid email or password"
    }
    ```
  - **狀態碼**: `500 Internal Server Error`

---

## 2. 店家相關 API

### 2.1 建立店家資訊

- **HTTP 方法**: `POST`
- **路徑**: `/api/stores`
- **功能描述**: 允許已註冊為 `STORE` 角色的使用者建立其店家資訊。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "name": "卡牌之家",
      "commission_cash": 10.00, 
      "commission_credit": 5.00
    }
    ```
  - **參數說明**:
    - `name` (string, required): 店家名稱。
    - `commission_cash` (float, required): 現金交易抽成比例 (0-100)。
    - `commission_credit` (float, required): 儲值金交易抽��比例 (0-100)。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "id": 1,
      "user_id": 1,
      "name": "卡牌之家",
      "commission_cash": 10.00,
      "commission_credit": 5.00,
      "created_at": "2025-08-19T10:00:00Z",
      "updated_at": "2025-08-19T10:00:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request` (例如：請求格式錯誤、參數不符合要求)
  - **狀態碼**: `401 Unauthorized` (例如：未提供 token 或 token 無效)
  - **狀態碼**: `403 Forbidden` (例如：非 `STORE` 角色或已建立店家)
  - **狀態碼**: `500 Internal Server Error`

---

## 3. 卡片資料庫相關 API

### 3.1 建立卡片資訊

- **HTTP 方法**: `POST`
- **路徑**: `/api/cards`
- **功能描述**: 允許店家新增標準化的卡片資訊到其資料庫。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "name": "青眼白龍",
      "series": "遊戲王",
      "rarity": "UR",
      "card_number": "SET-001"
    }
    ```
  - **參數說明**:
    - `name` (string, required): 卡片名稱。
    - `series` (string, optional): 卡片系列。
    - `rarity` (string, optional): 卡片稀有度。
    - `card_number` (string, optional): 卡片編號。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "id": 1,
      "store_id": 1,
      "name": "青眼白龍",
      "series": "遊戲王",
      "rarity": "UR",
      "card_number": "SET-001",
      "created_at": "2025-08-19T10:00:00Z",
      "updated_at": "2025-08-19T10:00:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：非 `STORE` 角色或使用者沒有店家)
  - **狀態碼**: `500 Internal Server Error`

### 3.2 取得單一卡片資訊

- **HTTP 方法**: `GET`
- **路徑**: `/api/cards/{id}`
- **功能描述**: 根據卡片 ID 取得單一卡片資訊。只有卡片所屬的店家或平台管理員可以查看。
- **需要的權限**: 需要登入 (`STORE` 或 `ADMIN` 角色)

- **請求 (Request)**:
  - **Header**: `Authorization: Bearer <token>`
  - **Path Parameters**:
    - `id` (integer, required): 卡片 ID。

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    {
      "id": 1,
      "store_id": 1,
      "name": "青眼白龍",
      "series": "遊戲王",
      "rarity": "UR",
      "card_number": "SET-001",
      "created_at": "2025-08-19T10:00:00Z",
      "updated_at": "2025-08-19T10:00:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request` (例如：無效的卡片 ID)
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：無權限查看此卡片)
  - **狀態碼**: `404 Not Found` (例如：卡片不存在)
  - **狀態碼**: `500 Internal Server Error`

### 3.3 列出所有卡片資訊

- **HTTP 方法**: `GET`
- **路徑**: `/api/cards`
- **功能描述**: 列出當前使用者所屬店家的所有卡片資訊。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: `Authorization: Bearer <token>`

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    [
      {
        "id": 1,
        "store_id": 1,
        "name": "青眼白龍",
        "series": "遊戲王",
        "rarity": "UR",
        "card_number": "SET-001",
        "created_at": "2025-08-19T10:00:00Z",
        "updated_at": "2025-08-19T10:00:00Z"
      },
      // ... 更多卡片
    ]
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `500 Internal Server Error`

### 3.4 更新卡片資訊

- **HTTP 方法**: `PUT`
- **路徑**: `/api/cards/{id}`
- **功能描述**: 更新指定卡片的資訊。只有卡片所屬的店家可以更新。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Path Parameters**:
    - `id` (integer, required): 卡片 ID。
  - **Body**: 
    ```json
    {
      "name": "青眼白龍 (新版)",
      "series": "遊戲王DM",
      "rarity": "SR",
      "card_number": "SET-001-V2"
    }
    ```
  - **參數說明**:
    - `name` (string, required): 卡片名稱。
    - `series` (string, optional): 卡片系列。
    - `rarity` (string, optional): 卡片稀有度。
    - `card_number` (string, optional): 卡片編號。

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    {
      "id": 1,
      "store_id": 1,
      "name": "青眼白龍 (新版)",
      "series": "遊戲王DM",
      "rarity": "SR",
      "card_number": "SET-001-V2",
      "created_at": "2025-08-19T10:00:00Z",
      "updated_at": "2025-08-19T10:30:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：無權限更新此卡片)
  - **狀態碼**: `404 Not Found`
  - **狀態碼**: `500 Internal Server Error`

### 3.5 刪除卡片資訊

- **HTTP 方法**: `DELETE`
- **路徑**: `/api/cards/{id}`
- **功能描述**: 刪除指定卡片資訊。只有卡片所屬的店家可以刪除。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: `Authorization: Bearer <token>`
  - **Path Parameters**:
    - `id` (integer, required): 卡片 ID。

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    {
      "message": "card deleted successfully"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：無權限刪除此卡片)
  - **狀態碼**: `404 Not Found`
  - **狀態碼**: `500 Internal Server Error`

---

## 4. 寄售相關 API

### 4.1 建立寄售申請

- **HTTP 方法**: `POST`
- **路徑**: `/api/consignments`
- **功能描述**: 玩���提交卡片寄售申請。
- **需要的權限**: 需要登入 (`PLAYER` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "store_id": 1,
      "card_id": 1,
      "quantity": 5
    }
    ```
  - **參數說明**:
    - `store_id` (integer, required): 寄售的店家 ID。
    - `card_id` (integer, required): 寄售的卡片 ID (該卡片必須存在於指定店家的卡片資料庫中)。
    - `quantity` (integer, required): 寄售數量，必須大於 0。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "id": 1,
      "player_id": 1,
      "store_id": 1,
      "card_id": 1,
      "quantity": 5,
      "status": "PENDING",
      "created_at": "2025-08-19T10:00:00Z",
      "updated_at": "2025-08-19T10:00:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：非 `PLAYER` 角色)
  - **狀態碼**: `409 Conflict` (例如：提供的卡片不屬於指定的店家)
  - **狀態碼**: `500 Internal Server Error`

### 4.2 列出寄售申請

- **HTTP 方法**: `GET`
- **路徑**: `/api/consignments`
- **功能描述**: 列出當前使用者相關的寄售申請。玩家可查看自己的，店家可查看其店家的所有寄售申請。
- **需要的權限**: 需要登入 (`PLAYER` 或 `STORE` 角色)

- **請求 (Request)**:
  - **Header**: `Authorization: Bearer <token>`

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    [
      {
        "id": 1,
        "player_id": 1,
        "store_id": 1,
        "card_id": 1,
        "quantity": 5,
        "status": "PENDING",
        "created_at": "2025-08-19T10:00:00Z",
        "updated_at": "2025-08-19T10:00:00Z"
      },
      // ... 更多寄售申請
    ]
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `500 Internal Server Error`

### 4.3 更新寄售狀態

- **HTTP 方法**: `PUT`
- **路徑**: `/api/consignments/{id}/status`
- **功能描述**: 更新指定寄售申請的狀態。只有相關店家可以更新。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Path Parameters**:
    - `id` (integer, required): 寄售申請 ID。
  - **Body**: 
    ```json
    {
      "status": "LISTED" // 或 "SOLD", "CLEARED"
    }
    ```
  - **參數說明**:
    - `status` (string, required): 新的寄售狀態，只能是 `LISTED`, `SOLD`, `CLEARED`。

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    {
      "id": 1,
      "player_id": 1,
      "store_id": 1,
      "card_id": 1,
      "quantity": 5,
      "status": "LISTED",
      "created_at": "2025-08-19T10:00:00Z",
      "updated_at": "2025-08-19T10:30:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：無權限更新此寄售申請)
  - **狀態碼**: `404 Not Found` (例如：寄售申請不存在)
  - **狀態碼**: `500 Internal Server Error`

---

## 5. 交易相關 API

### 5.1 建立交易紀錄

- **HTTP 方法**: `POST`
- **路徑**: `/api/transactions`
- **功能描述**: 店家為已售出的寄售卡片建立交易紀錄。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "consignment_id": 1,
      "price": 100.00,
      "payment_method": "CASH" // 或 "CREDIT"
    }
    ```
  - **參數說明**:
    - `consignment_id` (integer, required): 相關的寄售申請 ID。
    - `price` (float, required): 實際售出價格，必須大於 0。
    - `payment_method` (string, required): 支付方式，只能是 `CASH` 或 `CREDIT`。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "id": 1,
      "consignment_id": 1,
      "store_id": 1,
      "price": 100.00,
      "payment_method": "CASH",
      "commission_rate": 10.00, // 根據店家設定的抽成比例
      "created_at": "2025-08-19T10:00:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：非 `STORE` 角色或無權限操作此寄售)
  - **狀態碼**: `409 Conflict` (例如：寄售申請不存在、寄售已售出)
  - **狀態碼**: `500 Internal Server Error`

---

## 6. 清算相關 API

### 6.1 建立清算申請

- **HTTP 方法**: `POST`
- **路徑**: `/api/settlements`
- **功能描述**: 玩家申���清算其在指定店家的收益。
- **需要的權限**: 需要登入 (`PLAYER` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "store_id": 1
    }
    ```
  - **參數說明**:
    - `store_id` (integer, required): 申請清算的店家 ID。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "id": 1,
      "player_id": 1,
      "store_id": 1,
      "amount": 90.00, // 實際清算金額 (扣除抽成)
      "status": "REQUESTED",
      "created_at": "2025-08-19T10:00:00Z",
      "updated_at": "2025-08-19T10:00:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：非 `PLAYER` 角色)
  - **狀態碼**: `409 Conflict` (例如：沒有可供清算的交易)
  - **狀態碼**: `500 Internal Server Error`
