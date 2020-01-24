# go-redfish-test
Test stubs for Nordix/go-redfish library

# Steps to use
* Clone the Redfish repo 
```
go get -d github.com/manojkva/go-redfish
cd <rep>
go install .
```
* Clone the test repo
 ```
 git clone https://github.com/manojkva/go-redfish-test.git
 
 ```
* Make changes in go.mod to specify to the local directory
  Replace the path after ```=>``` to specify the local repo
```
  replace github.com/Nordix/go-redfish/client => /home/ekuamaj/go/src/github.com/manojkva/go-redfish/client
 ```
 * Build the file using the following command within the test repo dir
  ```
  go run .
  ``` 
  
