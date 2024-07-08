CREATE TABLE words(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    word TEXT NOT NULL UNIQUE,
    language TEXT NOT NULL,
    translation TEXT NOT NULL,
    example_sentence TEXT NOT NULL
);