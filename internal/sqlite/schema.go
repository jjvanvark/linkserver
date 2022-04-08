package sqlite

const dbSchema string = `

	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME DEFAULT NULL,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		password BLOB NOT NULL,
		key TEXT DEFAULT NULL
	);

	CREATE TABLE IF NOT EXISTS asset (
		id INTEGER PRIMARY KEY,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME DEFAULT NULL,
		cursor TEXT NOT NULL UNIQUE,
		file BLOB NOT NULL
	);

	CREATE TABLE IF NOT EXISTS link (
		id INTEGER PRIMARY KEY,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME DEFAULT NULL,
		cursor TEXT NOT NULL UNIQUE,
		user_id INTEGER NOT NULL,
		url TEXT NOT NULL UNIQUE,
		host TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES user(id)
	);

	CREATE TABLE IF NOT EXISTS article (
		id INTEGER PRIMARY KEY,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME DEFAULT NULL,
		link_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		by_line TEXT NOT NULL,
		content TEXT NOT NULL,
		text_content TEXT NOT NULL,
		length INTEGER NOT NULL,
		excerpt TEXT NOT NULL,
		site_name TEXT NOT NULL,
		image TEXT DEFAULT NULL,
		favicon TEXT DEFAULT NULL,
		FOREIGN KEY(link_id) REFERENCES link(id),
		FOREIGN KEY(image) REFERENCES asset(cursor),
		FOREIGN KEY(favicon) REFERENCES asset(cursor)
	);
	
`
