CREATE INDEX idx_checks_url_checked ON checks (target_id, checked_at DESC);
CREATE INDEX idx_checks_checked_at  ON checks (checked_at);