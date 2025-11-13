-- Migration pour ajouter la table Sessions
-- Peut être exécutée sans danger même si la table existe déjà

-- Table for storing user sessions
CREATE TABLE IF NOT EXISTS Sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Index pour améliorer les performances de recherche par token
CREATE INDEX IF NOT EXISTS idx_session_token ON Sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_session_expiry ON Sessions(expires_at);

-- Nettoyer les sessions expirées
DELETE FROM Sessions WHERE expires_at < NOW();
