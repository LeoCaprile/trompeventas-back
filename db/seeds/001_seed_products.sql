-- +goose Up
INSERT INTO "public"."categories" ("name") VALUES
('trumpets'),
('mouthpieces'),
('saxofon'),
('trombones'),
('case');

INSERT INTO "public"."products" ("name", "description", "price") VALUES
('Yamaha YTR-2330 Standard Bb Trumpet', 'Student trumpet with excellent intonation and durability. Yellow brass bell with clear lacquer finish.', 450000),
('Bach Stradivarius 180S37 Professional Trumpet', 'Professional trumpet with .459" bore, 37 bell, and silver plated finish. Industry standard for professionals.', 2850000),
('Jupiter JTR700 Performance Series Trumpet', 'Intermediate trumpet with rose brass leadpipe and stainless steel pistons. Great for advancing students.', 680000),
('Conn Vintage One Trumpet', 'Professional level trumpet with hand-hammered one-piece bell and vintage design. Rich, warm tone.', 1950000),
('Getzen 590S-S Capri Bb Trumpet', 'Silver-plated professional trumpet with shepherd''s crook bell design. Brilliant sound and response.', 1750000),

-- MOUTHPIECES (7 products - various instruments)
('Bach 7C Trumpet Mouthpiece', 'The most popular trumpet mouthpiece. Medium cup depth with standard rim. Ideal for all-around playing.', 45000),
('Schilke 14A4a Trumpet Mouthpiece', 'Professional trumpet mouthpiece with shallow cup. Excellent for lead and high register playing.', 58000),
('Denis Wick 5880-5AL Trombone Mouthpiece', 'Large bore trombone mouthpiece with deep cup. Rich, full tone for orchestral and solo work.', 62000),
('Bach 6-1/2AL Trombone Mouthpiece', 'Medium-large cup trombone mouthpiece. Standard choice for symphonic trombone players.', 48000),
('Yamaha 11C4-7C Trumpet Mouthpiece', 'Japanese-made trumpet mouthpiece with comfortable rim and medium-shallow cup.', 38000),
('Monette B2 S3 Trumpet Mouthpiece', 'High-end custom trumpet mouthpiece. Exceptional response and intonation throughout all registers.', 285000),
('Greg Black NY1 Trombone Mouthpiece', 'Modern design trombone mouthpiece with excellent projection and flexibility. Popular among jazz players.', 95000),

-- SAXOPHONES (4 products)
('Yamaha YAS-280 Alto Saxophone', 'Student alto saxophone with improved ergonomics and key layout. Durable and easy to play.', 1250000),
('Selmer Paris Serie III Jubilee Alto Saxophone', 'Professional alto saxophone with incredible tonal color and projection. Gold lacquer finish.', 6800000),
('Yanagisawa T-WO2 Tenor Saxophone', 'Professional tenor saxophone with bronze body. Warm, powerful sound with excellent intonation.', 5950000),
('Jupiter JAS710GNA Alto Saxophone', 'Intermediate alto saxophone with gold lacquer finish. Responsive action and rich tone.', 1580000),

-- TROMBONES (4 products)
('Yamaha YSL-354 Tenor Trombone', 'Student tenor trombone with F attachment. Lightweight and balanced design for beginners.', 1150000),
('Bach 42BO Stradivarius Trombone', 'Professional tenor trombone with open wrap F attachment. .547" bore with 8.5" bell.', 3250000),
('Conn 88HO Symphony Series Trombone', 'Professional tenor trombone with open wrap. Rich orchestral sound with excellent response.', 2980000),
('Edwards T350-A Tenor Trombone', 'Hand-crafted professional trombone with Alessi valve. Premium materials and construction.', 4500000),

-- CASES (5 products - for different instruments)
('ProTec PRO PAC Case for Trumpet', 'Lightweight contoured case with plush interior, exterior music pocket, and adjustable backpack straps.', 185000),
('BAM Classic Trumpet Case', 'High-quality French-made case with cushioned suspension system. Available in multiple colors.', 425000),
('Protec PB304CT Trombone Case', 'Heavy-duty contoured trombone case with wheels. Separate mouthpiece compartment and accessory pouch.', 295000),
('SKB Alto Saxophone Case', 'Hardshell molded case with TSA-approved locking latches. Maximum protection for your saxophone.', 380000),
('Gator GL-TENOR-SAX Tenor Saxophone Case', 'Lightweight rigid EPS foam case with plush interior. Comfortable carrying handle and shoulder strap.', 220000);


