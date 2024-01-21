create table if not exists customers
(
    user_id      varchar(255) constraint customers_pk primary key,
    login        varchar(255) not null,
    password     varchar(255) not null,
    name         varchar(255) not null,
    company_id   integer      not null
    constraint customers_customer_companies_company_id_fk
    references customer_companies,
    credit_cards TEXT[]
);

insert into customers (user_id, login,password,  name, company_id, credit_cards)
values ('ivan', 'ivan', '12345', 'Ivan Ivanovich', 1,  ARRAY ['*****-1234', '*****-5678']),
       ('petr', 'petr', '54321', 'Petr Petrovich', 2,  ARRAY ['*****-4321', '*****-8765']);

