package controllers

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofiber/fiber/v2"
	"github.com/lab/tests/app/services"
)

func SingleUploadImageHandler(c *fiber.Ctx) error {
	res := fiber.Map{
		"data":    fiber.Map{},
		"message": "",
	}

	menu := c.FormValue("menu")
	if menu == "" {
		res["message"] = "Menu IS NULL"
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	// รับไฟล์ภาพจากคำขอ
	file, err := c.FormFile("file")
	if err != nil {
		res["message"] = "File IS NULL"
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	// กำหนดค่าบัญชี Azure Storage
	accountName := "storagemasters"
	accountKey := "DY3+F5+OP8HTH6QloogHjKPnxV8Vo/bD8t+ZEPXgfTSHLp8bt8L/ZBsJsHrPp5r2FBzw1o2qoKRr+AStM4Z/1g=="

	// กำหนดค่าการรับรองความถูกต้อง
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	// เปิดไฟล์ภาพจากคำขอ
	srcFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	defer srcFile.Close()

	// ใช้ชื่อไฟล์เริ่มต้นหรือเลือกชื่อที่ต้องการ
	blobName := file.Filename

	file_type, containerName, err := services.DetectContentTypeFromHeader(file, menu)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	// อัปโหลดไฟล์จาก srcFile ไปยัง block blob
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName))
	blockBlobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	uploadOptions := azblob.UploadStreamToBlockBlobOptions{
		BufferSize: 8 * 1024 * 1024, // ขนาดของบล็อก (8 MB)
		MaxBuffers: 16,              // จำนวนการอัปโหลดบล็อกแบบพร้อมกัน
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: file_type,
		},
	}
	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), srcFile, blockBlobURL, uploadOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	res["message"] = os.Getenv("MESSAGE_UPLOAD_IMAGE_SUCCESS")

	return c.JSON(res)
}

/*

func GetBlobURL(containerURL azblob.ContainerURL, blobName string) string {
	blobURL := containerURL.NewBlobURL(blobName)
	blobURLString := blobURL.String()
	return blobURLString
}

func ListBlobsInContainer(c *fiber.Ctx) error {
	// กำหนดค่าบัญชี Azure Storage
	accountName := "storagemasters"
	accountKey := "DY3+F5+OP8HTH6QloogHjKPnxV8Vo/bD8t+ZEPXgfTSHLp8bt8L/ZBsJsHrPp5r2FBzw1o2qoKRr+AStM4Z/1g=="
	containerName := "images-project"

	// กำหนดค่าการรับรองความถูกต้อง
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// สร้าง URL สำหรับคอนเทนเนอร์
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	containerURL := azblob.NewContainerURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	// รับรายการของ blobs ในคอนเทนเนอร์
	ctx := context.TODO()
	listBlob, err := containerURL.ListBlobsFlatSegment(ctx, azblob.Marker{}, azblob.ListBlobsSegmentOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// สร้างโครงสร้างข้อมูลสำหรับ blobs พร้อม URL
	type BlobInfo struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	// สร้างรายชื่อของ blobs พร้อม URL
	var blobInfos []BlobInfo
	for _, blobInfo := range listBlob.Segment.BlobItems {
		blobName := blobInfo.Name
		blobURL := GetBlobURL(containerURL, blobName)
		blobInfos = append(blobInfos, BlobInfo{Name: blobName, URL: blobURL})
	}

	// ส่งรายชื่อของ blobs พร้อม URL ในรูปแบบ JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"blobs": blobInfos,
	})
}
*/