INSERT INTO "public"."products_category" ("product_id", "category_id") VALUES
((SELECT id FROM products WHERE name = 'Yamaha YTR-2330 Standard Bb Trumpet'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'Bach Stradivarius 180S37 Professional Trumpet'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'Jupiter JTR700 Performance Series Trumpet'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'Conn Vintage One Trumpet'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'Getzen 590S-S Capri Bb Trumpet'), (SELECT id FROM categories WHERE name = 'trumpets')),

-- Mouthpieces (mouthpiece + instrument category)
-- Trumpet mouthpieces
((SELECT id FROM products WHERE name = 'Bach 7C Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'mouthpieces')),
((SELECT id FROM products WHERE name = 'Bach 7C Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'Schilke 14A4a Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'mouthpieces')),
((SELECT id FROM products WHERE name = 'Schilke 14A4a Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'Yamaha 11C4-7C Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'mouthpieces')),
((SELECT id FROM products WHERE name = 'Yamaha 11C4-7C Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'Monette B2 S3 Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'mouthpieces')),
((SELECT id FROM products WHERE name = 'Monette B2 S3 Trumpet Mouthpiece'), (SELECT id FROM categories WHERE name = 'trumpets')),

-- Trombone mouthpieces
((SELECT id FROM products WHERE name = 'Denis Wick 5880-5AL Trombone Mouthpiece'), (SELECT id FROM categories WHERE name = 'mouthpieces')),
((SELECT id FROM products WHERE name = 'Denis Wick 5880-5AL Trombone Mouthpiece'), (SELECT id FROM categories WHERE name = 'trombones')),
((SELECT id FROM products WHERE name = 'Bach 6-1/2AL Trombone Mouthpiece'), (SELECT id FROM categories WHERE name = 'mouthpieces')),
((SELECT id FROM products WHERE name = 'Bach 6-1/2AL Trombone Mouthpiece'), (SELECT id FROM categories WHERE name = 'trombones')),
((SELECT id FROM products WHERE name = 'Greg Black NY1 Trombone Mouthpiece'), (SELECT id FROM categories WHERE name = 'mouthpieces')),
((SELECT id FROM products WHERE name = 'Greg Black NY1 Trombone Mouthpiece'), (SELECT id FROM categories WHERE name = 'trombones')),

-- Saxophones (only saxofon category)
((SELECT id FROM products WHERE name = 'Yamaha YAS-280 Alto Saxophone'), (SELECT id FROM categories WHERE name = 'saxofon')),
((SELECT id FROM products WHERE name = 'Selmer Paris Serie III Jubilee Alto Saxophone'), (SELECT id FROM categories WHERE name = 'saxofon')),
((SELECT id FROM products WHERE name = 'Yanagisawa T-WO2 Tenor Saxophone'), (SELECT id FROM categories WHERE name = 'saxofon')),
((SELECT id FROM products WHERE name = 'Jupiter JAS710GNA Alto Saxophone'), (SELECT id FROM categories WHERE name = 'saxofon')),

-- Trombones (only trombone category)
((SELECT id FROM products WHERE name = 'Yamaha YSL-354 Tenor Trombone'), (SELECT id FROM categories WHERE name = 'trombones')),
((SELECT id FROM products WHERE name = 'Bach 42BO Stradivarius Trombone'), (SELECT id FROM categories WHERE name = 'trombones')),
((SELECT id FROM products WHERE name = 'Conn 88HO Symphony Series Trombone'), (SELECT id FROM categories WHERE name = 'trombones')),
((SELECT id FROM products WHERE name = 'Edwards T350-A Tenor Trombone'), (SELECT id FROM categories WHERE name = 'trombones')),

-- Cases (case + instrument category)
-- Trumpet cases
((SELECT id FROM products WHERE name = 'ProTec PRO PAC Case for Trumpet'), (SELECT id FROM categories WHERE name = 'case')),
((SELECT id FROM products WHERE name = 'ProTec PRO PAC Case for Trumpet'), (SELECT id FROM categories WHERE name = 'trumpets')),
((SELECT id FROM products WHERE name = 'BAM Classic Trumpet Case'), (SELECT id FROM categories WHERE name = 'case')),
((SELECT id FROM products WHERE name = 'BAM Classic Trumpet Case'), (SELECT id FROM categories WHERE name = 'trumpets')),

-- Trombone case
((SELECT id FROM products WHERE name = 'Protec PB304CT Trombone Case'), (SELECT id FROM categories WHERE name = 'case')),
((SELECT id FROM products WHERE name = 'Protec PB304CT Trombone Case'), (SELECT id FROM categories WHERE name = 'trombones')),

-- Saxophone cases
((SELECT id FROM products WHERE name = 'SKB Alto Saxophone Case'), (SELECT id FROM categories WHERE name = 'case')),
((SELECT id FROM products WHERE name = 'SKB Alto Saxophone Case'), (SELECT id FROM categories WHERE name = 'saxofon')),
((SELECT id FROM products WHERE name = 'Gator GL-TENOR-SAX Tenor Saxophone Case'), (SELECT id FROM categories WHERE name = 'case')),
((SELECT id FROM products WHERE name = 'Gator GL-TENOR-SAX Tenor Saxophone Case'), (SELECT id FROM categories WHERE name = 'saxofon'));


INSERT INTO "public"."product_images" ("product_id", "image_url") VALUES
((SELECT id FROM products WHERE name = 'Yamaha YTR-2330 Standard Bb Trumpet'), 'https://media.sweetwater.com/m/products/image/ac23c66f9emramJaWtP4LviHC0FM3WODZzsha6u4.jpg'),
((SELECT id FROM products WHERE name = 'Bach Stradivarius 180S37 Professional Trumpet'), 'https://media.sweetwater.com/m/products/image/08cd8a1d841UIO1ya8qLbiDnoEB7jzbndSee37VU.jpg?quality=82&width=750&ha=08cd8a1d84e4cc53'),
((SELECT id FROM products WHERE name = 'Bach Stradivarius 180S37 Professional Trumpet'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQlKMJGkPjlwKRCL6tmXAfP-R6rxWbfGoYthA&s'),
((SELECT id FROM products WHERE name = 'Jupiter JTR700 Performance Series Trumpet'), 'https://i.ebayimg.com/images/g/vusAAOSwpchbBVbN/s-l1200.jpg'),
((SELECT id FROM products WHERE name = 'Conn Vintage One Trumpet'), 'https://valveandreed.co.uk/cdn/shop/products/IMG_7289_4000x@3x.progressive.jpg?v=1619447035'),
((SELECT id FROM products WHERE name = 'Getzen 590S-S Capri Bb Trumpet'), 'https://preview.redd.it/anyone-have-experience-with-a-getzen-capri-590s-v0-7etgzx0s0rcb1.jpg?width=1127&format=pjpg&auto=webp&s=ce877bf3bcfc5b9684df87c65473825f9ee65daf'),

-- Mouthpiece Images (Product IDs 6-12)
((SELECT id FROM products WHERE name = 'Bach 7C Trumpet Mouthpiece'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSloiAntUDb815cpTUgO644eTOv4quW-_KMhw&s'),
((SELECT id FROM products WHERE name = 'Schilke 14A4a Trumpet Mouthpiece'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQs9SjcW2ZgE_hsw7QstPTDJhFXCDMVRdltQw&s'),
((SELECT id FROM products WHERE name = 'Denis Wick 5880-5AL Trombone Mouthpiece'), 'https://m.media-amazon.com/images/I/71d33OvLvUL.jpg'),
((SELECT id FROM products WHERE name = 'Bach 6-1/2AL Trombone Mouthpiece'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSloiAntUDb815cpTUgO644eTOv4quW-_KMhw&s'),
((SELECT id FROM products WHERE name = 'Yamaha 11C4-7C Trumpet Mouthpiece'), 'https://www.adams-music.com/images/productpicture/5B/MS/CN/5BMSCNYM11C4_1_1024.jpg'),
((SELECT id FROM products WHERE name = 'Monette B2 S3 Trumpet Mouthpiece'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSj41Q7MUJ3C16jyr-xwpxrXry1MYbpGSkjDA&s'),
((SELECT id FROM products WHERE name = 'Greg Black NY1 Trombone Mouthpiece'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQH_Gj8DqppVgO1CUb_i5wXOpKwFiLxd8w_NQ&s'),

-- Saxophone Images (Product IDs 13-16)
((SELECT id FROM products WHERE name = 'Yamaha YAS-280 Alto Saxophone'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSpGCHmLjGLEqKviJSXeErZkENO5JMv6Ex27A&s'),
((SELECT id FROM products WHERE name = 'Yamaha YAS-280 Alto Saxophone'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSyHA9kLEZZ4qoKVTS_JqjkSLREj76IZNKrAQ&s'),
((SELECT id FROM products WHERE name = 'Selmer Paris Serie III Jubilee Alto Saxophone'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQbMoAini7urE-74Rd6ysz4PIl6uKsB8iMIRg&s'),
((SELECT id FROM products WHERE name = 'Yanagisawa T-WO2 Tenor Saxophone'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT3YIq-d4eIG2E89YxB2SJ0HJXi2yJwctt8ZQ&s'),
((SELECT id FROM products WHERE name = 'Jupiter JAS710GNA Alto Saxophone'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR57YJQ2M5RP-8amS3tOEzGJBGxCCXMjX9waA&s'),

-- Trombone Images (Product IDs 17-20)
((SELECT id FROM products WHERE name = 'Yamaha YSL-354 Tenor Trombone'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSBui9kJc-hH2vc8pWK2Z9ywZ7g1x1-zfJCew&s'),
((SELECT id FROM products WHERE name = 'Bach 42BO Stradivarius Trombone'), 'https://i.ebayimg.com/images/g/dEcAAOSwnk1nEV7H/s-l1200.jpg'),
((SELECT id FROM products WHERE name = 'Bach 42BO Stradivarius Trombone'), 'https://media.musicarts.com/is/image/MMGS7/483626000998988-06-2000x2000.jpg'),
((SELECT id FROM products WHERE name = 'Conn 88HO Symphony Series Trombone'), 'https://www.palenmusic.com/cdn/shop/files/Untitleddesign-2024-06-21T141848.800.png?v=1718997552&width=1600'),
((SELECT id FROM products WHERE name = 'Edwards T350-A Tenor Trombone'), 'https://i.ebayimg.com/images/g/q9QAAOSwinNmaQQA/s-l1200.jpg'),

-- Case Images (Product IDs 21-25)
((SELECT id FROM products WHERE name = 'ProTec PRO PAC Case for Trumpet'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRnrU5Qr5ll3swDT2LJJKIV9rR9njXn1Q6w8Q&s'),
((SELECT id FROM products WHERE name = 'BAM Classic Trumpet Case'), 'https://media.sweetwater.com/m/products/image/a76e155430SPhYUKcI0cS6eKIawrzYvB2lZ9BqrW.jpg?quality=85&badge=original&width=300&height=300&bg-color=ffffff&ha=a76e155430e96d40&fit=bounds&canvas=300%2C300'),
((SELECT id FROM products WHERE name = 'Protec PB304CT Trombone Case'), 'https://www.adams-music.com/images/productpicture/5B/HE/SX/5BHESXAPTV_4_1024.jpg'),
((SELECT id FROM products WHERE name = 'SKB Alto Saxophone Case'), 'https://m.media-amazon.com/images/I/61Gs0gG5EoL.jpg'),
((SELECT id FROM products WHERE name = 'Gator GL-TENOR-SAX Tenor Saxophone Case'), 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQrBRfoJICpkQyoN1xOKEQfb2XJ4SAfuZdnQA&s');

-- +goose Down

DELETE FROM products_category;
DELETE FROM categories;
DELETE FROM products;
