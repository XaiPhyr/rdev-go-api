CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT DEFAULT 0,
    path TEXT DEFAULT NULL,
    action VARCHAR(45) DEFAULT NULL,
    response_status INT DEFAULT 0,
    module_id TEXT DEFAULT NULL,
    module VARCHAR(45) DEFAULT NULL,
    before_change JSONB DEFAULT NULL,
    after_change JSONB DEFAULT NULL,
    ip_address INET,
    user_agent TEXT DEFAULT NULL,
    error_message TEXT DEFAULT NULL,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_audit_logs_module_id ON audit_logs (module, module_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at);
CREATE INDEX idx_audit_logs_user_id ON audit_logs (user_id);