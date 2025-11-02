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

-- Table for storing users
CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    domain VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    nboffer INT,
    phonenumber VARCHAR(255),
    date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

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