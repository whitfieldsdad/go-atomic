# DuckDB queries

- [ ] Search for tasks by step ID

## Search for tasks by step ID

```sql
SELECT 
    DISTINCT(task_id) 
FROM (
    SELECT 
        id AS task_id, 
        UNNEST(steps) AS step 
    FROM '*.json'
) WHERE step.id = '80a6c981-37b2-545b-b5e4-edaf360ba65d';
```
