# blog-to-rss-parser
Vibe-coded script that parses websites without RSS and generates a static RSS feed for the articles. Then, this RSS is hosted with [GitHub pages](https://fshchudlo.github.io/blog-to-rss-parser/feed.xml) and can be used to agrregate RSS.

## Development

1. **Install Dependencies**  
    Ensure you have the following installed on your system:
    - [Go](https://golang.org/dl/) (version 1.18 or later)
    - [Git](https://git-scm.com/)

2. **Set Up Your Environment**  
    - Clone the repository:
      ```bash
      git clone https://github.com/your-username/blog-to-rss-parser.git
      cd blog-to-rss-parser
      ```
    - Install Go modules:
      ```bash
      go mod tidy
      ```

3. **Run the Application**  
    - Start the application locally:
      ```bash
      go run main.go
      ```