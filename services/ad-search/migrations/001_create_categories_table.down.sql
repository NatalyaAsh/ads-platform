DROP TRIGGER IF EXISTS trigger_categories_updated_at ON categories;

DROP FUNCTION IF EXISTS update_categories_updated_at();

DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_name;
DROP TABLE IF EXISTS categories CASCADE;
