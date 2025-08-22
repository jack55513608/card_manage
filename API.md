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

(此區塊維持不變)

---

## 2. 店家相關 API

(此區塊維持不變)

---

## 3. 卡片資料庫相關 API

(此區塊維持不變)

---

## 4. 寄售相關 API

### 4.1 建立寄售申請

- **HTTP 方法**: `POST`
- **路徑**: `/api/consignments`
- **功能描述**: 玩家為一張或多張卡片提交一個寄售請求。
- **需要的權限**: 需要登入 (`PLAYER` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "store_id": 1,
      "card_ids": [101, 102, 102, 103] 
    }
    ```
  - **參數說明**:
    - `store_id` (integer, required): 寄售的店家 ID。
    - `card_ids` (array of integer, required): 欲寄售的卡片 ID 列表。如果想寄售多張相同的卡片，請在列表中包含多次該 ID。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "id": 1,
      "player_id": 1,
      "store_id": 1,
      "status": "PROCESSING",
      "items": [
        {
          "id": 1,
          "consignment_id": 1,
          "card_id": 101,
          "status": "PENDING",
          "rejection_reason": "",
          "created_at": "2025-08-21T15:00:00Z",
          "updated_at": "2025-08-21T15:00:00Z"
        },
        {
          "id": 2,
          "consignment_id": 1,
          "card_id": 102,
          "status": "PENDING",
          "rejection_reason": "",
          "created_at": "2025-08-21T15:00:00Z",
          "updated_at": "2025-08-21T15:00:00Z"
        }
      ],
      "created_at": "2025-08-21T15:00:00Z",
      "updated_at": "2025-08-21T15:00:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `500 Internal Server Error`

### 4.2 更新寄售品項狀態

- **HTTP 方法**: `PUT`
- **路徑**: `/api/consignments/items/{itemId}`
- **功能描述**: 店家核可或拒絕寄售請求中的單個品項。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Path Parameters**:
    - `itemId` (integer, required): 寄售品項 ID。
  - **Body**: 
    ```json
    {
      "status": "APPROVED",
      "reason": ""
    }
    ```
  - **參數說明**:
    - `status` (string, required): 新的狀態，只能是 `APPROVED` 或 `REJECTED`。
    - `reason` (string, optional): 當狀態為 `REJECTED` 時，可以提供拒絕原因。

- **成功回應 (Success Response)**:
  - **狀態碼**: `200 OK`
  - **Body**: 
    ```json
    {
      "id": 1,
      "consignment_id": 1,
      "card_id": 101,
      "status": "APPROVED",
      "rejection_reason": "",
      "created_at": "2025-08-21T15:00:00Z",
      "updated_at": "2025-08-21T15:05:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden` (例如：無權限更新此品項)
  - **狀態碼**: `404 Not Found` (例如：品項不存在)
  - **狀態碼**: `409 Conflict` (例如：狀態無法被更新)
  - **狀態碼**: `500 Internal Server Error`

---

## 5. 交易相關 API

### 5.1 建立交易紀錄

- **HTTP 方法**: `POST`
- **路徑**: `/api/transactions`
- **功能描述**: 店家為已售出的寄售品項建立交易紀錄。
- **需要的權限**: 需要登入 (`STORE` 角色)

- **請求 (Request)**:
  - **Header**: 
    - `Authorization: Bearer <token>`
    - `Content-Type: application/json`
  - **Body**: 
    ```json
    {
      "consignment_item_id": 1,
      "price": 100.00,
      "payment_method": "CASH"
    }
    ```
  - **參數說明**:
    - `consignment_item_id` (integer, required): 相關的寄售品項 ID。
    - `price` (float, required): 實際售出價格，必須大於 0。
    - `payment_method` (string, required): 支付方式，只能是 `CASH` 或 `CREDIT`。

- **成功回應 (Success Response)**:
  - **狀態碼**: `201 Created`
  - **Body**: 
    ```json
    {
      "id": 1,
      "consignment_item_id": 1,
      "store_id": 1,
      "price": 100.00,
      "payment_method": "CASH",
      "commission_rate": 10.00,
      "created_at": "2025-08-21T15:10:00Z"
    }
    ```

- **可能的錯誤回應 (Error Response)**:
  - **狀態碼**: `400 Bad Request`
  - **狀態碼**: `401 Unauthorized`
  - **狀態碼**: `403 Forbidden`
  - **狀態碼**: `404 Not Found` (例如：品項不存在)
  - **狀態碼**: `409 Conflict` (例如：品項未核可或已售出)
  - **狀態碼**: `500 Internal Server Error`

---

## 6. 清算相關 API

(此區塊維持不變)