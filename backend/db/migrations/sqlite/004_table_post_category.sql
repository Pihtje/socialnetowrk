-- +migrate Up

CREATE TABLE "post_category" (
    "post_id" INTEGER NOT NULL,
    "category_id" INTEGER NOT NULL,
    FOREIGN KEY ("post_id") REFERENCES "posts"("post_id")			
    FOREIGN KEY ("category_id") REFERENCES "categories"("category_id")
);

-- +migrate Down
DROP TABLE post_category;