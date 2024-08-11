## Initial CQL script

```sql
create keyspace dialogs_space
 with replication = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
 };
```

```sql
create table if not exists dialogs_space.messages(
 id UUID,
 author_uid UUID,
 receiver_uid UUID,
 message text,
 message_time timestamp,
 PRIMARY KEY((author_uid, receiver_uid), message_time, id)
) WITH CLUSTERING ORDER BY (message_time DESC);
```

# Useful queries

```sql

select * from dialogs_space.messages;
insert into dialogs_space.messages(id, author_uid, receiver_uid, message, message_time) values(uuid(), cc6691e2-d3b0-4271-bce1-f8387364aa0a, b452a917-fa9a-434d-b3db-b26ef6e5e8fe, 'Hello!', toTimestamp(now()));
insert into dialogs_space.messages(id, author_uid, receiver_uid, message, message_time) values(uuid(), b452a917-fa9a-434d-b3db-b26ef6e5e8fe, cc6691e2-d3b0-4271-bce1-f8387364aa0a, 'Hello!', toTimestamp(now()));

-- cheange keyspace
ALTER KEYSPACE dialogs_space
WITH REPLICATION = {
'class':'SimpleStrategy',
'replication_factor': 2
};

describe KEYSPACE dialogs_space;

-- the same token means that rows are stored on the same node
select token((author_uid, receiver_uid)), author_uid, receiver_uid, message, message_time from dialogs_space.messages;
```