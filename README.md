# Demo about how to handle large amount of parameters

This application represents the real case of a REST API which gets transactions from an ElasticSearch repository.

Here is the main path:
User call `/users/{id}/transactions` to get all his transactions. But with the time, he may have thousands of them. And for a frontend perspective, we need to paginate these transactions, even filter them (open vs paid, date range, etc...)

Unit tests and business logic in the reposity is removed in order to focus on the problem

## Issue

We have a struct defined for the http request model wanted. But the first intention was to separate the transport layer from the service and repository one.
Now we have multiple (too many) parameters in our methods. Do you know if there's a best practice in order to avoid that?
Use one parameter model per method in the service / repo ?
A generic "Option" model ?

Thanks for your help
