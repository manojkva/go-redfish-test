package main
import (
        "fmt"
        redfish "github.com/Nordix/go-redfish/client"
        "context"
        "net/http"
        "crypto/tls"
        "time"
        "encoding/json"
//        "reflect"
         "os"
         "github.com/antihax/optional"
         "net/url"
         "regexp"
)
func prettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}
var tr *http.Transport = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
        TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
        }

func createAPIClient(HeaderInfo map[string]string) *redfish.DefaultApiService {
//       tr := &http.Transport{
//	MaxIdleConns:       10,
//	IdleConnTimeout:    30 * time.Second,
//	DisableCompression: true,
//        TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
//        }
        client := &http.Client{Transport: tr}
        cfg := &redfish.Configuration{
                BasePath:      "https://32.68.220.144",
                DefaultHeader: make(map[string]string),
                UserAgent:     "go-redfish/client",
                HTTPClient: client,
        }

        if len(HeaderInfo) != 0 {
        
        for key,value := range HeaderInfo {
                cfg.DefaultHeader[key] = value
        }
        } 
        return redfish.NewAPIClient(cfg).DefaultApi
}

func getTask( ctx context.Context, taskID string ) {
       redfishApi := createAPIClient(make(map[string]string)) 
       sl, response,err := redfishApi.GetTask(ctx,taskID )
       fmt.Printf( "%+v %+v %+v", prettyPrint(sl),response, err)
}

func getVirtualMedia( ctx context.Context ) {
       redfishApi := createAPIClient(make(map[string]string)) 
       sl, response,err := redfishApi.GetManagerVirtualMedia(ctx,"iDRAC.Embedded.1","CD" )
       fmt.Printf( "%+v %+v %+v", prettyPrint(sl),response, err)
}

func updateService( ctx context.Context ){
       redfishApi := createAPIClient(make(map[string]string)) 
       // call the UpdateService and get the HttpPushURi
       sl, response,err := redfishApi.UpdateService(ctx)
       fmt.Printf( "%+v %+v %+v", prettyPrint(sl),response, err)

}

func httpUriDownload( ctx context.Context, filePath string , etag string) (*url.URL,error ) {
        //filehandle, err  := os.Open("/media/sf_Downloads/BIOS_9P3C0_WN64_2.4.8.EXE")
        filehandle, err  := os.Open(filePath)
        if err != nil {
            fmt.Println(err)
           }
        defer filehandle.Close()
        reqBody :=  redfish.FirmwareInventoryDownloadImageOpts{  SoftwareImage :  optional.NewInterface(filehandle) }
        headerInfo := make(map[string]string)
        headerInfo["if-match"] = etag
        //redfishApi = redfish.NewAPIClient(cfg).DefaultApi
       redfishApi := createAPIClient(headerInfo) 

	sl,response,err := redfishApi.FirmwareInventoryDownloadImage(ctx,&reqBody )
        fmt.Printf( "%+v %+v %+v", prettyPrint(sl),response, err)
        return response.Location()

}


func  getETagHttpURI ( ctx context.Context ) string {
       redfishApi := createAPIClient(make(map[string]string)) 
       sl, response, err := redfishApi.FirmwareInventory(ctx)
       fmt.Printf( "%+v %+v %+v", prettyPrint(sl),response, err)
       etag :=  response.Header["Etag"]
       fmt.Printf("%v", etag[0])
       return etag[0]
}

func simpleUpdateRequest( ctx context.Context, imageURI string) string {
        headerInfo := make(map[string]string)
        headerInfo["content-type"] = "application/json"
       redfishApi := createAPIClient(headerInfo) 
//        reqBody := redfish.SimpleUpdateRequestBody{ ImageURI: "https://32.68.220.144/redfish/v1/UpdateService/FirmwareInventory/Available-159-2.4.8" ,
        reqBody := redfish.SimpleUpdateRequestBody{ ImageURI: imageURI , 
           }
	sl,response,err := redfishApi.UpdateServiceSimpleUpdate(ctx, reqBody,)
        fmt.Printf( "%+v %+v %+v", prettyPrint(sl),response, err)
        jobID_location := response.Header["Location"]
        re := regexp.MustCompile(`JID_(.*)`)
        jobID :=  re.FindStringSubmatch(jobID_location[0])[1]
        return jobID
}


func main(){

        var auth = redfish.BasicAuth { UserName: "root", 
                                       Password: "Abc.1234",
                                       }
       ctx  :=  context.WithValue( context.Background(), redfish.ContextBasicAuth , auth,)

       getTask(ctx,"JID_797173652449")
/*
       //getVirtualMedia(ctx)
      updateService(ctx )
   
     etag  := getETagHttpURI(ctx)
     fmt.Println("%v", etag)
     availableEntryInUri , _ := httpUriDownload(ctx, "/home/ekuamaj/workspace/nordtest/iDRAC-with-Lifecycle-Controller_Firmware_4JCPK_WN64_4.00.00.00_A00.EXE", etag )
   
     //  availableEntryInUri := "https://32.68.220.144/redfish/v1/UpdateService/FirmwareInventory/Available-25227-4.00.00.00"
      jobID := simpleUpdateRequest(ctx, availableEntryInUri.String())

      fmt.Println("%v", jobID)
      
      // simpleUpdateRequest(ctx, availableEntryInUri)

*/

        
       
}
