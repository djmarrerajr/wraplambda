idempotency.AtMostOnce(<handler-func>, <table-name>, <msg-attr>)


- pull <msg-attr> from event data
- check <table-name> to see if a record exists with pk == <msg-attr>
    - if no, add the entry and invoke <handler-func>
    - if yes do not invoke <handler-func>
