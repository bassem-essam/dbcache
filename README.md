# dbcache
A simple demo of an in-memory DB cache with very fast single-item insertion speed

# Goal
Considering the item as a resource that should have a unique value:
It allows for creating Items simply by using NewItem(...) with neither the fear of duplicating the item in the database, nor the slow insertion overhead.

# How it works
A map (hashmap) is used to represent the Cache

A go channel is used to stream the newly created items into a pile (array) of the same class.

Whenever the pile reaches a certain size, batch insertion of items is used for maximum speed.

# Technical reasons
In SQL databases, insertion is slower than retrieval. Moreover in ORM's, insertion of large number of records is much slower than batch insertion.

The used ORM (gorm in this case) has efficient solutions for batch insertion which needs.

# Why
Maybe you need to create an Item and insert it into the db, then you may want to further interact with this item by relating it into another table in some way or another. 

