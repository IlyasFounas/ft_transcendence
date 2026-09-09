# apiFetch

`apiFetch` is a shared helper for making HTTP requests to the backend.


## What it does

`apiFetch`:

- calls fetch()
- parses the JSON response

## Returns 
- the HTTP status and parsed response data
- returns null if the request fails completely (ex network failed)


## Create Options for the function's parameters
### GET example
```
const options: RequestInit = {
  method: "GET",
};

const result = await apiFetch("/api/friends", options);
```

### POST example
```
const options: RequestInit = {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    email,
    password,
  }),
};

const result = await apiFetch("/api/login", options);
```

## Result
If the backend responds successfully:
```
{
  status: 200,
  data: {
    // backend response
  }
}
```
 If the request completely fails ex. network error return null

## Handling the result
```
const result = await apiFetch("/api/login", options);

if (result === null) {
  // Network/request error
  return;
}
if (result.status === 401) {
  // Authentication error
  console.log(result.data.message);
  return;
}
if (result.status === 200) {
  // Handle successful response
}
```

#### Note:
apiFetch only handles communication with the backend. 
It does not decide what the UI should do.
The component decides how to handle the result