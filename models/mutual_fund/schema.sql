create table mf_schemes (
	id int primary key,
	scheme_name varchar(512)
);

create table mf_investments (
	id serial primary key,
	scheme_id int,
	nav decimal(50, 4),
	units decimal(50, 3),
	invested_at date,
	created_at date
);

create table mf_nav_data (
	scheme_id int,
	nav_date date,
	nav decimal(50, 4),
	primary key (scheme_id, nav_date),
	created_at date
);