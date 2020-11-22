create table t_users
(
	id int not null primary key identity(1000,1),
	name varchar(16) not null,
	phone char(11) not null,
	headimg varchar(256) not null,
	city int not null,
	brithday date not null,
	ismale bit not null,
	password char(32) not null,
	updatetime date not null default getdate(),
	writetime datetime not null default getdate()
)