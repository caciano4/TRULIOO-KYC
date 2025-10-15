ALTER TABLE records 
ADD COLUMN am_reports_count INTEGER DEFAULT 0,
ADD COLUMN wl_reports_count INTEGER DEFAULT 0,
ADD COLUMN pep_reports_count INTEGER DEFAULT 0,
ADD COLUMN reports_data TEXT;
