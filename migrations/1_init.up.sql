# create schema if not exists task_runner;
create table if not exists workers
(
    uid            varchar(36)  not null primary key,
    name           varchar(150) null,
    last_connected datetime     null,
    status         tinyint      null
);

create table if not exists tasks
(
    uid          varchar(36)  not null primary key,
    status       varchar(50)  null,
    date_create  datetime     null,
    date_update  datetime     null,
    result       tinytext     null,
    subject      varchar(100) null,
    worker_uid   varchar(36)  null,
    retention    int    not null default 2592000, # 30 дней хранение по дефолту
    completed_at int    null,
    timeout      int          null,
    payload      text         null,
    error_msg    text         null,
    constraint tasks_workers_uid_fk foreign key (worker_uid) references workers (uid)
);