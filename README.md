# URL Shortener

## Installation using Docker
1. Pull docker image of urlshortener on your local system using: 

    `docker pull quaintdev/urlshortener`

2. Run urlshortener using

    `docker run -p 3000:3000 urlshortener`

This will start server listening on port `3000`.  Server will be accessible at http://localhost:3000/.


## Usage

The application provides two endpoints
1. `/shorten` - accepts json POST request as shown below. It responds with additional `shortUrl` parameter that can be used to access original webpage.
    ```
    {
      "LongUrl":"https://news.ycombinator.com/news"
    }
    ```
2. `/`- accepts short url as GET request and redirects to original webpage.

Application provides another endpoint on port `3001`.  It allows to take backup of url store that are then loaded when the application starts.

3. `/backup` is used to take backups of the url store.


Note that with VSCode, you can use `requests.http` to test the application.
