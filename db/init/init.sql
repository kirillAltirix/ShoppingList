CREATE TABLE users
(
    user_id SERIAL PRIMARY KEY,
    chat_id TEXT NOT NULL,
    username TEXT NOT NULL
);

CREATE TABLE lists
(
    list_id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE lists_owners
(
    list_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (list_id) REFERENCES lists (list_id)
);

CREATE TABLE lists_subowners
(
    list_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (list_id) REFERENCES lists (list_id)
);

CREATE TABLE items
(
    item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE lists_items
(
    list_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    FOREIGN KEY (list_id) REFERENCES lists (list_id),
    FOREIGN KEY (item_id) REFERENCES items (item_id)
);

CREATE TABLE lists_messages
(
    list_id INTEGER NOT NULL,
    message_id INTEGER NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (list_id) REFERENCES lists (list_id)
);

CREATE TABLE texts
(
    text_id SERIAL PRIMARY KEY,
    text TEXT NOT NULL
);