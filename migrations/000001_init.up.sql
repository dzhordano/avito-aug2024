CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    user_type TEXT NOT NULL
);

CREATE TABLE houses (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    year INTEGER NOT NULL,
    developer TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE flats (
    id SERIAL PRIMARY KEY,
    flat_number INTEGER NOT NULL UNIQUE,
    price INTEGER NOT NULL,
    rooms INTEGER NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE house_flats(
    house_id INTEGER REFERENCES houses (id) ON DELETE CASCADE NOT NULL,
    flat_id INTEGER REFERENCES flats (id) ON DELETE CASCADE NOT NULL,
    flat_number INTEGER REFERENCES flats (flat_number) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (house_id, flat_number)
);

CREATE TABLE house_subscriptions (
    house_id INTEGER REFERENCES houses (id) ON DELETE CASCADE NOT NULL,
    user_email TEXT REFERENCES users (email) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (house_id, user_email)
);