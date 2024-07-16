BEGIN;

CREATE TABLE IF NOT EXISTS users (
	userID INT GENERATED ALWAYS AS IDENTITY,
	discordUserID VARCHAR(20) NOT NULL,
	discordGuildID VARCHAR(20) NOT NULL,
	isMod BOOL NOT NULL,
	PRIMARY KEY(userID),
	UNIQUE(discordUserID, discordGuildID)
);

CREATE INDEX IF NOT exists index_discordUserID ON users(discordUserID);


create table if not exists strikes (
	strikeID INT GENERATED ALWAYS AS IDENTITY,
	userID INT NOT null,
	reason text,
	CONSTRAINT fk_user FOREIGN KEY(userID) REFERENCES users(userID)
);

CREATE TABLE if not exists exiles (
	exileID INT GENERATED ALWAYS AS IDENTITY,
	userID INT NOT null,
	reason TEXT,
	startTimestamp TIMESTAMP,
	endTimestamp TIMESTAMP,
	exileStatus INT NOT NULL,
	PRIMARY KEY(exileID), 
	CONSTRAINT fk_user FOREIGN KEY(userID) REFERENCES users(userID)
);

COMMIT;	