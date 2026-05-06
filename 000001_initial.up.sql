CREATE TYPE user_role AS ENUM ('admin', 'member');

CREATE TABLE users (
  id            uuid primary key default gen_random_uuid(),
  name          varchar(100) not null,
  email         text not null,
  role          user_role default 'member',
  phone_no      text not null unique,
  password_hash text not null,
  created_at    timestamp default now(),
  archived_at   timestamptz
);

CREATE TABLE user_sessions (
  id          uuid primary key default gen_random_uuid(),
  emp_id      uuid references users(id),
  created_at  timestamptz default now(),
  archived_at timestamptz
);

CREATE UNIQUE INDEX idx_unique_email ON users (email) WHERE archived_at IS NULL;

CREATE TABLE projects (
  id          uuid primary key default gen_random_uuid(),
  name        varchar(100) not null,
  description text,
  owner_id    uuid references users(id) on delete set null,
  created_at  timestamptz default now(),
  archived_at timestamptz
);

CREATE TABLE project_members (
  project_id  uuid references projects(id) on delete cascade,
  user_id     uuid references users(id) on delete cascade,
  PRIMARY KEY (project_id, user_id)
);

CREATE TABLE tasks (
  id          uuid primary key default gen_random_uuid(),
  title       varchar(200) not null,
  description text,
  project_id  uuid references projects(id) on delete cascade,
  assignee_id uuid references users(id) on delete set null,
  status      varchar(20) default 'todo' check (status in ('todo', 'in_progress', 'done')),
  due_date    timestamptz,
  created_at  timestamptz default now(),
  archived_at timestamptz
);