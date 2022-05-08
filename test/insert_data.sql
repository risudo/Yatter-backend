INSERT INTO `attachment` (`type`, `url`) VALUES
('image', 'a.png');

SET @attachment_id = LAST_INSERT_ID();

INSERT INTO `account` (`username`, `password_hash`) VALUES
('a', '$2a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @a_id = LAST_INSERT_ID();

INSERT INTO `status` (`account_id`, `content`) VALUES
(@a_id, 'from a');

INSERT INTO `account` (`username`, `password_hash`) VALUES
('b', '$2a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @id = LAST_INSERT_ID();

INSERT INTO `status_contain_attachment` (`status_id`, `attachment_id`) VALUES
(@a_id, @attachment_id);

INSERT INTO `relation` (`following_id`, `follower_id`) VALUES
(@a_id, @id);

INSERT INTO `status` (`account_id`, `content`) VALUES
(@id, 'from b'),
(@id, 'from b');

INSERT INTO `account` (`username`, `password_hash`) VALUES
('c', '$2a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @id = LAST_INSERT_ID();

INSERT INTO `relation` (`following_id`, `follower_id`) VALUES
(@a_id, @id);

INSERT INTO `status` (`account_id`, `content`) VALUES
(@id, 'from c');

INSERT INTO `account` (`username`, `password_hash`) VALUES
('d', '$2a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @id = LAST_INSERT_ID();

INSERT INTO `relation` (`following_id`, `follower_id`) VALUES
(@a_id, @id);

INSERT INTO `account` (`username`, `password_hash`) VALUES
('e', '$4a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @e_id = LAST_INSERT_ID();

INSERT INTO `account` (`username`, `password_hash`) VALUES
('f', '$4a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO'),
('g', '$4a$10$T3C9WgYroD2SWAQegbB0qOzVC4XbqnWHHd9srL5DQ2ixbSj.Y4MDO');

SET @id = LAST_INSERT_ID();

INSERT INTO `relation` (`following_id`, `follower_id`) VALUES
(@a_id, @id),
(@e_id, @id);

INSERT INTO `status` (`account_id`, `content`) VALUES
(@id, 'from g'),
(@id, 'from g');
