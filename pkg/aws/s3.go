package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)
type Data struct {
	Path string
	Size int
}


// functions to connection with AWS S3
// create a function to create a new bucket

// UploadFile function to upload a File to a Bucket
func UploadFile(filename string,  path string) {
	// intialize bucket name
	bucket :="symbyo-test"
	// create new session on aws aws
	sess, err := session.NewSession(&aws.Config{
		// add region
		Region:      aws.String(os.Getenv("Reg")),
		// add access key and secret access key to the to the credential
		Credentials: credentials.NewStaticCredentials(os.Getenv("Access_key_ID"),os.Getenv("Secret_access_key"),""),
	})
	currPath,  err :=os.Getwd()
	//check about any error
	if err != nil{
		// throw the error to error function
		exitErrorf("Error Accrue !",err)
	}

	// open file test
	file, err := os.Open(currPath+"\\test")
	// close the file
	defer file.Close()
	//check if the any error accrue
	if err != nil {
		// throw the error the error function
		exitErrorf("Unable to open file %q, %v", err)
	}

	if len(filename)!=0 {
		// open the file by the file name
		file, err := os.Open(currPath+"\\"+filename)
		// close the file
		defer file.Close()
		filename = strings.Replace(filename,`upload\`,"",-1)
		fmt.Println(filename)
		path = path + "/" + filename
		//check if the any error accrue
		if err != nil {
			// throw the error the error function
			exitErrorf("Unable to open file %q, %v", err)
		}


	}

	//initial the S3 upload Manager.
	uploader := s3manager.NewUploader(sess)
	// upload the file
	_, err   = uploader.Upload( &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(path),
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}

// DownloadFile function to download file from a bucket
func DownloadFile(filename string){
	// intialize bucket name
	bucket :="symbyo-test"
	// create new session on aws aws
	sess, err := session.NewSession(&aws.Config{
		// add region
		Region:      aws.String(os.Getenv("Reg")),
		// add access key and secret access key to the to the credential
		Credentials: credentials.NewStaticCredentials(os.Getenv("Access_key_ID"),os.Getenv("Secret_access_key"),""),
	})

	str := strings.Split(filename,"/")
	currPath,  err :=os.Getwd()

	// open the file by the file name
	file,err :=  os.Create(currPath+"\\download\\"+str[len(str)-1])
	// close the file
	defer file.Close()

	//check if the any error accrue
	if err != nil {
		// throw the error the error function
		exitErrorf("Unable to open file %q, %v", err)
	}

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", filename, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

// DeleteItem function to delete item in bucket
func DeleteItem(obj string){
	// intialize bucket name
	bucket :="symbyo-test"
	// create new session on aws aws
	sess, err := session.NewSession(&aws.Config{
		// add region
		Region:      aws.String(os.Getenv("Reg")),
		// add access key and secret access key to the to the credential
		Credentials: credentials.NewStaticCredentials(os.Getenv("Access_key_ID"),os.Getenv("Secret_access_key"),""),
	})
	//check about any error
	if err != nil{
		// throw the error to error function
		exitErrorf("Error Accrue !",err)
	}

	// specific configuration.
	svc := s3.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(obj)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", obj, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(obj),
	})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", obj, bucket, err)
	}
	fmt.Printf("Object %q successfully deleted\n", obj)
}

// ListItems function to list all items in bucket
func ListItems(id string) []Data {
	// intialize bucket name
	bucket :="symbyo-test"
	// create new session on aws aws
	sess, err := session.NewSession(&aws.Config{
		// add region
		Region:     aws.String(os.Getenv("Reg")),
		// add access key and secret access key to the to the credential
		Credentials: credentials.NewStaticCredentials(os.Getenv("Access_key_ID"),os.Getenv("Secret_access_key"),""),
	})
	//check about any error
	if err != nil{
		// throw the error to error function
		exitErrorf("Error Accrue !",err)
	}

	// specific configuration.
	svc := s3.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:              aws.String(bucket),
		//Delimiter:           aws.String("order/"+id+"/"),
	})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}
	//data := [...]string{}

	var arrdata []Data
	for _, item := range resp.Contents {
		if strings.Contains(*item.Key,"orders/"+id){
			arrdata = append(arrdata,Data{
				Path:     *item.Key ,
				Size: int(*item.Size),
			})
		}
	}
	return arrdata
}

// helps functions
// function to error
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

//function to create a random string
func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}
