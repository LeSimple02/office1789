/*-- Creating a database for the Office1789
CREATE DATABASE IF NOT EXISTS Office1789;

USE Office1789;

-- Table for storing users
CREATE TABLE IF NOT EXISTS Users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    domain VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    nboffer INT,
    phonenumber VARCHAR(255),
    date_joined DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME
);

-- Table for storing user drive data (files)
CREATE TABLE IF NOT EXISTS DriveFiles (
    file_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    file_size INT NOT NULL, -- Size in bytes
    file_type VARCHAR(50),
    date_uploaded DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Table for storing mail data
CREATE TABLE IF NOT EXISTS Emails (
    email_id INT AUTO_INCREMENT PRIMARY KEY,
    sender_user_id INT,
    recipient_user_id INT,
    subject VARCHAR(255),
    body TEXT,
    date_sent DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (sender_user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Table for storing email attachments (in case you want to handle attachments)
CREATE TABLE IF NOT EXISTS EmailAttachments (
    attachment_id INT AUTO_INCREMENT PRIMARY KEY,
    email_id INT,
    file_name VARCHAR(255),
    file_path VARCHAR(255),
    file_size INT,
    FOREIGN KEY (email_id) REFERENCES Emails(email_id) ON DELETE CASCADE
);

-- Table for storing user calendar events (Optional)
CREATE TABLE IF NOT EXISTS CalendarEvents (
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    event_title VARCHAR(255),
    event_description TEXT,
    event_start DATETIME,
    event_end DATETIME,
    location VARCHAR(255),
    date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);*/

/*-- Creating a database for the Office1789
CREATE DATABASE IF NOT EXISTS Office1789;

USE Office1789;

-- Table for storing users
CREATE TABLE IF NOT EXISTS Users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    domain VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    nboffer INT,
    phonenumber VARCHAR(255),
    date_joined DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME
);

-- Table for storing user drive data (files)
CREATE TABLE IF NOT EXISTS DriveFiles (
    file_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    file_size INT NOT NULL, -- Size in bytes
    file_type VARCHAR(50),
    date_uploaded DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Table for storing mail data
CREATE TABLE IF NOT EXISTS Emails (
    email_id INT AUTO_INCREMENT PRIMARY KEY,
    sender_user_id INT,
    recipient_user_id INT,
    subject VARCHAR(255),
    body TEXT,
    date_sent DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (sender_user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Table for storing email attachments (in case you want to handle attachments)
CREATE TABLE IF NOT EXISTS EmailAttachments (
    attachment_id INT AUTO_INCREMENT PRIMARY KEY,
    email_id INT,
    file_name VARCHAR(255),
    file_path VARCHAR(255),
    file_size INT,
    FOREIGN KEY (email_id) REFERENCES Emails(email_id) ON DELETE CASCADE
);

-- Table for storing user calendar events (Optional)
CREATE TABLE IF NOT EXISTS CalendarEvents (
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    event_title VARCHAR(255),
    event_description TEXT,
    event_start DATETIME,
    event_end DATETIME,
    location VARCHAR(255),
    date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);
*/
--Create Db version postgreSQL
-- Creating a database for the Office1789

/* Je retire ceci à cause du docker compose CREATE DATABASE Office1789;

-- Connecting to the database
\c office1789;*/

