## Register User:-
    
```http
  POST /api/reg
  
   JSON: {
    "username": "chetan",
    "email": "chetan.kumar@remotestate.com",
    "password": "abcd123"
}
```

## Login User:-
    
```http
  POST /api/login
  
  
   JSON: {
    "email": "chetan.kumar@remotestate.com",
    "password": "abcd123"
           }
```

#### After user Login User can create todos.


## User can Create Todo:-

```http
  POST /api/createTodo
  
   JSON: {
    "TodoName" : "test24",
    "TodoDescription" : "test23332 "
}
```


## User can Delete Todo By Id:-

```http
  POST /api/deleteById/12
```

## User can Find Todo By Id:-

```http
  POST /api/databyid/4
```




