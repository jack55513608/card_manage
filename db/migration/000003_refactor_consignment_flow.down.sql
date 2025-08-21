-- Drop the new foreign key on transactions.
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_consignment_item_id_fkey;

-- Rename the column back.
ALTER TABLE transactions RENAME COLUMN consignment_item_id TO consignment_id;

-- Drop the new tables.
DROP TABLE IF EXISTS consignment_items;
DROP TABLE IF EXISTS consignments;

-- Re-create the old consignments table.
CREATE TABLE consignments (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_id INT NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    card_id INT NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    status VARCHAR(20) CHECK (status IN ('PENDING','LISTED','SOLD','CLEARED')) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Re-add the foreign key to the old consignments table.
-- Note: This assumes the old consignments table had a constraint named transactions_consignment_id_fkey
-- which might not be true if it was named differently. This is a best-effort recreation.
ALTER TABLE transactions ADD CONSTRAINT transactions_consignment_id_fkey
FOREIGN KEY (consignment_id) REFERENCES consignments(id) ON DELETE CASCADE;
