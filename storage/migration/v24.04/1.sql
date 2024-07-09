create table if not exists logger.login_attempt
(
    id          serial,
    date        timestamp with time zone not null,
    domain_id   bigint,
    success     boolean                  not null,
    auth_type   text,
    user_id     bigint,
    user_ip     inet                     not null,
    user_agent  text,
    domain_name text                     not null,
    user_name   text                     not null,
    details     text
);

comment on table logger.login_attempt is 'stores all login attempts of the users';

alter table logger.login_attempt
    owner to opensips;

create index if not exists login_attempt_user_id_index
    on logger.login_attempt (user_id);

create index if not exists login_attempt_domain_id_index
    on logger.login_attempt (domain_id);

create index if not exists login_attempt_success_index
    on logger.login_attempt (success);

alter table logger.login_attempt
    add constraint login_attempt_pk
        primary key (id);