-- Table for storing organizations (for Enterprise accounts with sub-accounts)
-- MUST be created BEFORE Users table due to foreign key
CREATE TABLE IF NOT EXISTS Organizations (
    organization_id SERIAL PRIMARY KEY,
    organization_name VARCHAR(255) NOT NULL,
    owner_user_id INT, -- Will be set after user creation
    max_members INT DEFAULT 1, -- Maximum number of sub-accounts allowed
    custom_domain VARCHAR(255), -- Domaine personnalisé (ex: company.com)
    domain_verified BOOLEAN DEFAULT FALSE, -- Domaine vérifié via DNS
    domain_verification_token VARCHAR(255), -- Token pour vérification DNS TXT
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing users
CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    domain VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL, -- Email du compte mail (username@office1789.com)
    recovery_email VARCHAR(255), -- Email de récupération personnel (optionnel)
    recovery_email_verified BOOLEAN DEFAULT FALSE, -- Email vérifié de façon permanente
    phonenumber VARCHAR(255),
    phonenumber_verified BOOLEAN DEFAULT FALSE, -- Téléphone vérifié de façon permanente
    password_hash VARCHAR(255) NOT NULL,
    mail_password VARCHAR(255), -- Mot de passe chiffré AES-256-GCM pour Mail/Matrix
    nboffer INT,
    role VARCHAR(20) DEFAULT 'user', -- 'user' ou 'admin'
    stripe_customer_id VARCHAR(255), -- ID du client Stripe
    stripe_subscription_id VARCHAR(255), -- ID de l'abonnement Stripe actif
    organization_id INT, -- ID de l'organisation (si membre d'une organisation)
    parent_account_id INT, -- ID du compte parent (si sous-compte)
    account_type VARCHAR(20) DEFAULT 'personal', -- 'personal', 'organization_owner', 'organization_member'
    custom_domain VARCHAR(255), -- Domaine personnalisé pour compte individuel Pro/Enterprise
    domain_verified BOOLEAN DEFAULT FALSE, -- Domaine vérifié
    date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES Organizations(organization_id) ON DELETE SET NULL,
    FOREIGN KEY (parent_account_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Add foreign key constraint for Organizations.owner_user_id after Users table is created
ALTER TABLE Organizations ADD CONSTRAINT fk_owner_user 
    FOREIGN KEY (owner_user_id) REFERENCES Users(user_id) ON DELETE CASCADE;

-- Table for storing user drive data (files)
CREATE TABLE IF NOT EXISTS DriveFiles (
    file_id SERIAL PRIMARY KEY,
    user_id INT,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    file_size INT NOT NULL, -- Size in bytes
    file_type VARCHAR(255),
    date_uploaded TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS SharedFiles (
    share_id SERIAL PRIMARY KEY,
    file_id INT REFERENCES DriveFiles(file_id) ON DELETE CASCADE,
    shared_with_user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    shared_by_user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    permission VARCHAR(20) NOT NULL DEFAULT 'editor',
    active BOOLEAN DEFAULT true,
    date_shared TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(file_id, shared_with_user_id)
);

-- Création de la table Groups avec les membres intégrés
CREATE TABLE IF NOT EXISTS Groups (
    group_id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    member_id INT REFERENCES Users(user_id),
    UNIQUE (group_id, member_id)
);




-- Création de la table TypingIndicator
CREATE TABLE IF NOT EXISTS TypingIndicator (
    typing_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    group_id INT REFERENCES Groups(group_id),
    last_typing_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing mail data
CREATE TABLE IF NOT EXISTS Emails (
    email_id SERIAL PRIMARY KEY,
    sender_user_id INT,
    recipient_user_id INT,
    subject VARCHAR(255),
    body TEXT,
    date_sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (sender_user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Table for storing email attachments
CREATE TABLE IF NOT EXISTS EmailAttachments (
    attachment_id SERIAL PRIMARY KEY,
    email_id INT,
    file_name VARCHAR(255),
    file_path VARCHAR(255),
    file_size INT,
    FOREIGN KEY (email_id) REFERENCES Emails(email_id) ON DELETE CASCADE
);

-- Table for storing participants
CREATE TABLE IF NOT EXISTS participants (
    group_id INT,
    user_id INT,
    last_read TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE
);

-- Table for storing message
CREATE TABLE IF NOT EXISTS Message (
    message_id SERIAL PRIMARY KEY,
    group_id INT,
    sender_id INT,
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE
);




-- Table for storing user calendar events
CREATE TABLE IF NOT EXISTS CalendarEvents (
    event_id SERIAL PRIMARY KEY,
    user_id INT,
    event_title VARCHAR(255),
    event_description TEXT,
    event_start TIMESTAMP,
    event_end TIMESTAMP,
    location VARCHAR(255),
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Table for storing user sessions with persistent storage
CREATE TABLE IF NOT EXISTS sessions (
    session_token VARCHAR(255) PRIMARY KEY,
    user_id INT NOT NULL,
    username VARCHAR(255) NOT NULL,
    expiry TIMESTAMP NOT NULL,
    password_plain TEXT,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Index pour améliorer les performances de recherche et nettoyage
CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions(expiry);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_username ON sessions(username);

-- Table for storing user 2FA/TOTP settings
CREATE TABLE IF NOT EXISTS user_totp (
    user_id INT PRIMARY KEY,
    secret VARCHAR(255) NOT NULL,
    enabled BOOLEAN DEFAULT FALSE,
    backup_codes TEXT[], -- Array of backup codes (hashed)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Table for storing password reset tokens
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    token VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Indexes for password reset tokens
CREATE INDEX IF NOT EXISTS idx_reset_token ON password_reset_tokens(token);
CREATE INDEX IF NOT EXISTS idx_reset_user ON password_reset_tokens(user_id);

-- Table for storing verification codes (email/phone)
CREATE TABLE IF NOT EXISTS verification_codes (
    id SERIAL PRIMARY KEY,
    contact VARCHAR(255) NOT NULL, -- Email ou numéro de téléphone
    code VARCHAR(6) NOT NULL, -- Code à 6 chiffres
    type VARCHAR(10) NOT NULL, -- 'email' ou 'phone'
    verified BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for verification codes
CREATE INDEX IF NOT EXISTS idx_verification_contact ON verification_codes(contact);
CREATE INDEX IF NOT EXISTS idx_verification_code ON verification_codes(code);

CREATE INDEX IF NOT EXISTS idx_user_totp_enabled ON user_totp(enabled);