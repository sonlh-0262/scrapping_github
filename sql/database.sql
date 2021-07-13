DROP TABLE IF EXISTS github;
CREATE TABLE github (
  id            INT AUTO_INCREMENT NOT NULL,
  owner         VARCHAR(255),
  name          VARCHAR(255),
  star          VARCHAR(255),
  fork          VARCHAR(255),
  branch_count  INT(11),
  tag_count     INT(11),
  PRIMARY KEY (`id`)
);

DROP TABLE IF EXISTS scrapping_parameters;
CREATE TABLE scrapping_parameters (
  id            INT AUTO_INCREMENT NOT NULL,
  url           VARCHAR(255),
  parameter     VARCHAR(8000),
  PRIMARY KEY (`id`)
);

INSERT INTO scrapping_parameters
  (url, parameter)
VALUES
  ('https://github.com/rails/rails', ''),
  ('https://github.com/toptal/chewy', '');
