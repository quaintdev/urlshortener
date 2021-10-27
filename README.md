# URL Shortener

An URL shortener api implementation

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
