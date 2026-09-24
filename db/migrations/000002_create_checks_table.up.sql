Create table if not exists checks(
                       id 	BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                       target_id	 BIGINT NOT NULL REFERENCES target(id) ON DELETE CASCADE,
                       status_code  	 SMALLINT,
                       checked_at  	  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       latency_ms	INT
);