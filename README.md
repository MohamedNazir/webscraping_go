**webscraping_go**

   webscraping_go is a go lang based web site scrapper to scrape through the given url and display the meta basic information about the url.

**To Run as docker container**

1. Start your docker client.
2. Download/clone the source code.
3. Go to application root directory.
4. Run the command where you can find the Dockerfile.
   
   **docker build . -t go-webscraper**
5. After build success Run the command.
   
   **docker run -p 8080:8080 -d go-webscraper**
6. Open the web browser and hit the url "http://localhost:8080".



**To Build and run without docker.**

1. Make sure to install latest go.
2. Download/clone the source code
3. go to cmd directory inside the root of the application, where you can find main.go files
4. Run the command as "go run main.go"
5. Open the web browser and hit the url "http://localhost:8080".



**Project Struture**

      webscraping_go/
        
         ├──  app              # Contains app.go for server mux routing of request
         ├──  asset            # Contains the HTML files for view rendering
         ├──  cmd              # main.go file
         ├──  domain           # Contains the domains struct, in this web app it contains response format struct.
         ├──  handler          # Contains handler files.
         ├──  service          # Contains service files where core business logic resides.
         ├──  utils/mocks      # Contains a mock web client used for test cases.
         ├──  webclient        # Contains web client to access the http end point.







**Features**
1. Fully tested 
2. web based html template application.
