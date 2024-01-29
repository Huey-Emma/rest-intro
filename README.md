# What is a REST API?

REST stands for Representational State Transfer. It is a software architectural style created by Roy Fielding to guide the design of architecture for the web.

Any API that follows REST design principles is said to be RESTful. Simply put a RESTAPI is a medium for 2 computers to communicate over HTTP

## Some RESTful Design Principles Include

### Accept and respond in JSON

In the past, accepting and responding to API requests were done mostly in XML and in some cases HTML. These days however, JSON - which stands for JavaScript Object Notation is the standard for accepting and responding to API requests. Some of the reasons JSON has become the defactor format for interacting with APIs are its simplicity and support across multiple languages. To ensure that clients interpret JSON correctly, you should set the Content-Type in the Header response to application json `"{'Content-Type': 'application/json'}"` while making the request.

[Code sample](https://github.com/Huey-Emma/rest-intro/blob/main/d001)

### Use Nouns instead of Verbs in Endpoint Paths

RESTAPI endpoints should be written as nouns, not verbs because HTTP methods are verbs. The endpoint's intended action should be indicated by its HTTP request method.
HTTP request methods include:

-   GET retrieves resources
-   POST submits new data to the server
-   PUT updates existing data
-   PATCH partially updates existing data
-   DELETE removes data from the server

[Code sample](https://github.com/Huey-Emma/rest-intro/blob/main/d002)

### Use Nesting on Endpoints to Show Relationships

At times, different endpoints can be interlinked and therefore can be nested so they can be easily understood. For example, `https://www.sample.com/posts/23/comments` is an endpoint that is used to retrieve comments related to a certain post with id 23. As a side note, you should generally nesting that is more than 3 levels deep.

[Code sample](https://github.com/Huey-Emma/rest-intro/blob/main/d003)

### Handle Errors Gracefully and Return Standard Error Status Codes

When an error occurs while making a network call to an API endpoint, we should respond with appropriate status codes. The status codes give information to understand what the problem might be. Below is a table showing different HTTP status code ranges and their meaning

| Status Code Range | Meaning                                                                                  |
| ----------------- | ---------------------------------------------------------------------------------------- |
| 100-199           | Informational responses. For example, 102 indicates that the resource is being processed |
| 200-299           | Success. For example, 200 means OK; 201 means accepted                                   |
| 300-399           | Redirects. For example, 301 means moved permanently                                      |
| 400-499           | Client side errors. 400 means bad request; 404 means resource not found                  |
| 500-599           | Server side errors. For example, 500 means internal server error                         |

[Code sample](https://github.com/Huey-Emma/rest-intro/blob/main/d004)

### Use Filtering, Sorting and Pagination to Retrieve Requested Data

Sometimes an APIs database can get incredibly large, if this happens, retrieving data can become very slow. Filtering, sorting and pagination are actions that can be performed on large collections of data so it can be sent over the wire in small chunks and therefore is quicker

[Code sample](https://github.com/Huey-Emma/rest-intro/blob/main/d005)

### Cache data to improve performance

If you have a resource that is frequently requested for, you can improve your API latency by caching the resource to avoid querying your DB.

[Code sample](https://github.com/Huey-Emma/rest-intro/blob/main/d006)

### Versioning

API endpoints should have different versions to allow for backward compactibility. You typically do not want to force clients to migrate to a new version of your service abruptly.

[Code sample](https://github.com/Huey-Emma/rest-intro/blob/main/d007)
