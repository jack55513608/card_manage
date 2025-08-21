-- First, drop the foreign key in transactions that depends on the old consignments table.
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_consignment_id_fkey;

-- Then, drop the old consignments table.
DROP TABLE IF EXISTS consignments;

-- Create the new consignments table (for requests).
CREATE TABLE consignments (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_id INT NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    status VARCHAR(20) CHECK (status IN ('PROCESSING', 'COMPLETED')) NOT NULL DEFAULT 'PROCESSING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create the consignment_items table.
CREATE TABLE consignment_items (
    id SERIAL PRIMARY KEY,
    consignment_id INT NOT NULL REFERENCES consignments(id) ON DELETE CASCADE,
    card_id INT NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    status VARCHAR(20) CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'SOLD', 'CLEARED')) NOT NULL DEFAULT 'PENDING',
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Update the transactions table to reference consignment_items.
-- First, rename the old column to avoid confusion.
ALTER TABLE transactions RENAME COLUMN consignment_id TO consignment_item_id;

-- Then, add the new foreign key constraint.
ALTER TABLE transactions ADD CONSTRAINT transactions_consignment_item_id_fkey
FOREIGN KEY (consignment_item_id) REFERENCES consignment_items(id) ON DELETE CASCADE;
