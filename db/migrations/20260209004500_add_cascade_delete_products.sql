-- +goose Up
ALTER TABLE product_images DROP CONSTRAINT fk_product_images_id_products_id;
ALTER TABLE product_images ADD CONSTRAINT fk_product_images_id_products_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE products_category DROP CONSTRAINT fk_products_category_product_id_products_id;
ALTER TABLE products_category ADD CONSTRAINT fk_products_category_product_id_products_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE product_images DROP CONSTRAINT fk_product_images_id_products_id;
ALTER TABLE product_images ADD CONSTRAINT fk_product_images_id_products_id FOREIGN KEY (product_id) REFERENCES products(id);

ALTER TABLE products_category DROP CONSTRAINT fk_products_category_product_id_products_id;
ALTER TABLE products_category ADD CONSTRAINT fk_products_category_product_id_products_id FOREIGN KEY (product_id) REFERENCES products(id);
