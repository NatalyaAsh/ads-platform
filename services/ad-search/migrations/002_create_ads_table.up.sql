CREATE TABLE IF NOT EXISTS "ads" (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10,2) NOT NULL,
  user_id INT NOT NULL,
  category_id INT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  views INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_ads_category_id 
    FOREIGN KEY (category_id) REFERENCES categories(id) 
    ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_ads_user_id ON ads(user_id);
CREATE INDEX IF NOT EXISTS idx_ads_category_id ON ads(category_id);
CREATE INDEX IF NOT EXISTS idx_ads_status ON ads(status);
CREATE INDEX IF NOT EXISTS idx_ads_price ON ads(price);
CREATE INDEX IF NOT EXISTS idx_ads_created_at ON ads(created_at);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_ads_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_ads_updated_at 
    BEFORE UPDATE ON ads 
    FOR EACH ROW 
    EXECUTE FUNCTION update_ads_updated_at();