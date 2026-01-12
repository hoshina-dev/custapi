# Database Migrations

This directory contains SQL migration files following industry-standard naming conventions.

## Naming Convention

```
NNN_description.up.sql    # Migration forward
NNN_description.down.sql  # Migration rollback
```

Where:
- `NNN` is a 3-digit sequential number (001, 002, 003, etc.)
- `description` is a snake_case description of the migration
- `.up.sql` applies the migration
- `.down.sql` rolls back the migration

## Current Migrations

| # | Name | Description |
|---|------|-------------|
| 001 | create_organizations_table | Creates organizations table with UUID primary key, indexes, and triggers |
| 002 | create_users_table | Creates users table with foreign key to organizations, indexes, and triggers |

## Running Migrations

### Manual Execution

```bash
# Apply migration
psql -U postgres -d custapi -f internal/database/migrations/001_create_organizations_table.up.sql
psql -U postgres -d custapi -f internal/database/migrations/002_create_users_table.up.sql

# Rollback migration (in reverse order)
psql -U postgres -d custapi -f internal/database/migrations/002_create_users_table.down.sql
psql -U postgres -d custapi -f internal/database/migrations/001_create_organizations_table.down.sql
```

### Using Migration Tools (Optional)

For production, consider using migration tools like:
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [pressly/goose](https://github.com/pressly/goose)
- [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate)

## Migration Guidelines

1. **Never modify existing migrations** - Create new ones instead
2. **Always provide down migrations** - For rollback capability
3. **Test migrations** - On a copy of production data
4. **Keep migrations atomic** - One logical change per migration
5. **Version control** - All migrations should be in git

## Features Included

- UUID primary keys with automatic generation
- Soft deletes (deleted_at column)
- Automatic timestamp updates (updated_at trigger)
- Proper indexing for performance
- Foreign key constraints with cascade delete
