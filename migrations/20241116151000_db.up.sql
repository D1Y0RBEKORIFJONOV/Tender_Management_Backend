CREATE  TABLE  IF NOT EXISTS users (
             id VARCHAR(255) PRIMARY KEY,
             username VARCHAR(255) NOT NULL,
             password VARCHAR(255) NOT NULL,
             role VARCHAR(50) NOT NULL CHECK (role IN ('client', 'contractor')),
             email VARCHAR(255) NOT NULL UNIQUE
);

CREATE   TABLE IF NOT EXISTS tenders (
             id VARCHAR(255) PRIMARY KEY,
             client_id VARCHAR(255) NOT NULL,
             title VARCHAR(255) NOT NULL,
             description TEXT NOT NULL,
             deadline TIMESTAMP NOT NULL,
             budget NUMERIC(15, 2) NOT NULL CHECK (budget > 0),
             status VARCHAR(50) NOT NULL CHECK (status IN ('open', 'closed', 'awarded')),
             fileattachment TEXT,
             created_at TIMESTAMP NOT NULL,
             CONSTRAINT fk_tender_client FOREIGN KEY (client_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE  TABLE IF NOT EXISTS bids (
     id VARCHAR(255) PRIMARY KEY,
     tender_id VARCHAR(255) NOT NULL,
     contractor_id VARCHAR(255) NOT NULL,
     price NUMERIC(15, 2) NOT NULL CHECK (price > 0),
     delivery_time INT NOT NULL CHECK (delivery_time > 0),
     comments TEXT,
     status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected')),
     CONSTRAINT fk_bid_tender FOREIGN KEY (tender_id) REFERENCES tenders (id) ON DELETE CASCADE,
     CONSTRAINT fk_bid_contractor FOREIGN KEY (contractor_id) REFERENCES users (id) ON DELETE CASCADE
);
