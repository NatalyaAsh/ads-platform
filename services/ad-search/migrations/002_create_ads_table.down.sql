DROP TRIGGER IF EXISTS trigger_ads_updated_at ON ads;

DROP FUNCTION IF EXISTS update_ads_updated_at();

DROP INDEX IF EXISTS idx_ads_created_at;
DROP INDEX IF EXISTS idx_ads_price;
DROP INDEX IF EXISTS idx_ads_status;
DROP INDEX IF EXISTS idx_ads_category_id;
DROP INDEX IF EXISTS idx_ads_user_id;

DROP TABLE IF EXISTS ads CASCADE;
