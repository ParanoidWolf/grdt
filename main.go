package main

func main() {
	body := get_json("memes")
	image_list := get_images(body)
    downloadMultipleFiles(image_list)
}
