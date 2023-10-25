package controllers

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofiber/fiber/v2"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func SingleUploadImage(c *fiber.Ctx) error {
	// กำหนดค่าบัญชี Azure Storage
	accountName := "storagemasters"
	accountKey := "DY3+F5+OP8HTH6QloogHjKPnxV8Vo/bD8t+ZEPXgfTSHLp8bt8L/ZBsJsHrPp5r2FBzw1o2qoKRr+AStM4Z/1g=="
	containerName := "images-project"

	// กำหนดค่าการรับรองความถูกต้อง
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// รับไฟล์ภาพจากคำขอ
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// เปิดไฟล์ภาพจากคำขอ
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// ใช้ชื่อไฟล์เริ่มต้นหรือเลือกชื่อที่ต้องการ
	blobName := file.Filename

	// อัปโหลดไฟล์จาก srcFile ไปยัง block blob
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName))
	blockBlobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	uploadOptions := azblob.UploadStreamToBlockBlobOptions{
		BufferSize: 8 * 1024 * 1024, // ขนาดของบล็อก (8 MB)
		MaxBuffers: 16,              // จำนวนการอัปโหลดบล็อกแบบพร้อมกัน
	}
	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), srcFile, blockBlobURL, uploadOptions)
	handleError(err)

	return c.JSON(fiber.Map{"message": "อัปโหลดรูปภาพไปยัง Blob Storage สำเร็จ"})
}
