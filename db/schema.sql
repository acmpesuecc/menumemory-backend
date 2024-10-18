create table if not exists User (
  id integer primary key autoincrement,
  email varchar(255) not null,
  FirstName varchar(255),
  Lastname varchar(255),
  DOB date
);

create table if not exists Restaurant (
    id integer primary key autoincrement,
    Name varchar(255) not null,
    Area varchar(255),
    Address varchar(2048),
    MapsLink varchar(512) unique,
    MapsRatingOutOf5 float
);

create table if not exists Dish (
    id integer primary key autoincrement,
    Name varchar(255) unique not null
);

create table if not exists Visit (
    id integer primary key autoincrement,
    Date date not null,
    Time time,
    UserId integer,
    RestaurantId integer,
    foreign key (UserId) references user(id),
    foreign key (RestaurantId) references restaurant(id)
);


create table if not exists Orders (
    id integer primary key autoincrement,
    VisitId integer,
    DishId integer,
    Rating float,
    ReviewText varchar(2048),
    foreign key (VisitId) references visit(id),
    foreign key (DishId) references dish(id)
);

