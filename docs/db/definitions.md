# Db Table definitions

## User management

users(__id__, username, email, password)

```sql
DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD email ON user TYPE string;
DEFINE FIELD password ON user TYPE string;
```

stops(__id__, name, location)

```sql
DEFINE TABLE stop SCHEMAFULL;
DEFINE FIELD name ON stop TYPE string;
DEFINE FIELD location ON stop TYPE array;
```

lines(__id__, name)

```sql
DEFINE TABLE line SCHEMAFULL;
DEFINE FIELD name ON line TYPE string;
```

users->likes_lines->lines
users->likes_stops->stops
lines->stop_at->stops

```sql
DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD email ON user TYPE string;
DEFINE FIELD password ON user TYPE string;
DEFINE FIELD image ON user TYPE object;
DEFINE FIELD image.name ON user TYPE string;
DEFINE FIELD image.type ON user TYPE string;
DEFINE TABLE stop SCHEMAFULL;
DEFINE FIELD name ON stop TYPE string;
DEFINE FIELD location ON stop TYPE array;
DEFINE TABLE line SCHEMAFULL;
DEFINE FIELD name ON line TYPE string;
```
