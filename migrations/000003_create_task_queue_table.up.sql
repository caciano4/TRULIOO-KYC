CREATE TABLE task_queue (
    id SERIAL PRIMARY KEY,
    task_data JSONB NOT NULL,  -- Store task details as JSON
    status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, or failed
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
