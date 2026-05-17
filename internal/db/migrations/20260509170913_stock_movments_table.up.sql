CREATE TABLE stock_movements (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES products(id) ON DELETE CASCADE,
    change_amount INT NOT NULL, -- e.g., +10 (restock), -1 (sale)
    reason VARCHAR(50) NOT NULL, -- 'SALE', 'RESTOCK', 'ADJUSTMENT', 'RETURN', 'DAMAGE'
    reference_id VARCHAR(100),   -- Order ID or Invoice ID
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);