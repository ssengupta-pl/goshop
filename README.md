# goshop
A Go-based shopping list organizer.

# Database
At the root of the project a file named `.dsn` should be created. The file should
look like:
```
{
    "db": {
        "username": "<your_username>",
        "password": "<your_password>",
        "dbname": "goshopdb",
        "host": "localhost",
        "port": "5432",
        "sslmode": "disable"
    }
}
```

# Data Model
- Shopping List object
  - Name
  - Created On
  - Modified On
  - Creator
  - Collection of Item objects

- Item object
  - Name
  - Quantity
  - Unit of measurement (uom)
  - 1:1 mapping of a Store object

- Store object
  - Name
  - Address

# Services
- Shopping List
  - POST    /list       Create a new shopping list
  - GET     /list       Get a list of all shopping lists
  - GET     /list/$id   Get shopping list $listid
  - PUT     /list/$id   Update shopping $listid
  - DELETE  /list/$id   Delete shopping $listid

- Item
  - POST    /list/$listid/item          Create a new item in list $listid
  - GET     /list/$listid/item          Return all items in list $listid
  - GET     /list/$listid/item/$itemid  Return item $itemid in list $listid
  - PUT     /list/$listid/item/$itemid  Update item $itemid in list $listid
  - DELETE  /list/$listid/item          Delete all items in list $listid
  - DELETE  /list/$listid/item/$itemid  Delete item $itemid in list $listid

- Store
  - POST    /store          Create a new store
  - GET     /store          Get a list of all stores
  - DELETE  /store          Delete all stores
  - DELETE  /store/$storeid Delete store $storeid