### Requirements

- Multiple people should be able to add expanses
    - Authentication (verifing who is accessing)
    - Autherization  (who can access what)         /admin/...

- Adding expanses
    (title, money, tag)     tag = '#zomato', '#blinkit', '#rent'
    /expanses/add - {title, money, tag} (user_id)
                  - no output

- How much did I spend
    (from, to) ---> spends, graph
    /expanses/list - {from, to} (user_id)
                   - []Spends


## spends
user_id, title, money, tag

## users
id, first_name, last_name, email_id, password