-- SQLite

PRAGMA FOREIGN_KEYS=OFF; 

CREATE TABLE recordSpace (
    id      INTEGER NOT NULL,
    owner   INTEGER NOT NULL,
    name    TEXT    NOT NULL    UNIQUE,
    type    INTEGER NOT NULL    DEFAULT 0,
    PRIMARY KEY(id)
);

CREATE TABLE rsGroupMapping (
    id      INTEGER NOT NULL,
    rs      INTEGER NOT NULL,
    gp      INTEGER NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(rs) REFERENCES recordSpace(id)
);

CREATE TABLE imageStore (
    priv    INTEGER NOT NULL,
    rs      INTEGER NOT NULL,   -- record space
    url     TEXT NOT NULL,      -- base64 url
    name    TEXT,
    PRIMARY KEY(priv)
    FOREIGN KEY(rs) REFERENCES recordSpace(id)
);

CREATE TABLE rsUserMapping (
    id      INTEGER NOT NULL,
    dst     INTEGER NOT NULL,   -- user that grant access to modify 
    rs      INTEGER NOT NULL,   -- target recordSpace
    PRIMARY KEY(id)
    FOREIGN KEY(rs) REFERENCES recordsSpace(id)
);


PRAGMA FOREIGN_KEYS=ON;


-- Views

CREATE VIEW opUserGetRS AS
SELECT owner, name, type, dst as qq 
FROM recordSpace 
INNER JOIN rsUserMapping 
ON recordSpace.id = rsUserMapping.rs;

CREATE VIEW opGroupGetRS AS
SELECT recordSpace.id as id, name, type, gp
FROM recordSpace
INNER JOIN rsGroupMapping
ON recordSpace.id = rsGroupMapping.rs;