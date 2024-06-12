CREATE TABLE IF NOT EXISTS todo_lists (
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    description varchar(255)
);

CREATE TABLE IF NOT EXISTS todo_items (
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    description varchar(255),
    done boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS lists_items (
    id serial PRIMARY KEY,
    item_id integer REFERENCES todo_items(id) ON DELETE CASCADE NOT NULL,
    list_id integer REFERENCES todo_lists(id) ON DELETE CASCADE NOT NULL
);