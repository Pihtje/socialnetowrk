-- +migrate Up

CREATE TABLE optional_userdata (
    "user_id" INTEGER NOT NULL,
    "image_URL" TEXT,
    "nickname" TEXT,
    "about_me" TEXT,
    FOREIGN KEY ("user_id") REFERENCES users("user_id")
    FOREIGN KEY ("image_URL") REFERENCES images("image_URL")
);

-- +migrate Down

DROP TABLE optional_userdata;