create table users(
	id serial primary key,
	email varchar(64) not null,
	password varchar(72) not null
);

SELECT 'down SQL query';