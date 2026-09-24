create table if not exists target  (
                         id	BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                         url 	VARCHAR(2048) UNIQUE  NOT NULL ,
                         is_tracking  	 BOOLEAN NOT NULL DEFAULT TRUE,
                         interval_sec 	 INTEGER NOT NULL DEFAULT 60,
                         created_at	 TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at 	TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);