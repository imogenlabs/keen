# Sprints & Boards

## Boards

```bash
# List all boards
keen agile get-all-boards --agent

# Get board details
keen agile get-board 42 --agent

# Get board configuration
keen agile get-configuration 42 --agent

# Get board backlog
keen agile get-issues-for-backlog 42 --agent
```

## Sprints

```bash
# List sprints (active, future, closed)
keen agile get-all-sprints 42 --state active,future --agent

# Create sprint
keen agile create-sprint --stdin <<< '{"name":"Sprint 5","originBoardId":42}'

# Start sprint
keen agile partially-update-sprint 7 --stdin <<< '{"state":"active"}'

# Close sprint
keen agile partially-update-sprint 7 --stdin <<< '{"state":"closed"}'
```

## Moving Issues

```bash
# Move issues to sprint
keen agile move-issues-to-sprint-and-rank 7 --stdin <<< '{"issues":["PROJ-42","PROJ-43"]}'

# Move to backlog
keen agile move-issues-to-backlog --stdin <<< '{"issues":["PROJ-42"]}'

# Get sprint issues
keen agile get-issues-for-sprint 7 --agent
```
