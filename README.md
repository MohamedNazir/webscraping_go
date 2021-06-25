**webscraping_go**

   webscraping_go is a go lang based web site scrapper to scrape through the given url and display the meta basic information about the url.

**To Run as docker container**

1. Install and start your docker application.
2. download/clone the source code
3. go to home directory
4. Run the below command where you can find the Dockerfile.
   " docker build . -t go-webscraper "
5. After build success Run the below command.
   " docker run -p 8080:8080 -it go-webscraper "
6. Open the web browser and hit the url "https://localhost:8080".

**To Build and run without docker.**

1. make sure to install latest gol
2. download/clone the source code
3. go to cmd directory, where you can find main.go files
4. run the command as "go run main.go"


**Features**
1. Fully tested 
2. web based html templating application
3. 
