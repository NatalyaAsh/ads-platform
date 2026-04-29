CREATE TABLE IF NOT EXISTS "refresh_tokens" (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  token VARCHAR(512) NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE
);        

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id 
    ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at 
    ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked
        ON refresh_tokens(revoked);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id_revoked
        ON refresh_tokens(user_id, revoked);

